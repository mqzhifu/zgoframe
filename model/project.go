package model

//@description 项目详情
type Project struct {
	MODEL
	Name string `json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`                                //名称
	Type int    `json:"type" form:"type" db:"define:tinyint(1);comment:类型,1service 2frontend 3backend 4app;defaultValue:0"` //类型,1service 2frontend 3backend 4app
	Desc string `json:"desc" form:"desc" db:"define:varchar(255);comment:描述信息;defaultValue:''"`                            //描述信息

	SecretKey string `json:"secret_key" form:"secret_key" db:"define:varchar(100);comment:密钥;defaultValue:''"`     		//密钥,用于一些需要加密的场景
	Status    int    `json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭;defaultValue:0"`          		//状态1正常2关闭
	Access    string `json:"access" form:"access" db:"define:varchar(255);comment:baseAuth 认证KEY;defaultValue:''"` 		//baseAuth 认证KEY
	Lang 	  string `json:"lang" form:"lang" db:"define:tinyint(1);comment:实现语言1php2go3java4js;defaultValue:0"` 		//实现语言:1php2go3java4js5c++6c7c#
	Git       string `json:"git" form:"git" db:"define:varchar(255);comment:git仓地址;defaultValue:''"`               		//GIT代码仓URL地址
}

func (project *Project) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "服务/项目"

	return m
}
