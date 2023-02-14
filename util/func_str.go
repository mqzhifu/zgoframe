package util

//公共函数：字符串 操作
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"zgoframe/http/request"
)

//将字符串的首字母转大写
func StrFirstToLower(str string) string {
	if len(str) < 1 {
		return str
	}
	strArry := []rune(str)
	if strArry[0] >= 65 && strArry[0] <= 90 {
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
	if strArry[0] >= 97 && strArry[0] <= 122 {
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
func CheckStrEmpty(str string) bool {
	if str == "" {
		return true
	}
	str = strings.Trim(str, " ")
	if str == "" {
		return true
	}
	return false
}

//驼峰式 转 下划线 式,针对普通字符串
func CamelToSnake2(marshalled []byte) []byte {
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
func CamelToSnake(marshalled []byte) []byte {
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
func UriTurnPath(uri string) string {
	n := strings.Index(uri, "?")
	if n == -1 {
		return uri
	}
	uriByte := []byte(uri)
	path := uriByte[0:n]
	return string(path)
}

//字符串 下划线转中划线，同时每个单词首字母转大写，最后每个字符串开头再加上：X-
func StrCovertHttpHeader(str string) string {
	rsStr := ""
	arr := strings.Split(str, "_")
	if len(arr) <= 1 { //就一个单词，证明没有 _ ，直接把首字母大写就行了
		rsStr = StrFirstToUpper(str)
	} else {
		for _, v := range arr {
			rsStr += StrFirstToUpper(v) + "-"
		}
		rsStr = string([]byte(rsStr)[0 : len(rsStr)-1])
	}

	rsStr = "X-" + rsStr
	return rsStr
}

func HttpHeaderSureStructCovertSureMap(response request.HeaderResponse) (outMap map[string]string) {
	outMap = make(map[string]string)
	//先读取 输出的 struct 反射信息
	typeOfResponse := reflect.TypeOf(response)
	valueOfResponse := reflect.ValueOf(response)
	for i := 0; i < typeOfResponse.NumField(); i++ {
		fileName := typeOfResponse.Field(i).Name

		//fieldType := valueOfResponse.Elem().Field(i).Type()

		//输出的 struct 成员对象
		structFiled := typeOfResponse.Field(i)
		//从 struct 成员对象 的tag 中的 json 中读取 key信息
		structFiledTagName := structFiled.Tag.Get("json")
		//json里直接读取的字符串还不能用，得转换成http header格式，X-Abc-Def 格式
		headerKey := StrCovertHttpHeader(structFiledTagName)
		////outMap[headerKey] = "a"
		//

		valueFiledType := valueOfResponse.FieldByName(fileName).Type()

		//MyPrint(fileName, valueFiledType.String(), valueOfResponse.FieldByName(fileName).Interface())

		outMap[headerKey] = ""
		value := valueOfResponse.FieldByName(fileName).Interface()
		if valueFiledType.String() == "int" {
			outMap[headerKey] = strconv.Itoa(value.(int))
		} else if valueFiledType.String() == "string" {
			outMap[headerKey] = value.(string)
		} else {
			MyPrint("valueFiledType err:", valueFiledType)
			continue
		}
	}
	return outMap
}

func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func SHA1_1(s string) string {
	t := sha1.New()

	io.WriteString(t, s)
	sign := fmt.Sprintf("%x", t.Sum(nil))
	return sign
}

func SHA1_2(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	fmt.Printf("%x\n", bs)
	return string(bs)
}
