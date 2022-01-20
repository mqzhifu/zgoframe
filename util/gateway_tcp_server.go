package util

import (
	"bufio"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
)

type TcpServer struct {
	listener 	net.Listener
	Option 		TcpServerOption
}

type TcpServerOption struct{
	Ip 				string
	Port 			string
	Log 			*zap.Logger
	MsgSeparator 	string
	ProtocolManager *ProtocolManager	//给外部提供一个接口，用于将SOCKER FD 注册给外部
}

func NewTcpServer(tcpServerOption TcpServerOption)*TcpServer{
	tcpServerOption.Log.Info("NewTcpServer instance:")

	tcpServer := new (TcpServer)
	tcpServer.Option = tcpServerOption

	return tcpServer
}
//启动监听
func  (tcpServer *TcpServer)Start()error{
	ipPort := tcpServer.Option.Ip + ":" + tcpServer.Option.Port
	tcpServer.Option.Log.Info("tcpServer start:"+ipPort)
	//开启TCP监听
	listener,err :=net.Listen("tcp",ipPort)
	if err !=nil{
		errorMsg := "net.Listen tcp err:" + err.Error()
		tcpServer.Option.Log.Error(errorMsg)
		return errors.New(errorMsg)
	}
	tcpServer.Option.Log.Info("startTcpServer:")
	tcpServer.listener = listener
	go tcpServer.Accept()

	return nil
}

func   (tcpServer *TcpServer)Shutdown( ){
	//需要关闭的有两个： 1. accept  2 由accept 新产生的FD 的read 协程
	//此类只是最简单的封闭，可以直接关闭accept ，但是新产生的 read协程由一层调用 close 关闭
	//所以，调用此函数前，必须外层统一调用了close后，最终此函数只是简单的关闭accept
	tcpServer.Option.Log.Info("Shutdown tcpServer ")
	err := tcpServer.listener.Close()
	if err != nil{
		//mylog.Error("tcpServer.listener.Close err :",err)
	}
}
//接收新连接
func (tcpServer *TcpServer)Accept( ){
	for {
		//阻塞：获取新连接方FD
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
		tcpConn := NewTcpConn(conn,tcpServer.Option.Log,tcpServer)
		//myTcpServer.pool = append(myTcpServer.pool,tcpConn)
		go tcpConn.start(tcpServer )
	}
}
//====================================================================
type TcpConn struct {
	conn 				net.Conn	//真实的TCP FD
	MsgQueue 			[][]byte	//存储C端发送到来的消息内容
	Log				 	*zap.Logger
	CloseChan 			chan int
	TcpServer			*TcpServer
	callbackCloseHandle func(code int, text string)error	//客户端主动关闭连接时，通知外层调用者
}
//每一个TCP新的连接，就会创建一个结构体，统一管理
func NewTcpConn(conn net.Conn,log *zap.Logger ,tcpServer *TcpServer)*TcpConn{
	log.Info("TcpConnNew")
	tcpConn := new (TcpConn)
	tcpConn.conn 	= conn
	tcpConn.Log 	= log

	tcpConn.callbackCloseHandle = nil
	tcpConn.CloseChan = make(chan int)
	tcpConn.TcpServer = tcpServer
	return tcpConn
}

func  (tcpConn *TcpConn)start(tcpServer *TcpServer){
	tcpConn.Log.Info("TcpConn.start")
	//创建新协议：死循环 读取 C端发送过来的数据
	go tcpConn.readLoop()
	//先睡眠100毫秒，给上面的协程 读取消息的(处理)时间
	time.Sleep(time.Millisecond * 100)
	//将当前 FD 传给外部回调函数，由外部统一管理
	tcpServer.Option.ProtocolManager.tcpHandler(tcpConn)
}
//当C端主动断开连接后，外层需要通知，设置外层回调函数
func  (tcpConn *TcpConn)SetCloseHandler(h func(code int, text string)error) {
	tcpConn.callbackCloseHandle = h
}
//客户端主动关闭，服务端被动关闭，同时通知上层回调函数
func  (tcpConn *TcpConn)ClientClose()error{
	tcpConn.Log.Warn("tcpConn ClientClose")
	if tcpConn.callbackCloseHandle != nil{
		tcpConn.callbackCloseHandle(555,"close")
	}
	err := tcpConn.conn.Close()
	if err!=nil{
		tcpConn.Log.Error("tcpConn ClientClose err:"+err.Error())
		return err
	}

	return nil
}
//服务端主动关闭
func  (tcpConn *TcpConn)ServerClose( )error{
	err := tcpConn.conn.Close()
	if err!=nil{
		tcpConn.Log.Error("tcpConn close err:"+err.Error())
		return err
	}
	return nil
}
//死循环 - 从C端读取消息，并保存到队列中
func  (tcpConn *TcpConn)readLoop(){
	//defer func(ctx context.Context ) {
	//	//if err := recover(); err != nil {
	//	//	myNetWay.RecoverGoRoutine(tcpConn.readLoop,ctx,err)
	//	//}
	//}(CancelCtx)
	//创建一个IO 读取器
	reader := bufio.NewReader(tcpConn.conn)
	isBreak := 0
	for {
		//监听外部关闭信号
		select {
			case <- tcpConn.CloseChan:
				isBreak = 1
			default:
		}
		//从IO读取器中批量读取字节数据，以\n为分隔符
		msgSeparatorBytes := tcpConn.TcpServer.Option.MsgSeparator[0:1]
		onsMessage,err := reader.ReadBytes(msgSeparatorBytes[0])
		if err != nil{
			tcpConn.Log.Error("reader.ReadBytes err:"+err.Error())
			if err == io.EOF{
				tcpConn.Log.Warn("tcpConn readLoop close from io.EOF ")
				//对端已主动关闭
				tcpConn.ClientClose()
				break
			}
		}

		if len(onsMessage) <= 0{
			tcpConn.Log.Error("tcp readLoop reader empty")
		//	continue
		}

		MyPrint("tcp readLoop append msg queue  msgLen:",len(onsMessage))
		tcpConn.MsgQueue = append(tcpConn.MsgQueue,onsMessage)

		if isBreak == 1{
			tcpConn.Log.Warn("tcpConn readLoop close from chan.")
			break
		}
	}
}
//外层 - 读取消息
func  (tcpConn *TcpConn)ReadMessage()(messageType int, p []byte, err error){
	if len(tcpConn.MsgQueue) == 0 {
		str := ""
		return messageType,[]byte(str),nil
	}
	//从数组头部取出一条消息
	data := tcpConn.MsgQueue[0]
	//删除头部的这条消息
	tcpConn.MsgQueue = tcpConn.MsgQueue[1:]
	return TRAN_MESSAGE_TYPE_BINARY,data,nil
}
//外层 - 写入消息
func  (tcpConn *TcpConn)WriteMessage(messageType int, data []byte) error {
	tcpConn.conn.Write(data)
	return nil
}