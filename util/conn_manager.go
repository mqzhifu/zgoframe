package util

import (
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
	"errors"
	"context"
)

//连接 管理池


const(

	CONN_STATUS_INIT 	= 1	//初始化
	CONN_STATUS_EXECING = 2	//运行中
	CONN_STATUS_CLOSE 	= 3	//已关闭


	CLOSE_SOURCE_CLIENT 			= 1	//客户端-主动断开连接
	CLOSE_SOURCE_AUTH_FAILED 		= 21//客户端首次连接，登陆动作,服务端验证失败
	CLOSE_SOURCE_FD_READ_EMPTY 		= 22//客户端首次连接，登陆动作,服务端read信息为空
	CLOSE_SOURCE_FD_PARSE_CONTENT 	= 23//客户端首次连接，登陆动作,解析内容时出错
	CLOSE_SOURCE_FIRST_NO_LOGIN 	= 24//客户端首次连接，登陆动作,内容解出来了，但是action!=login
	CLOSE_SOURCE_CREATE 			= 3	//初始化 连接类失败，可能是连接数过大
	CLOSE_SOURCE_OPEN_PANIC			= 31//初始化 新连接创建成功后，上层要再重新做一次连接，结果未知panic
	CLOSE_SOURCE_OVERRIDE 			= 4	//创建新连接时，发现，该用户还有一个未关闭的连接,kickoff模式下，这条就没意义了
	CLOSE_SOURCE_TIMEOUT 			= 5	//最后更新时间 ，超时.后台守护协程触发
	CLOSE_SOURCE_SIGNAL_QUIT 		= 6 //接收到关闭信号，netWay.Quit触发
	CLOSE_SOURCE_CLIENT_WS_FD_GONE 	= 7	//S端读取连接消息时，异常了~可能是：客户端关闭了连接
	CLOSE_SOURCE_SEND_MESSAGE 		= 8 //S端给某个连接发消息，结果失败了，这里概率是连接已经断了
	CLOSE_SOURCE_CONN_SHUTDOWN 		= 11


	CONTENT_TYPE_JSON 		= 1		//内容类型 json
	CONTENT_TYPE_PROTOBUF 	= 2		//proto_buf


	PROTOCOL_TCP 		= 1
	PROTOCOL_UDP 		= 3
	PROTOCOL_WEBSOCKET 	= 2

)

type ConnManager struct {
	Pool 				map[int]*Conn //ws 连接池
	PoolRWLock 			*sync.RWMutex	//IO 连接池map时 的锁
	CloseChan			chan int		//<关闭>通道,主要是结束超时检查
	MaxClientConnNum 	int			//最大可连接数
	ConnTimeout			int			//连接超时时间
	Log 				*zap.Logger
	Ctx 				context.Context
	RecoverGo    		*RecoverGo

	DefaultContentType	int
	DefaultProtocol		int

	BackFunc			func(msg string,conn *Conn)
}
//IO 消息体
type ConnMsg struct{
	ActionId     int  `json:"action_id"`
	Action       string `json:"action"`
	Content      string `json:"content"`
	ContentType  int  `json:"content_type"`
	ProtocolType int  `json:"protocol_type"`
	SessionId    string `son:"session_id"`
}
//踢人 消息体
type ResponseKickOff struct {
	Time int64 `json:"time"`
}
//存储一个 连接FD
type Conn struct {
	AddTime			int
	UpTime 			int
	PlayerId		int
	Status  		int
	FD 			FDAdapter //socket : websocket tcp udp
	CloseChan 		chan int
	RTT 			int64
	//RTTCancelChan 	chan int
	//MsgInChan		chan string
	SessionId 		string
	UdpConn 		bool
	ContentType 	int		//本次会话内容类型
	ProtocolType	int		//本次会话协议类型
	Log 			*zap.Logger
}


type ProtocolCtrlInfo struct {
	ContentType int
	ProtocolType int
}

var myConnManager *ConnManager
//var CbFunc func(msg Msg,conn *Conn)error

func NewConnManager(maxClientConnNum int,connTimeout int, recoverGo *RecoverGo,log *zap.Logger,backFunc func(msg string,conn *Conn) )*ConnManager {
	connManager :=  new(ConnManager)
	//全局变量
	connManager.Pool 				= make(map[int]*Conn)
	connManager.MaxClientConnNum 	= maxClientConnNum
	connManager.ConnTimeout 		= connTimeout
	connManager.CloseChan 			= make(chan int)
	//connManager.Ctx 				= context.Background()
	connManager.PoolRWLock 			= &sync.RWMutex{}

	connManager.BackFunc			= backFunc
	connManager.DefaultContentType	= CONTENT_TYPE_JSON
	connManager.DefaultProtocol		= PROTOCOL_WEBSOCKET
	connManager.RecoverGo			= recoverGo
	connManager.Log					= log
	connManager.ConnTimeout			= 600

	myConnManager 					= connManager

	return connManager
}
func (connManager *ConnManager)Start(){
	//defer func(ctx context.Context ) {
	//	if err := recover(); err != nil {
	//		connManager.RecoverGo.RecoverGoRoutine(connManager.Start,ctx,err)
	//	}
	//}(ctx)

	connManager.Log.Warn("checkConnPoolTimeout start:")
	for{
		select {
		case   <-connManager.CloseChan:
			goto end
		default:
			pool := connManager.getPoolAll()
			for _,v := range pool{
				now :=  GetNowTimeSecondToInt()
				x := v.UpTime + connManager.ConnTimeout
				if now  > x {
					myConnManager.CloseOneConn(v, CLOSE_SOURCE_TIMEOUT)
				}
			}
			time.Sleep(time.Second * 1)
			//mySleepSecond(1,"checkConnPoolTimeout")
		}
	}
end:
	connManager.Log.Warn(" checkConnPoolTimeout close ")
}
func (connManager *ConnManager)Shutdown(){
	connManager.Log.Warn("shutdown connManager")
	connManager.CloseChan <- 1
	if len(connManager.Pool) <= 0{
		return
	}
	pool := connManager.getPoolAll( )
	for _,conn :=range pool{
		myConnManager.CloseOneConn(conn,CLOSE_SOURCE_CONN_SHUTDOWN)
	}
}

//创建一个新的连接结构体
func (connManager *ConnManager)CreateOneConn(connFd FDAdapter)(myConn *Conn,err error ){
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	connManager.Log.Info("Create one Conn  client struct")
	if len(connManager.Pool)   > connManager.MaxClientConnNum{
		connManager.Log.Error("more MaxClientConnNum")
		return myConn,errors.New("more MaxClientConnNum")
	}
	now := GetNowTimeSecondToInt()

	myConn = new (Conn)
	myConn.FD 		= connFd	//*websocket.Conn
	myConn.PlayerId 	= 0
	myConn.AddTime 		= now
	myConn.UpTime 		= now
	myConn.Status  		= CONN_STATUS_INIT
	//myConn.MsgInChan  	= make(chan string,5000)
	myConn.UdpConn    	= false
	myConn.Log			= connManager.Log
	//myConn.RTTCancelChan = make(chan int)

	myConn.ContentType 	= myConnManager.DefaultContentType
	myConn.ProtocolType = myConnManager.DefaultProtocol

	//MyPrint("FINISH.")
	connManager.Log.Debug("CreateOneConn finish.")

	return myConn,nil
}
//往POOL里添加一个新的连接
func  (connManager *ConnManager)AddConnPool( NewConn *Conn)error{
	oldConn ,exist := connManager.getConnPoolById(NewConn.PlayerId)
	if exist{
		//msg := strconv.Itoa(int(NewConn.PlayerId)) + " player has joined conn pool ,addTime : "+strconv.Itoa(int(v.AddTime)) + " , u can , kickOff old fd.?"
		msg := strconv.Itoa(int(NewConn.PlayerId)) + " kickOff old pid :"+strconv.Itoa(int(oldConn.PlayerId))
		connManager.Log.Warn(msg)
		return errors.New("has-exist")
	}
	connManager.Log.Info("addConnPool : "+ strconv.Itoa(int(NewConn.PlayerId)))
	connManager.PoolRWLock.Lock()
	defer connManager.PoolRWLock.Unlock()
	connManager.Pool[NewConn.PlayerId] = NewConn
	return nil
}


func (connManager *ConnManager)CloseOneConn(conn *Conn,source int){
	connManager.Log.Info("Conn close ,source : "+strconv.Itoa(source) + " pid:"  +strconv.Itoa(int(conn.PlayerId)))
	if conn.Status == CONN_STATUS_CLOSE {
		connManager.Log.Error("CloseOneConn error :Conn.Status == CLOSE")
		return
	}
	//通知同步服务，先做构造处理
	//mySync.CloseOne(conn)//这里可能还要再发消息

	//状态更新为已关闭，防止重复关闭
	conn.Status = CONN_STATUS_CLOSE
	//把后台守护  协程 先关了，不再收消息了
	conn.CloseChan <- 1
	//netWay.Players.delById(Conn.PlayerId)//这个不能删除，用于玩家掉线恢复的记录
	//先把玩家的在线状态给变更下，不然 mySync.close 里面获取房间在线人数，会有问题
	//myPlayerManager.upPlayerStatus(conn.PlayerId, PLAYER_STATUS_OFFLINE)
	err := conn.FD.Close()
	if err != nil{
		connManager.Log.Error("Conn.Conn.Close err:"+err.Error())
	}

	connManager.delConnPool(conn.PlayerId)
	//处理掉-已报名的玩家
	//myMatch.realDelOnePlayer(conn.PlayerId)
	//mySleepSecond(2,"CloseOneConn")
	//myMetrics.fastLog("total.fd.num",METRICS_OPT_DIM,0)
	//myMetrics.fastLog("history.fd.destroy",METRICS_OPT_INC,0)
}

func (connManager *ConnManager)getPoolAll()map[int]*Conn{
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	pool := make(map[int]*Conn)
	for k,v := range connManager.Pool{
		pool[k] = v
	}
	return pool
}

func (connManager *ConnManager)getConnPoolById(id int)(*Conn,bool){
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	conn,ok := connManager.Pool[id]
	return conn,ok
}
func  (connManager *ConnManager)delConnPool(uid int  ){
	connManager.Log.Warn("delConnPool uid :"+ strconv.Itoa(uid))
	connManager.PoolRWLock.Lock()
	defer connManager.PoolRWLock.Unlock()

	delete(connManager.Pool,uid)
}


//======================================================================================================================
