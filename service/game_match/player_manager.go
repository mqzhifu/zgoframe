package gamematch

import (
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
	RuleId         int            `json:"rule_id"`
	Weight         float32        `json:"weight"`
	WeightAttrs    map[string]int `json:"-"`
	GroupId        int            `json:"group_id"`
	ATime          int            `json:"a_time"`
	UTime          int            `json:"u_time"`
	SignTimeout    int            `json:"sign_timeout"`
	SuccessTimeout int            `json:"success_timeout"`
}

type PlayerManager struct {
	RedisKeySeparator  string
	RedisPrefix        string
	RedisTextSeparator string
	Log                *zap.Logger //log 实例
	Redis              *util.MyRedisGo
	GameMatch          *GameMatch //父类
	prefix             string
}

func NewPlayerManager(gameMatch *GameMatch) *PlayerManager {
	gameMatch.Option.Log.Info("NewPlayerManager")
	playerManager := new(PlayerManager)

	playerManager.Redis = gameMatch.Option.Redis
	playerManager.RedisTextSeparator = gameMatch.Option.RedisTextSeparator
	playerManager.RedisKeySeparator = gameMatch.Option.RedisKeySeparator
	playerManager.Log = gameMatch.Option.Log
	playerManager.RedisPrefix = gameMatch.Option.RedisPrefix
	playerManager.prefix = "PlayerManager"
	return playerManager
}

func (playerManager *PlayerManager) create() Player {
	newPlayer := Player{
		Id:          0,
		Status:      service.GAME_MATCH_PLAYER_STATUS_INIT,
		RuleId:      0,
		SignTimeout: 0,
		GroupId:     0,
	}
	return newPlayer
}

func (playerManager *PlayerManager) getRedisPrefixKey() string {
	return playerManager.RedisPrefix + playerManager.RedisKeySeparator + "player" + playerManager.RedisKeySeparator
}

func (playerManager *PlayerManager) getRedisPrefixByPid(playerId int) string {
	return playerManager.getRedisPrefixKey() + strconv.Itoa(playerId)
}

func (playerManager *PlayerManager) getRulePlayerPrefixByRuleId(ruleId int) string {
	return playerManager.getRedisPrefixKey() + "rule" + playerManager.RedisKeySeparator + strconv.Itoa(ruleId)
}

//索引，一个rule里包含的玩家信息，主要用于批量删除，也可用做一个rule的当前所有玩家列表
func (playerManager *PlayerManager) addOneRulePlayer(redisConn redis.Conn, playerId int, ruleId int) {
	key := playerManager.getRulePlayerPrefixByRuleId(ruleId)
	res, err := playerManager.Redis.Send(redisConn, "zadd", redis.Args{}.Add(key).Add(0).Add(playerId)...)
	//res,err := playerManager.Redis.RedisDo("zadd",redis.Args{}.Add(key).Add(0).Add(playerId)...)
	util.MyPrint("addRulePlayer:", res, err)
}

func (playerManager *PlayerManager) delOneRulePlayer(redisConn redis.Conn, playerId int, ruleId int) {
	key := playerManager.getRulePlayerPrefixByRuleId(ruleId)
	res, err := playerManager.Redis.Send(redisConn, "zrem", redis.Args{}.Add(key).Add(playerId)...)
	//res,err := playerManager.Redis.RedisDo("zrem",redis.Args{}.Add(key).Add(playerId)...)
	util.MyPrint("delOneRulePlayer:", res, err)
}
func (playerManager *PlayerManager) getOneRuleAllPlayer(ruleId int) []string {
	key := playerManager.getRulePlayerPrefixByRuleId(ruleId)
	//ZRANGEBYSCORE salary -inf +inf WITHSCORES
	res, err := redis.Strings(playerManager.Redis.RedisDo("ZRANGEBYSCORE", redis.Args{}.Add(key).Add("-inf").Add("+inf")...))
	util.MyPrint("getOneRuleAllPlayer:", res, err)
	return res
}

func (playerManager *PlayerManager) getOneRuleAllPlayerCnt(ruleId int) int {
	key := playerManager.getRedisPrefixByPid(ruleId)
	//ZRANGEBYSCORE salary -inf +inf WITHSCORES
	res, err := redis.Int(playerManager.Redis.RedisDo("ZCOUNT", redis.Args{}.Add(key).Add("-inf").Add("+inf")...))
	util.MyPrint("getOneRuleAllPlayer:", res, err)
	return res
}

//根据PID 获取一个玩家的状态信息
func (playerManager *PlayerManager) GetById(playerId int) (player Player, isEmpty int) {
	//var playerStatusElement PlayerStatusElement
	key := playerManager.getRedisPrefixByPid(playerId)
	res, err := redis.Values(playerManager.Redis.RedisDo("HGETALL", key))
	if err != nil {
		return player, 1
	}
	//playerStatusElement := &PlayerStatusElement{}
	if err := redis.ScanStruct(res, &player); err != nil {
		return player, 1
	}
	//playerStatusElement =  playerStatus.strToStruct(res)
	return player, 0
}

func (playerManager *PlayerManager) upInfo(player Player, redisConnFD redis.Conn) (bool, error) {

	playerManager.delOneRulePlayer(redisConnFD, player.Id, player.RuleId)
	playerManager.addOneRulePlayer(redisConnFD, player.Id, player.RuleId)

	player.UTime = util.GetNowTimeSecondToInt()
	playerManager.setInfo(redisConnFD, player)

	return true, nil
}

func (playerManager *PlayerManager) strToStruct(str string) Player {
	strArr := strings.Split(str, service.Separation)
	playerStatusElement := Player{
		Id:     util.Atoi(strArr[0]),
		Status: util.Atoi(strArr[1]),
		RuleId: util.Atoi(strArr[2]),
		//Weight			:zlib.Atoi(strArr[0]),
		GroupId:        util.Atoi(strArr[4]),
		ATime:          util.Atoi(strArr[5]),
		UTime:          util.Atoi(strArr[6]),
		SignTimeout:    util.Atoi(strArr[7]),
		SuccessTimeout: util.Atoi(strArr[8]),
	}

	return playerStatusElement
}

func (playerManager *PlayerManager) setInfo(conn redis.Conn, player Player) {
	key := playerManager.getRedisPrefixByPid(player.Id)
	res, err := playerManager.Redis.Send(conn, "HMSET", redis.Args{}.Add(key).AddFlat(&player)...)
	//res,err  := playerManager.Redis.RedisDo("HMSET",redis.Args{}.Add(key).AddFlat(&playerStatusElement)...)
	util.MyPrint("playerStatus setInfo : ", player, res, err)
}

//tmp process
func (playerManager *PlayerManager) delOneById(redisConn redis.Conn, playerId int) {
	playerStatusElement, isEmpty := playerManager.GetById(playerId)
	if isEmpty == 1 {
		playerManager.Log.Error(" getByid is empty!!!")
		return
	}
	key := playerManager.getRedisPrefixByPid(playerId)
	res, _ := playerManager.Redis.RedisDo("del", key)
	util.MyPrint("playerStatus delOneById , id : "+strconv.Itoa(playerId)+" , rs : ", res)

	playerManager.delOneRulePlayer(redisConn, playerId, playerStatusElement.RuleId)
	//启用事务后，这里先做 个补救
	playerManager.Redis.RedisDo("ping")
}

//删除所有玩家状态值
func (playerManager *PlayerManager) delAllPlayers() {
	playerManager.Log.Warn("delAllPlayers ")
	key := playerManager.getRedisPrefixKey()
	keys := key + "*"
	playerManager.Redis.RedisDelAllByPrefix(keys)
}

//检查报名超时
func (playerManager *PlayerManager) checkSignTimeout(player Player) (isTimeout bool) {
	now := util.GetNowTimeSecondToInt()
	if now > player.SignTimeout {
		return true
	}
	return false
}

func (playerManager *PlayerManager) getAllPlayers() (list map[int]Player, err error) {
	playerManager.Log.Warn("getAllPlayers ")
	key := playerManager.getRedisPrefixKey()
	keys := key + "*"

	res, err := redis.Strings(playerManager.Redis.RedisDo("keys", keys))
	if err != nil {
		util.ExitPrint("redis keys err :", err.Error())
	}
	playerManager.Log.Debug("all element will num :" + strconv.Itoa(len(res)))
	if len(res) <= 0 {
		playerManager.Log.Warn(" keys is null,no need del...")
		return list, err
	}
	list = make(map[int]Player)
	for _, p_key := range res {
		//onePlayerRedis,err := playerManager.Redis.RedisDo("get",v)
		//mylog.Debug("get one ",v , " ,  onePlayerRedis ",onePlayerRedis , " , err : ",err)
		if strings.Index(p_key, "rule") != -1 {
			continue
		}
		res, err := redis.Values(playerManager.Redis.RedisDo("HGETALL", p_key))
		if err != nil {
			util.ExitPrint("get oneplay error")
		}
		player := Player{}
		if err := redis.ScanStruct(res, &player); err != nil {
			util.ExitPrint("get oneplay error 2", p_key, err)
		}
		list[player.Id] = player
	}
	return list, err
}

func (playerManager *PlayerManager) TestRedisKey() {
	redisKey := playerManager.getRulePlayerPrefixByRuleId(1)
	util.MyPrint("playerManager test :", redisKey)

	redisKey = playerManager.getRedisPrefixByPid(1)
	util.MyPrint("playerManager test :", redisKey)

}

//func(playerManager *PlayerManager)  delOne(playerStatusElement PlayerStatusElement){
//	key := playerStatus.getRedisStatusPrefixByPid(playerStatusElement.PlayerId)
//	res,_ := playerManager.Redis.RedisDo("del",key)
//	mylog.Notice("playerStatus delOne , id : ",playerStatusElement.PlayerId , " rs : ",res)
//}
