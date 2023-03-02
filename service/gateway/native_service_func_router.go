package gateway

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"strconv"
	"strings"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/util"
)

func (gateway *Gateway) ListeningBridgeMsg() {
	for {
		select {
		case msg := <-gateway.ServiceBridge.NativeServiceList.Gateway:

			gateway.NativeServiceFuncRouter(msg)
		default:
			time.Sleep(time.Millisecond * service.BRIDGE_SLEEP_TIME)
		}
	}
}

//网关自解析的路由
func (gateway *Gateway) NativeServiceFuncRouter(msg pb.Msg) (data interface{}, err error) {
	prefix := "NativeServiceRouter "

	requestLogin := pb.Login{}
	requestClientPong := pb.PongRes{}
	requestClientPing := pb.PingReq{}
	requestClientHeartbeat := pb.Heartbeat{}
	requestProjectPushMsg := pb.ProjectPushMsg{}
	requestFDCloseEvent := pb.FDCloseEvent{}
	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))

	//ProtocolCtrlInfo := gateway.Netway.ConnManager.GetPlayerCtrlInfoById(msg.TargetUid)
	//util.MyPrint("ProtocolCtrlInfo:", ProtocolCtrlInfo.ContentType, ProtocolCtrlInfo.ProtocolType)
	//msg.ContentType = ProtocolCtrlInfo.ContentType
	//msg.ProtocolType = ProtocolCtrlInfo.ProtocolType

	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		err = util.ConnParserContentMsg(msg, &requestLogin)
	case "CS_Pong": //
		err = util.ConnParserContentMsg(msg, &requestClientPong)
	case "CS_Ping":
		err = util.ConnParserContentMsg(msg, &requestClientPing)
	case "CS_Heartbeat": //心跳
		err = util.ConnParserContentMsg(msg, &requestClientHeartbeat)
	case "SC_ProjectPushMsg": //某个服务想给其它服务推送消息
		err = util.ConnParserContentMsg(msg, &requestProjectPushMsg)
	case "SC_SendMsg": //这个方法特殊，msg.content 的值是可能是各种类型，虽然可以用msg.sourceService + msg.sourceFunc 解出来，但意义不在，不解了，后面也是直接发送出去了
		//err = util.ConnParserContentMsg(msg, &requestSendMsg)
	case "FdClose":
		err = util.ConnParserContentMsg(msg, &requestFDCloseEvent)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		//长连接建立成功后，首先就得登陆，成功后，等于该FD是正常的，会创:建新绑定关系  user<==>fd。  这里有个 BUG，LOGIN 函数只能在第一次调用，回头加个限定
		FDCreateEvent := pb.FDCreateEvent{UserId: msg.SourceUid, SourceUid: msg.SourceUid}
		requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(&FDCreateEvent, msg.SourceUid)
		msgFDCreateEvent, _, _ := gateway.Netway.ConnManager.MakeMsgByServiceFuncName(msg.SourceUid, "Gateway", "FdCreate", requestClientHeartbeatStrByte)
		//msgFDCreateEvent.SourceUid = conn.UserId
		//将消息广播给所有微服务
		gateway.BroadcastService("FdCreate", msgFDCreateEvent)

		loginRes := pb.LoginRes{
			Code:   200,
			ErrMsg: "",
			Uid:    msg.SourceUid,
		}
		conn, exist := gateway.Netway.ConnManager.GetConnPoolById(msg.SourceUid)
		if !exist {
			util.MyPrint("GetConnPoolById conn not exist ")
		}
		//告知玩家：登陆结果
		conn.SendMsgCompressByName("Gateway", "SC_Login", &loginRes)

	case "FdClose":
		gateway.BroadcastService("FdClose", msg)
	case "CS_Ping":
		//requestClientPing.SourceUid = conn.UserId
		gateway.clientPing(requestClientPing)
	case "CS_Pong":
		//requestClientPong.SourceUid = conn.UserId
		gateway.ClientPong(requestClientPong)
	case "CS_Heartbeat":
		util.MyPrint("", msg.SourceUid)
		requestClientHeartbeat.SourceUid = msg.SourceUid
		//网关自己要维护一个心跳，主要是更新原始FD的时间、计算RTT等
		gateway.heartbeat(requestClientHeartbeat)
		//心跳还要广播给后面的所有微服务
		//requestClientPing.SourceUid = conn.UserId
		requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(&requestClientHeartbeat, msg.SourceUid)
		msgN, _, _ := gateway.Netway.ConnManager.MakeMsgByServiceFuncName(msg.SourceUid, "Gateway", "CS_Heartbeat", requestClientHeartbeatStrByte)
		gateway.BroadcastService("CS_Heartbeat", msgN)
	case "SC_ProjectPushMsg":
		TargetUidArr := strings.Split(requestProjectPushMsg.TargetUids, ",")
		uids := []int{}
		gateway.Log.Debug("SC_ProjectPushMsg ,TargetUids:" + requestProjectPushMsg.TargetUids)
		for _, uidStr := range TargetUidArr {
			uid, _ := strconv.Atoi(uidStr)
			if uid <= 0 {
				err = errors.New("some uid <=0")
				gateway.Log.Error("some uid <=0")
				return data, err
			}
			uids = append(uids, uid)
		}

		for _, uid := range uids {
			conn, exist := gateway.Netway.ConnManager.GetConnPoolById(int32(uid))
			if !exist {
				errMsg := "msg.TargetUid not in conn pool :" + strconv.Itoa(int(msg.TargetUid))
				err = errors.New("msg.TargetUid not in conn pool :" + strconv.Itoa(int(msg.TargetUid)))
				gateway.Log.Error(errMsg)
				continue
			}
			conn.SendMsgCompressByName("Gateway", "SC_ProjectPushMsg", &requestProjectPushMsg)
		}

	case "SC_SendMsg":
		SC_SendMsgPrefix := "gateway router SC_SendMsg "
		sourceServiceFunc, empty := gateway.Netway.Option.ProtoMap.GetServiceFuncBySidFid(int(msg.SourceServiceId), int(msg.SourceFuncId))
		if empty {
			errMsg := SC_SendMsgPrefix + " GetServiceFuncBySidFid empty , msg.SourceServiceId:" + strconv.Itoa(int(msg.SourceServiceId)) + " msg.SourceFuncId: " + strconv.Itoa(int(msg.SourceFuncId))
			gateway.Log.Error(errMsg)
			return data, errors.New(errMsg)
		}
		conn, exist := gateway.Netway.ConnManager.GetConnPoolById(msg.TargetUid)
		if !exist {
			errMsg := SC_SendMsgPrefix + " GetServiceFuncBySidFid empty , msg.TargetUid:" + strconv.Itoa(int(msg.TargetUid))
			gateway.Log.Error(errMsg)
			return data, errors.New(errMsg)
		}
		gateway.RouterSendMsg(msg, sourceServiceFunc, conn)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

func (gateway *Gateway) MakeRouterErrNotFound(prefix string, funcName string, index string) string {
	errMsg := prefix + " , FuncName not found-" + index + " :" + funcName
	gateway.Log.Error(prefix + " , FuncName not found-" + index + " :" + funcName)
	return errMsg
}

//总：路由器，这里分成了两类：gateway 自解析 和 代理后方服务的请求
func (gateway *Gateway) RouterSendMsg(msg pb.Msg, sourceServiceFunc util.ProtoServiceFunc, conn *util.Conn) (data interface{}, err error) {
	//actionInfo, _ := gateway.NetWayOption.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	//gateway.Log.Info("service gateway router , ServiceName:" + actionInfo.ServiceName + " FuncName:" + actionInfo.FuncName + " SidFid:" + strconv.Itoa(int(msg.SidFid)))
	//serviceName := actionInfo.ServiceName
	switch sourceServiceFunc.ServiceName {
	case "FrameSync":
		data, err = gateway.RouterFuncSendMsgFrameSync(msg, sourceServiceFunc, conn)
	case "GameMatch":
		data, err = gateway.RouterFuncSendMsgGameMatch(msg, sourceServiceFunc, conn)
	case "TwinAgora":
		data, err = gateway.RouterFuncSendMsgTwinAgora(msg, sourceServiceFunc, conn)
	default:
		gateway.Log.Error("netWay Router err.")
		return nil, errors.New("netWay Router err.")
	}
	return data, err
}

func (gateway *Gateway) RouterFuncSendMsgFrameSync(msg pb.Msg, sourceServiceFunc util.ProtoServiceFunc, conn *util.Conn) (data interface{}, err error) {
	prefix := "RouterFuncSendMsgFrameSync "
	reqReadyTimeout := pb.ReadyTimeout{}
	reqEnterBattle := pb.EnterBattle{}
	reqLogicFrame := pb.LogicFrame{}
	reqRoomHistorySets := pb.RoomHistorySets{}
	reqRoomBaseInfo := pb.RoomBaseInfo{}
	reqOtherPlayerOffline := pb.OtherPlayerOffline{}
	reqPlayerOver := pb.PlayerOver{}
	reqPlayerResumeGame := pb.PlayerResumeGame{}
	reqStartBattle := pb.StartBattle{}
	reqRestartGame := pb.RestartGame{}
	reqGameOver := pb.GameOver{}
	reqPlayerState := pb.PlayerState{}
	reqReqHeartbeat := pb.Heartbeat{}

	switch sourceServiceFunc.FuncName {
	case "SC_ReadyTimeout":
		err = proto.Unmarshal([]byte(msg.Content), &reqReadyTimeout)
		//err = util.ConnParserContentMsg(msg, &reqReadyTimeout)
	case "SC_EnterBattle":
		err = proto.Unmarshal([]byte(msg.Content), &reqEnterBattle)
		//err = util.ConnParserContentMsg(msg, &reqEnterBattle)
	case "SC_LogicFrame":
		err = proto.Unmarshal([]byte(msg.Content), &reqLogicFrame)
		//err = util.ConnParserContentMsg(msg, &reqLogicFrame)
	case "SC_RoomHistory":
		err = proto.Unmarshal([]byte(msg.Content), &reqRoomHistorySets)
		//err = util.ConnParserContentMsg(msg, &reqRoomHistorySets)
	case "SC_RoomBaseInfo":
		err = proto.Unmarshal([]byte(msg.Content), &reqRoomBaseInfo)
		//err = util.ConnParserContentMsg(msg, &reqRoomBaseInfo)
	case "SC_OtherPlayerOffline":
		err = proto.Unmarshal([]byte(msg.Content), &reqOtherPlayerOffline)
		//err = util.ConnParserContentMsg(msg, &reqOtherPlayerOffline)
	case "SC_OtherPlayerOver":
		err = proto.Unmarshal([]byte(msg.Content), &reqPlayerOver)
		//err = util.ConnParserContentMsg(msg, &reqPlayerOver)
	case "SC_OtherPlayerResumeGame":
		err = proto.Unmarshal([]byte(msg.Content), &reqPlayerResumeGame)
		//err = util.ConnParserContentMsg(msg, &reqPlayerResumeGame)
	case "SC_StartBattle":
		err = proto.Unmarshal([]byte(msg.Content), &reqStartBattle)
		//err = util.ConnParserContentMsg(msg, &reqStartBattle)
	case "SC_RestartGame":
		err = proto.Unmarshal([]byte(msg.Content), &reqRestartGame)
		//err = util.ConnParserContentMsg(msg, &reqRestartGame)
	case "SC_GameOver":
		err = proto.Unmarshal([]byte(msg.Content), &reqGameOver)
		//err = util.ConnParserContentMsg(msg, &reqGameOver)
	case "SC_Heartbeat":
		err = proto.Unmarshal([]byte(msg.Content), &reqReqHeartbeat)
		//err = util.ConnParserContentMsg(msg, &reqReqHeartbeat)
	case "SC_PlayerState":
		err = proto.Unmarshal([]byte(msg.Content), &reqPlayerState)
		//err = util.ConnParserContentMsg(msg, &reqPlayerState)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, sourceServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
	}

	switch sourceServiceFunc.FuncName {
	case "SC_ReadyTimeout":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqReadyTimeout)
	case "SC_EnterBattle":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqEnterBattle)
	case "SC_LogicFrame":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqLogicFrame)
	case "SC_RoomHistory":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqRoomHistorySets)
	case "SC_RoomBaseInfo":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqRoomBaseInfo)
	case "SC_OtherPlayerOffline":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqOtherPlayerOffline)
	case "SC_OtherPlayerOver":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqPlayerOver)
	case "SC_OtherPlayerResumeGame":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqPlayerResumeGame)
	case "SC_StartBattle":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqStartBattle)
	case "SC_RestartGame":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqRestartGame)
	case "SC_GameOver":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqGameOver)
	case "SC_Heartbeat":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqReqHeartbeat)
	case "SC_PlayerState":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqPlayerState)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, sourceServiceFunc.FuncName, "2"))
	}

	return data, nil

}

func (gateway *Gateway) RouterFuncSendMsgTwinAgora(msg pb.Msg, sourceServiceFunc util.ProtoServiceFunc, conn *util.Conn) (data interface{}, err error) {
	prefix := "RouterFuncSendMsgTwinAgora "

	reqCallPeopleRes := pb.CallPeopleRes{}
	reqCancelCallPeopleReq := pb.CancelCallPeopleReq{}
	reqPeopleEntry := pb.PeopleEntry{}
	reqPeopleLeaveRes := pb.PeopleLeaveRes{}
	reqCallReply := pb.CallReply{}
	reqCallVote := pb.CallVote{}

	switch sourceServiceFunc.FuncName {
	case "SC_CallPeople":
		err = proto.Unmarshal([]byte(msg.Content), &reqCallPeopleRes)
	case "SC_CancelCallPeople":
		err = proto.Unmarshal([]byte(msg.Content), &reqCancelCallPeopleReq)
	case "SC_PeopleEntry":
		err = proto.Unmarshal([]byte(msg.Content), &reqPeopleEntry)
	case "SC_PeopleLeave":
		err = proto.Unmarshal([]byte(msg.Content), &reqPeopleLeaveRes)
	case "SC_CallReply":
		err = proto.Unmarshal([]byte(msg.Content), &reqCallReply)
	case "SC_CallPeopleAccept":
		err = proto.Unmarshal([]byte(msg.Content), &reqCallVote)
	case "SC_CallPeopleDeny":
		err = proto.Unmarshal([]byte(msg.Content), &reqCallVote)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, sourceServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
	}

	switch sourceServiceFunc.FuncName {
	case "SC_CallPeople":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqCallPeopleRes)
	case "SC_CancelCallPeople":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqCancelCallPeopleReq)
	case "SC_PeopleEntry":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqPeopleEntry)
	case "SC_PeopleLeave":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqPeopleLeaveRes)
	case "SC_CallReply":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqCallReply)
	case "SC_CallPeopleAccept":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqCallVote)
	case "SC_CallPeopleDeny":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqCallVote)

	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, sourceServiceFunc.FuncName, "2"))
	}

	return data, nil

}

func (gateway *Gateway) RouterFuncSendMsgGameMatch(msg pb.Msg, sourceServiceFunc util.ProtoServiceFunc, conn *util.Conn) (data interface{}, err error) {
	prefix := "RouterFuncSendMsgGameMatch "
	reqGameMatchOptResult := pb.GameMatchOptResult{}
	reqHeartbeat := pb.Heartbeat{}

	switch sourceServiceFunc.FuncName {
	case "SC_GameMatchOptResult":
		err = proto.Unmarshal([]byte(msg.Content), &reqGameMatchOptResult)
		//err = util.ConnParserContentMsg(msg, &reqGameMatchOptResult)
	case "SC_Heartbeat":
		err = proto.Unmarshal([]byte(msg.Content), &reqHeartbeat)
		//err = util.ConnParserContentMsg(msg, &reqHeartbeat)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, sourceServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
	}

	switch sourceServiceFunc.FuncName {
	case "SC_GameMatchOptResult":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqGameMatchOptResult)
	case "SC_Heartbeat":
		err = conn.SendMsgCompressByName(sourceServiceFunc.ServiceName, sourceServiceFunc.FuncName, &reqHeartbeat)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, sourceServiceFunc.FuncName, "2"))
	}

	return data, nil

}
