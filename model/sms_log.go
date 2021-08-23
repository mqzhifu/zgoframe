package model

type SmsLog struct {
	MODEL
	AppId       		int		`json:"app_id" db:"define:tinyint(1);comment:app_id;defaultValue:0"`
	Title       		string	`json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`
	Content       		string	`json:"content" db:"define:varchar(255);comment:内容;defaultValue:''"`
	RuleId       		int		`json:"rule_id" db:"define:tinyint(1);comment:app_id;defaultValue:0"`
	Uid       			int		`json:"uid" db:"define:INT;comment:保留字段;defaultValue:0"`
	Type    			int 	`json:"type" db:"define:tinyint(1);comment:0验证码1通知2营销3国际;defaultValue:0"`
	Status    			int		`json:"status" db:"define:tinyint(1);comment:1成功2失败3发送中4等待发送;defaultValue:0"`
	Ip 					string	`json:"ip" db:"define:varchar(20);comment:ip地址;defaultValue:''"`
	Mobile 				string	`json:"mobile" db:"define:varchar(20);comment:手机号;defaultValue:''"`
	OutNo 				string	`json:"out_no" db:"define:varchar(50);comment:3方ID;defaultValue:''"`
	Channel 			int		`json:"channel" db:"define:tinyint(1);comment:1阿里2腾讯;defaultValue:0"`
	ThirdBackInfo 		string	`json:"third_back_info" db:"define:varchar(255);comment:请示3方返回结果集;defaultValue:''"`
	ThirdCallbackStatus int		`json:"third_callback_status" db:"define:tinyint(1);comment:3方状态;defaultValue:0"`
	ThirdCallbackInfo   string	`json:"third_callback_info" db:"define:varchar(255);comment:3方回执-信息;defaultValue:''"`
	ThirdCallbackTime   string 	`json:"third_callback_time" db:"define:varchar(255);comment:3方回执-时间;defaultValue:''"`
}

func(smsLog *SmsLog) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "短信发送日志"

	return m
}


