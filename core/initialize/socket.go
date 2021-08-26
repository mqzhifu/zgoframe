package initialize

import (
	"encoding/json"
	"strconv"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	"zgoframe/util"
)

func initSocket(){
	global.V.ConnMng = util.NewConnManager(1000,5,global.V.RecoverGo,global.V.Zap,wsBackFunc)
	global.V.ConnPRotocol =util.NewConnProtocol(global.V.ProtobufMap,global.V.Zap)
	websocketOption :=util.WebsocketOption{
		WsUri: global.C.Websocket.Uri,
		Log: global.V.Zap,
		OpenNewConnBack :OpenNewConnBack,
	}
	global.V.Websocket = util.NewWebsocket(global.V.Gin,websocketOption)
	go global.V.ConnMng.Start()
}


func wsBackFunc(msg string,conn *util.Conn){
	//先把字符串转换成结构体
	connMsg,err  := global.V.ConnPRotocol.ParserContentProtocol(msg)
	if err != nil {

	}
	util.MyPrint("connMsg:",connMsg)
	//a := 0
	////最终的内容体，是有格式的，还得再解出来
	//err = global.V.ConnPRotocol.ParserContentMsg(connMsg,a)
	//if err != nil {
	//
	//}
}

type ResponseLoginRes struct {
	Code int
	Msg string
}

func OpenNewConnBack( connFD util.FDAdapter) {
	global.V.Zap.Info("OpenNewConnBack")
	//这里接收ws的连接fd
	NewConn, err := global.V.ConnMng.CreateOneConn(connFD)
	global.V.Zap.Debug("OpenNewConnBack CreateOneConn")
	if err != nil {
		global.V.Zap.Error("CreateOneConn:"+err.Error())
		//netWay.Option.Mylog.Error(err.Error())
		//NewConn.Write([]byte(err.Error()),websocket.TextMessage)
		//netWay.CloseOneConn(NewConn, CLOSE_SOURCE_CREATE)
		return
	}

	//defer func() {
	//	if err := recover(); err != nil {
	//		util.MyPrint("hit recover",err)
	//		global.V.Zap.Panic("OpenNewConn panic ")
	//		global.V.ConnMng.CloseOneConn(NewConn, util.CLOSE_SOURCE_OPEN_PANIC)
	//	}
	//}()

	var loginRes ResponseLoginRes
	//这里应该加一个超时机制
	content, err := NewConn.Read()
	global.V.Zap.Debug("first read :"+content)
	if err != nil {
		loginRes = ResponseLoginRes{
			Code: 500,
			Msg:  err.Error(),
		}
		global.V.ConnPRotocol.SendMsgCompressByConn(NewConn, "loginRes", loginRes)
		global.V.ConnMng.CloseOneConn(NewConn, util.CLOSE_SOURCE_FD_READ_EMPTY)
		return
	}
	//连接成功后，C端第一次发送消息给S端，内容是无格式的，一段字段串：(contentType + protocolType) + token
	header := []byte(content)[0:2]
	NewConn.ContentType, _ = strconv.Atoi(string(header[0]))
	NewConn.ProtocolType, _ = strconv.Atoi(string(header[1]))

	if NewConn.ContentType <= 0 || NewConn.ProtocolType <= 0 {
		loginRes = ResponseLoginRes{
			Code: 500,
			Msg:  "ContentType | ProtocolType <= 0",
		}

		loginResStr,_ := json.Marshal(loginRes)
		global.V.ConnPRotocol.Send(NewConn, string(loginResStr))

		global.V.ConnMng.CloseOneConn(NewConn, util.CLOSE_SOURCE_FD_READ_EMPTY)
		return
	}

	tokenBytes := []byte(content)[2:]
	requestHeader := request.Header {
		Token: string(tokenBytes),
		SourceType: 1,
	}

	global.V.Zap.Info("ContentType:"+ strconv.Itoa(NewConn.ContentType) + " , ProtocolType:"+ strconv.Itoa(NewConn.ProtocolType))
	parserTokenData, err := httpmiddleware.CheckToken(requestHeader)
	if err != nil {
		global.V.Zap.Error("CheckToken err:"+err.Error())
		loginRes = ResponseLoginRes{
			Code: 500,
			Msg:  "CheckToken err:"+err.Error(),
		}
		loginResStr,_ := json.Marshal(loginRes)
		global.V.ConnPRotocol.Send(NewConn, string(loginResStr))

		return
	}
	NewConn.PlayerId = parserTokenData.User.Id
	err = global.V.ConnMng.AddConnPool(NewConn)
	if err != nil {
		if err.Error() == "has-exist" {
			//err := errors.New(msg)
			responseKickOff := util.ResponseKickOff{
				Time: util.GetNowMillisecond(),
			}
			responseKickOffStr,_ := json.Marshal(responseKickOff)
			global.V.ConnPRotocol.Send(NewConn, string(responseKickOffStr))
		}

	}
	NewConn.IOLoop()
	//go ProcessMsgLoop(NewConn)
}




//func    RecoverProcessMsgLoop(ctx context.Context){
//	conn.Log.Warn("recover ReadLoop:")
//	go conn.ProcessMsgLoop(ctx)
//}

//func  ProcessMsgLoop(ctx context.Context,conn *util.Conn){
//	//defer func(ctx context.Context) {
//	//	if err := recover(); err != nil {
//	//		global.V.Zap.Panic("conn.ProcessMsgLoop panic defer :")
//	//		conn.RecoverProcessMsgLoop(ctx)
//	//	}
//	//}(ctx)
//
//	for{
//		ctxHasDone := 0
//		select{
//		case <-ctx.Done():
//			ctxHasDone = 1
//		case msgStr := <-conn.MsgInChan:
//			//global.V.Zap.Info("ProcessMsgLoop receive msg" + util.Action)
//			//myConnManager.BackFunc(msg,conn)
//		}
//		if ctxHasDone == 1{
//			goto end
//		}
//	}
//end :
//	global.V.Zap.Warn("ProcessMsgLoop receive signal: done.")
//}