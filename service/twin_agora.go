package service

import (
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

const (
	USER_ROLE   = 1 //普通用户
	USER_DOCTOR = 2 //专家
	USER_ALL    = 3 //所有人

	RTC_ROOM_STATUS_CALLING = 1
	RTC_ROOM_STATUS_EXECING = 2
	RTC_ROOM_STATUS_END     = 3

	RTC_ROOM_END_STATUS_TIMEOUT = 1
	RTC_ROOM_END_STATUS_QUIT    = 2
	RTC_ROOM_END_STATUS_DENY    = 3
	RTC_ROOM_END_STATUS_CANCEL  = 4
)

type TwinAgora struct {
	Gorm        *gorm.DB
	RTCUserList map[int]RTCUser
	RTCRoomPool map[string]RTCRoom
	CancelFunc  context.CancelFunc
	CancelCtx   context.Context
	CallTimeout int //一次呼叫，超时时间
	ExecTimeout int //一次通话，超时时间
	Separate    string
}

type RTCUser struct {
	Id     int
	RoomId string
	Uptime int
}

type RTCRoom struct {
	Id                string //唯一标识，UUID4
	Channel           string //频道名
	AddTime           int    //添加时间
	Uptime            int    //最后更新时间
	Status            int    //1发起呼叫，2正常通话中，3已结束
	EndStatus         int    //结束的状态：(1)超时，(2)某一方退出,(3)某一方拒绝(4)发起方主动取消呼叫
	CallUid           int    //发起通话的UID
	ReceiveUids       []int  //被呼叫的用户(专家)IDS
	ReceiveUidsAccept []int  //被呼叫的用户(专家)，接收了此次呼叫
	ReceiveUidsDeny   []int  //被呼叫的用户(专家)，拒绝了此次呼叫
	OnlineUids        []int  //当前在线/在房间通话的用户
	Uids              []int  //所有进入过频道的用户
	//Timeout           int   //超时
}

func NewTwinAgora(Gorm *gorm.DB) *TwinAgora {
	twinAgora := new(TwinAgora)
	twinAgora.Gorm = Gorm
	twinAgora.CallTimeout = 8
	twinAgora.ExecTimeout = 10
	twinAgora.Separate = "##"
	twinAgora.RTCRoomPool = make(map[string]RTCRoom)
	twinAgora.RTCUserList = make(map[int]RTCUser)

	twinAgora.CancelCtx, twinAgora.CancelFunc = context.WithCancel(context.Background())

	return twinAgora
}

func (twinAgora *TwinAgora) CheckTimeout() {
	for {
		select {
		case <-twinAgora.CancelCtx.Done():
			goto end
		default:
			if len(twinAgora.RTCRoomPool) > 0 {
				for _, v := range twinAgora.RTCRoomPool {
					if v.Status == RTC_ROOM_STATUS_END {
						twinAgora.MoveAndStore(v)
						continue
					}

					if util.GetNowTimeSecondToInt() > v.AddTime+twinAgora.CallTimeout && v.Status == RTC_ROOM_STATUS_CALLING {
						v.Status = RTC_ROOM_STATUS_END
						v.EndStatus = RTC_ROOM_END_STATUS_TIMEOUT
						twinAgora.MoveAndStore(v)
						continue
					}

					if util.GetNowTimeSecondToInt() > v.Uptime+twinAgora.ExecTimeout && v.Status == RTC_ROOM_STATUS_EXECING {
						v.Status = RTC_ROOM_STATUS_END
						v.EndStatus = RTC_ROOM_END_STATUS_TIMEOUT
						twinAgora.MoveAndStore(v)
						continue
					}
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
end:
	util.MyPrint("twinAgora CheckTimeout finish.")
}

func (twinAgora *TwinAgora) Start() {
	go twinAgora.CheckTimeout()
}

func (twinAgora *TwinAgora) Quit() {
	twinAgora.CancelFunc()
}

func (twinAgora *TwinAgora) MoveAndStore(RTCRoom RTCRoom) {
	twinAgora.StoreHistory(RTCRoom)
	for _, uid := range RTCRoom.Uids {
		delete(twinAgora.RTCUserList, uid)
	}
	delete(twinAgora.RTCRoomPool, RTCRoom.Channel)
}

func (twinAgora *TwinAgora) Heartbeat(heartbeat pb.Heartbeat, conn *util.Conn) {
	util.MyPrint("twinAgora Heartbeat data:", heartbeat)
	for _, room := range twinAgora.RTCRoomPool {
		for _, uid := range room.Uids {
			if uid == int(conn.UserId) {
				room.Uptime = util.GetNowTimeSecondToInt()
				break
			}
		}
	}
}

func (twinAgora *TwinAgora) ConnCloseCallback(connCloseEvent pb.FDCloseEvent, connManager *util.ConnManager) {
	util.MyPrint("TwinAgora ConnCloseCallback :", connCloseEvent)
	hasSearch := 0
	//已结束的会从map中删除，已超时的也会从map中删除
	for _, RTCRoomInfo := range twinAgora.RTCRoomPool {
		for _, uid := range RTCRoomInfo.Uids {
			if uid == int(connCloseEvent.UserId) {
				hasSearch = 1
			}
		}
		if hasSearch == 0 {
			continue
		}

		//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
		RTCRoomInfo.Status = RTC_ROOM_STATUS_END
		RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_QUIT

		util.MyPrint("RTCRoomInfo.Channel ", RTCRoomInfo.Channel, " , RTCRoomInfo.Uids:", RTCRoomInfo.Uids)
		for _, u := range RTCRoomInfo.Uids {
			if u == int(connCloseEvent.UserId) {
				//不要再给自己发了，因为：它已要断开连接了，发也是失败
				continue
			}
			callPeopleReq := pb.CallPeopleReq{}
			callPeopleReq.Uid = connCloseEvent.UserId
			callPeopleReq.Channel = RTCRoomInfo.Channel

			conn, ok := connManager.Pool[int32(u)]
			if ok {
				conn.SendMsgCompressByUid(int32(u), "SC_PeopleLeave", callPeopleReq)
			}
		}

		//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
		RTCRoomInfo.Status = RTC_ROOM_STATUS_END
		RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_QUIT
	}

}

func (twinAgora *TwinAgora) StoreHistory(RTCRoom RTCRoom) {
	ReceiveUidsStr := util.ArrCoverStr(RTCRoom.ReceiveUids, ",")
	ReceiveUidsAcceptStr := util.ArrCoverStr(RTCRoom.ReceiveUidsAccept, ",")
	ReceiveUidsDenyStr := util.ArrCoverStr(RTCRoom.ReceiveUidsDeny, ",")
	UidsStr := util.ArrCoverStr(RTCRoom.Uids, ",")

	str := RTCRoom.Channel + twinAgora.Separate + strconv.Itoa(RTCRoom.AddTime) + twinAgora.Separate + strconv.Itoa(RTCRoom.Status) + twinAgora.Separate + strconv.Itoa(RTCRoom.EndStatus) + twinAgora.Separate + strconv.Itoa(RTCRoom.CallUid) + twinAgora.Separate + ReceiveUidsStr + twinAgora.Separate + ReceiveUidsAcceptStr + UidsStr + ReceiveUidsAcceptStr + ReceiveUidsDenyStr
	util.MyPrint("StoreHistory:", str)

	var myTwinAgoraRoom model.TwinAgoraRoom
	myTwinAgoraRoom.Channel = RTCRoom.Channel
	myTwinAgoraRoom.CallUid = RTCRoom.CallUid
	myTwinAgoraRoom.Status = RTCRoom.Status
	myTwinAgoraRoom.EndStatus = RTCRoom.EndStatus
	myTwinAgoraRoom.ReceiveUids = ReceiveUidsStr
	myTwinAgoraRoom.ReceiveUidsAccept = ReceiveUidsAcceptStr
	myTwinAgoraRoom.ReceiveUidsDeny = ReceiveUidsDenyStr
	myTwinAgoraRoom.Uids = UidsStr

	err := twinAgora.Gorm.Create(&myTwinAgoraRoom).Error
	if err != nil {
		util.MyPrint("StoreHistory to mysql err:", err)
	}
}

//眼镜端发起呼叫
func (twinAgora *TwinAgora) CallPeople(callPeopleReq pb.CallPeopleReq, conn *util.Conn) {
	util.MyPrint("in func CallPeople:")
	callPeopleRes := pb.CallPeopleRes{}

	uid := int(conn.UserId)
	myRTCUser, ok := twinAgora.RTCUserList[uid]
	if ok && myRTCUser.RoomId != "" {
		RTCRoomInfo, ok := twinAgora.RTCRoomPool[callPeopleReq.Channel]
		if !ok {
			callPeopleRes.ErrCode = 511
			callPeopleRes.ErrMsg = "用户存在roomId,但是房间却不存在，可能程序有问题，请稍后重试"
			util.MyPrint(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_CALLING {
			callPeopleRes.ErrCode = 520
			callPeopleRes.ErrMsg = "已经存在一条记录：发起呼叫，请不要重复发起，或 等待超时 或 主动挂断"
			util.MyPrint(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_EXECING {
			callPeopleRes.ErrCode = 521
			callPeopleRes.ErrMsg = "已经存在一条记录：正在与其它人通话中...，不能再发起CALL了"
			util.MyPrint(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_END {
			twinAgora.MoveAndStore(RTCRoomInfo)
		}
	}

	if callPeopleReq.Uid <= 0 {
		callPeopleRes.ErrCode = 501
		callPeopleRes.ErrMsg = "callPeopleReq.Uid <= 0"
		util.MyPrint(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}

	if callPeopleReq.PeopleType != int32(USER_DOCTOR) {
		callPeopleRes.ErrCode = 502
		callPeopleRes.ErrMsg = "callPeopleReq.PeopleType !=  " + strconv.Itoa(USER_DOCTOR)
		util.MyPrint(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}
	var userDoctorList []model.User
	err := twinAgora.Gorm.Where(" role =  ?", USER_DOCTOR).Find(&userDoctorList).Error
	if err != nil {
		callPeopleRes.ErrCode = 503
		callPeopleRes.ErrMsg = "get user(doctor) by db: is empty"
		util.MyPrint(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}

	if callPeopleReq.TargetUid > 0 {
		callPeopleRes.ErrCode = 504
		callPeopleRes.ErrMsg = "暂时不支持 TargetUid > 0 的情况"
		util.MyPrint(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}
	//寻找在线的专家
	var onlineUserDoctorList []model.User
	for _, userConn := range conn.ConnManager.Pool {
		for _, user := range userDoctorList {
			if userConn.UserId == int32(user.Id) && userConn.Status == util.CONN_STATUS_EXECING {
				onlineUserDoctorList = append(onlineUserDoctorList, user)
			}
		}
	}

	if len(onlineUserDoctorList) <= 0 {
		callPeopleRes.ErrCode = 510
		callPeopleRes.ErrMsg = "所有专家，均不在线"
		util.MyPrint(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}
	myRTCRoom := twinAgora.CreateRTCRoom(callPeopleReq)
	//给所有在线的专家发送邀请通知
	var receiveUids []int
	receiveUidsStr := ""
	for _, user := range onlineUserDoctorList {
		receiveUidsStr += strconv.Itoa(user.Id) + ","
		receiveUids = append(receiveUids, user.Id)
		callReply := pb.CallReply{}
		callReply.RoomId = myRTCRoom.Id
		//pushMsgRes.MsgType = callPeopleReq.
		callReply.Content = strconv.Itoa(int(callPeopleReq.Uid)) + " 呼叫 视频连接...请进入频道:" + callPeopleReq.Channel
		conn.SendMsgCompressByUid(int32(user.Id), "SC_CallReply", callReply)
	}
	myRTCRoom.ReceiveUids = receiveUids

	//先给呼叫者返回消息，告知已成功请等待专家响应
	callPeopleRes.ErrCode = 200
	callPeopleRes.ErrMsg = "请求等待专家响应"
	callPeopleRes.ReceiveUid = receiveUidsStr[0 : len(receiveUidsStr)-1]
	util.MyPrint(callPeopleRes.ErrMsg)
	conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)

	return
}

//这里假设验证都成功了，不做二次验证了
func (twinAgora *TwinAgora) CreateRTCRoom(callPeopleReq pb.CallPeopleReq) RTCRoom {
	uids := []int{int(callPeopleReq.Uid)}
	RTCRoomOne := RTCRoom{
		AddTime: util.GetNowTimeSecondToInt(),
		CallUid: int(callPeopleReq.Uid),
		Status:  RTC_ROOM_STATUS_CALLING,
		Uids:    uids,
		Id:      uuid.NewV4().String(),
	}
	myRTCUser, ok := twinAgora.RTCUserList[int(callPeopleReq.Uid)]
	if ok {
		myRTCUser.RoomId = RTCRoomOne.Id
	} else {
		newRTCUser := RTCUser{}
		newRTCUser.RoomId = RTCRoomOne.Id
		newRTCUser.Uptime = util.GetNowTimeSecondToInt()
		twinAgora.RTCUserList[myRTCUser.Id] = newRTCUser
	}

	twinAgora.RTCRoomPool[callPeopleReq.Channel] = RTCRoomOne
	return RTCRoomOne
}

func (twinAgora *TwinAgora) CancelCallPeople(cancelCallPeopleReq pb.CancelCallPeopleReq, conn *util.Conn) error {
	if cancelCallPeopleReq.Uid <= 0 {
		return errors.New("cancelCallPeopleReq.Uid <= 0")
	}

	if cancelCallPeopleReq.RoomId == "" {
		return errors.New("RoomId empty")
	}
	myRTCUser, ok := twinAgora.RTCUserList[int(cancelCallPeopleReq.Uid)]
	if !ok {
		return errors.New("myRTCUser is not in map")
	}
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[cancelCallPeopleReq.RoomId]
	if !ok {
		return errors.New("roomId not in map")
	}
	//取消呼叫，只能由发起者自己取消
	if RTCRoomInfo.CallUid != int(cancelCallPeopleReq.Uid) {
		return errors.New("你不是发起呼叫者，不能取消操作")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New("房间状态错误：只能呼叫中才能取消")
	}
	//给所有专家端用户发送取消的消息
	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(cancelCallPeopleReq.Uid) == uid {
			continue
		}
		conn.SendMsgCompressByUid(cancelCallPeopleReq.Uid, "SC_CancelCallPeople", cancelCallPeopleReq)
	}
	RTCRoomInfo.Status = RTC_ROOM_STATUS_END
	RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_CANCEL
	return nil
}

func (twinAgora *TwinAgora) PeopleLeave(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
	myRTCUser, ok := twinAgora.RTCUserList[int(callPeopleReq.Uid)]
	if !ok {

	}
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callPeopleReq.RoomId]
	if !ok {
		return errors.New("channel 未找到rtcRoom-1")
	}

	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callPeopleReq.Channel]
	if !ok {
		return errors.New("channel 未找到rtcRoom-3")
	}
	//根据用户发出来的channel，判断该用户是否在此房间中
	//if RTCRoomInfo.CallUid != int(callPeopleReq.Uid) {
	hasSearch := 0
	for _, uid := range RTCRoomInfo.Uids {
		if int(callPeopleReq.Uid) == uid {
			hasSearch = 1
			break
		}
	}

	if hasSearch == 0 {
		//您并不在此频道中，请不要乱发消息
		return errors.New("您并不在此频道中，请不要乱发消息")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_EXECING {
		return errors.New("房间状态错误：只能运行中，才接收离开消息")
	}

	for _, uid := range RTCRoomInfo.Uids {
		if int(callPeopleReq.Uid) == uid {
			continue
		}

		//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
		RTCRoomInfo.Status = RTC_ROOM_STATUS_END
		RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_QUIT

		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_PeopleLeave", callPeopleReq)
	}

	return nil
}

func (twinAgora *TwinAgora) PeopleVote(callVote pb.CallVote) (error, RTCRoom) {
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callVote.RoomId]
	if !ok {
		return errors.New("channel 未找到rtcRoom-2"), RTCRoomInfo
	}
	hasSearch := 0
	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(callVote.Uid) == uid {
			hasSearch = 1
			break
		}
	}
	//并没有发消息让你判定是否接收
	if hasSearch == 0 {
		return errors.New("未发消息给你，请不要捣乱"), RTCRoomInfo
	}

	hasVote := 0
	for _, uid := range RTCRoomInfo.ReceiveUidsAccept {
		if int(callVote.Uid) == uid {
			hasVote = 1
			break
		}
	}
	if hasVote == 1 {
		//您已经投过票了，不要重复操作
		return errors.New("您已经投过票了1，请不要捣乱"), RTCRoomInfo
	}

	hasVote = 0
	for _, uid := range RTCRoomInfo.ReceiveUidsDeny {
		if int(callVote.Uid) == uid {
			hasVote = 1
			break
		}
	}
	if hasVote == 1 {
		//您已经投过票了，不要重复操作
		return errors.New("您已经投过票了2，请不要捣乱"), RTCRoomInfo
	}

	return nil, RTCRoomInfo
}

func (twinAgora *TwinAgora) CallPeopleAccept(callVote pb.CallVote, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New("房间状态错误：只能呼叫中，才接接收")
	}

	RTCRoomInfo.ReceiveUidsAccept = append(RTCRoomInfo.ReceiveUidsAccept, int(callVote.Uid))
	RTCRoomInfo.Uids = append(RTCRoomInfo.Uids, int(callVote.Uid))

	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleAccept", callVote)
	//目前是1v1视频，只要有一个接收，即把房间状态标识为运行中，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_EXECING
	return nil

}

func (twinAgora *TwinAgora) CallPeopleDeny(callVote pb.CallVote, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New("房间状态错误：只能呼叫中，才接拒绝")
	}

	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleDeny", callVote)
	//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_END
	RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_DENY
	return err

}
