package model

type Instance struct {
	MODEL
	Name        string  	`json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`
	Host 		string 		`json:"host" form:"host" db:"define:varchar(255);comment:主机地址;defaultValue:''"`
	Port        string    	`json:"port" form:"port" db:"define:varchar(50);comment:主机端口号;defaultValue:''"`
	Env 		string    	`json:"env" form:"env" db:"define:varchar(100);comment:环境变量;defaultValue:''"`
	User 		string    	`json:"user" form:"user" db:"define:varchar(100);comment:环境变量;defaultValue:''"`
	Ps 		string    		`json:"ps" form:"ps" db:"define:varchar(100);comment:环境变量;defaultValue:''"`
}


func(instance *Instance) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "3方服务-实例"

	return m
}
