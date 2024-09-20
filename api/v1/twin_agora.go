package v1

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
		OssAccessKeyId:     global.C.AliOss.AccessKeyId,
		OssBucket:          global.C.AliOss.Bucket,
		OssAccessKeySecret: global.C.AliOss.AccessKeySecret,
		OssEndpoint:        global.C.AliOss.Endpoint,
	}
	agora := util.NewMyAgora(op)
	return agora
}

// @Tags TwinAgora
// @Summary 获取用户的录屏记录列表
// @Description 获取用户的录屏记录列表
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/list [POST]
func TwinAgoraCloudRecordList(c *gin.Context) {
	uid, _ := request.GetUid(c)
	var list []model.AgoraCloudRecord
	err := global.V.Base.Gorm.Where("uid = ?", uid).Find(&list).Error
	util.MyPrint("err:", err, " list:", list)
	httpresponse.OkWithAll(list, "成功", c)
}

// @Tags TwinAgora
// @Summary 检查云端录制环境
// @Description 当有未停止的云端录制时，自动stop掉
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @accept application/json
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/check [POST]
func TwinAgoraCloudRecordCheck(c *gin.Context) {
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

	var agoraCloudRecord = model.AgoraCloudRecord{
		Status:     model.AGORA_CLOUD_RECORD_STATUS_END,
		EndTime:    util.GetNowTimeSecondToInt(),
		StopAction: model.AGORA_CLOUD_RECORD_STOP_ACTION_REENTER,
	}
	global.V.Base.Gorm.Where("channel_name = ? and listener_agora_uid = ? and status = ?", formData.Cname, formData.Uid, model.AGORA_CLOUD_RECORD_STATUS_START).
		Order("start_time desc").
		Limit(1).
		Updates(&agoraCloudRecord)
	httpresponse.OkWithAll(gin.H{"success": true}, "成功", c)
}

// @Tags TwinAgora
// @Summary 申请/创建 录屏资源Id
// @Description 录屏时，要先从声网，申请一个资源ID，之后，才能开始（声网限制：每秒最多请求10次）
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body util.AgoraAcquireReq false "基础信息"
// @Produce application/json
// @Success 200 {object} util.AgoraCloudRecordRes "结果"
// @Router /twin/agora/cloud/record/create/acquire [POST]
func TwinAgoraCloudRecordCreateAcquire(c *gin.Context) {
	thisUid, _ := request.GetUid(c)
	// GetUtilAgora().ExecBGOssFile()
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
		Scene:               0, // 非延迟转换
		// Scene:               2,//延迟转换
	}

	agoraCloudRecordRes, err := GetUtilAgora().CreateAcquire(formData)
	if err != nil {
		httpresponse.FailWithAll(err.Error(), "失败", c)
		return
	}

	if agoraCloudRecordRes.Code > 0 || agoraCloudRecordRes.HttpCode != 200 {
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}

	acquireConfig, err := json.Marshal(formData.ClientRequest)
	if err != nil {
		util.MyPrint("CreateAcquire json.Marsha err:", err)
	}

	agoraCloudRecord := model.AgoraCloudRecord{
		Uid:              thisUid,
		ListenerAgoraUid: util.Atoi(formData.Uid),
		AcquireConfig:    string(acquireConfig),
		ChannelName:      formData.Cname,
		ResourceId:       agoraCloudRecordRes.ResourceId,
		Status:           model.AGORA_CLOUD_RECORD_STATUS_RESOURCE,
		ServerStatus:     model.AGORA_CLOUD_RECORD_SERVER_STATUS_UNDO,
	}
	err = global.V.Base.Gorm.Create(&agoraCloudRecord).Error
	if err != nil {
		httpresponse.FailWithAll("CreateAcquire gorm err:"+err.Error(), "失败", c)
		return
	}
	agoraCloudRecordRes.Id = agoraCloudRecord.Id
	httpresponse.OkWithAll(agoraCloudRecordRes, "CreateAcquire-成功", c)

}

// @Tags TwinAgora
// @Summary 开始录屏
// @Description 根据上一步获取到的ResourceId，开始录屏，其数据会推送到3方的OSS上
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body request.TwinAgoraReq false "基础信息"
// @Produce application/json
// @Success 200 {object} util.AgoraCloudRecordRes "结果"
// @Router /twin/agora/cloud/record/start [POST]
func TwinAgoraCloudRecordStart(c *gin.Context) {
	var userFormData request.TwinAgoraReq
	c.ShouldBind(&userFormData)

	record, err := GetCloudRecordById(userFormData.RecordId)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	if record.Uid <= 0 || record.ListenerAgoraUid <= 0 || record.ChannelName == "" || record.ResourceId == "" {
		httpresponse.FailWithMessage("db record : Uid <=0 || AgoraUid <= 0 || ChannelName is empty || ResourceId empty !", c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_RESOURCE {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_RESOURCE", c)
		return
	}

	// 获取RTC TOKEN
	var form request.TwinAgoraToken
	form.Username = strconv.Itoa(record.ListenerAgoraUid)
	form.Channel = record.ChannelName
	token, _ := GetRtcToken(form)

	var formData util.AgoraRecordStartReq
	// formData.ClientRequest = make(map[string]interface{})
	formData.ClientRequest.Token = token

	formData.Uid = strconv.Itoa(record.ListenerAgoraUid)
	formData.Cname = record.ChannelName
	formData.ResourceId = record.ResourceId
	agoraCloudRecordRes, agoraCloudRecordResBack, err := GetUtilAgora().CloudRecordSingleStreamDelayTranscoding(formData)
	if err != nil {
		httpresponse.FailWithAll(err.Error(), "失败", c)
		return
	}
	if agoraCloudRecordRes.Code > 0 || agoraCloudRecordRes.HttpCode != 200 {
		CloudRecordErr(userFormData.RecordId, agoraCloudRecordRes)
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}
	ClientRequestBytes, err := json.Marshal(agoraCloudRecordResBack.ClientRequest)
	if err != nil {
		util.MyPrint("CloudRecordStart json.Marshal err:", err)
	}
	var agoraCloudRecord = model.AgoraCloudRecord{
		SessionId:  agoraCloudRecordRes.Sid,
		StartTime:  util.GetNowTimeSecondToInt(),
		Status:     model.AGORA_CLOUD_RECORD_STATUS_START,
		ConfigInfo: string(ClientRequestBytes),
	}
	// agoraCloudRecord.Id = formData.RecordId
	err = global.V.Base.Gorm.Where(" id = ?", userFormData.RecordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		util.MyPrint("CloudRecordStart gorm updates err:", err)
	}

	httpresponse.OkWithAll(agoraCloudRecord, "CloudRecordStart-成功", c)
}

func CloudRecordErr(recordId int, agoraCloudRecordRes util.AgoraCloudRecordRes) {
	errBytes, _ := json.Marshal(agoraCloudRecordRes)
	var agoraCloudRecord = model.AgoraCloudRecord{
		ErrLog: string(errBytes),
	}
	err := global.V.Base.Gorm.Where(" id = ?", recordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		util.MyPrint("CloudRecordStart gorm updates err:", err)
	}
}

// @Tags TwinAgora
// @Summary 录屏查询
// @Description 根据上一步获取到的ResourceId，
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Param rid path string true "rid"
// @Success 200 {object} util.AgoraCloudRecordRes "结果"
// @Router /twin/agora/cloud/record/query/{rid} [GET]
func TwinAgoraCloudRecordQuery(c *gin.Context) {
	recordId := util.Atoi(c.Param("rid"))
	if recordId <= 0 {
		httpresponse.FailWithMessage("RecordId <= 0", c)
		return
	}

	record, err := GetCloudRecordById(recordId)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_START {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_START", c)
		return
	}

	agoraCloudRecordRes, err := GetUtilAgora().CloudRecordQuery(record.ResourceId, record.SessionId)
	if agoraCloudRecordRes.Code > 0 || agoraCloudRecordRes.HttpCode != 200 {
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}
	httpresponse.OkWithAll(agoraCloudRecordRes, "Query-成功", c)
}

// @Tags TwinAgora
// @Summary 停止录屏
// @Description 各种异常情况都最好调一下stop，不然OSS要一直花钱呐....~~~~~
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param rid path string true "resource_id"
// @Param type path string true "类型"
// @Produce application/json
// @Success 200 {object} util.AgoraCloudRecordRes "结果"
// @Router /twin/agora/cloud/record/stop/{rid}/{type} [GET]
func TwinAgoraCloudRecordStop(c *gin.Context) {
	recordId := util.Atoi(c.Param("rid"))
	record, err := GetCloudRecordById(recordId)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}

	if record.Status != model.AGORA_CLOUD_RECORD_STATUS_START {
		httpresponse.FailWithMessage("record db status != AGORA_CLOUD_RECORD_STATUS_START", c)
		return
	}

	agoraCloudRecordRes, err := GetUtilAgora().CloudRecordStop(strconv.Itoa(record.ListenerAgoraUid), record.ChannelName, record.ResourceId, record.SessionId)
	// if agoraCloudRecordRes.Code > 0 || agoraCloudRecordRes.HttpCode != 200 {
	if agoraCloudRecordRes.Code > 0 {
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}

	ServerResponseBytes, err := json.Marshal(agoraCloudRecordRes.ServerResponse)
	util.MyPrint("stop ServerResponseBytes:", string(ServerResponseBytes), " err:", err)
	var agoraCloudRecord = model.AgoraCloudRecord{
		EndTime:     util.GetNowTimeSecondToInt(),
		Status:      model.AGORA_CLOUD_RECORD_STATUS_END,
		StopResInfo: string(ServerResponseBytes),
		StopAction:  util.Atoi(c.Param("type")),
	}
	err = global.V.Base.Gorm.Where(" id = ?", recordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		util.MyPrint("gorm updates err:", err)
	}

	httpresponse.OkWithAll(agoraCloudRecordRes, "stop-成功", c)
}

// @Tags TwinAgora
// @Summary 获取RTM-token
// @Description 使用RTM前，动态获取token，然后再登陆声网，才可正常使用声网的功能(token时效是一天，如果存在且未失效正常返回，否则创建新的)
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
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
	// 从redis中获取缓存的token
	redisElement, err := global.V.Base.Redis.GetElementByIndex("rtm_token", form.Username)
	if err != nil {
		return token, errors.New("GetElementByIndex <rtm_token> err:" + err.Error())
	}
	util.MyPrint("rtm redisElement:", redisElement)

	redisTokenStr, err := global.V.Base.Redis.Get(redisElement)
	util.MyPrint("rtm Redis.Get :", redisTokenStr, err)
	if err != nil && err != redis.Nil {
		return token, errors.New("redis get err:" + err.Error())
	}
	if err != redis.Nil && redisTokenStr != "" {
		util.MyPrint("return old token")
		return redisTokenStr, nil
	}

	util.MyPrint("create new token.")

	// appID := global.C.Agora.AppId
	// appCertificate := global.C.Agora.AppCertificate
	// expiredTs := uint32(util.GetNowTimeSecondToInt() + redisElement.Expire)
	// result, err := util.RTMBuildToken(appID, appCertificate, form.Username, util.RoleRtmUser, expiredTs)
	result, err := GetUtilAgora().GetRtmToken(form.Username, 0)
	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
	// if err != nil {
	//	return token, errors.New("BuildToken err:" + err.Error())
	// }
	// util.MyPrint(result)
	// token := util.AccessToken{}
	// token.FromString(result)
	// if token.Message[util.KLoginRtm] != expiredTs {
	//	httpresponse.FailWithMessage("expiredTs:"+err.Error(),c)
	//	return
	// }

	_, err = global.V.Base.Redis.SetEX(redisElement, result, 0)
	if err != nil {
		return token, errors.New("redis set err:" + err.Error())
	}
	return result, nil
}

func GetRtcToken(form request.TwinAgoraToken) (token string, err error) {
	// 从redis中获取缓存的token
	redisElement, err := global.V.Base.Redis.GetElementByIndex("rtc_token", form.Username, form.Channel)
	if err != nil {
		return token, errors.New("GetElementByIndex <rtc_token> err:" + err.Error())
	}
	util.MyPrint("rtc redisElement:", redisElement)

	redisTokenStr, err := global.V.Base.Redis.Get(redisElement)
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
	_, err = global.V.Base.Redis.SetEX(redisElement, result, 0)
	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}
	return result, nil
}

// @Tags TwinAgora
// @Summary 获取RTC-token
// @Description  使用RTC前，动态获取token，然后再登陆声网，才可正常使用声网的功能(token时效是一天，如果存在且未失效正常返回，否则创建新的)
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
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
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param rid path string true "rid"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/oss/files/{rid} [GET]
func TwinAgoraCloudRecordOssFiles(c *gin.Context) {
	recordId := util.Atoi(c.Param("rid"))
	err := GenerateCloudVideo(recordId)
	if err != nil {
		httpresponse.FailWithMessage(err.Error(), c)
		return
	}
	httpresponse.OkWithAll(true, "成功", c)
}

func GetCloudRecordById(rid int) (record model.AgoraCloudRecord, err error) {
	if rid <= 0 {
		return record, errors.New("RecordId <= 0")
	}

	err = global.V.Base.Gorm.First(&record, rid).Error
	if err != nil {
		errInfo := "db not found recordId:" + strconv.Itoa(rid)
		return record, errors.New(errInfo)
	}
	return record, nil
}

// GenerateCloudVideo 生成云端录制文件，并重新上传oss
func GenerateCloudVideo(recordId int) (err error) {
	if recordId <= 0 {
		return errors.New("RecordId <= 0")
	}

	updateRes := global.V.Base.Gorm.Where("id = ? and status = ? and server_status = ? and video_url = ''",
		recordId, model.AGORA_CLOUD_RECORD_STATUS_END, model.AGORA_CLOUD_RECORD_SERVER_STATUS_UNDO).
		Updates(&model.AgoraCloudRecord{
			ServerStatus: model.AGORA_CLOUD_RECORD_SERVER_STATUS_ING,
		})
	if updateRes.RowsAffected == 0 {
		return errors.New("recordId illegal:" + strconv.Itoa(recordId))
	}

	var record model.AgoraCloudRecord
	if err := global.V.Base.Gorm.First(&record, recordId).Error; err != nil {
		return err
	}

	// 错误处理时，将server_status改为4处理异常
	defer func() {
		if err != nil {
			_ = global.V.Base.Gorm.Where("id = ?", recordId).Updates(&model.AgoraCloudRecord{
				ServerStatus: model.AGORA_CLOUD_RECORD_SERVER_STATUS_ERR,
			}).Error
		}
	}()

	clientRequestStart := util.ClientRequestStart{}
	err = json.Unmarshal([]byte(record.ConfigInfo), &clientRequestStart)
	if err != nil {
		return err
	}
	pathPrefix := ""
	for _, v := range clientRequestStart.StorageConfig.FileNamePrefix {
		pathPrefix += v + "/"
	}
	// pathPrefix := "agoraRecord/ckck/1660733248/"

	// fileManager := global.GetUploadObj(1, "")
	localDiskPath := global.V.Util.DocsManager.GetLocalDiskDownloadBasePath() + "/" + pathPrefix
	util.MyPrint("pathPrefix:", pathPrefix, " , localDiskPath:", localDiskPath)
	listObjectsResult, err := global.V.Util.DocsManager.Option.AliOss.OssLs(pathPrefix)
	if len(listObjectsResult.Objects) <= 0 {
		return errors.New("path:" + pathPrefix + " is  empty,no files.")
	}

	_, err = util.PathExists(localDiskPath)
	if err != nil {
		err = os.MkdirAll(localDiskPath, 0666)
		if err != nil {
			return errors.New("Mkdir:" + localDiskPath + " err:" + err.Error())
		}
	}
	type ProcessFileInfo struct {
		LocalDiskPath string
		OssPath       string
		FileName      string
		Uid           string
		ExtName       string
	}
	processFileInfoList := []ProcessFileInfo{}
	// av := []string{}
	// av := []string{"8b666674134ffc392685e183d4b4e11f_ckck__uid_s_110__uid_e_av.m3u8","8b666674134ffc392685e183d4b4e11f_ckck__uid_s_44446__uid_e_av.m3u8","8b666674134ffc392685e183d4b4e11f_ckck__uid_s_44446__uid_e_av.mpd"}
	for _, v := range listObjectsResult.Objects {
		filePathArr := strings.Split(v.Key, "/")
		fileName := filePathArr[len(filePathArr)-1]
		fileNameSplitArr := strings.Split(fileName, ".")
		fileExtName := fileNameSplitArr[1]
		if fileExtName == "mkv" {
			continue
		}
		fileNameArr := strings.Split(fileNameSplitArr[0], "_")
		sid := fileNameArr[0]
		cname := fileNameArr[1]
		uid := fileNameArr[5]

		util.MyPrint("oss , Size:", v.Size, ", key:", v.Key, " , fileNameArr:", fileNameArr)

		fileCategory := ""
		fileIndex := ""
		if fileExtName == "mp4" {
			fileIndex = fileNameArr[9]
		} else {
			fileCategory = fileNameArr[9]
		}

		if fileCategory == "av" {
			processFileInfo := ProcessFileInfo{
				LocalDiskPath: localDiskPath,
				OssPath:       v.Key,
				FileName:      fileName,
				Uid:           uid,
				ExtName:       fileExtName,
			}

			processFileInfoList = append(processFileInfoList, processFileInfo)
		}

		util.MyPrint("fileName:", fileName, " , fileExtName:", fileExtName, " , sid:", sid, " , cname:", cname, " , uid:", uid, " , fileCategory:", fileCategory, " fileIndex:", fileIndex)
		if err := global.V.Util.DocsManager.Option.AliOss.DownloadFile(v.Key, localDiskPath+fileName); err != nil {
			return err
		}
	}
	// util.MyPrint(processFileInfoList)
	lastFiles := make(map[string]string)
	pathPrefix = pathPrefix[:len(pathPrefix)-1] // 去掉最后的/
	for _, v := range processFileInfoList {
		newFileName := v.Uid + "_" + v.ExtName + ".mkv"
		newFileFullName := v.LocalDiskPath + newFileName
		command := "ffmpeg -i " + v.LocalDiskPath + v.FileName + " -c copy " + newFileFullName + " -y"
		util.MyPrint(command)
		ctx := exec.Command("bash", "-c", command)

		output, err := ctx.CombinedOutput()
		strOutput := string(output)
		if err != nil {
			return errors.New("ExecShellCommand : <" + command + "> ,  has error , output:" + strOutput + err.Error())
		}
		// 重新上传至oss
		if err = global.V.Util.DocsManager.Option.AliOss.UploadOneByFile(newFileFullName, pathPrefix, newFileName); err != nil {
			return err
		}
		lastFiles[newFileName] = global.V.Util.DocsManager.Option.AliOss.Op.LocalDomain + "/" + pathPrefix + "/" + newFileName
	}

	// 更新数据库
	lastFilesJson, err := json.Marshal(lastFiles)
	util.MyPrint("lastFilesJson:", string(lastFilesJson), " err:", err)
	var agoraCloudRecord = model.AgoraCloudRecord{
		ServerStatus: model.AGORA_CLOUD_RECORD_SERVER_STATUS_OK,
		VideoUrl:     string(lastFilesJson),
	}
	err = global.V.Base.Gorm.Where(" id = ? and video_url = ''", recordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		return err
	}
	return nil
}

// @Tags TwinAgora
// @Summary 呼叫功能的，配置信息
// @Description 主要是超时时间的配置，C端需要使用
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/config [GET]
func TwinAgoraConfig(c *gin.Context) {
	config := make(map[string]string)
	config["call_timeout"] = strconv.Itoa(apiServices().TwinAgora.CallTimeout)
	config["exec_timeout"] = strconv.Itoa(apiServices().TwinAgora.ExecTimeout)
	config["user_heartbeat_timeout"] = strconv.Itoa(apiServices().TwinAgora.UserHeartbeatTimeout)
	config["res_accept_timeout"] = strconv.Itoa(apiServices().TwinAgora.ResAcceptTimeout)
	config["entry_timeout"] = strconv.Itoa(apiServices().TwinAgora.EntryTimeout)

	httpresponse.OkWithAll(config, "Query-成功", c)
}

// @Tags TwinAgora
// @Summary web-socket使用时，一些状态，
// @Description 如：房间、用户连接状态等
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/socket/tools [GET]
func TwinAgoraSocketTools(c *gin.Context) {
	config := make(map[string]interface{})
	config["rtc_room_pool"] = apiServices().TwinAgora.RTCRoomPool
	config["rtc_user_pool"] = apiServices().TwinAgora.RTCUserPool

	httpresponse.OkWithAll(config, "Query-成功", c)
}

// @Tags TwinAgora
// @Summary 获取推送事件的统计信息
// @Description 如：发送标注图次数、发送图片次数、发送视频次数
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param start_time path string true "start_time"
// @Param end_time path string true "end_time"
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/statistics/event/all [GET]
func TwinAgoraStatisticsEventAll(c *gin.Context) {
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	startDate, _ := time.Parse("20060102", startTime)
	endDate, _ := time.Parse("20060102", endTime)
	dateArray := make([]string, 0)
	dateArray = append(dateArray, startDate.Format("20060102"))
	for currDate := startDate.AddDate(0, 0, 1); currDate.Before(endDate); currDate = currDate.AddDate(0, 0, 1) {
		dateArray = append(dateArray, currDate.Format("20060102"))
	}
	dateArray = append(dateArray, endDate.Format("20060102"))

	// select count(*), `date`, `type` from project_push_msg where `date`>="20230505" and `date`<= "20230523" group by `date`
	rows, err := global.V.Base.Gorm.
		Model(&model.ProjectPushMsg{}).
		Select("count(*) as count", "date", "type as eventType").
		Where("date >= ? and date <= ?", startTime, endTime).
		Group("date").
		Group("type").
		Rows()
	if err != nil {
		httpresponse.FailWithAll(err.Error(), "失败", c)
		return
	}

	defer rows.Close()
	tmp := make(map[int]map[string]int)
	for _, i := range []int{6, 1, 2, 3, 4, 5} {
		tmp[i] = make(map[string]int)
	}
	for rows.Next() {
		var count int
		var date string
		var eventType int
		rows.Scan(&count, &date, &eventType)

		if _, ok := tmp[eventType]; ok {
			tmp[eventType][date] = count
		}
	}

	ret := make(map[int][]int)

	for eventType, _ := range tmp {
		for _, date := range dateArray {
			if _, ok := tmp[eventType][date]; !ok {
				ret[eventType] = append(ret[eventType], 0)
			} else {
				ret[eventType] = append(ret[eventType], tmp[eventType][date])
			}
		}
	}

	httpresponse.OkWithAll(ret, "成功", c)
}
