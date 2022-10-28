package gamematch

//
//import (
//	"github.com/gomodule/redigo/redis"
//	"strconv"
//	"strings"
//	"zgoframe/service"
//	"zgoframe/util"
//)
//
//type Group struct {
//	Id             int
//	LinkId         int     //关联ID，匹配成功后关联成功的那条记录的ID，正常报名用不上
//	Person         int     //小组人数
//	Weight         float32 //小组权重
//	MatchTimes     int     //已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃，不过没用到
//	SignTimeout    int     //多少秒后无人来取，即超时，更新用户状态，删除数据
//	SuccessTimeout int
//	SignTime       int //报名时间
//	SuccessTime    int //匹配成功时间
//	Players        []Player
//	Addition       string //请求方附加属性值
//	TeamId         int    //组队互相PK的时候，得成两个队伍
//	OutGroupId     int
//	CustomProp     string
//	MatchCode      string
//}
//
//func (gamematch *Gamematch) NewGroupStruct(rule Rule) Group {
//	group := Group{}
//	group.Id = gamematch.GetGroupIncId(rule.Id)
//	group.LinkId = 0
//	group.Person = 0
//	group.Weight = 0
//	group.MatchTimes = 0
//	group.SignTimeout = util.GetNowTimeSecondToInt()
//	group.SuccessTimeout = 0
//	group.SignTime = 0
//	group.SuccessTime = 0
//	group.Addition = ""
//	group.Players = nil
//	group.OutGroupId = 0
//	group.CustomProp = ""
//	group.MatchCode = ""
//	return group
//}
//
////组自增ID，因为匹配最小单位是基于组，而不是基于一个人，组ID就得做到全局唯一，很重要
//func (gamematch *Gamematch) getRedisGroupIncKey(ruleId int) string {
//	signClass := gamematch.GetContainerSignByRuleId(ruleId)
//	return signClass.getRedisCatePrefixKey() + service.RedisSeparation + "group_inc_id"
//}
//
////获取并生成一个自增GROUP-ID
//func (gamematch *Gamematch) GetGroupIncId(ruleId int) int {
//	key := gamematch.getRedisGroupIncKey(ruleId)
//	res, _ := redis.Int(myredis.RedisDo("INCR", key))
//	return res
//}
//
//func GroupStrToStruct(redisStr string) Group {
//	strArr := strings.Split(redisStr, service.Separation)
//	playersArr := strings.Split(strArr[9], service.IdsSeparation)
//	var players []Player
//	for _, v := range playersArr {
//		players = append(players, Player{Id: util.Atoi(v)})
//	}
//
//	element := Group{
//		Id:             util.Atoi(strArr[0]),
//		LinkId:         util.Atoi(strArr[1]),
//		Person:         util.Atoi(strArr[2]),
//		Weight:         util.StringToFloat(strArr[3]),
//		MatchTimes:     util.Atoi(strArr[4]),
//		SignTimeout:    util.Atoi(strArr[5]),
//		SuccessTimeout: util.Atoi(strArr[6]),
//		SignTime:       util.Atoi(strArr[7]),
//		SuccessTime:    util.Atoi(strArr[8]),
//		Players:        players,
//		Addition:       strArr[10],
//		TeamId:         util.Atoi(strArr[11]),
//		OutGroupId:     util.Atoi(strArr[12]),
//		CustomProp:     strArr[13],
//		MatchCode:      strArr[14],
//	}
//
//	return element
//}
//
//func GroupStructToStr(group Group) string {
//	//Weight	float32	//小组权重
//	//MatchTimes	int		//已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃
//
//	playersIds := ""
//	for _, v := range group.Players {
//		playersIds += strconv.Itoa(v.Id) + service.IdsSeparation
//	}
//	playersIds = playersIds[0 : len(playersIds)-1]
//	Weight := util.FloatToString(group.Weight, 3)
//
//	content :=
//		strconv.Itoa(group.Id) + service.Separation +
//			strconv.Itoa(group.LinkId) + service.Separation +
//			strconv.Itoa(group.Person) + service.Separation +
//			Weight + service.Separation +
//			strconv.Itoa(group.MatchTimes) + service.Separation +
//			strconv.Itoa(group.SignTimeout) + service.Separation +
//			strconv.Itoa(group.SuccessTimeout) + service.Separation +
//			strconv.Itoa(group.SignTime) + service.Separation +
//			strconv.Itoa(group.SuccessTime) + service.Separation +
//			playersIds + service.Separation +
//			group.Addition + service.Separation +
//			strconv.Itoa(group.TeamId) + service.Separation +
//			strconv.Itoa(group.OutGroupId) + service.Separation +
//			group.CustomProp + service.Separation +
//			group.MatchCode
//
//	return content
//}
