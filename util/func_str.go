package util
//公共函数：字符串 操作
import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"regexp"
	"strings"
	"unicode"
	"fmt"
)

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
//将一个字符串转换成MD5
func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}
//将一个字符串转换成MD5
func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
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


//驼峰式 转 下划线 式,针对普通字符串
func CamelToSnake2(marshalled []byte)[]byte{
	var keyMatchRegex = regexp.MustCompile(`(\w+)`)
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

//驼峰式 转 下划线 式,针对json串
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

//字符串 下划线转中划线，同时每个单词首字母转大写，最后每个字符串开头再加上：X-
func StrCovertHttpHeader(str string)string{
	rsStr := ""
	arr := strings.Split(str,"_")
	if len(arr) <= 1{//就一个单词，证明没有 _ ，直接把首字母大写就行了
		rsStr = StrFirstToUpper(str)
	}else{
		for _,v := range arr{
			rsStr += StrFirstToUpper(v) + "-"
		}
		rsStr = string([]byte(rsStr)[0:len(rsStr)-1])
	}

	rsStr = "X-"+rsStr
	return rsStr
}

