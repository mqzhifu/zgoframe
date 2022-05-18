package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags System
// @Summary 关闭 - 该服务进程
// @Description 关闭 - 该服务进程
// @Security ApiKeyAuth
// @Param data body request.SystemConfig true "用户名/密码"
// @Produce  application/json
// @Success 200 {string} string "成功"
// @Router /sys/quit [post]
func Quit(c *gin.Context) {
	var formData request.SystemConfig
	c.ShouldBind(&formData)

	if formData.Username == "opendoor" && formData.Password == "123456" {
		global.V.Process.RootQuitFunc(2)
		global.V.Process.CancelFunc()
		httpresponse.OkWithAll(global.C, "结束中...", c)
	} else {
		httpresponse.FailWithMessage("验证失败", c)
	}

}

// @Summary 服务进程 - 配置信息
// @Description 服务进程 - 配置信息
// @Tags System
// @Security ApiKeyAuth
// @Param data body request.SystemConfig true "用户名/密码"
// @Produce  application/json
// @Success 200 {string} string "成功"
// @Router /sys/config [post]
func Config(c *gin.Context) {
	util.MyPrint("im in sys.config")

	var formData request.SystemConfig
	c.ShouldBind(&formData)

	util.MyPrint(formData)
	if formData.Username == "opendoor" && formData.Password == "123456" {
		httpresponse.OkWithAll(global.C, "结束中...", c)
	} else {
		httpresponse.FailWithMessage("验证失败", c)
	}
	//str,_ := json.Marshal(global.C)

}

// @Summary 标量- 实时统计信息 ,未实现
// @Description 标量- 实时统计信息
// @Tags System
// @Security ApiKeyAuth
// @Param data body request.SystemConfig true "用户名/密码"
// @Produce  application/json
// @Success 200 {string} string "成功"
// @Router /metrics [post]
func Metrics(c *gin.Context) {

}
