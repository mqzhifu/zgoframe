package util

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/protobuf/pb"
)

type NetWayOption struct {
	ListenIp			string		`json:"listenIp"`		//程序启动时监听的IP
	OutIp				string		`json:"outIp"`			//对外访问的IP

	WsPort 				string		`json:"wsPort"`			//监听端口号
	TcpPort 			string		`json:"tcpPort"`		//监听端口号
	UdpPort				string 		`json:"udpPort"`		//UDP端口号

	//Protocol 			int32		`json:"protocol"`		//兼容协议：ws tcp udp
	WsUri				string		`json:"wsUri"`			//接HOST的后面的URL地址
	ContentType 		int32		`json:"contentType"`	//默认内容格式 ：json protobuf

	LoginAuthType		string		`json:"loginAuthType"`	//jwt
	LoginAuthSecretKey	string		`json:"login_auth_secret_key"`//密钥

	MaxClientConnNum	int32		`json:"maxClientConnMum"`//客户端最大连接数
	MsgContentMax		int32		`json:"msg_content_max"`//一条消息内容最大值
	IOTimeout			int64		`json:"io_timeout"`		//read write sock fd 超时时间
	ConnTimeout 		int32		`json:"connTimeout"`	//一个FD超时时间
	ProtobufMapPath		string		`json:"portobuf_map_path"`//协议号对应的函数名

	Log 				*zap.Logger		`json:"-"`
	//两种快速关闭方式，也可以直接调用shutdown函数
	OutCxt 				context.Context `json:"-"`			//调用方的CTX，用于所有协程的退出操作
	CloseChan 			chan int		`json:"-"`
	GrpcManager			*GrpcManager
	//ProtobufMap			*ProtobufMap	`json:"-"`

	//HttpdRootPath 	string 		`json:"httpdRootPath"`
	//HttpPort 			string 		`json:"httpPort"`		//短连接端口号
	//MapSize			int32		`json:"mapSize"`		//地址大小，给前端初始化使用
	//RoomPeople		int32		`json:"roomPeople"`		//一局游戏包含几个玩家
	//RoomTimeout 		int32 		`json:"roomTimeout"`	//一个房间超时时间
	//OffLineWaitTime	int32		`json:"offLineWaitTime"`//lockStep 玩家掉线后，其它玩家等待最长时间
	//LockMode  		int32 		`json:"lockMode"`		//锁模式，乐观|悲观
	//FPS 				int32 		`json:"fps"`			//frame pre second
	//RoomReadyTimeout	int32		`json:"roomReadyTimeout"`//一个房间的，玩家的准备，超时时间
	//Store 			int32 		`json:"store"`			//持久化：players room
	//LogOption 		zlib.LogOption `json:"-"`
	//OutCancelFunc		context.CancelFunc `json:"-"`
}

type NetWay struct {
	MyCancelCtx         context.Context
	MyCancelFunc		func()
	CloseChan       	chan int32

	Status 				int

	ProtocolManager 	*ProtocolManager
	ConnManager			*ConnManager
	Metrics 			*MyMetrics
	ProtobufMap			*ProtobufMap
	Option          	NetWayOption
}


//var myProtocolManager  *ProtocolManager
//var connManager *ConnManager
var myNetWay *NetWay//快捷变量，回头干掉
func NewNetWay(option NetWayOption)(*NetWay,error)  {
	option.Log.Info("New NetWay instance :")

	netWay := new(NetWay)
	//mylog = option.Mylog
	myNetWay = netWay

	netWay.Option = option
	netWay.Status = NETWAY_STATUS_INIT

	//协议管理适配器
	protocolManagerOption := ProtocolManagerOption{
		Ip				: netWay.Option.ListenIp,
		WsPort			: netWay.Option.WsPort,
		TcpPort			: netWay.Option.TcpPort,
		WsUri			: netWay.Option.WsUri,
		UdpPort			: netWay.Option.UdpPort,
		OpenNewConnBack	: netWay.OpenNewConn,
		Log				: option.Log,
	}
	netWay.ProtocolManager =  NewProtocolManager(protocolManagerOption)
	err := netWay.ProtocolManager.Start()
	if err != nil {
		return nil,err
	}
	//http 工具
	//httpdOption := HttpdOption {
	//	LogOption 	: netWay.Option.LogOption,
	//	RootPath 	: netWay.Option.HttpdRootPath,
	//	Ip			: netWay.Option.ListenIp,
	//	Port		: netWay.Option.HttpPort,
	//	ParentCtx 	: option.OutCxt,
	//}
	//myHttpd = NewHttpd(httpdOption)
	//ws conn 管理
	netWay.ConnManager = NewConnManager(option.MaxClientConnNum,option.ConnTimeout,option.Log,netWay.Option.ContentType,PROTOCOL_WEBSOCKET)
	//统计模块
	netWay.Metrics = NewMyMetrics(option.Log)
	//netWay.Metrics.CreateCounter("total.fd.num")
	netWay.Metrics.CreateCounter("create_fd_ok")
	netWay.Metrics.CreateCounter("create_fd_failed")
	netWay.Metrics.CreateCounter("close_fd_num")
	netWay.Metrics.CreateCounter("total_output_num")
	netWay.Metrics.CreateGauge("total_output_size")
	netWay.Metrics.CreateCounter("total_input_num")
	netWay.Metrics.CreateGauge("total_input_size")
	//player.fd.num
	//player.fd.size

	//在外层的CTX上，派生netway自己的根ctx
	//startupCtx ,cancel := context.WithCancel(netWay.Option.OutCxt)
	startupCtx , cancelFunc := context.WithCancel(context.Background())
	netWay.MyCancelCtx = startupCtx
	netWay.MyCancelFunc = cancelFunc


	netWay.Status = NETWAY_STATUS_START

	option.Log.Info("netway startup finish.")
	return netWay,nil
}
//一个新客户端连接请求进入
func(netWay *NetWay)OpenNewConn( connFD FDAdapter) {
	netWay.Option.Log.Info("OpenNewConn:")
	var loginRes pb.ResponseLoginRes

	if netWay.Status == NETWAY_STATUS_CLOSE{//当前网关已经关闭了，还有新的连接进来
		//记录创建FD失败次数
		netWay.Metrics.CounterInc("create_fd_failed")
		errMsg := "netWay closing... not accept new connect , sleep 1!"
		netWay.Option.Log.Error(errMsg)
		//直接给一个FD发送消息，不做任何封装
		netWay.WriteMessage(websocket.TextMessage,connFD,[]byte(errMsg))
		//这里暂停一会，保证上面的消息能发送成功
		time.Sleep(time.Millisecond * 200)
		//直接关闭一个FD，不做任何多余处理
		netWay.CloseFD(connFD,CLOSE_SOURCE_SERVER_HAS_CLOSE)

		return
	}
	//是否超过了，最大可连接数
	if int32(len(netWay.ConnManager.Pool))   > netWay.Option.MaxClientConnNum{
		errMsg  := "more MaxClientConnNum"
		netWay.Option.Log.Error(errMsg)
		netWay.WriteMessage(websocket.TextMessage,connFD,[]byte(errMsg))
		netWay.CloseFD(connFD,CLOSE_SOURCE_MAX_CLIENT)
		return
	}
	//创建一个连接元素，将WS FD 保存到该容器中
	NewConn := netWay.ConnManager.CreateOneConn(connFD)
	defer func() {
		if err := recover(); err != nil {
			netWay.Option.Log.Panic("OpenNewConn:")
			netWay.CloseOneConn(NewConn, CLOSE_SOURCE_OPEN_PANIC)
		}
	}()
	//开始-登陆验证
	jwtData,firstMsg,err := netWay.loginPre(NewConn)
	if err != nil{
		//这里不用发消息了，也不用关闭FD，因为loginPre内部已经处理过了
		return
	}
	//登陆验证已通过，开始添加各种状态及初始化
	NewConn.UserId = jwtData.Payload.Uid
	//将新的连接加入到连接池中，并且与玩家ID绑定
	netWay.ConnManager.addConnPool( NewConn)
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
	NewConn.Protocol 	= firstMsg.ProtocolType
	NewConn.ContentType = firstMsg.ContentType
	loginRes = pb.ResponseLoginRes{
		Code: 200,
		ErrMsg: "",
		Uid: NewConn.UserId,
	}
	//告知玩家：登陆结果
	netWay.SendMsgCompressByUid(jwtData.Payload.Uid,"loginRes",&loginRes)
	//统计 当前FD 数量/历史FD数量
	netWay.Metrics.CounterInc("create_fd_ok")
	//初始化即登陆成功的响应均完成后，开始该连接的 消息IO 协程
	go NewConn.IOLoop()
	//netWay.serverPingRtt(time.Duration(rttMinTimeSecond),NewWsConn,1)
	netWay.Option.Log.Info("wsHandler end ,player login success!!!")

}

func(netWay *NetWay)heartbeat(requestClientHeartbeat pb.RequestClientHeartbeat,conn *Conn){
	now := GetNowTimeSecondToInt()
	conn.UpTime = int32(now)
}
//=================================
//直接关闭一个FD，主要用于：登陆就失败了的情况
func(netWay *NetWay)CloseFD(connFD FDAdapter,source int){
	connFD.Close()
	//记录主动关闭FD次数
	netWay.Metrics.CounterInc("close_fd_num")
}
//关闭一个已登陆成功的FD
func (netWay *NetWay)CloseOneConn(conn *Conn,source int){
	netWay.Option.Log.Info("Conn close ,source : "+strconv.Itoa(source) + " , " + strconv.Itoa(int(conn.UserId)))
	if conn.Status == CONN_STATUS_CLOSE {
		netWay.Option.Log.Error("CloseOneConn error :Conn.Status == CLOSE")
	}
	//通知同步服务，先做构造处理
	//mySync.CloseOne(conn)//这里可能还要再发消息

	//状态更新为已关闭，防止重复关闭
	conn.Status = CONN_STATUS_CLOSE
	//把后台守护  协程 先关了，不再收消息了
	conn.CloseChan <- 1
	//netWay.Players.delById(Conn.PlayerId)//这个不能删除，用于玩家掉线恢复的记录
	//先把玩家的在线状态给变更下，不然 mySync.close 里面获取房间在线人数，会有问题
	//myPlayerManager.upPlayerStatus(conn.UserId, PLAYER_STATUS_OFFLINE)
	err := conn.Conn.Close()
	if err != nil{
		netWay.Option.Log.Error("Conn.Conn.Close err:"+err.Error())
	}

	netWay.ConnManager.delConnPool(conn.UserId)
	//处理掉-已报名的玩家
	//myMatch.realDelOnePlayer(conn.PlayerId)
	//mySleepSecond(2,"CloseOneConn")
	//myMetrics.fastLog("total.fd.num",METRICS_OPT_DIM,0)
	//myMetrics.fastLog("history.fd.destroy",METRICS_OPT_INC,0)
	//netWay.Metrics.CounterDec("total.fd.num")
	netWay.Metrics.CounterInc("close_fd_num")
}
//退出，目前能直接调用此函数的，就只有一种情况：
//MAIN 接收到了中断信号，并执行了：context-cancel()，然后，startup函数的守护监听到，调用些方法
func  (netWay *NetWay)Shutdown() {
	netWay.Option.Log.Warn("netWay.Shutdown")
	if netWay.Status == NETWAY_STATUS_CLOSE{
		netWay.Option.Log.Error("Quit err :netWay.Status ==  NETWAY_STATUS_CLOSE")
		return
	}
	netWay.Status = NETWAY_STATUS_CLOSE//更新状态为：关闭

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


func  (netWay *NetWay)loginPreFailed(msg string ,closeSource int,conn *Conn){
	loginRes := pb.ResponseLoginRes{
		Code : 500,
		ErrMsg:msg,
	}
	netWay.SendMsgCompressByConn(conn,"loginRes",loginRes)
	netWay.CloseOneConn(conn, closeSource)
	netWay.Option.Log.Error(msg)
}
//首次建立连接，登陆验证，预处理
func  (netWay *NetWay)loginPre(conn *Conn)(jwt JwtData,firstMsg pb.Msg,err error,){
	//var loginRes pb.ResponseLoginRes

	content,err := conn.Read()
	if err != nil{
		netWay.loginPreFailed(err.Error(),CLOSE_SOURCE_FD_READ_EMPTY,conn)
		return jwt,firstMsg,errors.New("conn read err:"+err.Error())
	}
	msg,err := netWay.ProtocolManager.parserContentProtocol(content)
	if err != nil{
		netWay.loginPreFailed(err.Error(),CLOSE_SOURCE_FD_PARSE_CONTENT,conn)
		return jwt,firstMsg,err
	}
	//这里有个问题，连接成功后，C端立刻就得发消息，不然就异常~bug
	if msg.Action != "login"{//进到这里，肯定是有新连接被创建且回调了公共函数
		netWay.loginPreFailed("first msg must login api!!",CLOSE_SOURCE_FIRST_NO_LOGIN,conn)
		return
	}
	//开始：登陆/验证 过程
	jwtDataInterface,err := netWay.Router(msg,conn)
	jwt = jwtDataInterface.(JwtData)
	if err != nil{
		netWay.loginPreFailed(err.Error(),CLOSE_SOURCE_AUTH_FAILED,conn)
		return jwt ,firstMsg,err
	}
	netWay.Option.Log.Info("login jwt auth ok~~")
	return jwt,msg,nil
}
//登陆验证token
func(netWay *NetWay)login(requestLogin pb.RequestLogin,conn *Conn)(jwtData JwtData,err error){
	token := ""
	if netWay.Option.LoginAuthType == "jwt"{
		token = requestLogin.Token
		jwtData,err := ParseJwtToken(netWay.Option.LoginAuthSecretKey,token)
		return jwtData,err
	}else{
		netWay.Option.Log.Error("LoginAuthType err")
	}

	return jwtData,err
}

func  (netWay *NetWay)RecoverGoRoutine(back func(ctx context.Context),ctx context.Context,err interface{}){
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