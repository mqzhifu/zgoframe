package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/model"
)

//type ParserTokenData struct {
//	Claims CustomClaims //解析后的token里面的值
//	User   model.User    //解析后的token，再反查userinfo
//	Token  string        //需要解析的TOKEN
//SourceType int           //需要解析的来源类型
//NewToken   string        //失效了，但在缓存期内，重新生成了一个新的token
//}

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

const (
	PLATFORM_MAC_PC_BROWSER = 11
	PLATFORM_MAC_APP        = 12

	PLATFORM_WIN_PC_BROWSER = 22
	PLATFORM_WIN_APP        = 23

	PLATFORM_ANDROID_H5_BROWSER = 31
	PLATFORM_ANDROID_APP        = 32

	PLATFORM_IOS_H5_BROWSER = 41
	PLATFORM_IOS_APP        = 42

	PLATFORM_UNKNOW = 99
)

func GetPlatformList() []int {
	list := []int{PLATFORM_MAC_PC_BROWSER, PLATFORM_WIN_PC_BROWSER, PLATFORM_ANDROID_H5_BROWSER, PLATFORM_IOS_H5_BROWSER, PLATFORM_ANDROID_APP, PLATFORM_IOS_APP, PLATFORM_MAC_APP, PLATFORM_WIN_APP, PLATFORM_UNKNOW}
	return list
}

func CheckPlatformExist(env int) bool {
	list := GetPlatformList()
	for _, v := range list {
		if v == env {
			return true
		}
	}
	return false
}

func GetMyHeader(c *gin.Context) Header {
	myHeaderInterface, exists := c.Get("myheader")
	if !exists {
		global.V.Zap.Error("myheader empty")
	}
	myHeader := myHeaderInterface.(Header)
	return myHeader
}

//func GetParserTokenData(c *gin.Context) (parserTokenData ParserTokenData, err error) {
//	parserTokenDataInter, exists := c.Get("parserTokenData")
//	if !exists {
//		global.V.Zap.Error("parserTokenData empty")
//		return parserTokenData, errors.New("parserTokenData empty")
//	}
//	parserTokenData = parserTokenDataInter.(ParserTokenData)
//	return parserTokenData, nil
//}

//1. 从token中解出来的值里获取
//2. 从DB中获取
func GetUid(c *gin.Context) (int, error) {
	user, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

//有4种方式获取：
//1. 从token解出来的结构体内获取
//2. 从token解出来的结构体内，再从DB中获取
//3. header中也可以取这个值
func GetProjectId(c *gin.Context) (int, error) {
	customClaims, err := GetClaims(c)
	if err != nil {
		return 0, errors.New("Claims key not exist")
	}

	return customClaims.ProjectId, nil
}

func GetSourceType(c *gin.Context) (int, error) {
	customClaims, err := GetClaims(c)
	if err != nil {
		return 0, errors.New("Claims key not exist")
	}

	return customClaims.SourceType, nil
}

func GetUser(c *gin.Context) (user model.User, err error) {
	u, exist := c.Get("user")
	if !exist {
		return user, errors.New("not exist")
	}
	user = u.(model.User)
	return user, nil
}

func GetClaims(c *gin.Context) (customClaims CustomClaims, err error) {
	cc, exist := c.Get("customClaims")
	if !exist {
		return customClaims, errors.New("not exist")
	}
	customClaims = cc.(CustomClaims)
	return customClaims, nil
}
