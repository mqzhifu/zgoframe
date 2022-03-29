package initialize

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

//var GateListenIp = "127.0.0.1"
//var GateWsPort = "1111"
//var GateWsUri = "/ws"
//var GateTcpPort = "2222"
var GateDefaultProtocol = int32(util.PROTOCOL_WEBSOCKET)
var GateDefaultContentType = int32(util.CONTENT_TYPE_PROTOBUF)

func InitGateway() (*util.Gateway, error) {
	netWayOption := util.NetWayOption{
		ListenIp: global.C.Gateway.ListenIp, //程序启动时监听的IP
		OutIp:    global.C.Gateway.OutIp,    //对外访问的IP

		WsPort:  global.C.Gateway.WsPort,  //监听端口号
		TcpPort: global.C.Gateway.TcpPort, //监听端口号
		//UdpPort				: "3333",		//UDP端口号

		WsUri:               global.C.Gateway.WsUri, //接HOST的后面的URL地址
		DefaultProtocolType: GateDefaultProtocol,    //兼容协议：ws tcp udp
		DefaultContentType:  GateDefaultContentType, //默认内容格式 ：json protobuf

		LoginAuthType:      "/jwt", //jwt
		LoginAuthSecretKey: "aaaa", //密钥

		MaxClientConnNum: 10,    //客户端最大连接数
		MsgContentMax:    10240, //一条消息内容最大值
		IOTimeout:        1,     //read write sock fd 超时时间
		ConnTimeout:      60,    //一个FD超时时间
		GrpcManager:      global.V.GrpcManager,
		Log:              global.V.Zap,
		ProtobufMap:      global.V.ProtobufMap,
		//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
		//两种快速关闭方式，也可以直接调用shutdown函数
		//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
		//CloseChan 			chan int		`json:"-"`
		FPS: 10,
	}
	gateway := util.NewGateway(global.V.GrpcManager, global.V.Zap)
	_, err := gateway.StartSocket(netWayOption)
	return gateway, err
}
