package msg_center

import "go.uber.org/zap"

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
