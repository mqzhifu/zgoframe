package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitConfigCenterRouter(Router *gin.RouterGroup) {
	configCenterRouter := Router.Group("config/center")
	{
		//以模块(文件)为单位，获取该模块(文件)下的所有配置信息
		configCenterRouter.POST("get/module", v1.ConfigCenterGetByModule)
		//以以模块(文件)+里面具体的key 为单位，获取配置信息
		configCenterRouter.POST("get/key", v1.ConfigCenterGetByModuleByKey)
		//以模块(文件)+里面具体的key 为单位，设置置信息(如果存在，覆盖)
		configCenterRouter.POST("set/key", v1.ConfigCenterSetByModuleByKey)
		//创建模块(文件)
		configCenterRouter.POST("create/module", v1.ConfigCenterCreateModule)
	}
}
