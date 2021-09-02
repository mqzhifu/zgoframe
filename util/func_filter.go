package util

import "regexp"

func CheckMobileRule(mobile string)bool{
	pattern := "^1[3|4|5|6|7|8|9][0-9]\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(mobile)
}
func CheckEmailRule(email string)bool{
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}


