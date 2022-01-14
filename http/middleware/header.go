package httpmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	httpresponse "zgoframe/http/response"
	"zgoframe/util"
)
//预处理header：每个HTTP-API请求，都得加上对应的header，解析出来
func ProcessHeader()gin.HandlerFunc{
	res := httpresponse.Response{}
	return func(c *gin.Context) {
		//string header map 映射到 request.Header 结构体中
		header := HttpHeaderSureMapCovertSureStruct(c.Request.Header)
		//验证SourceType
		if !request.CheckPlatformExist(header.SourceType){
			header.SourceType = request.PLATFORM_UNKNOW
			res.Code = 501
			res.Msg = "SourceType unknow"
			c.AbortWithStatusJSON(500,res)
			return
		}

		if header.ProjectId <= 0{
			res.Code = 502
			res.Msg = "ProjectId <= 0"
			c.AbortWithStatusJSON(500,res)
			return
		}

		if header.Access == ""{
			res.Code = 503
			res.Msg = "ACCESS empty"
			c.AbortWithStatusJSON(500,res)
			return
		}

		project , empty := global.V.ProjectMng.GetById(header.ProjectId)
		if empty{
			res.Code = 504
			res.Msg = "projectId  empty"
			c.AbortWithStatusJSON(500,res)
			return
		}
		fmt.Println(project.Access ," - ", header.Access)
		if project.Access != header.Access{
			res.Code = 505
			res.Msg = "ACCESS  error"
			c.AbortWithStatusJSON(500,res)
			return
		}

		header.AutoIp = c.Request.RemoteAddr

		if header.RequestId == ""{
			header.RequestId = CreateOneRequestId()
		}

		if header.TraceId == ""{
			header.TraceId = CreateOneTraceId()
		}

		c.Set("myheader",header)
		c.Next()
	}
}

func CreateOneRequestId()string{
	return uuid.NewV4().String()
}

func CreateOneTraceId()string{
	return uuid.NewV4().String()
}
/*
//给定一个空的struct ，再给定一个有值的map ， 根据struct的tag ， 把map值 映射到 空 struct 中
//问题：
	1目前仅支持一维
	2并不是真正的struct 转 map ， 还需要struct 元素中定义tag
	3map里的key 是http header 模式，也就是X-XXX 开头这种
*/
func HttpHeaderSureMapCovertSureStruct(inMap map[string][]string)request.Header{
	outStruct := request.Header{}
	ValueOfOutStruct := reflect.ValueOf(&outStruct)
	//先读取 输出的 struct 反射信息
	typeOfOutStructArgs := reflect.TypeOf(outStruct)
	for i:=0;i<typeOfOutStructArgs.NumField();i++{
		//输出的 struct 成员对象
		structFiled := typeOfOutStructArgs.Field(i)
		//从 struct 成员对象 的tag 中的 json 中读取 key信息
		structFiledTagName := structFiled.Tag.Get("json")
		//structFiledName := structFiled.Name
		//json里直接读取的字符串还不能用，得转换成http header格式，X-Abc-Def 格式
		headerKey := util.StrCovertHttpHeader(structFiledTagName)
		//util.MyPrint("json-tag:",structFiledTagName," headerKey:",headerKey , " structFiledName:",structFiledName)
		//根据计算出的：struct json key =>header key  ，再从map中读取最终的值
		headerOneValArr,ok := inMap[headerKey]//这里之所以是个数组，因为header map 之前就是这么存的，回头优化
		headerOneVal := ""
		if !ok {
			continue
		}else{
			headerOneVal = headerOneValArr[0]
			if headerOneVal == ""{
				//util.MyPrint("headerKey empty:",headerKey)
				continue
			}
		}
		//读取该struct 字段值的类型
		fieldType := ValueOfOutStruct.Elem().Field(i).Type()
		//header传输就只有字符串，最多还能转换个int ，所以这里只处理了int string
		if fieldType.String() == "int"{
			//将value 由 string => int
			fieldValue, _ := strconv.ParseInt(headerOneVal, 10, 64)
			ValueOfOutStruct.Elem().Field(i).SetInt(fieldValue)
		}else if fieldType.String() == "string"{
			ValueOfOutStruct.Elem().Field(i).SetString(headerOneVal)
		}else{
			util.MyPrint("HttpHeaderSureMapCovertSureStruct err:type err ")
		}
	}
	//util.PrintStruct(outStruct,":")
	return outStruct
}
