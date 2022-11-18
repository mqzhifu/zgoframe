package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func GameMatch(Router *gin.RouterGroup) {
	gameMatchRouter := Router.Group("/game/match")
	{
		gameMatchRouter.POST("sign", v1.GameMatchSign)
		gameMatchRouter.POST("sign/cancel", v1.GameMatchSignCancel)
		gameMatchRouter.GET("rule/:id", v1.GameMatchGetOneRule)
		gameMatchRouter.GET("lang", v1.GameMatchGetLang)
		gameMatchRouter.GET("config", v1.GameMatchConfig)

	}
}

func FrameSync(Router *gin.RouterGroup) {
	gameMatchRouter := Router.Group("/frame/sync")
	{
		gameMatchRouter.GET("config", v1.FrameSyncConfig)
	}
}
