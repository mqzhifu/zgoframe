package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
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

	if CheckID(formData) {
		global.V.Process.RootQuitFunc(2)
		httpresponse.OkWithAll(global.C, "信号已发出，结束中...请等待几秒", c)
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


	var formData request.SystemConfig
	c.ShouldBind(&formData)
	//这里也可以把global.C输出回去
	//global.C
	info := global.V.Process.InitBaseInfoCallbackFunc()
	//util.MyPrint("InitBaseInfoCallbackFunc:",info)
	//
	if  CheckID(formData)  {
		httpresponse.OkWithAll(info, "结束中...", c)
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
	//此方法主要是使用注解，生成文档给开发查看，实际在框架的初始化阶段，由GIN拦截了
}

func CheckID(form request.SystemConfig)bool{
	if form.Username == "opendoor" && form.Password == "123456" {
		return true
	}else{
		return false
	}
}
