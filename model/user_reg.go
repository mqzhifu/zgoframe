package model

type UserReg struct {
	MODEL
	ProjectId     int    `json:"project_id" db:"define:int;comment:project_id;defaultValue:0"  `
	SourceType    int    `json:"source_type" db:"define:tinyint(1);comment:来源类型;defaultValue:0"`
	Uid           int    `json:"uid" db:"define:int;comment:uid;defaultValue:0"`
	Type          int    `json:"type" db:"define:tinyint(1);comment:类型 1email2name3mobile3third4guest;defaultValue:0" `
	ThirdType     int    `json:"third_type" db:"define:varchar(50);comment:三方平台类型,参数常量USER_TYPE_THIRD;defaultValue:''"`
	Channel       int    `json:"channel" db:"define:tinyint(1);comment:推广渠道1平台自己;defaultValue:0"`
	Ip            string `json:"ip" db:"define:varchar(50);comment:请求方传输IP;defaultValue:''"`
	AutoIp        string `json:"auto_ip" db:"define:varchar(50);comment:程序自己计算的IP;defaultValue:''"`
	Province      int    `json:"province" db:"define:int;comment:project_id;defaultValue:0"`
	City          int    `json:"city" db:"define:int;comment:project_id;defaultValue:0"`
	County        int    `json:"county" db:"define:int;comment:project_id;defaultValue:0"`
	Town          int    `json:"town" db:"define:int;comment:project_id;defaultValue:0"`
	AreaDetail    string `json:"area_detail"  db:"define:varchar(255);comment:页面来源;defaultValue:"`
	AppVersion    string `json:"app_version" db:"define:varchar(50);comment:APP版本;defaultValue:''"`
	Os            string `json:"os" db:"define:string(50);comment:操作系统;defaultValue:0"`
	OsVersion     string `json:"os_version"  db:"define:varchar(50);comment:操作系统版本;defaultValue:''"`
	Device        string `json:"device" db:"define:varchar(50);comment:设备名称;defaultValue:''"`
	DeviceVersion string `json:"device_version" db:"define:varchar(50);comment:设备版本;defaultValue:''"`
	Lat           string `json:"lat" db:"define:varchar(50);comment:伟度;defaultValue:''"`
	Lon           string `json:"lon" db:"define:varchar(50);comment:经度;defaultValue:''"`
	DeviceId      string `json:"device_id" db:"define:varchar(50);comment:设备ID;defaultValue:''"`
	Dpi           string `json:"dpi" db:"define:varchar(50);comment:分辨率;defaultValue:''"`
	Referer       string `json:"referer" db:"define:varchar(255);comment:页面来源;defaultValue:''"`
}

func (userReg *UserReg) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户注册信息"

	return m
}
