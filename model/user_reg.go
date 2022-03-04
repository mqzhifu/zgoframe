package model

const (
	USER_REG_TYPE_EMAIL 	= 1
	USER_REG_TYPE_NAME 		= 2
	USER_REG_TYPE_MOBILE 	= 3
	USER_REG_TYPE_THIRD 	= 4
	USER_REG_TYPE_GUEST 	= 5

	USER_TYPE_THIRD_WEIBO 	= 1
	USER_TYPE_THIRD_WECHAT 	= 2
	USER_TYPE_THIRD_FACEBOOK = 3
	USER_TYPE_THIRD_GOOGLE 	= 4
	USER_TYPE_THIRD_TWITTER = 5
	USER_TYPE_THIRD_YOUTOBE = 6
	USER_TYPE_THIRD_QQ 		= 7


	CHANNEL_DEFAULT = 1
)

type UserReg struct {
	MODEL
	ProjectId      	int   			`json:"project_id" db:"define:int;comment:project_id;defaultValue:0"  `
	Uid        		int 			`json:"uid" db:"define:int;comment:uid;defaultValue:0"`
	Type 			int 			`json:"type" db:"define:tinyint(1);comment:类型;defaultValue:0" `
	ThirdType 		int				`json:"third_type" db:"define:tinyint(1);comment:三方平台用户ID;defaultValue:0"`
	Channel 		int 			`json:"channel" db:"define:tinyint(1);comment:推广渠道;defaultValue:0"`
	Ip 				string 			`json:"ip" db:"define:varchar(50);comment:请求方传输IP;defaultValue:''"`
	AutoIp			string			`json:"auto_ip" db:"define:varchar(50);comment:程序自己计算的IP;defaultValue:''"`
	AppVersion 		string			`json:"app_version" db:"define:varchar(50);comment:APP版本;defaultValue:''"`
	SourceType 		int				`json:"source_type" db:"define:tinyint(1);comment:来源类型;defaultValue:0"`//pc h5 ios android vr unknow
	Os 				int				`json:"os" db:"define:tinyint(1);comment:操作系统;defaultValue:0"`
	OsVersion 		string			`json:"os_version"  db:"define:varchar(50);comment:操作系统版本;defaultValue:''"`
	Device			string			`json:"device" db:"define:varchar(50);comment:设备名称;defaultValue:''"`
	DeviceVersion 	string			`json:"device_version" db:"define:varchar(50);comment:设备版本;defaultValue:''"`
	Lat 			string			`json:"lat" db:"define:varchar(50);comment:伟度;defaultValue:''"`
	Lon 			string			`json:"lon" db:"define:varchar(50);comment:经度;defaultValue:''"`
	DeviceId 		string		`json:"device_id" db:"define:varchar(50);comment:设备ID;defaultValue:''"`
	Dpi 			string		`json:"dpi" db:"define:varchar(50);comment:分辨率;defaultValue:''"`
	Referer 		string		`json:"referer" db:"define:varchar(255);comment:页面来源;defaultValue:''"`

}

func(userReg *UserReg) TableOptions()map[string]string{
	m := make(map[string]string)
	m["comment"] = "用户注册信息"

	return m
}
