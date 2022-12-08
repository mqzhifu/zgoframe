package gateway

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/service/frame_sync"
	gamematch "zgoframe/service/game_match"
	"zgoframe/service/seed_business"
	"zgoframe/util"
)

//这是个快捷变量，目前所有代码均在一起，直接挂在这个变量上即可，后期所有服务分拆出去，网关没那么多附加功能此变量就没用了
type MyServiceList struct {
	//Match      *Match
	GameMatch *gamematch.GameMatch
	FrameSync *frame_sync.FrameSync
	//RoomManage *frame_sync.RoomManager
	TwinAgora *seed_business.TwinAgora
}

type Gateway struct {
	GrpcManager           *util.GrpcManager              //通过GRPC反射代理其它微服务
	Log                   *zap.Logger                    //日志
	Netway                *util.NetWay                   //长连接公共类
	NetWayOption          util.NetWayOption              //长连接公共类的初始化参数
	MyServiceList         *MyServiceList                 //快捷访问内部微服务
	RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
}

/*
	网关，目前主要是分为3个主要功能：
	1. http 代理 grpc(中等)
	2. 长连接代理(重点)
	3. http 代理 http(鸡肋)
*/
func NewGateway(grpcManager *util.GrpcManager, log *zap.Logger, requestServiceAdapter *service.RequestServiceAdapter) *Gateway {
	gateway := new(Gateway)
	gateway.GrpcManager = grpcManager
	gateway.Log = log
	gateway.MyServiceList = &MyServiceList{}
	gateway.RequestServiceAdapter = requestServiceAdapter
	go gateway.ListeningMsg()
	return gateway
}
func (gateway *Gateway) ListeningMsg() {
	for {
		select {
		case GatewayMsg := <-gateway.RequestServiceAdapter.QueueGatewayMsg:
			conn, exist := gateway.Netway.ConnManager.GetConnPoolById(GatewayMsg.Uid)
			if !exist {
				gateway.Log.Error("ListeningMsg conn empty uid:" + strconv.Itoa(int(GatewayMsg.Uid)))
				break
			}
			conn.SendMsgCompressByUid(GatewayMsg.Uid, GatewayMsg.ActionName, GatewayMsg.Data)
		case ServiceMsg := <-gateway.RequestServiceAdapter.QueueServiceMsg:
			//工程太大不写了
			util.MyPrint("gateway ListeningMsg:", ServiceMsg)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}

//开启长连接监听
func (gateway *Gateway) StartSocket(netWayOption util.NetWayOption) (*util.NetWay, error) {
	gateway.Log.Info("gateway StartSocket:")
	netWayOption.RouterBack = gateway.Router //公共回调 路由器，用于给最底层的长连接公共类回调
	//创建长连接:底层-公共类
	gateway.NetWayOption = netWayOption
	netWay, err := util.NewNetWay(netWayOption)
	gateway.Netway = netWay
	return netWay, err
}

//广播给所有服务，如：心跳 PING PONG 关闭事件(不广播给gateway)
func (gateway *Gateway) BroadcastService(msg pb.Msg, conn *util.Conn) {
	gateway.Log.Debug("BroadcastService funcId:" + strconv.Itoa(int(msg.FuncId)))
	gateway.RouterServiceSync(msg, conn)
	//gateway.RouterServiceGameMatch(msg, conn)
	gateway.RouterServiceTwinAgora(msg, conn)
}
func (gateway *Gateway) MakeRouterErrNotFound(prefix string, funcName string, index string) string {
	errMsg := prefix + " , protoServiceFunc.FuncName not found-" + index + " :" + funcName
	gateway.Log.Error(prefix + " , protoServiceFunc.FuncName not found-" + index + " :" + funcName)
	return errMsg
}

//总：路由器，这里分成了两类：gateway 自解析 和 代理后方服务的请求
func (gateway *Gateway) Router(msg pb.Msg, conn *util.Conn) (data interface{}, err error) {
	actionInfo, _ := gateway.NetWayOption.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	gateway.Log.Info("service gateway router , ServiceName:" + actionInfo.ServiceName + " FuncName:" + actionInfo.FuncName + " SidFid:" + strconv.Itoa(int(msg.SidFid)))
	serviceName := actionInfo.ServiceName
	switch serviceName {
	case "Gateway":
		data, err = gateway.RouterServiceGateway(msg, conn)
	case "FrameSync":
		data, err = gateway.RouterServiceSync(msg, conn)
	case "GameMatch":
		data, err = gateway.RouterServiceGameMatch(msg, conn)
	case "TwinAgora":
		data, err = gateway.RouterServiceTwinAgora(msg, conn)
	default:
		gateway.Log.Error("netWay Router err.")
		return nil, errors.New("netWay Router err.")
	}
	return data, err
}

func (gateway *Gateway) RouterServiceTwinAgora(msg pb.Msg, conn *util.Conn) (data []byte, err error) {
	prefix := "RouterServiceTwinAgora"

	reqCallPeople := pb.CallPeopleReq{}
	reqFDCloseEvent := pb.FDCloseEvent{}
	reqHeartbeat := pb.Heartbeat{}
	reqFDCreateEvent := pb.FDCreateEvent{}
	reqCallVote := pb.CallVote{}
	reqRoomHeartbeat := pb.RoomHeartbeatReq{}
	reqCancelCallPeople := pb.CancelCallPeopleReq{}
	reqPeopleEntry := pb.PeopleEntry{}
	reqPeopleLeaveRes := pb.PeopleLeaveRes{}

	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))

	switch protoServiceFunc.FuncName {
	case "CS_CallPeople":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqCallPeople, conn.UserId)
	case "CS_PeopleLeave":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqPeopleLeaveRes, conn.UserId)
	case "CS_CancelCallPeople":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqCancelCallPeople, conn.UserId)
	case "FdClose":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqFDCloseEvent, 0)
	case "CS_PeopleEntry":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqPeopleEntry, conn.UserId)
	case "FdCreate":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqFDCreateEvent, conn.UserId)
	case "CS_RoomHeartbeat":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqRoomHeartbeat, conn.UserId)
	case "CS_Heartbeat":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqHeartbeat, conn.UserId)
	case "CS_CallPeopleAccept":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqCallVote, conn.UserId)
	case "CS_CallPeopleDeny":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqCallVote, conn.UserId)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_CallPeople":
		gateway.MyServiceList.TwinAgora.CallPeople(reqCallPeople)
	case "CS_Heartbeat":
		gateway.MyServiceList.TwinAgora.UserHeartbeat(reqHeartbeat)
	case "CS_PeopleLeave":
		gateway.MyServiceList.TwinAgora.PeopleLeave(reqPeopleLeaveRes)
	case "CS_RoomHeartbeat":
		gateway.MyServiceList.TwinAgora.RoomHeartbeat(reqRoomHeartbeat)
	case "FdClose":
		gateway.MyServiceList.TwinAgora.FDCloseEvent(reqFDCloseEvent)
	case "CS_PeopleEntry":
		gateway.MyServiceList.TwinAgora.PeopleEntry(reqPeopleEntry)
	case "CS_CancelCallPeople":
		gateway.MyServiceList.TwinAgora.CancelCallPeople(reqCancelCallPeople)
	case "FdCreate":
		gateway.MyServiceList.TwinAgora.FDCreateEvent(reqFDCreateEvent)
	case "CS_CallPeopleAccept":
		gateway.MyServiceList.TwinAgora.CallPeopleAccept(reqCallVote)
	case "CS_CallPeopleDeny":
		gateway.MyServiceList.TwinAgora.CallPeopleDeny(reqCallVote)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

//网关自解析的路由
func (gateway *Gateway) RouterServiceGateway(msg pb.Msg, conn *util.Conn) (data interface{}, err error) {
	prefix := "RouterServiceGateway "

	requestLogin := pb.Login{}
	requestClientPong := pb.PongRes{}
	requestClientPing := pb.PingReq{}
	requestClientHeartbeat := pb.Heartbeat{}
	requestProjectPushMsg := pb.ProjectPushMsg{}
	requestFDCloseEvent := pb.FDCloseEvent{}
	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))

	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestLogin, conn.UserId)
	case "CS_Pong": //
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestClientPong, conn.UserId)
	case "CS_Ping":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestClientPing, conn.UserId)
	case "CS_Heartbeat": //心跳
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestClientHeartbeat, conn.UserId)
	case "CS_ProjectPushMsg": //某个服务想给其它服务推送消息
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestProjectPushMsg, conn.UserId)
	case "FdClose":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestFDCloseEvent, conn.UserId)
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
		FDCreateEvent := pb.FDCreateEvent{UserId: conn.UserId, SourceUid: conn.UserId}
		requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(&FDCreateEvent, conn.UserId)
		msgFDCreateEvent, _, _ := gateway.Netway.ConnManager.MakeMsgByActionName(conn.UserId, "FdCreate", requestClientHeartbeatStrByte)
		msgFDCreateEvent.SourceUid = conn.UserId
		//将消息广播给所有微服务
		gateway.BroadcastService(msgFDCreateEvent, conn)
	case "FdClose":
		gateway.BroadcastService(msg, conn)
	case "CS_Ping":
		requestClientPing.SourceUid = conn.UserId
		gateway.clientPing(requestClientPing)
	case "CS_Pong":
		requestClientPong.SourceUid = conn.UserId
		gateway.ClientPong(requestClientPong, conn)
	case "CS_Heartbeat":
		//网关自己要维护一个心跳，主要是更新原始FD的时间、计算RTT等
		gateway.heartbeat(requestClientHeartbeat, conn)
		//心跳还要广播给后面的所有微服务
		requestClientPing.SourceUid = conn.UserId
		requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(&requestClientHeartbeat, conn.UserId)
		msg, _, _ := gateway.Netway.ConnManager.MakeMsgByActionName(conn.UserId, "CS_Heartbeat", requestClientHeartbeatStrByte)
		gateway.BroadcastService(msg, conn)
	case "CS_ProjectPushMsg":
		gateway.Log.Debug("CS_ProjectPushMsg message,but no implementation......")
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

func (gateway *Gateway) RouterServiceGameMatch(msg pb.Msg, conn *util.Conn) (data []byte, err error) {
	prefix := "RouterServiceGameMatch"

	requestPlayerMatchSign := pb.GameMatchSign{}
	requestPlayerMatchSignCancel := pb.GameMatchPlayerCancel{}
	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSign, conn.UserId)
	case "CS_PlayerMatchSignCancel":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSignCancel, conn.UserId)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}

	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		_, e := gateway.MyServiceList.GameMatch.PlayerJoin(requestPlayerMatchSign)
		if e != nil {
			gateway.Log.Debug("PlayerJoin return e: " + e.Error())
		}
	case "CS_PlayerMatchSignCancel":
		requestPlayerMatchSignCancel.SourceUid = conn.UserId
		gateway.MyServiceList.GameMatch.Cancel(requestPlayerMatchSignCancel)

	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

//帧同步的路由
func (gateway *Gateway) RouterServiceSync(msg pb.Msg, conn *util.Conn) (data []byte, err error) {
	prefix := "RouterServiceSync "
	requestLogicFrame := pb.LogicFrame{}
	requestPlayerResumeGame := pb.PlayerResumeGame{}
	requestPlayerReady := pb.PlayerReady{}
	requestPlayerOver := pb.PlayerOver{}
	requestRoomHistory := pb.ReqRoomHistory{}
	requestRoomBaseInfo := pb.RoomBaseInfo{}
	//requestPlayerMatchSign := pb.GameMatchSign{}
	//requestPlayerMatchSignCancel := pb.GameMatchPlayerCancel{}

	reqPlayerBase := pb.PlayerBase{}

	reqFDCreateEvent := pb.FDCreateEvent{}
	reqHeartbeat := pb.Heartbeat{}
	requestFDCloseEvent := pb.FDCloseEvent{}

	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	switch protoServiceFunc.FuncName {
	//case "CS_PlayerMatchSign":
	//	err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSign, conn.UserId)
	//case "CS_PlayerMatchSignCancel":
	//	err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSignCancel, conn.UserId)
	case "CS_PlayerOperations":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestLogicFrame, conn.UserId)
	case "CS_PlayerResumeGame":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerResumeGame, conn.UserId)
	case "CS_PlayerReady":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerReady, conn.UserId)
	case "CS_PlayerOver":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerOver, conn.UserId)
	case "CS_RoomHistory":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestRoomHistory, conn.UserId)
	case "CS_RoomBaseInfo":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestRoomBaseInfo, conn.UserId)
	case "FdClose":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestFDCloseEvent, conn.UserId)
	case "CS_Heartbeat":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqHeartbeat, conn.UserId)
	case "CS_PlayerState":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqPlayerBase, conn.UserId)
	case "FdCreate":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqFDCreateEvent, conn.UserId)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "1"))
	}
	if err != nil {
		gateway.Log.Error(prefix + " , ParserContentMsg err:" + err.Error())
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerOperations":
		requestLogicFrame.SourceUid = conn.UserId
		err = gateway.MyServiceList.FrameSync.ReceivePlayerOperation(requestLogicFrame)
	case "CS_PlayerResumeGame":
		requestPlayerResumeGame.SourceUid = conn.UserId
		err = gateway.MyServiceList.FrameSync.PlayerResumeGame(requestPlayerResumeGame)
	case "CS_PlayerReady":
		requestPlayerReady.SourceUid = conn.UserId
		err = gateway.MyServiceList.FrameSync.PlayerReady(requestPlayerReady)
	case "CS_PlayerOver":
		requestPlayerOver.SourceUid = conn.UserId
		err = gateway.MyServiceList.FrameSync.PlayerOver(requestPlayerOver)
	case "CS_RoomHistory":
		requestRoomHistory.SourceUid = conn.UserId
		err = gateway.MyServiceList.FrameSync.RoomHistory(requestRoomHistory)
	case "CS_RoomBaseInfo":
		requestRoomBaseInfo.SourceUid = conn.UserId
		err = gateway.MyServiceList.FrameSync.RoomManage.GetRoom(requestRoomBaseInfo)
	case "CS_PlayerState":
		gateway.MyServiceList.FrameSync.GetPlayerBase(reqPlayerBase)
	case "FdClose":
		err = gateway.MyServiceList.FrameSync.CloseFD(requestFDCloseEvent)
	case "CS_Heartbeat":
		err = gateway.MyServiceList.FrameSync.Heartbeat(reqHeartbeat)
	case "FdCreate":
		err = gateway.MyServiceList.FrameSync.CreateFD(reqFDCreateEvent)
	default:
		return data, errors.New(gateway.MakeRouterErrNotFound(prefix, protoServiceFunc.FuncName, "2"))
	}

	return data, err
}

func (gateway *Gateway) ClientPong(requestClientPong pb.PongRes, conn *util.Conn) {
	gateway.Log.Debug("ClientPong")
}

func (gateway *Gateway) heartbeat(requestClientHeartbeat pb.Heartbeat, conn *util.Conn) {
	now := util.GetNowTimeSecondToInt()
	now64 := util.GetNowMillisecond()
	conn.UpTime = int32(now)
	conn.RTT = now64 - requestClientHeartbeat.Time

	gateway.Log.Debug("gateway heartbeat , now64:", zap.Int64("now", now64), zap.Int64(" client_time", requestClientHeartbeat.Time), zap.Int64(" RTT:", conn.RTT))
	responseHeartbeat := pb.Heartbeat{
		Time:              now64,
		ReqTime:           requestClientHeartbeat.ClientReqTime,
		ClientReqTime:     requestClientHeartbeat.ClientReqTime,
		ServerReceiveTime: now64,
		RequestId:         requestClientHeartbeat.RequestId,
	}

	conn.SendMsgCompressByUid(conn.UserId, "SC_Heartbeat", &responseHeartbeat)
}

func (gateway *Gateway) clientPing(ping pb.PingReq) {
	responseServerPong := pb.PongRes{
		ClientReqTime:      ping.ClientReqTime,
		ClientReceiveTime:  ping.ClientReceiveTime,
		ServerReceiveTime:  util.GetNowMillisecond(),
		ServerResponseTime: util.GetNowMillisecond(),
	}
	gateway.RequestServiceAdapter.GatewaySendMsgByUid(ping.SourceUid, "SC_Pong", &responseServerPong)
}

//balanceFactor:负载均衡 方法
func (gateway *Gateway) HttpCallGrpc(serviceName string, funcName string, balanceFactor string, requestData []byte) (resJsonStr string, err error) {
	gateway.Log.Info("HttpCallGrpc ， serviceName:" + serviceName + " funcName:" + funcName + " balanceFactor:" + balanceFactor + " requestData:" + string(requestData))
	callGrpcResData, err := gateway.GrpcManager.CallGrpc(serviceName, funcName, balanceFactor, requestData)
	if err != nil {
		return resJsonStr, err
	}
	resJsonStrByte, err := json.Marshal(callGrpcResData)
	if err != nil {
		return resJsonStr, err
	}
	return string(resJsonStrByte), err
}
