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

type ParserTokenReq struct {
	Token string `json:"token" form:"token"`
}

// @Tags Base
// @Summary 解析一个TOKEN
// @Description 应用接到token后，要到后端再统计认证一下，确保准确
// @Param data body ParserTokenReq true "需要验证的token值"
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.CustomClaims
// @Router /base/parserToken [post]
func ParserToken(c *gin.Context) {
	//util.MyPrint("im in parserToken")
	var p ParserTokenReq
	c.ShouldBind(&p)

	//if err := util.Verify(sendSMS, util.SendSMSVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	//}
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

// @Tags Base
// @Summary header头结构体
// @Description 日常header里放一诸如验证类的东西，统一公示出来，方便使用
// @Param X-HeaderBaseInfo body request.HeaderBaseInfo true "客户端基础信息"
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} request.Header
// @Router /base/headerStruct [get]
func RequestHeaderStruct(c *gin.Context) {

}

// @Summary 发送验证码
// @Description 登陆、注册、通知等发送短信
// @Tags Base
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
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
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
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
	err, user := service.Login(U)
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

// @Summary 用户登陆三方
// @Description 用户登陆，验证，生成token
// @Tags Base
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.RLoginThird true "用户名, 密码, 验证码"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/loginThird [post]
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
	user, newReg, err := service.LoginThird(L, request.GetMyHeader(c))
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
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
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

// @Summary 所有常量列表
// @Description 常量列表
// @Tags Base
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce  application/json
// @Success 200 {object} httpresponse.Response
// @Router /base/constList [get]
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

var store = base64Captcha.DefaultMemStore

// @Tags Base
// @Summary 生成图片验证码
// @Description 防止有人恶意攻击，尝试破解密码
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {object} httpresponse.SysCaptchaResponse
// @Router /base/captcha [get]
func Captcha(c *gin.Context) {
	//字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.C.Captcha.ImgHeight, global.C.Captcha.ImgWidth, global.C.Captcha.NumberLength, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		global.V.Zap.Error("验证码获取失败!", zap.Any("err", err))
		httpresponse.FailWithMessage("验证码获取失败", c)
	} else {
		httpresponse.OkWithDetailed(httpresponse.SysCaptchaResponse{
			CaptchaId: id,
			PicPath:   b64s,
		}, "验证码获取成功", c)
	}
}
