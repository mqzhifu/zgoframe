package util

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"strconv"
	"zgoframe/protobuf/pb"
)
//总路由器，这里分成了两类：gateway 自解析 和 代理后方的请求
func(netWay *NetWay) Router(msg pb.Msg,conn *Conn)(data interface{},err error){
	actionInfo,_ := netWay.ProtobufMap.ActionMaps[int(msg.ActionId)]
	serviceName := actionInfo.ServiceName
	switch serviceName {
		case "Gateway":
			data ,err  = netWay.RouterServiceGateway(msg,conn)
		case "FrameSync":
			data ,err  = netWay.RouterServiceSync(msg,conn,actionInfo)
		default:
			netWay.Option.Log.Error("netWay Router err.")
	}
	return data,err
}
//帧同步的路由
func(netWay *NetWay) RouterServiceSync(msg pb.Msg,conn *Conn,actionMap ActionMap)(data []byte,err error){
	zgoframeClient ,err := netWay.Option.GrpcManager.GetZgoframeClient(actionMap.ServiceName,strconv.Itoa(int(conn.UserId)))
	ctx := context.Background()
	switch msg.Action {
	case "SayHello": //
		requestUser := pb.RequestUser{}
		proto.Unmarshal([]byte(msg.Content),&requestUser)
		//*ResponseUser
		dataClass,err := zgoframeClient.SayHello(ctx,&requestUser)
		if err != nil{
			return data,err
		}
		data ,err = proto.Marshal(dataClass)
	default:

	}

	return data,err
}
//网关自解析的路由
func(netWay *NetWay) RouterServiceGateway(msg pb.Msg,conn *Conn)(data interface{},err error){
	requestLogin := pb.RequestLogin{}
	//requestClientPong := pb.RequestClientPong{}
	//requestClientPing := pb.RequestClientPing{}
	//requestClientHeartbeat := pb.RequestClientHeartbeat{}
	//这里有个BUG，LOGIN 函数只能在第一次调用，回头加个限定
	switch msg.Action {
		case "ClientLogin": //
			err = netWay.ProtocolManager.parserContentMsg(msg, &requestLogin, conn.UserId)
		//case "clientPong": //
		//	err = netWay.ProtocolManager.parserContentMsg(msg, &requestClientPong, conn.UserId)
		//case "clientPing":
		//	err = netWay.ProtocolManager.parserContentMsg(msg, &requestClientPing, conn.UserId)
		//case "clientHeartbeat": //心跳
		//	err = netWay.ProtocolManager.parserContentMsg(msg, &requestClientHeartbeat, conn.UserId)
		default:
			netWay.Option.Log.Error("RouterServiceGateway Router err:")
			return data, errors.New("RouterServiceGateway Router err")
	}
	if err != nil {
		return data, err
	}
	netWay.Option.Log.Info("Router " + msg.Action)
	switch msg.Action {
		case "ClientLogin": //
			netWay.Option.Log.Info("requestLogin token:"+requestLogin.Token)
			data, err = netWay.login(requestLogin, conn)
			//return jwtData, err
		//case "clientPong": //
		//	//netWay.ClientPong(requestClientPong, conn)
		//case "clientHeartbeat": //心跳
		//	netWay.heartbeat(requestClientHeartbeat, conn)
		//case "clientPing": //
		//	netWay.clientPing(requestClientPing, conn)
	}
	return data,err
}

func(netWay *NetWay) ClientPong(requestClientPong pb.RequestClientPong,conn *Conn){

}

func(netWay *NetWay)clientPing(pingRTT pb.RequestClientPing,conn *Conn){
	responseServerPong := pb.ResponseServerPong{
		AddTime: pingRTT.AddTime,
		ClientReceiveTime :pingRTT.ClientReceiveTime,
		ServerResponseTime: GetNowMillisecond(),
		RttTimes: pingRTT.RttTimes,
		RttTimeout: pingRTT.RttTimeout,
	}
	conn.SendMsgCompressByUid(conn.UserId,"serverPong",&responseServerPong)
}