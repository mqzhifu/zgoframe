package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/util"
)

// @Tags Base
// @Summary header头结构体
// @Description 日常header里放一诸如验证类的东西，统一公示出来，仅是说明，方便测试，不是真实API，方便使用
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Success 200 {object} request.Header
// @Router /base/header/struct [get]
func HeaderStruct(c *gin.Context) {
	myheader := request.Header{}
	httpresponse.OkWithDetailed(myheader, "成功lalalalala", c)
}

var store = base64Captcha.DefaultMemStore

// @Tags Base
// @Summary 生成图片验证码
// @Description BASE64图片内容，防止有人恶意攻击，如：短信轰炸、暴力破解密码等
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {object} httpresponse.SysCaptchaResponse
// @Router /base/captcha [get]
func Captcha(c *gin.Context) {
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.C.Captcha.ImgHeight, global.C.Captcha.ImgWidth, global.C.Captcha.NumberLength, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		global.V.Zap.Error("验证码获取失败!", zap.Any("err", err))
		httpresponse.FailWithMessage("验证码获取失败", c)
	} else {
		httpresponse.OkWithDetailed(httpresponse.SysCaptchaResponse{
			Id:         id,
			PicContent: b64s,
		}, "验证码获取成功", c)
	}
}

// @Summary 发送短信
// @Description 登陆/注册/找回密码
// @Tags Base
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.SendSMS true "用户信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:否"
// @Router /base/send/sms [post]
func SendSms(c *gin.Context) {
	var sendSMSForm request.SendSMS
	c.ShouldBind(&sendSMSForm)
	if sendSMSForm.SendIp == "" {
		sendSMSForm.SendIp = c.Request.RemoteAddr
	}
	//
	projectId, _ := request.GetProjectId(c)
	err := global.V.MyService.SendSms.Send(projectId, sendSMSForm)
	if err != nil {
		httpresponse.FailWithMessage("失败了："+err.Error(), c)
	} else {
		httpresponse.OkWithMessage("成功喽~", c)
	}
	//if err != nil {
	//
	//}
}

// @Tags Base
// @Summary 重置密码 - 通过短信
// @Description 忘记密码后，可发送短信通知，重置密码
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.RestPasswordSms true "用户名, 原密码, 新密码"
// @Success 200 {bool} bool "true:成功 false:否"
// @Router /base/sms/reset/password [post]
func ResetPasswordSms(c *gin.Context) {
	var form request.RestPasswordSms
	_ = c.ShouldBindJSON(&form)

	if form.NewPassword == "" || form.NewPasswordConfirm == "" {
		httpresponse.OkWithMessage("NewPassword |NewPasswordConfirm empty", c)
		return
	}

	if form.NewPassword != form.NewPasswordConfirm {
		httpresponse.OkWithMessage("密码与确认密码不一致", c)
		return
	}

	err := global.V.MyService.SendSms.Verify(form.SmsRuleId, form.Mobile, form.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	uid, _ := request.GetUid(c)
	err = global.V.MyService.User.ChangePassword(uid, form.NewPassword)
	if err != nil {
		httpresponse.FailWithMessage("修改失败:"+err.Error(), c)
	} else {
		httpresponse.OkWithMessage("修改成功", c)
	}
}

// @Tags Base
// @Summary 用户注册账号
// @Description 普通注册，需要有：用户名 密码
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Base-Info header string false "客户端基础信息(json格式,参考request.HeaderBaseInfo)"
// @Param data body request.Register true "用户信息"
// @Success 200 {object} model.User
// @Router /base/register [post]
func Register(c *gin.Context) {
	var R request.Register
	_ = c.ShouldBind(&R)
	if err := util.Verify(R, util.RegisterVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	if R.Password != R.ConfirmPs {
		httpresponse.FailWithMessage("确认密码与密码不同", c)
		return
	}

	err, userInfo := global.V.MyService.User.RegisterByUsername(R, request.GetMyHeader(c))
	if err != nil {
		//global.V.Zap.Error("注册失败", zap.Any("err", err))
		httpresponse.FailWithDetailed(userInfo, "注册失败:"+err.Error(), c)
	} else {
		httpresponse.OkWithDetailed(userInfo, "注册成功", c)
	}
}

// @Tags Base
// @Summary 用户注册账号-通过手机验证码
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Base-Info header string false "客户端基础信息(json格式,参考request.HeaderBaseInfo)"
// @Param data body request.RegisterSms true "用户信息"
// @Success 200 {object} model.User
// @Router /base/register/sms [post]
func RegisterSms(c *gin.Context) {
	var registerSmsForm request.RegisterSms
	_ = c.ShouldBind(&registerSmsForm)

	user := model.User{
		Username: registerSmsForm.Mobile,
		Mobile:   registerSmsForm.Mobile,
		Guest:    model.USER_GUEST_FALSE,
	}

	err := global.V.MyService.SendSms.Verify(registerSmsForm.SmsRuleId, registerSmsForm.Mobile, registerSmsForm.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	err, userInfo := global.V.MyService.User.Register(user, request.GetMyHeader(c), service.UserRegInfo{})
	if err != nil {
		//global.V.Zap.Error("注册失败", zap.Any("err", err))
		httpresponse.FailWithDetailed(userInfo, "注册失败:"+err.Error(), c)
		return
	}

	httpresponse.OkWithDetailed(userInfo, "注册成功", c)

}

// @Tags Base
// @Summary 测试手机号：是否已存在，注册/绑定时会使用
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param data body request.CheckMobileExist true "用户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Success 200 {bool} bool "true:存在 false:不存在"
// @Router /base/check/mobile [post]
func CheckMobileExist(c *gin.Context) {
	var form request.CheckMobileExist
	_ = c.ShouldBind(&form)

	if form.Mobile == "" || form.Purpose <= 0 {
		httpresponse.FailWithMessage("form.Mobile == '' || form.Purpose <= 0", c)
		return
	}

	if !util.CheckMobileRule(form.Mobile) {
		httpresponse.FailWithMessage("mobile 格式 错误 ", c)
		return
	}

	_, empty, err := global.V.MyService.User.FindUserByMobile(form.Mobile)
	if err != nil {
		httpresponse.FailWithMessage("服务器错误，请等待或重试", c)
	} else {
		if !empty {
			httpresponse.FailWithDetailed(true, "成功", c)
		} else {
			httpresponse.FailWithDetailed(false, "成功", c)
		}
	}

}

// @Tags Base
// @Summary 测试邮件：是否已存在，注册/绑定时会使用
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CheckEmailExist true "用户信息"
// @Success 200 {bool} bool "true:存在 false:不存在"
// @Router /base/check/mobile [post]
func CheckEmailExist(c *gin.Context) {
	var form request.CheckMobileExist
	_ = c.ShouldBind(&form)

	if form.Mobile == "" || form.Purpose <= 0 {
		httpresponse.FailWithMessage("form.Mobile == '' || form.Purpose <= 0", c)
		return
	}

	if !util.CheckEmailRule(form.Mobile) {
		httpresponse.FailWithMessage("email 格式 错误 ", c)
		return
	}

	_, empty, err := global.V.MyService.User.FindUserByEmail(form.Mobile)
	if err != nil {
		httpresponse.FailWithMessage("服务器错误，请等待或重试", c)
	} else {
		if !empty {
			httpresponse.FailWithDetailed(true, "成功", c)
		} else {
			httpresponse.FailWithDetailed(false, "成功", c)
		}
	}

}

// @Tags Base
// @Summary 解析一个TOKEN
// @Description 应用接到token后，要到后端再统计认证一下，确保准确
// @Param data body request.ParserToken true "需要验证的token值"
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.CustomClaims
// @Router /base/parser/token [post]
func ParserToken(c *gin.Context) {
	//util.MyPrint("im in parserToken")
	var p request.ParserToken
	c.ShouldBind(&p)

	j := httpmiddleware.NewJWT()
	claims, err := j.ParseToken(p.Token)
	if err != nil {
		httpresponse.FailWithMessage("解析失败:"+err.Error(), c)
		return
	} else {
		httpresponse.OkWithDetailed(claims, "解析成功", c)
	}
	//util.MyPrint("claims sourceType:", claims.SourceType)
	//if err != nil {
	//	if err == httpmiddleware.TokenExpired {
	//		httpresponse.FailWithMessage("授权已过期", c)
	//		return
	//	}
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}

}

// @Summary 发送邮件
// @Description 登陆、注册、通知等发送
// @Tags Base
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.SendEmail true "基础信息"
// @Produce  application/json
// @Success 200 {bool} bool "true:成功 false:否"
// @Router /base/send/email [post]
func SendEmail(c *gin.Context) {
}

// @Summary 用户登陆
// @Description 用户名(手机邮箱)/密码登陆，验证成功后，生成token
// @Tags Base
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.Login true "用户名, 密码, 验证码"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/login [post]
func Login(c *gin.Context) {
	var L request.Login
	c.ShouldBind(&L)
	if err := util.Verify(L, util.LoginVerify); err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	//if store.Verify(L.CaptchaId, L.Captcha, true) {
	//先从DB中做比对
	U := &model.User{Username: L.Username, Password: L.Password}
	err, user := global.V.MyService.User.Login(U)
	if err != nil {
		httpresponse.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		//DB比较OK，开始做JWT处理
		loginResponse, err := tokenNext(c, user)
		if err != nil {
			httpresponse.FailWithMessage(err.Error(), c)
		} else {
			loginResponse.User = user
			httpresponse.OkWithDetailed(loginResponse, "登录成功", c)
		}
	}
	//} else {
	//	httpresponse.FailWithMessage("验证码错误", c)
	//}
}

// @Summary 短信登陆
// @Description 短信通知登陆，验证成功后，生成token
// @Tags Base
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.RLoginThird true "用户名, 密码, 验证码"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/login/sms [post]
func LoginSms(c *gin.Context) {
	var L request.LoginSMS
	c.ShouldBind(&L)

	err := global.V.MyService.SendSms.Verify(L.SmsRuleId, L.Mobile, L.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	user, err := global.V.MyService.User.LoginSms(L.Mobile)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	//DB比较OK，开始做JWT处理
	loginResponse, err := tokenNext(c, user)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		loginResponse.User = user
		httpresponse.OkWithDetailed(loginResponse, "登录成功", c)
	}

}

// @Summary 用户登陆三方
// @Description 3方平台登陆，验证成功后，生成token
// @Tags Base
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.RLoginThird true "用户名, 密码, 验证码"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/login/third [post]
func LoginThird(c *gin.Context) {
	var L request.RLoginThird
	c.ShouldBind(&L)
	//if err := util.Verify(L, util.LoginVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}

	//if !request.CheckPlatformExist(request.GetMyHeader(c).SourceType) {
	//	httpresponse.FailWithMessage("Header.SourceType unknow", c)
	//	return
	//}

	//if store.Verify(L.CaptchaId, L.Captcha, true) {
	//先从DB中做比对
	//U := &model.User{ThirdId: L.Code}
	user, newReg, err := global.V.MyService.User.LoginThird(L, request.GetMyHeader(c))
	if err != nil {
		httpresponse.FailWithMessage("用户名不存在或者密码错误 ,err:"+err.Error(), c)
	} else {
		//DB比较OK，开始做JWT处理
		loginResponse, err := tokenNext(c, user)
		loginResponse.IsNewReg = newReg
		if err != nil {
			httpresponse.FailWithMessage(err.Error(), c)
		} else {
			loginResponse.User = user
			httpresponse.OkWithDetailed(loginResponse, "登录成功", c)
		}
	}
	//} else {
	//	httpresponse.FailWithMessage("验证码错误", c)
	//}
}

// @Summary 项目列表
// @Description 每个项目的详细信息
// @Security ApiKeyAuth
// @Tags Base
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} model.Project
// @Router /base/project/list [post]
func ProjectList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	rs := global.V.ProjectMng.Pool

	httpresponse.OkWithDetailed(rs, "成功", c)
}

// @Summary 所有常量列表
// @Description 所有常量列表，方便调用与调试
// @Tags Base
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} httpresponse.Response
// @Router /base/const/list [get]
func ConstList(c *gin.Context) {
	var a model.Project
	c.ShouldBind(&a)

	list := make(map[string]interface{})

	list["PROJECT_TYPE_MAP"] = util.PROJECT_TYPE_MAP
	list["PlatformList"] = request.GetPlatformList()
	list["ThirdTypeList"] = model.GetUserThirdTypeList()
	list["UserRegTypeList"] = model.GetUserRegTypeList()
	list["UserRegTypeList"] = model.GetUserSexList()
	list["UserStatusList"] = model.GetUserStatusList()

	httpresponse.OkWithDetailed(list, "成功", c)

}
