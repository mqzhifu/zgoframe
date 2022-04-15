package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitToolsRouter(Router *gin.RouterGroup) {
	ToolsRouter := Router.Group("tools")
	{
		//牛课网
		ToolsRouter.GET("niuke/question/dir/list", v1.NiukeQuestionDirList)
		//header头结构体 - 用于测试
		ToolsRouter.GET("header/struct", v1.HeaderStruct)
		//所有常量列表
		ToolsRouter.GET("const/list", v1.ConstList)
		//所有常量列表 - 生成MYSQL 脚本
		ToolsRouter.GET("const/init/db", v1.ConstInitDb)
		//获取APP 列表
		ToolsRouter.POST("project/list", v1.ProjectList)
		//获取APP 列表
		ToolsRouter.POST("project/info/{id}", v1.ProjectOneInfo)
		//获取APP 列表
		ToolsRouter.POST("frame/sync/js/demo", v1.FrameSyncJsDemo)
		//
		ToolsRouter.GET("test/init/db", v1.ConstInitTestDb)


	}
}
