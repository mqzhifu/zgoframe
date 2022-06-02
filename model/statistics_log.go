package model

type StatisticsLog struct {
	MODEL
	HeaderCommon   	string    	`json:"header_common" db:"define:text;comment:http公共请求头信息;defaultValue:0"`//http公共请求头信息
	HeaderBase   	string 		`json:"header_base" db:"define:varchar(255);comment:http请求头客户端基础信息;defaultValue:''"`//http请求头客户端基础信息
	ProjectId       int    		`json:"project_id" db:"define:int;comment:项目ID;defaultValue:0"`    //项目ID
	Category    	int    		`json:"category" db:"define:int;comment:分类，暂未使用;defaultValue:0"`           //分类，暂未使用
	Action     		string 		`json:"action" db:"define:varchar(255);comment:动作标识;defaultValue:''"` //动作标识
	Uid      		int    		`json:"uid" db:"define:int;comment:用户ID;defaultValue:0"`                //uid
	Msg        		string 		`json:"msg" db:"define:varchar(255);comment:自定义消息体;defaultValue:''"`  //自定义消息体
}

func (statisticsLog *StatisticsLog) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "接收前端推送的统计日志"

	return m
}
