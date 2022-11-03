package gamematch

import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
	"zgoframe/service"
	"zgoframe/util"
)

type Result struct {
	//MatchCode   string
	Id          int
	RuleId      int
	ATime       int   //匹配成功的时间
	Timeout     int   //多少秒后无人来取，即超时，更新用户状态，删除数据
	Teams       []int //该结果，有几个 队伍，因为每个队伍是一个集合，要用来索引
	PlayerIds   []int
	GroupIds    []int
	PushId      int
	Groups      []Group     //该结果下包含的两个组详细信息，属性挂载，用于push payload
	PushElement PushElement //该结果下推送的详细信息，属性挂载
}

type QueueSuccess struct {
	Mutex                  sync.Mutex //锁
	Rule                   *Rule      //父类
	RedisKeySeparator      string
	RedisTextSeparator     string
	RedisIdSeparator       string
	RedisPayloadSeparation string
	Log                    *zap.Logger //log 实例
	Redis                  *util.MyRedisGo
	Err                    *util.ErrMsg
	CloseChan              chan int
	prefix                 string
}

func NewQueueSuccess(rule *Rule) *QueueSuccess {
	queueSuccess := new(QueueSuccess)
	queueSuccess.Rule = rule
	queueSuccess.Redis = rule.RuleManager.Option.GameMatch.Option.Redis
	queueSuccess.RedisTextSeparator = rule.RuleManager.Option.GameMatch.Option.RedisTextSeparator
	queueSuccess.RedisKeySeparator = rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
	queueSuccess.Log = rule.RuleManager.Option.GameMatch.Option.Log
	queueSuccess.Err = rule.RuleManager.Option.GameMatch.Err
	queueSuccess.CloseChan = make(chan int)
	queueSuccess.prefix = "success"
	queueSuccess.RedisIdSeparator = rule.RuleManager.Option.GameMatch.Option.RedisIdSeparator
	queueSuccess.RedisPayloadSeparation = rule.RuleManager.Option.GameMatch.Option.RedisPayloadSeparation
	return queueSuccess
}
func (queueSuccess *QueueSuccess) NewResult() Result {
	result := Result{
		Id:        queueSuccess.GetResultIncId(),
		RuleId:    queueSuccess.Rule.Id,
		ATime:     util.GetNowTimeSecondToInt(),
		Timeout:   util.GetNowTimeSecondToInt() + queueSuccess.Rule.SuccessTimeout,
		Teams:     nil,
		GroupIds:  nil,
		PlayerIds: nil,
		PushId:    0,
		//MatchCode: queueSuccess.Rule.CategoryKey,
		//Groups
		//PushElement
	}
	return result
}

func (queueSuccess *QueueSuccess) getRedisKeyResultPrefix() string {
	return queueSuccess.Rule.GetCommRedisKeyByModuleRuleId(queueSuccess.prefix, queueSuccess.Rule.Id) + "result"
}

func (queueSuccess *QueueSuccess) getRedisKeyResult(successId int) string {
	return queueSuccess.getRedisKeyResultPrefix() + queueSuccess.RedisKeySeparator + strconv.Itoa(successId)
}

func (queueSuccess *QueueSuccess) getRedisKeyTimeout() string {
	return queueSuccess.Rule.GetCommRedisKeyByModuleRuleId(queueSuccess.prefix, queueSuccess.Rule.Id) + "timeout"
}

func (queueSuccess *QueueSuccess) getRedisResultIncKey() string {
	return queueSuccess.Rule.GetCommRedisKeyByModuleRuleId(queueSuccess.prefix, queueSuccess.Rule.Id) + "inc_id"
}

//最简单的string：一个组的详细信息
func (queueSuccess *QueueSuccess) getRedisKeyGroupPrefix() string {
	return queueSuccess.Rule.GetCommRedisKeyByModuleRuleId(queueSuccess.prefix, queueSuccess.Rule.Id) + "group"
}
func (queueSuccess *QueueSuccess) getRedisKeyGroup(groupId int) string {
	return queueSuccess.getRedisKeyGroupPrefix() + queueSuccess.RedisKeySeparator + strconv.Itoa(groupId)
}

//=================上面均是操作redis key==============

func (queueSuccess *QueueSuccess) GetResultById(id int, isIncludeGroupInfo int, isIncludePushInfo int) (result Result, empty int) {
	key := queueSuccess.getRedisKeyResult(id)
	res, _ := redis.String(queueSuccess.Redis.RedisDo("get", key))
	if res == "" {
		queueSuccess.Log.Error("GetResultById is empty~~~")
		return result, 1
	}

	result = queueSuccess.strToStruct(res)
	//fmt.Printf("%+v",result)
	if isIncludeGroupInfo == 1 {
		var groups []Group
		for _, v := range result.GroupIds {
			group := queueSuccess.getGroupElementById(v)
			groups = append(groups, group)
		}
		result.Groups = groups
	}

	if isIncludePushInfo == 1 {
		//push := queueSuccess.Rule.pu.getContainerPushByRuleId(result.RuleId)
		//result.PushElement = push.getById(result.PushId)
		result.PushElement = queueSuccess.Rule.Push.getById(result.RuleId)
	}
	//fmt.Printf("%+v",result)
	return result, 0
}
func (queueSuccess *QueueSuccess) getGroupElementById(id int) (group Group) {
	key := queueSuccess.getRedisKeyGroup(id)
	res, _ := redis.String(queueSuccess.Redis.RedisDo("get", key))
	if res == "" {
		util.MyPrint(" getGroupElementById is empty!")
		return group
	}
	group = queueSuccess.Rule.RuleManager.Option.GameMatch.GroupStrToStruct(res)
	return group
}

//获取并生成一个自增GROUP-ID
func (queueSuccess *QueueSuccess) GetResultIncId() int {
	key := queueSuccess.getRedisResultIncKey()
	res, _ := redis.Int(queueSuccess.Redis.RedisDo("INCR", key))
	return res
}

//添加一条匹配成功记录
func (queueSuccess *QueueSuccess) addOne(redisConn redis.Conn, result Result) PushElement {
	//mymetrics.IncNode("matchingSuccess")

	queueSuccess.Log.Info("func : addOne")
	//添加元素超时信息
	key := queueSuccess.getRedisKeyTimeout()
	res, err := queueSuccess.Redis.Send(redisConn, "zadd", redis.Args{}.Add(key).Add(result.Timeout).Add(result.Id)...)
	//res,err := queueSuccess.Redis.RedisDo("zadd",redis.Args{}.Add(key).Add(result.Timeout).Add(result.Id)...)
	util.MyPrint("add timeout index rs : ", res, err)
	//这里注意下：pushId = 0
	resultStr := queueSuccess.structToStr(result)
	payload := strings.Replace(resultStr, queueSuccess.RedisTextSeparator, queueSuccess.RedisPayloadSeparation, -1)
	pushElement := queueSuccess.Rule.Push.addOnePush(redisConn, result.Id, service.PushCategorySuccess, payload)
	result.PushId = pushElement.Id
	queueSuccess.Log.Info("addOnePush , newId : " + strconv.Itoa(pushElement.Id))
	//添加一条元素
	key = queueSuccess.getRedisKeyResult(result.Id)
	//这里还得重新再  to str 一下，因为pushid 已经可以拿到了
	str := queueSuccess.structToStr(result)
	res, err = queueSuccess.Redis.Send(redisConn, "set", redis.Args{}.Add(key).Add(str)...)
	util.MyPrint("add successResult rs : ", res, err)

	return pushElement
	//mymetrics.FastLog("MatchSuccess",zlib.METRICS_OPT_INC,0)

}

//一条匹配成功记录，要包括N条组信息，这是添加一个组的记录
func (queueSuccess *QueueSuccess) addOneGroup(redisConn redis.Conn, group Group) {
	key := queueSuccess.getRedisKeyGroup(group.Id)
	content := queueSuccess.Rule.RuleManager.Option.GameMatch.GroupStructToStr(group)
	res, err := queueSuccess.Redis.Send(redisConn, "set", redis.Args{}.Add(key).Add(content)...)
	//res,err := queueSuccess.Redis.RedisDo("set",redis.Args{}.Add(key).Add(content)...)
	util.MyPrint("addOneGroup  success ", res, err)
}

func (queueSuccess *QueueSuccess) delOneGroup(redisConn redis.Conn, groupId int) {
	key := queueSuccess.getRedisKeyGroup(groupId)
	res, err := queueSuccess.Redis.Send(redisConn, "del", redis.Args{}.Add(key).Add(key)...)
	//res,err := queueSuccess.Redis.RedisDo("del",redis.Args{}.Add(key).Add(key)...)
	util.MyPrint("success delOneGroup : ", res, err)
	util.MyPrint("delOneGroup", groupId, res, err)
}

func (queueSuccess *QueueSuccess) strToStruct(redisStr string) Result {
	strArr := strings.Split(redisStr, queueSuccess.RedisTextSeparator)
	//fmt.Printf("%+v",strArr)
	Teams := strings.Split(strArr[4], queueSuccess.RedisIdSeparator)
	PlayerIds := strings.Split(strArr[5], queueSuccess.RedisIdSeparator)
	GroupIds := strings.Split(strArr[6], queueSuccess.RedisIdSeparator)
	result := Result{
		Id:     util.Atoi(strArr[0]),
		RuleId: util.Atoi(strArr[1]),
		//MatchCode: strArr[2],
		ATime:     util.Atoi(strArr[2]),
		Timeout:   util.Atoi(strArr[3]),
		Teams:     util.ArrStringCoverArrInt(Teams),
		PlayerIds: util.ArrStringCoverArrInt(PlayerIds),
		GroupIds:  util.ArrStringCoverArrInt(GroupIds),
		PushId:    util.Atoi(strArr[7]),
	}
	//fmt.Printf("%+v",result)
	return result
}

func (queueSuccess *QueueSuccess) structToStr(result Result) string {
	//fmt.Printf("%+v",result)
	//PushId		int
	//Groups		[]Group		//该结果下包含的两个组详细信息，属性挂载，用于push payload
	//PushElement	PushElement	//该结果下推送的详细信息，属性挂载

	TeamsStr := util.ArrCoverStr(result.Teams, queueSuccess.RedisIdSeparator)
	PlayerIds := util.ArrCoverStr(result.PlayerIds, queueSuccess.RedisIdSeparator)
	GroupIds := util.ArrCoverStr(result.GroupIds, queueSuccess.RedisIdSeparator)

	content :=
		strconv.Itoa(result.Id) + queueSuccess.RedisTextSeparator +
			strconv.Itoa(result.RuleId) + queueSuccess.RedisTextSeparator +
			//result.MatchCode + service.Separation +
			strconv.Itoa(result.ATime) + queueSuccess.RedisTextSeparator +
			strconv.Itoa(result.Timeout) + queueSuccess.RedisTextSeparator +
			TeamsStr + queueSuccess.RedisTextSeparator +
			PlayerIds + queueSuccess.RedisTextSeparator +
			GroupIds + queueSuccess.RedisTextSeparator +
			strconv.Itoa(result.PushId) + queueSuccess.RedisTextSeparator
	//Groups		[]Group			这两个是挂载的，先不管
	//PushElement	PushElement		这两个是挂载的，先不管
	return content
}

//删除所有：池里的报名组、玩家、索引等-有点暴力，尽量不用
//func  (queueSuccess *QueueSuccess)   delAll(){
//	key := queueSuccess.getRedisPrefixKey()
//	queueSuccess.Redis.RedisDo("del",key)
//
//	queueSuccess.Push.delAll()
//}
//
func (queueSuccess *QueueSuccess) delOneRule() {
	queueSuccess.Log.Info(" queueSuccess delOneRule ")
	keys := queueSuccess.Rule.GetCommRedisKeyByModuleRuleId("success", queueSuccess.Rule.Id) + "*"
	queueSuccess.Redis.RedisDelAllByPrefix(keys)
	//queueSuccess.delALLResult()
	//queueSuccess.delALLTimeout()
	//queueSuccess.delALLGroup()
}

//====================================================
func (queueSuccess *QueueSuccess) delOneResult(redisConn redis.Conn, id int, isIncludeGroupInfo int, isIncludePushInfo int, isIncludeTimeout int, isIncludePlayerStatus int) {
	util.MyPrint("delOneResult id :", id, isIncludeGroupInfo, isIncludePushInfo, isIncludeTimeout)
	element, isEmpty := queueSuccess.GetResultById(id, isIncludeGroupInfo, isIncludePushInfo)
	if isEmpty == 1 {
		queueSuccess.Log.Error("del failed ,id is empty~")
		queueSuccess.Log.Error("del failed ,id is empty~")
		return
	}
	key := queueSuccess.getRedisKeyResult(id)
	res, err := queueSuccess.Redis.Send(redisConn, "del", redis.Args{}.Add(key)...)
	//res,err :=   queueSuccess.Redis.RedisDo("del",redis.Args{}.Add(key)... )

	util.MyPrint(" delOneRuleOneResult res", res, err)
	util.MyPrint(" delOneRuleOneResult res", res, err)

	if isIncludePushInfo == 1 {
		//push := queueSuccess.Gamematch.getContainerPushByRuleId(element.RuleId)
		util.MyPrint("delOnePush", element.PushId)
		queueSuccess.Rule.Push.delOnePush(redisConn, element.PushId)
	}

	if isIncludeTimeout == 1 {
		queueSuccess.Log.Info("delOneTimeout" + strconv.Itoa(id))
		queueSuccess.delOneTimeout(redisConn, id)
	}

	if isIncludeGroupInfo == 1 {
		for _, groupId := range element.GroupIds {
			queueSuccess.delOneGroup(redisConn, groupId)
		}
	}

	if isIncludeTimeout == 1 {
		for _, playerId := range element.PlayerIds {
			queueSuccess.Log.Info("playerStatus.delOneById " + strconv.Itoa(playerId))
			queueSuccess.Rule.PlayerManager.delOneById(redisConn, playerId)
		}
	}

}

func (queueSuccess *QueueSuccess) delOneTimeout(redisConn redis.Conn, id int) {
	key := queueSuccess.getRedisKeyTimeout()
	res, _ := queueSuccess.Redis.Send(redisConn, "ZREM", redis.Args{}.Add(key).Add(id)...)
	//res,_ :=  queueSuccess.Redis.RedisDo("ZREM",redis.Args{}.Add(key).Add(id)... )
	util.MyPrint(" success delOneTimeout res", res)
}
func (queueSuccess *QueueSuccess) GetAllTimeoutCnt() int {
	key := queueSuccess.getRedisKeyTimeout()
	res, _ := redis.Int(queueSuccess.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add("-inf").Add("+inf")...))
	return res
}

func (queueSuccess *QueueSuccess) Demon() {
	queueSuccess.Log.Info(queueSuccess.prefix + " Demon")
	for {
		select {
		case signal := <-queueSuccess.CloseChan:
			queueSuccess.Log.Warn(queueSuccess.prefix + "Demon CloseChan receive :" + strconv.Itoa(signal))
			goto forEnd
		default:
			queueSuccess.CheckTimeout()
			time.Sleep(time.Millisecond * time.Duration(queueSuccess.Rule.RuleManager.Option.GameMatch.LoopSleepTime))
		}
	}
forEnd:
	//demonLog.Notice("MyDemon end : ",title)
	queueSuccess.Log.Warn(queueSuccess.prefix + " Demon end .")
}

func (queueSuccess *QueueSuccess) CheckTimeout() {
	keys := queueSuccess.getRedisKeyTimeout()
	push := queueSuccess.Rule.Push
	redisConnFD := queueSuccess.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()

	now := util.GetNowTimeSecondToInt()
	res, err := redis.IntMap(queueSuccess.Redis.RedisDo("ZREVRANGEBYSCORE", redis.Args{}.Add(keys).Add(now).Add("-inf").Add("WITHSCORES")...))
	if err != nil {
		queueSuccess.Log.Error("redis keys err :" + err.Error())
		return
	}
	if len(res) == 0 {
		if now%queueSuccess.Rule.DemonDebugTime == 0 { //每10秒 输出一次，避免日志过多
			queueSuccess.Log.Info(queueSuccess.prefix + " timeout empty , no need process")
		}
		return
	}
	queueSuccess.Log.Info("queueSuccess timeout group element total : " + strconv.Itoa(len(res)))
	for resultId, _ := range res {
		//queueSuccess.Redis.Send(redisConnFD,"multi")
		queueSuccess.Redis.Multi(redisConnFD)

		resultIdInt := util.Atoi(resultId)
		element, _ := queueSuccess.GetResultById(resultIdInt, 0, 0)
		//util.MyPrint("GetResultById", resultIdInt, element)
		//fmt.Printf("%+v",element)
		queueSuccess.delOneResult(redisConnFD, resultIdInt, 1, 1, 1, 1)

		payload := queueSuccess.structToStr(element)
		payload = strings.Replace(payload, queueSuccess.RedisTextSeparator, queueSuccess.RedisPayloadSeparation, -1)
		push.addOnePush(redisConnFD, resultIdInt, service.PushCategorySuccessTimeout, payload)

		queueSuccess.Redis.Exec(redisConnFD)
		//queueSuccess.Redis.ConnDo(redisConnFD,"exec")
	}
	//myGosched("success CheckTimeout")
	//	mySleepSecond(1," success CheckTimeout ")
	//}
}

func (queueSuccess *QueueSuccess) Close() {
	queueSuccess.Log.Info(queueSuccess.prefix + " Close")
	queueSuccess.CloseChan <- 1
}

func (queueSuccess *QueueSuccess) TestRedisKey() {
	redisKey := queueSuccess.getRedisKeyResultPrefix()
	util.MyPrint("queueSuccess test :", redisKey)

	redisKey = queueSuccess.getRedisKeyResult(1)
	util.MyPrint("queueSuccess test :", redisKey)

	redisKey = queueSuccess.getRedisKeyTimeout()
	util.MyPrint("queueSuccess test :", redisKey)

	redisKey = queueSuccess.getRedisKeyGroupPrefix()
	util.MyPrint("queueSuccess test :", redisKey)

	redisKey = queueSuccess.getRedisKeyGroup(1)
	util.MyPrint("queueSuccess test :", redisKey)
}

//删除一条规则的，所有分组详细信息
//func (queueSuccess *QueueSuccess)  delALLTimeout( ){
//	key := queueSuccess.getRedisKeyTimeout()
//	res,_ := queueSuccess.Redis.RedisDo("del",key)
//	mylog.Debug(" delALLTimeout res",res)
//}
//删除一条规则的，所有分组详细信息
//func (queueSuccess *QueueSuccess)  delALLResult( ){
//	prefix := queueSuccess.getRedisKeyResultPrefix()
//	res,_ := redis.Strings( queueSuccess.Redis.RedisDo("keys",prefix + "*"  ))
//	if len(res) == 0{
//		mylog.Notice(" delALLResult by keys(*) : is empty")
//		return
//	}
//	for _,v := range res{
//		res,_ := redis.Int(queueSuccess.Redis.RedisDo("del",v))
//	}
//}

////删除一条规则的，所有分组详细信息
//func (queueSuccess *QueueSuccess)  delALLGroup( ){
//	prefix := queueSuccess.getRedisKeyGroupPrefix()
//	res,_ := redis.Strings( queueSuccess.Redis.RedisDo("keys",prefix + "*"  ))
//	if len(res) == 0{
//		mylog.Notice(" delALLGroup by keys(*) : is empty")
//		return
//	}
//	for _,v := range res{
//		res,_ := redis.Int(queueSuccess.Redis.RedisDo("del",v))
//	}
//
//	queueSuccess.Push.delOneRule()
//}
