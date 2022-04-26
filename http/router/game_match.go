package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitGameMatchRouter(Router *gin.RouterGroup) {
	gameMatchRouter := Router.Group("/game/match")
	{
		//
		gameMatchRouter.POST("sign", v1.GameMatchSign)

		gameMatchRouter.POST("sign/cancel", v1.GameMatchSignCancel)
	}
}
