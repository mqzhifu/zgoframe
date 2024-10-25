package v1

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/core/global"
	httpmiddleware "zgoframe/http/middleware"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/service/user_center"
	"zgoframe/util"
)

// 图片验证码使用，主要是图片的ID得保存在内存(store)中
var store = base64Captcha.DefaultMemStore

// @Tags Base
// @Summary 生成图片验证码
// @Description 图片格式：BASE64，防止有人恶意攻击，如：短信轰炸、暴力破解密码等,前端使用方法：<img src="data:image/jpg;base64,接口获取的内容"/>
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.Captcha false "基础信息"
// @Produce application/json
// @Success 200 {object} httpresponse.Captcha "图片信息"
// @Router /base/captcha [POST]
func Captcha(c *gin.Context) {
	var form request.Captcha
	c.ShouldBind(&form)

	util.MyPrint(c.Request.Host, c.Request.URL)

	imgWidth := global.C.Captcha.ImgWidth
	imgHeight := global.C.Captcha.ImgHeight
	if (form.Width > 0 && form.Width < 1000) && (form.Height > 0 && form.Height < 1000) {
		imgWidth = form.Width
		imgHeight = form.Height
	}

	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(imgHeight, imgWidth, global.C.Captcha.NumberLength, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		global.V.Base.Zap.Error("验证码获取失败!", zap.Any("err", err))
		httpresponse.FailWithMessage("验证码获取失败", c)
	} else {
		httpresponse.OkWithAll(httpresponse.Captcha{
			Id:            id,
			PicContent:    b64s,
			ContentLength: global.C.Captcha.NumberLength,
		}, "验证码获取成功", c)
	}
}

// @Tags Base
// @Summary 发送短信
// @Description 登陆、注册、找回密码等，短信的内容由ruleId匹配（后台录入）
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.SendSMS true "基础信息"
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /base/send/sms [post]
func SendSms(c *gin.Context) {
	var sendSMSForm request.SendSMS
	c.ShouldBind(&sendSMSForm)
	if sendSMSForm.SendIp == "" {
		sendSMSForm.SendIp = c.Request.RemoteAddr
	}

	// if err := api.Verify(sendSMSForm, api.ApiVerify); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	// }

	projectId, _ := request.GetProjectId(c)
	dbNewId, err := ApiServices().Sms.Send(projectId, sendSMSForm)
	if err != nil {
		httpresponse.FailWithMessage("失败了："+err.Error(), c)
	} else {
		httpresponse.OkWithMessage(strconv.Itoa(dbNewId), c)
	}
}

// @Tags Base
// @Summary 发送邮件
// @Description 登陆、注册、找回密码等使用，目前不支持附件功能，邮件的内容由ruleId匹配（后台录入）
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.SendEmail true "基础信息"
// @Produce  application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /base/send/email [post]
func SendEmail(c *gin.Context) {
	var sendEmailForm request.SendEmail
	c.ShouldBind(&sendEmailForm)
	if sendEmailForm.SendIp == "" {
		sendEmailForm.SendIp = c.Request.RemoteAddr
	}

	projectId, _ := request.GetProjectId(c)
	dbNewId, err := ApiServices().Email.Send(projectId, sendEmailForm)
	if err != nil {
		httpresponse.FailWithMessage("失败了："+err.Error(), c)
	} else {
		httpresponse.OkWithMessage(strconv.Itoa(dbNewId), c)
	}
}

// @Tags Base
// @Summary 重置密码 - 通过短信
// @Description 忘记密码后，可发送短信通知，重置密码
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID" Enums(1,2,3,4) default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.RestPasswordSms true "用户名, 原密码, 新密码"
// @Success 200 {boolean} boolean "true:成功 false:否"
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

	err := ApiServices().Sms.Verify(form.SmsRuleId, form.Mobile, form.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	uid, _ := request.GetUid(c)
	err = ApiServices().User.ChangePassword(uid, form.NewPassword)
	if err != nil {
		httpresponse.FailWithMessage("修改失败:"+err.Error(), c)
	} else {
		httpresponse.OkWithMessage("修改成功", c)
	}
}

// @Tags Base
// @Summary 用户注册账号
// @Description 普通注册，需要有：用户名 密码
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Base-Info header string  true "客户端基础信息(json格式,参考request.HeaderBaseInfo)" default("{"app_version": "1.12.2","device": "iphone","device_id": "aaaabbbcccddd","device_version": "11.0.1","dpi": "1028x720","ip": "127.0.0.1","lat": "23.1123334455","lon": "45.11233311","os": 1,"os_version": "10.1","referer": "www.baidu.com"}")
// @Param data body request.Register true "基础信息"
// @Success 200 {object} model.User "用户结构体"
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

	header, _ := request.GetMyHeader(c)
	err, userInfo := ApiServices().User.RegisterByUsername(R, header)
	if err != nil {
		// global.V.Base.Zap.Error("注册失败", zap.Any("err", err))
		httpresponse.FailWithAll(userInfo, "注册失败:"+err.Error(), c)
	} else {
		httpresponse.OkWithAll(userInfo, "注册成功", c)
	}
}

// @Tags Base
// @Summary 用户注册账号-通过手机验证码
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param X-Base-Info header string  true "客户端基础信息(json格式,参考request.HeaderBaseInfo)" default("{"app_version": "1.12.2","device": "iphone","device_id": "aaaabbbcccddd","device_version": "11.0.1","dpi": "1028x720","ip": "127.0.0.1","lat": "23.1123334455","lon": "45.11233311","os": 1,"os_version": "10.1","referer": "www.baidu.com"}")
// @Param data body request.RegisterSms true "基础信息"
// @Success 200 {object} model.User "用户结构体"
// @Router /base/register/sms [post]
func RegisterSms(c *gin.Context) {
	var registerSmsForm request.RegisterSms
	_ = c.ShouldBind(&registerSmsForm)

	user := model.User{
		Username: registerSmsForm.Mobile,
		Mobile:   registerSmsForm.Mobile,
		Guest:    model.USER_GUEST_FALSE,
		Test:     model.USER_TEST_FALSE,
	}

	err := ApiServices().Sms.Verify(registerSmsForm.SmsRuleId, registerSmsForm.Mobile, registerSmsForm.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	header, _ := request.GetMyHeader(c)
	err, userInfo := ApiServices().User.Register(user, header, user_center.UserRegInfo{})
	if err != nil {
		// global.V.Base.Zap.Error("注册失败", zap.Any("err", err))
		httpresponse.FailWithAll(userInfo, "注册失败:"+err.Error(), c)
		return
	}

	httpresponse.OkWithAll(userInfo, "注册成功", c)

}

// @Tags Base
// @Summary 检测手机号：是否已存在
// @Description 注册/找加密码/登陆 使用
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CheckMobileExist true "基础信息"
// @Success 200 {boolean} boolean "true:存在 false:不存在"
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

	_, empty, err := ApiServices().User.FindUserByMobile(form.Mobile)
	util.MyPrint("CheckMobileExist empty:", empty)
	if err != nil {
		httpresponse.FailWithMessage("服务器错误，请等待或重试", c)
	} else {
		if !empty {
			httpresponse.OkWithAll(true, "成功", c)
		} else {
			httpresponse.OkWithAll(false, "成功", c)
		}
	}

}

// @Tags Base
// @Summary 检测用户名：是否已存在
// @Description 登陆 使用
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CheckUsernameExist true "用户信息"
// @Success 200 {boolean} boolean "true:存在 false:不存在"
// @Router /base/check/username [post]
func CheckUsernameExist(c *gin.Context) {
	var form request.CheckUsernameExist
	_ = c.ShouldBind(&form)

	if form.Username == "" || form.Purpose <= 0 {
		httpresponse.FailWithMessage("form.Username == '' || form.Purpose <= 0", c)
		return
	}

	if !util.CheckNameRule(form.Username) {
		httpresponse.FailWithMessage("username 格式 错误 ", c)
		return
	}

	_, empty, err := ApiServices().User.FindUserByUsername(form.Username)
	util.MyPrint("CheckMobileExist empty:", empty)
	if err != nil {
		httpresponse.FailWithMessage("服务器错误，请等待或重试", c)
	} else {
		if !empty {
			httpresponse.OkWithAll(true, "成功", c)
		} else {
			httpresponse.OkWithAll(false, "成功", c)
		}
	}

}

// @Tags Base
// @Summary 检测邮件：是否已存在
// @Description 注册/找加密码 使用
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.CheckEmailExist true "基础信息"
// @Success 200 {boolean} boolean "true:存在 false:不存在"
// @Router /base/check/email [post]
func CheckEmailExist(c *gin.Context) {
	print("dddd")
	var form request.CheckEmailExist
	_ = c.ShouldBind(&form)

	if form.Email == "" || form.Purpose <= 0 {
		httpresponse.FailWithMessage("form.Email == '' || form.Purpose <= 0", c)
		return
	}

	if !util.CheckEmailRule(form.Email) {
		httpresponse.FailWithMessage("email 格式 错误 ", c)
		return
	}

	if model.CheckConstInList(model.GetConstListPurpose(), form.Purpose) {
		httpresponse.FailWithMessage("form.Purpose err ", c)
		return
	}

	_, empty, err := ApiServices().User.FindUserByEmail(form.Email)
	if err != nil {
		httpresponse.FailWithMessage("服务器错误，请等待或重试", c)
	} else {
		if !empty {
			httpresponse.OkWithAll(true, "成功", c)
		} else {
			httpresponse.OkWithAll(false, "成功", c)
		}
	}

}

// @Tags Base
// @Summary 解析一个TOKEN
// @Description 单点登陆，各应用收到的接口都会带有token，可以到用户中心(微服务)再认证/解析一下，确保安全
// @accept application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.ParserToken true "需要验证的token值"
// @Produce  application/json
// @Success 200 {object} request.CustomClaims
// @Router /base/parser/token [post]
func ParserToken(c *gin.Context) {
	// util.MyPrint("im in parserToken")
	var p request.ParserToken
	c.ShouldBind(&p)

	j := httpmiddleware.NewJWT()
	claims, err := j.ParseToken(p.Token)
	if err != nil {
		httpresponse.FailWithAll(claims, "解析失败:"+err.Error(), c)
		return
	} else {
		httpresponse.OkWithAll(claims, "解析成功", c)
	}
}

// @Tags Base
// @Summary 普通登陆
// @Description 用户名/手机/邮箱+密码->登陆，验证成功后，生成token
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.Login true "基础信息"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/login [post]
func Login(c *gin.Context) {
	var L request.Login
	c.ShouldBind(&L)

	util.MyPrint(L)

	failedCnt, checkLoginFailedCntErr := ApiServices().User.CheckLoginFailedLimit(c.ClientIP(), L.Username, global.C.Login.MaxFailedCnt, global.C.Login.FailedLimitTime)
	if checkLoginFailedCntErr != nil {
		// httpresponse.FailWithMessage(checkLoginFailedCntErr.Error(), c)
		// return
	}
	// 先从DB中做比对
	U := &model.User{Username: L.Username, Password: L.Password}
	err, user := ApiServices().User.Login(U)
	if err != nil {
		ApiServices().User.IncrLoginFailedLimit(c.ClientIP(), L.Username)
		errMsg := "用户名不存在或者密码错误"
		if global.C.Login.MaxFailedCnt > 0 {
			balance := global.C.Login.MaxFailedCnt - failedCnt
			errMsg += "，还剩" + strconv.Itoa(balance) + "次机会"
		}
		httpresponse.FailWithMessage(errMsg, c)
	} else {
		loginType := ApiServices().User.TurnRegByUsername(L.Username)
		// DB比较OK，开始做JWT处理
		loginResponse, err := tokenNext(c, user, loginType)
		if err != nil {
			httpresponse.FailWithAll(loginResponse, err.Error(), c)
		} else {
			loginResponse.User = user
			httpresponse.OkWithAll(loginResponse, "登录成功", c)
		}
	}

}

// @Tags Base
// @Summary 短信(验证码)登陆
// @Description 登陆(验证)成功后，生成token
// @accept application/json
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.LoginSMS true "基础信息"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/login/sms [post]
func LoginSms(c *gin.Context) {
	var L request.LoginSMS
	c.ShouldBind(&L)

	err := ApiServices().Sms.Verify(L.SmsRuleId, L.Mobile, L.SmsAuthCode)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	user, err := ApiServices().User.LoginSms(L.Mobile)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	// DB比较OK，开始做JWT处理
	loginResponse, err := tokenNext(c, user, model.USER_REG_TYPE_MOBILE)
	if err != nil {
		httpresponse.FailWithAll(loginResponse, err.Error()+",(短信已使用，请重新发送一条)", c)
	} else {
		loginResponse.User = user
		httpresponse.OkWithAll(loginResponse, "登录成功", c)
	}

}

// @Tags Base
// @Summary 用户使用3方账号联合登陆
// @Description 3方平台登陆，验证成功后，生成token
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.RLoginThird true "基础信息"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/login/third [post]
func LoginThird(c *gin.Context) {
	var L request.RLoginThird
	c.ShouldBind(&L)
	// if err := util.Verify(L, util.LoginVerify); err != nil {
	//	httpresponse.FailWithMessage(err.Error(), c)
	//	return
	// }

	// if !request.CheckPlatformExist(request.GetMyHeader(c).SourceType) {
	//	httpresponse.FailWithMessage("Header.SourceType unknow", c)
	//	return
	// }

	// if store.Verify(L.CaptchaId, L.Captcha, true) {
	// 先从DB中做比对
	// U := &model.User{ThirdId: L.Code}
	header, _ := request.GetMyHeader(c)
	user, newReg, err := ApiServices().User.LoginThird(L, header)
	if err != nil {
		httpresponse.FailWithMessage("用户名不存在或者密码错误 ,err:"+err.Error(), c)
	} else {
		// DB比较OK，开始做JWT处理
		loginResponse, err := tokenNext(c, user, L.PlatformType)
		loginResponse.IsNewReg = newReg
		if err != nil {
			httpresponse.FailWithAll(loginResponse, err.Error(), c)
		} else {
			loginResponse.User = user
			httpresponse.OkWithAll(loginResponse, "登录成功", c)
		}
	}
	// } else {
	//	httpresponse.FailWithMessage("验证码错误", c)
	// }
}

// @Tags Base
// @Summary 项目获取 ACCCESS-TOKEN
// @Description 项目没有用户名+密码，只有密钥，拿到 ACCCESS-TOKEN 后，就跟正常用户登陆成功一样，可访问大部分API接口.sign值规则：md5(SecretKey+Timestamp+Access)
// @Produce  application/json
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.AccessToken true "基础信息"
// @Success 200 {object} httpresponse.LoginResponse
// @Router /base/access/token [post]
func AccessToken(c *gin.Context) {
	var L request.AccessToken
	c.ShouldBind(&L)

	if L.Sign == "" || L.Timestamp <= 0 {

	}
	projectId := request.GetProjectIdByHeader(c)
	if projectId <= 0 {
		httpresponse.FailWithMessage("projectId <= 0 ", c)
		return
	}

	projectInfo, empty := global.V.Util.ProjectMng.GetById(projectId)
	if empty {
		httpresponse.FailWithMessage("项目不存在", c)
		return
	}
	signStr := projectInfo.SecretKey + strconv.Itoa(L.Timestamp) + projectInfo.Access

	m := md5.New()
	m.Write([]byte(signStr))
	signMd5 := hex.EncodeToString(m.Sum(nil))

	util.MyPrint("signStr:", signStr, " , signMd5:", signMd5)
	if signMd5 != L.Sign {
		httpresponse.FailWithMessage("签名错误", c)
		return
	}

	U := &model.User{Username: projectInfo.Name, Password: "123456"}
	err, user := ApiServices().User.Login(U)
	if err != nil {
		httpresponse.FailWithMessage("未找到该用户", c)
		return
	}
	loginType := ApiServices().User.TurnRegByUsername(projectInfo.Name)
	// DB比较OK，开始做JWT处理
	loginResponse, err := tokenNext(c, user, loginType)
	if err != nil {
		httpresponse.FailWithAll(loginResponse, err.Error(), c)
	} else {
		loginResponse.User = user
		httpresponse.OkWithAll(loginResponse, "登录成功", c)
	}
}
