package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	"zgoframe/model"
)

type ParserTokenData struct {
	Claims *CustomClaims	//解析后的token里面的值
	User *model.User	//解析后的token，再反查userinfo
	Token string		//需要解析的TOKEN
	SourceType int		//需要解析的来源类型
	NewToken string		//失效了，但在缓存期内，重新生成了一个新的token
}

type Header struct {
	Access     		string  `json:"access"`	//使用网关时，不允许随意访问，得有key
	RequestId 		string	`json:"request_id"`
	TraceId 		string	`json:"trace_id"`
	SourceType 		int		`json:"source_type"`	//pc h5 ios android vr spider unknow
	ProjectId		int		`json:"project_id"`
	Token 			string	`json:"token"`
	AppVersion 		string	`json:"app_version"`	//app版本/前端版本
	OS 				int		`json:"os"`				//win mac android ios
	OSVersion 		string	`json:"os_version"`		//win7 win9 mac10 android9
	Device			string	`json:"device"`			// ipad iphone huawei mi
	DeviceVersion 	string	`json:"device_version"`
	Lat 			string	`json:"lat"`			//纬度
	Lon 			string	`json:"lon"`			//经度
	DeviceId 		string	`json:"device_id"`
	DPI 			string	`json:"dpi"`			//分辨率
	Ip 				string	`json:"ip"`
	AutoIp 			string	`json:"auto_ip"`		//获取不到请求方IP时，系统自动生成
	Referer 		string	`json:"referer"`		//页面来源
}

const  (
	PLATFORM_PC = 1
	PLATFORM_H5 = 2
	PLATFORM_ANDROID = 3
	PLATFORM_IOS = 4
	PLATFORM_UNKNOW = 5
)

func GetPlatformList()[]int{
	list := []int{PLATFORM_PC, PLATFORM_H5, PLATFORM_ANDROID,PLATFORM_IOS}
	return list
}
func CheckPlatformExist(env int)bool{
	list := GetPlatformList()
	for _,v :=range list{
		if v == env{
			return true
		}
	}
	return false
}

func GetMyHeader(c *gin.Context)Header{
	myHeaderInterface , exists := c.Get("myheader")
	if !exists{
		global.V.Zap.Error("myheader empty")
	}
	myHeader := myHeaderInterface.(Header)
	return myHeader
}

func GetParserTokenData(c *gin.Context)(parserTokenData ParserTokenData,err error){
	parserTokenDataInter , exists := c.Get("parserTokenData")
	if !exists{
		global.V.Zap.Error("parserTokenData empty")
		return parserTokenData,errors.New("parserTokenData empty")
	}
	parserTokenData = parserTokenDataInter.(ParserTokenData)
	return parserTokenData,nil
}

func GetUid(c *gin.Context)(int,error){
	user,err := GetUser(c)
	if err != nil{
		return 0,err
	}
	return user.Id,nil
}
//有4种方式获取：
//1. 从token解出来的结构体内获取
//2. 从token解出来的结构体内，再从DB中获取
//3. header中也可以取这个值
//4. 请求方的body中直接附加此值
func GetAppId(c *gin.Context) (int,error) {
	CustomClaims ,err := GetClaims(c)
	if err != nil{
		return 0,errors.New("从Gin的Context中获取从jwt解析出来的user_appID失败, 请检查路由是否使用jwt中间件")
	}

	return CustomClaims.ProjectId,nil
}

func GetUser(c *gin.Context)(user *model.User,err error){
	parserTokenData ,err := GetParserTokenData(c)
	if err != nil{
		return user,err
	}
	return parserTokenData.User,nil
}

func GetClaims(c *gin.Context)(customClaims *CustomClaims,err error){
	parserTokenData ,err := GetParserTokenData(c)
	if err != nil{
		return customClaims,err
	}
	return parserTokenData.Claims,nil
}

