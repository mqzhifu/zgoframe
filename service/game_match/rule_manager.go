package gamematch

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"zgoframe/service"
	"zgoframe/util"
)

type Rule struct {
	model.GameMatchRule
	QueueSign     *QueueSign     `json:"-"`
	QueueSuccess  *QueueSuccess  `json:"-"`
	Push          *Push          `json:"-"`
	Match         *Match         `json:"Match"`
	PlayerManager *PlayerManager `json:"-"`
	RuleManager   *RuleManager   `json:"-"`
	Status        int            `json:"status"`
}

type RuleManagerOption struct {
	Gorm      *gorm.DB
	GameMatch *GameMatch //父类
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

//从 3方容器中读取出 所有 rule 的配置信息
func (ruleManager *RuleManager) InitData() (err error) {
	ruleManager.Log.Info(ruleManager.prefix + " init data RuleDataSourceType:" + strconv.Itoa(ruleManager.Option.GameMatch.Option.RuleDataSourceType))
	var list []model.GameMatchRule
	switch ruleManager.Option.GameMatch.Option.RuleDataSourceType {
	case service.GAME_MATCH_DATA_SOURCE_TYPE_ETCD:
		return errors.New("not support etcd.")
	case service.GAME_MATCH_DATA_SOURCE_TYPE_DB:
		list, err = ruleManager.GetDataByDb()
	case service.GAME_MATCH_DATA_SOURCE_TYPE_SERVICE:
		return errors.New("not support GAME_MATCH_DATA_SOURCE_TYPE_SERVICE.")
	default:
		return errors.New("dataSourceType err")
	}
	if err != nil {
		return err
	}
	//上面读取的是基础配置信息的数据，现在要给该条 rule 挂载 具体的实现类
	for _, v := range list {
		oneRule := Rule{}
		oneRule.Status = service.GAME_MATCH_RULE_STATUS_INIT
		oneRule.RuleManager = ruleManager
		oneRule.GameMatchRule = v

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

		//util.ExitPrint(33)
		ruleManager.pool = append(ruleManager.pool, &oneRule)
		oneRule.Status = service.GAME_MATCH_RULE_STATUS_EXEC
	}

	return nil
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

//开启一条rule的所有守护协程，
//虽然有4个，但是只有match是最核心、最复杂的，另外3个算是辅助
func (ruleManager *RuleManager) startOneRuleDemon(rule *Rule) {
	ruleManager.Log.Info(ruleManager.prefix + "  startOneRuleDemon:" + strconv.Itoa(rule.Id))
	go rule.QueueSign.Demon()
	//报名成功
	go rule.QueueSuccess.Demon()
	//推送
	go rule.Push.Demon()
	//匹配
	go rule.Match.Demon()
}

func (ruleManager *RuleManager) Quit() {
	for _, v := range ruleManager.pool {
		//从内存池中删除该rule info
		v.Close()
	}
}

//redis公共前缀+模块名
func (rule *Rule) GetCommRedisKeyByModule(module string) string {
	return rule.RuleManager.Option.GameMatch.Option.RedisPrefix + rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator + module + rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
}

//redis公共前缀+模块名+ruleId
func (rule *Rule) GetCommRedisKeyByModuleRuleId(module string, ruleId int) string {
	return rule.GetCommRedisKeyByModule(module) + strconv.Itoa(ruleId) + rule.RuleManager.Option.GameMatch.Option.RedisKeySeparator
}

func (rule *Rule) Close() {
	rule.Status = service.GAME_MATCH_RULE_STATUS_CLOSE
	rule.QueueSign.Close()
	rule.QueueSuccess.Close()
	rule.Push.Close()
	rule.Match.Close()
}

//删除一条rule,执行此方法前，要先停掉当前rule的各种操作守护协程
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

//一条 rule 是后台录入的，属于其它的项目，最终会持久化到某个地方
//而匹配服务读取进来，并不确定该条记录是正确的
//所以需要，验证正确，才能使用
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

	if rule.TeamMaxPeople > ruleManager.Option.GameMatch.RuleTeamMaxPeople {
		return ruleManager.Err.New(615)
	}

	if rule.ConditionPeople <= 0 {
		return ruleManager.Err.New(610)
	}

	if rule.Type == service.RULE_TYPE_TEAM_VS {
		if rule.ConditionPeople%2 != 0 {
			return errors.New("组队互相撕杀，仅支持两个队伍，那么满足条件总人数肯定是偶数")
		}

		if rule.TeamMaxPeople*2 > rule.ConditionPeople {
			return errors.New("组队互相撕杀，仅支持两个队伍，每个队伍最大支持5人， 剩2 肯定是 < 10人的")
		}
	} else if rule.Type == service.RULE_TYPE_TEAM_EACH_OTHER {
		if rule.ConditionPeople > ruleManager.Option.GameMatch.RulePersonConditionMax {
			return ruleManager.Err.NewReplace(611, ruleManager.Err.MakeOneStringReplace(strconv.Itoa(ruleManager.Option.GameMatch.RulePersonConditionMax)))
		}
	} else {
		return ruleManager.Err.New(607)
	}

	if rule.MatchTimeout < ruleManager.Option.GameMatch.RuleMatchTimeoutMin || rule.MatchTimeout > ruleManager.Option.GameMatch.RuleMatchTimeoutMax {
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(ruleManager.Option.GameMatch.RuleMatchTimeoutMin)
		msg[1] = strconv.Itoa(ruleManager.Option.GameMatch.RuleMatchTimeoutMax)
		return ruleManager.Err.NewReplace(612, msg)
	}

	if rule.SuccessTimeout < ruleManager.Option.GameMatch.RuleSuccessTimeoutMin || rule.SuccessTimeout > ruleManager.Option.GameMatch.RuleSuccessTimeoutMax {
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(ruleManager.Option.GameMatch.RuleSuccessTimeoutMin)
		msg[1] = strconv.Itoa(ruleManager.Option.GameMatch.RuleSuccessTimeoutMax)
		return ruleManager.Err.NewReplace(613, msg)
	}

	if rule.Formula != "" {
		if rule.WeightScoreMin > rule.WeightScoreMax {
			return ruleManager.Err.New(617)
		}

		if rule.WeightScoreMax > ruleManager.Option.GameMatch.WeightMaxValue {
			return errors.New("rule > WeightMaxValue")
		}
	}

	return nil
}
