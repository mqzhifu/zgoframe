package util

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)
//这个函数只是懒......
func MyPrint(a ...interface{}) (n int, err error) {
	return fmt.Println(a)
}
//这个函数只是懒......debug 调试使用
func ExitPrint(a ...interface{})   {
	fmt.Println(a)
	//fmt.Println("ExitPrint...22")
	os.Exit(999)
}
//输出复杂类型的数据，如：结构体
func MyComplexPrint(a ...interface{}) (n int, err error) {
	return fmt.Printf("%+v",a)
}
//一次获取N个空格，用于测试时 输出时 加些空格格式化内容
func GetSpaceStr(n int)string{
	str := ""
	for i:=0;i<n;i++{
		str += " "
	}
	return str
}
//获取一个随机数：int
func GetRandIntNum(max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
//获取一个随机数：int32
func GetRandInt32Num(max int32) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(int(max))
}
//获取一个随机整数:可设置范围
func GetRandIntNumRange(min int ,max int) int{
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
//自带的<strconv.Atoi>函数返回的是两个参数，很麻烦，这里简化，只返回一个参数
func Atoi(str string)int{
	num, _ := strconv.Atoi(str)
	return num
}
//类型转换：浮点转字符串
func FloatToString(number float32,little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(number), 'f', little, 64)
}
//类型转换：float64 -> 字符串
func Float64ToString(number float64,little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat( number, 'f', little, 64)
}
//类型转换：string -> float32
func StringToFloat(str string)float32{
	v1,_ := strconv.ParseFloat(str, 32)
	number  := float32(v1)
	return number
}
//指令行映射，根据使用者提供的map从指令行读取取，有查错，并映射进去
//练习了两个知识点：1给定一个struct，和一堆string，反射struct成员值，把string映射进map里 2从一个struct的tag中解析数据
func CmsArgs(data interface{})(argMap map[string]string,err error){
	//读取 data 类型 反射
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
		memVar := typeOfCmsArgs.Field(i)//获取结构体中的一个成员对象
		sqeNum := memVar.Tag.Get("seq")//读取出该成员的tag
		num , _:=strconv.Atoi(sqeNum)//转换成字符串
		cmsArg[memVar.Name] = os.Args[num]//根据成员名，写入map中
	}
	return cmsArg,nil
}
//获取本机的Ip地址
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

//func MapCovertStruct(inMap map[string]interface{},outStruct interface{})interface{}{
//	//fmt.Printf("%+v",inMap)
//	//fmt.Printf("%+v",outStruct)
//
//	setFiledValue := func(	outStruct interface{},name string , v interface{}) {
//		//MyPrint(name)
//
//		structValue := reflect.ValueOf(outStruct).Elem()
//		structFieldValue := structValue.FieldByName(name)
//
//		structFieldType := structFieldValue.Type() //结构体的类型
//		val := reflect.ValueOf(v)              //map值的反射值
//
//		var err error
//		//判断 结构体的元素类型 和 map元素的值类型
//		if structFieldType != val.Type() {
//			//类型不同，需要进行转换
//			val, err = TypeConversion(fmt.Sprintf("%v", v), structFieldValue.Type().Name()) //类型转换
//			if err != nil {
//				ExitPrint(err.Error())
//			}
//		}
//		//MyPrint(val,val.Type(),v)
//
//		structFieldValue.Set(val)
//	}
//	for k,v := range inMap{
//		//MyPrint("MapCovertStruct for range:",outStruct,k,v)
//		setFiledValue(outStruct,k,v)
//	}
//	//outStructV := reflect.ValueOf(outStruct)
//	//outStructT := reflect.TypeOf(outStruct)
//
//	return outStruct
//}
////类型转换
//func TypeConversion(value string, ntype string) (reflect.Value, error) {
//	if ntype == "string" {
//		return reflect.ValueOf(value), nil
//	} else if ntype == "time.Time" {
//		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
//		return reflect.ValueOf(t), err
//	} else if ntype == "Time" {
//		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
//		return reflect.ValueOf(t), err
//	} else if ntype == "int" {
//		i, err := strconv.Atoi(value)
//		return reflect.ValueOf(i), err
//	} else if ntype == "int8" {
//		i, err := strconv.ParseInt(value, 10, 64)
//		return reflect.ValueOf(int8(i)), err
//	} else if ntype == "int32" {
//		i, err := strconv.ParseInt(value, 10, 64)
//		return reflect.ValueOf(int64(i)), err
//	} else if ntype == "int64" {
//		i, err := strconv.ParseInt(value, 10, 64)
//		return reflect.ValueOf(i), err
//	} else if ntype == "float32" {
//		i, err := strconv.ParseFloat(value, 64)
//		return reflect.ValueOf(float32(i)), err
//	} else if ntype == "float64" {
//		i, err := strconv.ParseFloat(value, 64)
//		return reflect.ValueOf(i), err
//	}
//
//	//else if .......增加其他一些类型的转换
//
//	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
//}

