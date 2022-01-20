package util

import (
	"github.com/gorilla/websocket"
)

//============适配：WS TCP 传输内容=====================
type FDAdapter interface {
	//连接断开后，最底层的连接代码会最先捕获，处理完后，需要告知：上层函数
	SetCloseHandler(h func(code int, text string) error)
	//写入一条消息
	WriteMessage(messageType int, data []byte) error
	//读取一条消息
	ReadMessage()(messageType int, p []byte, err error)
	//主动关闭一个FD
	Close()error
	//C端信息
	RemoteAddr()string
}
//===========================

//实现了FDAdapter接口，用于WS协议的内容操作
type WebsocketConnImp struct {
	FD	*websocket.Conn
}
func WebsocketConnImpNew(FD *websocket.Conn)*WebsocketConnImp{
	websocketConnImp := new (WebsocketConnImp)
	websocketConnImp.FD = FD
	return websocketConnImp
}

func (websocketConnImp *WebsocketConnImp)SetCloseHandler(h func(code int, text string)error){
	websocketConnImp.FD.SetCloseHandler(h)
}

func (websocketConnImp *WebsocketConnImp)WriteMessage(messageType int, data []byte) error{
	return websocketConnImp.FD.WriteMessage(messageType,data)
}

func (websocketConnImp *WebsocketConnImp)Close()error{
	return websocketConnImp.FD.Close()
}

func (websocketConnImp *WebsocketConnImp)ReadMessage()(messageType int, p []byte, err error){
	return websocketConnImp.FD.ReadMessage()
}

func (websocketConnImp *WebsocketConnImp)RemoteAddr()(info string){
	net :=  websocketConnImp.FD.RemoteAddr()
	return net.String() + " " + net.Network()
}

//=========================
//实现了FDAdapter接口，用于TCP协议的内容操作
type TcpConnImp struct {
	FD	*TcpConn
}
func TcpConnImpNew(FD *TcpConn)*TcpConnImp{
	tcpConnImp := new (TcpConnImp)
	tcpConnImp.FD = FD
	return tcpConnImp
}

func (tcpConnImp *TcpConnImp)SetCloseHandler(h func(code int, text string)error){
	tcpConnImp.FD.SetCloseHandler(h)
}

func (tcpConnImp *TcpConnImp)WriteMessage(messageType int, data []byte) error{
	return tcpConnImp.FD.WriteMessage(messageType,data)
}

func (tcpConnImp *TcpConnImp)Close()error{
	tcpConnImp.FD.CloseChan <- 1
	return tcpConnImp.FD.ServerClose()
}

func (tcpConnImp *TcpConnImp)ReadMessage()(messageType int, p []byte, err error){
	return tcpConnImp.FD.ReadMessage()
}

func (tcpConnImp *TcpConnImp)RemoteAddr()(info string){
	net := tcpConnImp.FD.conn.RemoteAddr()
	return net.String() + " " + net.Network()
}
