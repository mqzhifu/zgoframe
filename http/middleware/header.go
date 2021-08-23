package httpmiddleware

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
	"zgoframe/http/request"
	"zgoframe/util"
)

func ProcessHeader()gin.HandlerFunc{
	return func(c *gin.Context) {

		header := HttpHeaderSureMapCovertSureStruct(c.Request.Header)

		if !request.CheckPlatformExist(header.SourceType){
			header.SourceType = request.PLATFORM_UNKNOW
		}
		header.AutoIp = c.Request.RemoteAddr

		c.Set("myheader",header)
		//request.Header{
		//	RequestId: c.GetHeader("request_id"),
		//}
		c.Next()
	}
}
//字符串 下划线转中划线
func StrCovertHttpHeader(str string)string{
	rsStr := ""
	arr := strings.Split(str,"_")
	if len(arr) <= 1{
		rsStr = util.StrFirstToUpper(str)
	}else{
		for _,v := range arr{
			rsStr += util.StrFirstToUpper(v) + "-"
		}
		rsStr = string([]byte(rsStr)[0:len(rsStr)-1])
	}

	rsStr = "X-"+rsStr
	return rsStr
}
//确定一个map 转换成 一个确定的struct
//map的key取值，从struce的json里取
//问题：目前仅支持一维
func HttpHeaderSureMapCovertSureStruct(inMap map[string][]string)request.Header{
	outStruct := request.Header{}

	//util.MyPrint("inMap:",inMap)
	stringMap := make(map[string]string)
	indexMap := make(map[int]string)
	typeOfOutStructArgs := reflect.TypeOf(outStruct)
	for i:=0;i<typeOfOutStructArgs.NumField();i++{
		structFiled := typeOfOutStructArgs.Field(i)
		structFiledTagName := structFiled.Tag.Get("json")
		structFiledName := structFiled.Name
		headerKey := StrCovertHttpHeader(structFiledTagName)
		//util.MyPrint("json-tag:",structFiledTagName," headerKey:",headerKey , " structFiledName:",structFiledName)
		valArr,ok := inMap[headerKey]

		stringMap[structFiledName] = ""
		indexMap[i] = ""
		if ok {
			stringMap[structFiledName] = valArr[0]
			indexMap[i] = structFiledName
		}
	}
	//util.MyPrint("stringMap:",stringMap)
	//typeOfOutStruct := reflect.TypeOf(outStruct)
	//util.MyPrint(outStruct)
	ValueOfOutStruct := reflect.ValueOf(&outStruct)
	for i:=0;i<ValueOfOutStruct.Elem().NumField();i++{
		fieldType := ValueOfOutStruct.Elem().Field(i).Type()
		if fieldType.String() == "int"{
			fieldValue, _ := strconv.ParseInt(stringMap[indexMap[i]], 10, 64)
			ValueOfOutStruct.Elem().Field(i).SetInt(fieldValue)
		}else if fieldType.String() == "string"{
			ValueOfOutStruct.Elem().Field(i).SetString(stringMap[indexMap[i]])
		}else{
			util.MyPrint("HttpHeaderSureMapCovertSureStruct err:type err ")
		}
		//util.MyPrint("type:",fieldType, " key:",ValueOfOutStruct.na)
	}


	//ValueOfOutStruct.Elem().FieldByName("RequestId").SetString("aaaaa")

	//for i := 0; i < typeOfOutStruct.NumField(); i++ {
	//	MyPrint(typeOfOutStruct.Field(i).Name)
	//}
	//cmsArg := make(map[string]string)
	//for i:=0;i<typeOfCmsArgs.NumField();i++{
	//	memVar := typeOfCmsArgs.Field(i)
	//	sqeNum := memVar.Tag.Get("seq")
	//	num , _:=strconv.Atoi(sqeNum)
	//	cmsArg[memVar.Name] = os.Args[num]
	//}
	//util.MyPrint(outStruct)
	//util.ExitPrint(1111)
	return outStruct
}
