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
	USER_ROLE_NORMAL = 1 //普通用户
	USER_ROLE_DOCTOR = 2 //专家

	CALL_USER_PEOPLE_ALL     = 1 //呼叫所有人
	CALL_USER_PEOPLE_GROUP   = 2 //按照<小组>呼叫
	CALL_USER_PEOPLE_PROVIDE = 3 //用户自己指定呼叫的人

	RTC_ROOM_STATUS_CALLING = 1 //房间状态：呼叫中
	RTC_ROOM_STATUS_EXECING = 2 //房间状态：运行中
	RTC_ROOM_STATUS_END     = 3 //房间状态：已结束

	RTC_ROOM_END_STATUS_TIMEOUT_CALLING = 10 //房间结束状态标识：超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_TIMEOUT_EXEC    = 11 //房间结束状态标识：超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_QUIT            = 2  //房间结束状态标识：用户退出
	RTC_ROOM_END_STATUS_DENY            = 3  //房间结束状态标识：被呼叫人拒绝
	RTC_ROOM_END_STATUS_CANCEL          = 4  //房间结束状态标识：呼叫者取消
)

type TwinAgora struct {
	Gorm        *gorm.DB
	RTCUserPool map[int]*RTCUser
	RTCRoomPool map[string]*RTCRoom
	CancelFunc  context.CancelFunc
	CancelCtx   context.Context
	CallTimeout int //一次呼叫，超时时间
	ExecTimeout int //一次通话，超时时间
	Separate    string
}

//创建连接的FD管理池：用户基础信息
type RTCUser struct {
	Id      int    `json:"id"`
	RoomId  string `json:"room_id"`
	Uptime  int    `json:"uptime"`
	AddTime int    `json:"add_time"`
}

type RTCRoom struct {
	Id                string `json:"id"`                  //唯一标识，UUID4
	Channel           string `json:"channel"`             //频道名
	AddTime           int    `json:"add_time"`            //添加时间
	Uptime            int    `json:"uptime"`              //最后更新时间
	Status            int    `json:"status"`              //1发起呼叫，2正常通话中，3已结束
	EndStatus         int    `json:"end_status"`          //结束的状态：(1)超时，(2)某一方退出,(3)某一方拒绝(4)发起方主动取消呼叫
	CallUid           int    `json:"call_uid"`            //发起通话的UID
	ReceiveUids       []int  `json:"receive_uids"`        //被呼叫的用户(专家)IDS
	ReceiveUidsAccept []int  `json:"receive_uids_accept"` //被呼叫的用户(专家)，接收了此次呼叫
	ReceiveUidsDeny   []int  `json:"receive_uids_deny"`   //被呼叫的用户(专家)，拒绝了此次呼叫
	OnlineUids        []int  `json:"online_uids"`         //当前在线并且在房间通话的用户
	Uids              []int  `json:"uids"`                //CallUid + ReceiveUids ,只是记录下，方便函数调用
}

func NewTwinAgora(Gorm *gorm.DB) *TwinAgora {
	twinAgora := new(TwinAgora)
	twinAgora.Gorm = Gorm
	twinAgora.CallTimeout = 8  //呼叫过程的超时时间
	twinAgora.ExecTimeout = 10 //房间运行的超时时间
	twinAgora.Separate = "##"
	twinAgora.RTCRoomPool = make(map[string]*RTCRoom) //房间池
	twinAgora.RTCUserPool = make(map[int]*RTCUser)    //用户池

	twinAgora.CancelCtx, twinAgora.CancelFunc = context.WithCancel(context.Background())

	return twinAgora
}

//开启RTC房间监控.这里有2个主要的功能：
//1. 检查房间各种超时
//2. 心跳更新房间的：最后操作时间(目前是被动接收)
func (twinAgora *TwinAgora) Start() {
	go twinAgora.CheckTimeout()
}

//退出，做善后处理
func (twinAgora *TwinAgora) Quit() {
	twinAgora.CancelFunc() //发送关闭信息
}

//检查房间超时：呼叫超时、运行超时(连接断开)
func (twinAgora *TwinAgora) CheckTimeout() {
	for {
		select {
		case <-twinAgora.CancelCtx.Done():
			goto end
		default:
			if len(twinAgora.RTCRoomPool) <= 0 {
				time.Sleep(time.Millisecond * 100)
				break
			}
			//循环遍历每个房间
			for _, v := range twinAgora.RTCRoomPool {
				//呼叫过程超时
				if util.GetNowTimeSecondToInt() > v.AddTime+twinAgora.CallTimeout && v.Status == RTC_ROOM_STATUS_CALLING {
					twinAgora.RoomEnd(v.Id, RTC_ROOM_END_STATUS_TIMEOUT_CALLING, 10)
					twinAgora.MoveAndStore(v)
					continue
				}
				//房间运行中超时()
				if util.GetNowTimeSecondToInt() > v.Uptime+twinAgora.ExecTimeout && v.Status == RTC_ROOM_STATUS_EXECING {
					twinAgora.RoomEnd(v.Id, RTC_ROOM_END_STATUS_TIMEOUT_EXEC, 11)
					twinAgora.MoveAndStore(v)
					continue
				}
			}

		}
	}
end:
	util.MyPrint("twinAgora CheckTimeout finish.")
}

func (twinAgora *TwinAgora) RoomEnd(roomId string, endStatus int, source int) {
	roomInfo, ok := twinAgora.RTCRoomPool[roomId]
	if !ok {
		errors.New("roomId not in map")
		return
	}
	roomInfo.Status = RTC_ROOM_STATUS_END
	roomInfo.EndStatus = endStatus

	util.MyPrint("RoomEnd , roomId: ", roomId, " , endStatus:", endStatus, " , source:", source)
}

//已结束的房间要做:
//1. 持久化
//2. 房间池内删除该元素
//3. 更新用户池内：在线用户的房间ID清除
func (twinAgora *TwinAgora) MoveAndStore(RTCRoom *RTCRoom) error {
	err := twinAgora.StoreHistory(RTCRoom)
	if err != nil {
		return err
	}

	for _, uid := range RTCRoom.OnlineUids {
		myRTCUser, ok := twinAgora.GetUserById(uid)
		if ok {
			myRTCUser.RoomId = ""
		}
	}
	delete(twinAgora.RTCRoomPool, RTCRoom.Channel)
	return nil
}

func (twinAgora *TwinAgora) DelUserById(uid int) {
	util.MyPrint("DelUserById:", uid)
	delete(twinAgora.RTCUserPool, uid)
}

func (twinAgora *TwinAgora) GetUserById(uid int) (mmRTCUserRTCUser *RTCUser, rs bool) {
	if uid <= 0 {
		util.MyPrint("GetUserById err : uid <= 0 ")
		return mmRTCUserRTCUser, false
	}
	//util.MyPrint("GetUserById uid:", uid, twinAgora.RTCUserPool)
	myRTCUser, ok := twinAgora.RTCUserPool[uid]
	util.MyPrint(myRTCUser, ok)
	if !ok {
		util.MyPrint("GetUserById empty:", uid, myRTCUser, ok)
	}
	return myRTCUser, ok
}

//心跳，更新房间的最后更新时间
func (twinAgora *TwinAgora) Heartbeat(heartbeat pb.Heartbeat, conn *util.Conn) {
	util.MyPrint("twinAgora Heartbeat data:", heartbeat)
	myRTCUser, ok := twinAgora.GetUserById(int(conn.UserId))
	if !ok {
		//这里也是个异常，按说连接创建，网关会推送消息，会立刻创建一个记录
		myRTCUser = &RTCUser{
			Id:      int(conn.UserId),
			AddTime: util.GetNowTimeSecondToInt(),
			RoomId:  "",
			Uptime:  util.GetNowTimeSecondToInt(),
		}
		return
	}

	if myRTCUser.RoomId == "" {
		return
	}

	myRTCRoom, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
	if !ok { //这是种异常的情况，用户基础信息里roomId存在 ,但是在池里已经不存在了，可能是其它协程已经操作了，但是没有清空RTCUser的ROOMID
		myRTCUser.RoomId = ""
		return
	}
	//这里是个异常，按说房间已经结束，用户基础信息应该把roomId清掉
	if myRTCRoom.Status == RTC_ROOM_STATUS_END {
		myRTCUser.RoomId = ""
		//这里只做用户的roomId的清理，房间的回收交给专门的后台守护协程来做
		return
	}
	//只能最后一种房间状态了：RTC_ROOM_STATUS_EXEC
	myRTCUser.Uptime = util.GetNowTimeSecondToInt()
	myRTCRoom.Uptime = util.GetNowTimeSecondToInt()
}
func (twinAgora *TwinAgora) FDCreateEvent(FDCreateEvent pb.FDCreateEvent, conn *util.Conn) {
	myRTCUser, ok := twinAgora.GetUserById(int(FDCreateEvent.UserId))
	if ok {
		util.MyPrint("错误：重复创建RTCUser....", myRTCUser)
		return
	}
	util.MyPrint("FDCreateEvent ,uid:", FDCreateEvent.UserId)
	NewRTCUser := RTCUser{
		Id:      int(FDCreateEvent.UserId),
		AddTime: util.GetNowTimeSecondToInt(),
		Uptime:  util.GetNowTimeSecondToInt(),
		RoomId:  "",
	}
	twinAgora.RTCUserPool[int(FDCreateEvent.UserId)] = &NewRTCUser
}

//对端连接断开(超时)，网关推送事件
func (twinAgora *TwinAgora) ConnCloseCallback(connCloseEvent pb.FDCloseEvent, connManager *util.ConnManager) {
	util.MyPrint("TwinAgora ConnCloseCallback :", connCloseEvent)
	myRTCUser, ok := twinAgora.GetUserById(int(connCloseEvent.UserId))
	if !ok {
		return
	}

	if myRTCUser.RoomId == "" {
		twinAgora.DelUserById(int(connCloseEvent.UserId))
		return
	}
	myRTCRoom, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
	if !ok { //这是种异常的情况，用户基础信息里roomId存在 ,但是在池里已经不存在了，可能是其它协程已经操作了，但是没有清空RTCUser的ROOMID
		twinAgora.DelUserById(int(connCloseEvent.UserId))
		return
	}

	if myRTCRoom.Status == RTC_ROOM_STATUS_END {
		//这也是异常情况，池子里虽然有个房间，但是状态是已经结束了，可能后台协程也没有来得及处理
		twinAgora.DelUserById(int(connCloseEvent.UserId))
		return
	}
	util.MyPrint("RTCRoomInfo.Channel ", myRTCRoom.Channel, " , RTCRoomInfo.OnlineUids:", myRTCRoom.OnlineUids)
	for _, u := range myRTCRoom.OnlineUids {
		if u == int(connCloseEvent.UserId) {
			//不要再给自己发了，因为：它已要断开连接了，发也是失败
			continue
		}
		peopleLeaveRes := pb.PeopleLeaveRes{}
		peopleLeaveRes.Uid = connCloseEvent.UserId
		peopleLeaveRes.Channel = myRTCRoom.Channel
		peopleLeaveRes.RoomId = myRTCRoom.Id

		conn, ok := connManager.Pool[int32(u)]
		if ok {
			conn.SendMsgCompressByUid(int32(u), "SC_PeopleLeave", peopleLeaveRes)
		}
	}

	//目前是1v1视频，只要有一个人拒绝|断线，即结束，这里后期优化一下吧
	twinAgora.RoomEnd(myRTCRoom.Id, RTC_ROOM_END_STATUS_QUIT, 21)
	twinAgora.MoveAndStore(myRTCRoom) //注意：它的内部会把RTCUser roomId清空，但不删除
	twinAgora.DelUserById(int(connCloseEvent.UserId))

}

//持久化到DB中
func (twinAgora *TwinAgora) StoreHistory(RTCRoom *RTCRoom) error {
	var twinAgoraRoomRow model.TwinAgoraRoom
	twinAgora.Gorm.Where("room_id = ? ", RTCRoom.Id).First(&twinAgoraRoomRow)
	if twinAgoraRoomRow.Id > 0 {
		return errors.New("db has exist ,do not repeat opt")
	}
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
	return nil
}

//眼镜端发起呼叫
func (twinAgora *TwinAgora) CallPeople(callPeopleReq pb.CallPeopleReq, conn *util.Conn) {
	util.MyPrint("in func CallPeople:")
	callPeopleRes := pb.CallPeopleRes{}

	if callPeopleReq.Uid <= 0 {
		callPeopleRes.ErrCode = 501
		callPeopleRes.ErrMsg = "callPeopleReq.Uid <= 0"
		util.MyPrint(callPeopleRes.ErrMsg)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}

	if callPeopleReq.PeopleType != int32(USER_ROLE_DOCTOR) {
		callPeopleRes.ErrCode = 502
		callPeopleRes.ErrMsg = "callPeopleReq.PeopleType !=  " + strconv.Itoa(USER_ROLE_DOCTOR)
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}

	if callPeopleReq.TargetUid > 0 {
		callPeopleRes.ErrCode = 504
		callPeopleRes.ErrMsg = "暂时不支持 TargetUid > 0 的情况"
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}

	uid := int(conn.UserId)
	myRTCUser, ok := twinAgora.GetUserById(uid)
	if ok && myRTCUser.RoomId != "" {
		RTCRoomInfo, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
		if !ok {
			callPeopleRes.ErrCode = 511
			callPeopleRes.ErrMsg = "用户存在roomId,但是房间却不存在，可能程序有问题，请稍后重试"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_CALLING {
			callPeopleRes.ErrCode = 520
			callPeopleRes.ErrMsg = "已经存在一条记录：发起呼叫，请不要重复发起，或 等待超时 或 主动挂断"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}

		if RTCRoomInfo.Status == RTC_ROOM_STATUS_EXECING {
			callPeopleRes.ErrCode = 521
			callPeopleRes.ErrMsg = "已经存在一条记录：正在与其它人通话中...，不能再发起CALL了"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}
		//该房间状态已经结束了，但未做清算处理(持久化)，这里做个容错吧
		if RTCRoomInfo.Status == RTC_ROOM_STATUS_END {
			callPeopleRes.ErrCode = 522
			callPeopleRes.ErrMsg = "房间结束了，但未清楚，请等待后台进程清算"
			twinAgora.MakeError(callPeopleRes.ErrMsg)
			conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
			return
		}
	}

	var userDoctorList []model.User
	err := twinAgora.Gorm.Where(" role =  ?", USER_ROLE_DOCTOR).Find(&userDoctorList).Error
	if err != nil {
		callPeopleRes.ErrCode = 503
		callPeopleRes.ErrMsg = "get user(doctor) by db: is empty"
		twinAgora.MakeError(callPeopleRes.ErrMsg)
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
		twinAgora.MakeError(callPeopleRes.ErrMsg)
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)
		return
	}
	myRTCRoom := twinAgora.CreateRTCRoom(callPeopleReq, onlineUserDoctorList)
	receiveUidsStr := ""
	for _, user := range onlineUserDoctorList {
		receiveUidsStr += strconv.Itoa(user.Id) + "," //专家接收列表
		callReply := pb.CallReply{}
		callReply.RoomId = myRTCRoom.Id
		callReply.Content = strconv.Itoa(int(callPeopleReq.Uid)) + " 呼叫 视频连接...请进入频道:" + callPeopleReq.Channel
		conn.SendMsgCompressByUid(int32(user.Id), "SC_CallReply", callReply)
	}

	myRTCUser.RoomId = myRTCRoom.Id
	//先给呼叫者返回消息，告知已成功请等待专家响应
	callPeopleRes.RoomId = myRTCRoom.Id
	callPeopleRes.ErrCode = 200
	callPeopleRes.ErrMsg = "请求等待专家响应"
	callPeopleRes.ReceiveUid = receiveUidsStr[0 : len(receiveUidsStr)-1]

	conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_CallPeople", callPeopleRes)

	return
}

//这里假设验证都成功了，不做二次验证了
func (twinAgora *TwinAgora) CreateRTCRoom(callPeopleReq pb.CallPeopleReq, onlineUserDoctorList []model.User) RTCRoom {
	//给所有在线的专家发送邀请通知
	var receiveUids []int
	//receiveUidsStr := ""
	for _, user := range onlineUserDoctorList {
		//receiveUidsStr += strconv.Itoa(user.Id) + "," //专家接收列表
		receiveUids = append(receiveUids, user.Id) //专家接收列表

	}
	//myRTCRoom.ReceiveUids = receiveUids
	//uids := []int{int(callPeopleReq.Uid)}
	RTCRoomOne := RTCRoom{
		AddTime:     util.GetNowTimeSecondToInt(),
		CallUid:     int(callPeopleReq.Uid),
		Status:      RTC_ROOM_STATUS_CALLING,
		Uids:        append(receiveUids, int(callPeopleReq.Uid)),
		Id:          uuid.NewV4().String(),
		ReceiveUids: receiveUids,
	}
	myRTCUser, ok := twinAgora.GetUserById(int(callPeopleReq.Uid))
	if ok {
		myRTCUser.RoomId = RTCRoomOne.Id
	} else {
		newRTCUser := RTCUser{}
		newRTCUser.RoomId = RTCRoomOne.Id
		newRTCUser.Uptime = util.GetNowTimeSecondToInt()
		twinAgora.RTCUserPool[myRTCUser.Id] = &newRTCUser
	}

	twinAgora.RTCRoomPool[RTCRoomOne.Id] = &RTCRoomOne
	return RTCRoomOne
}

//发起方，取消呼叫
func (twinAgora *TwinAgora) CancelCallPeople(cancelCallPeopleReq pb.CancelCallPeopleReq, conn *util.Conn) error {
	if cancelCallPeopleReq.Uid <= 0 {
		return errors.New("cancelCallPeopleReq.Uid <= 0")
	}

	if cancelCallPeopleReq.RoomId == "" {
		return errors.New("RoomId empty")
	}

	myRTCUser, ok := twinAgora.GetUserById(int(cancelCallPeopleReq.Uid))
	if !ok {
		return errors.New("myRTCUser is not in map")
	}

	if myRTCUser.RoomId != cancelCallPeopleReq.RoomId {
		return errors.New("myRTCUser.RoomId != cancelCallPeopleReq.RoomId")
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
	twinAgora.MoveAndStore(RTCRoomInfo)
	return nil
}

func (twinAgora *TwinAgora) MakeError(errMsg string) error {
	util.MyPrint("*********=====MakeError : ", errMsg)
	return errors.New(errMsg)
}

func (twinAgora *TwinAgora) PeopleEntry(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
	util.MyPrint("in func PeopleEntry:")
	myRTCUser, ok := twinAgora.GetUserById(int(callPeopleReq.Uid))
	if !ok {
		return twinAgora.MakeError("GetUserById empty ")
	}
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
	if !ok {
		return twinAgora.MakeError("roomId NOT in map")
	}

	hasSearch := 0
	for _, uid := range RTCRoomInfo.OnlineUids {
		if int(callPeopleReq.Uid) == uid {
			hasSearch = 1
			break
		}
	}

	if hasSearch == 1 {
		//您并不在此频道中，请不要乱发消息
		return errors.New("您已经此频道中，请不要乱发消息")
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_EXECING {
		return errors.New("房间状态错误：只能运行中，才接收进入房间消息")
	}

	for _, uid := range RTCRoomInfo.OnlineUids {
		if int(callPeopleReq.Uid) == uid {
			continue
		}
		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_PeopleEntry", callPeopleReq)
	}
	RTCRoomInfo.OnlineUids = append(RTCRoomInfo.OnlineUids, int(callPeopleReq.Uid))

	return nil
}

//某用户离开了房间
func (twinAgora *TwinAgora) PeopleLeave(callPeopleReq pb.CallPeopleReq, conn *util.Conn) error {
	myRTCUser, ok := twinAgora.GetUserById(int(callPeopleReq.Uid))
	if !ok {
		return errors.New("GetUserById empty ")
	}
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
	if !ok {
		return errors.New("roomId NOT in map")
	}

	hasSearch := 0
	for _, uid := range RTCRoomInfo.OnlineUids {
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

	for _, uid := range RTCRoomInfo.OnlineUids {
		if int(callPeopleReq.Uid) == uid {
			continue
		}

		//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
		RTCRoomInfo.Status = RTC_ROOM_STATUS_END
		RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_QUIT
		twinAgora.MoveAndStore(RTCRoomInfo)

		conn.SendMsgCompressByUid(callPeopleReq.Uid, "SC_PeopleLeave", callPeopleReq)
	}

	return nil
}

//被呼叫者，接收/拒绝 公共验证
func (twinAgora *TwinAgora) PeopleVote(callVote pb.CallVote) (error, *RTCRoom) {
	util.MyPrint("PeopleVote :", callVote)
	RTCRoomInfo, ok := twinAgora.RTCRoomPool[callVote.RoomId]
	if !ok {
		return twinAgora.MakeError("room id not in map"), RTCRoomInfo
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
		return twinAgora.MakeError("未发消息给你，请不要捣乱"), RTCRoomInfo
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
		return twinAgora.MakeError("您已经投过票了1，请不要捣乱"), RTCRoomInfo
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
		return twinAgora.MakeError("您已经投过票了2，请不要捣乱"), RTCRoomInfo
	}

	return nil, RTCRoomInfo
}

//被呼叫者，接收呼叫
func (twinAgora *TwinAgora) CallPeopleAccept(callVote pb.CallVote, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return twinAgora.MakeError("房间状态错误：只能呼叫中，才接接收")
	}

	RTCUser, _ := twinAgora.GetUserById(int(conn.UserId))
	RTCUser.RoomId = RTCRoomInfo.Id

	RTCRoomInfo.ReceiveUidsAccept = append(RTCRoomInfo.ReceiveUidsAccept, int(callVote.Uid))
	util.MyPrint("RTCRoomInfo.ReceiveUidsAccept:", RTCRoomInfo.ReceiveUidsAccept)
	//RTCRoomInfo.Uids = append(RTCRoomInfo.Uids, int(callVote.Uid))
	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleAccept", callVote)
	//目前是1v1视频，只要有一个接收，即把房间状态标识为运行中，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_EXECING
	return nil

}

//被呼叫者，拒绝呼叫
func (twinAgora *TwinAgora) CallPeopleDeny(callVote pb.CallVote, conn *util.Conn) error {
	err, RTCRoomInfo := twinAgora.PeopleVote(callVote)
	if err != nil {
		return err
	}

	if RTCRoomInfo.Status != RTC_ROOM_STATUS_CALLING {
		return twinAgora.MakeError("房间状态错误：只能呼叫中，才接拒绝")
	}

	conn.SendMsgCompressByUid(int32(RTCRoomInfo.CallUid), "SC_CallPeopleDeny", callVote)
	//目前是1v1视频，只要有一个人拒绝，即结束，这里后期优化一下吧
	RTCRoomInfo.Status = RTC_ROOM_STATUS_END
	RTCRoomInfo.EndStatus = RTC_ROOM_END_STATUS_DENY
	twinAgora.MoveAndStore(RTCRoomInfo)
	return err

}
