package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
)

//获取resourceId
type AgoraAcquireReq struct {
	Cname         string                `json:"cname"`         //频道名称
	Uid           string                `json:"uid"`           //声网uid, 最好用类似：99999，而不是用普通用户的UID
	ClientRequest AgoraAcquireClientReq `json:"clientRequest"` //详细的配置信息
}

//获取resourceId
type AgoraAcquireClientReq struct {
	Region              string `json:"region"`              //区域，这个写死：CN，即可
	ResourceExpiredHour int    `json:"resourceExpiredHour"` //单位为小时，录制周期，从成功开启云端录制并获得 sid （录制 ID）后开始计算。默认72
	Scene               int    `json:"scene"`               //0：（默认）实时音视频录制或延时混音,1：页面录制,2：延时转码
}

//开始录制
type AgoraRecordStartReq struct {
	Cname         string                 `json:"cname"`         //频道
	Uid           string                 `json:"uid"`           //uid
	ClientRequest map[string]interface{} `json:"clientRequest"` //
	ResourceId    string                 `json:"resource_id"`
	RecordId      int                    `json:"record_id"` //平台自己生成和自增ID，上一步获取中得到的
}

//所有请求声网的返回结果集合
type AgoraCloudRecordRes struct {
	Id             int                       `json:"id"`
	Code           int                       `json:"code"`
	Reason         string                    `json:"reason"`
	ResourceId     string                    `json:"resourceId"`
	Sid            string                    `json:"sid"`
	Message        string                    `json:"message"`        //这里有值，证明头部验证失败了
	ServerResponse AgoraRecordServerResponse `json:"serverResponse"` //给query接口使用
}

//给query接口使用
type AgoraRecordServerResponse struct {
	FileListMode         string                           `json:"fileListMode"`
	UploadingStatus      string                           `json:"uploadingStatus"`
	FileList             []AgoraFileList                  `json:"fileList"`
	Command              string                           `json:"command"`
	SubscribeModeBitmask int                              `json:"subscribeModeBitmask"`
	Vid                  string                           `json:"vid"`
	Payload              AgoraRecordServerPayloadResponse `json:"payload"`
}

//给query接口使用
type AgoraFileList struct {
	FileName       string `json:"fileName"`
	TrackType      string `json:"trackType"`
	Uid            int    `json:"uid"`
	MixedAllUser   bool   `json:"mixedAllUser"`
	IsPlayable     bool   `json:"isPlayable"`
	SliceStartTime int64  `json:"sliceStartTime"`
}

type AgoraRecordServerPayloadResponse struct {
	Message string `json:"message"`
}

//=================end================

//停止录制
type AgoraRecordStopReq struct {
	Cname         string                 `json:"cname"`         //频道
	Uid           string                 `json:"uid"`           //uid
	ClientRequest map[string]interface{} `json:"clientRequest"` //这个不是下划线模式，主要是对端的agora就这么定义的
	ResourceId    string                 `json:"resource_id"`
	Sid           string
	RecordId      int `json:"record_id"` //平台自己生成和自增ID，上一步获取中得到的
}

//开始录制 - 配置信息
type AgoraRecordingConfig struct {
	MaxIdleTime int `json:"maxIdleTime"` //最长空闲频道时间，默认30，5秒 < time < 30天。如果频道内没有用户加入的时间超过此值，录制自动结束
	//default：默认模式。录制过程中音频转码，分别生成 M3U8 音频索引文件和视频索引文件。
	//standard：标准模式。(推荐使用该模式)录制过程中音频转码，分别生成 M3U8 音频索引文件、视频索引文件和合并的音视频索引文件。如果在 Web 端使用 VP8 编码，则生成一个合并的 MPD 音视频索引文件。
	//original：原始编码模式。适用于单流音频不转码录制。仅订阅音频时（streamTypes 为 0）时该参数生效，录制过程中音频不转码，生成 M3U8 音频索引文件。
	//分析看：original好像没啥用，主要是 default 和 standard，default多了一个延迟转码的功能，而standard虽然不能有延迟转码后的MP4文件，但是在收到文件后，直接用特殊播放器可直接观看
	StreamMode        string `json:"streamMode"`
	StreamTypes       int    `json:"streamTypes"`       //0：仅订阅音频,1：仅订阅视频,2：（默认）订阅音频和视频
	ChannelType       int    `json:"channelType"`       //0:（默认）通信场景,1场场景，注：普通用户选择RTC的场景必须与此一致
	SubscribeUidGroup int    `json:"subscribeUidGroup"` //(选填）Number 类型，预估的订阅人数峰值。在单流模式下，为必填参数。
	VideoStreamType   int    `json:"videoStreamType"`   //0原始视频流大小 1小流大小，默认为0

	SubscribeVideoUids []string `json:"subscribeVideoUids"` //默认是该频道中的所有人
	SubscribeAudioUids []string `json:"subscribeAudioUids"` //默认是该频道中的所有人

	//TranscodingConfig  TwinAgoraRecordingConfigTranscodingConfig `json:"transcodingConfig"`
}

//延迟转码
type AgoraRecordingConfigAppsCollection struct {
	CombinationPolicy string `json:"combinationPolicy"`
}

//视频转码的详细设置。仅适用于合流模式，单流模式下不能设置该参数
type TwinAgoraRecordingConfigTranscodingConfig struct {
	Height           int    `json:"height"`
	Width            int    `json:"width"`
	Bitrate          int    `json:"bitrate"`
	Fps              int    `json:"fps"`
	MixedVideoLayout int    `json:"mixedVideoLayout"`
	BackgroundColor  string `json:"backgroundColor"`
}

//===================end====================
//声网将录制好的视频推送给阿里OSS
type AgoraStorageConfig struct {
	AccessKey      string   `json:"accessKey"`
	Region         int      `json:"region"`
	Bucket         string   `json:"bucket"`
	SecretKey      string   `json:"secretKey"`
	Vendor         int      `json:"vendor"`
	FileNamePrefix []string `json:"fileNamePrefix"`
}

type MyAgoraOption struct {
	AppId              string
	AppCertificate     string
	TokenExpire        int
	Domain             string
	HttpSecret         string
	HttpKey            string
	OssAccessKeyId     string
	OssBucket          string
	OssAccessKeySecret string
	OssEndpoint        string
}

type MyAgora struct {
	Option MyAgoraOption
}

func NewMyAgora(option MyAgoraOption) *MyAgora {
	MyAgora := new(MyAgora)
	MyAgora.Option = option
	return MyAgora
}

func (myAgora *MyAgora) GetCommonHTTPHeader() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json;charset=utf-8"
	headers["Authorization"] = "Basic " + myAgora.GetHTTPBaseAuth(myAgora.Option.HttpKey, myAgora.Option.HttpSecret)
	return headers
}

func (myAgora *MyAgora) GetCommTokenExpire(expire int) (int, error) {
	if expire <= 0 {
		if myAgora.Option.TokenExpire <= 0 {
			return 0, errors.New("agora expire<=0 ")
		}
		expire = myAgora.Option.TokenExpire
	}
	return expire, nil
}

func (myAgora *MyAgora) FormatAgoraRes(res string) (agoraCloudRecordRes AgoraCloudRecordRes, err error) {
	err = json.Unmarshal([]byte(res), &agoraCloudRecordRes)
	if err != nil {
		return agoraCloudRecordRes, err
	}
	return agoraCloudRecordRes, nil
}
func (myAgora *MyAgora) CreateAcquire(agoraAcquireReq AgoraAcquireReq) (agoraCloudRecordRes AgoraCloudRecordRes, err error) {
	url := myAgora.Option.Domain + myAgora.Option.AppId + "/cloud_recording/acquire"
	httpCurl := NewHttpCurl(url, myAgora.GetCommonHTTPHeader())
	res, err := httpCurl.PostJson(agoraAcquireReq)
	if err != nil {
		return agoraCloudRecordRes, err
	}
	return myAgora.FormatAgoraRes(res)
}

func (myAgora *MyAgora) GetRtcToken(username string, channel string, expire int) (token string, err error) {
	expire, err = myAgora.GetCommTokenExpire(expire)
	if err != nil {
		return token, err
	}
	MyPrint("create new token.")

	appID := myAgora.Option.AppId
	appCertificate := myAgora.Option.AppCertificate
	expiredTs := uint32(GetNowTimeSecondToInt() + expire)
	result, err := RTCBuildTokenWithUserAccount(appID, appCertificate, channel, username, RoleRtmUser, expiredTs)

	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}

	return result, nil
}

func (myAgora *MyAgora) GetRtmToken(username string, expire int) (token string, err error) {
	expire, err = myAgora.GetCommTokenExpire(expire)
	if err != nil {
		return token, err
	}
	MyPrint("create new token.")

	appID := myAgora.Option.AppId
	appCertificate := myAgora.Option.AppCertificate
	expiredTs := uint32(GetNowTimeSecondToInt() + expire)
	result, err := RTMBuildToken(appID, appCertificate, username, RoleRtmUser, expiredTs)

	if err != nil {
		return token, errors.New("BuildToken err:" + err.Error())
	}

	return result, nil
}

// 基于 Golang 实现的 HTTP 基本认证示例，使用 RTC 的服务端 RESTful API
func (myAgora *MyAgora) GetHTTPBaseAuth(customerKey string, customerSecret string) string {
	// 客户 ID
	//customerKey :=
	// 客户密钥
	//customerSecret := "a8b9fd618edb4061a7d8abd8f734ccaf"

	// 拼接客户 ID 和客户密钥并使用 base64 进行编码
	plainCredentials := customerKey + ":" + customerSecret
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(plainCredentials))

	//MyPrint("-------------------base64Credentials:", base64Credentials)
	return base64Credentials
}

//单流录制
func (myAgora *MyAgora) CloudRecordSingleStreamDelayTranscoding(formData AgoraRecordStartReq) (agoraCloudRecordRes AgoraCloudRecordRes, err error) {

	TwinAgoraRecordingConfig := AgoraRecordingConfig{}
	TwinAgoraRecordingConfig.VideoStreamType = 0
	TwinAgoraRecordingConfig.MaxIdleTime = 300
	//TwinAgoraRecordingConfig.StreamMode = "default"//延迟转码
	TwinAgoraRecordingConfig.StreamMode = "standard" //非延迟转码
	TwinAgoraRecordingConfig.StreamTypes = 2
	TwinAgoraRecordingConfig.ChannelType = 0
	TwinAgoraRecordingConfig.SubscribeUidGroup = 5
	//TwinAgoraRecordingConfig.SubscribeVideoUids = []string{"44446", "33311"}
	//TwinAgoraRecordingConfig.SubscribeAudioUids = []string{"44446", "33311"}
	TwinAgoraRecordingConfig.SubscribeVideoUids = []string{"#allstream#"}
	TwinAgoraRecordingConfig.SubscribeAudioUids = []string{"#allstream#"}

	//TwinAgoraRecordingConfigTranscodingConfig := request.TwinAgoraRecordingConfigTranscodingConfig{
	//	Height:           640,
	//	Width:            360,
	//	Bitrate:          500,
	//	Fps:              15,
	//	MixedVideoLayout: 1,
	//	BackgroundColor:  "#FF0000",
	//}
	//TwinAgoraRecordingConfig.TranscodingConfig = TwinAgoraRecordingConfigTranscodingConfig

	//持久化配置
	formData.ClientRequest["storageConfig"] = myAgora.CloudRecordGetStorageConfig(formData.Cname)
	//录屏 - 配置
	formData.ClientRequest["recordingConfig"] = TwinAgoraRecordingConfig
	//如需使用延时转码，则将 combinationPolicy 字段设置为 postpone_transcoding。设置该场景后，录制服务会在录制后 24 小时内对录制文件进行转码生成 MP4 文件，并将 MP4 文件上传至你指定的第三方云存储（不支持七牛云）。
	//twinAgoraRecordingConfigAppsCollection := AgoraRecordingConfigAppsCollection{
	//	CombinationPolicy: "postpone_transcoding",
	//}
	//formData.ClientRequest["appsCollection"] = twinAgoraRecordingConfigAppsCollection

	//mode:
	//1. individual: 单流,分开录制频道内每个 UID 的音频流和视频流，每个 UID 均有其对应的音频文件和视频文件。
	//2. mix: 合流,（默认模式）频道内所有 UID 的音视频混合录制为一个音视频文件
	//3. web: web页面,
	url := myAgora.Option.Domain + myAgora.Option.AppId + "/cloud_recording/resourceid/" + formData.ResourceId + "/mode/individual/start"
	httpCurl := NewHttpCurl(url, myAgora.GetCommonHTTPHeader())
	res, err := httpCurl.PostJson(formData)
	if err != nil {
		return agoraCloudRecordRes, errors.New("httpCurl err:" + err.Error())
	}

	return myAgora.FormatAgoraRes(res)
}
func (myAgora *MyAgora) ExecBGOssFile() {
	ossUtilCommand := "/data/www/golang/ossutilmac64 --endpoint " + myAgora.Option.OssEndpoint + "  --access-key-id " + myAgora.Option.OssAccessKeyId + " --access-key-secret " + myAgora.Option.OssAccessKeySecret + "   ls oss://" + myAgora.Option.OssBucket + "/agoraRecord"
	c := exec.Command("bash", "-c", ossUtilCommand)

	output, err := c.CombinedOutput()
	strOutput := string(output)
	if err != nil {
		MyPrint("ExecShellCommand : <"+ossUtilCommand+"> ,  has error , output:", strOutput, err.Error())
		return
	}

	MyPrint("ExecShellCommand : <"+ossUtilCommand+"> ,  success , output:", strOutput)
}

//注：channelName 做OOS路径的时候，不允许有下划线
func (myAgora *MyAgora) CloudRecordGetStorageConfig(channelName string) AgoraStorageConfig {
	twinAgoraStorageConfig := AgoraStorageConfig{
		AccessKey:      myAgora.Option.OssAccessKeyId,
		Bucket:         myAgora.Option.OssBucket,
		SecretKey:      myAgora.Option.OssAccessKeySecret,
		Vendor:         2,
		Region:         3,
		FileNamePrefix: []string{"agoraRecord", channelName, strconv.Itoa(GetNowTimeSecondToInt())},
	}
	return twinAgoraStorageConfig
}
func (myAgora *MyAgora) CloudRecordQuery(ResourceId string, SessionId string) (agoraCloudRecordRes AgoraCloudRecordRes, err error) {
	url := myAgora.Option.Domain + myAgora.Option.AppId + "/cloud_recording/resourceid/" + ResourceId + "/sid/" + SessionId + "/mode/individual/query"
	httpCurl := NewHttpCurl(url, myAgora.GetCommonHTTPHeader())
	res, err := httpCurl.Get()
	if err != nil {
		return agoraCloudRecordRes, err
	}
	return myAgora.FormatAgoraRes(res)
}

func (myAgora *MyAgora) CloudRecordStop(uid string, channel string, ResourceId string, SessionId string) (agoraCloudRecordRes AgoraCloudRecordRes, err error) {
	formData := AgoraRecordStopReq{}
	formData.Uid = uid
	formData.Cname = channel
	formData.ClientRequest = make(map[string]interface{})
	//twinAgoraAcquireStruct := request.TwinAgoraAcquireStruct{}
	//twinAgoraAcquireStruct.ClientRequest = make(map[string]interface{})
	//twinAgoraAcquireStruct.Uid = formData.Uid
	//twinAgoraAcquireStruct.Cname = formData.Cname
	//twinAgoraAcquireStruct["clientRequest"] = false

	url := myAgora.Option.Domain + myAgora.Option.AppId + "/cloud_recording/resourceid/" + ResourceId + "/sid/" + SessionId + "/mode/individual/stop"
	httpCurl := NewHttpCurl(url, myAgora.GetCommonHTTPHeader())
	res, err := httpCurl.PostJson(formData)
	if err != nil {
		return agoraCloudRecordRes, err
	}
	return myAgora.FormatAgoraRes(res)
}
