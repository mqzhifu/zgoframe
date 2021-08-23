package router

import (
	"zgoframe/api/v1"
	"github.com/gin-gonic/gin"
	httpmiddleware "zgoframe/http/middleware"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(httpmiddleware.OperationRecord())
	{
		//UserRouter.POST("changePassword", v1.ChangePassword)     // 修改密码
		//UserRouter.POST("getUserList", v1.GetUserList)           // 分页获取用户列表
		//UserRouter.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户权限
		//UserRouter.DELETE("deleteUser", v1.DeleteUser)           // 删除用户
		//UserRouter.PUT("setUserInfo", v1.SetUserInfo)            // 设置用户信息
		UserRouter.PUT("logout", v1.Logout)            // 退出
	}
}


func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("user").Use(httpmiddleware.OperationRecord()).Use(httpmiddleware.RateMiddleware()).Use(httpmiddleware.ProcessHeader())
	{
		//登陆
		BaseRouter.POST("login", v1.Login)
		//BaseRouter.POST("login/sms", v1.Login)
		//BaseRouter.POST("login/third", v1.Login)
		//检查token正确性
		BaseRouter.POST("checktoken", v1.Checktoken)
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

