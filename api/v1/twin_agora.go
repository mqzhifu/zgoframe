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
// @accept application/json
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {boolean} boolean "true:成功 false:否"
// @Router /twin/agora/cloud/record/list [POST]
func TwinAgoraCloudRecordList(c *gin.Context) {
	uid, _ := request.GetUid(c)
	var list []model.AgoraCloudRecord
	err := global.V.Gorm.Where("uid = ?", uid).Find(&list).Error
	util.MyPrint("err:", err, " list:", list)
	httpresponse.OkWithAll(list, "成功", c)
}

// @Tags TwinAgora
// @Summary 申请/创建 录屏资源Id
// @Description 录屏时，要先从声网，申请一个资源ID，之后，才能开始（声网限制：每秒最多请求10次）
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body util.AgoraAcquireReq false "基础信息"
// @Produce application/json
// @Success 200 {object} util.AgoraCloudRecordRes "结果"
// @Router /twin/agora/cloud/record/create/acquire [POST]
func TwinAgoraCloudRecordCreateAcquire(c *gin.Context) {
	thisUid, _ := request.GetUid(c)
	//GetUtilAgora().ExecBGOssFile()
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
	err = global.V.Gorm.Create(&agoraCloudRecord).Error
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
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
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
	//agoraCloudRecord.Id = formData.RecordId
	err = global.V.Gorm.Where(" id = ?", userFormData.RecordId).Updates(&agoraCloudRecord).Error
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
	err := global.V.Gorm.Where(" id = ?", recordId).Updates(&agoraCloudRecord).Error
	if err != nil {
		util.MyPrint("CloudRecordStart gorm updates err:", err)
	}
}

// @Tags TwinAgora
// @Summary 录屏查询
// @Description 根据上一步获取到的ResourceId，
// @accept application/json
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
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
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param rid path string true "rid"
// @Produce application/json
// @Success 200 {object} util.AgoraCloudRecordRes "结果"
// @Router /twin/agora/cloud/record/stop/{rid} [GET]
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
	if agoraCloudRecordRes.Code > 0 || agoraCloudRecordRes.HttpCode != 200 {
		httpresponse.FailWithAll(agoraCloudRecordRes, "失败", c)
		return
	}

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

	httpresponse.OkWithAll(agoraCloudRecordRes, "stop-成功", c)
}

// @Tags TwinAgora
// @Summary 获取RTM-token
// @Description 使用RTM前，动态获取token，然后再登陆声网，才可正常使用声网的功能(token时效是一天，如果存在且未失效正常返回，否则创建新的)
// @accept application/json
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" default(11)
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

	//fileManager := global.GetUploadObj(1, "")
	localDiskPath := global.V.DocsManager.GetLocalDiskDownloadBasePath() + "/" + pathPrefix
	util.MyPrint("pathPrefix:", pathPrefix, " , localDiskPath:", localDiskPath)
	listObjectsResult, err := global.V.DocsManager.Option.AliOss.OssLs(pathPrefix)
	if len(listObjectsResult.Objects) <= 0 {
		httpresponse.FailWithMessage("path:"+pathPrefix+" is  empty,no files.", c)
		return
	}

	_, err = util.PathExists(localDiskPath)
	if err != nil {
		err = os.MkdirAll(localDiskPath, 0666)
		if err != nil {
			util.MyPrint("Mkdir:", localDiskPath, " err:", err)
			return
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
	//av := []string{}
	//av := []string{"8b666674134ffc392685e183d4b4e11f_ckck__uid_s_110__uid_e_av.m3u8","8b666674134ffc392685e183d4b4e11f_ckck__uid_s_44446__uid_e_av.m3u8","8b666674134ffc392685e183d4b4e11f_ckck__uid_s_44446__uid_e_av.mpd"}
	for _, v := range listObjectsResult.Objects {
		filePathArr := strings.Split(v.Key, "/")
		fileName := filePathArr[len(filePathArr)-1]
		fileNameSplitArr := strings.Split(fileName, ".")
		fileExtName := fileNameSplitArr[1]
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
		global.V.DocsManager.Option.AliOss.DownloadFile(v.Key, localDiskPath+fileName)
	}
	//util.MyPrint(processFileInfoList)
	for _, v := range processFileInfoList {
		newFileName := v.LocalDiskPath + v.Uid + "_" + v.ExtName + ".mkv"
		command := "ffmpeg -i " + v.LocalDiskPath + v.FileName + " -c copy " + newFileName
		util.MyPrint(command)
		ctx := exec.Command("bash", "-c", command)

		output, err := ctx.CombinedOutput()
		strOutput := string(output)
		if err != nil {
			util.MyPrint("ExecShellCommand : <"+command+"> ,  has error , output:", strOutput, err.Error())
		} else {
			util.MyPrint("ExecShellCommand : <"+command+"> ,  success , output:", strOutput)
		}
	}

}
func GetCloudRecordById(rid int) (record model.AgoraCloudRecord, err error) {
	if rid <= 0 {
		return record, errors.New("RecordId <= 0")
	}

	err = global.V.Gorm.First(&record, rid).Error
	if err != nil {
		errInfo := "db not found recordId:" + strconv.Itoa(rid)
		return record, errors.New(errInfo)
	}
	return record, nil
}
