package service

import (
	"context"
	"errors"
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
	RTCRoomPool map[string]RTCRoom
	CancelFunc  context.CancelFunc
	CancelCtx   context.Context
	Timeout     int
	Separate    string
}

type RTCRoom struct {
	Channel           string
	AddTime           int
	Status            int   //1发起呼叫，2正常通话中，3已结束
	EndStatus         int   //结束的状态：(1)超时，(2)某一方退出,(3)某一方拒绝(4)发起方主动取消呼叫
	CallUid           int   //发起通话的UID
	ReceiveUids       []int //被呼叫的用户IDS
	ReceiveUidsAccept []int //被呼叫的用户IDS，接收了此次呼叫
	ReceiveUidsDeny   []int //被呼叫的用户IDS，拒绝了此次呼叫
	Uids              []int //ReceiveUidsAccept+CallUid
	//Timeout           int   //超时
}

func NewTwinAgora(Gorm *gorm.DB) *TwinAgora {
	twinAgora := new(TwinAgora)
	twinAgora.Gorm = Gorm
	twinAgora.Timeout = 60
	twinAgora.Separate = "##"
	twinAgora.RTCRoomPool = make(map[string]RTCRoom)
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

					if util.GetNowTimeSecondToInt() > v.AddTime+twinAgora.Timeout && v.Status == RTC_ROOM_STATUS_CALLING {
						v.Status = RTC_ROOM_STATUS_END
						v.EndStatus = RTC_ROOM_END_STATUS_TIMEOUT
						twinAgora.MoveAndStore(v)
						continue
					}
				}
			}
			time.Sleep(time.Second * 1)
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
	delete(twinAgora.RTCRoomPool, RTCRoom.Channel)
}

func (twinAgora *TwinAgora) ConnCloseCallback(conn *util.Conn, source int) {
	hasSearch := 0
	//已结束的会从map中删除，已超时的也会从map中删除
	for _, RTCRoomInfo := range twinAgora.RTCRoomPool {
		for _, uid := range RTCRoomInfo.Uids {
			if uid == int(conn.UserId) {
				hasSearch = 1
			}
		}
		if hasSearch == 0 {
			continue
		}

		//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
		RTCRoomInfo.Status = RTC_ROOM_STATUS_END
		RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_QUIT

		for _, u := range RTCRoomInfo.Uids {
			if u == int(conn.UserId) {
				//不要再给自己发了，因为：它已要断开连接了，发也是失败
				continue
			}
			callPeopleReq := pb.CallPeopleReq{}
			callPeopleReq.Uid = conn.UserId
			callPeopleReq.Channel = RTCRoomInfo.Channel
			conn.SendMsgCompressByUid(int32(u), "SC_PeopleLeave", callPeopleReq)
		}
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
func (twinAgora *TwinAgora) CallPeople(callPeopleReq pb.CallPeopleReq, conn *util.Conn) {
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callPeopleReq.Channel]
	callPeopleRes := pb.CallPeopleRes{}
	if ok {
		if RTCRoomInfo.Status == RTC_ROOM_STATUS_CALLING {
			callPeopleRes.ErrCode = 520
			callPeopleRes.ErrMsg = "已经存在一条记录：发起呼叫，请不要重复发起，或等待超时"
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_EXECING {
			callPeopleRes.ErrCode = 521
			callPeopleRes.ErrMsg = "已经存在一条记录：正常通话中...，不能再发起CALL了"
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_EXECING {
			twinAgora.MoveAndStore(RTCRoomInfo)
		}
	}

	if callPeopleReq.Uid <= 0 {
		callPeopleRes.ErrCode = 501
		callPeopleRes.ErrMsg = "callPeopleReq.Uid <= 0"
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
		return
	}

	if callPeopleReq.PeopleType != int32(USER_DOCTOR) {
		callPeopleRes.ErrCode = 502
		callPeopleRes.ErrMsg = "callPeopleReq.PeopleType != 1 "
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
		return
	}
	var userDoctorList []model.User
	err := twinAgora.Gorm.Where(" role =  ?", USER_DOCTOR).Find(&userDoctorList).Error
	if err != nil {
		callPeopleRes.ErrCode = 503
		callPeopleRes.ErrMsg = "get user by db is empty"
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
		return
	}

	if callPeopleReq.TargetUid > 0 {
		callPeopleRes.ErrCode = 504
		callPeopleRes.ErrMsg = "暂时不支持 TargetUid > 0 的情况"
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
		return
	}

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
		callPeopleRes.ErrMsg = "onlineUserDoctorList is empty"
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleReq)
		return
	}

	var receiveUids []int
	for _, user := range onlineUserDoctorList {
		receiveUids = append(receiveUids, user.Id)
		pushMsgRes := pb.PushMsgRes{}
		pushMsgRes.MsgType = 1
		pushMsgRes.Content = strconv.Itoa(int(callPeopleReq.Uid)) + " 呼叫 视频连接...请进入频道:" + callPeopleReq.Channel
		conn.SendMsgCompressByUid(int32(user.Id), "SC_PushMsg", pushMsgRes)
	}
	twinAgora.CreateRTCRoom(callPeopleReq, receiveUids)
	return
}

//这里假设验证都成功了，不做二次验证了
func (twinAgora *TwinAgora) CreateRTCRoom(callPeopleReq pb.CallPeopleReq, receiveUserIds []int) {
	RTCRoomOne := RTCRoom{
		AddTime: util.GetNowTimeSecondToInt(),
		//Timeout:     util.GetNowTimeSecondToInt() + 60,
		CallUid:     int(callPeopleReq.Uid),
		Status:      RTC_ROOM_STATUS_CALLING,
		ReceiveUids: receiveUserIds,
		Uids:        []int{int(callPeopleReq.Uid)},
	}
	twinAgora.RTCRoomPool[callPeopleReq.Channel] = RTCRoomOne
}

func (twinAgora *TwinAgora) CancelCallPeople(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callPeopleReq.Channel]
	if !ok {
		return errors.New("channel 未找到rtcRoom-1")
	}
	//取消呼叫，只能由发起者自己取消
	if RTCRoomInfo.CallUid != int(callPeopleReq.Uid) {
		return errors.New("你不是发起呼叫者，不能取消操作")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New("房间状态错误：只能呼叫中才能取消")
	}

	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(callPeopleReq.Uid) == uid {
			continue
		}
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CancelCallPeople", callPeopleReq)
	}
	RTCRoomInfo.Status = RTC_ROOM_STATUS_END
	RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_CANCEL
	return nil
}

func (twinAgora *TwinAgora) PeopleLeave(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
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

func (twinAgora *TwinAgora) PeopleVote(callPeopleReq pb.CallPeopleReq) (error, RTCRoom) {
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callPeopleReq.Channel]
	if !ok {
		return errors.New("channel 未找到rtcRoom-2"), RTCRoomInfo
	}
	hasSearch := 0
	for _, uid := range RTCRoomInfo.ReceiveUids {
		if int(callPeopleReq.Uid) == uid {
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
		if int(callPeopleReq.Uid) == uid {
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
		if int(callPeopleReq.Uid) == uid {
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

func (twinAgora *TwinAgora) CallPeopleAccept(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callPeopleReq)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New("房间状态错误：只能呼叫中，才接接收")
	}

	RTCRoomInfo.ReceiveUidsAccept = append(RTCRoomInfo.ReceiveUidsAccept, int(callPeopleReq.Uid))
	RTCRoomInfo.Uids = append(RTCRoomInfo.Uids, int(callPeopleReq.Uid))

	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleAccept", callPeopleReq)
	//目前是1v1视频，只要有一个接收，即把房间状态标识为运行中，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_EXECING
	return nil

}

func (twinAgora *TwinAgora) CallPeopleDeny(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callPeopleReq)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return errors.New("房间状态错误：只能呼叫中，才接拒绝")
	}

	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleDeny", callPeopleReq)
	//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_END
	RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_DENY
	return err

}