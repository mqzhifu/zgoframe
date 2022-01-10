package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
)

func GatewayService(c *gin.Context) {
	serviceName := c.Param("name")
	funcName := c.Param("func")

	data ,_ := c.GetRawData()

	//geteway := util.NewGateway(global.V.GrpcManager,global.V.Zap)
	backData,err := global.V.Gateway.HttpCallGrpc(serviceName,funcName,"",data)
	if err != nil{

	}

	fmt.Println(backData)
}
