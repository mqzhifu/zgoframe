package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)


func InitPersistenceRouter(Router *gin.RouterGroup) {
	persistenceRouter := Router.Group("persistence")
	{
		persistenceRouter.POST("log/push", v1.LogPush)
		persistenceRouter.POST("file/upload", v1.Upload)
	}
}