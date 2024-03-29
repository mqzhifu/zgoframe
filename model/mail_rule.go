package model

//@description 站内信 - 配置规则
type MailRule struct {
	MODEL
	ProjectId   int    `json:"app_id" db:"define:tinyint(1);comment:项目ID;defaultValue:0"`           //项目ID
	Title       string `json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`            //标题
	Content     string `json:"content" db:"define:varchar(255);comment:模板内容,可变量替换;defaultValue:''"` //模板内容,可变量替换
	Type        int    `json:"type" db:"define:tinyint(1);comment:分类,1验证码2通知3营销;defaultValue:0"`    //分类,1验证码2通知3营销
	DayTimes    int    `json:"day_times" db:"define:int;comment:一天最多发送次数;defaultValue:0"`           //每天最多发送次数
	Period      int    `json:"period" db:"define:int;comment:周期时间-秒;defaultValue:0"`                //周期
	PeriodTimes int    `json:"period_times" db:"define:int;comment:周期时间内-发送次数;defaultValue:0" `     //周期内最多可发送次数
	ExpireTime  int    `json:"expire_time" db:"define:int;comment:验证码要有失效时间;defaultValue:0" `       //验证码的失效时间
	Memo        string `json:"memo" db:"define:varchar(255);comment:描述，主要是给3方审核用;defaultValue:''"`  //备注

	PeopleType int `json:"people_type" db:"define:tinyint(1);comment:接收人群，1单发2群发3指定group4指定tag5指定UIDS;defaultValue:0"` //项目ID

}

func (mailRule *MailRule) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "站内信 - 发送规则配置"

	return m
}
