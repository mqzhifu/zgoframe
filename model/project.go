package model

//@description 项目详情
type Project struct {
	MODEL
	Name string `json:"name" form:"name" db:"define:varchar(50);comment:项目名称;defaultValue:''"`                         //项目名称
	Type int    `json:"type" form:"type" db:"define:tinyint(1);comment:类型,SERVIC frontend backend app;defaultValue:0"` //类型
	Desc string `json:"desc" form:"desc" db:"define:varchar(255);comment:描述;defaultValue:''"`                          //描述

	SecretKey string `json:"secret_key" form:"secret_key" db:"define:varchar(100);comment:密钥;defaultValue:''"` //密钥
	Status    int    `json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭;defaultValue:0"`      //状态
	Git       string `json:"git" form:"git" db:"define:string(255);comment:git仓地址;defaultValue:''"`            //GIT代码仓URL地址
	Access    string `json:"access" form:"access" db:"define:string(255);comment:gate访问权限;defaultValue:''"`    //简单baseAuth 认证KEY
}

func (project *Project) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "所有项目集合"

	return m
}
