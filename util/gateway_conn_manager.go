package util

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
	"zgoframe/protobuf/pb"
)


type ConnManager struct {
	Pool 				map[int32]*Conn //ws 连接池
	PoolRWLock 			*sync.RWMutex
	Close 				chan int
	MaxClientConnNum 	int32
	ConnTimeout			int32
	Log 				*zap.Logger

	DefaultContentType 	int32
	DefaultProtocol		int32
}

type Conn struct {
	AddTime			int32
	UpTime 			int32
	UserId			int32
	Status  		int
	Conn 			FDAdapter //TCP/WS Conn FD
	CloseChan 		chan int
	RTT 			int64
	SessionId 		string
	ConnManager		*ConnManager
	MsgInChan		chan pb.Msg
	ContentType		int32
	Protocol		int32
	//RTTCancelChan 	chan int
	//UdpConn 		bool
}

func NewConnManager(maxClientConnNum int32,connTimeout int32,log *zap.Logger,defaultContentType int32 ,defaultProtocol int32)*ConnManager {
	log.Info("NewConnManager instance:")
	connManager :=  new(ConnManager)
	//全局变量
	connManager.Pool 				= make(map[int32]*Conn)
	connManager.MaxClientConnNum 	= maxClientConnNum
	connManager.ConnTimeout 		= connTimeout
	connManager.Close 				= make(chan int)
	connManager.PoolRWLock 			= &sync.RWMutex{}
	connManager.Log 				= log
	connManager.DefaultContentType	= defaultContentType
	connManager.DefaultProtocol  	= defaultProtocol

	return connManager
}
//创建一个新的连接结构体
func (connManager *ConnManager)CreateOneConn(connFd FDAdapter)(myConn *Conn  ){
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	connManager.Log.Info("Create one Conn  client struct")

	now := int32( GetNowTimeSecondToInt())

	myConn = new (Conn)
	myConn.Conn 		= connFd
	myConn.UserId 		= 0
	myConn.AddTime 		= now
	myConn.UpTime 		= now
	myConn.Status  		= CONN_STATUS_INIT
	myConn.MsgInChan  	= make(chan pb.Msg,5000)
	myConn.ConnManager	= connManager
	//myConn.UdpConn    	= false
	//myConn.RTTCancelChan = make(chan int)

	connManager.Log.Info("reg conn callback CloseHandler")

	return myConn
}

func (connManager *ConnManager)getConnPoolById(id int32)(*Conn,bool){
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	conn,ok := connManager.Pool[id]
	return conn,ok
}
func (connManager *ConnManager)GetPlayerCtrlInfoById(userId int32)ProtocolCtrlInfo{
	var contentType  int32
	var protocolType int32
	if userId == 0{
		contentType = connManager.DefaultContentType
		protocolType = connManager.DefaultProtocol
	}else{
		conn ,empty := connManager.getConnPoolById(userId)
		//mylog.Debug("GetContentTypeById player",player)
		if empty{
			contentType = connManager.DefaultContentType
			protocolType = connManager.DefaultProtocol
		}else{
			contentType = conn.ContentType
			protocolType = conn.Protocol
		}
	}

	protocolCtrlInfo := ProtocolCtrlInfo{
		ContentType: contentType,
		ProtocolType: protocolType,
	}

	connManager.Log.Info("GetPlayerCtrlInfo uid : "+strconv.Itoa(int(userId))+" contentType:"+ strconv.Itoa(int(contentType)) + " , protocolType:" + strconv.Itoa(int(protocolType)))

	return protocolCtrlInfo
}

//往POOL里添加一个新的连接
func  (connManager *ConnManager)addConnPool( NewConn *Conn)error{
	if NewConn.UserId <= 0{
		connManager.Log.Error("addConnPool NewConn.UserId <= 0 ")
	}
	oldConn ,exist := connManager.getConnPoolById(NewConn.UserId)
	if exist{//该UID已经创建了连接，可能是在别处登陆，直接踢掉旧的连接
		msg := strconv.Itoa(int(NewConn.UserId)) + " kickOff old pid :"+strconv.Itoa(int(oldConn.UserId))
		connManager.Log.Warn(msg)
		//err := errors.New(msg)
		responseKickOff := pb.ResponseKickOff{
			Time: GetNowMillisecond(),
		}
		//给旧连接发送消息通知
		myNetWay.SendMsgCompressByConn(oldConn,"kickOff",&responseKickOff)
		time.Sleep(time.Millisecond * 200)
		myNetWay.CloseOneConn(oldConn,CLOSE_SOURCE_OVERRIDE)
	}
	connManager.Log.Info("addConnPool : " + strconv.Itoa(int(NewConn.UserId)))
	connManager.PoolRWLock.Lock()
	defer connManager.PoolRWLock.Unlock()
	connManager.Pool[NewConn.UserId] = NewConn

	return nil
}

func  (connManager *ConnManager)delConnPool(uid int32  ){
	connManager.Log.Warn("delConnPool uid :"+strconv.Itoa(int(uid)))
	connManager.PoolRWLock.Lock()
	defer connManager.PoolRWLock.Unlock()

	delete(connManager.Pool,uid)
}

func   (conn *Conn)Write(content []byte,messageType int){
	defer func() {
		if err := recover(); err != nil {
			conn.ConnManager.Log.Error("conn.Conn.WriteMessage failed:")
			myNetWay.CloseOneConn(conn,CLOSE_SOURCE_SEND_MESSAGE)
		}
	}()

	//myMetrics.fastLog("total.output.num",METRICS_OPT_INC,0)
	//myMetrics.fastLog("total.output.size",METRICS_OPT_PLUS,len(content))
	myNetWay.Metrics.CounterInc("total.output.num")
	myNetWay.Metrics.GaugeAdd("total.output.size",float64(StringToFloat(strconv.Itoa(len(content)))))

	//pid := strconv.Itoa(int(conn.UserId))
	//myMetrics.fastLog("player.fd.num."+pid,METRICS_OPT_INC,0)
	//myMetrics.fastLog("player.fd.size."+pid,METRICS_OPT_PLUS,len(content))

	conn.Conn.WriteMessage(messageType,[]byte(content))
}
func   (conn *Conn)UpLastTime(){
	conn.UpTime = int32( GetNowTimeSecondToInt() )
}

func   (conn *Conn)ReadBinary()(content []byte,err error){
	messageType , dataByte  , err  := conn.Conn.ReadMessage()
	if err != nil{
		conn.ConnManager.Log.Error("conn.Conn.ReadMessage err: "+err.Error())
		return content,err
	}
	conn.ConnManager.Log.Debug("conn.ReadMessage Binary messageType:"+ strconv.Itoa(messageType) +" len :"+strconv.Itoa(len(dataByte)) +" data:"  + string(dataByte))
	//content = string(dataByte)
	return dataByte,nil
}
//从FD中读取一条消息
func   (conn *Conn)Read()(content string,err error){
	// 设置消息的最大长度 - 暂无
	//conn.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(mynetWay.Option.IOTimeout)))
	messageType , dataByte  , err  := conn.Conn.ReadMessage()
	//_ , dataByte  , err  := conn.Conn.ReadMessage()
	if err != nil{
		//myMetrics.fastLog("total.input.err.num",METRICS_OPT_INC,0)
		conn.ConnManager.Log.Error("conn.Conn.ReadMessage err: " + err.Error())
		return content,err
	}
	myNetWay.Metrics.CounterInc("total.input.num")
	myNetWay.Metrics.GaugeAdd("total.input.size",float64(StringToFloat(strconv.Itoa(len(dataByte)))))

	//pid := strconv.Itoa(int(conn.UserId))
	//myMetrics.fastLog("player.fd.num."+pid,METRICS_OPT_INC,0)
	//myMetrics.fastLog("player.fd.size."+pid,METRICS_OPT_PLUS,len(content))

	conn.ConnManager.Log.Debug("conn.ReadMessage messageType:" + strconv.Itoa(messageType) +" len :"+strconv.Itoa(len(dataByte)) + " data:" +string(dataByte))
	content = string(dataByte)
	return content,nil
}

func  (conn *Conn)IOLoop(){
	conn.ConnManager.Log.Info("conn IOLoop:")
	conn.ConnManager.Log.Info("set conn status :"+strconv.Itoa(CONN_STATUS_EXECING)+ " make close chan")
	conn.Status = CONN_STATUS_EXECING
	conn.CloseChan = make(chan int)
	ctx,cancel := context.WithCancel(myNetWay.Option.OutCxt)
	go conn.ReadLoop(ctx)//读取消息，拆包，然后投入到队列中
	go conn.ProcessMsgLoop(ctx)//从队列中取出已拆包的值，做下一步处理：router
	//进入阻塞，监听外部<取消>操作
	<- conn.CloseChan
	conn.ConnManager.Log.Warn("IOLoop receive chan quit~~~")
	//取消上面两个协程
	cancel()
}
//一个协程挂了，再给拉起来
func  (conn *Conn) RecoverReadLoop(ctx context.Context){
	conn.ConnManager.Log.Warn("recover ReadLoop:")
	go conn.ReadLoop(ctx)
}
//死循环，从FD中读取消息
func  (conn *Conn)ReadLoop(ctx context.Context){
	defer func(ctx context.Context) {
		if err := recover(); err != nil {
			conn.ConnManager.Log.Panic("conn.ReadLoop panic defer :")
			conn.RecoverReadLoop(ctx)
		}
	}(ctx)
	for{
		select{
			case <-ctx.Done():
				conn.ConnManager.Log.Warn("connReadLoop receive signal: ctx.Done.")
				goto end
			default:
				//从ws 读取 数据
				content,err :=  conn.Read()
				if err != nil{
					IsUnexpectedCloseError := websocket.IsUnexpectedCloseError(err)
					conn.ConnManager.Log.Warn("connReadLoop connRead err:"+err.Error()+"IsUnexpectedCloseError:")
					if IsUnexpectedCloseError{
						myNetWay.CloseOneConn(conn, CLOSE_SOURCE_CLIENT_WS_FD_GONE)
						goto end
					}else{
						continue
					}
				}

				if content == ""{
					continue
				}
				//最后更新时间
				conn.UpLastTime()
				//解析消息内容
				msg,err  := myNetWay.ProtocolManager.parserContentProtocol(content)
				if err !=nil{
					conn.ConnManager.Log.Warn("parserContent err :" + err.Error())
					continue
				}
				//写入队列，等待其它协程处理，继续死循环
				conn.MsgInChan <- msg
		}
	}
end :
	conn.ConnManager.Log.Warn("connReadLoop receive signal: done.")
}
func  (conn *Conn) RecoverProcessMsgLoop(ctx context.Context){
	conn.ConnManager.Log.Warn("recover ReadLoop:")
	go conn.ProcessMsgLoop(ctx)
}
//从：FD里读取的消息（缓存队列），拿出来，做分发路由，处理
func  (conn *Conn)ProcessMsgLoop(ctx context.Context){
	defer func(ctx context.Context) {
		if err := recover(); err != nil {
			conn.ConnManager.Log.Panic("conn.ProcessMsgLoop panic defer :")
			conn.RecoverProcessMsgLoop(ctx)
		}
	}(ctx)

	for{
		ctxHasDone := 0
		select{
		case <-ctx.Done():
			ctxHasDone = 1
		case msg := <-conn.MsgInChan:
			conn.ConnManager.Log.Info("ProcessMsgLoop receive msg" + msg.Action)
			myNetWay.Router(msg,conn)
		}
		if ctxHasDone == 1{
			goto end
		}
	}
end :
	conn.ConnManager.Log.Warn("ProcessMsgLoop receive signal: done.")
}
//监听到某个FD被关闭后，回调函数
func  (conn *Conn)CloseHandler(code int, text string) error{
	myNetWay.CloseOneConn(conn, CLOSE_SOURCE_CLIENT)
	return nil
}
func (connManager *ConnManager)Shutdown(){
	connManager.Log.Warn("shutdown connManager")
	connManager.Close <- 1
	if len(connManager.Pool) <= 0{
		return
	}
	pool := connManager.getPoolAll( )
	for _,conn :=range pool{
		myNetWay.CloseOneConn(conn,CLOSE_SOURCE_CONN_SHUTDOWN)
	}
}
func (connManager *ConnManager)getPoolAll()map[int32]*Conn{
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	pool := make(map[int32]*Conn)
	for k,v := range connManager.Pool{
		pool[k] = v
	}
	return pool
}
func (connManager *ConnManager)Start(ctx context.Context){
	defer func(ctx context.Context ) {
		if err := recover(); err != nil {
			myNetWay.RecoverGoRoutine(connManager.Start,ctx,err)
		}
	}(ctx)

	connManager.Log.Warn("checkConnPoolTimeout start:")
	for{
		select {
		case   <-connManager.Close:
			goto end
		default:
			pool := connManager.getPoolAll()
			for _,v := range pool{
				now := int32 (GetNowTimeSecondToInt())
				x := v.UpTime + connManager.ConnTimeout
				if now  > x {
					myNetWay.CloseOneConn(v, CLOSE_SOURCE_TIMEOUT)
				}
			}
			time.Sleep(time.Second * 1)
			//mySleepSecond(1,"checkConnPoolTimeout")
		}
	}
end:
	connManager.Log.Warn(CTX_DONE_PRE+"checkConnPoolTimeout close")
}

