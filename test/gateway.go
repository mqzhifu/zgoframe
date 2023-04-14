package test

import (
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
	"zgoframe/core/global"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

func Gateway() {
	//GateServer()
	//GateClientWebsocket()
	util.MyPrint("test Gateway:===================")
	GateClientTcp()
}

var GateListenIp = "127.0.0.1"
var GateWsPort = "1111"
var GateWsUri = "/ws"
var GateTcpPort = "3333"
var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_PROTOBUF)

//func GateServer() {
//	netWayOption := util.NetWayOption{
//		ListenIp: GateListenIp, //程序启动时监听的IP
//		OutIp:    GateListenIp, //对外访问的IP
//
//		WsPort:  GateWsPort,  //监听端口号
//		TcpPort: GateTcpPort, //监听端口号
//		//UdpPort				: "3333",		//UDP端口号
//
//		WsUri:               GateWsUri,              //接HOST的后面的URL地址
//		DefaultProtocolType: GateDefaultProtocol,    //兼容协议：ws tcp udp
//		DefaultContentType:  GateDefaultContentType, //默认内容格式 ：json protobuf
//
//		LoginAuthType:      "/jwt", //jwt
//		LoginAuthSecretKey: "aaaa", //密钥
//
//		MaxClientConnNum: 10,    //客户端最大连接数
//		MsgContentMax:    10240, //一条消息内容最大值
//		IOTimeout:        1,     //read write sock fd 超时时间
//		ConnTimeout:      60,    //一个FD超时时间
//		GrpcManager:      global.V.GrpcManager,
//		Log:              global.V.Zap,
//		ProtobufMap:      global.V.ProtobufMap,
//		//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
//		//两种快速关闭方式，也可以直接调用shutdown函数
//		//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
//		//CloseChan 			chan int		`json:"-"`
//	}
//	gateway := util.NewGateway(global.V.GrpcManager, global.V.Zap)
//	gateway.StartSocket(netWayOption)
//
//}
//
func GetSendLoginMsg() []byte {
	funcName := "CS_Login"
	serviceName := "Gateway"

	requestLogin := pb.Login{
		Token: "aaaa",
	}
	requestLoginMarshal, err := proto.Marshal(&requestLogin)
	if err != nil {
		global.V.Zap.Fatal("proto.Marshal err:" + err.Error())
	}

	actionMap, empty := global.V.ProtoMap.GetServiceByName(serviceName, funcName)
	if empty {
		global.V.Zap.Panic("GetActionId empty.")
	}

	util.MyPrint("actionMap funcId:", actionMap.FuncId, " , serviceId:", actionMap.ServiceId)

	protocol := GateDefaultProtocol
	contentType := GateDefaultContentType

	connManager := GetConnManager()

	msg := pb.Msg{
		ContentType:  contentType,
		ProtocolType: protocol,
		ServiceId:    int32(actionMap.ServiceId),
		FuncId:       int32(actionMap.FuncId),
		Content:      string(requestLoginMarshal),
	}

	contentBytes := connManager.PackContentMsg(msg)
	util.MyPrint("contentBytes len:", len(contentBytes))
	//util.MyPrint(contentBytes)
	return contentBytes
}

func GetConnManager() *util.ConnManager {
	connManagerOption := util.ConnManagerOption{
		Log:                 global.V.Zap,
		DefaultContentType:  GateDefaultProtocol,
		DefaultProtocolType: GateDefaultProtocol,
	}
	connManager := util.NewConnManager(connManagerOption)
	return connManager
}

//
//func GateClientWebsocket() {
//	dns := GateListenIp + ":" + GateWsPort
//	u := url.URL{Scheme: "ws", Host: dns, Path: GateWsUri}
//	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//	if err != nil {
//		global.V.Zap.Fatal("dial:" + err.Error())
//	}
//	//defer c.Close()
//	connManager := GetConnManager()
//
//	contentBytes := GetSendLoginMsg()
//	err = c.WriteMessage(websocket.BinaryMessage, contentBytes)
//	if err != nil {
//		global.V.Zap.Error("write:" + err.Error())
//		return
//	}
//
//	for {
//		_, message, err := c.ReadMessage()
//		if err != nil {
//			global.V.Zap.Error("read:" + err.Error())
//			return
//		}
//
//		msg, err := connManager.ParserContentProtocol(string(message))
//		if err != nil {
//			global.V.Zap.Error("ParserContentProtocol:" + err.Error())
//			return
//		}
//
//		util.PrintStruct(msg, ":")
//		time.Sleep(time.Second * 1)
//	}
//
//}
//
func GateClientTcp() {
	dns := GateListenIp + ":" + GateTcpPort
	fd, err := net.Dial("tcp", dns)
	if err != nil {
		global.V.Zap.Fatal("net.Listen err :" + err.Error())
	}
	util.MyPrint("dns: ", dns)

	contentBytes := GetSendLoginMsg()
	util.MyPrint("fd.Write len:", len(contentBytes))
	n, err := fd.Write(contentBytes)
	if err != nil {
		global.V.Zap.Fatal("fd.write err :" + err.Error())
	}
	global.V.Zap.Info("write n:" + strconv.Itoa(n))

	//input := bufio.NewReader(os.Stdin)
	for {
		util.MyPrint("once....")
		//reader := bufio.NewReader(fd)
		//reader.Size()
		buf := make([]byte, 1024)
		r_len, err := fd.Read(buf)
		if err != nil {
			global.V.Zap.Fatal("read line faild err:%v\n" + err.Error())
		}

		//bytes, _, err := reader.ReadLine("\f")

		//str := string(bytes)
		util.MyPrint("loop read len:", r_len, "bytes:", buf[0:r_len], " ,str:", string(buf))

		time.Sleep(time.Millisecond * 100)
	}
}
