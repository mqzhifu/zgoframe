package test

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

func Gateway(){
	GateServer()
	//GateClientWebsocket()
	//GateClientTcp()
}

var GateListenIp 			= "127.0.0.1"
var GateWsPort 				= "1111"
var GateWsUri 				= "/ws"
var GateTcpPort 			= "2222"
var GateDefaultProtocol 	= int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType 	= int32(util.CONTENT_TYPE_PROTOBUF)

func GateServer(){
	netWayOption := util.NetWayOption{
		ListenIp			: GateListenIp,	//程序启动时监听的IP
		OutIp				: GateListenIp,	//对外访问的IP

		WsPort 				: GateWsPort,		//监听端口号
		TcpPort 			: GateTcpPort,		//监听端口号
		//UdpPort				: "3333",		//UDP端口号

		WsUri				: GateWsUri,				//接HOST的后面的URL地址
		DefaultProtocolType	: GateDefaultProtocol,		//兼容协议：ws tcp udp
		DefaultContentType	: GateDefaultContentType,	//默认内容格式 ：json protobuf

		LoginAuthType		: "/jwt",	//jwt
		LoginAuthSecretKey	: "aaaa",	//密钥

		MaxClientConnNum	: 10,		//客户端最大连接数
		MsgContentMax		: 10240,	//一条消息内容最大值
		IOTimeout			: 1,		//read write sock fd 超时时间
		ConnTimeout 		: 60,		//一个FD超时时间
		GrpcManager			: global.V.GrpcManager,
		Log 				: global.V.Zap,
		ProtobufMap			: global.V.ProtobufMap,
		//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
		//两种快速关闭方式，也可以直接调用shutdown函数
		//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
		//CloseChan 			chan int		`json:"-"`
	}
	gateway := util.NewGateway(global.V.GrpcManager,global.V.Zap)
	gateway.StartSocket(netWayOption)

}

func GetSendLoginMsg()[]byte{
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

	protocol 			:= GateDefaultProtocol
	contentType 		:= GateDefaultContentType

	connManager := GetConnManager()

	msg := pb.Msg{
		ContentType	: contentType,
		ProtocolType: protocol,
		Action		: actionName,
		ActionId	: int32(actionMap.Id),
		ServiceId	: int32( actionMap.ServiceId),
		Content		: string(requestLoginMarshal),
	}

	contentBytes := connManager.PackContentMsg(msg)
	util.MyPrint("contentBytes len:",len(contentBytes))
	//util.MyPrint(contentBytes)
	return contentBytes
}

func GetConnManager()*util.ConnManager{
	connManagerOption := util.ConnManagerOption {
		Log			: global.V.Zap,
		ProtobufMap	: global.V.ProtobufMap,
		DefaultContentType: GateDefaultProtocol,
		DefaultProtocolType: GateDefaultProtocol,
	}
	connManager := util.NewConnManager(connManagerOption)
	return connManager
}

func GateClientWebsocket(){
	dns := GateListenIp+":"+GateWsPort
	u := url.URL{Scheme: "ws", Host: dns, Path: GateWsUri}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		global.V.Zap.Fatal("dial:" + err.Error())
	}
	//defer c.Close()
	connManager := GetConnManager()

	contentBytes := GetSendLoginMsg()
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

		msg ,err := connManager.ParserContentProtocol(string(message))
		if err != nil{
			global.V.Zap.Error("ParserContentProtocol:"+err.Error())
			return
		}

		util.PrintStruct(msg,":")
		time.Sleep(time.Second * 1)
	}

}

func GateClientTcp(){
	dns := GateListenIp+":"+GateTcpPort
	fd , err := net.Dial("tcp", dns)
	if err != nil{
		global.V.Zap.Fatal("net.Listen err :"+err.Error())
	}


	contentBytes := GetSendLoginMsg()
	n ,err := fd.Write(contentBytes)
	if err != nil{
		global.V.Zap.Fatal("fd.write err :"+err.Error())
	}
	global.V.Zap.Info("write n:"+strconv.Itoa(n))

	input := bufio.NewReader(os.Stdin)
	for {
		bytes, _, err := input.ReadLine()
		if err != nil {
			global.V.Zap.Fatal("read line faild err:%v\n"+ err.Error())
		}
		str := string(bytes)
		util.MyPrint("read:",string(str))
		//n, err := conn.Write(bytes)
		//if err != nil {
		//	fmt.Printf("send data faild err:%v\n", err)
		//} else {
		//	fmt.Printf("send data length %d\n", n)
		//}
		time.Sleep(time.Second * 1)
	}
}
