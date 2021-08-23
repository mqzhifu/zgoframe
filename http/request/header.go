package request

import "github.com/gin-gonic/gin"

type Header struct {
	RequestId 		string	`json:"request_id"`
	TraceId 		string	`json:"trace_id"`
	SourceType 		int		`json:"source_type"`	//pc h5 ios android vr spider unknow
	AppId			int		`json:"app_id"`
	Token 			string	`json:"token"`
	AppVersion 		string	`json:"app_version"`
	OS 				int		`json:"os"`
	OSVersion 		string	`json:"os_version"`
	Device			string	`json:"device"`
	DeviceVersion 	string	`json:"device_version"`
	Lat 			string	`json:"lat"`
	Lon 			string	`json:"lon"`
	DeviceId 		string	`json:"device_id"`
	DPI 			string	`json:"dpi"`
	Ip 				string	`json:"ip"`
	AutoIp 			string	`json:"auto_ip"`
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
	myHeaderInterface , _ := c.Get("myheader")
	myHeader := myHeaderInterface.(Header)
	return myHeader
}