package util

//应用程序主动向3方发送报警
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

const (
	ALERT_METHOD_SYNC = 1
	ALERT_METHOD_ASYNC = 2
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
	SendMethod int
}
//level :报警级别，  SMS EMAIL FEISHU WEIXIN DINGDING ， -1：为全部
//template:报警发送的内容(模板)，还需要动态替换变量值
//template:报警发送的标题
//SendMethod:同步/异步
func NewAlertHook(level int ,template string,Title string,SendMethod int,log *zap.Logger)*AlertHook{
	log.Info("NewAlertHook:")

	alertHook := new (AlertHook)
	alertHook.Log = log
	if level == -1{
		levelList := GetAlertHookLevelList()
		allLevel := 0
		for _,v :=range levelList{
			allLevel += v
		}
		alertHook.Level = allLevel
	}else{
		alertHook.Level = level
	}

	alertHook.Template 		= template
	alertHook.Title 		= Title
	alertHook.SendMethod 	= SendMethod

	return alertHook
}


func (alertHook *AlertHook)Alert(msg string){
	alertHook.Log.Info("Alert one")

	if alertHook.Level & ALERT_LEVEL_SMS == ALERT_LEVEL_SMS{

	}

	if alertHook.Level & ALERT_LEVEL_EMAIL == ALERT_LEVEL_EMAIL{
		alertHook.Log.Info(" alert in email...")
		body := strings.Replace(alertHook.Template , "#body#",msg,-1)
		if alertHook.SendMethod == ALERT_METHOD_SYNC {
			alertHook.Email.SendOneEmailSync( "mqzhifu@sina.com" ,alertHook.Title,body)
		}
	}

	if alertHook.Level & ALERT_LEVEL_FEISHU == ALERT_LEVEL_FEISHU{


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