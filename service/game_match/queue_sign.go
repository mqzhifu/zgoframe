package gamematch

//此类主要是操作：报名的 redis queue 数据
import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"zgoframe/service"
	"zgoframe/util"
)

type QueueSign struct {
	Mutex                  sync.Mutex //计算匹配结果时，要加锁  1、阻塞住报名，2、阻塞其它匹配的协程
	RedisCommExpireTime    int        //redis 数据 公共失效时间，用于兜底，如果使用记不失效，进程异常退出 ，再次启动检查timeout又异常，那么用户就永远不能再参与匹配机制了
	Rule                   *Rule      //父类
	RedisKeySeparator      string
	RedisTextSeparator     string
	RedisIdSeparator       string
	RedisPayloadSeparation string
	Log                    *zap.Logger //log 实例
	Redis                  *util.MyRedisGo
	Err                    *util.ErrMsg
	CloseChan              chan int
	Prefix                 string
}

func NewQueueSign(rule *Rule) *QueueSign {
	queueSign := new(QueueSign)
	queueSign.Rule = rule
	queueSign.RedisCommExpireTime = rule.MatchTimeout * 6 //匹配超时的时间 * 6
	queueSign.Redis = rule.RuleManager.Option.GameMatch.Option.Redis
	queueSign.RedisTextSeparator = rule.RuleManager.Option.GameMatch.Option.RedisTextSeparator
	queueSign.RedisKeySeparator = rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
	queueSign.RedisIdSeparator = rule.RuleManager.Option.GameMatch.Option.RedisIdSeparator
	queueSign.RedisPayloadSeparation = rule.RuleManager.Option.GameMatch.Option.RedisPayloadSeparation
	queueSign.Log = rule.RuleManager.Option.GameMatch.Option.Log
	queueSign.Err = rule.RuleManager.Option.GameMatch.Err
	queueSign.CloseChan = make(chan int)
	queueSign.Prefix = rule.Prefix + "_sign"
	return queueSign
}

func (queueSign *QueueSign) Demon() {
	//超时检查，是在匹配的守护协程中调用
	queueSign.Log.Info(queueSign.Prefix + " Demon")
}

//有序集合：每个组的权重值。也可以做(组索引)，所有报名的组，都在这个集合中，weight => group_id
func (queueSign *QueueSign) getRedisKeyWeight() string {
	return queueSign.Rule.GetCommRedisKeyByModuleRuleId(queueSign.Prefix, queueSign.Rule.Id) + "group_weight"
}

//有序集合：组的人数索引，每个规则的池，允许N人成组，其中，每个组里有多少个人，就是这个索引
func (queueSign *QueueSign) getRedisKeyPersonIndexPrefix() string {
	return queueSign.Rule.GetCommRedisKeyByModuleRuleId(queueSign.Prefix, queueSign.Rule.Id) + "group_person"
}

//有序集合：小组人数=>小组id ，如：1人 => groupId ,  2人  => groupId,  3人  => groupId,  4人  => groupId
func (queueSign *QueueSign) getRedisKeyPersonIndex(personNum int) string {
	return queueSign.getRedisKeyPersonIndexPrefix() + queueSign.RedisKeySeparator + strconv.Itoa(personNum)
}

//最简单的string：一个组的详细信息。一个组的详细信息
func (queueSign *QueueSign) getRedisKeyGroupElementPrefix() string {
	return queueSign.Rule.GetCommRedisKeyByModuleRuleId(queueSign.Prefix, queueSign.Rule.Id) + "group_element"
}

//一个组，简单的KV类型，即：set name string
func (queueSign *QueueSign) getRedisKeyGroupElement(id int) string {
	return queueSign.getRedisKeyGroupElementPrefix() + queueSign.RedisKeySeparator + strconv.Itoa(id)
}

//有序集合：一个小组，包含的所有玩家ID.  groupId => playerId
func (queueSign *QueueSign) getRedisKeyGroupPlayer() string {
	return queueSign.Rule.GetCommRedisKeyByModuleRuleId(queueSign.Prefix, queueSign.Rule.Id) + "group_player"
}

//有序集合:组的超时索引.  超时时间 => groupId
func (queueSign *QueueSign) getRedisKeyGroupSignTimeout() string {
	return queueSign.Rule.GetCommRedisKeyByModuleRuleId(queueSign.Prefix, queueSign.Rule.Id) + "timeout"
}

//===============================以上是 redis key 相关============================

//获取当前所有，已报名的，所有组，总数
func (queueSign *QueueSign) getAllGroupsWeightCnt() int {
	return queueSign.getGroupsWeightCnt("-inf", "+inf")
}

//获取当前所有，已报名的，所有组，总数
func (queueSign *QueueSign) getGroupsWeightCnt(rangeStart string, rangeEnd string) int {
	key := queueSign.getRedisKeyWeight()
	res, err := redis.Int(queueSign.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add(rangeStart).Add(rangeEnd)...))
	if err != nil {
		util.ExitPrint("ZCOUNT err", err.Error())
	}
	return res
}

//获取当前所有，已报名的，组，总数
func (queueSign *QueueSign) getAllGroupPersonCnt() map[int]int {
	groupPersonNum := make(map[int]int)
	for i := 1; i <= queueSign.Rule.TeamMaxPeople; i++ {
		groupPersonNum[i] = queueSign.getAllGroupPersonIndexCnt(i)
	}
	return groupPersonNum
}
func (queueSign *QueueSign) getAllGroupPersonIndexCnt(personNum int) int {
	return queueSign.getGroupPersonIndexCnt(personNum, "-inf", "+inf")
}

//获取当前所有，已报名的，组，总数
func (queueSign *QueueSign) getGroupPersonIndexCnt(personNum int, rangeStart string, rangeEnd string) int {
	key := queueSign.getRedisKeyPersonIndex(personNum)
	res, err := redis.Int(queueSign.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add(rangeStart).Add(rangeEnd)...))
	if err != nil {
		util.ExitPrint("ZCOUNT err", err.Error())
	}
	return res
}

//获取当前所有，已报名的，玩家，总数
func (queueSign *QueueSign) getAllPlayersCnt() int {
	return queueSign.getPlayersCnt("-inf", "+inf")
}

//获取当前所有，已报名的，玩家，总数,这个是基于groupId
func (queueSign *QueueSign) getPlayersCnt(rangeStart string, rangeEnd string) int {
	key := queueSign.getRedisKeyGroupPlayer()
	res, err := redis.Int(queueSign.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add(rangeStart).Add(rangeEnd)...))
	if err != nil {
		util.ExitPrint("ZCOUNT err", err.Error())
	}
	return res
}

//获取当前所有，已报名的，玩家，总数,这个是基于 权重
func (queueSign *QueueSign) getPlayersCntTotalByWeight(rangeStart string, rangeEnd string) int {
	total := 0
	for i := 1; i <= queueSign.Rule.TeamMaxPeople; i++ {
		oneCnt := queueSign.getGroupPersonIndexCnt(i, rangeStart, rangeEnd)
		total += oneCnt * i

	}
	return total
}

func (queueSign *QueueSign) getPlayersCntByWeight(rangeStart string, rangeEnd string) map[int]int {
	groupPersonNum := make(map[int]int)
	for i := 1; i <= queueSign.Rule.TeamMaxPeople; i++ {
		groupPersonNum[i] = queueSign.getGroupPersonIndexCnt(i, rangeStart, rangeEnd)
	}
	return groupPersonNum
}

func (queueSign *QueueSign) cancelByPlayerId(playerId int) {
	//queueSign.delOneByPlayerId(playerId)
}

func (queueSign *QueueSign) CancelByGroupId(groupId int) error {
	group := queueSign.getGroupElementById(groupId)
	queueSign.Log.Info("cancelByGroupId groupId : " + strconv.Itoa(groupId))
	//这里是偷懒了，判断是否为空，按说应该返回2个参数，但这个方法调用 地方多，先这样
	if group.Id == 0 || group.Person == 0 || len(group.Players) == 0 {
		msg := queueSign.Err.MakeOneStringReplace(strconv.Itoa(groupId))
		return queueSign.Err.NewReplace(750, msg)
	}
	redisConn := queueSign.Redis.GetNewConnFromPool()
	defer redisConn.Close()
	//检查超时，这种概率比较低，一般守护协程会提前处理的
	if group.SignTime > util.GetNowTimeSecondToInt() {
		queueSign.Log.Warn("group timeout:" + strconv.Itoa(groupId))
	}
	////检查每个玩家的当前状态 ，是否已经超时了
	//for _, player := range group.Players {
	//	playerElement, err := queueSign.Rule.RuleManager.GetById(player.Id)
	//	if err != nil { //证明为空，未取到
	//		if playerElement.Status != service.GAME_MATCH_PLAYER_STATUS_SIGN {
	//			msg := queueSign.Err.MakeOneStringReplace(strconv.Itoa(playerElement.Status))
	//			return queueSign.Err.NewReplace(623, msg)
	//		}
	//	}
	//}
	//开始真正删除一个小组
	queueSign.delOneRuleOneGroup(redisConn, groupId, 1)
	return nil
}

//============================以上是 各种获取总数===============================

func (queueSign *QueueSign) getGroupElementById(id int) (group Group) {
	key := queueSign.getRedisKeyGroupElement(id)
	res, _ := redis.String(queueSign.Redis.RedisDo("get", key))
	if res == "" {
		queueSign.Log.Error(" getGroupElementById is empty," + strconv.Itoa(id))
		return group
	}
	group = queueSign.Rule.RuleManager.Option.GameMatch.GroupStrToStruct(res)
	return group
}

func (queueSign *QueueSign) getGroupIdByPlayerId(playerId int) int {
	key := queueSign.getRedisKeyGroupPlayer()
	memberIndex, _ := redis.Int(queueSign.Redis.RedisDo("zrank", key, playerId))
	if memberIndex == 0 {
		util.MyPrint(" getGroupElementById is empty!")
	}
	id, _ := redis.Int(queueSign.Redis.RedisDo("zrange", redis.Args{}.Add(key).Add(memberIndex).Add(memberIndex).Add("WITHSCORES")...))
	//group := queueSign.strToStruct(res)
	return id
}

//============================以上都是根据ID，获取一条===============================
func (queueSign *QueueSign) getGroupPersonIndexList(personNum int, rangeStart string, rangeEnd string, limitOffset int, limitCnt int, isDel bool) (ids map[int]int) {
	key := queueSign.getRedisKeyPersonIndex(personNum)
	//这里有个问题，person=>id,person是重复的，如果带着分值一并返回，MAP会有重复值的情况
	argc := redis.Args{}.Add(key).Add(rangeEnd).Add(rangeStart).Add("limit").Add(limitOffset).Add(limitCnt)
	res, err := redis.Ints(queueSign.Redis.RedisDo("ZREVRANGEBYSCORE", argc...))
	if err != nil {
		queueSign.Log.Error("getPersonRangeList err" + err.Error())
	}
	if len(res) <= 0 {
		return ids
	}

	rs := util.ArrCovertMap(res)

	if isDel {
		redisConn := queueSign.Redis.GetNewConnFromPool()
		defer redisConn.Close()
		util.MyPrint("getGroupPersonIndexList del group:", res)
		for _, groupId := range res {
			queueSign.delOneRuleOneGroupIndex(redisConn, groupId, personNum)
		}
	}

	return rs
}

//============================以上都是范围性的取值===============================

//获取超时索引集合的，所有成员
//这里间接等于：获取当前所有groupId,取个巧
func (queueSign *QueueSign) GetGroupSignTimeoutAll() []int {
	//queueSign.Redis.Multi(redisConn)
	//defer queueSign.Redis.
	groupSignTimeoutKey := queueSign.getRedisKeyGroupSignTimeout()
	list, _ := redis.Ints(queueSign.Redis.RedisDo("ZRANGE", redis.Args{}.Add(groupSignTimeoutKey).Add(0).Add(-1)...))
	//queueSign.Redis.Exec(redisConn)
	return list
}

//报名
func (queueSign *QueueSign) AddOne(group Group, connFd redis.Conn) {
	queueSign.Log.Debug(queueSign.Prefix + " , add sign one group , start :")
	group.SignTime = util.GetNowTimeSecondToInt()

	groupIndexKey := queueSign.getRedisKeyWeight()
	queueSign.Redis.Send(connFd, "zadd", redis.Args{}.Add(groupIndexKey).Add(group.Weight).Add(group.Id)...)
	//res,err := queueSign.Redis.RedisDo("zadd",redis.Args{}.Add(groupIndexKey).Add(group.Weight).Add(group.Id)...)

	//util.MyPrint("add GroupWeightIndex rs : ", res, err)

	PersonIndexKey := queueSign.getRedisKeyPersonIndex(group.Person)
	queueSign.Redis.Send(connFd, "zadd", redis.Args{}.Add(PersonIndexKey).Add(group.Weight).Add(group.Id)...)
	//res,err = queueSign.Redis.RedisDo("zadd",redis.Args{}.Add(PersonIndexKey).Add(group.Weight).Add(group.Id)...)
	//util.MyPrint("add GroupPersonIndex ( ", group.Person, " ) rs : ", res, err)

	groupSignTimeoutKey := queueSign.getRedisKeyGroupSignTimeout()
	queueSign.Redis.Send(connFd, "zadd", redis.Args{}.Add(groupSignTimeoutKey).Add(group.SignTimeout).Add(group.Id)...)
	//res,err = queueSign.Redis.RedisDo("zadd",redis.Args{}.Add(groupSignTimeoutKey).Add(group.SignTimeout).Add(group.Id)...)
	//util.MyPrint("add GroupSignTimeout rs : ", res, err)

	groupPlayersKey := queueSign.getRedisKeyGroupPlayer()
	for _, v := range group.Players {
		queueSign.Redis.Send(connFd, "zadd", redis.Args{}.Add(groupPlayersKey).Add(group.Id).Add(v.Id)...)
		//res,err = queueSign.Redis.RedisDo("zadd",redis.Args{}.Add(groupPlayersKey).Add(group.Id).Add(v.Id)...)
		//util.MyPrint("add player rs : ", res, err)
	}

	groupElementRedisKey := queueSign.getRedisKeyGroupElement(group.Id)
	content := queueSign.Rule.RuleManager.Option.GameMatch.GroupStructToStr(group)
	queueSign.Redis.Send(connFd, "SET", redis.Args{}.Add(groupElementRedisKey).Add(content)...)
	//res,err = queueSign.Redis.RedisDo("set",redis.Args{}.Add(groupElementRedisKey).Add(content)...)
	//util.MyPrint("add groupElement rs : ", res, err)

	queueSign.Log.Debug(queueSign.Prefix + " , add sign one group , finish . ")
}

//==============以下均是 删除操作======================================

////删除 所有Rule：池里的报名组、玩家、索引等-有点暴力，尽量不用
//func  (queueSign *QueueSign)  delAll(){
//	key := queueSign.getRedisPrefixKey()
//	queueSign.Redis.RedisDo("del",key)
//}
//删除一条规则的所有匹配信息
func (queueSign *QueueSign) delOneRule() {
	queueSign.Log.Info(" queueSign delOneRule : ")
	keys := queueSign.Rule.GetCommRedisKeyByModuleRuleId("sign", queueSign.Rule.Id) + "*"
	queueSign.Redis.RedisDelAllByPrefix(keys)

	//queueSign.delOneRuleALLGroupElement()
	//queueSign.delOneRuleALLPersonIndex()
	//queueSign.delOneRuleAllGroupSignTimeout()
	//queueSign.delOneRuleALLPlayers()
	//queueSign.delOneRuleALLWeight()
}

//====================================================

////删除一条规则的，所有玩家索引
//func  (queueSign *QueueSign)  delOneRuleALLPlayers( ){
//	key := queueSign.getRedisKeyGroupPlayer()
//	res,_ := redis.Int(queueSign.Redis.RedisDo("del",key))
//	mylog.Debug("delOneRuleALLPlayers : ",res)
//}
////删除一条规则的，所有权重索引
//func  (queueSign *QueueSign)  delOneRuleALLWeight( ){
//	key := queueSign.getRedisKeyWeight()
//	res,_ := redis.Int(queueSign.Redis.RedisDo("del",key))
//	mylog.Debug("delOneRuleALLPlayers : ",res)
//}
//删除一条规则的，所有人数分组索引
//func  (queueSign *QueueSign)  delOneRuleALLPersonIndex( ){
//	for i:=1 ; i <= RuleGroupPersonMax;i++{
//		queueSign.delOneRuleOnePersonIndex(i)
//	}
//}
//
////删除一条规则的，所有分组详细信息
//func  (queueSign *QueueSign)  delOneRuleALLGroupElement( ){
//	prefix := queueSign.getRedisKeyGroupElementPrefix()
//	res,_ := redis.Strings( queueSign.Redis.RedisDo("keys",prefix + "*"  ))
//	if len(res) == 0{
//		mylog.Notice(" GroupElement by keys(*) : is empty")
//		return
//	}
//	//zlib.ExitPrint(res,-200)
//	for _,v := range res{
//		res,_ := redis.Int(queueSign.Redis.RedisDo("del",v))
//		mylog.Debug("del group element v :",res)
//	}
//}

////删除一条规则的，所有人组超时索引
//func  (queueSign *QueueSign)  delOneRuleAllGroupSignTimeout( ){
//	key := queueSign.getRedisKeyGroupSignTimeout()
//	res,_ := redis.Int(queueSign.Redis.RedisDo("del",key))
//	mylog.Debug("delOneRuleALLPlayers : ",res)
//}

//删除一条规则的，某一人数各类的，所有人数分组索引
func (queueSign *QueueSign) delOneRuleOnePersonIndex(personNum int) {
	key := queueSign.getRedisKeyPersonIndex(personNum)
	res, _ := redis.Int(queueSign.Redis.RedisDo("del", key))
	util.MyPrint("delOneRuleALLPlayers : ", res)
}

//====================================================

func (queueSign *QueueSign) delOneRuleOnePersonIndexById(redisConn redis.Conn, personNum int, id int) {
	key := queueSign.getRedisKeyPersonIndex(personNum)
	queueSign.Redis.Send(redisConn, "ZREM", redis.Args{}.Add(key).Add(id)...)
}

//删除一个组的所有玩家信息
func (queueSign *QueueSign) delOneRuleOneGroupPlayers(redisConn redis.Conn, id int) {
	key := queueSign.getRedisKeyGroupPlayer()
	queueSign.Redis.Send(redisConn, "ZREMRANGEBYSCORE", redis.Args{}.Add(key).Add(id).Add(id)...)

}

//删除一条规则的一个组的详细信息
func (queueSign *QueueSign) delOneRuleOneGroupSignTimeout(redisConn redis.Conn, id int) {
	key := queueSign.getRedisKeyGroupSignTimeout()
	queueSign.Redis.Send(redisConn, "ZREM", redis.Args{}.Add(key).Add(id)...)
}

//删除一条规则的权限分组索引
func (queueSign *QueueSign) delOneRuleOneWeight(redisConn redis.Conn, id int) {
	key := queueSign.getRedisKeyWeight()
	queueSign.Redis.Send(redisConn, "ZREM", redis.Args{}.Add(key).Add(id)...)
}

//删除一条规则的一个组的详细信息
func (queueSign *QueueSign) delOneRuleOneGroupElement(redisConn redis.Conn, id int) {
	key := queueSign.getRedisKeyGroupElement(id)
	queueSign.Redis.Send(redisConn, "del", key)
	//res,_ := queueSign.Redis.RedisDo("del",key)
	//util.MyPrint("delOneGroupElement : ", res)
	//util.MyPrint("delOneGroupElement ", id)
}

////删除一个玩家的报名
//func (queueSign *QueueSign) delOneByPlayerId( playerId int){
//	id := queueSign.getGroupIdByPlayerId(playerId)
//	if id <=0 {
//		mylog.Error("getGroupIdByPlayerId is empty !!!")
//		return
//	}
//	queueSign.delOneRuleOneGroup(id,1)
//	//这里还要更新当前组所有用户的，状态信息
//}

//====================================================
//删除一个组
func (queueSign *QueueSign) delOneRuleOneGroup(redisConn redis.Conn, id int, isDelPlayerStatus int) error {
	queueSign.Log.Warn("queueSign delOneRuleOneGroup id:" + strconv.Itoa(id) + " isDelPlayerStatus:" + strconv.Itoa(isDelPlayerStatus))

	group := queueSign.getGroupElementById(id)
	//这里是偷懒了，判断是否为空，按说应该返回2个参数，但这个方法调用 地方多，先这样
	if group.Id == 0 || group.Person == 0 || len(group.Players) == 0 {
		msg := queueSign.Err.MakeOneStringReplace(strconv.Itoa(id))
		nowerr := queueSign.Err.NewReplace(750, msg)

		queueSign.Log.Error("750" + nowerr.Error())
		queueSign.Log.Error("750" + nowerr.Error())

		return nowerr
	}

	queueSign.delOneRuleOneGroupElement(redisConn, id)
	queueSign.delOneRuleOneWeight(redisConn, id)
	queueSign.delOneRuleOnePersonIndexById(redisConn, group.Person, id)
	queueSign.delOneRuleOneGroupPlayers(redisConn, id)
	queueSign.delOneRuleOneGroupSignTimeout(redisConn, id)

	if isDelPlayerStatus == 1 {
		//这里是直接删除，比较好的方法是：先查询下再删除，另外：删除不如修改该键的状态~
		for _, player := range group.Players {
			queueSign.Rule.PlayerManager.delOneById(redisConn, player.Id)
		}
	}

	queueSign.Log.Warn("queueSign delOneRuleOneGroup finish, id:" + strconv.Itoa(id))

	return nil
}

////删除一个组的，所有索引信息,反向看：除了小组基础信息外，其余均删除
//func  (queueSign *QueueSign)  delOneGroupIndex(groupId int){
//	zlib.MyPrint(" delOneGroupIndex : ",queueSign.getRedisCatePrefixKey())
//	group := queueSign.getGroupElementById(groupId)
//
//	queueSign.delOneRuleOneWeight(groupId)
//	queueSign.delOneRuleOnePersonIndexById(group.Person,groupId)
//	queueSign.delOneRuleOneGroupPlayers(groupId)
//	queueSign.delOneRuleOneGroupSignTimeout(groupId)
//}
//添加一个组的:<权限+人数>,<人数+权重>的索引
//-匹配时，会把小于的索引先删掉，避免重复，最终结果，可能有些小组还要再PUSH BACK 回来
func (queueSign *QueueSign) addOneGroupIndex(redisConn redis.Conn, groupId int, personNum int, weight float32) {
	//zlib.MyPrint(" addOneGroupIndex : ",queueSign.getRedisCatePrefixKey())
	//group := queueSign.getGroupElementById(groupId)

	groupIndexKey := queueSign.getRedisKeyWeight()
	res, err := queueSign.Redis.RedisDo("zadd", redis.Args{}.Add(groupIndexKey).Add(weight).Add(groupId)...)
	util.MyPrint("add GroupWeightIndex rs : ", res, err)

	PersonIndexKey := queueSign.getRedisKeyPersonIndex(personNum)
	res, err = queueSign.Redis.RedisDo("zadd", redis.Args{}.Add(PersonIndexKey).Add(weight).Add(groupId)...)
	util.MyPrint("add GroupPersonIndex ( ", personNum, " ) rs : ", res, err)
}

//func  (queueSign *QueueSign)  addOneGroupIndex(groupId int){
//	zlib.MyPrint(" addOneGroupIndex : ",queueSign.getRedisCatePrefixKey())
//	group := queueSign.getGroupElementById(groupId)
//
//	groupIndexKey := queueSign.getRedisKeyWeight( )
//	res,err := queueSign.Redis.RedisDo("zadd",redis.Args{}.Add(groupIndexKey).Add(group.Weight).Add(group.Id)...)
//	mylog.Debug("add GroupWeightIndex rs : ",res,err)
//
//	PersonIndexKey := queueSign.getRedisKeyPersonIndex(  group.Person)
//	res,err = queueSign.Redis.RedisDo("zadd",redis.Args{}.Add(PersonIndexKey).Add(group.Weight).Add(group.Id)...)
//	mylog.Debug("add GroupPersonIndex ( ",group.Person," ) rs : ",res,err)
//}

//删除一个组的:<权限+人数>,<人数+权重>的索引
//-匹配的时候，前期计算均是基于索引，在getList时，要把这个索引先删了，等到计算出最终结果，再决定是否全删
//-如果有些索引要放回去，直接恢复这两个索引即可
func (queueSign *QueueSign) delOneRuleOneGroupIndex(redisConn redis.Conn, groupId int, personNum int) {
	queueSign.Log.Info("action : delOneRuleOneGroupIndex id:" + strconv.Itoa(groupId))
	//group := queueSign.getGroupElementById(id)
	//mylog.Debug(group)

	queueSign.delOneRuleOneWeight(redisConn, groupId)
	queueSign.delOneRuleOnePersonIndexById(redisConn, personNum, groupId)
}

//检测 小组 超时（因为player是包含在组里，所以未对player级别细粒度检查）
func (queueSign *QueueSign) CheckTimeout() {
	push := queueSign.Rule.Push
	keys := queueSign.getRedisKeyGroupSignTimeout()

	now := util.GetNowTimeSecondToInt()
	res, err := redis.IntMap(queueSign.Redis.RedisDo("ZREVRANGEBYSCORE", redis.Args{}.Add(keys).Add(now).Add("-inf").Add("WITHSCORES")...))
	if err != nil {
		queueSign.Log.Error("redis ZREVRANGEBYSCORE err :" + err.Error())
		return
	}

	if len(res) == 0 {
		queueSign.Rule.NothingToDoLog("queueSign CheckTimeout empty , no need process")
		return
	}
	//走到这里，证明，redis 中已有 玩家 报名失效了，需要做处理了
	queueSign.Log.Info("sign timeout group element total : " + strconv.Itoa(len(res)))
	redisConn := queueSign.Redis.GetNewConnFromPool()
	defer redisConn.Close()
	for groupIdStr, _ := range res {
		queueSign.Redis.Multi(redisConn)

		groupId := util.Atoi(groupIdStr)                //redis 数据全是string 转成int
		group := queueSign.getGroupElementById(groupId) //根据ID获取小组的详细信息
		queueSign.Log.Warn(queueSign.Prefix + " group timeout , id:" + groupIdStr + " , playerIds:" + queueSign.Rule.RuleManager.Option.GameMatch.GetGroupPlayerIds(group) + " timeout:" + strconv.Itoa(group.MatchTimes))

		payload := queueSign.Rule.RuleManager.Option.GameMatch.GroupStructToStr(group) //小组是个结构体，要存redis得转成字符串
		payload = strings.Replace(payload, queueSign.RedisTextSeparator, queueSign.RedisPayloadSeparation, -1)

		pushElement := push.addOnePush(redisConn, groupId, service.PushCategorySignTimeout, payload) //添加一条推送消息
		queueSign.Log.Info("delOneRuleOneGroup")
		queueSign.delOneRuleOneGroup(redisConn, groupId, 1)

		queueSign.Redis.Exec(redisConn)
		queueSign.Rule.RuleManager.Option.GameMatch.PersistenceRecordSuccessPush(pushElement, queueSign.Rule.Id)

	}
	//rootAndSingToLogInfoMsg(queueSign,"queueSign checkTimeOut finish in oneTime")
	queueSign.Log.Info("queueSign checkTimeOut finish in oneTime")
	//myGosched("sign CheckTimeout")
	//mySleepSecond(1," sign CheckTimeout ")
}

//用于 测试用例
type SignTotalCnt struct {
	GroupsWeightCnt     int
	GroupPersonIndexCnt int
	GroupSignTimeoutCnt int
	GroupPlayerCnt      int
}

//用于 测试用例
func (queueSign *QueueSign) TestTotalCnt() SignTotalCnt {
	signTotalCnt := SignTotalCnt{}
	data := make(map[string]int)
	signTotalCnt.GroupsWeightCnt = queueSign.getGroupsWeightCnt("-inf", "+inf")
	GroupPersonIndexTotal := 0
	for i := 1; i <= 5; i++ {
		GroupPersonIndexCnt := queueSign.getGroupPersonIndexCnt(i, "-inf", "+inf")
		GroupPersonIndexTotal += GroupPersonIndexCnt
		//data["groupPersonIndexCnt_"+strconv.Itoa(i)] = GroupPersonIndexCnt
	}
	signTotalCnt.GroupPersonIndexCnt = GroupPersonIndexTotal

	GetGroupSignTimeoutAllArr := queueSign.GetGroupSignTimeoutAll()
	signTotalCnt.GroupSignTimeoutCnt = len(GetGroupSignTimeoutAllArr)

	key := queueSign.getRedisKeyGroupElement(0)
	key = key[0 : len(key)-1]
	res, _ := redis.Strings(queueSign.Redis.RedisDo("keys", key+"*"))
	data["groupElement"] = len(res)

	groupPlayerCnt := 0
	for _, v := range res {
		split := strings.Split(v, "_")
		groupIdStr := split[len(split)-1]
		groupId, _ := strconv.Atoi(groupIdStr)
		groupInfo := queueSign.getGroupElementById(groupId)
		//zlib.MyPrint(len(groupInfo.Players))
		groupPlayerCnt += len(groupInfo.Players)
	}
	signTotalCnt.GroupPlayerCnt = groupPlayerCnt
	return signTotalCnt
	//queueSign.getRedisKeyGroupPlayer()
	//gameMatchInstance.PlayerStatus.getOneRuleAllPlayerCnt
}

func (queueSign *QueueSign) Close() {
	queueSign.CloseChan <- 1
}

func (queueSign *QueueSign) TestRedisKey() {
	redisKey := queueSign.getRedisKeyWeight()
	util.MyPrint("queueSign test :", redisKey)

	redisKey = queueSign.getRedisKeyPersonIndexPrefix()
	util.MyPrint("queueSign test :", redisKey)

	redisKey = queueSign.getRedisKeyPersonIndex(1)
	util.MyPrint("queueSign test :", redisKey)

	redisKey = queueSign.getRedisKeyGroupElementPrefix()
	util.MyPrint("queueSign test :", redisKey)

	redisKey = queueSign.getRedisKeyGroupElement(1)
	util.MyPrint("queueSign test :", redisKey)

	redisKey = queueSign.getRedisKeyGroupPlayer()
	util.MyPrint("queueSign test :", redisKey)

	redisKey = queueSign.getRedisKeyGroupSignTimeout()
	util.MyPrint("queueSign test :", redisKey)

}
