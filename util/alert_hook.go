package util

//
import (
	"go.uber.org/zap"
	"strings"
)

const (
	ALERT_LEVEL_ALL = -1

	ALERT_LEVEL_SMS = 1
	ALERT_LEVEL_EMAIL = 2
	ALERT_LEVEL_FEISHU = 4
	ALERT_LEVEL_WEIXIN = 8
	ALERT_LEVEL_DINGDING = 16
)


func GetAlertHookLevelList()map[string]int{
	mm := make(map[string]int)
	mm["ALERT_LEVEL_SMS"] = 1
	mm["ALERT_LEVEL_EMAIL"] = 2
	mm["ALERT_LEVEL_FEISHU"] = 4
	mm["ALERT_LEVEL_WEIXIN"] = 8
	mm["ALERT_LEVEL_DINGDING"] = 16

	return mm
}

type AlertHook struct {
	Level int
	Email *MyEmail
	EmailOption EmailOption
	Log *zap.Logger
	Template string
	Title 	string
}



func NewAlertHook(level int ,template string,Title string,log *zap.Logger)*AlertHook{
	alertHook := new (AlertHook)
	alertHook.Log = log
	ExitPrint(level)
	if level == -1{
		levelList := GetAlertHookLevelList()
		allLevel := 0
		for _,v :=range levelList{
			allLevel += v
		}
		ExitPrint(allLevel)
		alertHook.Level = allLevel
	}else{
		alertHook.Level = level
	}

	alertHook.Template = template
	alertHook.Title = Title

	log.Info("NewAlertHook")

	return alertHook
}


func (alertHook *AlertHook)Alert(msg string){
	if alertHook.Level & ALERT_LEVEL_SMS == ALERT_LEVEL_SMS{

	}

	if alertHook.Level & ALERT_LEVEL_EMAIL == ALERT_LEVEL_EMAIL{

	}

	if alertHook.Level & ALERT_LEVEL_FEISHU == ALERT_LEVEL_FEISHU{
		body := strings.Replace(alertHook.Template , "#body#",msg,-1)
		alertHook.Email.SendOneEmailSync( "mqzhifu@sina.com" ,alertHook.Title,body)
	}

	if alertHook.Level & ALERT_LEVEL_WEIXIN == ALERT_LEVEL_WEIXIN{

	}
}


func (alertHook *AlertHook)SendSMS(){

}

func (alertHook *AlertHook)SendEmail(){

}

func (alertHook *AlertHook)SendFeishu(){

}

func (alertHook *AlertHook)SendWeixin(){

}