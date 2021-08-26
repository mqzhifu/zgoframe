package router

import (
	"github.com/gin-gonic/gin"
	"zgoframe/api/v1"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		//UserRouter.POST("changePassword", v1.ChangePassword)     // 修改密码
		//UserRouter.POST("getUserList", v1.GetUserList)           // 分页获取用户列表
		//UserRouter.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户权限
		//UserRouter.DELETE("deleteUser", v1.DeleteUser)           // 删除用户
		//UserRouter.PUT("setUserInfo", v1.SetUserInfo)            // 设置用户信息
		UserRouter.PUT("logout", v1.Logout)            // 退出
		//检查token正确性
		UserRouter.POST("checktoken", v1.Checktoken)
	}
}

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		//检查token正确性
		BaseRouter.POST("parserToken", v1.ParserToken)
		//登陆
		BaseRouter.POST("login", v1.Login)
		//BaseRouter.POST("login/sms", v1.Login)
		//BaseRouter.POST("login/third", v1.Login)
		//验证码
		BaseRouter.POST("captcha", v1.Captcha)
		//注册
		BaseRouter.POST("register", v1.Register)
		//BaseRouter.POST("register/sms", v1.Register)
		//BaseRouter.POST("register/third", v1.Register)
		//发送短信
		//BaseRouter.POST("sendsms", v1.SendSMS)
	}
}

func InitLogslaveRouter(Router *gin.RouterGroup) {
	LogsalveRouter := Router.Group("logslave")
	{
		LogsalveRouter.POST("receive", v1.Receive)
		//ws 连接信息
		LogsalveRouter.GET("getWsServer", v1.GetWsServer)
		//
		LogsalveRouter.GET("getApiList", v1.GetApiList)
	}
}


