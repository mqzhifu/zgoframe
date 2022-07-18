package model

type StatisticsLog struct {
	MODEL
	HeaderCommon  string `json:"header_common" db:"define:text;comment:http公共请求头信息"`                  //http公共请求头信息
	HeaderBase    string `json:"header_base" db:"define:text;comment:http请求头客户端基础信息;defaultValue:''"` //http请求头客户端基础信息
	ProjectId     int    `json:"project_id" db:"define:int;comment:项目ID;defaultValue:0"`              //项目ID
	Category      int    `json:"category" db:"define:int;comment:分类，暂未使用;defaultValue:0"`             //分类，暂未使用
	Action        string `json:"action" db:"define:varchar(255);comment:动作标识;defaultValue:''"`        //动作标识
	Uid           int    `json:"uid" db:"define:int;comment:用户ID;defaultValue:0"`                     //uid
	Msg           string `json:"msg" db:"define:varchar(255);comment:自定义消息体;defaultValue:''"`         //自定义消息体
	Sn            string `json:"sn" db:"define:varchar(100);comment:设备-序列号"`
	SystemVersion string `json:"system_version" db:"define:varchar(100);comment:设备-版本号"`
	RecordTime    int    `json:"record_time" db:"define:int;comment:记录时间"`
	PackageName   string `json:"package_name" db:"define:varchar(100);comment:包名"`
	AppVersion    string `json:"app_version" db:"define:varchar(100);comment:应用版本号"`
	AppName       string `json:"app_name" db:"define:varchar(100);comment:应用名"`
}

func (statisticsLog *StatisticsLog) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "接收前端推送的统计日志"

	return m
}
