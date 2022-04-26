package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

// @Tags GameMatch
// @Summary 生成图片验证码
// @Description BASE64图片内容，防止有人恶意攻击，如：短信轰炸、暴力破解密码等,<img src="data:image/jpg;base64,内容"/>
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {object} httpresponse.SysCaptchaResponse
// @Router /base/captcha [get]
func GameMatchSign(c *gin.Context) {
	//var httpReqBusiness gamematch.HttpReqBusiness
	//c.ShouldBind(&sendSMSForm)


	////httpd.Log.Info(" routing in signHandler : ")
	//
	//errCode , httpReqBusiness:= httpd.businessCheckData(postJsonStr)
	//if errCode != 0{
	//	errs := myerr.NewErrorCode(errCode)
	//	errInfo := zlib.ErrInfo{}
	//	json.Unmarshal([]byte(errs.Error()),&errInfo)
	//	return errInfo.Code,errInfo.Msg
	//}
	//errs := httpd.Gamematch.CheckHttpSignData(httpReqBusiness)
	//if errs != nil{
	//	errInfo := zlib.ErrInfo{}
	//	json.Unmarshal([]byte(errs.Error()),&errInfo)
	//
	//	return errInfo.Code,errInfo.Msg
	//}
	//signRsData, errs := httpd.Gamematch.Sign(httpReqBusiness)
	//if errs != nil{
	//	errInfo := zlib.ErrInfo{}
	//	json.Unmarshal([]byte(errs.Error()),&errInfo)
	//
	//	return errInfo.Code,errInfo.Msg
	//}
	//return 200,signRsData
	//
	//
	////// 生成默认数字的driver
	////driver := base64Captcha.NewDriverDigit(global.C.Captcha.ImgHeight, global.C.Captcha.ImgWidth, global.C.Captcha.NumberLength, 0.7, 80)
	////cp := base64Captcha.NewCaptcha(driver, store)
	////if id, b64s, err := cp.Generate(); err != nil {
	////	global.V.Zap.Error("验证码获取失败!", zap.Any("err", err))
	////	httpresponse.FailWithMessage("验证码获取失败", c)
	////} else {
	////	httpresponse.OkWithDetailed(httpresponse.SysCaptchaResponse{
	////		Id:         id,
	////		PicContent: b64s,
	////	}, "验证码获取成功", c)
	////}
}

// @Tags GameMatch
// @Summary 生成图片验证码
// @Description BASE64图片内容，防止有人恶意攻击，如：短信轰炸、暴力破解密码等,<img src="data:image/jpg;base64,内容"/>
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {object} httpresponse.SysCaptchaResponse
// @Router /base/captcha [get]
func GameMatchSignCancel(c *gin.Context) {
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

//}else if uri == "/success/del"{//匹配成功记录，不想要了，删除一掉
//code,msg = httpd.successDelHandler(postJsonStr)
//}else if uri == "/config"{//
//code,msg = httpd.ConfigHandler(postJsonStr)
//}else if uri == "/rule/add" {//添加一条rule
////code,msg = httpd.ruleAddOne(postDataMap)
//}else if uri == "/tools/getErrorInfo" {//所有错误码列表
//code,msg = httpd.getErrorInfoHandler()
//}else if uri == "/tools/clearRuleByCode"{//清空一条rule的所有数组，用于测试
//code,msg = httpd.clearRuleByCodeHandler(postJsonStr)
//}else if uri == "/tools/getNormalMetrics"{//html api
//code,msg = httpd.normalMetrics()
//}else if uri == "/tools/getRedisMetrics"{//html api
//code,msg = httpd.redisMetrics()
//}else if uri == "/tools/RedisStoreDb"{//html api
//code,msg = httpd.RedisStoreDb()
//}else if uri == "/tools/getHttpReqBusiness"{//html api

//
////通用 业务型  请求 数据  检查
//func  businessCheckData(postJsonStr string )(errCode int,httpReqBusiness HttpReqBusiness){
//	httpd.Log.Info(" businessCheckData : ")
//	if postJsonStr == ""{
//		return 802,httpReqBusiness
//	}
//	var jsonUnmarshalErr error
//	jsonUnmarshalErr = json.Unmarshal([]byte(postJsonStr),&httpReqBusiness)
//	if jsonUnmarshalErr != nil{
//		httpd.Log.Error(jsonUnmarshalErr)
//		mylog.Error(jsonUnmarshalErr)
//		return 459,httpReqBusiness
//	}
//	if httpReqBusiness.MatchCode == ""{
//		return 450,httpReqBusiness
//	}
//	rule ,err := httpd.Gamematch.RuleConfig.getByCategory(httpReqBusiness.MatchCode)
//	if err !=nil{
//		return 806,httpReqBusiness
//	}
//	httpReqBusiness.RuleId = rule.Id
//	_,err  = httpd.checkHttpdState(httpReqBusiness.RuleId)
//	if err != nil{
//		return 804,httpReqBusiness
//	}
//
//	return 0 ,httpReqBusiness
//}
