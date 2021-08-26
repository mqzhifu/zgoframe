package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)


type LogData struct {
	AppId 	int		`json:"app_id" form:"app_id"`
	Uid 	int		`json:"uid" form:"uid"`
	Level 	int		`json:"level" form:"level"`
	Msg 	string	`json:"msg" form:"msg"`

}

// @Tags Logslave
// @Summary 接收日志
// @Produce  application/json
// @Param data body v1.LogData true "level,msg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"已收录"}"
// @Router /logslave/receive [post]
func Receive(c *gin.Context) {
	var L LogData
	_ = c.ShouldBind(&L)
	if err := util.Verify(L, util.LogReceiveVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	L.Uid ,_ = request.GetUid(c)
	L.AppId,_ =request.GetAppId(c)


	str,_ := json.Marshal(L)
	global.V.Zap.Info(string(str))

	httpresponse.OkWithDetailed("", "已收录", c)
}

type WsServer struct {
	Env		string	`json:"env"`
	Ip		string	`json:"ip"`
	Port 	string	`json:"port"`
	Uri 	string	`json:"uri"`
}

// @Summary ParserToken
// @Description ParserToken
// @Security ApiKeyAuth
// @Tags Logslave
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /logslave/getWsServer [GET]
func GetWsServer(c *gin.Context ) {
	//var L LogData
	//_ = c.ShouldBind(&L)
	//if err := util.Verify(L, util.LogReceiveVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}

	w := WsServer{
		Env: "",
		Ip : "127.0.0.1",
		Port : "5555",
		Uri :global.C.Websocket.Uri,
	}
	//httpresponse.OkWithDetailed(httpresponse.SysUserResponse{User: userReturn}, "注册成功", c)
	httpresponse.OkWithDetailed(w, "成功", c)
	return
}

// @Summary 获取长连接API映射表
// @Description 获取长连接API映射表
// @Security ApiKeyAuth
// @Tags Logslave
// @Produce  application/json
// @Router /logslave/getWsServer [GET]
func GetApiList(c *gin.Context){
	actionMap := global.V.ProtobufMap.GetActionMap()
	httpresponse.OkWithDetailed(actionMap, "成功", c)
}

func WsNewFdCallback(connFD util.FDAdapter){

}