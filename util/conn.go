package util

import (
	"github.com/gorilla/websocket"
	"strconv"
)
//往一个连接里写入数据
func   (conn *Conn)Write(content []byte,messageType int){
	defer func() {
		if err := recover(); err != nil {
			conn.Log.Error("conn.Conn.WriteMessage failed:")
			myConnManager.CloseOneConn(conn,CLOSE_SOURCE_SEND_MESSAGE)
		}
	}()

	conn.FD.WriteMessage(messageType,[]byte(content))
}
//读取二进制信息
func   (conn *Conn)ReadBinary()(content []byte,err error){
	messageType , dataByte  , err  := conn.FD.ReadMessage()
	if err != nil{
		conn.Log.Error("conn.Conn.ReadMessage err: " + err.Error())
		return content,err
	}
	conn.Log.Debug("conn.ReadMessage Binary messageType:"+strconv.Itoa(messageType)  + " len :" + strconv.Itoa(len(dataByte)) +" data:" +string(dataByte))
	//content = string(dataByte)
	return dataByte,nil
}
//读取信息
func   (conn *Conn)Read()(content string,err error){
	// 设置消息的最大长度 - 暂无
	//conn.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(mynetWay.Option.IOTimeout)))
	//messageType , dataByte  , err  := conn.Conn.ReadMessage()
	msgType , dataByte  , err  := conn.FD.ReadMessage()
	MyPrint("conn read , msgType:" ,msgType ,  " dataByte:" , string(dataByte))
	if err != nil{
		//myMetrics.fastLog("total.input.err.num",METRICS_OPT_INC,0)
		conn.Log.Error("conn.Conn.ReadMessage err: " + err.Error())
		return content,err
	}
	//myMetrics.fastLog("total.input.num",METRICS_OPT_INC,0)
	//myMetrics.fastLog("total.input.size",METRICS_OPT_PLUS,len(dataByte))
	//
	//pid := strconv.Itoa(int(conn.PlayerId))
	//myMetrics.fastLog("player.fd.num."+pid,METRICS_OPT_INC,0)
	//myMetrics.fastLog("player.fd.size."+pid,METRICS_OPT_PLUS,len(content))

	//mylog.Debug("conn.ReadMessage messageType:",messageType , " len :",len(dataByte) , " data:" ,string(dataByte))
	content = string(dataByte)
	return content,nil
}
//更新最后操作时间
func   (conn *Conn)UpLastTime(){
	conn.UpTime = int( GetNowTimeSecondToInt() )
}
func  (conn *Conn)IOLoop(){
	conn.Log.Info("IOLoop . set conn status :"+ strconv.Itoa(CONN_STATUS_EXECING) +" make close chan")
	conn.Status = CONN_STATUS_EXECING
	conn.CloseChan = make(chan int)
	//ctx,cancel := context.WithCancel(myConnManager.Ctx)
	go conn.ReadLoop()
	//go conn.ProcessMsgLoop(ctx)
	//<- conn.CloseChan
	//conn.Log.Warn("IOLoop receive chan quit~~~")
	//cancel()
}
//func  (conn *Conn) RecoverReadLoop(ctx context.Context){
//	conn.Log.Warn("recover ReadLoop:")
//	go conn.ReadLoop(ctx)
//}
func  (conn *Conn)ReadLoop( ){
	//defer func(ctx context.Context) {
	//	if err := recover(); err != nil {
	//		conn.Log.Panic("conn.ReadLoop panic defer :")
	//		conn.RecoverReadLoop(ctx)
	//	}
	//}(ctx)
	for{
		select{
		case <-conn.CloseChan:
			conn.Log.Warn("connReadLoop receive signal: ctx.Done.")
			goto end
		default:
			//从ws 读取 数据
			content,err :=  conn.Read()
			if err != nil{
				IsUnexpectedCloseError := websocket.IsUnexpectedCloseError(err)
				conn.Log.Warn("connReadLoop connRead err:" + err.Error() + "IsUnexpectedCloseError:")
				if IsUnexpectedCloseError{
					myConnManager.CloseOneConn(conn, CLOSE_SOURCE_CLIENT_WS_FD_GONE)
					goto end
				}else{
					continue
				}
			}

			if content == ""{
				continue
			}

			conn.UpLastTime()
			//msg,err  := myProtocolManager.parserContentProtocol(content)
			//if err !=nil{
			//	conn.Log.Warn("parserContent err :" +err.Error())
			//	continue
			//}
			//conn.MsgInChan <- content
			myConnManager.BackFunc(content,conn)
		}
	}
end :
	conn.Log.Warn("connReadLoop receive signal: done.")
}

//监听到某个FD被关闭后，回调函数
func  (conn *Conn)CloseHandler(code int, text string) error{
	MyPrint("CloseHandler:",code,text)
	myConnManager.CloseOneConn(conn, CLOSE_SOURCE_CLIENT)
	return nil
}

//获取一个连接的 头信息
func (conn *Conn)GetCtrlInfoById()ProtocolCtrlInfo{
	var contentType  int
	var protocolType int
	if conn.PlayerId == 0{
		contentType = myConnManager.DefaultContentType
		protocolType = myConnManager.DefaultProtocol
	}else{
		contentType = conn.ContentType
		protocolType = conn.ProtocolType
	}

	protocolCtrlInfo := ProtocolCtrlInfo{
		ContentType: contentType,
		ProtocolType: protocolType,
	}
	return protocolCtrlInfo
}

//===================================================================================================================



