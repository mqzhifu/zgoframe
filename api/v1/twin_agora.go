package v1

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

func GetUtilAgora() *util.MyAgora {
	op := util.MyAgoraOption{
		AppId:              global.C.Agora.AppId,
		AppCertificate:     global.C.Agora.AppCertificate,
		TokenExpire:        24 * 60 * 60,
		Domain:             global.C.Agora.Domain,
		HttpKey:            global.C.Agora.HttpKey,
		HttpSecret:         global.C.Agora.HttpSecret,
		OssAccessKeyId:     global.C.Oss.AccessKeyId,
		OssBucket:          global.C.Oss.Bucket,
		OssAccessKeySecret: global.C.Oss.AccessKeySecret,
		OssEndpoint:        global.C.Oss.Endpoint,
	}
	agora := util.NewMyAgora(op)
	return agora
}

// @Tags TwinAgora
// @Summary 申请/创建 录屏资源Id
// @Description 录屏时，要先从声网，申请一个资源ID，之后，才能开始（声网限制：每秒最多请求10次）
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body util.AgoraAcquireReq false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/create/acquire [POST]
func TwinAgoraCloudRecordCreateAcquire(c *gin.Context) {
	//GetUtilAgora().ExecBGOssFile()
	//util.ExitPrint(33)
	//oss := global.GetUploadObj(111, "")
	//oss.OssLs()
	//util.ExitPrint(1111)

	var formData util.AgoraAcquireReq
	c.ShouldBind(&formData)

	if util.Atoi(formData.Uid) <= 0 {
		httpresponse.FailWithMessage("uid <= 0", c)
		return
	}

	if formData.Cname == "" {
		httpresponse.FailWithMessage("cname empty", c)
		return
	}
	formData.ClientRequest = util.AgoraAcquireClientReq{
		Region:              "CN",
		ResourceExpiredHour: 72,
		Scene:               0, //非延迟转换
		//Scene:               2,//延迟转换

	}

	agoraCloudRecordRes, err := GetUtilAgora().CreateAcquire(formData)
	if err != nil {
		httpresponse.FailWithAll(err.Error(), "失败", c)
		return
	}
	if agoraCloudRecordRes.Code > 0 {
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}

	acquireConfig, err := json.Marshal(formData.ClientRequest)
	agoraCloudRecord := model.AgoraCloudRecord{
		ListenerAgoraUid: util.Atoi(formData.Uid),
		AcquireConfig:    string(acquireConfig),
		ChannelName:      formData.Cname,
		ResourceId:       agoraCloudRecordRes.ResourceId,
		Status:           model.AGORA_CLOUD_RECORD_STATUS_RESOURCE,
	}
	err = global.V.Gorm.Create(&agoraCloudRecord).Error
	if err != nil {
		httpresponse.FailWithAll("gorm err:"+err.Error(), "失败", c)
		return
	}
	agoraCloudRecordRes.Id = agoraCloudRecord.Id

	httpresponse.OkWithAll(agoraCloudRecordRes, "RTC-acquire-成功", c)

}

// @Tags TwinAgora
// @Summary 开始录屏
// @Description 根据上一步获取到的ResourceId，开始录屏，其数据会推送到3方的OSS上
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraReq false "基础信息"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/start [POST]
func TwinAgoraCloudRecordStart(c *gin.Context) {
	var userFormData request.TwinAgoraReq
	c.ShouldBind(&userFormData)

	if userFormData.RecordId <= 0 {
		httpresponse.FailWithMessage("RecordId <= 0", c)
		return
	}
	var record model.AgoraCloudRecord
	err := global.V.Gorm.First(&record, userFormData.RecordId).Error
	if err != nil {
		errInfo := "db not found recordId:" + strconv.Itoa(userFormData.RecordId)
		httpresponse.FailWithMessage(errInfo, c)
		return
	}

	if record.ListenerAgoraUid <= 0 || record.ChannelName == "" || record.ResourceId == "" {
		httpresponse.FailWithMessage("db record : AgoraUid <= 0 || ChannelName is empty || ResourceId empty !", c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_RESOURCE {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_RESOURCE", c)
		return
	}

	//获取RTC TOKEN
	var form request.TwinAgoraToken
	form.Username = strconv.Itoa(record.ListenerAgoraUid)
	form.Channel = record.ChannelName
	token, _ := GetRtcToken(form)

	var formData util.AgoraRecordStartReq
	//formData.ClientRequest = make(map[string]interface{})
	formData.ClientRequest.Token = token

	formData.Uid = strconv.Itoa(record.ListenerAgoraUid)
	formData.Cname = record.ChannelName
	formData.ResourceId = record.ResourceId
	agoraCloudRecordRes, agoraCloudRecordResBack, err := GetUtilAgora().CloudRecordSingleStreamDelayTranscoding(formData)
	if err != nil {
		httpresponse.FailWithAll(err.Error(), "失败", c)
		return
	}
	if agoraCloudRecordRes.Code > 0 {
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}
	//storageConfig, err1 := agoraCloudRecordResBack.ClientRequest["storageConfig"].(util.AgoraStorageConfig)
	//recordingConfig, err2 := agoraCloudRecordResBack.ClientRequest["recordingConfig"].(util.AgoraRecordingConfig)
	//storageConfigBytes, err3 := json.Marshal(storageConfig)
	//recordingConfigBytes, err4 := json.Marshal(recordingConfig)
	//util.MyPrint("CloudRecordSingleStreamDelayTranscoding storageConfigStr:", string(storageConfigBytes), " recordingConfigStr", string(recordingConfigBytes), " errList:", err1, err2, err3, err4)
	//ClientRequestArray := []string{string(storageConfigBytes), string(recordingConfigBytes)}
	//ClientRequestBytes, err := json.Marshal(ClientRequestArray)
	ClientRequestBytes, err := json.Marshal(agoraCloudRecordResBack.ClientRequest)
	var agoraCloudRecord = model.AgoraCloudRecord{
		SessionId:  agoraCloudRecordRes.Sid,
		StartTime:  util.GetNowTimeSecondToInt(),
		Status:     model.AGORA_CLOUD_RECORD_STATUS_START,
		ConfigInfo: string(ClientRequestBytes),
	}
	//agoraCloudRecord.Id = formData.RecordId
	err = global.V.Gorm.Where(" id = ?", userFormData.RecordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		util.MyPrint("gorm updates err:", err)
	}
	//util.MyPrint("agoraRecord:", agoraRecord)

	httpresponse.OkWithAll(agoraCloudRecord, "RTC-acquire-成功", c)
}

// @Tags TwinAgora
// @Summary 录屏查询
// @Description 根据上一步获取到的ResourceId，
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param rid path string true "rid"
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/query/{rid} [GET]
func TwinAgoraCloudRecordQuery(c *gin.Context) {
	recordId := util.Atoi(c.Param("rid"))
	if recordId <= 0 {
		httpresponse.FailWithMessage("RecordId <= 0", c)
		return
	}

	var record model.AgoraCloudRecord
	err := global.V.Gorm.First(&record, recordId).Error
	if err != nil {
		errInfo := "db not found recordId:" + strconv.Itoa(recordId)
		httpresponse.FailWithMessage(errInfo, c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_START {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_START", c)
		return
	}

	agoraCloudRecordRes, err := GetUtilAgora().CloudRecordQuery(record.ResourceId, record.SessionId)

	httpresponse.OkWithAll(agoraCloudRecordRes, "RTC-query-成功", c)
}

// @Tags TwinAgora
// @Summary 停止录屏
// @Description 各种异常情况都最好调一下stop，不然OSS要一直花钱呐....~~~~~
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param rid path string true "rid"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/stop/{rid} [GET]
func TwinAgoraCloudRecordStop(c *gin.Context) {
	recordId := util.Atoi(c.Param("rid"))
	if recordId <= 0 {
		httpresponse.FailWithMessage("RecordId <= 0", c)
		return
	}

	var record model.AgoraCloudRecord
	err := global.V.Gorm.First(&record, recordId).Error
	if err != nil {
		errInfo := "db not found recordId:" + strconv.Itoa(recordId)
		httpresponse.FailWithMessage(errInfo, c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_START {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_START", c)
		return
	}

	agoraCloudRecordRes, err := GetUtilAgora().CloudRecordStop(strconv.Itoa(record.ListenerAgoraUid), record.ChannelName, record.ResourceId, record.SessionId)
	ServerResponseBytes, err := json.Marshal(agoraCloudRecordRes.ServerResponse)
	util.MyPrint("stop ServerResponseBytes:", string(ServerResponseBytes), " err:", err)
	var agoraCloudRecord = model.AgoraCloudRecord{
		EndTime:     util.GetNowTimeSecondToInt(),
		Status:      model.AGORA_CLOUD_RECORD_STATUS_END,
		StopResInfo: string(ServerResponseBytes),
	}
	err = global.V.Gorm.Where(" id = ?", recordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		util.MyPrint("gorm updates err:", err)
	}

	httpresponse.OkWithAll(agoraCloudRecordRes, "RTC-acquire-成功", c)
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

	//appID := global.C.Agora.AppId
	//appCertificate := global.C.Agora.AppCertificate
	//expiredTs := uint32(util.GetNowTimeSecondToInt() + redisElement.Expire)
	//result, err := util.RTMBuildToken(appID, appCertificate, form.Username, util.RoleRtmUser, expiredTs)
	result, err := GetUtilAgora().GetRtmToken(form.Username, 0)
	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
	//if err != nil {
	//	return token, errors.New("BuildToken err:" + err.Error())
	//}
	//util.MyPrint(result)
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
	result, err := GetUtilAgora().GetRtcToken(form.Username, form.Channel, 0)
	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
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

// @Tags TwinAgora
// @Summary 处理阿里云OSS上的录屏文件
// @Description 将小文件，合并成一个大文件
// @accept application/json
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param rid path string true "rid"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/oss/files/{rid} [GET]
func TwinAgoraCloudRecordOssFiles(c *gin.Context) {
	recordId := util.Atoi(c.Param("rid"))
	if recordId <= 0 {
		httpresponse.FailWithMessage("RecordId <= 0", c)
		return
	}

	var record model.AgoraCloudRecord
	err := global.V.Gorm.First(&record, recordId).Error
	if err != nil {
		errInfo := "db not found recordId:" + strconv.Itoa(recordId)
		httpresponse.FailWithMessage(errInfo, c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_END {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_END", c)
		return
	}
	clientRequestStart := util.ClientRequestStart{}
	err = json.Unmarshal([]byte(record.ConfigInfo), &clientRequestStart)
	if err != nil {
		util.MyPrint("record.ConfigInfo Unmarshal err:", err)
	}
	pathPrefix := ""
	for _, v := range clientRequestStart.StorageConfig.FileNamePrefix {
		pathPrefix += v + "/"
	}
	//pathPrefix := "agoraRecord/ckck/1660733248/"

	util.MyPrint(pathPrefix)
	upload := global.GetUploadObj(1, "")
	listObjectsResult, err := upload.OssLs(pathPrefix)
	if len(listObjectsResult.Objects) <= 0 {
		httpresponse.FailWithMessage("path:"+pathPrefix+" is  empty,no files.", c)
		return
	}

	for _, v := range listObjectsResult.Objects {
		util.MyPrint("Size:", v.Size, ", key:", v.Key)

	}

}
