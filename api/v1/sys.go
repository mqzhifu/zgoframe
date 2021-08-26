package v1

import (
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Summary Quit
// @Description Quit
// @Security ApiKeyAuth
// @Tags User
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /sys/quit [POST]
func Quit(c *gin.Context ) {
	SysCheck()
	global.V.Process.RootQuitFunc(2)
	//httpresponse.OkWithDetailed("", "结束中...", c)
	util.MyPrint("CancelFunc")
	global.V.Process.CancelFunc()
	util.MyPrint("CancelFunc finish")
}

// @Summary Config
// @Description Config
// @Security ApiKeyAuth
// @Tags User
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /sys/config [POST]
func Config(c *gin.Context ) {
	SysCheck()
	//str,_ := json.Marshal(global.C)
	httpresponse.OkWithDetailed(global.C, "结束中...", c)
}

func SysCheck(){

}