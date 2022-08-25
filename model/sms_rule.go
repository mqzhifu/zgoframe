package model

//@description 短信配置规则
type SmsRule struct {
	MODEL
	ProjectId         int    `json:"app_id" db:"define:tinyint(1);comment:项目ID;defaultValue:0"`                    //项目ID
	Title             string `json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`                     //标题
	Content           string `json:"content" db:"define:varchar(255);comment:模板内容,可变量替换;defaultValue:''"`          //模板内容,可变量替换
	Type              int    `json:"type" db:"define:tinyint(1);comment:分类,1验证码2通知3营销;defaultValue:0"`             //分类,1验证码2通知3营销
	DayTimes          int    `json:"day_times" db:"define:int;comment:一天最多发送次数;defaultValue:0"`                    //每天最多发送次数
	Period            int    `json:"period" db:"define:int;comment:周期时间-秒;defaultValue:0"`                         //周期
	PeriodTimes       int    `json:"period_times" db:"define:int;comment:周期时间内-发送次数;defaultValue:0" `              //周期内最多可发送次数
	ExpireTime        int    `json:"expire_time" db:"define:int;comment:验证码要有失效时间;defaultValue:0" `                //验证码的失效时间
	Memo              string `json:"memo" db:"define:varchar(255);comment:描述，主要是给3方审核用;defaultValue:''"`           //备注
	Purpose           int    `json:"purpose" db:"define:tinyint(1);comment:用途,参考代码常量;defaultValue:0"`              //用途,参考代码常量
	Channel           int    `json:"channel" db:"define:tinyint(1);comment:1阿里2腾讯;defaultValue:0"`                 //渠道，1阿里2腾讯
	ThirdBackInfo     string `json:"third_back_info" db:"define:varchar(255);comment:请示3方返回结果集;defaultValue:''"`   //请示3方返回结果集
	ThirdTemplateId   string `json:"third_template_id" db:"define:varchar(100);comment:3方模板ID;defaultValue:''"`    //3方模板ID
	ThirdStatus       int    `json:"third_status" db:"define:tinyint(1);comment:3方状态;defaultValue:0"`              //3方状态
	ThirdReason       string `json:"third_reason" db:"define:varchar(255);comment:3方模板审核失败，理由信息;defaultValue:''"`  //3方模板审核失败，理由信息
	ThirdCallbackInfo string `json:"third_callback_info" db:"define:varchar(255);comment:3方回执-信息;defaultValue:''"` //:3方回执-信息
	ThirdCallbackTime string `json:"third_callback_time" db:"define:varchar(255);comment:3方回执-时间;defaultValue:''"` //3方回执-时间
}

func (smsRule *SmsRule) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "短信发送规则配置"

	return m
}
