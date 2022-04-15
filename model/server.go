package model

type Server struct {
	MODEL
	Name           string `json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`
	Platform       string `json:"platform" form:"platform" db:"define:int;comment:平台类型1自有2阿里3腾讯4华为;defaultValue:0"`
	OutIp          string `json:"out_ip" form:"out_ip" db:"define:varchar(15);comment:外网IP;defaultValue:''"`
	InnerIp        string `json:"inner_ip" form:"inner_ip" db:"define:varchar(15);comment:内网IP;defaultValue:''"`
	Env            int `json:"env" form:"env" db:"define:int;comment:环境变量,1本地2开发3测试4预发布5线上;defaultValue:0"`
	Status         int    `json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭;defaultValue:0"`
	Ext            string `json:"ext" form:"ext" db:"define:varchar(255);comment:自定义配置信息;defaultValue:''"`
	ChargeUserName string `json:"charge_user_name" form:"charge_user_name" db:"define:varchar(50);comment:负责人姓名;defaultValue:''"`
	StartTime      int    `json:"start_time" form:"start_time" db:"define:int;comment:开始时间;defaultValue:0"`
	EndTime        int    `json:"end_time" form:"end_time" db:"define:int;comment:结束时间;defaultValue:0"`
	Price          int    `json:"price" form:"price" db:"define:int;comment:价格;defaultValue:0"`
}

func (server *Server) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "服务器"

	return m
}
