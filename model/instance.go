package model

type Instance struct {
	MODEL
	Platform       string `json:"platform" form:"platform" db:"define:int;comment:平台类型1自有2阿里3腾讯4华为;defaultValue:0"`
	Name           string `json:"name" form:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`
	Host           string `json:"host" form:"host" db:"define:varchar(255);comment:主机地址;defaultValue:''"`
	Port           string `json:"port" form:"port" db:"define:varchar(50);comment:主机端口号;defaultValue:''"`
	Env            string `json:"env" form:"env" db:"define:varchar(100);comment:环境变量;defaultValue:''"`
	User           string `json:"user" form:"user" db:"define:varchar(100);comment:用户名;defaultValue:''"`
	Ps             string `json:"ps" form:"ps" db:"define:varchar(100);comment:密码;defaultValue:''"`
	Ext            string `json:"ext" form:"ext" db:"define:varchar(255);comment:自定义配置信息;defaultValue:''"`
	Status         int    `json:"status" form:"status" db:"define:tinyint(1);comment:状态1正常2关闭3异常;defaultValue:0"`
	ChargeUserName string `json:"charge_user_name" form:"charge_user_name" db:"define:varchar(50);comment:负责人姓名;defaultValue:''"`
	StartTime      int    `json:"start_time" form:"start_time" db:"define:int;comment:开始时间;defaultValue:0"`
	EndTime        int    `json:"end_time" form:"end_time" db:"define:int;comment:结束时间;defaultValue:0"`
	Price          int    `json:"price" form:"price" db:"define:int;comment:价格;defaultValue:0"`
}

func (instance *Instance) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "服务-实例"

	return m
}
