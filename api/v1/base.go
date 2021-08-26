package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/util"
)

// @Summary checktoken
// @Description checktoken
// @Security ApiKeyAuth
// @Tags User
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/checktoken [post]
func Checktoken(c *gin.Context ) {
	//httpresponse.OkWithDetailed(httpresponse.LoginResponse{
	//	Token:     myHeader.Token,
	//}, "检测成功", c)
}

type ParserTokenReq struct {
	Token string	`json:"token" form:"token"`
}

// @Summary ParserToken
// @Description ParserToken
// @Security ApiKeyAuth
// @Tags Base
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/parserToken [post]
func ParserToken(c *gin.Context ) {
	var  p ParserTokenReq
	c.ShouldBind(&p)

	//if err := util.Verify(sendSMS, util.SendSMSVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	j := httpmiddleware.NewJWT()
	claims, err := j.ParseToken(p.Token)
	if err != nil{
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	util.MyPrint("claims sourceType:",claims.SourceType)
	if err != nil {
		if err == httpmiddleware.TokenExpired {
			httpresponse.FailWithMessage("授权已过期", c)
			return
		}
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	httpresponse.OkWithDetailed(claims, "解析成功", c)
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
// @Summary 用户登陆
// @Description 用户登陆，验证，生成token
// @Tags Base
// @Produce  application/json
// @Param data body request.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/login [post]
func Login(c *gin.Context) {
	var L request.Login
	c.ShouldBind(&L)
	if err := util.Verify(L, util.LoginVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	if !request.CheckPlatformExist(request.GetMyHeader(c).SourceType){
		httpresponse.FailWithMessage("Header.SourceType unknow", c)
		return
	}

	//if store.Verify(L.CaptchaId, L.Captcha, true) {
	U := &model.User{Username: L.Username, Password: L.Password,AppId: L.AppId}
	err, user := service.Login(U)
	if  err != nil {
		global.V.Zap.Error("登陆失败! 用户名不存在或者密码错误", zap.Any("err", err))
		httpresponse.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		tokenNext(c, *user)
	}
	//} else {
	//	httpresponse.FailWithMessage("验证码错误", c)
	//}
}

