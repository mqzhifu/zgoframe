package test

import (
	"encoding/json"
	"zgoframe/core/global"
	"zgoframe/service/msg_center"
	"zgoframe/util"
)

func Sms() {

	SendAlertSmsT := msg_center.SendAlertSmsT{
		ProjectId: "6",
		Level:     "debug",
		Content:   "支付时，出现订单号不匹配了",
	}
	SendAlertSmsTBytes, _ := json.Marshal(SendAlertSmsT)

	backInfo, err := global.V.Util.AliSms.Send("13522536459", "SMS_273495087", "正负无限科技", string(SendAlertSmsTBytes))
	util.MyPrint("err:", err)
	util.MyPrint("backInfo", backInfo)
}
