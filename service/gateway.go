package service

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

//这是个快捷变量，目前所有代码均在一起，直接挂在这个变量上即可，后期所有服务分拆出去，网关没那么多附加功能此变量就没用了
type MyServiceList struct {
	Match      *Match
	FrameSync  *FrameSync
	RoomManage *RoomManager
	TwinAgora  *TwinAgora
}

type Gateway struct {
	GrpcManager   *util.GrpcManager //通过GRPC反射代理其它微服务
	Log           *zap.Logger       //日志
	Netway        *util.NetWay      //长连接公共类
	NetWayOption  util.NetWayOption //长连接公共类的初始化参数
	MyServiceList *MyServiceList    //快捷访问内部微服务
}

//网关，目前主要是分为2部分
//1. http 代理 grpc(中等)
//2. 长连接代理(重点)
//3. http 代理 http(鸡肋)
func NewGateway(grpcManager *util.GrpcManager, log *zap.Logger) *Gateway {
	gateway := new(Gateway)
	gateway.GrpcManager = grpcManager
	gateway.Log = log
	gateway.MyServiceList = &MyServiceList{}
	return gateway
}

//开启长连接监听
func (gateway *Gateway) StartSocket(netWayOption util.NetWayOption) (*util.NetWay, error) {
	netWayOption.RouterBack = gateway.Router //公共回调 路由器，用于给最底层的长连接公共类回调
	//创建底层长连接公共类
	gateway.NetWayOption = netWayOption
	netWay, err := util.NewNetWay(netWayOption)
	gateway.Netway = netWay
	//监听，长连接公共类 - FD Close 事件(只能被动监听)
	go gateway.ListenCloseEvent()
	//微服务内部无法直接发送消息(在没有conn的情况下)，回头我想想怎么处理
	gateway.MyServiceList.TwinAgora.ConnManager = gateway.Netway.ConnManager

	return netWay, err
}

//监听长连接 - 关闭事件(只能被动监听)
func (gateway *Gateway) ListenCloseEvent() {
	for {
		select {
		case connCloseEvent := <-gateway.Netway.ConnManager.CloseEventQueue:
			gateway.Log.Debug("ListenCloseEvent connCloseEvent:" + util.StructToJsonStr(connCloseEvent))
			//随便取一个conn，给到下层服务，因为：下层服务可能还要继续给其它人发消息
			msg := gateway.MakeMsgCloseEventInfo(connCloseEvent)
			gateway.BroadcastService(msg, nil)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

//网关自己创建一条长连接消息，发送给service
func (gateway *Gateway) MakeMsgCloseEventInfo(connCloseEvent pb.FDCloseEvent) pb.Msg {
	//msg := pb.Msg{}
	//msg.ServiceId = 90
	//msg.FuncId = 120
	//msg.SidFid = 90120
	//msg.ContentType = int32(connCloseEvent.ContentType)
	//msg.ProtocolType = int32(connCloseEvent.ProtocolType)
	//
	//FDCloseEvent := pb.FDCloseEvent{}
	//FDCloseEvent.UserId = connCloseEvent.UserId
	//FDCloseEvent.Source = int32(connCloseEvent.Source)
	//var reqContentStr string
	//if msg.ContentType == util.CONTENT_TYPE_PROTOBUF {
	//	requestClientHeartbeatStrByte, _ := proto.Marshal(&FDCloseEvent)
	//	reqContentStr = string(requestClientHeartbeatStrByte)
	//} else {
	//	requestClientHeartbeatStrByte, _ := json.Marshal(FDCloseEvent)
	//	reqContentStr = string(requestClientHeartbeatStrByte)
	//}
	//
	//msg.Content = reqContentStr

	//FDCloseEvent := pb.FDCloseEvent{}
	//FDCloseEvent.UserId = connCloseEvent.UserId
	//FDCloseEvent.Source = int32(connCloseEvent.Source)

	requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressNormalContent(connCloseEvent, int(connCloseEvent.ContentType))
	msg, _, _ := gateway.Netway.ConnManager.MakeMsgByActionName(connCloseEvent.UserId, "FDCloseEvent", requestClientHeartbeatStrByte)
	return msg
}

//网关自己创建一条长连接消息，发送给service
func (gateway *Gateway) MakeMsgHeartbeat(requestClientHeartbeat pb.Heartbeat, conn *util.Conn) pb.Msg {
	//msg := pb.Msg{}
	//msg.ServiceId = 90
	//msg.ContentType = conn.ContentType
	//msg.ProtocolType = conn.ProtocolType
	//msg.FuncId = 106
	//msg.SidFid = 90106
	//
	//var reqContentStr string
	//if msg.ContentType == util.CONTENT_TYPE_PROTOBUF {
	//	requestClientHeartbeatStrByte, _ := proto.Marshal(&requestClientHeartbeat)
	//	reqContentStr = string(requestClientHeartbeatStrByte)
	//} else {
	//	requestClientHeartbeatStrByte, _ := json.Marshal(requestClientHeartbeat)
	//	reqContentStr = string(requestClientHeartbeatStrByte)
	//}
	//
	//msg.Content = reqContentStr

	requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(requestClientHeartbeat, conn.UserId)
	msg, _, _ := gateway.Netway.ConnManager.MakeMsgByActionName(conn.UserId, "CS_Heartbeat", requestClientHeartbeatStrByte)

	return msg
}

//网关自己创建一条长连接消息，发送给service
func (gateway *Gateway) MakeMsgFDCreateEventInfo(FDCreateEvent pb.FDCreateEvent, conn *util.Conn) pb.Msg {
	//msg := pb.Msg{}
	//msg.ServiceId = 90
	//msg.ContentType = conn.ContentType
	//msg.ProtocolType = conn.ProtocolType
	//msg.FuncId = 122
	//msg.SidFid = 90122
	//
	//var reqContentStr string
	//if msg.ContentType == util.CONTENT_TYPE_PROTOBUF {
	//	requestClientHeartbeatStrByte, _ := proto.Marshal(&FDCreateEvent)
	//	reqContentStr = string(requestClientHeartbeatStrByte)
	//} else {
	//	requestClientHeartbeatStrByte, _ := json.Marshal(FDCreateEvent)
	//	reqContentStr = string(requestClientHeartbeatStrByte)
	//}
	//
	//msg.Content = reqContentStr

	requestClientHeartbeatStrByte, _ := gateway.Netway.ConnManager.CompressContent(FDCreateEvent, conn.UserId)
	msg, _, _ := gateway.Netway.ConnManager.MakeMsgByActionName(conn.UserId, "FdCreate", requestClientHeartbeatStrByte)

	return msg
}

//广播给所有服务，如：心跳 PING PONG 关闭事件
func (gateway *Gateway) BroadcastService(msg pb.Msg, conn *util.Conn) {
	//gateway.RouterServiceSync(msg, conn)
	//gateway.RouterServiceGameMatch(msg, conn)
	gateway.RouterServiceTwinAgora(msg, conn)
}

//总路由器，这里分成了两类：gateway 自解析 和 代理后方服务的请求
func (gateway *Gateway) Router(msg pb.Msg, conn *util.Conn) (data interface{}, err error) {
	actionInfo, _ := gateway.NetWayOption.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	gateway.Log.Info("service gateway router:" + actionInfo.ServiceName + " " + actionInfo.FuncName)
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
		gateway.Netway.Option.Log.Error("netWay Router err.")
	}
	return data, err
}

func (gateway *Gateway) RouterServiceTwinAgora(msg pb.Msg, conn *util.Conn) (data []byte, err error) {
	//util.MyPrint(msg, "-------RouterServiceTwinAgora msg.Content:", msg.Content)
	requestCallPeopleReq := pb.CallPeopleReq{}
	requestFDCloseEvent := pb.FDCloseEvent{}
	reqHeartbeat := pb.Heartbeat{}
	reqFDCreateEvent := pb.FDCreateEvent{}
	reqCallVote := pb.CallVote{}
	reqRoomHeartbeat := pb.RoomHeartbeatReq{}
	cancelCallPeopleReq := pb.CancelCallPeopleReq{}
	reqPeopleEntry := pb.PeopleEntry{}
	reqPeopleLeaveRes := pb.PeopleLeaveRes{}
	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	//util.MyPrint("RouterServiceTwinAgora protoServiceFunc:", protoServiceFunc)
	switch protoServiceFunc.FuncName {
	case "CS_CallPeople":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestCallPeopleReq, conn.UserId)
	case "CS_PeopleLeave":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &reqPeopleLeaveRes, conn.UserId)
	case "CS_CancelCallPeople":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &cancelCallPeopleReq, conn.UserId)
	case "FdClose":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestFDCloseEvent, 0)
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
		gateway.Netway.Option.Log.Error("RouterServiceTwinAgora Router err-1:")
		return data, errors.New("RouterServiceTwinAgora Router err-1")
	}

	if err != nil {
		util.MyPrint(err)
		//util.ExitPrint(err)
		return data, err
	}
	//util.MyPrint("=======2,", reqHeartbeat)
	switch protoServiceFunc.FuncName {
	case "CS_CallPeople":
		gateway.MyServiceList.TwinAgora.CallPeople(requestCallPeopleReq, conn)
	case "CS_Heartbeat":
		gateway.MyServiceList.TwinAgora.UserHeartbeat(reqHeartbeat, conn)
	case "CS_PeopleLeave":
		gateway.MyServiceList.TwinAgora.PeopleLeave(reqPeopleLeaveRes, conn)
	case "CS_RoomHeartbeat":
		gateway.MyServiceList.TwinAgora.RoomHeartbeat(reqRoomHeartbeat, conn)
	case "FdClose":
		gateway.MyServiceList.TwinAgora.FDCloseEvent(requestFDCloseEvent, gateway.Netway.ConnManager)
	case "CS_PeopleEntry":
		gateway.MyServiceList.TwinAgora.PeopleEntry(reqPeopleEntry, conn)
	case "CS_CancelCallPeople":
		gateway.MyServiceList.TwinAgora.CancelCallPeople(cancelCallPeopleReq, conn)
	case "FdCreate":
		gateway.MyServiceList.TwinAgora.FDCreateEvent(reqFDCreateEvent, conn)
	case "CS_CallPeopleAccept":
		gateway.MyServiceList.TwinAgora.CallPeopleAccept(reqCallVote, conn)
	case "CS_CallPeopleDeny":
		gateway.MyServiceList.TwinAgora.CallPeopleDeny(reqCallVote, conn)
	default:
		gateway.Netway.Option.Log.Error("RouterServiceTwinAgora Router err-2:")
		return data, errors.New("RouterServiceTwinAgora Router err-2")
	}

	return data, err
}

func (gateway *Gateway) RouterServiceGameMatch(msg pb.Msg, conn *util.Conn) (data []byte, err error) {
	requestPlayerMatchSign := pb.PlayerMatchSign{}
	requestPlayerMatchSignCancel := pb.PlayerMatchSignCancel{}
	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSign, conn.UserId)
	case "CS_PlayerMatchSignCancel":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSignCancel, conn.UserId)
	default:
		gateway.Netway.Option.Log.Error("RouterServiceGameMatch Router err:")
		return data, errors.New("RouterServiceGameMatch Router err")
	}

	if err != nil {
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		gateway.MyServiceList.Match.AddOnePlayer(requestPlayerMatchSign, conn)
	case "CS_PlayerMatchSignCancel":
		gateway.MyServiceList.Match.CancelOnePlayer(requestPlayerMatchSignCancel, conn)

	default:
		gateway.Netway.Option.Log.Error("RouterServiceGameMatch Router err:")
		return data, errors.New("RouterServiceGameMatch Router err")
	}

	return data, err
}

//帧同步的路由
func (gateway *Gateway) RouterServiceSync(msg pb.Msg, conn *util.Conn) (data []byte, err error) {
	requestLogicFrame := pb.LogicFrame{}
	requestPlayerResumeGame := pb.PlayerResumeGame{}
	requestPlayerReady := pb.PlayerReady{}
	requestPlayerOver := pb.PlayerOver{}
	requestRoomHistory := pb.ReqRoomHistory{}
	requestRoomBaseInfo := pb.RoomBaseInfo{}
	requestPlayerMatchSign := pb.PlayerMatchSign{}
	requestPlayerMatchSignCancel := pb.PlayerMatchSignCancel{}

	protoServiceFunc, _ := gateway.Netway.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	switch protoServiceFunc.FuncName {
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
	case "CS_PlayerMatchSign":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSign, conn.UserId)
	case "CS_PlayerMatchSignCancel":
		err = gateway.Netway.ProtocolManager.ParserContentMsg(msg, &requestPlayerMatchSignCancel, conn.UserId)
	default:
		gateway.Netway.Option.Log.Error("RouterServiceSync Router err:")
		return data, errors.New("RouterServiceSync Router err")
	}
	if err != nil {
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerOperations":
		err = gateway.MyServiceList.FrameSync.ReceivePlayerOperation(requestLogicFrame, conn)
	case "CS_PlayerResumeGame":
		err = gateway.MyServiceList.FrameSync.PlayerResumeGame(requestPlayerResumeGame, conn)
	case "CS_PlayerReady":
		err = gateway.MyServiceList.FrameSync.PlayerReady(requestPlayerReady, conn)
	case "CS_PlayerOver":
		err = gateway.MyServiceList.FrameSync.PlayerOver(requestPlayerOver, conn)
	case "CS_RoomHistory":
		err = gateway.MyServiceList.FrameSync.RoomHistory(requestRoomHistory, conn)
	case "CS_RoomBaseInfo":
		err = gateway.MyServiceList.RoomManage.GetRoom(requestRoomBaseInfo, conn)
	default:
		gateway.Netway.Option.Log.Error("RouterServiceSync Router err:")
		return data, errors.New("RouterServiceSync Router err")
	}

	return data, err
}

//网关自解析的路由
func (gateway *Gateway) RouterServiceGateway(msg pb.Msg, conn *util.Conn) (data interface{}, err error) {
	requestLogin := pb.Login{}
	requestClientPong := pb.PongRes{}
	requestClientPing := pb.PingReq{}
	requestClientHeartbeat := pb.Heartbeat{}
	requestProjectPushMsg := pb.ProjectPushMsg{}

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
		util.MyPrint("requestProjectPushMsg: ", requestProjectPushMsg)
		err = errors.New("CS_ProjectPushMsg no implementation。")
	default:
		gateway.Netway.Option.Log.Error("RouterServiceGateway Router err:")
		return data, errors.New("RouterServiceGateway Router err")
	}
	if err != nil {
		return data, err
	}
	gateway.Netway.Option.Log.Info("Router " + protoServiceFunc.FuncName)
	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		//这里有个BUG，LOGIN 函数只能在第一次调用，回头加个限定
		cc, err := gateway.Netway.Login(requestLogin, conn)
		data = cc
		if err == nil {
			FDCreateEvent := pb.FDCreateEvent{UserId: int32(cc.Id)}
			pbMsg := gateway.MakeMsgFDCreateEventInfo(FDCreateEvent, conn)
			gateway.BroadcastService(pbMsg, conn)
		}

	case "CS_Ping":
		gateway.clientPing(requestClientPing, conn)
	case "CS_Pong":
		gateway.ClientPong(requestClientPong, conn)
	case "CS_Heartbeat":

		gateway.heartbeat(requestClientHeartbeat, conn)
		msg := gateway.MakeMsgHeartbeat(requestClientHeartbeat, conn)
		gateway.BroadcastService(msg, conn)
	case "CS_ProjectPushMsg":
		util.MyPrint("RouterServiceGateway CS_ProjectPushMsg")
	}

	return data, err
}

func (gateway *Gateway) ClientPong(requestClientPong pb.PongRes, conn *util.Conn) {

}

func (gateway *Gateway) heartbeat(requestClientHeartbeat pb.Heartbeat, conn *util.Conn) {
	now := util.GetNowTimeSecondToInt()
	conn.UpTime = int32(now)

	//responseHeartbeat := pb.Heartbeat{
	//	Time: int64(now),
	//}

	//conn.SendMsgCompressByUid(conn.UserId, "SC_Heartbeat", &responseHeartbeat)
}

func (gateway *Gateway) clientPing(ping pb.PingReq, conn *util.Conn) {
	responseServerPong := pb.PongRes{
		ClientReqTime:      ping.ClientReqTime,
		ClientReceiveTime:  ping.ClientReceiveTime,
		ServerReceiveTime:  util.GetNowMillisecond(),
		ServerResponseTime: util.GetNowMillisecond(),
	}
	conn.SendMsgCompressByUid(conn.UserId, "SC_Pong", &responseServerPong)
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
