package gamematch

//游戏匹配
//强依赖redis(1. 持久化  2. 集合运算)
import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/util"
)

type GameMatchOption struct {
	Log   *zap.Logger     //log 实例
	Redis *util.MyRedisGo //redis 实例
	Gorm  *gorm.DB        //mysql 实例
	//Service            *util.Service          //服务 实例
	Metrics            *util.MyMetrics        //统计 实例
	ServiceDiscovery   *util.ServiceDiscovery //服务发现 实例
	StaticPath         string                 //静态文件公共目录
	RuleDataSourceType int                    //rule的数据来源类型
	RedisPrefix        string                 //redis公共的前缀，主要是怕key重复
	RedisTextSeparator string                 //结构体不能直接存到redis中，得手动分隔存进去。不存JSON是因为浪费空间
	RedisKeySeparator  string                 //redis key 的分隔符号
	ProjectId          int
	//Etcd             *util.MyEtcd
}

type GameMatch struct {
	Option                 GameMatchOption //初始化 配置参数
	Err                    *util.ErrMsg    //错误处理类
	RuleManager            *RuleManager    //管理一条 rule 的所有控制，如：报名 匹配 推送 等等
	PlayerManager          *PlayerManager  //管理所有用户的状态信息等
	prefix                 string          //前缀字符串
	LoopSleepTime          int             //有些死循环需要睡眠的场景,毫秒
	FormulaFirst           string          //游戏匹配-计算权重公式-前缀
	FormulaEnd             string          //游戏匹配-计算权重公式-后缀
	RuleTeamMaxPeople      int             //一个小组允许最大人数
	RulePersonConditionMax int             //N人组团，最大人数
	RuleMatchTimeoutMax    int             //报名，最大超时时间
	RuleMatchTimeoutMin    int             //报名，最小时间
	RuleSuccessTimeoutMax  int             //匹配成功后，最大超时时间
	RuleSuccessTimeoutMin  int             //匹配成功后，最短超时时间
	WeightMaxValue         int
	//RuleTeamVSPersonMax    int             //组队互相PK，每个队最多人数
}

func NewGameMatch(option GameMatchOption) (*GameMatch, error) {
	option.Log.Info("NewGameMatch : ")

	util.MyPrint(option.Redis)

	gameMatch := new(GameMatch)
	gameMatch.Option = option

	gameMatch.prefix = "gameMatch"
	gameMatch.LoopSleepTime = 100
	gameMatch.FormulaFirst = "<"           //游戏匹配-计算权重公式-前缀
	gameMatch.FormulaEnd = ">"             //游戏匹配-计算权重公式-后缀
	gameMatch.RuleTeamMaxPeople = 5        //一个小组允许最大人数
	gameMatch.RulePersonConditionMax = 100 //N人组团，最大人数
	gameMatch.RuleMatchTimeoutMax = 400    //报名，最大超时时间
	gameMatch.RuleMatchTimeoutMin = 3      //报名，最小时间
	gameMatch.RuleSuccessTimeoutMax = 600  //匹配成功后，最大超时时间
	gameMatch.RuleSuccessTimeoutMin = 10   //匹配成功后，最短超时时间
	gameMatch.WeightMaxValue = 100
	//gameMatch.RuleTeamVSPersonMax = 10     //组队互相PK，每个队最多人数

	lang, err := util.NewErrMsg(option.Log, option.StaticPath+"/data/game_match_cn.lang")
	if err != nil {
		util.ExitPrint(err)
	}
	gameMatch.Err = lang
	ruleManagerOption := RuleManagerOption{
		Gorm:      option.Gorm,
		GameMatch: gameMatch,
	}

	gameMatch.RuleManager, err = NewRuleManager(ruleManagerOption)
	if err != nil {
		util.MyPrint(err)
		return gameMatch, err
	}
	gameMatch.PlayerManager = NewPlayerManager(gameMatch)
	gameMatch.PlayerManager.TestRedisKey()

	gameMatch.RuleManager.StartupAll()

	return gameMatch, nil
}

//退出
func (gameMatch *GameMatch) Quit(source int) {
	gameMatch.RuleManager.Quit()
}

//启动后台守护-协程
func (gameMatch *GameMatch) Startup() {
	gameMatch.RuleManager.StartupAll()
}

//
////睡眠 - 协程
//func (gameMatch *GameMatch) mySleepSecond(second time.Duration, msg string) {
//	gameMatch.Option.Log.Info(msg + " sleep second " + strconv.Itoa(int(second)))
//	time.Sleep(second * time.Second)
//}

//删除 全部 redis 数据，这个是方便测试，线上业务不能用
func (gameMatch *GameMatch) DelRedisData() {
	gameMatch.Option.Log.Warn(gameMatch.prefix + " action :  DelRedisData")
	keys := gameMatch.Option.RedisPrefix + "*"
	gameMatch.Option.Redis.RedisDelAllByPrefix(keys)
}

//删除 一条 rule的 redis 数据
func (gameMatch *GameMatch) DelOneRuleRedisDataById(ruleId int) {
	gameMatch.Option.Log.Warn(gameMatch.prefix + " action :  DelOneRuleRedisDataById , " + strconv.Itoa(ruleId))
	gameMatch.RuleManager.delOneRuleRedisData(ruleId)
}
