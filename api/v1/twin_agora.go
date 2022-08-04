package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"net/http"
	"strings"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)

// @Tags TwinAgora
// @Summary 申请录屏资源Id
// @Description 录屏时，要先申请一个资源ID，才能开始
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraAcquireStruct false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/rtc/get/cloud/record/acquire [POST]
func TwinAgoraRTCGetCloudRecordAcquire(c *gin.Context) {
	client := &http.Client{}
	//var data = request.TwinAgoraAcquireStruct{
	//	Cname:         "ckck",
	//	Uid:           "999999",
	//	ClientRequest: make(map[string]string),
	//}
	var formData request.TwinAgoraAcquireStruct
	formData.ClientRequest = make(map[string]string)
	c.ShouldBind(&formData)

	dataStrByte, _ := json.Marshal(formData)
	base64Credentials := util.GetHTTPBaseAuth()
	util.MyPrint("base64Credentials:", base64Credentials, " data:", formData, " dataStr:", string(dataStrByte))
	url := "https://api.agora.io/v1/apps/8ff429463a234c7bae327d74941a5956/cloud_recording/acquire"
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataStrByte)))

	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Authorization", "Basic "+base64Credentials)

	response, err := client.Do(request)
	defer response.Body.Close()

	util.MyPrint("response:", response.Body, " err:", err)
	b, err := ioutil.ReadAll(response.Body)
	bodyStr := string(b)
	util.MyPrint("response read body:", bodyStr, " err:", err)

	type ResourceIdAgora struct {
		ResourceId string `json:"resourceId"`
	}

	resourceIdAgora := ResourceIdAgora{}
	json.Unmarshal(b, &resourceIdAgora)
	util.MyPrint("resourceIdAgora:", resourceIdAgora)

	httpresponse.OkWithAll(resourceIdAgora.ResourceId, "RTC-acquire-成功", c)

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

	//从redis中获取缓存的token
	redisElement, err := global.V.Redis.GetElementByIndex("rtm_token", form.Username)
	if err != nil {
		httpresponse.FailWithMessage("GetElementByIndex <rtm_token> err:"+err.Error(), c)
		return
	}
	util.MyPrint("rtm redisElement:", redisElement)

	redisTokenStr, err := global.V.Redis.Get(redisElement)
	util.MyPrint("rtm Redis.Get :", redisTokenStr, err)
	if err != nil && err != redis.Nil {
		httpresponse.FailWithMessage("redis get err:"+err.Error(), c)
		return
	}
	if err != redis.Nil && redisTokenStr != "" {
		util.MyPrint("return old token")
		httpresponse.OkWithAll(redisTokenStr, "RTM-"+form.Username+"-成功", c)
		return
	}

	util.MyPrint("create new token.")

	appID := global.C.Agora.AppId
	appCertificate := global.C.Agora.AppCertificate
	expiredTs := uint32(1446455471)
	result, err := util.RTMBuildToken(appID, appCertificate, form.Username, util.RoleRtmUser, expiredTs)

	if err != nil {
		httpresponse.FailWithMessage("BuildToken err:"+err.Error(), c)
		return
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
		httpresponse.FailWithMessage("redis set err:"+err.Error(), c)
		return
	}

	httpresponse.OkWithAll(result, "RTM-"+form.Username+"-成功", c)

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
	//从redis中获取缓存的token
	redisElement, err := global.V.Redis.GetElementByIndex("rtc_token", form.Username, form.Channel)
	if err != nil {
		httpresponse.FailWithMessage("GetElementByIndex <rtc_token> err:"+err.Error(), c)
		return
	}
	util.MyPrint("rtc redisElement:", redisElement)

	redisTokenStr, err := global.V.Redis.Get(redisElement)
	util.MyPrint("rtc Redis.Get :", redisTokenStr, err)
	if err != nil && err != redis.Nil {
		httpresponse.FailWithMessage("redis get err:"+err.Error(), c)
		return
	}
	if err != redis.Nil && redisTokenStr != "" {
		util.MyPrint("return old token")
		httpresponse.OkWithAll(redisTokenStr, "RTC-"+form.Username+"-成功", c)
		return
	}

	util.MyPrint("create new token.")

	appID := global.C.Agora.AppId
	appCertificate := global.C.Agora.AppCertificate
	expiredTs := uint32(util.GetNowTimeSecondToInt() + redisElement.Expire)
	result, err := util.RTCBuildTokenWithUserAccount(appID, appCertificate, form.Channel, form.Username, util.RoleRtmUser, expiredTs)

	if err != nil {
		httpresponse.FailWithMessage("BuildToken err:"+err.Error(), c)
		return
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
		httpresponse.FailWithMessage("redis set err:"+err.Error(), c)
		return
	}

	httpresponse.OkWithAll(result, "RTC-"+form.Username+"-成功", c)
}
