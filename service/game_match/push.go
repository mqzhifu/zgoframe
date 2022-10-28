package gamematch

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
	"zgoframe/http/request"
	"zgoframe/service"
	"zgoframe/util"
)

//推送给3方，支持重试
type PushElement struct {
	Id       int
	ATime    int
	UTime    int //最后更新的时间
	LinkId   int
	Status   int    //状态：1未推送2推送失败，等待重试3推送成功4推送失败，不再重试
	Times    int    //已推送次数
	Category int    //1:报名超时 2匹配成功 3成功结果超时
	Payload  string //自定义的载体
}

type Push struct {
	Mutex              sync.Mutex
	Rule               *Rule //父类
	RedisKeySeparator  string
	RedisTextSeparator string
	Log                *zap.Logger //log 实例
	Redis              *util.MyRedisGo
	Err                *util.ErrMsg
	CloseChan          chan int
	prefix             string
	RetryPeriod        []int
}

func NewPush(rule *Rule) *Push {
	push := new(Push)

	push.Rule = rule
	push.Redis = rule.RuleManager.Option.GameMatch.Option.Redis
	push.RedisTextSeparator = rule.RuleManager.Option.GameMatch.Option.RedisTextSeparator
	push.RedisKeySeparator = rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
	push.Log = rule.RuleManager.Option.GameMatch.Option.Log
	push.Err = rule.RuleManager.Option.GameMatch.Err
	push.CloseChan = make(chan int)
	push.prefix = "push"
	//var PushRetryPeriod = []int{10,30,60,600}
	push.RetryPeriod = []int{5, 10, 15} //方便测试
	return push
}

func (push *Push) TestRedisKey() {
	redisKey := push.getRedisPushIncKey()
	util.MyPrint("push test :", redisKey)

	redisKey = push.getRedisKeyPushStatus()
	util.MyPrint("push test :", redisKey)

	redisKey = push.getRedisKeyPush(1)
	util.MyPrint("push test :", redisKey)

}

func (push *Push) getRedisPushIncKey() string {
	return push.Rule.GetCommRedisKeyByModuleRuleId(push.prefix, push.Rule.Id) + "inc_id"
}

func (push *Push) getRedisKeyPushStatus() string {
	return push.Rule.GetCommRedisKeyByModuleRuleId(push.prefix, push.Rule.Id) + "status"
}

func (push *Push) getRedisKeyPush(id int) string {
	return push.Rule.GetCommRedisKeyByModuleRuleId(push.prefix, push.Rule.Id) + strconv.Itoa(id)
}

func (push *Push) GetPushIncId() int {

	key := push.getRedisPushIncKey()
	res, _ := redis.Int(push.Redis.RedisDo("INCR", key))
	return res
}

//
func (push *Push) RuntimeSuccess(httpReqGameMatchPlayerSign request.HttpReqGameMatchPlayerSign, ruleId int) {
	//mylog.Debug("RuntimeSuccess : " , httpReqBusiness ,strconv.Itoa( util.GetNowTimeSecondToInt()))
	//time.Sleep(time.Second * 1)
	//mylog.Debug("RuntimeSuccess sleep wake up , " + strconv.Itoa(util.GetNowTimeSecondToInt()))

	pushElement := PushElement{
		//Id  		:0,
		//ATime 		:zlib.GetNowTimeSecondToInt(),
		//UTime   	:zlib.GetNowTimeSecondToInt(),
		//LinkId		:0,
		Status: 1,
		//Times  		:1,
		Category: service.PushCategorySuccess,
		//Payload 	:"",
	}
	var playerIds []int
	for _, v := range httpReqGameMatchPlayerSign.PlayerList {
		playerIds = append(playerIds, v.Uid)
	}
	now := util.GetNowTimeSecondToInt()
	newGroups := Group{
		//MatchCode:      httpReqBusiness.MatchCode,
		//CustomProp:     httpReqBusiness.CustomProp,
		OutGroupId:     httpReqGameMatchPlayerSign.GroupId,
		MatchTimes:     -1,
		Weight:         -99,
		TeamId:         1,
		Addition:       httpReqGameMatchPlayerSign.Addition,
		Person:         len(httpReqGameMatchPlayerSign.PlayerList),
		SignTime:       now,
		SignTimeout:    now,
		SuccessTime:    now,
		SuccessTimeout: now,
	}

	result := Result{
		Id:     -11,
		RuleId: ruleId,
		//MatchCode:   httpReqBusiness.MatchCode,
		ATime:       now + 1,
		Timeout:     now + 1,
		Teams:       []int{httpReqGameMatchPlayerSign.GroupId},
		PlayerIds:   playerIds,
		GroupIds:    []int{httpReqGameMatchPlayerSign.GroupId},
		PushId:      -22,
		Groups:      []Group{newGroups},
		PushElement: pushElement,
	}

	//queueSuccess :=  gamematch.getContainerSuccessByRuleId(ruleId)
	//resultStr := queueSuccess.structToStr(result)
	//
	//thirdMethodUri,postData ,err  := push.getServiceUri(pushElement,resultStr)
	//if err != nil{
	//	push.Log.Error("push.getServiceUri err:",err)
	//}
	//httpRs,err := myservice.HttpPost(SERVICE_MSG_SERVER,thirdMethodUri,postData)
	myServiceDiscovery := push.Rule.RuleManager.Option.GameMatch.Option.ServiceDiscovery
	projectId := push.Rule.RuleManager.Option.GameMatch.Option.ProjectId
	myService, _ := myServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(service.SERVICE_MSG_SERVER, "")
	serviceHttp := util.NewServiceHttp(projectId, service.SERVICE_MSG_SERVER, myService.Ip, myService.Port, myService.ServiceId)

	thirdMethodUri := "v1/match/succ"
	//httpRs,err := myservice.HttpPost(SERVICE_MSG_SERVER,thirdMethodUri,result)
	httpRs, err := serviceHttp.Post(thirdMethodUri, result)
	util.MyPrint("RuntimeSuccess finish: ", httpRs, "|", err)
}

func (push *Push) getById(id int) (element PushElement) {
	key := push.getRedisKeyPush(id)
	res, _ := redis.String(push.Redis.RedisDo("get", key))
	if res == "" {
		return element
	}

	element = push.pushStrToStruct(res)
	return element
}

func (push *Push) addOnePush(redisConn redis.Conn, linkId int, category int, payload string) int {
	push.Log.Debug("addOnePush" + strconv.Itoa(linkId) + " " + strconv.Itoa(category) + " " + payload)
	id := push.GetPushIncId()
	key := push.getRedisKeyPush(id)
	pushElement := PushElement{
		Id:       id,
		ATime:    util.GetNowTimeSecondToInt(),
		Status:   1,
		UTime:    util.GetNowTimeSecondToInt(),
		Times:    0,
		LinkId:   linkId,
		Category: category,
		Payload:  payload,
	}
	pushStr := push.pushStructToStr(pushElement)
	res, err := push.Redis.Send(redisConn, "set", redis.Args{}.Add(key).Add(pushStr)...)
	//res,err := push.Redis.RedisDo("set",redis.Args{}.Add(key).Add(pushStr)...)
	util.MyPrint("addOnePush rs : ", res, err)
	push.Log.Info("addOnePush ,cate : " + strconv.Itoa(category) + payload + ",payload")

	pushKey := push.getRedisKeyPushStatus()
	res, err = push.Redis.Send(redisConn, "zadd", redis.Args{}.Add(pushKey).Add(service.PushStatusWait).Add(id)...)
	//res,err = push.Redis.RedisDo("zadd",redis.Args{}.Add(pushKey).Add(PushStatusWait).Add(id)...)
	util.MyPrint("addOnePush status : ", res, err)
	push.Log.Info("addOnePush status")

	return id
}

func (push *Push) pushStrToStruct(redisStr string) PushElement {
	strArr := strings.Split(redisStr, service.Separation)
	result := PushElement{
		Id:       util.Atoi(strArr[0]),
		LinkId:   util.Atoi(strArr[1]),
		ATime:    util.Atoi(strArr[2]),
		Status:   util.Atoi(strArr[3]),
		UTime:    util.Atoi(strArr[4]),
		Times:    util.Atoi(strArr[5]),
		Category: util.Atoi(strArr[6]),
		Payload:  strArr[7],
	}
	return result
}

func (push *Push) pushStructToStr(pushElement PushElement) string {
	str :=
		strconv.Itoa(pushElement.Id) + service.Separation +
			strconv.Itoa(pushElement.LinkId) + service.Separation +
			strconv.Itoa(pushElement.ATime) + service.Separation +
			strconv.Itoa(pushElement.Status) + service.Separation +
			strconv.Itoa(pushElement.UTime) + service.Separation +
			strconv.Itoa(pushElement.Times) + service.Separation +
			strconv.Itoa(pushElement.Category) + service.Separation +
			pushElement.Payload + service.Separation

	return str
}

//func (push *Push)   delAll(){
//	key := push.getRedisPrefixKey()
//	push.Redis.RedisDo("del",key)
//}

func (push *Push) delOneRule() {
	push.Log.Debug(" push delOneRule : ")
	key := push.Rule.GetCommRedisKeyByModuleRuleId("push", push.Rule.Id) + "*"
	push.Redis.RedisDelAllByPrefix(key)
	//push.delAllPush()
	//push.delAllStatus()
}

//func  (push *Push)  delAllPush( ){
//	prefix := push.getRedisCatePrefixKey()
//	res,_ := redis.Strings( push.Redis.RedisDo("keys",prefix + "*"  ))
//	if len(res) == 0{
//		mylog.Notice(" GroupElement by keys(*) : is empty")
//		return
//	}
//	//zlib.ExitPrint(res,-200)
//	for _,v := range res{
//		res,_ := redis.Int(push.Redis.RedisDo("del",v))
//		zlib.MyPrint("del group element v :",res)
//	}
//}
//
//func  (push *Push)  delAllStatus( ){
//	key := push.getRedisKeyPushStatus()
//	res,_ := redis.Strings( push.Redis.RedisDo("del",key ))
//	mylog.Debug("delAllStatus :",res)
//}

//PushMatchSuccessOk		int	//推送：匹配成功记录-对方正常接收
//PushMatchSuccessDrop	int//推送：匹配成功记录-对方接收失败
//
//PushMatchSuccessTimeoutOk 	int//推送：匹配成功记录，但对方一直拒绝接收，PUSH也没有超过重度次数，记录本身超时了
//PushMatchSuccessTimeoutDrop int//
//
//PushSignOk		int	//推送：报名超时记录-对方正常接收
//PushSignDrop	int	//推送：报名超时记录-对方接收失败

func (push *Push) delOneStatus(redisConn redis.Conn, pushId int) {
	key := push.getRedisKeyPushStatus()
	res, err := push.Redis.Send(redisConn, "ZREM", redis.Args{}.Add(key).Add(pushId)...)
	//res,err :=  push.Redis.RedisDo("ZREM",redis.Args{}.Add(key).Add(pushId)... )
	util.MyPrint(" delOne PushStatus index res", res, err)
	//push.Log.Info(" delOne PushStatus index res",res,err)
}
func (push *Push) metrics(elementCategory int, action string) {
	//actionStr := action
	//if action == "Ok" || action == "Drop"{
	//
	//}else{
	//	mylog.Error("push metrics action err:" + action)
	//	return
	//}
	//
	//key := ""
	//if elementCategory == PushCategorySignTimeout{
	//	key = "PushSign" + actionStr
	//}else if elementCategory == PushCategorySuccess{
	//	key = "PushMatchSuccess" + actionStr
	//}else if elementCategory == PushCategorySuccessTimeout{
	//	key = "PushMatchSuccessTimeout" + actionStr
	//}else{
	//	util.MyPrint("push metrics category err:",elementCategory)
	//	return
	//}

	//mymetrics.FastLog(key,zlib.METRICS_OPT_INC,0)
}

//失败且需要重试的PUSH-ELEMENT
func (push *Push) upRetryPushInfo(element PushElement) {
	redisConnFD := push.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()

	push.Redis.Send(redisConnFD, "multi")
	element.Status = service.PushStatusRetry
	element.UTime = util.GetNowTimeSecondToInt()
	element.Times = element.Times + 1
	key := push.getRedisKeyPush(element.Id)
	pushStr := push.pushStructToStr(element)
	res, err := push.Redis.Send(redisConnFD, "set", redis.Args{}.Add(key).Add(pushStr)...)
	//res,err := push.Redis.RedisDo("set",redis.Args{}.Add(key).Add(pushStr)...)

	util.MyPrint("upRetryPushElementInfo , ", element)
	//这里有个麻烦点，元素信息 和 索引信息，是分开放的，元素的变更比较简单，索引是一个集合，改起来有点麻烦
	//那就直接先删了，再重新添加一条
	statuskey := push.getRedisKeyPushStatus()
	util.MyPrint("del pushStatus index ,pushId : ", element.Id)
	push.delOneStatus(redisConnFD, element.Id)
	res, err = push.Redis.Send(redisConnFD, "zadd", redis.Args{}.Add(statuskey).Add(service.PushStatusRetry).Add(element.Id)...)
	//res,err = push.Redis.RedisDo("zadd",redis.Args{}.Add(statuskey).Add(PushStatusRetry).Add(element.Id)...)

	util.MyPrint("add  new pushStatus index : ", res, err)
	//mylog.Info("add  new pushStatus index : ",res,err)

	push.Redis.Exec(redisConnFD)
}

//在业务里，删除一条push
//走到里，前置条件肯定是PUSH成功了
func (push *Push) delOneByIdInBusiness(redisConn redis.Conn, id int) {
	push.Redis.Send(redisConn, "multi")
	element := push.getById(id)
	push.delOnePush(redisConn, id)
	if element.Category == service.PushCategorySuccess || element.Category == service.PushCategorySuccessTimeout {
		push.Log.Info("delOneResult")
		//success := push.Gamematch.getContainerSuccessByRuleId(push.Rule.Id)
		success := push.Rule.QueueSuccess
		success.delOneResult(redisConn, element.LinkId, 1, 1, 1, 1)
	}
	push.metrics(element.Category, "Ok")
	push.Redis.ConnDo(redisConn, "exec")
}

func (push *Push) Demon() {
	push.Log.Info(push.prefix + " Demon start")
	for {
		select {
		case signal := <-push.CloseChan:
			push.Log.Warn(push.prefix + "Demon CloseChan receive :" + strconv.Itoa(signal))
			goto forEnd
		default:
			push.checkStatus()
			time.Sleep(time.Millisecond * time.Duration(push.Rule.RuleManager.Option.GameMatch.LoopSleepTime))
		}
	}
forEnd:
	push.Log.Warn(push.prefix + "  Demon end .")
}

//检查需要抢着的数据：待推送、重试推送
func (push *Push) checkStatus() {
	//mylog.Info("one rule checkStatus : start ")
	//push.Log.Info("one rule checkStatus : start ")
	key := push.getRedisKeyPushStatus()

	push.checkOneByStatus(key, service.PushStatusWait)
	push.checkOneByStatus(key, service.PushStatusRetry)
	//push.Log.Info("one rule checkStatus : finish ")

}

func (push *Push) getAllCnt() int {
	key := push.getRedisKeyPushStatus()
	res, _ := redis.Int(push.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add("-inf").Add("+inf")...))
	return res
}

func (push *Push) getStatusCnt(status int) int {
	key := push.getRedisKeyPushStatus()
	res, _ := redis.Int(push.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add(status).Add(status)...))
	return res
}

//status:待推送、重试推送
func (push *Push) checkOneByStatus(key string, status int) {
	//mylog.Info("checkOneByStatus :",status)
	res, err := redis.Ints(push.Redis.RedisDo("ZREVRANGEBYSCORE", redis.Args{}.Add(key).Add(status).Add(status)...))
	if err != nil {
		push.Log.Error("redis keys err :" + err.Error())
		push.Log.Error("redis keys err :" + err.Error())
		return
	}
	now := util.GetNowTimeSecondToInt()
	if len(res) == 0 {
		if now%10 == 0 {
			push.Log.Info(push.prefix + " checkOneByStatus :" + strconv.Itoa(status) + " empty , no need process")
		}
		return
	}
	push.Log.Info("push need process element total : " + strconv.Itoa(len(res)))
	for _, id := range res {
		push.processOne(id, status)
	}
}

func (push *Push) processOne(id int, status int) {
	//mylog.Info(" action hook , push id : " ,id ," status : ",status)
	push.Log.Info(" action processOne , push id : " + strconv.Itoa(id) + " status : " + strconv.Itoa(status))
	element := push.getById(id)
	//fmt.Printf("%+v", element)
	if status == service.PushStatusWait {
		push.Log.Info("element first push")
		push.pushAndUpInfo(element, service.PushStatusRetry)
	} else {
		push.Log.Info("element retry ,element.Times:" + strconv.Itoa(element.Times) + " len(PushRetryPeriod):" + strconv.Itoa(len(push.RetryPeriod)))
		if element.Times >= len(push.RetryPeriod) {
			//已超过，最大重试次数
			push.metrics(element.Category, "Drop")

			redisConnFD := push.Redis.GetNewConnFromPool()
			defer redisConnFD.Close()

			push.Log.Warn(" push retry time > maxRetryTime , drop this msg.")
			push.delOneByIdInBusiness(redisConnFD, id)
		} else {
			time := push.RetryPeriod[element.Times]
			util.MyPrint("retry rule : ", push.RetryPeriod, " this time : ", time)
			d := util.GetNowTimeSecondToInt() - element.UTime
			//mylog.Info("this time : ",time,"now :",zlib.GetNowTimeSecondToInt() , " - element.UTime ",element.UTime , " = ",d)
			util.MyPrint("this time : ", time, "now :", util.GetNowTimeSecondToInt(), " - element.UTime ", element.UTime, " = ", d)
			if d >= time {
				push.pushAndUpInfo(element, service.PushStatusRetry)
			} else {
				//重试的时间间隔 未满足
				push.Log.Named("The time is too short to trigger the Push!!! ")
			}
		}
	}
	push.Log.Info("processOne finish")
}
func getAnyType() (a interface{}) {
	return a
}

func (push *Push) getServiceUri(element PushElement, payload string) (uri string, post interface{}, err error) {
	postData := getAnyType()
	thirdMethodUri := ""
	//success := push.Gamematch.getContainerSuccessByRuleId(push.Rule.Id)
	success := push.Rule.QueueSuccess
	if element.Category == service.PushCategorySignTimeout {
		push.Log.Debug("element.Category == PushCategorySignTimeout")
		postData = GroupStrToStruct(payload)

		thirdMethodUri = "v1/match/error"
	} else if element.Category == service.PushCategorySuccessTimeout {
		push.Log.Debug("element.Category == PushCategorySuccessTimeout")
		postData = push.Rule.QueueSuccess.strToStruct(payload)
		thirdMethodUri = "v1/match/error"
	} else if element.Category == service.PushCategorySuccess {
		push.Log.Debug("element.Category == PushCategorySuccess")
		thisResult := success.strToStruct(payload)
		//fmt.Printf("%+v", thisResult)
		postData, _ = success.GetResultById(thisResult.Id, 1, 0)
		thirdMethodUri = "v1/match/succ"
	} else {
		push.Log.Error("element.Category error.")
		return uri, post, errors.New("element.Category error")
	}
	return thirdMethodUri, postData, nil
}

func (push *Push) pushAndUpInfo(element PushElement, upStatus int) {
	//mylog.Debug("pushAndUpInfo",element,upStatus)
	util.MyPrint("pushAndUpInfo", element, " upStatus: ", upStatus)
	var httpRs util.ResponseMsgST
	var err error

	payload := strings.Replace(element.Payload, service.PayloadSeparation, service.Separation, -1)
	thirdMethodUri, postData, err := push.getServiceUri(element, payload)
	if err != nil {
		return
	}
	util.MyPrint("push third service ,  uri : ", thirdMethodUri, " , postData : ", postData)

	myServiceDiscovery := push.Rule.RuleManager.Option.GameMatch.Option.ServiceDiscovery
	projectId := push.Rule.RuleManager.Option.GameMatch.Option.ProjectId
	myService, _ := myServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(service.SERVICE_MSG_SERVER, "")
	serviceHttp := util.NewServiceHttp(projectId, service.SERVICE_MSG_SERVER, myService.Ip, myService.Port, myService.ServiceId)

	//httpRs,err = myservice.HttpPost(SERVICE_MSG_SERVER,thirdMethodUri,postData)
	httpRs, err = serviceHttp.Post(thirdMethodUri, postData)
	util.MyPrint("push third service , httpRs : ", httpRs, " err : ", err)

	if err != nil {
		push.upRetryPushInfo(element)
		msg := push.Err.MakeOneStringReplace(err.Error())
		push.Err.NewReplace(911, msg)
		push.Log.Error("push third service " + err.Error())
		return
	}

	push.hook(element, httpRs)

}
func (push *Push) hook(element PushElement, httpRs util.ResponseMsgST) {
	redisConnFD := push.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()

	if httpRs.Code == 0 { // 0 即是200 ，推送成功~
		//mymetrics.IncNode("pushSuccess")
		push.delOneByIdInBusiness(redisConnFD, element.Id)
		//push.delOnePush(element.Id)
		//if element.Category == PushCategorySuccess || element.Category == PushCategorySuccessTimeout{
		//	push.Log.Info("delOneResult")
		//	push.QueueSuccess.delOneResult(element.LinkId,1,1,1,1)
		//}
		return
	}

	if element.Category == service.PushCategorySignTimeout {
		if httpRs.Code == 116 || httpRs.Code == 119 {
			push.metrics(element.Category, "Ok")
			push.delOnePush(redisConnFD, element.Id)
		} else {
			push.upRetryPushInfo(element)

			httpRsJsonStr, _ := json.Marshal(httpRs)
			msg := push.Err.MakeOneStringReplace(string(httpRsJsonStr))
			//msg := myerr.MakeOneStringReplace(string(httpRsJsonStr))
			push.Err.NewReplace(700, msg)
			return
		}
	} else if element.Category == service.PushCategorySuccessTimeout {
		push.Log.Info("delOneResult")
		push.delOneByIdInBusiness(redisConnFD, element.Id)
		//push.delOnePush(element.Id)
		//push.QueueSuccess.delOneResult(element.LinkId,1,1,1,1)

	} else if element.Category == service.PushCategorySuccess {
		if httpRs.Code == 108 || httpRs.Code == 102 {
			push.upRetryPushInfo(element)

			httpRsJsonStr, _ := json.Marshal(httpRs)
			msg := push.Err.MakeOneStringReplace(string(httpRsJsonStr))
			//msg := myerr.MakeOneStringReplace(string(httpRsJsonStr))
			push.Err.NewReplace(700, msg)
			return
		} else {
			push.Log.Info("delOneResult")
			push.delOneByIdInBusiness(redisConnFD, element.Id)
			//push.delOnePush(element.Id)
			//push.QueueSuccess.delOneResult(element.LinkId,1,1,1,1)
			return
		}
	} else {
		push.Log.Error("pushAndUpInfo element.Category not found!!!")
		push.Log.Error("pushAndUpInfo element.Category not found!!!")
	}
}
func (push *Push) delOnePush(redisConn redis.Conn, id int) {
	key := push.getRedisKeyPush(id)
	push.Log.Info("delOnePush action" + strconv.Itoa(id) + " key:" + key)
	res, err := push.Redis.Send(redisConn, "del", redis.Args{}.Add(key)...)
	//res,err :=   push.Redis.RedisDo("del",redis.Args{}.Add(key)... )
	util.MyPrint(" delOnePush (", id, ")", res, err)
	//util.MyPrint(" delOnePush (",id,")",res,err)

	push.delOneStatus(redisConn, id)
}

func (push *Push) Close() {
	push.CloseChan <- 1
}
