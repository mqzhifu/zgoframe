package v1

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags TwinAgora
// @Summary 申请录屏资源Id
// @Description 录屏时，要先申请一个资源ID，才能开始（声网限制：每秒最多请求10次）
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraAcquireStruct false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/rtc/get/cloud/record/acquire [POST]
func TwinAgoraRTCGetCloudRecordAcquire(c *gin.Context) {
	//client := &http.Client{}
	var formData request.TwinAgoraAcquireStruct
	formData.ClientRequest = make(map[string]interface{})
	c.ShouldBind(&formData)

	//formData.Uid = "99999" //如果是申请rid，最好用类似：99999，不能用视频中的UID

	url := global.C.Agora.Domain + global.C.Agora.AppId + "/cloud_recording/acquire"
	httpCurl := util.NewHttpCurl(url, GetAgoraCommonHTTPHeader())
	res, _ := httpCurl.PostJson(formData)
	agoraRecord := httpresponse.AgoraRecord{}
	err := json.Unmarshal([]byte(res), &agoraRecord)
	if err != nil {
		util.MyPrint("json.Unmarshal err:", err)
	}
	util.MyPrint("agoraRecord:", agoraRecord)
	httpresponse.OkWithAll(agoraRecord, "RTC-acquire-成功", c)

}

// @Tags TwinAgora
// @Summary 开始录屏
// @Description 根据上一步获取到的ResourceId，开始录屏，其数据会推送到3方的OSS上
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraRecordStartStruct false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/rtc/cloud/record/start [POST]
func TwinAgoraRTCCloudRecordStart(c *gin.Context) {

	var formData request.TwinAgoraRecordStartStruct
	formData.ClientRequest = make(map[string]interface{})
	c.ShouldBind(&formData)

	var form request.TwinAgoraToken
	form.Username = formData.Uid
	form.Channel = "ckck"
	//util.MyPrint(form)
	token, _ := GetRtcToken(form)
	//util.ExitPrint(token)
	formData.Token = token
	storageConfig := request.TwinAgoraStorageConfig{
		AccessKey:      global.C.Oss.AccessKeyId,
		Region:         0,
		Bucket:         global.C.Oss.Bucket,
		SecretKey:      global.C.Oss.AccessKeySecret,
		Vendor:         2,
		FileNamePrefix: []string{"imagora"},
	}
	formData.ClientRequest["storageConfig"] = storageConfig

	url := global.C.Agora.Domain + global.C.Agora.AppId + "/cloud_recording/resourceid/" + formData.ResourceId + "/mode/individual/start"
	httpCurl := util.NewHttpCurl(url, GetAgoraCommonHTTPHeader())
	res, _ := httpCurl.PostJson(formData)
	agoraRecord := httpresponse.AgoraRecord{}
	err := json.Unmarshal([]byte(res), &agoraRecord)
	if err != nil {
		util.MyPrint("json.Unmarshal err:", err)
	}
	util.MyPrint("agoraRecord:", agoraRecord)

	httpresponse.OkWithAll(agoraRecord, "RTC-acquire-成功", c)
}

// @Tags TwinAgora
// @Summary 停止录屏
// @Description 各种异常情况都最好调一下stop，不然OSS要一直花钱呐....~~~~~
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraRecordStopStruct false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/rtc/cloud/record/stop [POST]
func TwinAgoraRTCCloudRecordStop(c *gin.Context) {
	var formData request.TwinAgoraRecordStopStruct
	formData.ClientRequest = make(map[string]interface{})
	c.ShouldBind(&formData)

	twinAgoraAcquireStruct := request.TwinAgoraAcquireStruct{}
	twinAgoraAcquireStruct.ClientRequest = make(map[string]interface{})
	twinAgoraAcquireStruct.Uid = formData.Uid
	twinAgoraAcquireStruct.Cname = formData.Cname
	//twinAgoraAcquireStruct["clientRequest"] = false

	url := global.C.Agora.Domain + global.C.Agora.AppId + "/cloud_recording/resourceid/" + formData.ResourceId + "/sid/" + formData.Sid + "/mode/individual/stop"
	httpCurl := util.NewHttpCurl(url, GetAgoraCommonHTTPHeader())
	res, _ := httpCurl.PostJson(twinAgoraAcquireStruct)
	resourceIdAgora := httpresponse.AgoraRecord{}
	err := json.Unmarshal([]byte(res), &resourceIdAgora)
	if err != nil {
		util.MyPrint("json.Unmarshal err:", err)
	}

	httpresponse.OkWithAll(resourceIdAgora, "RTC-acquire-成功", c)
}

// @Tags TwinAgora
// @Summary 获取RTM-token
// @Description 使用RTM前，动态获取token，然后再登陆声网，才可正常使用声网的功能(token时效是一天，如果存在且未失效正常返回，否则创建新的)
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraToken false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/rtm/get/token [POST]
func TwinAgoraRTMGetToken(c *gin.Context) {
	var form request.TwinAgoraToken
	c.ShouldBind(&form)
	util.MyPrint("form:", form)
	if form.Username == "" {
		httpresponse.FailWithMessage("username is empty", c)
		return
	}
	result, err := GetRtmToken(form)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(result, "RTM-"+form.Username+"-成功", c)
	}

}

func GetRtmToken(form request.TwinAgoraToken) (token string, err error) {
	//从redis中获取缓存的token
	redisElement, err := global.V.Redis.GetElementByIndex("rtm_token", form.Username)
	if err != nil {
		return token, errors.New("GetElementByIndex <rtm_token> err:" + err.Error())
	}
	util.MyPrint("rtm redisElement:", redisElement)

	redisTokenStr, err := global.V.Redis.Get(redisElement)
	util.MyPrint("rtm Redis.Get :", redisTokenStr, err)
	if err != nil && err != redis.Nil {
		return token, errors.New("redis get err:" + err.Error())
	}
	if err != redis.Nil && redisTokenStr != "" {
		util.MyPrint("return old token")
		return redisTokenStr, nil
	}

	util.MyPrint("create new token.")

	appID := global.C.Agora.AppId
	appCertificate := global.C.Agora.AppCertificate
	expiredTs := uint32(1446455471)
	result, err := util.RTMBuildToken(appID, appCertificate, form.Username, util.RoleRtmUser, expiredTs)

	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
	util.MyPrint(result)

	//token := util.AccessToken{}
	//token.FromString(result)
	//if token.Message[util.KLoginRtm] != expiredTs {
	//	httpresponse.FailWithMessage("expiredTs:"+err.Error(),c)
	//	return
	//}

	_, err = global.V.Redis.SetEX(redisElement, result, 0)
	if err != nil {
		return token, errors.New("redis set err:" + err.Error())
	}
	return result, nil
}
func GetRtcToken(form request.TwinAgoraToken) (token string, err error) {
	//从redis中获取缓存的token
	redisElement, err := global.V.Redis.GetElementByIndex("rtc_token", form.Username, form.Channel)
	if err != nil {
		return token, errors.New("GetElementByIndex <rtc_token> err:" + err.Error())
	}
	util.MyPrint("rtc redisElement:", redisElement)

	redisTokenStr, err := global.V.Redis.Get(redisElement)
	util.MyPrint("rtc Redis.Get :", redisTokenStr, err)
	if err != nil && err != redis.Nil {
		return token, errors.New("redis get err:" + err.Error())
	}
	if err != redis.Nil && redisTokenStr != "" {
		util.MyPrint("return old token")
		return redisTokenStr, nil
	}

	util.MyPrint("create new token.")

	appID := global.C.Agora.AppId
	appCertificate := global.C.Agora.AppCertificate
	expiredTs := uint32(util.GetNowTimeSecondToInt() + redisElement.Expire)
	result, err := util.RTCBuildTokenWithUserAccount(appID, appCertificate, form.Channel, form.Username, util.RoleRtmUser, expiredTs)

	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
	//token := util.AccessToken{}
	//token.FromString(result)
	//if token.Message[util.KJoinChannel] != expiredTs {
	//	errors.New("no kJoinChannel ts")
	//}
	//
	//if token.Message[util.KPublishVideoStream] != 0 {
	//	errors.New("should not have publish video stream privilege")
	//}

	_, err = global.V.Redis.SetEX(redisElement, result, 0)
	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
	return result, nil
}

// @Tags TwinAgora
// @Summary 获取RTC-token
// @Description  使用RTC前，动态获取token，然后再登陆声网，才可正常使用声网的功能(token时效是一天，如果存在且未失效正常返回，否则创建新的)
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraToken false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/rtc/get/token [POST]
func TwinAgoraRTCGetToken(c *gin.Context) {
	var form request.TwinAgoraToken
	c.ShouldBind(&form)
	util.MyPrint("form:", form)
	if form.Username == "" {
		httpresponse.FailWithMessage("username is empty", c)
		return
	}

	if form.Channel == "" {
		httpresponse.FailWithMessage("channel is empty", c)
		return
	}
	result, err := GetRtcToken(form)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
	} else {
		httpresponse.OkWithAll(result, "RTC-"+form.Username+"-成功", c)
	}

}

func GetAgoraCommonHTTPHeader() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json;charset=utf-8"

	headers["Authorization"] = "Basic " + util.GetHTTPBaseAuth(global.C.Agora.HttpKey, global.C.Agora.HttpSecret)
	return headers
}
