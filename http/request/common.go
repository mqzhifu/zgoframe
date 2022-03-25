//http 请求公共处理
package request

import "github.com/dgrijalva/jwt-go"

//@description 解析token
type ParserToken struct {
	Token string `json:"token" form:"token"`
}

//@description http客户端请求头
type Header struct {
	Access     string         `json:"access"`      //使用网关时，不允许随意访问，得有key
	RequestId  string         `json:"request_id"`  //每次请求的唯一标识，响应时也会返回，如果请求方没有，后端会默认生成一个
	TraceId    string         `json:"trace_id"`    //追踪ID，主要用于链路追踪，如果请求方没有，后端会默认生成一个，跟request略像，但给后端使用
	SourceType int            `json:"source_type"` //请求方来源类型(pc h5 ios android vr spider unknow)，不同类型，不同JWT，原因：1手机端登陆后，PC端再登陆，互踢，无法共存。2越权，有些接口不允许互相访问
	ProjectId  int            `json:"project_id"`  //项目ID，所有的服务/项目/前端/App，均要先向管理员申请一个账号，才能用于日常请求
	Token      string         `json:"token"`       //JWT用户登陆令牌(HS256 对称算法，共享一个密钥)
	AutoIp     string         `json:"auto_ip"`     //获取不到请求方IP时，系统自动获取生成
	BaseInfo   HeaderBaseInfo `json:"base_info"`   //收集客户端的一些基础信息，json
}

//@description http客户端请求头-基础信息
type HeaderBaseInfo struct {
	AppVersion    string `json:"app_version"`    //app/前端/服务/项目 版本号
	OS            int    `json:"os"`             //win mac android ios
	OSVersion     string `json:"os_version"`     //win7 win9 mac10 android9
	Device        string `json:"device"`         //ipad iphone huawei mi chrome firefox ie
	DeviceVersion string `json:"device_version"` //mi8 hongmi7 ios8 ios9 ie8 ie9
	Lat           string `json:"lat"`            //纬度
	Lon           string `json:"lon"`            //经度
	DeviceId      string `json:"device_id"`      //设备ID
	DPI           string `json:"dpi"`            //分辨率
	Ip            string `json:"ip"`             //请求方的IP
	Referer       string `json:"referer"`        //页面来源
}

//@description jwt 容器内容
type CustomClaims struct {
	ProjectId  int    `json:"project_id"`  //项目ID
	SourceType int    `json:"source_type"` //来源
	Id         int    `json:"id"`          //用户ID
	Username   string `json:"username"`    //用户名
	NickName   string `json:"nick_name"`   //用户昵称

	jwt.StandardClaims `swaggerignore:"true"`

	//BufferTime int64 `json:"buffer_time"`
	//AuthorityId string
	//UUID        uuid.UUID
}

//@description 系统管理员操作，需要2次验证
type SystemConfig struct {
	Username string `json:"username" form:"username"` //用户名
	Password string `json:"password" form:"password"` //密码
}

//@description 3方登陆
type RLoginThird struct {
	Register
	ThirdId      string //3方平台用户ID
	PlatformType int    //3方平台类型
}

//@description 分页
type PageInfo struct {
	Page     int `json:"page" form:"page"`         //当前页数
	PageSize int `json:"pageSize" form:"pageSize"` //每页多少条记录
}

//type GetById struct {
//	Id float64 `json:"id" form:"id"`
//}
//
//type IdsReq struct {
//	Ids []int `json:"ids" form:"ids"`
//}
//
//// Get role by id structure
//type GetAuthorityId struct {
//	AuthorityId string
//}

type Empty struct{}

// Casbin info structure
type CasbinInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// Casbin structure for input parameters
type CasbinInReceive struct {
	AuthorityId string       `json:"authorityId"`
	CasbinInfos []CasbinInfo `json:"casbinInfos"`
}
