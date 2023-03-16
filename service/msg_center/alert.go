package msg_center

import (
	"go.uber.org/zap"
	"strconv"
	"zgoframe/util"
)

const (
	ALERT_CHANNEL_SMS      = 1
	ALERT_CHANNEL_EMAIL    = 2
	ALERT_CHANNEL_FEISHU   = 4
	ALERT_CHANNEL_DINGDING = 8
	ALERT_CHANNEL_WECHAT   = 16
)

type SendAlertSmsT struct {
	ProjectId string `json:"project_id"`
	Level     string `json:"level"`
	Content   string `json:"content"`
}

type AlertOption struct {
	SendMsgChannel    int
	MsgTemplateRuleId int
	SendSync          bool
	Log               *zap.Logger
	Sms               *Sms
	SmsReceiver       []string
	Email             *Email
	EmailReceiver     []string
	SendUid           int
}

type Alert struct {
	Op AlertOption
}

func NewAlert(alertOption AlertOption) (*Alert, error) {
	alert := new(Alert)
	alert.Op = alertOption
	return alert, nil
}

func (alert Alert) LogSend(ProjectId int, Content string, Level string) {

}

func (alert Alert) Send(ProjectId int, Content string, Level string) {
	//localIp, _ := util.GetLocalIp()

	replaceVarMap := make(map[string]string)
	replaceVarMap["project_id"] = strconv.Itoa(ProjectId)
	replaceVarMap["content"] = Content
	replaceVarMap["level"] = Level
	//if alert.Op.SendMsgChannel&ALERT_CHANNEL_SMS == ALERT_CHANNEL_SMS {
	//	for _, receiver := range alert.Op.SmsReceiver {
	//		SendSMS := request.SendSMS{
	//			Receiver:   receiver,
	//			SendIp:     localIp,
	//			RuleId:     alert.Op.MsgTemplateRuleId,
	//			ReplaceVar: replaceVarMap,
	//			SendUid:    alert.Op.SendUid,
	//		}
	//		recordNewId, err := alert.Op.Sms.Send(ProjectId, SendSMS)
	//		util.MyPrint(recordNewId, err)
	//	}
	//}

	//if alert.Op.SendMsgChannel&ALERT_CHANNEL_EMAIL == ALERT_CHANNEL_EMAIL {
	//	for _, receiver := range alert.Op.EmailReceiver {
	//		SendEmail := request.SendEmail{
	//			Receiver:   receiver,
	//			SendIp:     localIp,
	//			RuleId:     alert.Op.MsgTemplateRuleId,
	//			ReplaceVar: replaceVarMap,
	//			SendUid:    alert.Op.SendUid,
	//		}
	//		recordNewId, err := alert.Op.Email.Send(ProjectId, SendEmail)
	//		util.MyPrint(recordNewId, err)
	//	}
	//}

	if alert.Op.SendMsgChannel&ALERT_CHANNEL_FEISHU == ALERT_CHANNEL_FEISHU {
		url := "https://open.feishu.cn/open-apis/bot/v2/hook/db2f1cf2-cf33-41f0-bf16-abe0c39060b1"

		httpHeader := make(map[string]string)
		httpHeader["Content-Type"] = "application/json"
		httpCurl := util.NewHttpCurl(url, httpHeader)

		type FeishuMsgContentT struct {
			Text string `json:"text"`
		}

		type FeishuMsgT struct {
			MsgType string            `json:"msg_type"`
			Content FeishuMsgContentT `json:"content"`
		}

		feishuMsgContentT := FeishuMsgContentT{
			Text: "报警，程序出错了",
		}

		feishuMsgT := FeishuMsgT{
			MsgType: "text",
			Content: feishuMsgContentT,
		}

		//content := {"msg_type":"text","content":{"text":"request example"}}
		httpCurl.PostJson(feishuMsgT)

	}

}
