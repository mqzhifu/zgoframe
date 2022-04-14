package router

import (
	"github.com/gin-gonic/gin"
	v1 "zgoframe/api/v1"
)

func InitMailRouter(Router *gin.RouterGroup) {
	mailRouter := Router.Group("mail")
	{
		//发送站一条站内信
		mailRouter.POST("send", v1.MailSend)
		//站内信列表
		mailRouter.POST("list", v1.MailList)
		//一条站内信详情
		mailRouter.POST("info", v1.MailInfo)
		//未读站内信总数
		mailRouter.GET("unread", v1.MailUnread)
	}
}
