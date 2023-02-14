package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func Tools(Router *gin.RouterGroup) {
	ToolsRouter := Router.Group("tools")
	{
		//获取 - httpHeader 结构体
		ToolsRouter.GET("header/struct", v1.HeaderStruct)
		//所有常量列表
		ToolsRouter.GET("const/list", v1.ConstList)
		//所有常量列表 - 生成MYSQL 脚本
		ToolsRouter.GET("const/init/db", v1.ConstInitDb)
		//获取APP 列表
		ToolsRouter.POST("project/list", v1.ProjectList)
		//获取APP 列表
		ToolsRouter.GET("project/info/:id", v1.ProjectOneInfo)
		ToolsRouter.GET("init/db/structure", v1.InitDbStructure)
		ToolsRouter.GET("init/db/data", v1.InitDbData)
		//生成 MYSQL 测试 数据
		ToolsRouter.GET("test/init/db", v1.ConstInitTestDb)

		ToolsRouter.GET("test/migu/api/para", v1.TestMiguAPIGetPara)
		ToolsRouter.GET("test/migu/api/back/data ", v1.ReceiveMiguBackData)

	}
}
