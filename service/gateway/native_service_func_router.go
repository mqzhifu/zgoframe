package gateway

import (
	"errors"
	"strconv"
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
	requestSendMsg := pb.Msg{}
	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))

	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		err = util.ConnParserContentMsg(msg, &requestLogin)
	case "CS_Pong": //
		err = util.ConnParserContentMsg(msg, &requestClientPong)
	case "CS_Ping":
		err = util.ConnParserContentMsg(msg, &requestClientPing)
	case "CS_Heartbeat": //心跳
		err = util.ConnParserContentMsg(msg, &requestClientHeartbeat)
	case "CS_ProjectPushMsg": //某个服务想给其它服务推送消息
		err = util.ConnParserContentMsg(msg, &requestProjectPushMsg)
	case "SC_SendMsg":
		err = util.ConnParserContentMsg(msg, &requestSendMsg)
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
		//网关自己要维护一个心跳，主要是更新原始FD的时间、计算RTT等
		gateway.heartbeat(requestClientHeartbeat)
		//心跳还要广播给后面的所有微服务
		//requestClientPing.SourceUid = conn.UserId
		requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(&requestClientHeartbeat, msg.SourceUid)
		msg, _, _ := gateway.Netway.ConnManager.MakeMsgByServiceFuncName(msg.SourceUid, "Gateway", "CS_Heartbeat", requestClientHeartbeatStrByte)
		gateway.BroadcastService("CS_Heartbeat", msg)
	case "CS_ProjectPushMsg":
		//gateway.Log.Debug("CS_ProjectPushMsg message,but no implementation......")
		conn, exist := gateway.Netway.ConnManager.GetConnPoolById(msg.TargetUid)
		if !exist {
			err = errors.New("msg.TargetUid not in conn pool :" + strconv.Itoa(int(msg.TargetUid)))
			return data, err
		}
		conn.SendMsgCompressByName(protoServiceFunc.ServiceName, protoServiceFunc.FuncName, msg.Content)
	case "SC_SendMsg":
		conn, _ := gateway.Netway.ConnManager.GetConnPoolById(requestSendMsg.TargetUid)
		conn.SendMsgCompressByName(protoServiceFunc.ServiceName, protoServiceFunc.FuncName, msg.Content)
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

////总：路由器，这里分成了两类：gateway 自解析 和 代理后方服务的请求
//func (gateway *Gateway) Router(msg pb.Msg) (data interface{}, err error) {
//	actionInfo, _ := gateway.NetWayOption.ProtoMap.GetServiceFuncById(int(msg.SidFid))
//	gateway.Log.Info("service gateway router , ServiceName:" + actionInfo.ServiceName + " FuncName:" + actionInfo.FuncName + " SidFid:" + strconv.Itoa(int(msg.SidFid)))
//	serviceName := actionInfo.ServiceName
//	switch serviceName {
//	case "Gateway":
//		data, err = gateway.RouterServiceGateway(msg)
//	case "FrameSync":
//		data, err = gateway.RouterServiceSync(msg)
//	case "GameMatch":
//		data, err = gateway.RouterServiceGameMatch(msg)
//	case "TwinAgora":
//		data, err = gateway.RouterServiceTwinAgora(msg)
//	default:
//		gateway.Log.Error("netWay Router err.")
//		return nil, errors.New("netWay Router err.")
//	}
//	return data, err
//}
//
