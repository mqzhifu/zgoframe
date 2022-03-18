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

type ParserTokenReq struct {
	Token string `json:"token" form:"token"`
}

// @Tags Base
// @Summary 解析一个TOKEN
// @Description 应用接到token后，要到后端再统计认证一下，确保准确
// @Param data body ParserTokenReq true "需要验证的token值"
// @Param X-Source-Type header string true "来源" default(1)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} httpresponse.Response
// @Router /base/parserToken [post]
func ParserToken(c *gin.Context) {
	var p ParserTokenReq
	c.ShouldBind(&p)

	//if err := util.Verify(sendSMS, util.SendSMSVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
	j := httpmiddleware.NewJWT()
	claims, err := j.ParseToken(p.Token)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	util.MyPrint("claims sourceType:", claims.SourceType)
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

// @Tags Base
// @Summary header头结构体
// @Description 日常header里放一诸如验证类的东西，统一公示出来，方便使用
// @Param X-Source-Type header string true "来源" default(1)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.Header
// @Router /base/headerStruct [get]
func RequestHeaderStruct(c *gin.Context) {

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

	if !request.CheckPlatformExist(request.GetMyHeader(c).SourceType) {
		httpresponse.FailWithMessage("Header.SourceType unknow", c)
		return
	}

	//if store.Verify(L.CaptchaId, L.Captcha, true) {
	//先从DB中做比对
	U := &model.User{Username: L.Username, Password: L.Password, ProjectId: L.AppId}
	err, user := service.Login(U)
	if err != nil {
		global.V.Zap.Error("登陆失败! 用户名不存在或者密码错误", zap.Any("err", err))
		httpresponse.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		//DB比较OK，开始做JWT处理
		tokenNext(c, *user)
	}
	//} else {
	//	httpresponse.FailWithMessage("验证码错误", c)
	//}
}

// @Summary 项目列表
// @Description 每个项目的详细信息
// @Security ApiKeyAuth
// @Tags Base
// @Param X-Source-Type header string true "来源" default(1)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /base/projectList [post]
func ProjectList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	rs := global.V.ProjectMng.Pool

	httpresponse.OkWithDetailed(rs, "成功", c)
}

// @Summary 项目的类型
// @Description 项目的类型
// @Security ApiKeyAuth
// @Tags Base
// @Param X-Source-Type header string true "来源" default(1)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/projectTypeList [get]
func ProjectTypeList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	httpresponse.OkWithDetailed(util.PROJECT_TYPE_MAP, "成功", c)
}

// @Summary 获取平台类型列表
// @Description 因为所有请求的hedaer 里必须得，所以动态获取
// @Security ApiKeyAuth
// @Tags Base
// @Param X-Source-Type header string true "来源" default(1)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/platformList [get]
func PlatformList(c *gin.Context) {
	httpresponse.OkWithDetailed(request.GetPlatformList(), "成功", c)
}
