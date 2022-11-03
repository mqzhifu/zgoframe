package model

//@description 项目详情
type GameMatchRule struct {
	MODEL
	Name                  string `json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`                                                                        //名称
	GameId                int    `json:"game_id" form:"game_id" db:"define:tinyint(1);comment:游戏关联ID;defaultValue:0"`                                                                //类型,1service 2frontend 3backend 4app 5 unity 6 cocos
	Status                int    `json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭;defaultValue:0"`                                                                //状态1正常2关闭
	MatchTimeout          int    `json:"match_timeout" form:"match_timeout" db:"define:tinyint(1);comment:匹配超时时间(秒);defaultValue:0"`                                                 //状态1正常2关闭
	SuccessTimeout        int    `json:"success_timeout" form:"success_timeout" db:"define:tinyint(1);comment:匹配成功后，对方未接收超时时间(秒);defaultValue:0"`                                    //状态1正常2关闭
	Type                  int    `json:"type" form:"team_type" db:"define:tinyint(1);comment:1.N(TEAM)VS N(TEAM)2.N人够了就行(吃鸡模式);defaultValue:0"`                                      //1.N(TEAM)VS N(TEAM)类型王者2.N人够了就行(吃鸡模式)
	TeamMaxPeople         int    `json:"team_max_people" form:"team_max_people" db:"define:tinyint(1);comment:队伍最大人数;defaultValue:0"`                                                //多人组队一起玩游戏，一个队伍最大人数
	ConditionPeople       int    `json:"condition_people" form:"condition_people" db:"define:int;comment:多少人，可开始一局游戏;defaultValue:0"`                                                //多少人，可开始一局游戏
	Formula               string `json:"formula" form:"formula" db:"define:varchar(100);comment:权限计算公式;defaultValue:''"`                                                             //权重公式
	WeightTeamAggregation string `json:"weight_team_aggregation" form:"weight_team_aggregation" db:"define:varchar(50);comment:每个小组的最终权重计算聚合方法 sum min max average;defaultValue:''"` //权重的计算是以：人，为单位，但报名是以组为单位，当计算好每个人的权重后，最终求整组的权重值
	WeightScoreMin        int    `json:"weight_score_min" form:"weight_score_min" db:"define:int;comment:权重最小值;defaultValue:0"`                                                      //权重值范围：最小值，范围： 1-100
	WeightScoreMax        int    `json:"weight_score_max" form:"weight_score_max" db:"define:int;comment:权重最大值;defaultValue:0"`                                                      //权重值范围：最大值，范围： 1-100
	WeightAutoAssign      bool   `json:"WeightAutoAssign" form:"WeightAutoAssign" db:"define:tinyint(1);comment:权重自动匹配;defaultValue:0"`                                              //当权重值范围内，没有任何玩家，是否接收，自动调度分配，这样能提高匹配成功率
}

func (gameMatchRule *GameMatchRule) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "游戏匹配-规则配置"

	return m
}
