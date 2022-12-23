package seed_business

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/util"
)

type TwinAgora struct {
	Op                   TwinAgoraOption
	Gorm                 *gorm.DB
	RTCUserPool          map[int]*RTCUser    //用户池
	RTCRoomPool          map[string]*RTCRoom //房间池
	CallTimeout          int                 //一次呼叫，超时时间
	ExecTimeout          int                 //一次通话，超时时间
	ResAcceptTimeout     int                 //专家端收到  确定  取消  请求后，超时时间
	EntryTimeout         int                 //专家端同意了通话，此时对端迟迟未进入房间
	UserHeartbeatTimeout int                 //一个用户建立的长连接，超时时间
	Separate             string              //一个房间信息转换成字符串的：分隔符
	//RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
	CancelFunc context.CancelFunc
	CancelCtx  context.Context
	//ConnManager           *util.ConnManager
	//Err                  map[int]string
	Log  *zap.Logger
	Lang *util.ErrMsg
}

//创建连接的FD管理池：用户基础信息
type RTCUser struct {
	Id            int    `json:"id"`             //用户ID
	RoomId        string `json:"room_id"`        //用户所有房间ID
	RoomHeartbeat int    `json:"room_heartbeat"` //检测一个用户，是否有发送room heartbeat
	Uptime        int    `json:"uptime"`         //最后更新时间
	AddTime       int    `json:"add_time"`       //添加时间
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
	OnlineUids        []int        `json:"online_uids"`         //已进入房间(在线)通话的用户
	Uids              []int        `json:"uids"`                //CallUid + ReceiveUids ,只是记录下，方便函数调用
	RWLock            sync.RWMutex `json:"-"`                   //变更状态的时候使用
}

type TwinAgoraOption struct {
	ProtoMap   *util.ProtoMap `json:"-"`
	Gorm       *gorm.DB
	Log        *zap.Logger
	StaticPath string
	//RequestServiceAdapter *service.RequestServiceAdapter
	ServiceBridge *service.Bridge
}

func NewTwinAgora(op TwinAgoraOption) (*TwinAgora, error) {
	twinAgora := new(TwinAgora)
	twinAgora.Gorm = op.Gorm           //房间数据持久化
	twinAgora.CallTimeout = 8          //呼叫过程的超时时间
	twinAgora.ExecTimeout = 10         //房间运行的超时时间，room_heartbeat 也使用此值
	twinAgora.UserHeartbeatTimeout = 3 //一个用户建立的长连接，超时时间
	twinAgora.ResAcceptTimeout = 60
	twinAgora.EntryTimeout = 60
	twinAgora.Separate = "##"                         //一个房间信息转换成字符串的：分隔符
	twinAgora.RTCRoomPool = make(map[string]*RTCRoom) //房间池
	twinAgora.RTCUserPool = make(map[int]*RTCUser)    //用户池
	//twinAgora.RequestServiceAdapter = op.RequestServiceAdapter

	twinAgora.Log = op.Log
	twinAgora.Op = op

	twinAgora.CancelCtx, twinAgora.CancelFunc = context.WithCancel(context.Background())

	//错误码 文案 管理（还未用起来，后期优化）
	lang, err := util.NewErrMsg(op.Log, op.StaticPath+"/data/twin_agora.en.lang")
	if err != nil {
		twinAgora.MakeError(err.Error())
		return twinAgora, err
	}
	twinAgora.Lang = lang

	go twinAgora.ListeningBridgeMsg()
	return twinAgora, nil
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

//守护协程，检查房间超时：呼叫超时、运行超时(连接断开)
func (twinAgora *TwinAgora) CheckTimeout() {
	twinAgora.Log.Debug("twinAgora CheckTimeout demon.")
	for {
		select {
		case <-twinAgora.CancelCtx.Done():
			goto end
		default:
			//循环遍历每个房间
			for _, room := range twinAgora.RTCRoomPool {
				//呼叫过程超时
				if util.GetNowTimeSecondToInt() > room.AddTime+twinAgora.CallTimeout && room.Status == service.RTC_ROOM_STATUS_CALLING {
					twinAgora.RoomEnd(room.Id, service.RTC_ROOM_END_STATUS_TIMEOUT_CALLING)
					continue
				}
				//房间运行中超时
				if util.GetNowTimeSecondToInt() > room.Uptime+twinAgora.ExecTimeout && room.Status == service.RTC_ROOM_STATUS_EXECING {
					twinAgora.RoomEnd(room.Id, service.RTC_ROOM_END_STATUS_TIMEOUT_EXEC)
					continue
				}
			}
			//检查每个用户长连接是否超时
			for _, user := range twinAgora.RTCUserPool {
				if util.GetNowTimeSecondToInt() > user.Uptime+twinAgora.UserHeartbeatTimeout {
					twinAgora.ConnCloseProcess(user, "demon")
				}
			}
			time.Sleep(time.Millisecond * 100)
			break

		}
	}
end:
	twinAgora.Log.Debug("twinAgora CheckTimeout finish.")
}

//网关监控到有C端连接，并通过了登陆验证后，会推送事件
func (twinAgora *TwinAgora) FDCreateEvent(FDCreateEvent pb.FDCreateEvent) {
	if FDCreateEvent.UserId <= 0 {
		twinAgora.MakeError(twinAgora.Lang.NewString(400))
		return
	}

	_, ok := twinAgora.GetUserById(int(FDCreateEvent.UserId))
	if ok {
		msgInfo := twinAgora.Lang.NewReplaceOneString(405, strconv.Itoa(int(FDCreateEvent.UserId)))
		twinAgora.MakeError(msgInfo)
		twinAgora.PushMsg(int(FDCreateEvent.UserId), 500, service.RTC_PUSH_MSG_EVENT_FD_CREATE_REPEAT, msgInfo)
		return
	}

	twinAgora.Log.Warn("FDCreateEvent ,uid:" + strconv.Itoa(int(FDCreateEvent.UserId)))
	NewRTCUser := RTCUser{
		Id:      int(FDCreateEvent.UserId),
		AddTime: util.GetNowTimeSecondToInt(),
		Uptime:  util.GetNowTimeSecondToInt(),
		RoomId:  "",
	}
	twinAgora.RTCUserPool[int(FDCreateEvent.UserId)] = &NewRTCUser
}

//网关监控到用户连接断了(超时)，会回调通知服务
func (twinAgora *TwinAgora) FDCloseEvent(connCloseEvent pb.FDCloseEvent) {
	twinAgora.Log.Warn("TwinAgora ConnCloseCallback :", zap.Int("userId", int(connCloseEvent.UserId)), zap.Int("source", int(connCloseEvent.Source)))
	myRTCUser, ok := twinAgora.GetUserById(int(connCloseEvent.UserId))
	if !ok {
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(connCloseEvent.UserId))))
		return
	}
	twinAgora.ConnCloseProcess(myRTCUser, "FDCloseEvent")
}

//用户长连接 - 心跳，更新房间的最后更新时间
func (twinAgora *TwinAgora) UserHeartbeat(heartbeat pb.Heartbeat) {
	twinAgora.Log.Info("twinAgora Heartbeat , ", zap.Int64("time", heartbeat.Time), zap.Int32("uid", heartbeat.SourceUid))
	myRTCUser, ok := twinAgora.GetUserById(int(heartbeat.SourceUid))
	if !ok {
		msgInfo := twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(heartbeat.SourceUid)))
		twinAgora.MakeError(msgInfo)
		twinAgora.PushMsg(int(heartbeat.SourceUid), 500, service.RTC_PUSH_MSG_EVENT_UID_NOT_IN_MAP, msgInfo)
		return
	}

	myRTCUser.Uptime = util.GetNowTimeSecondToInt()

	if myRTCUser.RoomId == "" {
		return
	}
	myRTCRoom, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	//myRTCRoom, ok := twinAgora.RTCRoomPool[myRTCUser.RoomId]
	if err != nil { //这是种异常的情况，用户基础信息里roomId存在 ,但是在池里已经不存在了，可能是其它协程已经操作了，但是没有清空RTCUser的ROOMID
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(501, myRTCUser.RoomId))
		myRTCUser.RoomId = ""
		return
	}
	//这里是个异常，按说房间已经结束，用户基础信息应该把roomId清掉
	if myRTCRoom.Status == service.RTC_ROOM_STATUS_END {
		//交给后台守护协程处理，roomId会被清空的
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(510, myRTCUser.RoomId))
		return
	}

}

//每个房间的心跳，因为音视频使用的是声网，监控不到，就得单独再加一个心跳
func (twinAgora *TwinAgora) RoomHeartbeat(heartbeat pb.RoomHeartbeatReq) {
	twinAgora.Log.Info("twinAgora RoomHeartbeat , ", zap.Int64("time", heartbeat.Time), zap.Int32("uid", heartbeat.SourceUid))
	myRTCUser, ok := twinAgora.GetUserById(int(heartbeat.Uid))
	if !ok {
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(int(heartbeat.Uid))))
		return
	}

	myRTCRoom, err := twinAgora.GetRoomById(myRTCUser.RoomId)
	if err != nil {
		return
	}

	myRTCRoom.Uptime = util.GetNowTimeSecondToInt()
	myRTCUser.RoomHeartbeat = util.GetNowTimeSecondToInt()

	if myRTCRoom.Status != service.RTC_ROOM_STATUS_EXECING {
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(511, heartbeat.RoomId))
		return
	}

}

//给客户端推送消息，主要是一些错误信息
func (twinAgora *TwinAgora) PushMsg(uid int, code int, eventId int, content string) {
	//util.MyPrint("PushMsg uid:", uid, ", code ", code, " , eventId:", eventId, " , content:", content)
	//pushMsg := pb.PushMsg{
	//	Code:    int32(code),
	//	Uid:     int32(uid),
	//	EventId: int32(eventId),
	//	Content: content,
	//}
	//data, _ := proto.Marshal(&pushMsg)
	////twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(uid), "SC_PushMsg", pushMsg)
	//twinAgora.Op.ServiceBridge.CallByName("Gateway", "SC_PushMsg", string(data), "", 0)
	////conn.GatewaySendMsgByUid(int32(uid), "SC_PushMsg", pushMsg)
}

//连接断开或超时处理
func (twinAgora *TwinAgora) ConnCloseProcess(rtcUserRTCUser *RTCUser, source string) {
	twinAgora.Log.Warn("ConnCloseProcess source: " + source + " , uid:" + strconv.Itoa(rtcUserRTCUser.Id) + " roomId:" + rtcUserRTCUser.RoomId)
	if rtcUserRTCUser.RoomId == "" {
		twinAgora.DelUserById(rtcUserRTCUser.Id)
		return
	}
	myRTCRoom, err := twinAgora.GetRoomById(rtcUserRTCUser.RoomId)
	if err != nil { //这是种异常的情况，用户基础信息里roomId存在 ,但是在池里已经不存在了，可能是其它协程已经操作了，但是没有清空RTCUser的ROOMID
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(501, rtcUserRTCUser.RoomId))
		twinAgora.DelUserById(rtcUserRTCUser.Id)
		return
	}

	if myRTCRoom.Status == service.RTC_ROOM_STATUS_END {
		//这也是异常情况，池子里虽然有个房间，但是状态是已经结束了，可能后台协程也没有来得及处理
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(510, rtcUserRTCUser.RoomId))
		twinAgora.DelUserById(rtcUserRTCUser.Id)
		return
	}
	//目前是1v1视频，只要有一个人拒绝|断线，即结束，这里后期优化一下吧
	twinAgora.DelUserById(rtcUserRTCUser.Id) //这个得放在上面，可能有人断网后立即重连，并行的话有问题
	twinAgora.RoomEnd(myRTCRoom.Id, service.RTC_ROOM_END_STATUS_CONN_CLOSE)

}

//已结束的房间要做:
//1. 持久化
//2. 房间池内删除该元素
//3. 更新用户池内：在线用户的房间ID清除
func (twinAgora *TwinAgora) RoomEnd(roomId string, endStatus int) {
	twinAgora.Log.Warn("RoomEnd id:" + roomId + " , endStatus:" + strconv.Itoa(endStatus))
	roomInfo, err := twinAgora.GetRoomById(roomId)
	if err != nil {
		return
	}
	//要修改房间状态，要持久化，最终还要删除内存池中的记录，所以要加写锁
	roomInfo.RWLock.Lock()
	defer roomInfo.RWLock.Unlock()

	roomInfo.Status = service.RTC_ROOM_STATUS_END
	roomInfo.EndStatus = endStatus

	twinAgora.StoreHistory(roomInfo)

	for _, uid := range roomInfo.OnlineUids {
		myRTCUser, ok := twinAgora.GetUserById(uid)
		if ok {
			//conn, ok2 := twinAgora.ConnManager.Pool[int32(uid)]
			//if ok2 {
			peopleLeaveRes := pb.PeopleLeaveRes{}
			peopleLeaveRes.Uid = int32(uid)
			peopleLeaveRes.Channel = roomInfo.Channel
			peopleLeaveRes.RoomId = roomInfo.Id

			//data, _ := proto.Marshal(&peopleLeaveRes)
			//twinAgora.Op.ServiceBridge.CallByName("Gateway", "SC_PeopleLeave", string(data), "", 0)
			callGatewayMsg := service.CallGatewayMsg{ServiceName: "TwinAgora", FunName: "SC_PeopleLeave", SourceUid: int32(uid), Data: &peopleLeaveRes}
			twinAgora.Op.ServiceBridge.CallGateway(callGatewayMsg)

			//twinAgora.RequestServiceAdapter.GatewaySendMsgByUid(int32(uid), "SC_PeopleLeave", peopleLeaveRes)
			//conn.GatewaySendMsgByUid(int32(uid), "SC_PeopleLeave", peopleLeaveRes)
			//}
			myRTCUser.RoomId = ""
		}
	}
	twinAgora.Log.Warn("delete room:" + roomInfo.Id)
	delete(twinAgora.RTCRoomPool, roomInfo.Id)

	twinAgora.Log.Warn("RoomEnd ok , roomId: " + roomId + " , endStatus:" + strconv.Itoa(endStatus))
}

//持久化到DB中
func (twinAgora *TwinAgora) StoreHistory(RTCRoom *RTCRoom) error {
	var twinAgoraRoomRow model.TwinAgoraRoom
	twinAgora.Gorm.Where("room_id = ? ", RTCRoom.Id).First(&twinAgoraRoomRow)
	if twinAgoraRoomRow.Id > 0 {
		return twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(520, RTCRoom.Id))
	}
	ReceiveUidsStr := util.ArrCoverStr(RTCRoom.ReceiveUids, ",")
	ReceiveUidsAcceptStr := util.ArrCoverStr(RTCRoom.ReceiveUidsAccept, ",")
	ReceiveUidsDenyStr := util.ArrCoverStr(RTCRoom.ReceiveUidsDeny, ",")
	UidsStr := util.ArrCoverStr(RTCRoom.Uids, ",")

	str := RTCRoom.Id + twinAgora.Separate + RTCRoom.Channel + twinAgora.Separate + strconv.Itoa(RTCRoom.AddTime) + twinAgora.Separate + strconv.Itoa(RTCRoom.Status) + twinAgora.Separate + strconv.Itoa(RTCRoom.EndStatus) + twinAgora.Separate + strconv.Itoa(RTCRoom.CallUid) + twinAgora.Separate + ReceiveUidsStr + twinAgora.Separate + ReceiveUidsAcceptStr + UidsStr + ReceiveUidsAcceptStr + ReceiveUidsDenyStr
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
	twinAgora.Log.Warn("DelUserById:" + strconv.Itoa(uid))
	delete(twinAgora.RTCUserPool, uid)
}

func (twinAgora *TwinAgora) GetUserById(uid int) (mmRTCUserRTCUser *RTCUser, rs bool) {
	if uid <= 0 {
		twinAgora.MakeError(twinAgora.Lang.NewString(400))
		return mmRTCUserRTCUser, false
	}
	//util.MyPrint("GetUserById uid:", uid, twinAgora.RTCUserPool)
	myRTCUser, ok := twinAgora.RTCUserPool[uid]
	if !ok {
		twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(401, strconv.Itoa(uid)))
	}
	return myRTCUser, ok
}

func (twinAgora *TwinAgora) MakeError(errMsg string) error {
	twinAgora.Log.Error("*********=====MakeError : " + errMsg)
	return errors.New(errMsg)
}

func (twinAgora *TwinAgora) GetRoomById(id string) (room *RTCRoom, err error) {
	if id == "" {
		return room, twinAgora.MakeError(twinAgora.Lang.NewString(500))
	}

	room, ok := twinAgora.RTCRoomPool[id]
	if !ok {
		return room, twinAgora.MakeError(twinAgora.Lang.NewReplaceOneString(501, id))
	}

	return room, nil
}

//func (twinAgora *TwinAgora) InitErrorMsg() {
//	errMap := make(map[int]string)
//	errMap[501] = ERR_UID_ZERO
//	errMap[502] = "callPeopleReq.PeopleType err , now only support: calling doctor (PeopleType= 2) "
//	errMap[504] = "not support : TargetUid > 0 "
//	errMap[511] = ERR_ROOM_ID_NOT_IN_MAP
//	errMap[520] = "exist <callPeople> record : ，don't repeat opt"
//	errMap[521] = "exist <room talking> record : ，don't repeat opt"
//	errMap[522] = "The room has end ,wait demon coroutines process recycle"
//	errMap[503] = "DB not have role=doctor user"
//	errMap[510] = "All doctor user not online..."
//
//	twinAgora.Err = errMap
//}
