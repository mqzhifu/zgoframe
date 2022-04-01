package model

//站内信
type MailGroup struct {
	MODEL
	RuleId     int    `json:"rule_id" db:"define:tinyint(1);comment:规则ID;defaultValue:0"`                                 //规则ID
	PeopleType int    `json:"people_type" db:"define:tinyint(1);comment:接收人群，1单发2群发3指定group4指定tag5指定UIDS;defaultValue:0"` //项目ID
	Title      string `json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`                                   //标题
	Content    string `json:"content" db:"define:varchar(255);comment:模板内容,可变量替换;defaultValue:''"`                        //模板内容,可变量替换
	Receiver   string `json:"receiver" db:"define:varchar(50);comment:接收者，groupId，tagId , all ;defaultValue:''"`          //标题
	SendUid    int    `json:"send_uid" db:"define:int;comment:发送者UID，管理员是9999，未知8888;defaultValue:0"`                     //发送者UID，管理员是9999，未知8888
	SendIp     string `json:"send_ip" db:"define:varchar(50);comment:发送者的IP;defaultValue:''"`                             //发送者IP，如为空系统默认取：请求方的IP,最好给真实的，一但被刷，会使用此值

}

func (mailGroup *MailGroup) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "站内信 - 群发记录"

	return m
}
