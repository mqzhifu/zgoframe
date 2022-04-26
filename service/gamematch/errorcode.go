package gamematch

type MyErrorCode struct {
	Code 	int
	Msg  	string
	Flag 	string
	MsgCn 	string
}
func getErrorCode() (container []string) {
	//redis相关
	container = append(container, "300,redis connect err,error,redis连接超时")
	container = append(container, "301,redisDo has err {0},error,执行redis指令失败")
	//log类相关
	container = append(container, "350,init log class has err : {0},error,初始化公共库：日志，失败")

	//报名相关
	container = append(container, "400,ruleId not in redis db,error,ruleId错误不在REDIS-DB中")
	container = append(container, "401,players len is : 0,error,报名的玩家数为空")
	container = append(container, "402,sign err Status:PlayerStatusSign  id:{0} ,error,玩家状态：报名中~不能重复报名")
	container = append(container, "403,sign err Status:PlayerStatusSuccess  id:{0} ,error,玩家状态为已成功匹配，等待ROOM服务接收成功数据")
	container = append(container, "405,sign err Status:PlayerStatusSign id:{0} but timeout and group person > 1 ,error,玩家状态为报名中，且已失效，等待后台任务回收，并且该玩家所有组的人数大于1")
	container = append(container, "406,players is timeout : 0,error,报名的玩家数为空")
	container = append(container, "407,sign err Status:PlayerStatusSign timeout  id:{0} ,error,玩家状态：~不能重复报名，已超时,等待后台协程回收或者可以取消报名")
	container = append(container, "408,sign err :player len  > rule.max person id:{0},error,rule规定是N人-报名的小组人数大于此数")
	container = append(container, "409,sign err :check groupId uniq...:{0},error,外部groupId错误")
	container = append(container, "410,sign err :player len  > rule.PersonCondition ,max person :{0},error,rule规定是N人-报名的小组人数大于此数")
	container = append(container, "411,sign err :max person > 5,报名每组人数最大为5")

	container = append(container, "450,matchCode is null,error,matchCode为空")
	container = append(container, "451,matchCode not exist in db,error,matchCode在DB中找不见")
	container = append(container, "452,groupId is null or <= 0,error,groupId为空")
	container = append(container, "453,customProp is null,error,customProp为空")
	container = append(container, "454,playerList is null,error,playerList为空")
	container = append(container, "455, playerList.(map[string]map[string]interface{}) error ,error,playerList格式错误")
	container = append(container, "456,some player id is null,error,某个玩家ID为空")
	container = append(container, "457, player id is null or group id is null,error,某个玩家ID为空")
	container = append(container, "458,oldPlayerStatusElement == newPlayerStatusElement,error,更新玩家信息出错，新值与旧值相同")
	container = append(container, "459,httpReqBusiness decode json ,error,un")
	container = append(container, "460,httpReqBusiness successId = 0 ,error,un")
	//gamematch
	container = append(container, "500,getErrorCode list is empty ,error,错误描述文件列表为空")
	container = append(container, "501,pidFilePath is empty,error,进程ID保存文件位置为空")
	container = append(container, "502,init initPid open file err : {0},error,进程ID保存文件-打开失败")
	container = append(container, "503,pidFilePath io.WriteString err :{0},error,进程ID写入出错")

	//container = append(container, "555, ,error,公众错误类中，有一个找不到错误码的情况")

	//rule
	container = append(container, "600,rule config HGETALL err:{0},error,rule数据库中为空")
	container = append(container, "601,RuleConfig is null  ,error,rule数据库中为空")
	container = append(container, "602,rule id not exist :{0},error,ruleId为空")
	container = append(container, "603,addRuleOne redis err,error,添加一条rule发生错误")

	container = append(container, "604,id <= 0,error,id为空")
	container = append(container, "605,AppId <= 0,error,appId为空")
	container = append(container, "606,Flag <= 0,error,flag <= 0")
	container = append(container, "607,Flag must : 1 or 2,error,flag值错误~只允许有：1 和2 ")
	container = append(container, "608,TeamVSPerson <= 0,error,TeamVSPerson<=0")
	container = append(container, "609,TeamVSPerson > {0},error,TeamVSPerson")
	container = append(container, "610,PersonCondition <= 0,error,PersonCondition")
	container = append(container, "611,PersonCondition > {0},error,PersonCondition")
	container = append(container, "612,MatchTimeout < {0} or MatchTimeout > {1},error,MatchTimeout")
	container = append(container, "613,SuccessTimeout < {0} or SuccessTimeout > {1},error,MatchTimeout")
	container = append(container, "614,GroupPersonMax <= 0,error,GroupPersonMax")
	container = append(container, "615,GroupPersonMax > {0},error,GroupPersonMax")
	container = append(container, "616,CategoryKey is empty,error,CategoryKey为空")
	container = append(container, "617,PlayerWeight.ScoreMin > rule.PlayerWeight.ScoreMax,error,ScoreMin>ScoreMax")


	container = append(container, "620, etcd rule matchCode is empty ,error,matchCode为空 ")
	container = append(container, "621, etcd rule value is empty,error,rule结构为空 ")
	container = append(container, "622, etcd rule json.Unmarshal err {0},error, 解析rule字符串到json失败 ")
	container = append(container, "623, playerStatus not equal Sign {0},error,un ")
	container = append(container, "624, etcd rule status != online ,error,un ")
	container = append(container, "625, rule config is null ,error,从etcd里可以读到ruleList ，但是用户监听的ruleId list 对不上 ")
	//push相关
	container = append(container, "700, push respone code err {0},error,un ")
	//group相关
	container = append(container, "750, groupId not in db {0},error,un ")

	//http 相关
	container = append(container, "800,http content = 0 ,post data is empty ,error,该接口需要POST数据，但数据为空")
	container = append(container, "801,http no route this uri ,error,请求URI无法路由到具体方法")
	container = append(container, "802,post data is empty,error,数据为空")
	container = append(container, "803,HttpdRuleState error,error,该rule未初始化")
	container = append(container, "804,state != HTTPD_RULE_STATE_OK,error,该rule的httpd接口，未开启")
	container = append(container, "805, json decode to struct failed : {0},error,un ")
	container = append(container, "806, MatchCode not in map ,error,unknow ")
	container = append(container, "807, successId not in db ,error,unknow ")

	//系统级别错误
	container = append(container, "900,get ErrorCode container is null ,error,获取错误码配置表为空")
	//etcd相关
	container = append(container, "910,init etcd error {0} ,error,初始化etcd 发生错误")
	container = append(container, "911,service HttpPost etcd  {0} ,error,向微服务发送请求时，出现错误")


	return container
}