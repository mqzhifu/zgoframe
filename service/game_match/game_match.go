package gamematch

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/util"
)

/*
	游戏匹配
	强依赖redis(1. 持久化  2. 集合运算)
*/

/*	需要优化的问题：
	1. 所有的 redis 未加失效时间，是靠其它守护协程轮询检查失效，回收过程。如果只是简单的 key 可以直接加expire ，但是有些集合 直接失效不太好，我再想想
	2. 目前团队最大仅支持最大 5 人，这个是被 组合排队 的算法难住了
	3. 目前仅支持两个队互相掑杀，2队以上，目前还不支持，程序中的代码也是写死的
	4. 成功后超时 与  匹配成功重试  冲突
	5. groupId目前是使用外部的ID，是否考虑换成内部
	6. 查看当时进程有多少个协程，是否有metrics的方式，且最可以UI可视化
	7. 各协程之间是否需要加锁
	8. pnic 异常机制处理
	9. 压测 redis etcd log 匹配 http 报名
	10. 目前：每个rule，都有一组守护协程，和一组redis key 存数据，所以rule之间不需要加锁了
		但是单个ruleId,不加锁，可能出现的问题：
		1、报名中，有取消指令发出，但是匹配协程依然还在计算，并且计算成功了
		2、匹配中的玩家已超时，但是超时协程未执行，匹配协程先执行了，匹配成功.....
		3、匹配成功已超时，但是超时协程未执行，PUSH协程先执行了....
		4、超时的2个协程挂了，但是 报名 PUSH 匹配 协程均是正常，那匹配依然会成功，PUSH依然还是会推送

		这里有2个维度的问题：
		1、如何保证所有守护协程是正常的？进程的健康可以由外部shell控制，协程呢？
		2、如果保证上面一条是正常的，那核心点就是匹配协程了
	)
	11. 负载均衡，
		1. 因为每个rule有单独的redis数据，互不影响，那么：每台机器负载监听具体的ruleId
		2. 单个ruleId，每台机器均可以启动，只要redis换个HOST即可
	)
*/

type GameMatchOption struct {
	Log   *zap.Logger     //log 实例
	Redis *util.MyRedisGo //redis 实例
	Gorm  *gorm.DB        //mysql 实例
	//Service            *util.Service          //服务 实例
	Metrics                *util.MyMetrics        //统计 实例
	ServiceDiscovery       *util.ServiceDiscovery //服务发现 实例
	StaticPath             string                 //静态文件公共目录
	RuleDataSourceType     int                    //rule的数据来源类型
	RedisPrefix            string                 //redis公共的前缀，主要是怕key重复
	RedisTextSeparator     string                 //结构体不能直接存到redis中，得手动分隔存进去。不存JSON是因为浪费空间
	RedisKeySeparator      string                 //redis key 的分隔符号
	RedisIdSeparator       string                 //一个大的结构被转化成字符串后，有些元素是复合结构，如IDS。得有个分隔符
	RedisPayloadSeparation string                 //也是redis内容的分隔符，但它包含了其它的内容，与其它内容的分隔符冲突，所以得新起一个
	ProjectId              int
	PersistenceType        int //持久化类型，0关闭
	//Etcd             *util.MyEtcd
}

type GameMatch struct {
	Option                 GameMatchOption //初始化 配置参数
	Err                    *util.ErrMsg    //错误处理类
	RuleManager            *RuleManager    //管理一条 rule 的所有控制，如：报名 匹配 推送 等等
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
	WeightMaxValue         int             //玩家权重上限值
	//RuleTeamVSPersonMax    int             //组队互相PK，每个队最多人数
	//PlayerManager          *PlayerManager  //管理所有用户的状态信息等
}

func NewGameMatch(option GameMatchOption) (*GameMatch, error) {
	option.Log.Info("NewGameMatch : ")

	util.MyPrint(option.Redis)

	gameMatch := new(GameMatch)
	gameMatch.Option = option

	gameMatch.prefix = "gameMatch"
	gameMatch.LoopSleepTime = 200
	gameMatch.FormulaFirst = "<"           //游戏匹配-计算权重公式-前缀
	gameMatch.FormulaEnd = ">"             //游戏匹配-计算权重公式-后缀
	gameMatch.RuleTeamMaxPeople = 5        //一个小组允许最大人数
	gameMatch.RulePersonConditionMax = 100 //N人组团，最大人数
	gameMatch.RuleMatchTimeoutMax = 400    //报名，最大超时时间
	gameMatch.RuleMatchTimeoutMin = 3      //报名，最小时间
	gameMatch.RuleSuccessTimeoutMax = 600  //匹配成功后，最大超时时间
	gameMatch.RuleSuccessTimeoutMin = 10   //匹配成功后，最短超时时间
	gameMatch.WeightMaxValue = 100
	gameMatch.Option.PersistenceType = service.PERSISTENCE_TYPE_MYSQL
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
	//gameMatch.PlayerManager = NewPlayerManager(gameMatch)
	//gameMatch.PlayerManager.TestRedisKey()

	gameMatch.RuleManager.StartupAll()

	return gameMatch, nil
}

//退出
func (gameMatch *GameMatch) Quit(source int) {
	gameMatch.RuleManager.Quit()
}

//持久化数据 - 组
func (gameMatch *GameMatch) PersistenceRecordGroup(group Group, ruleId int) {
	if gameMatch.Option.PersistenceType == service.PERSISTENCE_TYPE_MYSQL {
		pids := ""
		for _, v := range group.Players {
			pids += strconv.Itoa(v.Id) + ","
		}
		pids = pids[0 : len(pids)-1]
		gameMatchSign := model.GameMatchGroup{
			RuleId:         ruleId,
			SelfId:         group.Id,
			Type:           group.Type,
			Person:         group.Person,
			Weight:         util.FloatToString(group.Weight, 2),
			MatchTimes:     group.MatchTimes,
			SignTimeout:    group.SignTimeout,
			SuccessTimeout: group.SuccessTimeout,
			SignTime:       group.SignTime,
			SuccessTime:    group.SuccessTime,
			PlayerIds:      pids,
			Addition:       group.Addition,
			TeamId:         group.TeamId,
			OutGroupId:     group.OutGroupId,
		}
		gameMatch.Option.Gorm.Create(&gameMatchSign)
	}
}

//持久化数据 - 匹配成功结果
func (gameMatch *GameMatch) PersistenceRecordSuccessResult(result Result, ruleId int) {
	if gameMatch.Option.PersistenceType == service.PERSISTENCE_TYPE_MYSQL {
		gameMatchSuccess := model.GameMatchSuccess{
			RuleId:     ruleId,
			ATime:      result.ATime,
			Timeout:    result.Timeout,
			Teams:      util.ArrCoverStr(result.Teams, ","),
			PlayerIds:  util.ArrCoverStr(result.PlayerIds, ","),
			PushSelfId: result.PushId,
			GroupIds:   util.ArrCoverStr(result.GroupIds, ","),
		}
		gameMatch.Option.Gorm.Create(&gameMatchSuccess)
	}
}

//持久化数据 - 组
func (gameMatch *GameMatch) PersistenceRecordSuccessPush(pushElement PushElement, ruleId int) {
	if gameMatch.Option.PersistenceType == service.PERSISTENCE_TYPE_MYSQL {
		gameMatchPush := model.GameMatchPush{
			RuleId:   ruleId,
			ATime:    pushElement.ATime,
			SelfId:   pushElement.Id,
			LinkId:   pushElement.LinkId,
			Status:   pushElement.Status,
			Times:    pushElement.Times,
			Category: pushElement.Category,
			Payload:  pushElement.Payload,
		}
		gameMatch.Option.Gorm.Create(&gameMatchPush)
	}
}

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

//启动后台守护-协程
//func (gameMatch *GameMatch) Startup() {
//
//}
//
////睡眠 - 协程
//func (gameMatch *GameMatch) mySleepSecond(second time.Duration, msg string) {
//	gameMatch.Option.Log.Info(msg + " sleep second " + strconv.Itoa(int(second)))
//	time.Sleep(second * time.Second)
//}
