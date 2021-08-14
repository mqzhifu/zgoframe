package util

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"fmt"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}



//这个函数只是懒......
func MyPrint(a ...interface{}) (n int, err error) {
	if LogLevelFlag == LOG_LEVEL_DEBUG{
		return fmt.Println(a)
	}
	return
}

func PanicPrint(a ...interface{})   {
	if LogLevelFlag == LOG_LEVEL_DEBUG{
		fmt.Println(a)
	}
	panic(a[0])
}

//debug 调试使用
func ExitPrint(a ...interface{})   {
	fmt.Println(a)
	//fmt.Println("ExitPrint...22")
	os.Exit(-22)
}
//获取一个随机整数
func GetRandIntNum(max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
func GetRandInt32Num(max int32) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(int(max))
}
//获取一个随机整数 范围
func GetRandIntNumRange(min int ,max int) int{
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
//判断一个字符串是否为空，包括  空格
func CheckStrEmpty(str string)bool{
	if str == ""{
		return true
	}
	str = strings.Trim(str," ")
	if str == ""{
		return true
	}
	return false
}
//检查一个文件是否已存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename);os.IsNotExist(err){
		exist = false
	}
	return exist
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MapCovertArr( myMap map[int]int) (arr []int){
	for _,v := range myMap {
		arr = append(arr,v)
	}
	return arr
}

func ArrCovertMap(arr []int )map[int]int{
	mapArr := make(map[int]int)
	for k,v := range arr {
		mapArr[k] = v
	}
	return mapArr
}

func ArrStringCoverArrInt(arr []string )(arr2 []int){
	for i:=0;i<len(arr);i++{
		arr2 = append(arr2, Atoi(arr[i]))
	}
	return arr2
}
func GetSpaceStr(n int)string{
	str := ""
	for i:=0;i<n;i++{
		str += " "
	}
	return str
}
//检查已经make过的，二维map int 类型，是否为空
func CheckMap2IntIsEmpty(hashMap map[int]map[int]int)bool{
	if len(hashMap) == 0{
		return true
	}

	for _,v := range hashMap{
		if len(v) > 0{
			return false
		}
	}
	return true

}

func ArrCoverStr(arr []int,IdsSeparation string)string{
	if len(arr) == 0{
		ExitPrint("ArrCoverStr arr len = 0")
	}
	str := ""
	for _,v := range arr{
		str +=  strconv.Itoa(v) + IdsSeparation
	}
	str = str[0:len(str)-1]
	return str
}

func StructCovertMap(inStruct interface{})interface{}{
	jsonStr ,_:= json.Marshal(inStruct)
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {

	}
	return mapResult
}
//strconv.Atoi 返回的是两个参数，很麻烦，这里简化，只返回一个参数
func Atoi(str string)int{
	num, _ := strconv.Atoi(str)
	return num
}
//将字符串的首字母转大写
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return str
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122  {
		strArry[0] = strArry[0] - 32
	}
	return string(strArry)
}
// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
//func trunVarJson(marshalled []byte)string{
//	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
//	//marshalled, err := json.Marshal(room)
//	converted := keyMatchRegex.ReplaceAllFunc(
//		marshalled,
//		func(match []byte) []byte {
//			matchStr := string(match)
//			key := matchStr[1 : len(matchStr)-2]
//			resKey := Lcfirst(Case2Camel(key))
//			return []byte(`"` + resKey + `":`)
//		},
//	)
//
//	return string(converted)
//}
func GetNowDateMonth()string{
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//str := strconv.Itoa(now.Year()) +  strconv.Itoa(int(now.Month())) + strconv.Itoa( now.Day())
	str := time.Now().Format("200601")
	return str
}

func GetNowDate()string{
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//str := strconv.Itoa(now.Year()) +  strconv.Itoa(int(now.Month())) + strconv.Itoa( now.Day())
	str := time.Now().Format("20060102")
	return str
}
func GetNowDateHour()string{
	//ss := time.Now().Format("20060102")
	//ExitPrint(ss)
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//date := GetNowDate()
	//str := date + strconv.Itoa(now.Hour())
	str := time.Now().Format("2006010215")
	return str
}

func GetDateHour(now int64)string{
	//ss := time.Now().Format("20060102")
	//ExitPrint(ss)
	//cstZone := time.FixedZone("CST", 8*3600)
	//now := time.Now().In(cstZone)
	//now := time.Now()
	//date := GetNowDate()
	//str := date + strconv.Itoa(now.Hour())
	tm := time.Unix(now, 0)
	str := tm.Format("2006010215")
	return str
}




//驼峰式 转 下划线 式
func CamelToSnake(marshalled []byte)[]byte{
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted
}

//将字符串的首字母转大写
func StrFirstToLower(str string) string {
	if len(str) < 1 {
		return str
	}
	strArry := []rune(str)
	if strArry[0] >= 65 && strArry[0] <= 90  {
		strArry[0] = strArry[0] + 32
	}
	return string(strArry)
}
func CmsArgs(data interface{})(argMap map[string]string,err error){
	typeOfCmsArgs := reflect.TypeOf(data)
	if len(os.Args) < typeOfCmsArgs.NumField() + 1{
		errInfo := "os.Args len < "+ strconv.Itoa(typeOfCmsArgs.NumField()) + " , eg:"
		for i:=0;i<typeOfCmsArgs.NumField();i++{
			memVar := typeOfCmsArgs.Field(i)
			errInfo += memVar.Tag.Get("err") + " ,"
		}
		return argMap,errors.New(errInfo)
	}
	cmsArg := make(map[string]string)
	for i:=0;i<typeOfCmsArgs.NumField();i++{
		memVar := typeOfCmsArgs.Field(i)
		sqeNum := memVar.Tag.Get("seq")
		num , _:=strconv.Atoi(sqeNum)
		cmsArg[memVar.Name] = os.Args[num]
	}
	return cmsArg,nil
}
func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func ReadLine(fileName string) ([]string,error){
	f, err := os.Open(fileName)
	if err != nil {
		return nil,err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result,nil
			}
			return nil,err
		}
		result = append(result,line)
	}
	return result,nil
}
//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
func GetNowMillisecond()int64{
	return time.Now().UnixNano() / 1e6
}
func MapCovertStruct(inMap map[string]interface{},outStruct interface{})interface{}{
	fmt.Printf("%+v",inMap)
	fmt.Printf("%+v",outStruct)

	setFiledValue := func(	outStruct interface{},name string , v interface{}) {
		MyPrint(name)
		structValue := reflect.ValueOf(outStruct).Elem()
		structFieldValue := structValue.FieldByName(name)

		structFieldType := structFieldValue.Type() //结构体的类型
		val := reflect.ValueOf(v)              //map值的反射值

		var err error
		if structFieldType != val.Type() {
			val, err = TypeConversion(fmt.Sprintf("%v", v), structFieldValue.Type().Name()) //类型转换
			if err != nil {
				ExitPrint(err.Error())
			}
		}
		//MyPrint(val,val.Type(),v)

		structFieldValue.Set(val)
	}
	for k,v := range inMap{
		//MyPrint("MapCovertStruct for range:",outStruct,k,v)
		setFiledValue(outStruct,k,v)
	}
	//outStructV := reflect.ValueOf(outStruct)
	//outStructT := reflect.TypeOf(outStruct)

	return outStruct
}


//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
//判断一个元素，在一个数组中的位置
func ElementInArrIndex(arr []int ,element int )int{
	for i:=0;i<len(arr);i++{
		if arr[i] == element{
			return i
		}
	}
	return -1
}

func ElementStrInArrIndex(arr []string ,element string )int{
	for i:=0;i<len(arr);i++{
		if arr[i] == element{
			return i
		}
	}
	return -1
}

func FloatToString(number float32,little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(number), 'f', little, 64)
}

func Float64ToString(number float64,little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat( number, 'f', little, 64)
}

func GetNowTimeSecondToInt()int{
	return int( time.Now().Unix() )
}

func GetNowTimeSecondToInt64()int64{
	return time.Now().Unix()
}

func StringToFloat(str string)float32{
	v1,_ := strconv.ParseFloat(str, 32)
	number  := float32(v1)
	return number
}
//func ParseStuctDesc(mystruct interface{})map[string]string{
//	//httpReqBusinessStruct := gamematch.HttpReqBusiness{}
//	rs := make(map[string]string)
//	types := reflect.TypeOf(&mystruct)
//	for i:=0 ; i < types.Elem().NumField() ; i++{
//		field := types.Elem().Field(i)
//		tagName1 := field.Tag.Get("desc")
//		rs[field.Name] = tagName1
//	}
//	return rs
//}
//var GoRoutineList = make( map[string]int )
//func AddRoutineList(name string){
//	GoRoutineList[name] = GetNowTimeSecondToInt()
//}
func GetLocalIp()(ip string,err error){
	netInterfaces, err := net.Interfaces()
	//MyPrint(netInterfaces, err)
	if err != nil {
		return ip,errors.New("net.Interfaces failed, err:" +  err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(),nil
					}
				}
			}
		}
	}

	return ip,nil
}
//在一个：一维数组中，找寻最大数
func FindMaxNumInArrFloat32(arr []float32  )float32{
	number := arr[0]
	for _,v := range arr{
		if v > number{
			number = v
		}
	}
	return number
}
//将一个完整的URL地址，以 <?>号分开，取<?>后面的参数
func UriTurnPath (uri string)string{
	n := strings.Index(uri,"?")
	if  n ==  - 1{
		return uri
	}
	uriByte := []byte(uri)
	path := uriByte[0:n]
	return string(path)
}
//在一个：一维数组中，找寻最小数
func FindMinNumInArrFloat32(arr []float32  )float32{
	number := arr[0]
	for _,v := range arr{
		if v < number{
			number = v
		}
	}
	return number
}
//从一个固定数字容器内，随机抽取出2个元素，且不重复
func getSomeRandNumByContainer(){

}
//一个结构体转成字符串，一般用于输出调度
func PrintStruct(mystruct interface{},separator string ){
	t := reflect.TypeOf(mystruct)
	v := reflect.ValueOf(mystruct)
	//str := ""
	for k := 0; k < t.NumField(); k++ {
		//str += t.Field(k).Name + separator + string(v.Field(k).Interface())
		//fmt.Printf("%s -- %v \n", t.Field(k).Name, v.Field(k).Interface())
		MyPrint(t.Field(k).Name,separator,v.Field(k).Interface())
	}
}
//遍历一个目录的所有文件列表，但 子目录不处理
func GetFileListByDir(path string)[]string {
	var fileList []string
	fs,err := ioutil.ReadDir(path)
	if err != nil{
		MyPrint("GetFileListByDir err:",err.Error())
		return fileList
	}
	for _,file:=range fs{
		if file.IsDir(){
			//fmt.Println(path+file.Name())
			//GetFileListByDir(path+file.Name()+"/")
		}else{
			//fmt.Println(path+file.Name())
			fileList = append(fileList,file.Name())
		}
	}
	return fileList
}

//4舍5入，保留2位小数
//func round(x float32)string{
//	numberStr := FloatToString(x,4)
//	numberStrSplit :=  strings.Split(numberStr,".")
//	if len(numberStrSplit) == 1{
//		return numberStrSplit[0]
//	}
//	numberLittleStrByte := []byte(numberStrSplit[1])
//	numberLittle := Atoi(numberStrSplit[1])
//	if len(numberLittleStrByte) == 4{//4位小数
//		three := numberLittleStrByte[0] + numberLittleStrByte[1] + numberLittleStrByte[2]
//		if strByte[4] >= 5{
//			three = Atoi(three) + 1
//		}
//	}else{//3位小数
//
//	}
//
//}