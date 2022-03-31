package util

import (
	"errors"
	"zgoframe/protobuf/pb"
)

//总路由器，这里分成了两类：gateway 自解析 和 代理后方的请求
func (netWay *NetWay) Router(msg pb.Msg, conn *Conn) (data interface{}, err error) {
	//serviceActionIds, _ := strconv.Atoi(strconv.Itoa(int(msg.ServiceId)) + strconv.Itoa(int(msg.ActionId)))
	actionInfo, _ := netWay.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	serviceName := actionInfo.ServiceName
	//MyPrint("Router:", actionInfo, "serviceActionIds:", serviceActionIds, " msg info:", msg)
	switch serviceName {
	case "Gateway":
		data, err = netWay.RouterServiceGateway(msg, conn)
	case "FrameSync":
		data, err = netWay.RouterServiceSync(msg, conn, actionInfo)
	default:
		netWay.Option.Log.Error("netWay Router err.")
	}
	return data, err
}

//帧同步的路由
func (netWay *NetWay) RouterServiceSync(msg pb.Msg, conn *Conn, actionMap ProtoServiceFunc) (data []byte, err error) {
	//zgoframeClient, err := netWay.Option.GrpcManager.GetZgoframeClient(actionMap.ServiceName, strconv.Itoa(int(conn.UserId)))
	//ctx := context.Background()
	//switch msg.Action {
	//case "SayHello": //
	//	requestUser := pb.RequestUser{}
	//	proto.Unmarshal([]byte(msg.Content), &requestUser)
	//	//*ResponseUser
	//	dataClass, err := zgoframeClient.SayHello(ctx, &requestUser)
	//	if err != nil {
	//		return data, err
	//	}
	//	data, err = proto.Marshal(dataClass)
	//default:
	//
	//}

	return data, err
}

//网关自解析的路由
func (netWay *NetWay) RouterServiceGateway(msg pb.Msg, conn *Conn) (data interface{}, err error) {
	requestLogin := pb.Login{}
	requestClientPong := pb.Pong{}
	requestClientPing := pb.Ping{}
	requestClientHeartbeat := pb.Heartbeat{}

	protoServiceFunc, _ := netWay.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		err = netWay.ProtocolManager.parserContentMsg(msg, &requestLogin, conn.UserId)
	case "CS_Pong": //
		err = netWay.ProtocolManager.parserContentMsg(msg, &requestClientPong, conn.UserId)
	case "CS_Ping":
		err = netWay.ProtocolManager.parserContentMsg(msg, &requestClientPing, conn.UserId)
	case "CS_Heartbeat": //心跳
		err = netWay.ProtocolManager.parserContentMsg(msg, &requestClientHeartbeat, conn.UserId)
	default:
		netWay.Option.Log.Error("RouterServiceGateway Router err:")
		return data, errors.New("RouterServiceGateway Router err")
	}
	if err != nil {
		return data, err
	}
	netWay.Option.Log.Info("Router " + protoServiceFunc.FuncName)
	switch protoServiceFunc.FuncName {
	case "CS_Login": //
		//这里有个BUG，LOGIN 函数只能在第一次调用，回头加个限定
		data, err = netWay.login(requestLogin, conn)
	case "CS_Ping":
		netWay.clientPing(requestClientPing, conn)
	case "CS_Pong":
		netWay.ClientPong(requestClientPong, conn)
	case "CS_Heartbeat":
		netWay.heartbeat(requestClientHeartbeat, conn)
	}
	return data, err
}

func (netWay *NetWay) ClientPong(requestClientPong pb.Pong, conn *Conn) {

}

func (netWay *NetWay) clientPing(pingRTT pb.Ping, conn *Conn) {
	responseServerPong := pb.Pong{
		AddTime:            pingRTT.AddTime,
		ClientReceiveTime:  pingRTT.ClientReceiveTime,
		ServerResponseTime: GetNowMillisecond(),
		RttTimes:           pingRTT.RttTimes,
		RttTimeout:         pingRTT.RttTimeout,
	}
	conn.SendMsgCompressByUid(conn.UserId, "SC_Pong", &responseServerPong)
}
