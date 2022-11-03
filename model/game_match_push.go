package model

////推送给3方，支持重试
//type PushElement struct {
//	Id       int
//	Category int    //1:报名超时 2匹配成功 3成功结果超时
//	Payload  string //自定义的载体
//}

//@description 游戏匹配-推送消息
type GameMatchPush struct {
	MODEL
	RuleId   int    `json:"rule_id" form:"rule_id" db:"define:int;comment:rule_id;defaultValue:0"`                                         //rule_id
	SelfId   int    `json:"self_id" form:"self_id" db:"define:int;comment:用redis自生成的自增ID;defaultValue:0"`                                  //用redis自生成的自增ID
	ATime    int    `json:"type" form:"type" db:"define:int;comment:添加时间;defaultValue:0"`                                                  //添加时间
	LinkId   int    `json:"person" form:"person" db:"define:int;comment:小组总人数;defaultValue:0"`                                             //关联调方用的ID
	Status   int    `json:"match_times" form:"match_times" db:"define:tinyint(1);comment:状态：1未推送2推送失败，等待重试3推送成功4推送失败，不再重试;defaultValue:0"` //状态：1未推送2推送失败，等待重试3推送成功4推送失败，不再重试
	Times    int    `json:"sign_timeout" form:"sign_timeout" db:"define:int;comment:已推送次数;defaultValue:0"`                                 //已推送次数
	Category int    `json:"success_timeout" form:"success_timeout" db:"define:int;comment:1:报名超时 2匹配成功 3成功结果超时;defaultValue:0"`            //1:报名超时 2匹配成功 3成功结果超时
	Payload  string `json:"addition" form:"addition" db:"define:text;comment:自定义的载体"`                                                      //自定义的载体
}

func (GameMatchPush *GameMatchPush) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "游戏匹配-推送消息"

	return m
}
