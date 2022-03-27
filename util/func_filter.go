package util

import (
	"regexp"
	"strconv"
)

func CheckMobileRule(mobile string) bool {
	pattern := "^1[3|4|5|6|7|8|9][0-9]\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(mobile)
}
func CheckEmailRule(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func CheckIp4Rule(ip string) bool {
	pattern := `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(ip)
}

func CheckNameRule(ip string) bool { //帐号是否合法(字母开头，允许7-50字节，允许字母数字下划线)
	pattern := `^[a-zA-Z][a-zA-Z0-9_]{7,50}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(ip)
}

//unix stamp 10位，秒为单位
func CheckUnixStampSecondRule(time int) bool {
	//pattern := `^\d{10}$`
	pattern := `^1\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(strconv.Itoa(time))
}

//unix stamp 13位，毫秒为单位
func CheckUnixStampMicroSecondRule(time int) bool {
	pattern := `^1\d{12}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(strconv.Itoa(time))
}

//身份证号
func CheckIdNumberRule(idNumber string) bool {
	pattern := `^\d{17}[0-9Xx]|\d{15}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(idNumber)

}
