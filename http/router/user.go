package router

import (
	"github.com/gin-gonic/gin"
	"zgoframe/api/v1"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.POST("changePassword", v1.ChangePassword) // 修改密码
		//UserRouter.POST("getUserList", v1.GetUserInfoList) // 分页获取用户列表
		//UserRouter.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户权限
		//UserRouter.DELETE("deleteUser", v1.DeleteUser)           // 删除用户
		UserRouter.PUT("setUserInfo", v1.SetUserInfo) // 设置用户信息
		UserRouter.PUT("logout", v1.Logout)           // 退出
		UserRouter.GET("getUserInfo", v1.GetUserInfo)
	}
}

func InitSysRouter(Router *gin.RouterGroup) {
	SysRouter := Router.Group("sys")
	{
		//
		SysRouter.POST("quit", v1.Quit)
		//
		SysRouter.POST("config", v1.Config)
	}
}

func InitGatewayRouter(Router *gin.RouterGroup) {
	GatewayRouter := Router.Group("service")
	{
		GatewayRouter.POST(":name/:func", v1.GatewayService)
		GatewayRouter.GET("getConfig", v1.GatewayService)
	}
}

func InitLogslaveRouter(Router *gin.RouterGroup) {
	LogsalveRouter := Router.Group("logslave")
	{
		LogsalveRouter.POST("push", v1.Push)
	}
}
