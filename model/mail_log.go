package model

type MailLog struct {
	MODEL
	ProjectId  int    `json:"project_id" db:"define:int;comment:项目ID;defaultValue:0"`      //项目ID
	Title      string `json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`    //标题
	Content    string `json:"content" db:"define:varchar(255);comment:内容;defaultValue:''"` //模板内容,可变量替换
	RuleId     int    `json:"rule_id" db:"define:tinyint(1);comment:规则ID;defaultValue:0"`  //规则ID
	Receiver   string `json:"receiver" db:"define:varchar(50);comment:接收者邮件地址;defaultValue:''"`
	ExpireTime int    `json:"expire_time" db:"define:int;comment:失效时间;defaultValue:0"` //失效时间
	AuthCode   string `json:"auth_code" db:"define:varchar(50);comment:验证码;defaultValue:''"`
	AuthStatus int    `json:"auth_status" db:"define:tinyint(1);comment:1未使用2已使用3已超时;defaultValue:0"`
	SendUid    int    `json:"send_uid" db:"define:int;comment:发送者UID，管理员是9999，未知8888;defaultValue:0"` //发送者UID，管理员是9999，未知8888
	SendIp     string `json:"send_ip" db:"define:varchar(50);comment:发送者的IP;defaultValue:''"`         //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值
	Status     int    `json:"status" db:"define:tinyint(1);comment:1成功2失败3发送中4等待发送;defaultValue:0"`

	ReceiverRead int `json:"receiver_read" db:"define:tinyint(1);comment:接收者已读;defaultValue:0"` //接收者已读
	ReceiverDel  int `json:"receiver_del" db:"define:tinyint(1);comment:接收者已删除;defaultValue:0"` //接收者已删除
	SendDel      int `json:"send_del" db:"define:tinyint(1);comment:发送者已删除;defaultValue:0"`     //发送者已删除
	MailGroupId  int `json:"mail_group_id" db:"define:int;comment:群发的ID;defaultValue:0"`        //群发的ID
}

func (ailLog *MailLog) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "站内信-日志"

	return m
}
