package util

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/http/request"
	"zgoframe/protobuf/pb"
)

type NetWayOption struct {
	ListenIp  string `json:"listenIp"` //程序启动时监听的IP
	OutIp     string `json:"outIp"`    //对外访问的IP
	OutDomain string `json:"outDomain"`
	WsPort    string `json:"wsPort"`  //ws监听端口号
	TcpPort   string `json:"tcpPort"` //tcp监听端口号
	UdpPort   string `json:"udpPort"` //udp端口号

	WsUri               string `json:"wsUri"`                 //ws接HOST的后面的URL地址
	DefaultProtocolType int32  `json:"default_protocol_type"` //默认响应协议：ws tcp udp
	DefaultContentType  int32  `json:"default_content_type"`  //默认响应内容格式 ：json protobuf

	LoginAuthType      string `json:"loginAuthType"`         //jwt登陆验证
	LoginAuthSecretKey string `json:"login_auth_secret_key"` //jwt登陆验证-密钥

	MaxClientConnNum int32 `json:"maxClientConnMum"` //客户端最大连接数
	MsgContentMax    int32 `json:"msg_content_max"`  //一条消息内容最大值,byte,ps:最大10KB
	IOTimeout        int64 `json:"io_timeout"`       //read write sock fd 超时时间
	ConnTimeout      int32 `json:"connTimeout"`      //一个FD超时时间

	GrpcManager *GrpcManager `json:"-"` //外部链接,grpc反代
	ProtoMap    *ProtoMap    `json:"-"` //外部链接,协议号转换
	Log         *zap.Logger  `json:"-"` //外部链接,日志

	//ProtobufMapPath		string		`json:"portobuf_map_path"`	//协议号对应的函数名

	//两种关闭方式：
	//OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
	//CloseChan 			chan int		`json:"-"`			//调用方的CTX，用于所有协程的退出操作
	//http 查看状态
	//HttpdRootPath 	string 		`json:"httpdRootPath"`
	//HttpPort 			string 		`json:"httpPort"`		//短连接端口号

	//以下都是帧同步的配置信息：
	MapSize int32 `json:"mapSize"` //地址大小，给前端初始化使用，测试使用
	//RoomPeople		int32		`json:"roomPeople"`		//一局游戏包含几个玩家
	//RoomTimeout 		int32 		`json:"roomTimeout"`	//一个房间超时时间
	OffLineWaitTime int32 `json:"offLineWaitTime"` //lockStep 玩家掉线后，其它玩家等待最长时间
	//LockMode  		int32 		`json:"lockMode"`		//锁模式，乐观|悲观
	FPS              int32 `json:"fps"`              //frame pre second
	RoomReadyTimeout int32 `json:"roomReadyTimeout"` //一个房间的，玩家的准备，超时时间
	//Store 			int32 		`json:"store"`			//持久化：players room

	RouterBack func(msg pb.Msg, conn *Conn) (data interface{}, err error) `json:"-"`
}

var myMetrics *MyMetrics

type NetWay struct {
	//CancelCtx         	context.Context
	//CancelFunc			func()
	//CloseChan       	chan int32
	Status int
	Prefix string
	//MetricsPool     []MyMetricsPoolItem
	ProtocolManager *ProtocolManager //协议管理器
	ConnManager     *ConnManager     //连接管理 器
	Metrics         *MyMetrics       //metric管理 器
	ProtoMap        *ProtoMap        //protoBuf 管理器

	Option NetWayOption
}

//var myNetWay *NetWay//快捷变量，回头干掉
func NewNetWay(option NetWayOption) (*NetWay, error) {
	option.Log.Info("New NetWay instance :")

	netWay := new(NetWay)
	netWay.Prefix = "netway "
	//统计模块
	netWay.Metrics = netWay.InitMetrics(option.Log)
	myMetrics = netWay.Metrics
	//单条消息最大值
	netWay.Option = option
	if option.MsgContentMax > 10240 {
		option.MsgContentMax = 10240 //最大10KB
	}
	//设置状态为：初始化
	netWay.Status = NETWAY_STATUS_INIT

	netWay.ProtoMap = option.ProtoMap
	//协议管理适配器
	protocolManagerOption := ProtocolManagerOption{
		Ip:              option.ListenIp,
		WsPort:          option.WsPort,
		TcpPort:         option.TcpPort,
		WsUri:           option.WsUri,
		UdpPort:         option.UdpPort,
		IOTimeout:       option.IOTimeout,
		OpenNewConnBack: netWay.OpenNewConn, //回调函数
		Log:             option.Log,
	}
	netWay.ProtocolManager = NewProtocolManager(protocolManagerOption)
	err := netWay.ProtocolManager.Start()
	if err != nil {
		return nil, err
	}

	//player.fd.num
	//player.fd.size

	//http 工具
	//httpdOption := HttpdOption {
	//	LogOption 	: netWay.Option.LogOption,
	//	RootPath 	: netWay.Option.HttpdRootPath,
	//	Ip			: netWay.Option.ListenIp,
	//	Port		: netWay.Option.HttpPort,
	//	ParentCtx 	: option.OutCxt,
	//}
	//myHttpd = NewHttpd(httpdOption)

	//conn FD 管理
	connManagerOption := ConnManagerOption{
		maxClientConnNum:    option.MaxClientConnNum,
		connTimeout:         option.ConnTimeout,
		Log:                 option.Log,
		DefaultProtocolType: netWay.Option.DefaultProtocolType,
		DefaultContentType:  netWay.Option.DefaultContentType,
		Metrics:             netWay.Metrics,
		ProtoMap:            option.ProtoMap,
		NetWay:              netWay,
		MsgContentMax:       option.MsgContentMax,
	}
	netWay.ConnManager = NewConnManager(connManagerOption)
	go netWay.ConnManager.CheckTimeout()

	netWay.Status = NETWAY_STATUS_START

	option.Log.Info("netway startup finish.")
	return netWay, nil
}
func (netWay *NetWay) InitMetrics(log *zap.Logger) *MyMetrics {

	metrics := NewMyMetrics(MyMetricsOption{Log: log, DescPrefix: netWay.Prefix})

	metrics.CreateGauge("startup_time", "启动时间") //启动时间

	metrics.CreateCounter("ws_ok_fd", "websocket 成功建立FD 数量")           //websocket 成功建立FD 数量
	metrics.CreateCounter("ws_server_close_fd", "websocket 主动关闭FD 数量") //websocket 主动关闭FD 数量
	metrics.CreateCounter("ws_client_close_fd", "websocket 被动关闭FD 数量") //websocket 被动关闭FD 数量
	metrics.CreateCounter("tcp_ok_fd", "tcp 成功建立FD 数量")                //tcp 成功建立FD 数量
	metrics.CreateCounter("tcp_server_close_fd", "tcp 主动关闭FD 数量")      //tcp 主动关闭FD 数量
	metrics.CreateCounter("tcp_client_close_fd", "tcp 被动关闭FD 数量")      //tcp 被动关闭FD 数量
	//以上均是 最底层 TCP WS  的统计信息

	//以下有点偏向应用层的统计
	metrics.CreateCounter("new_fd", "接收来自 tcp/ws 新FD 数量") //netway 接收来自 tcp/ws 新FD 数量

	metrics.CreateCounter("create_fd_ok", "验证通过，成功创建的FD") //验证通过，成功创建的FD
	metrics.CreateCounter("create_fd_failed", "验证失败，FD")  //验证失败，FD
	metrics.CreateCounter("server_close_fd", "主动关闭FD")    //主动关闭FD
	metrics.CreateCounter("client_close_fd", "被动关闭FD")    //被动关闭FD

	metrics.CreateCounter("total_output_num", "总发送消息 次数") //总发送消息 次数
	metrics.CreateGauge("total_output_size", "总发送消息 大小")  //总发送消息 大小
	metrics.CreateCounter("total_input_num", "总接收消息 次数")  //总接收消息 次数
	metrics.CreateGauge("total_input_size", "总接收消息 大小")   //总接收消息 大小

	now := GetNowTimeSecondToInt64()
	metrics.GaugeSet("startup_time", float64(now))

	return metrics
}

//func (netWay *NetWay) MetricsTotal() {
//	item := netWay.MetricsPool[0]
//}

//一个新客户端连接请求进入
func (netWay *NetWay) OpenNewConn(connFD FDAdapter) {
	myMetrics.CounterInc("ws_ok_fd")
	netWay.Option.Log.Info("OpenNewConn:" + connFD.RemoteAddr())

	var loginRes pb.LoginRes

	if netWay.Status == NETWAY_STATUS_CLOSE { //当前网关已经关闭了，还有新的连接进来
		//记录：创建FD失败次数
		netWay.Metrics.CounterInc("create_fd_failed")
		errMsg := "netWay closing... not accept new connect , sleep 1!"
		netWay.Option.Log.Error(errMsg)
		//直接给一个FD发送消息，不做任何封装
		netWay.WriteMessage(int(netWay.Option.DefaultContentType), connFD, []byte(errMsg))
		//这里暂停一会，保证上面的消息能发送成功
		time.Sleep(time.Millisecond * 200)
		//直接关闭一个FD，不做任何多余处理
		connFD.Close()
		return
	}
	//是否超过了，最大可连接数
	if int32(len(netWay.ConnManager.Pool)) > netWay.Option.MaxClientConnNum {
		netWay.Metrics.CounterInc("create_fd_failed")

		errMsg := "more MaxClientConnNum"
		netWay.Option.Log.Error(errMsg)
		netWay.WriteMessage(int(netWay.Option.DefaultContentType), connFD, []byte(errMsg))
		connFD.Close()
		return
	}
	//创建一个新的连接结体体，将 FD 保存到该容器中
	NewConn := netWay.ConnManager.CreateOneConn(connFD, netWay)
	//defer func() {
	//	if err := recover(); err != nil {
	//		netWay.Option.Log.Panic("OpenNewConn:")
	//		netWay.CloseOneConn(NewConn, CLOSE_SOURCE_OPEN_PANIC)
	//	}
	//}()
	//开始-登陆验证
	jwtData, firstMsg, err := netWay.loginPre(NewConn)
	if err != nil {
		//这里不用发消息了，也不用关闭FD，因为loginPre内部已经处理过了
		return
	}
	//登陆验证已通过，开始添加各种状态及初始化
	NewConn.UserId = int32(jwtData.Id)
	//将新的连接加入到连接池中，并且与玩家ID绑定
	netWay.ConnManager.addConnPool(NewConn)
	//if err != nil{
	//	loginRes = pb.ResponseLoginRes{
	//		Code: 500,
	//		ErrMsg: err.Error(),
	//	}
	//	netWay.SendMsgCompressByUid(jwtData.Payload.Uid,"loginRes",&loginRes)
	//	netWay.CloseOneConn(NewConn, CLOSE_SOURCE_OVERRIDE)
	//	return
	//}
	//更新当前连接的属性值
	NewConn.ProtocolType = firstMsg.ProtocolType
	NewConn.ContentType = firstMsg.ContentType
	loginRes = pb.LoginRes{
		Code:   200,
		ErrMsg: "",
		Uid:    NewConn.UserId,
	}
	//告知玩家：登陆结果
	NewConn.SendMsgCompressByUid(NewConn.UserId, "SC_Login", &loginRes)
	//统计 当前FD 数量/历史FD数量
	netWay.Metrics.CounterInc("create_fd_ok")

	//具体的执行过程，要走一遍gateway 的router ,开始：登陆/验证 过程
	netWay.Router(firstMsg, NewConn)

	//初始化即登陆成功的响应均完成后，开始该连接的 消息IO 协程
	go NewConn.IOLoop()
	//netWay.serverPingRtt(time.Duration(rttMinTimeSecond),NewWsConn,1)
	netWay.Option.Log.Info("wsHandler end ,player login success!!!")

}

func (netWay *NetWay) heartbeat(requestClientHeartbeat pb.Heartbeat, conn *Conn) {
	now := GetNowTimeSecondToInt()
	conn.UpTime = int32(now)

	responseHeartbeat := pb.Heartbeat{
		Time: int64(now),
	}

	conn.SendMsgCompressByUid(conn.UserId, "SC_Headerbeat", &responseHeartbeat)
}

//=================================
//直接关闭一个FD，主要用于：登陆就失败了的情况
func (netWay *NetWay) CloseFD(connFD FDAdapter, source int) {
	connFD.Close()
	//记录主动关闭FD次数
	netWay.Metrics.CounterInc("close_fd_num")
}

//退出，目前能直接调用此函数的，就只有一种情况：
//MAIN 接收到了中断信号，并执行了：context-cancel()，然后，startup函数的守护监听到，调用些方法
func (netWay *NetWay) Shutdown() {
	netWay.Option.Log.Warn("netWay.Shutdown")
	if netWay.Status != NETWAY_STATUS_START {
		netWay.Option.Log.Error("Quit err :netWay.Status !=  NETWAY_STATUS_START")
		return
	}
	netWay.Status = NETWAY_STATUS_CLOSE //更新状态为：关闭

	//myHttpd.shutdown()
	//myMatch.Shutdown()
	//mySync.Shutdown()
	//myPlayerManager.Shutdown()
	netWay.ConnManager.Shutdown()
	netWay.ProtocolManager.Shutdown()
	netWay.Metrics.Shutdown()
	//go netWay.PlayerManager.checkOfflineTimeout(startupCtx)
	//netWay.Option.OutCancelFunc()
}

func (netWay *NetWay) loginPreFailedSendMsg(msg string, closeSource int, conn *Conn) {
	code := 500
	loginRes := pb.LoginRes{
		Code:   500,
		ErrMsg: msg,
	}
	netWay.Metrics.CounterInc("create_fd_failed")
	netWay.Option.Log.Error("loginPreFailed:" + strconv.Itoa(code) + " " + msg)
	conn.SendMsgCompressByConn("SC_Login", &loginRes)
	netWay.Option.Log.Info("sleep Millisecond * 500 wait msg sending...")
	time.Sleep(time.Millisecond * 500) //这里休息半秒，保证普通消息先发出去，且前端正常收到，不然可能：<fd关闭>消息早于消息到达前端
	conn.CloseOneConn(closeSource)
	//netWay.Option.Log.Error(msg)
}

//首次建立连接，登陆验证，预处理
func (netWay *NetWay) loginPre(conn *Conn) (jwt request.CustomClaims, firstMsg pb.Msg, err error) {
	content, err := conn.Read() //先从socket FD中读取一次数据
	if err != nil {
		netWay.loginPreFailedSendMsg(err.Error(), CLOSE_SOURCE_FD_READ_EMPTY, conn)
		return jwt, firstMsg, errors.New("conn read err:" + err.Error())
	}
	msg, err := conn.ConnManager.ParserContentProtocol(content)
	if err != nil {
		netWay.loginPreFailedSendMsg(err.Error(), CLOSE_SOURCE_FD_PARSE_CONTENT, conn)
		return jwt, firstMsg, err
	}
	//这里可能有个极端问题，连接成功后，C端立刻就得发消息，FD 读取消息可能会出现延迟，因为READ是异步，可能第一时间没有读到C端发来的数据
	protoServiceFunc, _ := netWay.Option.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	if protoServiceFunc.FuncName != "CS_Login" { //进到这里，肯定是有新连接被创建且回调了公共函数
		netWay.loginPreFailedSendMsg("first msg must : action=CS_Login api!!", CLOSE_SOURCE_FIRST_NO_LOGIN, conn)
		return
	}
	//jwt, err := netWay.Login(requestLogin, conn)
	////具体的执行过程，要走一遍gateway 的router ,开始：登陆/验证 过程
	//jwtDataInterface, err := netWay.Router(msg, conn)
	requestLogin := pb.Login{}
	err = netWay.ProtocolManager.ParserContentMsg(msg, &requestLogin, conn.UserId)
	if err != nil {
		netWay.loginPreFailedSendMsg(err.Error(), CLOSE_SOURCE_FIRST_PARSER_LOGIN, conn)
		return jwt, firstMsg, err
	}

	jwt, err = netWay.Login(requestLogin, conn)
	if err != nil {
		netWay.loginPreFailedSendMsg(err.Error(), CLOSE_SOURCE_CONN_LOGIN_ROUTER_ERR, conn)
		return
	}
	//jwt = jwtDataInterface.(request.CustomClaims)
	//if err != nil {
	//	netWay.loginPreFailedSendMsg(err.Error(), CLOSE_SOURCE_AUTH_FAILED, conn)
	//	return jwt, firstMsg, err
	//}
	netWay.Option.Log.Info("login jwt auth ok~~")
	return jwt, msg, nil
}

//登陆验证token
func (netWay *NetWay) Login(requestLogin pb.Login, conn *Conn) (customClaims request.CustomClaims, err error) {
	netWay.Option.Log.Info("netWay Login , token:" + requestLogin.Token)
	if conn.UserId > 0 {
		msg := " don't repeat login." + strconv.Itoa(int(conn.UserId))
		netWay.Option.Log.Error(msg)
		return customClaims, errors.New(msg)
	}
	token := ""
	if netWay.Option.LoginAuthType == "jwt" {
		token = requestLogin.Token
		MyPrint("token:", token)
		tokenParseWithClaims, err := jwt.ParseWithClaims(token, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
			return []byte(netWay.Option.LoginAuthSecretKey), nil
		})
		if err != nil {
			netWay.Option.Log.Error("gateway_netway Login jwt.ParseWithClaims err:  " + err.Error())
			return customClaims, err
		}
		claims, _ := tokenParseWithClaims.Claims.(*request.CustomClaims)
		//jwtData, err := ParseJwtToken(netWay.Option.LoginAuthSecretKey, token)
		return *claims, err
	} else {
		errMsg := "LoginAuthType err"
		netWay.Option.Log.Error(errMsg)
		return customClaims, errors.New(errMsg)
	}

	return customClaims, err
}

//直接给一个FD发送消息，基本上不用，只是特殊报错的时候，直接使用
//transmissionType : 1字符 2二进制
func (netWay *NetWay) WriteMessage(transmissionType int, connFD FDAdapter, content []byte) {
	myMetrics.CounterInc("total_output_num")
	myMetrics.GaugeAdd("total_output_size", float64(len(content)))

	err := connFD.WriteMessage(transmissionType, content)
	if err != nil {
		netWay.Option.Log.Error("WriteMessage err:" + err.Error())
	}
}

func (netWay *NetWay) RecoverGoRoutine(back func(ctx context.Context), ctx context.Context, err interface{}) {
	//pc, file, lineNo, ok := runtime.Caller(3)
	//if !ok {
	//	netWay.Option.Log.Error("runtime.Caller ok is false :",ok)
	//}
	//funcName := runtime.FuncForPC(pc).Name()
	//netWay.Option.Log.Info(" RecoverGoRoutine  panic in defer  :"+ funcName + " "+file + " "+ strconv.Itoa(lineNo))
	//RecoverGoRoutineRetryTimesRWLock.RLock()
	//retryTimes , ok := RecoverGoRoutineRetryTimes[funcName]
	//RecoverGoRoutineRetryTimesRWLock.RUnlock()
	//if ok{
	//	if retryTimes > 3{
	//		mylog.Error("retry than max times")
	//		panic(err)
	//		return
	//	}else{
	//		RecoverGoRoutineRetryTimesRWLock.Lock()
	//		RecoverGoRoutineRetryTimes[funcName]++
	//		RecoverGoRoutineRetryTimesRWLock.Unlock()
	//		mylog.Info("RecoverGoRoutineRetryTimes = ",RecoverGoRoutineRetryTimes[funcName])
	//	}
	//}else{
	//	mylog.Info("RecoverGoRoutineRetryTimes = 1")
	//	RecoverGoRoutineRetryTimesRWLock.Lock()
	//	RecoverGoRoutineRetryTimes[funcName] = 1
	//	RecoverGoRoutineRetryTimesRWLock.Unlock()
	//}
	//go back(ctx)
}
