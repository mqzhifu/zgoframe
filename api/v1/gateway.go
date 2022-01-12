package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/util"
)

func GatewayService(c *gin.Context) {
	prefix := "Gateway Http , "
	serviceName := c.Param("name")
	funcName := c.Param("func")

	data ,err  := c.GetRawData()
	if err != nil{
		util.ExitPrint(prefix + " GetRawData err:"+err.Error())
	}


	fmt.Println(prefix + " ServiceName:"+serviceName, " funcName:"+funcName + " data:"+string(data))
	util.ExitPrint(111)
	//geteway := util.NewGateway(global.V.GrpcManager,global.V.Zap)
	backData,err := global.V.Gateway.HttpCallGrpc(serviceName,funcName,"",data)
	if err != nil{

	}

	fmt.Println(backData)
}
