package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags Gateway
// @Summary 网关 - 短连接
// @Description 通过网关调取后端服务(grpc)
// @Security ApiKeyAuth
// @Param service_name path string true "服务名"
// @Param func_name path string true "函数名"
// @Param data body request.Empty true "任意，请参考.proto"
// @Success 200 {bool} bool "true:成功 false:否"
// @Router /gateway/service/{service_name}/{func_name} [post]
func GatewayService(c *gin.Context) {
	prefix := "Gateway Http , "
	serviceName := c.Param("service_name")
	funcName := c.Param("func_name")

	data, err := c.GetRawData()
	util.MyPrint("c.GetRawData data:", data, " err:", err)
	//if err != nil {
	//	errMsg := prefix + " GetRawData err:" + err.Error()
	//	//util.ExitPrint(prefix + " GetRawData err:"+err.Error())
	//	httpresponse.FailWithMessage(errMsg, c)
	//	c.Abort()
	//}

	fmt.Println(prefix+" ServiceName:"+serviceName, " funcName:"+funcName+" data:"+string(data))
	backData, err := global.V.Gateway.HttpCallGrpc(serviceName, funcName, "", data)
	if err != nil {
		fmt.Println(err)
		httpresponse.FailWithMessage(err.Error(), c)
		c.Abort()
	}

	fmt.Println(backData)
}

// @Tags Gateway
// @Summary 获取网关配置信息
// @Description 主要是长连接的配置(端口|协议)
// @Security ApiKeyAuth
// @Success 200 {object} util.NetWayOption
// @Router /gateway/config [get]
func GatewayConfig(c *gin.Context) {
	httpresponse.OkWithDetailed(global.V.Gateway.NetWayOption, "ok", c)

}

// @Tags Gateway
// @Summary 获取所有服务的.proto 配置文件
// @Description proto接口及GRPC微服务函数的信息等
// @Security ApiKeyAuth
// @Router /gateway/proto [get]
func GatewayProto(c *gin.Context) {

}

// @Tags Gateway
// @Summary 网关 - 长连接 websocket
// @Description 通过网关调取后端服务(ws)
// @Security ApiKeyAuth
// @Router /gateway/service/ws [get]
func GatewayWS(c *gin.Context) {

}
