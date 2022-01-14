package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

func GatewayService(c *gin.Context) {
	prefix := "Gateway Http , "
	serviceName := c.Param("name")
	funcName := c.Param("func")

	data ,err  := c.GetRawData()
	if err != nil{
		errMsg := prefix + " GetRawData err:"+err.Error()
		//util.ExitPrint(prefix + " GetRawData err:"+err.Error())
		httpresponse.FailWithMessage(errMsg,c)
		c.Abort()
	}

	fmt.Println(prefix + " ServiceName:"+serviceName, " funcName:"+funcName + " data:"+string(data))
	backData,err := global.V.Gateway.HttpCallGrpc(serviceName,funcName,"",data)
	if err != nil{
		fmt.Println(err)
		httpresponse.FailWithMessage(err.Error(),c)
		c.Abort()
	}

	fmt.Println(backData)
}
