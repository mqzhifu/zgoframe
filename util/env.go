package util

const(
	ENV_LOCAL					= "local"//开发环境
	ENV_DEV					= "dev"//开发环境
	ENV_TEST				= "test"//测试环境
	ENV_PRE					= "pre"//预发布环境
	ENV_ONLINE				= "online"//线上环境
)

func GetEnvList()[]string{
	list := []string{ENV_LOCAL, ENV_DEV, ENV_TEST, ENV_PRE, ENV_ONLINE}
	return list
}
func CheckEnvExist(env string)bool{
	list := []string{ENV_LOCAL, ENV_DEV, ENV_TEST, ENV_PRE, ENV_ONLINE}
	for _,v :=range list{
		if v == env{
			return true
		}
	}
	return false
}

