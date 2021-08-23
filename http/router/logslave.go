package router

import (
	"zgoframe/core/global"
	"zgoframe/service"
)

func logsalve(){
	global.V.Gin.GET("/receive",service.Receive)
}
