package model

//@description 游戏匹配-小组信息
type GameSyncRoom struct {
	MODEL
	RuleId            int    `json:"rule_id" form:"rule_id" db:"define:int;comment:rule_id;defaultValue:0"`                               //rule_id
	AddTime           int    `json:"add_time" form:"add_time" db:"define:int;comment:添加时间;defaultValue:0"`                                //添加时间
	StartTime         int    `json:"start_time" form:"start_time" db:"define:int;comment:开始游戏时间;defaultValue:0"`                          //开始游戏时间
	EndTime           int    `json:"end_time" form:"end_time" db:"define:int;comment:游戏结束时间;defaultValue:0"`                              //游戏结束时间
	ReadyTimeout      int    `json:"ready_timeout" form:"ready_timeout" db:"define:tinyint(1);comment:准备超时时间;defaultValue:0"`             //准备超时时间
	Status            int    `json:"status" form:"status" db:"define:int;comment:状态;defaultValue:0"`                                      //状态
	SequenceNumber    int    `json:"sequence_number" form:"sequence_number" db:"define:int;comment:匹配成功后，无人来取，超时;defaultValue:0"`         //当前逻辑帧号
	RandSeek          int    `json:"rand_seek" form:"rand_seek" db:"define:int;comment:当前逻辑帧号;defaultValue:0"`                            //随机数种子
	WaitPlayerOffline int    `json:"wait_player_offline" form:"wait_player_offline" db:"define:int;comment:玩家掉线等待时间;defaultValue:0"`      //玩家掉线等待时间
	PlayerIds         string `json:"player_ids" form:"player_ids" db:"define:varchar(100);comment:玩家列表;defaultValue:''"`                  //玩家列表
	PlayersAckList    string `json:"players_ack_list" form:"players_ack_list" db:"define:varchar(100);comment:最后一帧的确认情况;defaultValue:''"` //最后一帧的确认情况
	EndTotal          string `json:"end_total" form:"end_total" db:"define:varchar(255);comment:结算信息;defaultValue:''"`                    //结算信息
	LogicFrameHistory string `json:"logic_frame_history" form:"logic_frame_history" db:"define:text;comment:玩家的历史所有记录"`                   //玩家的历史所有记录

}

func (GameSyncRoom *GameSyncRoom) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "游戏匹配-小组信息"

	return m
}
