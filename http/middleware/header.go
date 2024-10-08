package httpmiddleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/util"
)

// 预处理header：每个HTTP-API请求，都得加上对应的header，解析出来
func Header() gin.HandlerFunc {

	return func(c *gin.Context) {
		global.V.Base.Zap.Debug("http middleware <header> start:")

		//string header map 映射到 request.Header 结构体中
		header := HttpHeaderSureMapCovertSureStruct(c.Request.Header)
		header.AutoIp = c.Request.RemoteAddr

		//h := request.HeaderBaseInfo{}
		//s, _ := json.Marshal(h)

		header.ServerReceiveTime = util.GetNowTimeSecondToInt()
		if header.RequestId == "" {
			header.RequestId = CreateOneRequestId()
		}

		if header.TraceId == "" {
			header.TraceId = CreateOneTraceId()
		}

		c.Set("myHeader", header)

		global.V.Base.Zap.Debug("http middleware <header>  finish.")

		c.Next()
	}
}

func CreateOneRequestId() string {
	return uuid.NewV4().String()
}

func CreateOneTraceId() string {
	return uuid.NewV4().String()
}

/*
//给定一个空的struct ，再给定一个有值的map ， 根据struct的tag ， 把map值 映射到 空 struct 中
//问题：
 1. 目前仅支持一维
 2. 并不是真正的struct 转 map ， 还需要struct 元素中定义tag
 3. map里的key 是http header 模式，也就是X-XXX 开头这种
*/
func HttpHeaderSureMapCovertSureStruct(inMap map[string][]string) request.HeaderRequest {
	outStruct := request.HeaderRequest{}
	ValueOfOutStruct := reflect.ValueOf(&outStruct)
	//先读取 输出的 struct 反射信息
	typeOfOutStructArgs := reflect.TypeOf(outStruct)
	for i := 0; i < typeOfOutStructArgs.NumField(); i++ {
		//输出的 struct 成员对象
		structFiled := typeOfOutStructArgs.Field(i)
		//从 struct 成员对象 的tag 中的 json 中读取 key信息
		structFiledTagName := structFiled.Tag.Get("json")
		//structFiledName := structFiled.Name
		//json里直接读取的字符串还不能用，得转换成http header格式，X-Abc-Def 格式
		headerKey := util.StrCovertHttpHeader(structFiledTagName)
		//util.MyPrint("json-tag:",structFiledTagName," headerKey:",headerKey , " structFiledName:",structFiledName)
		//根据计算出的：struct json key =>header key  ，再从map中读取最终的值
		headerOneValArr, ok := inMap[headerKey] //这里之所以是个数组，因为header map 之前就是这么存的，回头优化
		headerOneVal := ""
		if !ok {
			continue
		} else {
			headerOneVal = headerOneValArr[0]
			if headerOneVal == "" {
				//util.MyPrint("headerKey empty:",headerKey)
				continue
			}
		}
		//读取该struct 字段值的类型
		fieldType := ValueOfOutStruct.Elem().Field(i).Type()
		//header传输就只有字符串，最多还能转换个int ，所以这里只处理了int string
		if fieldType.String() == "int" {
			//将value 由 string => int
			fieldValue, _ := strconv.ParseInt(headerOneVal, 10, 64)
			ValueOfOutStruct.Elem().Field(i).SetInt(fieldValue)
		} else if fieldType.String() == "string" {
			ValueOfOutStruct.Elem().Field(i).SetString(headerOneVal)
		} else if fieldType.String() == "request.HeaderBaseInfo" {
			hbi := request.HeaderBaseInfo{}
			err := json.Unmarshal([]byte(headerOneVal), &hbi)
			util.MyPrint("HeaderBaseInfo : err ", err, " hbi")
			if err == nil {
				ValueOfOutStruct.Elem().Field(i).Set(reflect.ValueOf(hbi))
			}
		} else {

			util.MyPrint("HttpHeaderSureMapCovertSureStruct err:type err ", fieldType.String())
		}
	}
	//util.PrintStruct(outStruct,":")
	return outStruct
}
