//游戏-用户匹配
package gamematch

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/service/frame_sync"
	"zgoframe/util"
)

/*
	游戏匹配
	强依赖redis(1. 持久化  2. 集合运算)

	负载均衡
		1. 每个 rule 有单独的 redis 数据，互不影响(即使 redis 用一个实例也不影响  )，那么：每台机器启动进程的时候，配置好参数，监听具体的 ruleId，那么 rule 之间就可以负载。如果分机器：建议：redis 实例还是分开较好，或者 redis 分库
		2. 单个 rule 的负载，每台机器启动程序时，监听同一个 ruleId，但 redis 实例得用不同的，不然，多机器上的程序一起访问同一个 redis 实例，数据会乱。
			前端如何访问？
			1. 后端提供一个接口，前端请求接口时带上一个参数，如：地域、UID、联通移动 等，后端实时给返回一个 地址
			2. 还是一个接口，但后端报名的入口得加个代理，帮忙找到那个服务
*/

/*	需要优化的问题：
	1. 目前 N VS N 模式下，团队最大仅支持最大 5 人，这个是被 组合排队 的算法难住了
	2. 目前仅支持两个队互相掑杀，2队以上，目前还不支持，程序中的代码也是写死的
	>这两条如果从业务上看，也可以不用改，因为：主流游戏都是两个阵营互相撕杀，而且队伍中基本都是3人 或 5人
	3. rule 的数据现在改成了从mysql 中取，如果发生变化，得有个接口来接收，并更新 rule 状态
	4. 匹配成功后，超时机制 与  push-<匹配成功>-重试  冲突(概率极低，两个协程同时唤醒，一边超时了，另一个又一直没推送成功，碰巧这次重试跟超时成功碰上了)
	5. groupId目前是使用外部的ID，是否考虑换成内部
	6. UI可视化工具
	7. pnic 异常机制处理
	8. 压测 redis etcd log 匹配 http 报名

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

	)
*/

type GameMatchOption struct {
	ProjectId              int                   `json:"project_id"`
	StaticPath             string                `json:"static_path"`               //静态文件公共目录,用于读取语言包
	RuleDataSourceType     int                   `json:"rule_data_source_type"`     //rule的数据来源类型
	PersistenceType        int                   `json:"persistence_type"`          //持久化类型，0关闭
	RedisPrefix            string                `json:"redis_prefix"`              //redis公共的前缀，主要是怕key重复
	RedisTextSeparator     string                `json:"redis_text_separator"`      //结构体不能直接存到redis中，得手动分隔存进去。不存JSON是因为浪费空间
	RedisKeySeparator      string                `json:"redis_key_separator"`       //redis key 的分隔符号
	RedisIdSeparator       string                `json:"redis_id_separator"`        //一个大的结构被转化成字符串后，有些元素是复合结构，如IDS。得有个分隔符
	RedisPayloadSeparation string                `json:"redis_payload_separation"`  //也是redis内容的分隔符，但它包含了其它的内容，与其它内容的分隔符冲突，所以得新起一个
	LoopSleepTime          int                   `json:"loop_sleep_time"`           //有些死循环守护协程，需要睡眠的场景,毫秒
	FormulaFirst           string                `json:"formula_first"`             //游戏匹配-计算权重公式-前缀
	FormulaEnd             string                `json:"formula_end"`               //游戏匹配-计算权重公式-后缀
	WeightMaxValue         int                   `json:"weight_max_value"`          //玩家权重上限值
	RuleDebugShow          int                   `json:"rule_debug_show"`           //守护协程，在没有处理数据时，需要输出日志，太多，每X秒输出一次
	RuleTeamMaxPeople      int                   `json:"rule_team_max_people"`      //一个小组允许最大人数
	RulePersonConditionMax int                   `json:"rule_person_condition_max"` //N人组团，最大人数
	RuleMatchTimeoutMax    int                   `json:"rule_match_timeout_max"`    //报名，最大超时时间
	RuleMatchTimeoutMin    int                   `json:"rule_match_timeout_min"`    //报名，最小时间
	RuleSuccessTimeoutMax  int                   `json:"rule_success_timeout_max"`  //匹配成功后，最大超时时间
	RuleSuccessTimeoutMin  int                   `json:"rule_success_timeout_min"`  //匹配成功后，最短超时时间
	FrameSync              *frame_sync.FrameSync `json:"-"`                         //帧同步
	ServiceBridge          *service.Bridge
	//RequestServiceAdapter  *service.RequestServiceAdapter `json:"-"`                         //请求3方服务 适配器
	Log              *zap.Logger            `json:"-"` //log 实例
	Redis            *util.MyRedisGo        `json:"-"` //redis 实例
	Gorm             *gorm.DB               `json:"-"` //mysql 实例
	Metrics          *util.MyMetrics        `json:"-"` //统计 实例
	ServiceDiscovery *util.ServiceDiscovery `json:"-"` //服务发现 实例
	ProtoMap         *util.ProtoMap         `json:"-"`
}

type GameMatch struct {
	Option      GameMatchOption //初始化 配置参数
	Err         *util.ErrMsg    //错误处理类
	RuleManager *RuleManager    //管理一条 rule 的所有控制，如：报名 匹配 推送 等等
	prefix      string          //前缀字符串
	//RuleTeamVSPersonMax    int             //组队互相PK，每个队最多人数
}

func NewGameMatch(option GameMatchOption) (*GameMatch, error) {
	option.Log.Info("NewGameMatch : ")

	gameMatch := new(GameMatch)

	gameMatch.prefix = "gameMatch"
	option.LoopSleepTime = 200
	option.RuleDebugShow = 5
	option.FormulaFirst = "<"                               //游戏匹配-计算权重公式-前缀
	option.FormulaEnd = ">"                                 //游戏匹配-计算权重公式-后缀
	option.RuleTeamMaxPeople = 5                            //一个小组允许最大人数
	option.RulePersonConditionMax = 100                     //N人组团，最大人数
	option.RuleMatchTimeoutMax = 100                        //报名，最大超时时间
	option.RuleMatchTimeoutMin = 3                          //报名，最小时间
	option.RuleSuccessTimeoutMax = 300                      //匹配成功后，最大超时时间
	option.RuleSuccessTimeoutMin = 10                       //匹配成功后，最短超时时间
	option.WeightMaxValue = 100                             //权限最终的值，不能大于 100
	option.PersistenceType = service.PERSISTENCE_TYPE_MYSQL //数据 - 持久化
	gameMatch.Option = option
	//语言包
	lang, err := util.NewErrMsg(option.Log, option.StaticPath+"/data/game_match_cn.lang")
	if err != nil {
		util.ExitPrint(err)
	}
	gameMatch.Err = lang
	ruleManagerOption := RuleManagerOption{
		Gorm:      option.Gorm,
		GameMatch: gameMatch,
		//RequestServiceAdapter: option.RequestServiceAdapter,
		ServiceBridge: option.ServiceBridge,
	}

	gameMatch.RuleManager, err = NewRuleManager(ruleManagerOption)
	if err != nil {
		util.MyPrint(err)
		return gameMatch, err
	}
	//gameMatch.PlayerManager = NewPlayerManager(gameMatch)
	//gameMatch.PlayerManager.TestRedisKey()

	gameMatch.RuleManager.StartupAll()

	go gameMatch.ListeningBridgeMsg()
	return gameMatch, nil
}

//退出
func (gameMatch *GameMatch) Quit(source int) {
	gameMatch.RuleManager.Quit()
}

//获取语言包
func (gameMatch *GameMatch) GetLang() map[int]util.ErrInfo {
	return gameMatch.Err.Pool
}

//获取配置信息
func (gameMatch *GameMatch) GetOption() GameMatchOption {
	return gameMatch.Option
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
