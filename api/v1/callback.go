package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/model"
	"zgoframe/util"
)

type AgoraCloudCallbackPayloadDetailReq struct {
	ErrorCode  int         `json:"errorCode"`
	ErrorLevel int         `json:"errorLevel"`
	ErrorMsg   string      `json:"errorMsg"`
	Module     int         `json:"module"`
	MsgName    string      `json:"msgName"`
	Stat       int         `json:"stat"`
	Status     int         `json:"status"`
	FileList   []AgoraFile `json:"fileList"`
	ExitStatus int         `json:"exitStatus"`
	LeaveCode  int         `json:"leaveCode"`
}
type AgoraFile struct {
	FileName       string `json:"fileName"`
	TrackType      string `json:"trackType"`
	Uid            string `json:"uid"`
	MixedAllUser   bool   `json:"mixedAllUser"`
	IsPlayable     bool   `json:"isPlayable"`
	SliceStartTime int64  `json:"sliceStartTime"`
}

type AgoraCloudCallbackPayloadReq struct {
	Cname        string                             `json:"cname"`
	Sendts       int64                              `json:"sendts"`
	Sequence     int                                `json:"sequence"`
	ServiceType  int                                `json:"serviceType"`
	Sid          string                             `json:"sid"`
	Uid          string                             `json:"uid"`
	ServiceScene string                             `json:"serviceScene"`
	Details      AgoraCloudCallbackPayloadDetailReq `json:"details"`
}

type AgoraCloudCallbackReq struct {
	NoticeId  string                       `json:"noticeId"`
	ProductId int                          `json:"productId"`
	EventType int                          `json:"eventType"`
	NotifyMs  int64                        `json:"notifyMs"`
	Payload   AgoraCloudCallbackPayloadReq `json:"payload"`
}

// ===================
type AgoraRtcCallbackPayloadReq struct {
	ChannelName string `json:"channelName"`
	Platform    int    `json:"platform"`
	Reason      int    `json:"reason"`
	Ts          int    `json:"ts"`
	Uid         int    `json:"uid"`
}

type AgoraRtcCallbackReq struct {
	NoticeId  string                     `json:"noticeId"`
	ProductId int                        `json:"productId"`
	EventType int                        `json:"eventType"`
	NotifyMs  int64                      `json:"notifyMs"`
	Payload   AgoraRtcCallbackPayloadReq `json:"payload"`
}

// @Tags Callback
// @Summary 声网 - 回调
// @Description 订阅什么事件就回调什么事件
// @Param Agora-Signature header string true "签名" default(26a4fa1ec3df450caad3d8a4b907efe5476124da)
// @Param Agora-Signature-V2 header string true "签名" default(60216b719ca4a21701fcea43373370671d1401e4a8e408e2a550aa1a041fbe1c)
// @Produce application/json
// @Param data body AgoraRtcCallbackReq true " "
// @Success 200 {string} string "成功"
// @Router /callback/agora/rtc [post]
func AgoraCallbackRTC(c *gin.Context) {
	//prefix := "AgoraCallbackRTC "
	//for k, v := range c.Request.Header {
	//	util.MyPrint(prefix, "header ", k, v)
	//}
	//util.MyPrint("=======================")
	//util.MyPrint(prefix, "url:", c.Request.URL)
	//bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	//util.MyPrint(prefix, "ReadAll body:", string(bodyBytes), " err:", err)

	var form AgoraRtcCallbackReq
	err := c.ShouldBind(&form)
	util.MyPrint("form:", form, " err:", err)

	NotifyMsStr := strconv.FormatInt(form.NotifyMs, 10)
	payloadBytes, _ := json.Marshal(form.Payload)
	util.MyPrint("NotifyMsStr:", NotifyMsStr, " payloadBytes:", string(payloadBytes))
	agoraCallbackRecord := model.AgoraCallbackRecord{
		EventType:   form.EventType,
		NoticeId:    form.NoticeId,
		ProductId:   form.ProductId,
		NotifyMs:    NotifyMsStr,
		ChannelName: form.Payload.ChannelName,
		Payload:     string(payloadBytes),
	}
	global.V.Base.Gorm.Create(&agoraCallbackRecord)

	httpresponse.OkWithAll("回调成功", "ok", c)
}

// @Tags Callback
// @Summary 声网 - 云端录制 - 回调
// @Description 订阅什么事件就回调什么事件
// @Param Agora-Signature header string true "签名" default(26a4fa1ec3df450caad3d8a4b907efe5476124da)
// @Param Agora-Signature-V2 header string true "签名" default(60216b719ca4a21701fcea43373370671d1401e4a8e408e2a550aa1a041fbe1c)
// @Produce application/json
// @Param data body AgoraCloudCallbackReq true " "
// @Success 200 {string} string "成功"
// @Router /callback/agora/cloud [post]
func AgoraCallbackCloud(c *gin.Context) {
	//录制需要注意的，eventType-id: 1 2 3 11 30 31 32 40 41 80 81 90 1001
	//prefix := "AgoraCallbackCloud "
	//for k, v := range c.Request.Header {
	//	util.MyPrint(prefix, "header ", k, v)
	//}
	//util.MyPrint("=======================")
	//util.MyPrint(prefix, "url:", c.Request.URL)
	//bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	//util.MyPrint(prefix, "ReadAll body:", string(bodyBytes), " err:", err)

	var form AgoraCloudCallbackReq
	err := c.ShouldBind(&form)
	util.MyPrint("form:", form, " err:", err)

	NotifyMsStr := strconv.FormatInt(form.NotifyMs, 10)
	payloadBytes, _ := json.Marshal(form.Payload)
	util.MyPrint("NotifyMsStr:", NotifyMsStr, " payloadBytes:", string(payloadBytes))
	agoraCallbackRecord := model.AgoraCallbackRecord{
		EventType:   form.EventType,
		NoticeId:    form.NoticeId,
		ProductId:   form.ProductId,
		NotifyMs:    NotifyMsStr,
		ChannelName: form.Payload.Cname,
		SessionId:   form.Payload.Sid,
		Payload:     string(payloadBytes),
	}
	global.V.Base.Gorm.Create(&agoraCallbackRecord)

	if form.EventType == model.CallbackEventAllUploaded {
		go func() {
			var record model.AgoraCloudRecord
			err := global.V.Base.Gorm.First(&record).Where("session_id = ?", form.Payload.Sid).Error
			if err != nil {
				return
			}
			_ = GenerateCloudVideo(record.Id)
		}()
	}

	if form.EventType == model.CallbackEventRecordExit {
		go func() {
			var agoraCloudRecord = model.AgoraCloudRecord{
				EndTime:    util.GetNowTimeSecondToInt(),
				Status:     model.AGORA_CLOUD_RECORD_STATUS_END,
				StopAction: model.AGORA_CLOUD_RECORD_STOP_ACTION_CALLBACK,
			}
			err = global.V.Base.Gorm.Where("status != ? and session_id = ?", model.AGORA_CLOUD_RECORD_STATUS_END, agoraCallbackRecord.SessionId).Updates(&agoraCloudRecord).Error
			if err != nil {
				util.MyPrint("gorm updates err:", err)
			}
		}()
	}

	httpresponse.OkWithAll("回调成功", "ok", c)
}
