package model

type App struct {
	MODEL
	Name        string  	`json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`
	Type 		int 		`json:"type" form:"type" db:"define:tinyint(1);comment:类型,SERVIC frontend backend app;defaultValue:0"`
	Desc 		string 		`json:"desc" form:"desc" db:"define:varchar(255);comment:描述;defaultValue:''"`
	Key        	string    	`json:"key" form:"key" db:"define:varchar(50);comment:key;defaultValue:''"`
	SecretKey 	string    	`json:"secret_key" form:"secret_key" db:"define:varchar(100);comment:密钥;defaultValue:''"`
	Status 		int 		`json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭;defaultValue:0"`
}


func(app *App) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "所有项目集合"

	return m
}
