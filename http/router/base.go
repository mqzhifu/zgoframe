package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		//header头结构体 - 用于测试
		BaseRouter.GET("niuke/question/dir/list", v1.NiukeQuestionDirList)
		//header头结构体 - 用于测试
		BaseRouter.GET("header/struct", v1.HeaderStruct)
		//所有常量列表
		BaseRouter.GET("const/list", v1.ConstList)
		//图形 - 验证码
		BaseRouter.GET("captcha", v1.Captcha)
		//发送短信 登陆/注册/找回密码
		BaseRouter.POST("send/sms", v1.SendSms)
		//发送邮件 登陆/注册/找回密码
		BaseRouter.POST("send/email", v1.SendEmail)
		//登陆 - 用户名/密码
		BaseRouter.POST("login", v1.Login)
		//登陆 - 三方平台(无密码，且可自动注册)
		BaseRouter.POST("login/third", v1.LoginThird)
		//登陆 - 短信(无密码)
		BaseRouter.POST("login/sms", v1.LoginSms)
		//注册 - 用户名/密码
		BaseRouter.POST("register", v1.Register)
		//注册 - 短信
		BaseRouter.POST("register/sms", v1.RegisterSms)
		//重置密码 - 通过短信
		BaseRouter.POST("sms/reset/password", v1.ResetPasswordSms)
		//检查手机号是否存在 登陆/注册/找回密码
		BaseRouter.POST("check/mobile", v1.CheckMobileExist)
		//检查邮件是否存在 登陆/注册/找回密码
		BaseRouter.POST("check/email", v1.CheckEmailExist)
		//获取APP 列表
		BaseRouter.POST("project/list", v1.ProjectList)
		//检查token正确性
		BaseRouter.POST("parser/token", v1.ParserToken)

	}
}
