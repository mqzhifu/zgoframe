package test

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

func Gateway(){
	GateServer()
	//GateClientWebsocket()
}

func GetGatewayInstance(){

}

func GateServer(){
	//a := int32(1)
	//ab := byte(a)
	//util.ExitPrint(ab)

	netWayOption := util.NetWayOption{
		ListenIp			: "127.0.0.1",	//程序启动时监听的IP
		OutIp				: "127.0.0.1",	//对外访问的IP

		WsPort 				: "1111",		//监听端口号
		TcpPort 			: "2222",		//监听端口号
		UdpPort				: "3333",		//UDP端口号

		WsUri				: "/ws",			//接HOST的后面的URL地址
		Protocol 			:util.PROTOCOL_WEBSOCKET,		 	//兼容协议：ws tcp udp
		ContentType 		: util.CONTENT_TYPE_PROTOBUF,	//默认内容格式 ：json protobuf

		LoginAuthType		: "/jwt",	//jwt
		LoginAuthSecretKey	: "aaaa",	//密钥

		MaxClientConnNum	: 10,		//客户端最大连接数
		MsgContentMax		: 10240,		//一条消息内容最大值
		IOTimeout			: 1,				//read write sock fd 超时时间
		ConnTimeout 		: 60,			//一个FD超时时间
		GrpcManager			: global.V.GrpcManager,
		Log 				: global.V.Zap,
		ProtobufMap			: global.V.ProtobufMap,
		//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
		//两种快速关闭方式，也可以直接调用shutdown函数
		//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
		//CloseChan 			chan int		`json:"-"`
	}

	gateway := util.NewGateway(global.V.GrpcManager,global.V.Zap)
	//global.V.Gateway = gateway
	gateway.StartSocket(netWayOption)


}

func GateClientWebsocket(){
	dns := "127.0.0.1:1111"
	u := url.URL{Scheme: "ws", Host: dns, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		global.V.Zap.Fatal("dial:" + err.Error())
	}
	//defer c.Close()

	actionName := "ClientLogin"

	requestLogin := pb.RequestLogin{
		Token: "aaaa",
	}
	requestLoginMarshal ,err := proto.Marshal(&requestLogin)
	if err != nil{
		global.V.Zap.Fatal("proto.Marshal err:" + err.Error())
	}

	actionMap ,empty:= global.V.ProtobufMap.GetActionId(actionName)
	if empty{
		global.V.Zap.Panic("GetActionId empty.")
	}

	protocol 			:= util.PROTOCOL_WEBSOCKET
	contentType 		:= util.CONTENT_TYPE_PROTOBUF

	protocolManagerOption := util.ProtocolManagerOption {
		Log: global.V.Zap,
	}
	protocolManager := util.NewProtocolManager(protocolManagerOption)

	msg := pb.Msg{
		ContentType: int32(contentType),
		ProtocolType: int32(protocol),
		Action: actionName,
		ActionId: int32(actionMap.Id),
		ServiceId:int32( actionMap.ServiceId),
		Content:string(requestLoginMarshal),
	}

	contentBytes := protocolManager.PackContentMsg(msg)
	util.MyPrint(contentBytes)
	err = c.WriteMessage(websocket.BinaryMessage,contentBytes)
	if err != nil {
		global.V.Zap.Error("write:"+err.Error())
		return
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			global.V.Zap.Error("read:"+err.Error())
			return
		}
		global.V.Zap.Info("recv:"+string(message))
		time.Sleep(time.Second * 1)
	}

}

func GateClientTcp(){

}
