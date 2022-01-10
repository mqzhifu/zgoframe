package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func Gateway(){
	//netWayOption := util.NetWayOption{
	//	ListenIp		:"127.0.0.1",	//程序启动时监听的IP
	//	OutIp			:"127.0.0.1",	//对外访问的IP
	//
	//	WsPort 			:"1111",		//监听端口号
	//	TcpPort 		:"2222",		//监听端口号
	//	UdpPort			:"3333",		//UDP端口号
	//
	//	WsUri			:"/ws",			//接HOST的后面的URL地址
	//	//Protocol 		:	1,		 	//兼容协议：ws tcp udp
	//	ContentType 	:util.CONTENT_TYPE_PROTOBUF,	//默认内容格式 ：json protobuf
	//
	//	LoginAuthType		:"/jwt",	//jwt
	//	LoginAuthSecretKey	:"aaaa",	//密钥
	//
	//	MaxClientConnNum	:10,		//客户端最大连接数
	//	MsgContentMax		:10240,		//一条消息内容最大值
	//	IOTimeout		:1,				//read write sock fd 超时时间
	//	ConnTimeout 	:60,			//一个FD超时时间
	//	GrpcManager		: global.V.GrpcManager,
	//	Log 			:global.V.Zap,
	//	//ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名
	//	//两种快速关闭方式，也可以直接调用shutdown函数
	//	//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
	//	//CloseChan 			chan int		`json:"-"`
	//}

	gateway := util.NewGateway(global.V.GrpcManager,global.V.Zap)
	global.V.Gateway = gateway
	//gateway.StartSocket(netWayOption)
}
