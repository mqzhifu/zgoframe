package util

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
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
	IOTime 			int64
	MsgContentMax	int32
	OutCxt 			context.Context
	Log				*zap.Logger
	ProtocolManager *ProtocolManager	//给外部提供一个接口，用于将SOCKER FD 注册给外部
}

func NewWebsocket( option WebsocketOption)*Websocket{
	websocket := new(Websocket)
	websocket.Option = option

	dns := option.ListenerIp + ":" + option.Port

	websocket.httpServer = & http.Server{
		Addr:dns,
		//ErrorLog: logger,
	}

	return websocket
}
//关闭
func (ws *Websocket )Shutdown(){
	//直接关闭HTTP 守护协程即可，也不需要断开连接
	ws.httpServer.Close()
}
//启动http 监听
func (ws *Websocket )Start(){
	ws.Option.Log.Info("start websocket  dns:"+ws.httpServer.Addr + " uri:"+ ws.Option.WsUri)
	http.HandleFunc(ws.Option.WsUri, ws.HttpHandler)
	//这里开始阻塞，直到接收到停止信号（shutdown会触发）
	err := ws.httpServer.ListenAndServe()
	if err != nil {
		ws.Option.Log.Error("ws.httpServer.ListenAndServe()")
	}
}
//所有的HTTP-WS请求会回调此函数
func (ws *Websocket )HttpHandler(w http.ResponseWriter, r *http.Request){
	//配置ws
	var httpUpGrader = websocket.Upgrader{
		//HandshakeTimeout: ws.Option.IOTime,//read write 超时时间
		//下面这两个是创建buffer大小 ，间接等于一条消息的大小，这里不做限制了，直接交给上层再做处理
		//ReadBufferSize: int( ws.Option.MsgContentMax),
		//WriteBufferSize: int( ws.Option.MsgContentMax),

		// 允许所有的CORS 跨域请求，正式环境可以关闭
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//将此HTTP请求升级到ws格式
	wsConnFD, err := httpUpGrader.Upgrade(w,r, nil)
	ws.Option.Log.Info("ws HttpHandler Upgrade this http req to websocket ,http remote:"+r.RemoteAddr)
	if err != nil {
		ws.Option.Log.Error("Upgrade websocket failed: " + err.Error())
		return
	}
	myMetrics.CounterInc("ws_ok_fd")
	//将ws fd 回调给更上层，做更多操作
	ws.Option.ProtocolManager.websocketHandler(wsConnFD)
}