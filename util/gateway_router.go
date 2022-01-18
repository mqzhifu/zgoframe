package util

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"strconv"
	"zgoframe/protobuf/pb"
)
//这里分成了两类：gateway 自解析 和 代理后方的请求
func(netWay *NetWay) Router(msg pb.Msg,conn *Conn)(data interface{},err error){
	actionInfo,_ := netWay.ProtobufMap.ActionMaps[int(msg.ActionId)]
	serviceName := actionInfo.ServiceName
	switch serviceName {
	case "Gateway":
		data ,err  = netWay.RouterServiceGateway(msg,conn)
	case "FrameSync":
		data ,err  = netWay.RouterServiceSync(msg,conn,actionInfo)
	default:

	}
	return data,nil
}
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
			netWay.Option.Log.Error("Router err:")
			return data, nil
	}
	if err != nil {
		return data, err
	}
	netWay.Option.Log.Info("Router " + msg.Action)
	switch msg.Action {
		case "ClientLogin": //
			data, err = netWay.login(requestLogin, conn)
			//return jwtData, err
		//case "clientPong": //
		//	//netWay.ClientPong(requestClientPong, conn)
		//case "clientHeartbeat": //心跳
		//	netWay.heartbeat(requestClientHeartbeat, conn)
		//case "clientPing": //
		//	netWay.clientPing(requestClientPing, conn)
	}
	return data,nil
}
//直接给一个FD发送消息，基本上不用，只是特殊报错的时候，直接使用
func(netWay *NetWay)WriteMessage(TextMessage int, connFD FDAdapter,content []byte){
	err := connFD.WriteMessage(websocket.BinaryMessage,content)
	if err != nil{
		netWay.Option.Log.Error( "WriteMessage err:"+err.Error() )
	}
}
//发送一条消息给一个玩家，根据conn，同时将消息内容进行编码与压缩
//大部分通信都是这个方法
func(netWay *NetWay)SendMsgCompressByConn(conn *Conn,actionName string , contentStruct interface{}){
	netWay.Option.Log.Info("SendMsgCompressByConn  actionName:"+actionName)
	//conn.UserId=0 时，由函数内部做兼容，主要是用来取content type ,protocol type
	contentByte ,err := netWay.ProtocolManager.CompressContent(contentStruct,conn.UserId)
	if err != nil{
		return
	}
	netWay.SendMsg(conn,actionName,contentByte)
}
//发送一条消息给一个玩家，根据UserId，同时将消息内容进行编码与压缩
func(netWay *NetWay)SendMsgCompressByUid(UserId int32,action string , contentStruct interface{}){
	netWay.Option.Log.Info("SendMsgCompressByUid UserId:"+strconv.Itoa(int(UserId))  +  " action:" + action)
	contentByte ,err := netWay.ProtocolManager.CompressContent(contentStruct,UserId)
	if err != nil{
		return
	}
	netWay.SendMsgByUid(UserId,action,contentByte)
}
//发送一条消息给一个玩家,根据UserId,且不做压缩处理
func(netWay *NetWay)SendMsgByUid(UserId int32,action string , content []byte){
	conn,ok := netWay.ConnManager.getConnPoolById(UserId)
	if !ok {
		netWay.Option.Log.Error("conn not in pool,maybe del.")
		return
	}
	netWay.SendMsg(conn,action,content)
}
//发送一条消息给一个玩家,根据UserId,且不做压缩处理
func(netWay *NetWay)SendMsgByConn(conn *Conn,action string , content []byte){
	netWay.SendMsg(conn,action,content)
}

func(netWay *NetWay)SendMsg(conn *Conn,action string,content []byte){
	//获取协议号结构体
	actionMap,empty := netWay.ProtobufMap.GetActionId(action)
	if empty{
		netWay.Option.Log.Error("GetActionId empty:"+action)
		return
	}
	netWay.Option.Log.Info("SendMsg , actionId:"+ strconv.Itoa(actionMap.Id )+ " , userId:" + strconv.Itoa( int(conn.UserId))  + " , actionName:" + action)

	if conn.Status == CONN_STATUS_CLOSE {
		netWay.Option.Log.Error("Conn status =CONN_STATUS_CLOSE.")
		return
	}

	protocolCtrlInfo := myNetWay.ConnManager.GetPlayerCtrlInfoById(conn.UserId)
	msg := pb.Msg{
		Content: string(content),
		ServiceId: int32(actionMap.ServiceId),
		ActionId: int32(actionMap.Id),
		Action: actionMap.Action,
		ContentType:protocolCtrlInfo.ContentType,
		ProtocolType: protocolCtrlInfo.ProtocolType,
	}
	contentBytes := netWay.ProtocolManager.PackContentMsg(msg)

	if protocolCtrlInfo.ContentType == CONTENT_TYPE_PROTOBUF {
		conn.Write(contentBytes,websocket.BinaryMessage)
	}else{
		conn.Write(contentBytes,websocket.TextMessage)
	}
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
	netWay.SendMsgCompressByUid(conn.UserId,"serverPong",&responseServerPong)
}