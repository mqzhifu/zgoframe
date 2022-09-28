package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"sync"
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

	RTC_ROOM_END_STATUS_TIMEOUT_CALLING = 10 //房间结束状态标识：呼叫超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_TIMEOUT_EXEC    = 11 //房间结束状态标识：运行超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_QUIT            = 2  //房间结束状态标识：用户退出
	RTC_ROOM_END_STATUS_DENY            = 3  //房间结束状态标识：被呼叫人拒绝
	RTC_ROOM_END_STATUS_CANCEL          = 4  //房间结束状态标识：呼叫者取消
)

type TwinAgora struct {
	Gorm                 *gorm.DB
	RTCUserPool          map[int]*RTCUser
	RTCRoomPool          map[string]*RTCRoom
	CancelFunc           context.CancelFunc
	CancelCtx            context.Context
	CallTimeout          int //一次呼叫，超时时间
	ExecTimeout          int //一次通话，超时时间
	UserHeartbeatTimeout int //一个用户建立的长连接，超时时间
	Separate             string
	ConnManager          *util.ConnManager
}

//创建连接的FD管理池：用户基础信息
type RTCUser struct {
	Id      int    `json:"id"`       //用户ID
	RoomId  string `json:"room_id"`  //用户所有房间ID
	Uptime  int    `json:"uptime"`   //最后更新时间
	AddTime int    `json:"add_time"` //添加时间
}

type RTCRoom struct {
	Id                string       `json:"id"`                  //唯一标识，UUID4
	Channel           string       `json:"channel"`             //频道名
	AddTime           int          `json:"add_time"`            //添加时间
	Uptime            int          `json:"uptime"`              //最后更新时间
	Status            int          `json:"status"`              //1发起呼叫，2正常通话中，3已结束
	EndStatus         int          `json:"end_status"`          //结束的状态：(1)超时，(2)某一方退出,(3)某一方拒绝(4)发起方主动取消呼叫
	CallUid           int          `json:"call_uid"`            //发起通话的UID
	ReceiveUids       []int        `json:"receive_uids"`        //被呼叫的用户(专家)IDS
	ReceiveUidsAccept []int        `json:"receive_uids_accept"` //被呼叫的用户(专家)，接收了此次呼叫
	ReceiveUidsDeny   []int        `json:"receive_uids_deny"`   //被呼叫的用户(专家)，拒绝了此次呼叫
	OnlineUids        []int        `json:"online_uids"`         //当前在线并且在房间通话的用户
	Uids              []int        `json:"uids"`                //CallUid + ReceiveUids ,只是记录下，方便函数调用
	RWLock            sync.RWMutex `json:"-"`                   //变更状态的时候使用
}

func NewTwinAgora(Gorm *gorm.DB) *TwinAgora {
	twinAgora := new(TwinAgora)
	twinAgora.Gorm = Gorm                             //房间数据持久化
	twinAgora.CallTimeout = 8                         //呼叫过程的超时时间
	twinAgora.ExecTimeout = 10                        //房间运行的超时时间，room_heartbeat 也使用此值
	twinAgora.UserHeartbeatTimeout = 3                //一个用户建立的长连接，超时时间
	twinAgora.Separate = "##"                         //一个房间信息转换成字符串的：分隔符
	twinAgora.RTCRoomPool = make(map[string]*RTCRoom) //房间池
	twinAgora.RTCUserPool = make(map[int]*RTCUser)    //用户池

	twinAgora.CancelCtx, twinAgora.CancelFunc = context.WithCancel(context.Background())

	return twinAgora
}

var ERR_ROOM_ID_NOT_IN_MAP = "roomId not in map : "
var ERR_ROOM_STATUS_END_WAIT_DEMON = "RTC_ROOM_STATUS_END , waiting demon coroutines process..."

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

//守护协程，检查房间超时：呼叫超时、运行超时(连接断开)
func (twinAgora *TwinAgora) CheckTimeout() {
	for {
		select {
		case <-twinAgora.CancelCtx.Done():
			goto end
		default:
			//循环遍历每个房间
			for _, room := range twinAgora.RTCRoomPool {
				//呼叫过程超时
				if util.GetNowTimeSecondToInt() > room.AddTime+twinAgora.CallTimeout && room.Status == RTC_ROOM_STATUS_CALLING {
					twinAgora.RoomEnd(room.Id, RTC_ROOM_END_STATUS_TIMEOUT_CALLING, 10)
					continue
				}
				//房间运行中超时()
				if util.GetNowTimeSecondToInt() > room.Uptime+twinAgora.ExecTimeout && room.Status == RTC_ROOM_STATUS_EXECING {
					twinAgora.RoomEnd(room.Id, RTC_ROOM_END_STATUS_TIMEOUT_EXEC, 11)
					continue
				}
			}

			for _, user := range twinAgora.RTCUserPool {
				if util.GetNowTimeSecondToInt() > user.Uptime+twinAgora.UserHeartbeatTimeout {
					twinAgora.ConnCloseProcess(user)
				}
			}
			time.Sleep(time.Millisecond * 100)
			break

		}
	}
end:
	util.MyPrint("twinAgora CheckTimeout finish.")
}

//网关监控到有C端连接，并通过了登陆验证后，会推送事件
func (twinAgora *TwinAgora) FDCreateEvent(FDCreateEvent pb.FDCreateEvent, conn *util.Conn) {
	_, ok := twinAgora.GetUserById(int(FDCreateEvent.UserId))
	if ok {
		msgInfo := "错误：已有存在RTCUser，请不要重复连接....UID:" + strconv.Itoa(int(FDCreateEvent.UserId))
		twinAgora.MakeError(msgInfo)
		twinAgora.PushMsg(conn, int(FDCreateEvent.UserId), 500, 1, msgInfo)
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

//网关监控到用户连接断了(超时)，会回调通知服务
func (twinAgora *TwinAgora) FDCloseEvent(connCloseEvent pb.FDCloseEvent, connManager *util.ConnManager) {
	util.MyPrint("TwinAgora ConnCloseCallback :", connCloseEvent)
	myRTCUser, ok := twinAgora.GetUserById(int(connCloseEvent.UserId))
	if !ok {
		return
	}
	twinAgora.ConnCloseProcess(myRTCUser)
}

//用户长连接 - 心跳，更新房间的最后更新时间
func (twinAgora *TwinAgora) UserHeartbeat(heartbeat pb.Heartbeat, conn *util.Conn) {
	util.MyPrint("twinAgora Heartbeat data:", heartbeat)
	myRTCUser, ok := twinAgora.GetUserById(int(conn.UserId))
	if !ok {
		msgInfo := "错误：未找到RTCUser....UID:" + strconv.Itoa(int(conn.UserId))
		twinAgora.MakeError(msgInfo)
		twinAgora.PushMsg(conn, int(conn.UserId), 500, 1, msgInfo)
		return
	}

	myRTCUser.Uptime = util.GetNowTimeSecondToInt()

	if myRTCUser.RoomId == "" {
		return
	}
	myRTCRoom, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	//myRTCRoom, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
	if err != nil { //这是种异常的情况，用户基础信息里roomId存在 ,但是在池里已经不存在了，可能是其它协程已经操作了，但是没有清空RTCUser的ROOMID
		twinAgora.MakeError(ERR_ROOM_ID_NOT_IN_MAP + myRTCUser.RoomId)
		myRTCUser.RoomId = ""
		return
	}
	//这里是个异常，按说房间已经结束，用户基础信息应该把roomId清掉
	if myRTCRoom.Status == RTC_ROOM_STATUS_END {
		//交给后台守护协程处理，roomId会被清空的
		twinAgora.MakeError(ERR_ROOM_STATUS_END_WAIT_DEMON + myRTCUser.RoomId)
		return
	}

}

//每个房间的心跳，因为音视频使用的是声网，监控不到，就得单独再加一个心跳
func (twinAgora *TwinAgora) RoomHeartbeat(heartbeat pb.RoomHeartbeatReq, conn *util.Conn) {
	myRTCUser, ok := twinAgora.GetUserById(int(heartbeat.Uid))
	if !ok {
		return
	}

	if myRTCUser.RoomId == "" {
		return
	}

	myRTCRoom, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	if err != nil {
		return
	}

	if myRTCRoom.Status != RTC_ROOM_STATUS_EXECING {
		twinAgora.MakeError("RoomHeartbeat myRTCRoom.Status != RTC_ROOM_STATUS_EXECING ")
		return
	}

	myRTCRoom.Uptime = util.GetNowTimeSecondToInt()

}

//给客户端推送消息，主要是一些错误信息
func (twinAgora *TwinAgora) PushMsg(conn *util.Conn, uid int, code int, eventId int, content string) {
	pushMsg := pb.PushMsg{
		Code:    int32(code),
		Uid:     int32(uid),
		EventId: int32(eventId),
		Content: content,
	}
	conn.SendMsgCompressByUid(int32(uid), "SC_PushMsg", pushMsg)
}

//连接断开或超时处理
func (twinAgora *TwinAgora) ConnCloseProcess(rtcUserRTCUser *RTCUser) {
	if rtcUserRTCUser.RoomId == "" {
		twinAgora.DelUserById(rtcUserRTCUser.Id)
		return
	}
	myRTCRoom, err := twinAgora.GetRoomById(rtcUserRTCUser.RoomId)
	if err != nil { //这是种异常的情况，用户基础信息里roomId存在 ,但是在池里已经不存在了，可能是其它协程已经操作了，但是没有清空RTCUser的ROOMID
		twinAgora.MakeError(ERR_ROOM_ID_NOT_IN_MAP + rtcUserRTCUser.RoomId)
		twinAgora.DelUserById(rtcUserRTCUser.Id)
		return
	}

	if myRTCRoom.Status == RTC_ROOM_STATUS_END {
		//这也是异常情况，池子里虽然有个房间，但是状态是已经结束了，可能后台协程也没有来得及处理
		twinAgora.MakeError(ERR_ROOM_STATUS_END_WAIT_DEMON + rtcUserRTCUser.RoomId)
		twinAgora.DelUserById(rtcUserRTCUser.Id)
		return
	}
	util.MyPrint("RTCRoomInfo.Channel ", myRTCRoom.Channel, " , RTCRoomInfo.OnlineUids:", myRTCRoom.OnlineUids)
	//for _, u := range myRTCRoom.OnlineUids {
	//	if u == rtcUserRTCUser.Id {
	//		//不要再给自己发了，因为：它已要断开连接了，发也是失败
	//		continue
	//	}
	//	peopleLeaveRes := pb.PeopleLeaveRes{}
	//	peopleLeaveRes.Uid = int32(rtcUserRTCUser.Id)
	//	peopleLeaveRes.Channel = myRTCRoom.Channel
	//	peopleLeaveRes.RoomId = myRTCRoom.Id
	//
	//	conn, ok := twinAgora.ConnManager.Pool[int32(u)]
	//	if ok {
	//		conn.SendMsgCompressByUid(int32(u), "SC_PeopleLeave", peopleLeaveRes)
	//	}
	//}

	//目前是1v1视频，只要有一个人拒绝|断线，即结束，这里后期优化一下吧
	twinAgora.RoomEnd(myRTCRoom.Id, RTC_ROOM_END_STATUS_QUIT, 21)
	twinAgora.DelUserById(rtcUserRTCUser.Id)
}

//已结束的房间要做:
//1. 持久化
//2. 房间池内删除该元素
//3. 更新用户池内：在线用户的房间ID清除
func (twinAgora *TwinAgora) RoomEnd(roomId string, endStatus int, source int) {
	roomInfo, err := twinAgora.GetRoomById(roomId)
	if err != nil {
		errors.New("roomId not in map")
		return
	}
	//要修改房间状态，要持久化，最终还要删除内存池中的记录，所以要加写锁
	roomInfo.RWLock.Lock()
	defer roomInfo.RWLock.Unlock()

	roomInfo.Status = RTC_ROOM_STATUS_END
	roomInfo.EndStatus = endStatus

	twinAgora.StoreHistory(roomInfo)

	for _, uid := range roomInfo.OnlineUids {
		myRTCUser, ok := twinAgora.GetUserById(uid)
		if ok {
			conn, ok2 := twinAgora.ConnManager.Pool[int32(uid)]
			if ok2 {
				peopleLeaveRes := pb.PeopleLeaveRes{}
				peopleLeaveRes.Uid = int32(uid)
				peopleLeaveRes.Channel = roomInfo.Channel
				peopleLeaveRes.RoomId = roomInfo.Id

				conn.SendMsgCompressByUid(int32(uid), "SC_PeopleLeave", peopleLeaveRes)
			}
			myRTCUser.RoomId = ""
		}
	}
	util.MyPrint("delete room:", roomInfo.Id)
	delete(twinAgora.RTCRoomPool, roomInfo.Id)

	util.MyPrint("RoomEnd , roomId: ", roomId, " , endStatus:", endStatus, " , source:", source)
}

//持久化到DB中
func (twinAgora *TwinAgora) StoreHistory(RTCRoom *RTCRoom) error {
	var twinAgoraRoomRow model.TwinAgoraRoom
	twinAgora.Gorm.Where("room_id = ? ", RTCRoom.Id).First(&twinAgoraRoomRow)
	if twinAgoraRoomRow.Id > 0 {
		return twinAgora.MakeError("db has exist ,do not repeat opt")
	}
	ReceiveUidsStr := util.ArrCoverStr(RTCRoom.ReceiveUids, ",")
	ReceiveUidsAcceptStr := util.ArrCoverStr(RTCRoom.ReceiveUidsAccept, ",")
	ReceiveUidsDenyStr := util.ArrCoverStr(RTCRoom.ReceiveUidsDeny, ",")
	UidsStr := util.ArrCoverStr(RTCRoom.Uids, ",")

	str := RTCRoom.Channel + twinAgora.Separate + strconv.Itoa(RTCRoom.AddTime) + twinAgora.Separate + strconv.Itoa(RTCRoom.Status) + twinAgora.Separate + strconv.Itoa(RTCRoom.EndStatus) + twinAgora.Separate + strconv.Itoa(RTCRoom.CallUid) + twinAgora.Separate + ReceiveUidsStr + twinAgora.Separate + ReceiveUidsAcceptStr + UidsStr + ReceiveUidsAcceptStr + ReceiveUidsDenyStr
	util.MyPrint("StoreHistory:", str)

	var myTwinAgoraRoom model.TwinAgoraRoom
	myTwinAgoraRoom.RoomId = RTCRoom.Id
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
		twinAgora.MakeError("StoreHistory to mysql err:" + err.Error())
	}
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

func (twinAgora *TwinAgora) MakeError(errMsg string) error {
	util.MyPrint("*********=====MakeError : ", errMsg)
	return errors.New(errMsg)
}

func (twinAgora *TwinAgora) GetRoomById(id string) (room *RTCRoom, err error) {
	if id == "" {
		return room, twinAgora.MakeError("roomId is empty")
	}

	room, ok := twinAgora.RTCRoomPool[id]
	if !ok {
		return room, twinAgora.MakeError("id not in room map.")
	}

	return room, nil
}
