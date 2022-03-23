package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Summary Quit
// @Description 关闭该服务进程
// @Security ApiKeyAuth
// @Tags System
// @Param data body request.SystemConfig true "用户名/密码"
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /sys/quit [POST]
func Quit(c *gin.Context) {
	var formData request.SystemConfig
	c.ShouldBind(&formData)

	if formData.Username == "opendoor" && formData.Password == "123456" {
		global.V.Process.RootQuitFunc(2)
		global.V.Process.CancelFunc()
		httpresponse.OkWithDetailed(global.C, "结束中...", c)
	} else {
		httpresponse.FailWithMessage("验证失败", c)
	}

}

// @Summary Config
// @Description Config
// @Tags System
// @Security ApiKeyAuth
// @Param data body request.SystemConfig true "用户名/密码"
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /sys/config [POST]
func Config(c *gin.Context) {
	util.MyPrint("im in sys.config")

	var formData request.SystemConfig
	c.ShouldBind(&formData)

	util.MyPrint(formData)
	if formData.Username == "opendoor" && formData.Password == "123456" {
		httpresponse.OkWithDetailed(global.C, "结束中...", c)
	} else {
		httpresponse.FailWithMessage("验证失败", c)
	}
	//str,_ := json.Marshal(global.C)

}
