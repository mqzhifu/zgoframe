package model

type CICDPublish struct {
	MODEL
	ServiceId	int 		`json:"service_id" form:"service_id" db:"define:int;comment:服务ID;defaultValue:0"`
	ServerId	int 		`json:"server_id" form:"server_id" db:"define:int;comment:服务器ID;defaultValue:0"`
	Status 		int 		`json:"status" form:"status" db:"define:tinyint(1);comment:1发布中2发布失败3发布成功;defaultValue:0"`
	ServiceInfo string  	`json:"service_info" form:"service_info" db:"define:varchar(50);comment:名称;defaultValue:''"`
	ServerInfo  string  	`json:"server_info" form:"server_info" db:"define:varchar(50);comment:名称;defaultValue:''"`
}


func(cICDPublish *CICDPublish) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "cicd 发布记录"

	return m
}

