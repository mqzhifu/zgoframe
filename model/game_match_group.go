package model

//@description 游戏匹配-小组信息
type GameMatchGroup struct {
	MODEL
	RuleId         int    `json:"rule_id" form:"rule_id" db:"define:int;comment:rule_id;defaultValue:0"`                                              //rule_id
	SelfId         int    `json:"self_id" form:"self_id" db:"define:int;comment:用redis自生成的自增ID;defaultValue:0"`                                       //用redis自生成的自增ID
	Type           int    `json:"type" form:"type" db:"define:tinyint(1);comment:报名跟报名成功会各创建一条group记录，1：报名，2匹配成功;defaultValue:0"`                     //报名跟报名成功会各创建一条group记录，1：报名，2匹配成功
	Person         int    `json:"person" form:"person" db:"define:tinyint(1);comment:小组总人数;defaultValue:0"`                                           //小组总人数
	Weight         string `json:"weight" form:"weight" db:"define:varchar(50);comment:小组权重;defaultValue:''"`                                          //小组权重
	MatchTimes     int    `json:"match_times" form:"match_times" db:"define:tinyint(1);comment:已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃，不过没用到;defaultValue:0"` //已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃，不过没用到
	SignTimeout    int    `json:"sign_timeout" form:"sign_timeout" db:"define:int;comment:多少秒后无人来取，即超时，更新用户状态，删除数据;defaultValue:0"`                   //多少秒后无人来取，即超时，更新用户状态，删除数据
	SuccessTimeout int    `json:"success_timeout" form:"success_timeout" db:"define:int;comment:匹配成功后，无人来取，超时;defaultValue:0"`                        //匹配成功后，无人来取，超时
	SignTime       int    `json:"sign_time" form:"sign_time" db:"define:int;comment:报名时间;defaultValue:0"`                                             //报名时间
	SuccessTime    int    `json:"success_time" form:"success_time" db:"define:int;comment:匹配成功时间;defaultValue:0"`                                     //匹配成功时间
	PlayerIds      string `json:"player_ids" form:"player_ids" db:"define:varchar(100);comment:用户列表;defaultValue:''"`                                 //用户列表
	Addition       string `json:"addition" form:"addition" db:"define:varchar(100);comment:请求方附加属性值;defaultValue:''"`                                 //请求方附加属性值
	TeamId         int    `json:"team_id" form:"team_id" db:"define:tinyint(1);comment:组队互相PK的时候，得成两个队伍;defaultValue:0"`                              //组队互相PK的时候，得成两个队伍
	OutGroupId     int    `json:"out_group_id" form:"out_group_id" db:"define:int;comment:报名时，客户端请求时，自带的一个ID;defaultValue:0"`                         //报名时，客户端请求时，自带的一个ID
}

func (GameMatchGroup *GameMatchGroup) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "游戏匹配-小组信息"

	return m
}
