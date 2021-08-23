package model

type SmsRule struct {
	MODEL
	AppId       		int		`json:"app_id" db:"define:tinyint(1);comment:app_id;defaultValue:0"`
	Title       		string	`json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`
	Content       		string	`json:"content" db:"define:varchar(255);comment:内容;defaultValue:''"`
	Type    			int 	`json:"type" db:"define:tinyint(1);comment:0验证码1通知2营销3国际;defaultValue:0"`
	DayTimes    		int		`json:"day_times" db:"define:int;comment:一天最多发送次数;defaultValue:0"`
	Period    			int		`json:"period" db:"define:int;comment:周期时间-秒;defaultValue:0"`
	PeriodTimes    		int		`json:"period_times" db:"define:int;comment:周期时间内-发送次数;defaultValue:0" `
	Memo 				string	`json:"memo" db:"define:varchar(255);comment:描述，主要是给3方审核用;defaultValue:''"`
	Channel 			int		`json:"channel" db:"define:tinyint(1);comment:1阿里2腾讯;defaultValue:0"`
	ThirdBackInfo 		string	`json:"third_back_info" db:"define:varchar(255);comment:请示3方返回结果集;defaultValue:''"`
	ThirdTemplateId 	string 	`json:"third_template_id" db:"define:varchar(100);comment:3方模板ID;defaultValue:''"`
	ThirdStatus 		int		`json:"third_status" db:"define:tinyint(1);comment:3方状态;defaultValue:0"`
	ThirdReason 		string	`json:"third_reason" db:"define:varchar(255);comment:3方模板审核失败，理由信息;defaultValue:''"`

	ThirdCallbackInfo   string	`json:"third_callback_info" db:"define:varchar(255);comment:3方回执-信息;defaultValue:''"`
	ThirdCallbackTime   string 	`json:"third_callback_time" db:"define:varchar(255);comment:3方回执-时间;defaultValue:''"`
}




func(smsRule *SmsRule) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "短信发送规则配置"

	return m
}

