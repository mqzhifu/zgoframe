package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func Test(Router *gin.RouterGroup) {
	testRouter := Router.Group("test")
	{
		//testRouter.POST("binary/tree", v1.BinaryTree)
		testRouter.GET("binary/tree/list/:flag", v1.BinaryTreeListByFlag)

		testRouter.GET("binary/tree/insert/one/:keyword", v1.BinaryTreeInsertOne)
		testRouter.GET("binary/tree/each/deep", v1.BinaryTreeEachDeepByBreadthFirst)

	}
}
