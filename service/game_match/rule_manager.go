package gamematch

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"zgoframe/service/bridge"
	"zgoframe/util"
)

type Rule struct {
	model.GameMatchRule
	Prefix             string         `json:"prefix"`
	Status             int            `json:"status"`
	DemonDebugTime     int            `json:"demon_debug_time"`      //守护协程，在没有处理数据时，需要输出日志，太多，每X秒输出一次
	DemonDebugShowTime int            `json:"demon_debug_show_time"` //配合上面一起使用，用于记录一次的输出时间
	QueueSign          *QueueSign     `json:"-"`
	QueueSuccess       *QueueSuccess  `json:"-"`
	Push               *Push          `json:"-"`
	Match              *Match         `json:"-"`
	PlayerManager      *PlayerManager `json:"-"`
	RuleManager        *RuleManager   `json:"-"`
}

type RuleManagerOption struct {
	Gorm           *gorm.DB
	GameMatch      *GameMatch //父类
	MonitorRuleIds []int      //负载时使用，当某个 rule 负载较大时，拆分到其它机器上，其它机器上启动的进程仅监听此 rule 即可
	//RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
	ServiceBridge *bridge.Bridge
}
type RuleManager struct {
	Option RuleManagerOption
	pool   []*Rule
	Log    *zap.Logger
	prefix string
	Err    *util.ErrMsg
	//WatcherCancelFunc context.CancelFunc
}

func NewRuleManager(option RuleManagerOption) (*RuleManager, error) {
	option.GameMatch.Option.Log.Info("NewRuleManager:")

	ruleManager := new(RuleManager)
	ruleManager.prefix = "RuleManager"
	ruleManager.Option = option
	ruleManager.Log = ruleManager.Option.GameMatch.Option.Log
	ruleManager.Err = ruleManager.Option.GameMatch.Err

	err := ruleManager.InitData()
	return ruleManager, err
}

// 从 3方容器中读取出 所有 rule 的配置信息
func (ruleManager *RuleManager) InitData() (err error) {
	ruleManager.Log.Info(ruleManager.prefix + " init data RuleDataSourceType:" + strconv.Itoa(ruleManager.Option.GameMatch.Option.RuleDataSourceType))
	var list []model.GameMatchRule
	switch ruleManager.Option.GameMatch.Option.RuleDataSourceType {
	case GAME_MATCH_DATA_SOURCE_TYPE_ETCD:
		return errors.New("not support etcd.")
	case GAME_MATCH_DATA_SOURCE_TYPE_DB:
		list, err = ruleManager.GetDataByDb()
	case GAME_MATCH_DATA_SOURCE_TYPE_SERVICE:
		return errors.New("not support GAME_MATCH_DATA_SOURCE_TYPE_SERVICE.")
	default:
		return errors.New("dataSourceType err")
	}
	if err != nil {
		return err
	}
	//上面读取的是基础配置信息的数据，现在要给该条 rule 挂载 具体的实现类
	for _, rule := range list {
		if len(ruleManager.Option.MonitorRuleIds) > 0 {
			hasSearchRuleId := 0
			for _, MonitorRuleId := range ruleManager.Option.MonitorRuleIds {
				if rule.Id == MonitorRuleId {
					hasSearchRuleId = 1
					break
				}

			}
			if hasSearchRuleId == 0 {
				continue
			}
		}

		oneRule := Rule{}
		oneRule.DemonDebugTime = ruleManager.Option.GameMatch.Option.RuleDebugShow
		oneRule.Prefix = ruleManager.Option.GameMatch.prefix + "_rule_" + strconv.Itoa(rule.Id)
		oneRule.Status = GAME_MATCH_RULE_STATUS_INIT
		oneRule.RuleManager = ruleManager
		oneRule.GameMatchRule = rule
		oneRule.DemonDebugTime = 5
		err = ruleManager.CheckRule(oneRule)
		if err != nil {
			return err
		}
		oneRule.QueueSign = NewQueueSign(&oneRule)
		oneRule.QueueSign.TestRedisKey()
		oneRule.QueueSuccess = NewQueueSuccess(&oneRule)
		oneRule.QueueSuccess.TestRedisKey()
		oneRule.PlayerManager = NewPlayerManager(&oneRule)
		oneRule.PlayerManager.TestRedisKey()
		oneRule.Push = NewPush(&oneRule)
		oneRule.Push.TestRedisKey()
		oneRule.Match = NewMatch(&oneRule)

		ruleManager.pool = append(ruleManager.pool, &oneRule)
	}

	return nil
}

// 每10秒 输出一次，避免日志过多
func (rule *Rule) NothingToDoLog(msg string) {
	now := util.GetNowTimeSecondToInt()
	if now%rule.DemonDebugTime == 0 && now != rule.DemonDebugShowTime {
		rule.DemonDebugShowTime = now
		rule.RuleManager.Log.Info(msg)
	}
}
func (ruleManager *RuleManager) GetDataByDb() (list []model.GameMatchRule, err error) {
	err = ruleManager.Option.Gorm.Where("status = 1").Find(&list).Error
	if err != nil {
		return list, err
	}

	if len(list) <= 0 {
		return list, errors.New("is empty")
	}

	return list, nil
}
func (ruleManager *RuleManager) GetById(id int) (rule *Rule, err error) {
	for _, v := range ruleManager.pool {
		if v.Id == id {
			return v, nil
		}
	}

	return rule, errors.New("is empty")
}

func (ruleManager *RuleManager) StartupAll() error {
	queueLen := len(ruleManager.pool)
	ruleManager.Log.Info(ruleManager.prefix + " , StartupAll  ,  rule total : " + strconv.Itoa(queueLen))
	if queueLen <= 0 {
		msg := "RuleConfig list is empty!!!"
		ruleManager.Log.Error(msg)
		return errors.New(msg)
	}
	//开始每个rule
	for _, rule := range ruleManager.pool {
		ruleManager.startOneRuleDemon(rule)
	}
	//后台守护协程均已开启完毕，可以开启前端HTTPD入口了
	//gamematch.StartHttpd(gamematch.Option.HttpdOption)
	return nil
}

// 开启一条rule的所有守护协程，
// 虽然有4个，但是只有match是最核心、最复杂的，另外3个算是辅助
func (ruleManager *RuleManager) startOneRuleDemon(rule *Rule) {
	ruleManager.Log.Info(ruleManager.prefix + "  startOneRuleDemon:" + strconv.Itoa(rule.Id))
	//启动守护协程后，要先执行一下，超时检测，都完成后，该 rule 的状态才是正常，才能接收新的报名，不然，redis 里数据会出现混乱
	//这里必须是同步，不能异步，且必须得是顺序执行
	rule.QueueSign.CheckTimeout()
	rule.QueueSuccess.CheckTimeout()
	rule.Push.checkStatus()

	go rule.QueueSign.Demon()
	//报名成功
	go rule.QueueSuccess.Demon()
	//推送
	go rule.Push.Demon()
	//匹配
	go rule.Match.Demon()

	rule.Status = GAME_MATCH_RULE_STATUS_EXEC
}

func (ruleManager *RuleManager) Quit() {
	for _, v := range ruleManager.pool {
		//从内存池中删除该rule info
		v.Close()
	}
}

// redis公共前缀+模块名
func (rule *Rule) GetCommRedisKeyByModule(module string) string {
	return rule.RuleManager.Option.GameMatch.Option.RedisPrefix + rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator + module + rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
}

// redis公共前缀+模块名+ruleId
func (rule *Rule) GetCommRedisKeyByModuleRuleId(module string, ruleId int) string {
	return rule.GetCommRedisKeyByModule(module) + strconv.Itoa(ruleId) + rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
}

func (rule *Rule) Close() {
	rule.Status = GAME_MATCH_RULE_STATUS_CLOSE
	rule.QueueSign.Close()
	rule.QueueSuccess.Close()
	rule.Push.Close()
	rule.Match.Close()
}

// 删除一条rule,执行此方法前，要先停掉当前rule的各种操作守护协程
func (ruleManager *RuleManager) delOneRuleRedisData(ruleId int) {
	rule, err := ruleManager.GetById(ruleId)
	if err != nil {
		ruleManager.Log.Error("ruleConfig.GetByI is empty~")
		return
	}
	//删除报名池信息
	rule.QueueSign.delOneRule()
	//删除推送池信息
	rule.Push.delOneRule()
	//删除匹配成功池信息
	rule.QueueSuccess.delOneRule()
	//清空该池的玩家信息
	playerIds, _ := rule.PlayerManager.getOneRuleAllPlayer()

	redisConnFD := ruleManager.Option.GameMatch.Option.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()
	ruleManager.Option.GameMatch.Option.Redis.Send(redisConnFD, "Multi")
	for _, playerId := range playerIds {
		rule.PlayerManager.delOneById(redisConnFD, util.Atoi(playerId))
	}
	ruleManager.Option.GameMatch.Option.Redis.ConnDo(redisConnFD, "exec")
}

// 一条 rule 是后台录入的，属于其它的项目，最终会持久化到某个地方
// 而匹配服务读取进来，并不确定该条记录是正确的
// 所以需要，验证正确，才能使用
func (ruleManager *RuleManager) CheckRule(rule Rule) error {
	if rule.Id <= 0 {
		return ruleManager.Err.New(604)
	}
	//if rule.AppId <= 0 {
	//	return false, myerr.New(605)
	//}
	//if rule.CategoryKey == "" {
	//	return false, myerr.New(616)
	//}
	if rule.Type <= 0 {
		return ruleManager.Err.New(606)
	}

	if rule.TeamMaxPeople <= 0 {
		return ruleManager.Err.New(608) //614
	}

	if rule.TeamMaxPeople > ruleManager.Option.GameMatch.Option.RuleTeamMaxPeople {
		return ruleManager.Err.New(615)
	}

	if rule.ConditionPeople <= 0 {
		return ruleManager.Err.New(610)
	}

	if rule.Type == RULE_TYPE_TEAM_VS {
		if rule.ConditionPeople%2 != 0 {
			return errors.New("组队互相撕杀，仅支持两个队伍，那么满足条件总人数肯定是偶数")
		}

		if rule.TeamMaxPeople > rule.ConditionPeople {
			return errors.New("组队互相撕杀，仅支持两个队伍，每个队伍最大支持5人， 剩2 肯定是 < 10人的")
		}
	} else if rule.Type == RULE_TYPE_TEAM_EACH_OTHER {
		if rule.ConditionPeople > ruleManager.Option.GameMatch.Option.RulePersonConditionMax {
			return ruleManager.Err.NewReplace(611, ruleManager.Err.MakeOneStringReplace(strconv.Itoa(ruleManager.Option.GameMatch.Option.RulePersonConditionMax)))
		}
	} else {
		return ruleManager.Err.New(607)
	}

	if rule.MatchTimeout < ruleManager.Option.GameMatch.Option.RuleMatchTimeoutMin || rule.MatchTimeout > ruleManager.Option.GameMatch.Option.RuleMatchTimeoutMax {
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(ruleManager.Option.GameMatch.Option.RuleMatchTimeoutMin)
		msg[1] = strconv.Itoa(ruleManager.Option.GameMatch.Option.RuleMatchTimeoutMax)
		return ruleManager.Err.NewReplace(612, msg)
	}

	if rule.SuccessTimeout < ruleManager.Option.GameMatch.Option.RuleSuccessTimeoutMin || rule.SuccessTimeout > ruleManager.Option.GameMatch.Option.RuleSuccessTimeoutMax {
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(ruleManager.Option.GameMatch.Option.RuleSuccessTimeoutMin)
		msg[1] = strconv.Itoa(ruleManager.Option.GameMatch.Option.RuleSuccessTimeoutMax)
		return ruleManager.Err.NewReplace(613, msg)
	}

	if rule.Formula != "" {
		if rule.WeightScoreMin > rule.WeightScoreMax {
			return ruleManager.Err.New(617)
		}

		if rule.WeightScoreMax > ruleManager.Option.GameMatch.Option.WeightMaxValue {
			return errors.New("rule > WeightMaxValue")
		}
	}

	return nil
}
