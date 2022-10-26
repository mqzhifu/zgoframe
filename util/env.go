package util

import "strconv"

func CheckEnvExist(env int) bool {
	list := GetConstListEnv()
	for _, v := range list {
		if v == env {
			return true
		}
	}
	return false
}

func GetConstListEnv() map[string]int {
	list := make(map[string]int)
	list["本地"] = ENV_LOCAL_INT
	list["开发"] = ENV_DEV_INT
	list["测试"] = ENV_TEST_INT
	list["预发布"] = ENV_PRE_INT
	list["线上"] = ENV_ONLINE_INT

	return list
}

func GetConstListEnvStr() map[int]string {
	list := make(map[int]string)
	list[ENV_LOCAL_INT] = ENV_LOCAL_STR
	list[ENV_DEV_INT] = ENV_DEV_STR
	list[ENV_TEST_INT] = ENV_TEST_STR
	list[ENV_PRE_INT] = ENV_PRE_STR
	list[ENV_ONLINE_INT] = ENV_ONLINE_STR

	return list
}

func ConstListEnvToStr() string {
	list := GetConstListEnv()
	str := ""
	for k, v := range list {
		str += strconv.Itoa(v) + k
	}
	return str
}
