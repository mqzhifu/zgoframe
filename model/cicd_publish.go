package model

type CicdPublish struct {
	MODEL
	ServiceId   int    `json:"service_id" form:"service_id" db:"define:int;comment:服务ID;defaultValue:0"`
	ServerId    int    `json:"server_id" form:"server_id" db:"define:int;comment:服务器ID;defaultValue:0"`
	Status      int    `json:"status" form:"status" db:"define:tinyint(1);comment:1发布中2发布失败3发布成功;defaultValue:0"`
	ServiceInfo string `json:"service_info" form:"service_info" db:"define:varchar(255);comment:服务信息-备份;defaultValue:''"`
	ServerInfo  string `json:"server_info" form:"server_info" db:"define:varchar(255);comment:服务器信息-备份;defaultValue:''"`
	Log         string `json:"log" form:"log" db:"define:text;comment:日志;defaultValue:''"`
	ErrInfo     string `json:"err_info" form:"err_info" db:"define:varchar(255);comment:错误日志;defaultValue:''"`
	ExecTime    int    `json:"exec_time" form:"exec_time" db:"define:int;comment:执行时间;defaultValue:0"`
}

func (cicdPublish *CicdPublish) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "cicd发布记录"

	return m
}
