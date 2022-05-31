package router

import (
	"github.com/gin-gonic/gin"
	"zgoframe/api/v1"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.POST("set/password", v1.SetPassword) // 设置/修改密码
		//UserRouter.POST("getUserList", v1.GetUserInfoList) // 分页获取用户列表
		//UserRouter.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户权限
		UserRouter.DELETE("delete", v1.DeleteUser) // 删除用户
		UserRouter.PUT("set/mobile", v1.SetMobile) //绑定手机号
		UserRouter.PUT("set/email", v1.SetEmail)   //绑定邮箱

		UserRouter.POST("set/info", v1.SetUserInfo) // 设置用户信息
		UserRouter.PUT("logout", v1.Logout)         // 退出
		UserRouter.GET("info", v1.GetUserInfo)
	}
}

func InitSysRouter(Router *gin.RouterGroup) {
	SysRouter := Router.Group("sys")
	{
		//
		SysRouter.POST("quit", v1.Quit)
		//
		SysRouter.POST("config", v1.Config)
		//
		SysRouter.POST("metrics", v1.Metrics)

	}
}

func InitGatewayRouter(Router *gin.RouterGroup) {
	GatewayRouter := Router.Group("gateway")
	{
		GatewayRouter.POST("service/:service_name/:func_name", v1.GatewayService)
		GatewayRouter.GET("proto", v1.GatewayProto)
		GatewayRouter.GET("action/map", v1.ActionMap)
		GatewayRouter.GET("config", v1.GatewayConfig)
	}
}

