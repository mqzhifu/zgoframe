package util

import "strconv"

const (
	ENV_LOCAL_STR  = "local"  //开发环境
	ENV_DEV_STR    = "dev"    //开发环境
	ENV_TEST_STR   = "test"   //测试环境
	ENV_PRE_STR    = "pre"    //预发布环境
	ENV_ONLINE_STR = "online" //线上环境

	ENV_LOCAL_INT  = 1  //开发环境
	ENV_DEV_INT    = 2    //开发环境
	ENV_TEST_INT   = 3   //测试环境
	ENV_PRE_INT    = 4    //预发布环境
	ENV_ONLINE_INT = 5 //线上环境

)

//func GetEnvList() []string {
//	list := []string{ENV_LOCAL, ENV_DEV, ENV_TEST, ENV_PRE, ENV_ONLINE}
//	return list
//}
//func CheckEnvExist(env string) bool {
//	list := []string{ENV_LOCAL, ENV_DEV, ENV_TEST, ENV_PRE, ENV_ONLINE}
//	for _, v := range list {
//		if v == env {
//			return true
//		}
//	}
//	return false
//}


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

func ConstListEnvToStr()string{
	list := GetConstListEnv()
	str := ""
	for k, v := range list {
		str += strconv.Itoa(v) + k
	}
	return str
}

func GetConstListEnvStr() map[int]string {
	list := make(map[int]string)
	list[ENV_LOCAL_INT] = ENV_LOCAL_STR
	list[ENV_DEV_INT] =ENV_DEV_STR
	list[ENV_TEST_INT] = ENV_TEST_STR
	list[ENV_PRE_INT] = ENV_PRE_STR
	list[ENV_ONLINE_INT] = ENV_ONLINE_STR

	return list
}

//func GetConstListEnvString() map[string]string {
//	list := make(map[string]string)
//	list["本地"] = ENV_LOCAL
//	list["开发"] = ENV_DEV
//	list["测试"] = ENV_TEST
//	list["预发布"] = ENV_PRE
//	list["线上"] = ENV_ONLINE
//
//	return list
//}
