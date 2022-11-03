package model

//@description 游戏匹配-成功
type GameMatchSuccess struct {
	MODEL
	RuleId     int    `json:"rule_id" form:"rule_id" db:"define:int;comment:rule_id;defaultValue:0"`                           //rule_id
	SelfId     int    `json:"self_id" form:"self_id" db:"define:int;comment:用redis自生成的自增ID;defaultValue:0"`                    //用redis自生成的自增ID
	ATime      int    `json:"a_time" form:"a_time" db:"define:int;comment:匹配成功的时间;defaultValue:0"`                             //匹配成功的时间
	Timeout    int    `json:"timeout" form:"timeout" db:"define:int;comment:多少秒后无人来取后超时，更新用户状态，删除数据;defaultValue:0"`           //多少秒后无人来取后超时，更新用户状态，删除数据
	Teams      string `json:"teams" form:"teams" db:"define:varchar(50);comment:该结果，有几个 队伍，因为每个队伍是一个集合，要用来索引;defaultValue:''"` //该结果，有几个 队伍，因为每个队伍是一个集合，要用来索引
	PlayerIds  string `json:"player_ids" form:"player_ids" db:"define:varchar(100);comment:玩家ID列表;defaultValue:''"`            //玩家ID列表
	GroupIds   string `json:"group_ids" form:"group_ids" db:"define:varchar(100);comment:小组ID列表;defaultValue:''"`              //小组ID列表
	PushSelfId int    `json:"push_self_id" form:"push_self_id" db:"define:int;comment:推送ID,用redis生成的自增ID;defaultValue:0"`      //推送ID,用redis生成的自增ID
}

func (GameMatchSuccess *GameMatchSuccess) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "游戏匹配-成功"

	return m
}
