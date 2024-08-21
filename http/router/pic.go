package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func Pic(Router *gin.RouterGroup) {
	picRouter := Router.Group("pic")
	{
		//testRouter.POST("binary/tree", v1.BinaryTree)
		picRouter.POST("split", v1.Split)
	}
}
