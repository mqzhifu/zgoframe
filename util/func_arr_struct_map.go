package util
//公用函数：数组、集合、结构体、字节
import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"reflect"
)
//打印：一个结构体（仅支持一维结构体）
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
//将一个map转换成一个数组
func MapCovertArr( myMap map[int]int) (arr []int){
	for _,v := range myMap {
		arr = append(arr,v)
	}
	return arr
}
//数组转换成map
func ArrCovertMap(arr []int )map[int]int{
	mapArr := make(map[int]int)
	for k,v := range arr {
		mapArr[k] = v
	}
	return mapArr
}
//将一个：一给数组(string)转成成 数组(int)
func ArrStringCoverArrInt(arr []string )(arr2 []int){
	for i:=0;i<len(arr);i++{
		arr2 = append(arr2, Atoi(arr[i]))
	}
	return arr2
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
//把一维int 数组，转换成一个字符串
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
//结构体转map ，这里实际是借用了json类中转
func StructCovertMap(inStruct interface{})interface{}{
	jsonStr ,_:= json.Marshal(inStruct)
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {

	}
	return mapResult
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
//判断一个元素 int，在一个数组中的位置
func ElementInArrIndex(arr []int ,element int )int{
	for i:=0;i<len(arr);i++{
		if arr[i] == element{
			return i
		}
	}
	return -1
}
//判断一个元素 string，在一个数组中的位置
func ElementStrInArrIndex(arr []string ,element string )int{
	for i:=0;i<len(arr);i++{
		if arr[i] == element{
			return i
		}
	}
	return -1
}
//BytesToInt32
func BytesToInt32(bys []byte) int {
	byteBuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(byteBuff, binary.BigEndian, &data)
	return int(data)
}
//Int32ToBytes
func Int32ToBytes(n int32) []byte {
	//x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}
