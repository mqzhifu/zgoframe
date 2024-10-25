package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
)

// @Tags Mail
// @Summary 发送一条站内信
// @Description 注意参数
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @accept application/json
// @Param data body request.SendMail true "参数信息,参考model"
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /mail/send [post]
func MailSend(c *gin.Context) {
	var form request.SendMail
	_ = c.ShouldBind(&form)
	// projectId, _ := request.GetProjectId(c)
	recordNewId, err := ApiServices().Mail.Send(form)
	if err != nil {
		httpresponse.FailWithMessage("失败了："+err.Error(), c)
	} else {
		httpresponse.OkWithMessage(strconv.Itoa(recordNewId), c)
	}
}

// @Tags Mail
// @Summary 获取用户站内信列表
// @Description 注意参数
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @accept application/json
// @Param data body request.MailList true "参数信息,参考model"
// @Success 200 {object} []model.MailLog
// @Router /mail/list [post]
func MailList(c *gin.Context) {
	var form request.MailList
	_ = c.ShouldBind(&form)
	// projectId, _ := request.GetProjectId(c)
	uid, _ := request.GetUid(c)
	mailList, err := ApiServices().Mail.GetUserListByUid(uid, form)
	if err != nil {
		httpresponse.FailWithMessage("失败了："+err.Error(), c)
	} else {
		httpresponse.OkWithAll(mailList, "ok", c)
	}
}

// @Tags Mail
// @Summary 获取用户一条信息的详情
// @Description 注意参数
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @accept application/json
// @Param data body request.MailInfo true "参数信息,参考model"
// @Success 200 {object} model.MailLog
// @Router /mail/info [post]
func MailInfo(c *gin.Context) {
	var form request.MailInfo
	_ = c.ShouldBind(&form)
	uid, _ := request.GetUid(c)
	info, err := ApiServices().Mail.GetOneByUid(uid, form)
	if err != nil {
		httpresponse.FailWithMessage("失败了："+err.Error(), c)
	} else {
		httpresponse.OkWithAll(info, "ok", c)
	}
}

// @Tags Mail
// @Summary 站内信未读总数
// @Description 注意参数
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @accept application/json
// @Success 200 {integer} int "155"
// @Router /mail/unread [get]
func MailUnread(c *gin.Context) {
	uid, _ := request.GetUid(c)
	cnt := ApiServices().Mail.GetUnreadCnt(uid)
	httpresponse.OkWithAll(cnt, "ok", c)
}
