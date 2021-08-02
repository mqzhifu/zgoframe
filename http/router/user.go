package router

import (
	"zgoframe/api/v1"
	"github.com/gin-gonic/gin"
	httpmiddleware "zgoframe/http/middleware"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(httpmiddleware.OperationRecord())
	{
		UserRouter.POST("register", v1.Register)
		UserRouter.POST("changePassword", v1.ChangePassword)     // 修改密码
		UserRouter.POST("getUserList", v1.GetUserList)           // 分页获取用户列表
		UserRouter.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户权限
		UserRouter.DELETE("deleteUser", v1.DeleteUser)           // 删除用户
		UserRouter.PUT("setUserInfo", v1.SetUserInfo)            // 设置用户信息
	}
}


func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("user").Use(httpmiddleware.OperationRecord())
	{
		BaseRouter.POST("login", v1.Login)
		BaseRouter.POST("captcha", v1.Captcha)
	}
}

