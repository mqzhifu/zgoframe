package v1

import (
	"github.com/gin-gonic/gin"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Summary checktoken
// @Description checktoken
// @Tags Base
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/checktoken [post]
func Checktoken(c *gin.Context ) {
	httpmiddleware.CheckToken(request.GetMyHeader(c))
}


// @Summary 发送验证码
// @Description 登陆、注册、通知等发送短信
// @Tags Base
// @Produce  application/json
// @Param data body request.SendSMS true "手机号, 规则ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /base/sendSMS [post]
func SendSMS(c *gin.Context) {
	var sendSMS request.SendSMS
	c.ShouldBind(&sendSMS)
	if err := util.Verify(sendSMS, util.SendSMSVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
}