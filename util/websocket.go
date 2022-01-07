package util

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Websocket struct {
	httpServer  *http.Server
	Option  WebsocketOption
}

type WebsocketOption struct{
	WsUri			string	`json:"wsUri"`
	Port  			string	`json:"wsPort"`
	ListenerIp 		string	`json:"listenerIp"`
	OutIp 			string	`json:"outIp"`
	OutCxt 			context.Context
	Log				*zap.Logger
	//OpenNewConnBack	func ( connFD FDAdapter)	//来了新连接后，回调函数
	ProtocolManager *ProtocolManager	//给外部提供一个接口，用于将SOCKER FD 注册给外部
}

func NewWebsocket( option WebsocketOption)*Websocket{
	websocket := new(Websocket)
	websocket.Option = option
	//httpGin.GET(option.WsUri,websocket.HttpHandler)
	//ctx,cancelFunc := context.WithCancel(context.Background())

	dns := option.ListenerIp + ":" + option.Port

	websocket.httpServer = & http.Server{
		Addr:dns,
		//ErrorLog: logger,
	}

	return websocket
}
func (ws *Websocket )Shutdown(){
	ws.httpServer.Close()
}

func (ws *Websocket )Start(){
	http.HandleFunc(ws.Option.WsUri, ws.HttpHandler)
	//这里开始阻塞，直到接收到停止信号
	err := ws.httpServer.ListenAndServe()
	if err != nil {
		ws.Option.Log.Error("ws.httpServer.ListenAndServe()")
		if strings.Index(err.Error(),"Server closed") == -1{
			//zlib.PanicPrint("httpd:"+err.Error())
		}
		//mylog.Error(" httpd ListenAndServe err:", err.Error())
	}
}

func (ws *Websocket )HttpHandler(w http.ResponseWriter, r *http.Request){
	var httpUpGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有的CORS 跨域请求，正式环境可以关闭
		CheckOrigin: func(r *http.Request) bool {
			return true
		},

	}

	wsConnFD, err := httpUpGrader.Upgrade(w,r, nil)
	ws.Option.Log.Info("Upgrade this http req to websocket")
	if err != nil {
		ws.Option.Log.Error("Upgrade websocket failed: " + err.Error())
		return
	}

	ws.Option.ProtocolManager.websocketHandler(wsConnFD)
	//imp := WebsocketConnImpNew(wsConnFD)
	//ws.Option.OpenNewConnBack(imp)
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
