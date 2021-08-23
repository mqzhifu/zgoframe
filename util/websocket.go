package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type Websocket struct {
	//Ctx context.Context
	Option  WebsocketOption
}

var httpUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的CORS 跨域请求，正式环境可以关闭
	CheckOrigin: func(r *http.Request) bool {
		return true
	},

}

type WebsocketOption struct{
	WsUri			string	`json:"wsUri"`
	Port  			string	`json:"wsPort"`
	ListenerIp 		string	`json:"listenerIp"`
	OutIp 			string	`json:"outIp"`
	OutCxt 			context.Context
	Log				*zap.Logger
	OpenNewConnBack	func ( connFD FDAdapter)
}

func NewWebsocket(httpGin *gin.Engine , option WebsocketOption)*Websocket{
	websocket := new(Websocket)
	websocket.Option = option
	httpGin.GET(option.WsUri,websocket.HttpHandler)
	//ctx,cancelFunc := context.WithCancel(context.Background())
	return websocket
}
func (ws *Websocket )Quit(){

}

func (ws *Websocket )HttpHandler(c *gin.Context){
	wsConnFD, err := httpUpGrader.Upgrade(c.Writer, c.Request, nil)
	ws.Option.Log.Info("Upgrade this http req to websocket")
	if err != nil {
		ws.Option.Log.Error("Upgrade websocket failed: " + err.Error())
		return
	}
	//
	//newCtx,cancelFunc := context.WithCancel(context.Background())
	//wsConn := new (WsConn)
	//wsConn.FD = wsConnFD
	//wsConn.Ctx = newCtx
	//wsConn.Log = ws.Option.Log
	//wsConn.CtxCancelFunc = cancelFunc
	//wsConn.FD.SetCloseHandler(wsConn.ConnCloseHandler)

	imp := WebsocketConnImpNew(wsConnFD)
	ws.Option.OpenNewConnBack(imp)
	//go wsConn.ReadLoop()

}

type WsConn struct {
	FD  *websocket.Conn
	Ctx context.Context
	CtxCancelFunc context.CancelFunc
	Log				*zap.Logger
}
//func   (wsConn *WsConn) ConnCloseHandler( code int, text string )error{
//	wsConn.Log.Info("ConnCloseHandler")
//	wsConn.CtxCancelFunc()
//	return nil
//}
//
//func  (wsConn *WsConn)ReadLoop( ){
//	wsConn.Log.Info("new ReadLoop")
//	for{
//		select{
//		case <-wsConn.Ctx.Done():
//			wsConn.Log.Warn("wsConn ReadLoop  receive signal: ctx.Done.")
//			goto end
//		default:
//			//从ws 读取 数据
//			messageType , dataByte  , err  := wsConn.FD.ReadMessage()
//			if err != nil{
//
//			}
//			content := string(dataByte)
//			if content == ""{
//				continue
//			}
//			wsConn.Log.Info("read msg:" + content + strconv.Itoa(messageType))
//			//service.ProcessOne(dataByte)
//		}
//	}
//end :
//	wsConn.Log.Warn("wsConn ReadLoop: end.")
//}

type FDAdapter interface {
	//连接断开后，最底层的连接代码会最先捕获，处理完后，需要告知：上层函数
	SetCloseHandler(h func(code int, text string) error)
	//写入一条消息
	WriteMessage(messageType int, data []byte) error
	//读取一条消息
	ReadMessage()(messageType int, p []byte, err error)
	//主动关闭一个FD
	Close()error
}

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