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
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/util"
)

//推送给3方，支持重试
type PushElement struct {
	Id       int    `json:"id"`
	ATime    int    `json:"a_time"`   //添加时间
	UTime    int    `json:"u_time"`   //最后更新的时间
	LinkId   int    `json:"link_id"`  //关联调方用的ID
	Status   int    `json:"status"`   //状态：1未推送2推送失败，等待重试3推送成功4推送失败，不再重试
	Times    int    `json:"times"`    //已推送次数
	Category int    `json:"category"` //1:报名超时 2匹配成功 3成功结果超时
	Payload  string `json:"payload"`  //自定义的载体
}

type Push struct {
	Mutex                  sync.Mutex
	Service                string
	Rule                   *Rule //父类
	RedisKeySeparator      string
	RedisTextSeparator     string
	RedisIdSeparator       string
	RedisPayloadSeparation string
	Log                    *zap.Logger //log 实例
	Redis                  *util.MyRedisGo
	Err                    *util.ErrMsg
	CloseChan              chan int
	prefix                 string
	RetryPeriod            []int
}

func NewPush(rule *Rule) *Push {
	push := new(Push)

	push.Rule = rule
	push.Redis = rule.RuleManager.Option.GameMatch.Option.Redis
	push.RedisTextSeparator = rule.RuleManager.Option.GameMatch.Option.RedisTextSeparator
	push.RedisKeySeparator = rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
	push.RedisIdSeparator = rule.RuleManager.Option.GameMatch.Option.RedisIdSeparator
	push.RedisPayloadSeparation = rule.RuleManager.Option.GameMatch.Option.RedisPayloadSeparation
	push.Log = rule.RuleManager.Option.GameMatch.Option.Log
	push.Err = rule.RuleManager.Option.GameMatch.Err
	push.CloseChan = make(chan int)
	push.Service = "roomService"
	push.prefix = rule.Prefix + "_push"
	//var PushRetryPeriod = []int{10,30,60,600}
	push.RetryPeriod = []int{2, 4, 6} //方便测试
	return push
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

func (push *Push) getById(id int) (element PushElement) {
	key := push.getRedisKeyPush(id)
	res, _ := redis.String(push.Redis.RedisDo("get", key))
	if res == "" {
		return element
	}

	element = push.pushStrToStruct(res)
	return element
}

func (push *Push) addOnePush(redisConn redis.Conn, linkId int, category int, payload string) PushElement {
	push.Log.Debug("addOnePush linkId:" + strconv.Itoa(linkId) + " category:" + strconv.Itoa(category))
	id := push.GetPushIncId()
	key := push.getRedisKeyPush(id)
	pushElement := PushElement{
		Id:       id,
		ATime:    util.GetNowTimeSecondToInt(),
		Status:   service.PUSH_STATUS_WAIT,
		UTime:    util.GetNowTimeSecondToInt(),
		Times:    0,
		LinkId:   linkId,
		Category: category,
		Payload:  payload,
	}
	pushStr := push.pushStructToStr(pushElement)
	push.Redis.Send(redisConn, "set", redis.Args{}.Add(key).Add(pushStr)...)

	pushKey := push.getRedisKeyPushStatus()
	push.Redis.Send(redisConn, "zadd", redis.Args{}.Add(pushKey).Add(service.PUSH_STATUS_WAIT).Add(id)...)
	pushElementBytes, _ := json.Marshal(&pushElement)
	push.Log.Debug("addOnePush finish ,info:" + string(pushElementBytes))
	return pushElement
}

func (push *Push) pushStrToStruct(redisStr string) PushElement {
	strArr := strings.Split(redisStr, push.RedisTextSeparator)
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
		strconv.Itoa(pushElement.Id) + push.RedisTextSeparator +
			strconv.Itoa(pushElement.LinkId) + push.RedisTextSeparator +
			strconv.Itoa(pushElement.ATime) + push.RedisTextSeparator +
			strconv.Itoa(pushElement.Status) + push.RedisTextSeparator +
			strconv.Itoa(pushElement.UTime) + push.RedisTextSeparator +
			strconv.Itoa(pushElement.Times) + push.RedisTextSeparator +
			strconv.Itoa(pushElement.Category) + push.RedisTextSeparator +
			pushElement.Payload + push.RedisTextSeparator

	return str
}

func (push *Push) delOneRule() {
	push.Log.Debug(" push delOneRule : ")
	key := push.Rule.GetCommRedisKeyByModuleRuleId("push", push.Rule.Id) + "*"
	push.Redis.RedisDelAllByPrefix(key)
	//push.delAllPush()
	//push.delAllStatus()
}

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
	push.Log.Info("delOnePush " + strconv.Itoa(pushId) + " key:" + key)
	push.Redis.Send(redisConn, "ZREM", redis.Args{}.Add(key).Add(pushId)...)
}

//失败且需要重试的PUSH-ELEMENT
func (push *Push) upRetryPushInfo(element PushElement) {
	redisConnFD := push.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()

	push.Log.Debug("upRetryPushElementInfo , id: " + strconv.Itoa(element.Id) + " , oldTimes:" + strconv.Itoa(element.Times) + " , oldStatus: " + strconv.Itoa(element.Status))

	push.Redis.Send(redisConnFD, "multi")
	element.Status = service.PUSH_STATUS_RETRY
	element.UTime = util.GetNowTimeSecondToInt()
	element.Times = element.Times + 1 //重试次数+1
	key := push.getRedisKeyPush(element.Id)
	pushStr := push.pushStructToStr(element)
	push.Redis.Send(redisConnFD, "set", redis.Args{}.Add(key).Add(pushStr)...)

	//这里有个麻烦点，元素信息 和 索引信息，是分开放的，元素的变更比较简单，索引是一个集合，改起来有点麻烦,那就直接先删了，再重新添加一条
	push.delOneStatus(redisConnFD, element.Id)
	statusKey := push.getRedisKeyPushStatus()
	//util.MyPrint("del pushStatus index ,pushId : ", element.Id)
	push.Redis.Send(redisConnFD, "zadd", redis.Args{}.Add(statusKey).Add(service.PUSH_STATUS_RETRY).Add(element.Id)...)
	//res,err = push.Redis.RedisDo("zadd",redis.Args{}.Add(statuskey).Add(PushStatusRetry).Add(element.Id)...)
	//util.MyPrint("add  new pushStatus index : ", res, err)
	//mylog.Info("add  new pushStatus index : ",res,err)

	push.Redis.Exec(redisConnFD)
}

//在业务里，删除一条push
//走到里，前置条件肯定是HTTP-PUSH成功了
//除了删除push相关的数据外，还得删除连带着的业务数据，这个有点烦
func (push *Push) delOneByIdInBusiness(redisConn redis.Conn, id int) {
	push.Log.Debug("delOneByIdInBusiness id:" + strconv.Itoa(id))
	push.Redis.Send(redisConn, "multi")
	element := push.getById(id)
	push.delOnePush(redisConn, id)
	if element.Category == service.PushCategorySuccess || element.Category == service.PushCategorySuccessTimeout {
		push.Log.Info("delOneResult")
		success := push.Rule.QueueSuccess
		//删除连带的业务信息，这里正确的理解应该是：业务的：回调函数
		success.delOneResult(redisConn, element.LinkId, 1, 1, 1, 1)
	} else {
		push.Log.Error("delOneByIdInBusiness Category err")
	}
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
			time.Sleep(time.Millisecond * time.Duration(push.Rule.RuleManager.Option.GameMatch.Option.LoopSleepTime))
		}
	}
forEnd:
	push.Log.Warn(push.prefix + "  Demon end .")
}

//检查需要抢着的数据：待推送、重试推送
func (push *Push) checkStatus() {
	key := push.getRedisKeyPushStatus()
	//redis 集合里的数据，其实一次均可以取出来，但拆分成两个状态：优先计算 正常推送，而重试的推送放在后面
	push.checkOneByStatus(key, service.PUSH_STATUS_WAIT)
	push.checkOneByStatus(key, service.PUSH_STATUS_RETRY)
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
	res, err := redis.Ints(push.Redis.RedisDo("ZREVRANGEBYSCORE", redis.Args{}.Add(key).Add(status).Add(status)...))
	if err != nil {
		push.Log.Error("redis keys err :" + err.Error())
		return
	}
	if len(res) == 0 {
		push.Rule.NothingToDoLog(push.prefix + " checkOneByStatus :" + strconv.Itoa(status) + " empty , no need process")
		return
	}
	push.Log.Info("push need process element total : " + strconv.Itoa(len(res)) + " status: " + strconv.Itoa(status))
	for _, id := range res {
		if status == service.PUSH_STATUS_WAIT { //正常/首次推送
			push.processWaitOne(id, status)
		} else { //推送失败的，重试推送
			push.processRetryOne(id, status)
		}

	}
}

//处理首次推送
func (push *Push) processWaitOne(id int, status int) {
	element := push.getById(id)
	push.Log.Info("processWaitOne , push id : " + strconv.Itoa(id) + " status : " + strconv.Itoa(status) + "(retry) category:" + strconv.Itoa(element.Category))
	//httpRs, err := push.ServiceDiscoveryRequest(element, service.PUSH_STATUS_RETRY)
	httpRs, err := push.ServiceDiscoveryRequestUser(element)

	if err != nil {
		push.upRetryPushInfo(element)
		msg := push.Err.MakeOneStringReplace(err.Error())
		push.Err.NewReplace(911, msg)
		push.Log.Error("push third service-1 " + err.Error())
		return
	}
	push.hook(element, httpRs)
	push.Log.Info("processOne finish")
}

//处理重试推送
func (push *Push) processRetryOne(id int, status int) {
	element := push.getById(id)
	push.Log.Info("processRetryOne , push id : " + strconv.Itoa(id) + " status : " + strconv.Itoa(status) + "(wait) category:" + strconv.Itoa(element.Category))
	if element.Times >= len(push.RetryPeriod) {
		//已超过，最大重试次数
		redisConnFD := push.Redis.GetNewConnFromPool()
		defer redisConnFD.Close()

		push.Log.Warn(" push retry time(" + strconv.Itoa(element.Times) + ") > maxRetryTime(" + strconv.Itoa(len(push.RetryPeriod)) + ") , drop this msg.")
		push.delOneByIdInBusiness(redisConnFD, id)
		return
	}

	time := push.RetryPeriod[element.Times]
	d := util.GetNowTimeSecondToInt() - element.UTime
	//util.MyPrint("this time : ", time, "now :", util.GetNowTimeSecondToInt(), " - element.UTime ", element.UTime, " = ", d)
	if d >= time {
		push.Log.Info("trigger retry Push!!! ")
		//httpRs, err := push.ServiceDiscoveryRequest(element, service.PUSH_STATUS_RETRY)
		httpRs, err := push.ServiceDiscoveryRequestUser(element)
		if err != nil {
			push.upRetryPushInfo(element)
			msg := push.Err.MakeOneStringReplace(err.Error())
			push.Err.NewReplace(911, msg)
			push.Log.Error("push third service-1 " + err.Error())
			return
		}
		push.hook(element, httpRs)
	} else {
		//重试的时间间隔 未满足
		push.Log.Info("The retry time is too short ,no trigger the Push!!! ")
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
		//push.Log.Debug("element.Category == PushCategorySignTimeout")
		postData = push.Rule.RuleManager.Option.GameMatch.GroupStrToStruct(payload)
		thirdMethodUri = "v1/match/error"
	} else if element.Category == service.PushCategorySuccessTimeout {
		//push.Log.Debug("element.Category == PushCategorySuccessTimeout")
		postData = push.Rule.QueueSuccess.strToStruct(payload)
		thirdMethodUri = "v1/match/error"
	} else if element.Category == service.PushCategorySuccess {
		//push.Log.Debug("element.Category == PushCategorySuccess")
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

//没有单独的 room 服务做聚合，直接返回给用户
func (push *Push) ServiceDiscoveryRequestUser(element PushElement) (httpRs util.ResponseMsgST, err error) {

	payload := strings.Replace(element.Payload, push.RedisPayloadSeparation, push.RedisTextSeparator, -1)
	success := push.Rule.QueueSuccess
	if element.Category == service.PushCategorySignTimeout {
		groupInfo := push.Rule.RuleManager.Option.GameMatch.GroupStrToStruct(payload)
		playerIds := push.Rule.RuleManager.Option.GameMatch.GetGroupPlayerIds(groupInfo)
		gameMatchOptResult := pb.GameMatchOptResult{
			GroupId: int32(groupInfo.Id),
			RoomId:  "",
			Code:    400,
			Msg:     "sign timeout",
		}
		push.Rule.RuleManager.Option.RequestServiceAdapter.GatewaySendMsgByUids(playerIds, "SC_GameMatchOptResult", gameMatchOptResult)
	} else if element.Category == service.PushCategorySuccessTimeout {
		return httpRs, nil
	} else if element.Category == service.PushCategorySuccess {
		thisResult := success.strToStruct(payload)
		resultInfo, _ := success.GetResultById(thisResult.Id, 1, 0)

		newRoom := push.Rule.RuleManager.Option.GameMatch.Option.FrameSync.RoomManage.NewEmptyRoom()
		util.MyPrint("newRoom:", newRoom)
		newRoom.RuleId = int32(push.Rule.Id)
		for _, uid := range resultInfo.PlayerIds {
			newRoom.AddPlayer(uid)
			gameMatchOptResult := pb.GameMatchOptResult{
				//GroupId: resultInfo.Groups,
				RoomId: newRoom.Id,
				Code:   200,
				Msg:    "success",
			}
			push.Rule.RuleManager.Option.RequestServiceAdapter.GatewaySendMsgByUid(int32(uid), "SC_GameMatchOptResult", gameMatchOptResult)
		}

	} else {
		push.Log.Error("element.Category error.")
		return httpRs, errors.New("element.Category error")
	}

	return httpRs, nil
}

func (push *Push) ServiceDiscoveryRequest(element PushElement, upStatus int) (httpRs util.ResponseMsgST, err error) {
	push.Log.Debug("ServiceDiscoveryRequest id:" + strconv.Itoa(element.Id) + " status:" + strconv.Itoa(element.Status) + " , upStatus:" + strconv.Itoa(upStatus) + " category:" + strconv.Itoa(element.Category))
	//var httpRs util.ResponseMsgST
	//var err error

	payload := strings.Replace(element.Payload, push.RedisPayloadSeparation, push.RedisTextSeparator, -1)
	thirdMethodUri, postData, err := push.getServiceUri(element, payload)
	if err != nil {
		return httpRs, err
	}
	push.Log.Debug("push third service ,  uri : " + thirdMethodUri)

	myServiceDiscovery := push.Rule.RuleManager.Option.GameMatch.Option.ServiceDiscovery
	projectId := push.Rule.RuleManager.Option.GameMatch.Option.ProjectId
	myService, err := myServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(push.Service, "")
	if err != nil {
		push.Log.Error("myServiceDiscovery err1:" + err.Error())
		return httpRs, err
	}

	serviceHttp := util.NewServiceHttp(projectId, push.Service, myService.Ip, myService.Port, myService.ServiceId)
	httpRs, err = serviceHttp.Post(thirdMethodUri, postData)
	return httpRs, nil
}

//到了这一步，一定是发送了http请求，并拿到了http的响应数据
func (push *Push) hook(element PushElement, httpRs util.ResponseMsgST) {
	push.Log.Debug("hook id:" + strconv.Itoa(element.Id) + " status:" + strconv.Itoa(element.Status) + " category:" + strconv.Itoa(element.Category))
	redisConnFD := push.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()

	if httpRs.Code == 0 { // 0 即是200 ，推送成功~
		push.delOneByIdInBusiness(redisConnFD, element.Id)
		return
	}
	push.Log.Debug("hook upRetryPushInfo")
	push.upRetryPushInfo(element)
	//if element.Category == service.PushCategorySignTimeout {
	//	if httpRs.Code == 116 || httpRs.Code == 119 {
	//		//push.metrics(element.Category, "Ok")
	//		push.delOnePush(redisConnFD, element.Id)
	//	} else {
	//		push.upRetryPushInfo(element)
	//
	//		httpRsJsonStr, _ := json.Marshal(httpRs)
	//		msg := push.Err.MakeOneStringReplace(string(httpRsJsonStr))
	//		//msg := myerr.MakeOneStringReplace(string(httpRsJsonStr))
	//		push.Err.NewReplace(700, msg)
	//		return
	//	}
	//} else if element.Category == service.PushCategorySuccessTimeout {
	//	push.Log.Info("delOneResult")
	//	push.delOneByIdInBusiness(redisConnFD, element.Id)
	//	//push.delOnePush(element.Id)
	//	//push.QueueSuccess.delOneResult(element.LinkId,1,1,1,1)
	//
	//} else if element.Category == service.PushCategorySuccess {
	//	if httpRs.Code == 108 || httpRs.Code == 102 {
	//		push.upRetryPushInfo(element)
	//
	//		httpRsJsonStr, _ := json.Marshal(httpRs)
	//		msg := push.Err.MakeOneStringReplace(string(httpRsJsonStr))
	//		//msg := myerr.MakeOneStringReplace(string(httpRsJsonStr))
	//		push.Err.NewReplace(700, msg)
	//		return
	//	} else {
	//		push.Log.Info("delOneResult")
	//		push.delOneByIdInBusiness(redisConnFD, element.Id)
	//		//push.delOnePush(element.Id)
	//		//push.QueueSuccess.delOneResult(element.LinkId,1,1,1,1)
	//		return
	//	}
	//} else {
	//	push.Log.Error("hook element.Category not found!!!")
	//	push.delOnePush(redisConnFD, element.Id)
	//}
}
func (push *Push) delOnePush(redisConn redis.Conn, id int) {
	key := push.getRedisKeyPush(id)
	push.Log.Info("delOnePush " + strconv.Itoa(id) + " key:" + key)
	push.Redis.Send(redisConn, "del", redis.Args{}.Add(key)...)

	push.delOneStatus(redisConn, id)
}

func (push *Push) Close() {
	push.CloseChan <- 1
}

func (push *Push) TestRedisKey() {
	redisKey := push.getRedisPushIncKey()
	util.MyPrint("push test :", redisKey)

	redisKey = push.getRedisKeyPushStatus()
	util.MyPrint("push test :", redisKey)

	redisKey = push.getRedisKeyPush(1)
	util.MyPrint("push test :", redisKey)

}

//func (push *Push) metrics(elementCategory int, action string) {
//	actionStr := action
//	if action == "Ok" || action == "Drop" {
//
//	} else {
//		mylog.Error("push metrics action err:" + action)
//		return
//	}
//
//	key := ""
//	if elementCategory == PushCategorySignTimeout {
//		key = "PushSign" + actionStr
//	} else if elementCategory == PushCategorySuccess {
//		key = "PushMatchSuccess" + actionStr
//	} else if elementCategory == PushCategorySuccessTimeout {
//		key = "PushMatchSuccessTimeout" + actionStr
//	} else {
//		util.MyPrint("push metrics category err:", elementCategory)
//		return
//	}
//
//	mymetrics.FastLog(key, zlib.METRICS_OPT_INC, 0)
//}

//func (push *Push)   delAll(){
//	key := push.getRedisPrefixKey()
//	push.Redis.RedisDo("del",key)
//}
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

//func (push *Push) pushAndUpInfo(element PushElement, upStatus int) error {
//	push.Log.Debug("pushAndUpInfo id:" + strconv.Itoa(element.Id) + " status:" + strconv.Itoa(element.Status) + "upStatus:" + strconv.Itoa(upStatus) + " category:" + strconv.Itoa(element.Category))
//	var httpRs util.ResponseMsgST
//	var err error
//
//	payload := strings.Replace(element.Payload, push.RedisPayloadSeparation, push.RedisTextSeparator, -1)
//	thirdMethodUri, postData, err := push.getServiceUri(element, payload)
//	if err != nil {
//		return err
//	}
//	push.Log.Debug("push third service ,  uri : " + thirdMethodUri)
//
//	myServiceDiscovery := push.Rule.RuleManager.Option.GameMatch.Option.ServiceDiscovery
//	projectId := push.Rule.RuleManager.Option.GameMatch.Option.ProjectId
//	myService, err := myServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(push.Service, "")
//	if err != nil {
//		push.Log.Error("myServiceDiscovery err1:" + err.Error())
//		return err
//	}
//
//	serviceHttp := util.NewServiceHttp(projectId, push.Service, myService.Ip, myService.Port, myService.ServiceId)
//	httpRs, err = serviceHttp.Post(thirdMethodUri, postData)
//	push.Log.Debug("push third service , httpCode : " + strconv.Itoa(httpRs.Code) + " err : " + err.Error())
//	if err != nil {
//		return err
//	}
//
//	push.hook(element, httpRs)
//	return nil
//}
