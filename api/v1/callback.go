package v1

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags Callback
// @Summary 声网 - 回调
// @Description 订阅什么事件就回调什么事件
// @Security ApiKeyAuth
// @Param data body request.SystemConfig true "用户名/密码"
// @Produce  application/json
// @Success 200 {string} string "成功"
// @Router /callback/agora/rtc [post]
func AgoraCallbackRTC(c *gin.Context) {
	util.MyPrint("header:", c.Request.Header)
	util.MyPrint("header-foreach:")
	for k, v := range c.Request.Header {
		util.MyPrint(k, v)
	}
	util.MyPrint("=======================")
	util.MyPrint("url:", c.Request.URL)
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	util.MyPrint("ReadAll body:", string(bodyBytes), " err:", err)

	util.MyPrint("im in callback")
	httpresponse.OkWithAll("回调成功", "ok", c)

}
