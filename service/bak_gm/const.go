package gamematch

//
//import "zgoframe/service"
//
//type CmdArgs struct {
//	Env         string `seq:"1" err:"env=local" desc:"环境变量： local test dev pre online"`
//	LogBasePath string `seq:"2" err:"log_base_path=/golang/logs" desc:"日志文件保存位置"`
//	LogLevel    string `seq:"3" err:"log_level=" desc:"日志级别"`
//	BaseUrl     string `seq:"4" err:"base_url=" desc:"获取配置URL"`
//}
//
//func GetEnvList() []string {
//	list := []string{service.ENV_DEV, service.ENV_TEST, service.ENV_PRE, service.ENV_ONLINE}
//	return list
//}
//func CheckEnvExist(env string) bool {
//	list := []string{service.ENV_DEV, service.ENV_TEST, service.ENV_PRE, service.ENV_ONLINE}
//	for _, v := range list {
//		if v == env {
//			return true
//		}
//	}
//	return false
//}
//
///*
//	匹配类型 - 规则
//	1. N人匹配 ，只要满足N人，即成功
//	2. N人匹配 ，划分为2个队，A队满足N人，B队满足N人，即成功
//
//	权重		：根据某个用户上的某个特定属性值，计算出权重，优先匹配
//	组		：ABC是一个组，一起参与匹配，那这3个人匹配的时候是不能分开的
//	游戏属性	：游戏类型，也可以有子类型，如：不同的赛制。最终其实是分队列。不同的游戏忏悔分类，有不同的分类
//*/
//
//func getSignalDesc(signal int) string {
//	switch signal {
//	case service.SIGNAL_GOROUTINE_EXEC_EXIT:
//		return "请执行协程退出"
//	case service.SIGNAL_GOROUTINE_EXIT_FINISH:
//		return "协程已退出"
//	case service.SIGNAL_QUIT_SOURCE:
//		return "退出来源1"
//	case service.SIGNAL_QUIT_SOURCE_RULE_WATCH:
//		return "退出来源~rule发生变更"
//	default:
//		return "signal错误"
//	}
//}
//
////GameMatchMetrics	所有需要统计的信息
//type GMMData struct {
//	StartUpTime       int //启动时间
//	ShutdownStartTime int //接收到结束信号的时间
//	ShutdownTime      int //关闭时间
//	InitEndTime       int //初始化结束时间
//
//	HttpSign        int //http请求报名数
//	HttpSignSuccess int //http请求报名成功数
//	HttpSignFiled   int //http请求报名失败数
//
//	HttpCancel        int //http请求报名数
//	HttpCancelSuccess int //http请求报名成功数
//	HttpCancelFiled   int //http请求报名失败数
//
//	SignTimeout  int //报名玩家超时
//	MatchSuccess int //匹配成功数
//	Rule         int //rule 总数
//
//	PushMatchSuccessOk   int //推送：匹配成功记录-对方正常接收
//	PushMatchSuccessDrop int //推送：匹配成功记录-对方接收失败
//
//	PushMatchSuccessTimeoutOk   int //推送：匹配成功记录，但对方一直拒绝接收，PUSH也没有超过重度次数，记录本身超时了
//	PushMatchSuccessTimeoutDrop int //
//
//	PushSignOk   int //推送：报名超时记录-对方正常接收
//	PushSignDrop int //推送：报名超时记录-对方接收失败
//
//}
