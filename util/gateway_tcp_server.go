package util

import (
	"bufio"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
	"context"
)

type TcpServer struct {
	listener net.Listener
	Option TcpServerOption
	CancelCtx context.Context
	CancelFunc context.CancelFunc
}

type TcpServerOption struct{
	OutCxt   context.Context
	Ip 		string
	Port 	string
	Log 	*zap.Logger
	ProtocolManager *ProtocolManager	//给外部提供一个接口，用于将SOCKER FD 注册给外部
}

func NewTcpServer(tcpServerOption TcpServerOption)*TcpServer{
	tcpServerOption.Log.Info("NewTcpServer instance:")
	tcpServer := new (TcpServer)
	tcpServer.Option = tcpServerOption


	cancelCtx,canclFunc := context.WithCancel(context.Background())

	tcpServer.CancelCtx = cancelCtx
	tcpServer.CancelFunc = canclFunc

	return tcpServer
}

//func (tcpServer *TcpServer)ListeningClose(){
//<- tcpServer.OutCxt.Done()
//tcpServer.Shutdown()
//}
//outCtx:这里没用上，因为accept是阻塞的模式，只能用另外的方式close
func  (tcpServer *TcpServer)Start()error{
	ipPort := tcpServer.Option.Ip + ":" + tcpServer.Option.Port
	tcpServer.Option.Log.Info("tcpServer start:"+ipPort)

	listener,err :=net.Listen("tcp",ipPort)
	if err !=nil{
		errorMsg := "net.Listen tcp err:" + err.Error()
		tcpServer.Option.Log.Error(errorMsg)
		return errors.New(errorMsg)
		//zlib.PanicPrint("tcp net.Listen tcp err:"+err.Error())
	}
	tcpServer.Option.Log.Info("startTcpServer:")
	tcpServer.listener = listener
	go tcpServer.Accept()

	return nil
}

func   (tcpServer *TcpServer)Shutdown( ){
	//先停掉 死循环的每个FD的读取协程
	tcpServer.CancelFunc()
	tcpServer.Option.Log.Info("Shutdown tcpServer ")
	err := tcpServer.listener.Close()
	if err != nil{
		//mylog.Error("tcpServer.listener.Close err :",err)
	}
}

func (tcpServer *TcpServer)Accept( ){
	for {
		//阻塞：获取连接方FD
		conn,err := tcpServer.listener.Accept()
		if err == nil{
			tcpServer.Option.Log.Info("listener.Accept new conn:")
		}else{
			tcpServer.Option.Log.Error("listener.Accept err :"+err.Error())
			//if strings.Contains(err.Error(), "use of closed network connection") {
				//mylog.Warning("TcpAccept end.")
				//break
			//}else{
			//
			//}
			//continue
			break
		}
		tcpConn := NewTcpConn(conn,tcpServer.Option.Log)
		//myTcpServer.pool = append(myTcpServer.pool,tcpConn)
		go tcpConn.start(tcpServer )
	}
}
//====================================================================
type TcpConn struct {
	conn 				net.Conn
	MsgQueue 			[][]byte
	callbackCloseHandle func(code int, text string)error
	Log				 	*zap.Logger
}

func NewTcpConn(conn net.Conn,log *zap.Logger)*TcpConn{
	log.Info("TcpConnNew")
	tcpConn := new (TcpConn)
	tcpConn.conn = conn
	tcpConn.Log = log
	tcpConn.callbackCloseHandle = nil
	return tcpConn
}

func  (tcpConn *TcpConn)start(tcpServer *TcpServer){
	//mylog.Info("TcpConn.start")

	go tcpConn.readLoop(tcpServer.CancelCtx)
	//先睡眠100毫秒，给上面的协程 读取消息的时间
	time.Sleep(time.Millisecond * 100)
	if tcpServer != nil{
		//回调：将新的FD连接交给主层统一管理
		tcpServer.Option.ProtocolManager.tcpHandler(tcpConn)//将当前socker FD 传给外部
	}
}

func  (tcpConn *TcpConn)SetCloseHandler(h func(code int, text string)error) {
	tcpConn.callbackCloseHandle = h
}

func  (tcpConn *TcpConn)Close()error{
	//myTcpServer.pool[]
	tcpConn.realClose(1)
	return nil
}

func  (tcpConn *TcpConn)realClose(source int){
	if tcpConn.callbackCloseHandle != nil{
		tcpConn.callbackCloseHandle(555,"close")
	}
	//mylog.Warning("realClose :",source)
	err := tcpConn.conn.Close()
	if err!=nil{

	}
	//mylog.Error("tcpConn.conn.Close:",err)
}

func  (tcpConn *TcpConn)readLoop(CancelCtx context.Context ){
	defer func(ctx context.Context ) {
		//if err := recover(); err != nil {
		//	myNetWay.RecoverGoRoutine(tcpConn.readLoop,ctx,err)
		//}
	}(CancelCtx)
	//创建一个IO 读取器
	reader := bufio.NewReader(tcpConn.conn)
	isBreak := 0
	//loopReadCnt := 0
	for {
		select {
		case <- CancelCtx.Done():
			isBreak = 1
		default:

		}
		//从IO读取器中批量读取字节数据，以\n为分隔符
		onsMessage,err := reader.ReadBytes('\n')
		if err != nil{
			tcpConn.Log.Error("reader.ReadBytes err:"+err.Error())
			if err == io.EOF{
				continue
			}
		}

		if len(onsMessage) <= 0{
			tcpConn.Log.Error("tcp readLoop reader empty")
			continue
		}

		MyPrint("tcp readLoop append msg queue  msgLen:",len(onsMessage))
		tcpConn.MsgQueue = append(tcpConn.MsgQueue,onsMessage)

		if isBreak == 1{
			break
		}
	}
}

func  (tcpConn *TcpConn)ReadMessage()(messageType int, p []byte, err error){
	if len(tcpConn.MsgQueue) == 0 {
		str := ""
		return messageType,[]byte(str),nil
	}
	//从数组头部取出一条消息
	data := tcpConn.MsgQueue[0]
	//删除头部的这条消息
	tcpConn.MsgQueue = tcpConn.MsgQueue[1:]
	return messageType,data,nil
}

func  (tcpConn *TcpConn)WriteMessage(messageType int, data []byte) error {
	tcpConn.conn.Write(data)
	return nil
}