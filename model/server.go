package model

type Server struct {
	MODEL
	Name           string `json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`
	Platform       string `json:"platform" form:"platform" db:"define:varchar(50);comment:平台类型1自有2阿里3腾讯4华为;defaultValue:''"`
	OutIp          string `json:"out_ip" form:"out_ip" db:"define:varchar(255);comment:外部IP;defaultValue:''"`
	InnerIp        string `json:"inner_ip" form:"inner_ip" db:"define:varchar(50);comment:内容IP;defaultValue:''"`
	Status         int    `json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭3异常;defaultValue:0"`
	Env            string `json:"env" form:"env" db:"define:varchar(50);comment:环境变量,local dev test pre online;defaultValue:0"` //环境变量,local dev test pre online;defaultValue:0
	Ext            string `json:"ext" form:"ext" db:"define:varchar(255);comment:自定义配置信息;defaultValue:''"`
	StartTime      int    `json:"start_time" form:"start_time" db:"define:int;comment:开始时间;defaultValue:0"`
	EndTime        int    `json:"end_time" form:"end_time" db:"define:int;comment:结束时间;defaultValue:0"`
	Price          int    `json:"price" form:"price" db:"define:int;comment:价格;defaultValue:0"`
	ChargeUserName string `json:"charge_user_name" form:"charge_user_name" db:"define:varchar(50);comment:负责人姓名;defaultValue:''"`
}

func (server *Server) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "服务器"

	return m
}
