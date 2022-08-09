package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

type Gateway struct {
	GrpcManager  *util.GrpcManager
	Log          *zap.Logger
	NetWayOption util.NetWayOption
	Netway       *util.NetWay
	MyService    *Service
}

//网关，目前主要是分为2部分
//1. http 代理 grpc
//2. 长连接代理，这里才是重点
func NewGateway(grpcManager *util.GrpcManager, log *zap.Logger) *Gateway {
	gateway := new(Gateway)
	gateway.GrpcManager = grpcManager
	gateway.Log = log
	return gateway
}

func (gateway *Gateway) HttpCallGrpc(serviceName string, funcName string, balanceFactor string, requestData []byte) (resJsonStr string, err error) {
	fmt.Print("HttpCallGrpc :", serviceName, funcName, balanceFactor, requestData)
	//gateway.Log.Info("HttpCallGrpc:")
	callGrpcResData, err := gateway.GrpcManager.CallGrpc(serviceName, funcName, balanceFactor, requestData)
	if err != nil {
		return resJsonStr, err
	}
	resJsonStrByte, err := json.Marshal(callGrpcResData)
	if err != nil {
		return resJsonStr, err
	}
	return string(resJsonStrByte), err
	//return resJsonStr,err
}

func (gateway *Gateway) StartSocket(netWayOption util.NetWayOption) (*util.NetWay, error) {
	netWayOption.RouterBack = gateway.Router
	gateway.NetWayOption = netWayOption
	netWay, err := util.NewNetWay(netWayOption)
	gateway.Netway = netWay
	return netWay, err
	//if err != nil {
	//	//errMsg := "NewNetWay err:" + err.Error()
	//	return netWay, err
	//}
	//for {
	//	time.Sleep(time.Second * 1)
	//}
	//netWay.Shutdown()
	//
	//roomId := "aabbccdd"
	//ZgoframeClient ,err := gateway.GrpcManager.GetZgoframeClient(roomId)

}

//总路由器，这里分成了两类：gateway 自解析 和 代理后方的请求
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
	default:
		gateway.Netway.Option.Log.Error("netWay Router err.")
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
		gateway.Netway.Option.Log.Error("RouterServiceGateway Router err:")
		return data, errors.New("RouterServiceGateway Router err")
	}

	if err != nil {
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerMatchSign":
		gateway.MyService.Match.AddOnePlayer(requestPlayerMatchSign, conn)
	case "CS_PlayerMatchSignCancel":
		gateway.MyService.Match.CancelOnePlayer(requestPlayerMatchSignCancel, conn)

	default:
		gateway.Netway.Option.Log.Error("RouterServiceGateway Router err:")
		return data, errors.New("RouterServiceGateway Router err")
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
		gateway.Netway.Option.Log.Error("RouterServiceGateway Router err:")
		return data, errors.New("RouterServiceGateway Router err")
	}
	if err != nil {
		return data, err
	}

	switch protoServiceFunc.FuncName {
	case "CS_PlayerOperations":
		err = gateway.MyService.FrameSync.ReceivePlayerOperation(requestLogicFrame, conn)
	case "CS_PlayerResumeGame":
		err = gateway.MyService.FrameSync.PlayerResumeGame(requestPlayerResumeGame, conn)
	case "CS_PlayerReady":
		err = gateway.MyService.FrameSync.PlayerReady(requestPlayerReady, conn)
	case "CS_PlayerOver":
		err = gateway.MyService.FrameSync.PlayerOver(requestPlayerOver, conn)
	case "CS_RoomHistory":
		err = gateway.MyService.FrameSync.RoomHistory(requestRoomHistory, conn)
	case "CS_RoomBaseInfo":
		err = gateway.MyService.RoomManage.GetRoom(requestRoomBaseInfo, conn)
	default:
		gateway.Netway.Option.Log.Error("RouterServiceGateway Router err:")
		return data, errors.New("RouterServiceGateway Router err")
	}

	return data, err
}

//网关自解析的路由
func (gateway *Gateway) RouterServiceGateway(msg pb.Msg, conn *util.Conn) (data interface{}, err error) {
	requestLogin := pb.Login{}
	requestClientPong := pb.Pong{}
	requestClientPing := pb.Ping{}
	requestClientHeartbeat := pb.Heartbeat{}

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
		data, err = gateway.Netway.Login(requestLogin, conn)
	case "CS_Ping":
		gateway.clientPing(requestClientPing, conn)
	case "CS_Pong":
		gateway.ClientPong(requestClientPong, conn)
	case "CS_Heartbeat":
		gateway.heartbeat(requestClientHeartbeat, conn)
	}
	return data, err
}

func (gateway *Gateway) ClientPong(requestClientPong pb.Pong, conn *util.Conn) {

}

func (gateway *Gateway) heartbeat(requestClientHeartbeat pb.Heartbeat, conn *util.Conn) {
	now := util.GetNowTimeSecondToInt()
	conn.UpTime = int32(now)

	responseHeartbeat := pb.Heartbeat{
		Time: int64(now),
	}

	conn.SendMsgCompressByUid(conn.UserId, "SC_Headerbeat", &responseHeartbeat)
}

func (gateway *Gateway) clientPing(pingRTT pb.Ping, conn *util.Conn) {
	responseServerPong := pb.Pong{
		AddTime:            pingRTT.AddTime,
		ClientReceiveTime:  pingRTT.ClientReceiveTime,
		ServerResponseTime: util.GetNowMillisecond(),
		RttTimes:           pingRTT.RttTimes,
		RttTimeout:         pingRTT.RttTimeout,
	}
	conn.SendMsgCompressByUid(conn.UserId, "SC_Pong", &responseServerPong)
}
