package gamematch

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"zgoframe/service"
	"zgoframe/util"
)

// 1. 等待计算   2计算中   3匹配成功
type Player struct {
	Id             int            `json:"id"`
	Status         int            `json:"status"`
	GroupId        int            `json:"group_id"`
	Weight         float32        `json:"weight"`
	ATime          int            `json:"a_time"`
	UTime          int            `json:"u_time"`
	SignTimeout    int            `json:"sign_timeout"`
	SuccessTimeout int            `json:"success_timeout"`
	WeightAttrs    map[string]int `json:"weight_attrs" redis:"-"`
	PlayerManager  *PlayerManager `json:"-" redis:"-"`
	//RuleId         int            `json:"rule_id"`
}

type PlayerManager struct {
	RedisKeySeparator  string
	RedisPrefix        string
	RedisTextSeparator string
	Log                *zap.Logger //log 实例
	Redis              *util.MyRedisGo
	GameMatch          *GameMatch //父类
	Rule               *Rule      //父类
	Err                *util.ErrMsg
	CloseChan          chan int
	prefix             string
}

func NewPlayerManager(rule *Rule) *PlayerManager {
	playerManager := new(PlayerManager)

	playerManager.Rule = rule
	playerManager.Redis = rule.RuleManager.Option.GameMatch.Option.Redis
	playerManager.RedisTextSeparator = rule.RuleManager.Option.GameMatch.Option.RedisTextSeparator
	playerManager.RedisKeySeparator = rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
	playerManager.Log = rule.RuleManager.Option.GameMatch.Option.Log
	playerManager.Err = rule.RuleManager.Option.GameMatch.Err
	playerManager.CloseChan = make(chan int)
	playerManager.prefix = "PlayerManager"
	return playerManager
}

func (playerManager *PlayerManager) createEmptyPlayer() Player {
	newPlayer := Player{
		Id:          0,
		Status:      service.GAME_MATCH_PLAYER_STATUS_INIT,
		SignTimeout: 0,
		GroupId:     0,
		ATime:       util.GetNowTimeSecondToInt(),
		UTime:       util.GetNowTimeSecondToInt(),
		//RuleId:      playerManager.Rule.Id,
	}
	return newPlayer
}

func (playerManager *PlayerManager) UpStatus(pid int, status int, SuccessTimeout int, redisFDConn redis.Conn) {
	util.MyPrint("playerManager UpStatus , pid:", pid, " status:", status, " SuccessTimeout:", SuccessTimeout)
	now := util.GetNowTimeSecondToInt()
	player, _ := playerManager.GetById(pid)
	player.Status = status
	player.UTime = now
	player.SuccessTimeout = SuccessTimeout

	key := playerManager.getRedisPrefixKeyByPid(player.Id)
	//res, err := playerManager.Redis.Send(redisFDConn, "HMSET",
	//	redis.Args{}.Add(key).Add("Status").Add(status).Add("SuccessTimeout").Add(SuccessTimeout).Add("UTime").Add(now))
	res, err := playerManager.Redis.Send(redisFDConn, "HMSET", redis.Args{}.Add(key).AddFlat(&player)...)
	util.MyPrint(res, err)
}

//设置一个redis hash 结构茶杯，如不存在，则创建一个新的
func (playerManager *PlayerManager) Create(player Player, redisConnFD redis.Conn) (bool, error) {
	//playerManager.delOneRulePlayer(redisConnFD, player.Id)
	//playerManager.addOneRulePlayer(redisConnFD, player.Id)
	player.UTime = util.GetNowTimeSecondToInt()

	key := playerManager.getRedisPrefixKeyByPid(player.Id)
	playerManager.Redis.Send(redisConnFD, "HMSET", redis.Args{}.Add(key).AddFlat(&player)...)
	//res,err  := playerManager.Redis.RedisDo("HMSET",redis.Args{}.Add(key).AddFlat(&playerStatusElement)...)
	playerBytes, _ := json.Marshal(&player)
	playerManager.Log.Info("playerManager Create : " + string(playerBytes))

	return true, nil
}

//redis hash 结构
func (playerManager *PlayerManager) getRedisPrefix() string {
	return playerManager.Rule.GetCommRedisKeyByModuleRuleId(playerManager.prefix, playerManager.Rule.Id) + "player_"
}

//redis hash 结构
func (playerManager *PlayerManager) getRedisPrefixKeyByPid(pid int) string {
	return playerManager.getRedisPrefix() + strconv.Itoa(pid)
}

func (playerManager *PlayerManager) delOneRulePlayer(redisConn redis.Conn, playerId int) {
	key := playerManager.getRedisPrefixKeyByPid(playerId)
	playerManager.Redis.Send(redisConn, "del", redis.Args{}.Add(key).Add(playerId)...)
	//res,err := playerManager.Redis.RedisDo("zrem",redis.Args{}.Add(key).Add(playerId)...)
	//playerManager.Log.Warn("playerManager delOneRulePlayer")
	//util.MyPrint("delOneRulePlayer:", key, res, err)
}

func (playerManager *PlayerManager) getOneRuleAllPlayer() (list []string, err error) {
	key := playerManager.getRedisPrefix()
	res, err := redis.Strings(playerManager.Redis.RedisDo("keys", key+"*"))
	if err != nil {
		return list, errors.New("redis keys err :" + err.Error())
	}

	if len(res) <= 0 {
		playerManager.Log.Warn(" keys is null,no need del...")
		return
	}
	for _, v := range res {
		arr := strings.Split(v, playerManager.RedisKeySeparator)
		pidStr := arr[len(arr)-1]
		list = append(list, pidStr)

	}
	//res, err := redis.Strings(playerManager.Redis.RedisDo("ZRANGEBYSCORE", redis.Args{}.Add(key).Add("-inf").Add("+inf")...))
	util.MyPrint("getOneRuleAllPlayer:", res, err)
	return list, nil
}

func (playerManager *PlayerManager) getOneRuleAllPlayerCnt() (int, error) {
	key := playerManager.getRedisPrefix()
	res, err := redis.Strings(playerManager.Redis.RedisDo("keys", key+"*"))
	if err != nil {
		return 0, errors.New("redis keys err :" + err.Error())
	}

	if len(res) <= 0 {
		playerManager.Log.Warn(" keys is null,no need ...")
		return 0, errors.New(" keys is null,no need ...")
	}
	return len(res), nil
}

//根据PID 获取一个玩家的状态信息
func (playerManager *PlayerManager) GetById(playerId int) (player Player, isEmpty int) {
	//var playerStatusElement PlayerStatusElement
	key := playerManager.getRedisPrefixKeyByPid(playerId)
	res, err := redis.Values(playerManager.Redis.RedisDo("HGETALL", key))
	if err != nil {
		return player, 1
	}

	if res == nil || len(res) <= 0 {
		return player, 1
	}
	//playerStatusElement := &PlayerStatusElement{}
	if err := redis.ScanStruct(res, &player); err != nil {
		return player, 1
	}
	//playerStatusElement =  playerStatus.strToStruct(res)
	return player, 0
}

//tmp process
func (playerManager *PlayerManager) delOneById(redisConn redis.Conn, playerId int) {
	playerManager.Log.Warn("playerManager delOneById:" + strconv.Itoa(playerId))
	_, isEmpty := playerManager.GetById(playerId)
	if isEmpty == 1 {
		playerManager.Log.Warn(" playerManager getById is empty , id:" + strconv.Itoa(playerId))
		return
	}
	//key := playerManager.getRedisPrefixKeyByPid(playerId)
	//res, _ := playerManager.Redis.RedisDo("del", key)
	//util.MyPrint("playerStatus delOneById , id : "+strconv.Itoa(playerId)+" , rs : ", res)

	playerManager.delOneRulePlayer(redisConn, playerId)
	////启用事务后，这里先做 个补救
	//playerManager.Redis.RedisDo("ping")
}

//检查报名超时
func (playerManager *PlayerManager) checkSignTimeout(player Player) (isTimeout bool) {
	now := util.GetNowTimeSecondToInt()
	if now > player.SignTimeout {
		return true
	}
	return false
}

//删除所有玩家状态值
func (playerManager *PlayerManager) delAllPlayers() {
	playerManager.Log.Warn("delAllPlayers ")
	key := playerManager.getRedisPrefix()
	keys := key + "*"
	playerManager.Redis.RedisDelAllByPrefix(keys)
}

func (playerManager *PlayerManager) TestRedisKey() {
	//time.Sleep(time.Second * 3)
	redisKey := playerManager.getRedisPrefix()
	util.MyPrint("playerManager test :", redisKey)

	redisKey = playerManager.getRedisPrefixKeyByPid(1)
	util.MyPrint("playerManager test :", redisKey)

	//key := "gm_PlayerManager_1_player_1"
	//type Aaa struct {
	//	A  string         `json:"a" redis:"a"`
	//	B  string         `json:"b" redis:"b"`
	//	Cc map[string]int `json:"cc" redis:"-"`
	//}
	//
	//dd := make(map[string]int)
	//aaa := Aaa{
	//	A:  "111",
	//	Cc: dd,
	//}
	//
	//player := Player{
	//	Id:      1,
	//	Status:  service.GAME_MATCH_RULE_STATUS_INIT,
	//	GroupId: 0,
	//	Weight:  1,
	//	//WeightAttrs:    WeightAttrs,
	//	ATime:          0,
	//	UTime:          util.GetNowTimeSecondToInt(),
	//	SignTimeout:    1,
	//	SuccessTimeout: 0,
	//}

	//redisConnFD := playerManager.Redis.GetNewConnFromPool()
	//playerManager.Redis.Multi(redisConnFD)
	//res, err := playerManager.Redis.RedisDo("HMSET", redis.Args{}.Add(key).AddFlat(&aaa)...)
	//util.MyPrint(key)
	//util.MyPrint("playerStatus setInfo : ", res, err)
	//time.Sleep(time.Second * 1)
	//util.ExitPrint(33)
}

//func (playerManager *PlayerManager) strToStruct(str string) Player {
//	strArr := strings.Split(str, service.Separation)
//	playerStatusElement := Player{
//		Id:     util.Atoi(strArr[0]),
//		Status: util.Atoi(strArr[1]),
//		//RuleId: util.Atoi(strArr[2]),
//		//Weight			:zlib.Atoi(strArr[0]),
//		GroupId:        util.Atoi(strArr[4]),
//		ATime:          util.Atoi(strArr[5]),
//		UTime:          util.Atoi(strArr[6]),
//		SignTimeout:    util.Atoi(strArr[7]),
//		SuccessTimeout: util.Atoi(strArr[8]),
//	}
//
//	return playerStatusElement
//}
//func (playerManager *PlayerManager) getAllPlayers() (list map[int]Player, err error) {
//	playerManager.Log.Warn("getAllPlayers ")
//	key := playerManager.getRedisPrefixKey()
//	keys := key + "*"
//
//	res, err := redis.Strings(playerManager.Redis.RedisDo("keys", keys))
//	if err != nil {
//		util.ExitPrint("redis keys err :", err.Error())
//	}
//	playerManager.Log.Debug("all element will num :" + strconv.Itoa(len(res)))
//	if len(res) <= 0 {
//		playerManager.Log.Warn(" keys is null,no need del...")
//		return list, err
//	}
//	list = make(map[int]Player)
//	for _, p_key := range res {
//		//onePlayerRedis,err := playerManager.Redis.RedisDo("get",v)
//		//mylog.Debug("get one ",v , " ,  onePlayerRedis ",onePlayerRedis , " , err : ",err)
//		if strings.Index(p_key, "rule") != -1 {
//			continue
//		}
//		res, err := redis.Values(playerManager.Redis.RedisDo("HGETALL", p_key))
//		if err != nil {
//			util.ExitPrint("get oneplay error")
//		}
//		player := Player{}
//		if err := redis.ScanStruct(res, &player); err != nil {
//			util.ExitPrint("get oneplay error 2", p_key, err)
//		}
//		list[player.Id] = player
//	}
//	return list, err
//}
//
//func(playerManager *PlayerManager)  delOne(playerStatusElement PlayerStatusElement){
//	key := playerStatus.getRedisStatusPrefixByPid(playerStatusElement.PlayerId)
//	res,_ := playerManager.Redis.RedisDo("del",key)
//	mylog.Notice("playerStatus delOne , id : ",playerStatusElement.PlayerId , " rs : ",res)
//}

//
////索引，一个rule里包含的玩家信息，主要用于批量删除，也可用做一个rule的当前所有玩家列表
//func (playerManager *PlayerManager) addOneRulePlayer(redisConn redis.Conn, playerId int) {
//	key := playerManager.getRedisPrefixKeyByPid(playerId)
//	res, err := playerManager.Redis.Send(redisConn, "zadd", redis.Args{}.Add(key).Add(0).Add(playerId)...)
//	//res,err := playerManager.Redis.RedisDo("zadd",redis.Args{}.Add(key).Add(0).Add(playerId)...)
//	util.MyPrint("addRulePlayer:", res, err)
//}
