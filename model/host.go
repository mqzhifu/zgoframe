package model

type Server struct {
	MODEL
	Name        string  	`json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`
	Platform    string  	`json:"platform" form:"platform" db:"define:varchar(50);comment:平台类型;defaultValue:''"`
	OutIp 		string 		`json:"out_ip" form:"out_ip" db:"define:varchar(255);comment:外部IP;defaultValue:''"`
	InnerIp     string    	`json:"inner_ip" form:"inner_ip" db:"define:varchar(50);comment:内容IP;defaultValue:''"`
	Status 		int 		`json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭3异常;defaultValue:0"`
	Env 		string 		`json:"env" form:"env" db:"define:varchar(50);comment:环境变量;defaultValue:0"`
}


func(server *Server) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "主机列表"

	return m
}
