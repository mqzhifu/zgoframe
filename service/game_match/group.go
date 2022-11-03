package gamematch

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"zgoframe/util"
)

type Group struct {
	Id             int      `json:"id"`
	Type           int      `json:"type"`            //报名跟报名成功会各创建一条group记录，1：报名，2匹配成功
	Person         int      `json:"person"`          //小组总人数
	Weight         float32  `json:"weight"`          //小组权重
	MatchTimes     int      `json:"match_times"`     //已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃，不过没用到
	SignTimeout    int      `json:"sign_timeout"`    //多少秒后无人来取，即超时，更新用户状态，删除数据
	SuccessTimeout int      `json:"success_timeout"` //匹配成功后，无人来取，超时
	SignTime       int      `json:"sign_time"`       //报名时间
	SuccessTime    int      `json:"success_time"`    //匹配成功时间
	Players        []Player `json:"players"`         //用户列表
	Addition       string   `json:"addition"`        //请求方附加属性值
	TeamId         int      `json:"team_id"`         //组队互相PK的时候，得成两个队伍
	OutGroupId     int      `json:"out_group_id"`    //报名时，客户端请求时，自带的一个ID
	//CustomProp     string
	//MatchCode      string
	//LinkId         int     //关联ID，匹配成功后关联成功的那条记录的ID，正常报名用不上
}

func (gamematch *GameMatch) NewGroupStruct(rule *Rule) Group {
	group := Group{}
	//group.Id = gamematch.GetGroupIncId(rule.Id)
	group.Person = 0
	group.Weight = 0
	group.MatchTimes = 0
	group.SignTimeout = 0
	group.SuccessTimeout = 0
	group.SignTime = util.GetNowTimeSecondToInt()
	group.SuccessTime = 0
	group.Addition = ""
	group.Players = nil
	group.OutGroupId = 0
	//group.LinkId = 0
	//group.CustomProp = ""
	//group.MatchCode = ""
	return group
}

//组自增ID，因为匹配最小单位是基于组，而不是基于一个人，组ID就得做到全局唯一，很重要
func (gameMatch *GameMatch) getRedisGroupIncKey(ruleId int) string {
	return gameMatch.Option.RedisKeySeparator + "group_inc_id" + gameMatch.Option.RedisKeySeparator + strconv.Itoa(ruleId)
}

//获取并生成一个自增GROUP-ID
func (gameMatch *GameMatch) GetGroupIncId(ruleId int) int {
	key := gameMatch.getRedisGroupIncKey(ruleId)
	res, _ := redis.Int(gameMatch.Option.Redis.RedisDo("INCR", key))
	return res
}

func (gameMatch *GameMatch) GroupStrToStruct(redisStr string) Group {
	strArr := strings.Split(redisStr, gameMatch.Option.RedisTextSeparator)
	playersArr := strings.Split(strArr[8], gameMatch.Option.RedisIdSeparator)
	var players []Player
	for _, v := range playersArr {
		players = append(players, Player{Id: util.Atoi(v)})
	}

	element := Group{
		Id:             util.Atoi(strArr[0]),
		Person:         util.Atoi(strArr[1]),
		Weight:         util.StringToFloat(strArr[2]),
		MatchTimes:     util.Atoi(strArr[3]),
		SignTimeout:    util.Atoi(strArr[4]),
		SuccessTimeout: util.Atoi(strArr[5]),
		SignTime:       util.Atoi(strArr[6]),
		SuccessTime:    util.Atoi(strArr[7]),
		Players:        players,
		Addition:       strArr[9],
		TeamId:         util.Atoi(strArr[10]),
		OutGroupId:     util.Atoi(strArr[11]),
		//LinkId:         util.Atoi(strArr[1]),
		//CustomProp:     strArr[13],
		//MatchCode:      strArr[14],
	}

	return element
}
func (gameMatch *GameMatch) GetGroupPlayerIds(group Group) string {
	playersIds := ""
	for _, v := range group.Players {
		playersIds += strconv.Itoa(v.Id) + gameMatch.Option.RedisIdSeparator
	}
	playersIds = playersIds[0 : len(playersIds)-1]
	return playersIds
}

//理论上直接用 redis 的 hash 更简单，但是一些结构的元素又包含其它结构体，是2维3维的，redis就处理不了，只能用这种略笨的方法
func (gameMatch *GameMatch) GroupStructToStr(group Group) string {
	//Weight	float32	//小组权重
	//MatchTimes	int		//已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃
	playersIds := gameMatch.GetGroupPlayerIds(group)
	Weight := util.FloatToString(group.Weight, 3)

	content :=
		strconv.Itoa(group.Id) + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.Person) + gameMatch.Option.RedisTextSeparator +
			Weight + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.MatchTimes) + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.SignTimeout) + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.SuccessTimeout) + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.SignTime) + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.SuccessTime) + gameMatch.Option.RedisTextSeparator +
			playersIds + gameMatch.Option.RedisTextSeparator +
			group.Addition + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.TeamId) + gameMatch.Option.RedisTextSeparator +
			strconv.Itoa(group.OutGroupId) + gameMatch.Option.RedisTextSeparator

	return content
}
