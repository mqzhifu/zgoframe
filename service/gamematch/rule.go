package gamematch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"zgoframe/util"
)

//匹配规则 - 配置 ,像： 附加属性
type Rule struct {
	Id 				int
	AppId			int
	CategoryKey 	string	//分类名，严格来说是：队列的KEY的值，也就是说，该字符串，可以有多级，如：c1_c2_c3,gameId_gameType
	MatchTimeout 	int		//匹配 - 超时时间
	SuccessTimeout	int		//匹配成功后，一直没有人来取，超时时间
	Flag			int 	//匹配机制-类型，目前两种：1. N（TEAM） VS N（TEAM）  ，  2. N人够了就行(吃鸡模式)
	PersonCondition	int		//触发 满足人数的条件,即匹配成功, 字段 flag = 2 ，该变量才有用
	TeamVSPerson	int		//如果是N VS N 的类型，得确定 每个队伍几个人,必须上面的flag = 1 ，该变量才有用,目前最大是5
	TeamVSNumber	int		//保留字段，还没想好，初衷：正常就是一个队伍跟一个队伍PK，但是可能会有多队伍互相PK，不一定是N V N
	PlayerWeight	PlayerWeight	//权重，目前是以最小单位：玩家属性，如果是小组/团队，是计算平均值
	GroupPersonMax	int		//玩家以组为单位，一个小组最大报名人数,暂定最大为：5，注：flag = 1
	//IsSupportGroup 	int		//保留字段，是否支持，以组为单位报名，疑问：所有的类型按说都应该支持这以组为单位的报名
}

type PlayerWeight struct {
	ScoreMin 	int		//权重值范围：最小值，范围： 1-100
	ScoreMax	int		//权重值范围：最大值，范围： 1-100
	AutoAssign	bool	//当权重值范围内，没有任何玩家，是否接收，自动调度分配，这样能提高匹配成功率
	Formula		string	//属性计算公式，由玩家的N个属性，占比，最终计算出权重值
	Aggregation string  //sum average min max 默认为：average
}
//Flag		int		//1、计算权重平均的区间的玩家，2、权重越高的匹配越快

//这个是后台的录入rule结构体
type GamesMatchConfig struct {
	Id 		int		`json:"id"`
	GamesId    int    `json:"games_id"`    //游戏ID
	Name       string `json:"name"`        //规则名称
	Status     int    `json:"status"`      //规则状态：1上线 2下线 3删除
	MatchCode  string `json:"match_code"`  //匹配代码
	TeamType   int    `json:"team_type"`   //1. N（TEAM） VS N（TEAM）  ，  2. N人够了就行(吃鸡模式) (后面这个注释有问题但又不敢删：团队类型 1各自为战 2对称团队战
	MaxPlayers int    `json:"max_players"` //匹配最大人数 如团队战代表每个队伍人数
	Rule	   string `json:"rule"`
	//RuleStruct PlayerWeight `json:"rule_struct"`        //表达式匹配规则
	Timeout    int    `json:"timeout"`     //匹配超时时间
	Fps        int    `json:"fps"`         //帧率
	SuccessTimeout int `json:"success_timeout"`
}

type RuleConfig struct {
	Data map[int]Rule
	gamematch *Gamematch	//父类/主类
	WatcherCancelFunc context.CancelFunc
}
//rule redis key 前缀
func (ruleConfig *RuleConfig) getRedisKey()string{
	return redisPrefix + redisSeparation + "rule"
}

func (ruleConfig *RuleConfig) getRedisIncKey()string{
	return redisPrefix + redisSeparation + "rule" + redisSeparation + "inc"
}
//初始化
//monitorRuleIds:可以指定监控哪些RULE,主要可以 负载均衡
func NewRuleConfig (gamematch *Gamematch,monitorRuleIds []int)(*RuleConfig,error){
	mylog.Info("NewRuleConfig , monitorRuleIds: ")
	rule := new (RuleConfig)
	rule.Data = make(map[int]Rule)
	//获取基配置数据，从etcd
	ruleList ,err:= rule.GetDataByEtcd( )
	if err != nil{
		return nil,err
	}
	if len(ruleList) <= 0 {
		return rule,myerr.New(601)
	}
	//过滤，只监听指定的rule
	if(len(monitorRuleIds) > 0 ){
		monitorRule  :=  make( map[int]Rule)
		for _,ruleOne := range ruleList{
			for _,monitorRuleId := range monitorRuleIds{
				if ruleOne.Id == monitorRuleId{
					monitorRule[ruleOne.Id] = ruleOne
					break
				}
			}
		}
		ruleList = monitorRule
	}
	if len(ruleList) <= 0 {
		return rule,myerr.New(625)
	}

	mylog.Info("rule final cnt : " + strconv.Itoa(len(ruleList)))
	for _,v := range ruleList{
		mylog.Info("match code : " + v.CategoryKey + " , id: " + strconv.Itoa(v.Id))
	}
	rule.gamematch = gamematch
	rule.Data = ruleList
	go rule.WatchEtcdChange()
	return rule,nil
}
//监听rule 基数据(配置)变更
func (ruleConfig *RuleConfig)WatchEtcdChange( ){
	prefix := "rule etcd watching"
	ctx , cancelFunc := context.WithCancel(context.Background())
	watchChann := myetcd.Watch(ctx,RuleEtcdConfigPrefix)
	ruleConfig.WatcherCancelFunc = cancelFunc
	//mylog.Notice(prefix , " , new key : ",RuleEtcdConfigPrefix)
	//watchChann := myetcd.Watch("/testmatch")
	for wresp := range watchChann{
		for _, ev := range wresp.Events{
			action := ev.Type.String()
			key := string(ev.Kv.Key)
			val := string(ev.Kv.Value)
			mylog.Warn(prefix  +  " chan has event , action: " + action  + " key:"+key  +  " val: "+val)
			//mylog.Warn(prefix  +  " key : " + key)
			//mylog.Warn(prefix  +  " val : " + val)
			//zlib.MyPrint(ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
			//fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			//etcd的存储，会带有前缀字符串，先过滤掉
			matchCode := strings.Replace(key,RuleEtcdConfigPrefix,"",-1)
			//过滤掉空格
			matchCode = strings.Trim(matchCode," ")
			matchCode = strings.Trim(matchCode,"/")
			mylog.Warn(prefix  + " matchCode : " + matchCode)

			rule,ruleErr := ruleConfig.getByCategory(matchCode)
			if ev.Type.String() == "PUT"{//添加/编辑
				mylog.Warn(prefix  + " action PUT : ruleErr ," + ruleErr.Error())
				if ruleErr != nil{
					mylog.Error("etcd event DELETE rule ,but Value no match rule , maybe add new rule~")
				}else{
					//关闭当前已启动rule的所有操作协程，重新加载变更后的rule基数据
					ruleConfig.gamematch.closeOneRuleDemonRoutine(rule.Id)
					//清空该rule的所有redis数据
					ruleConfig.delOne(rule.Id)
				}
				newRule,err := ruleConfig.parseOneConfigByEtcd(key,val)
				if err != nil{
					mylog.Error("etcd monitor:" + err.Error())
				}else{
					//新添加一条rule
					ruleConfig.Data[newRule.Id] = newRule
					//重新启动所有守护协程
					ruleConfig.gamematch.startOneRuleDomon(newRule)
					//ruleConfig.gamematch.HttpdRuleState[newRule.Id] = HTTPD_RULE_STATE_OK
					//mySleepSecond(3,"testtest")
					//zlib.ExitPrint(111111)
				}
			}else if ev.Type.String() == "DELETE"{
				//mylog.Warn(prefix ," dvent = DELETE : ruleErr ,",ruleErr)
				if ruleErr != nil{
					mylog.Error("etcd event DELETE rule ,but Value no match rule!!!")
				}else{
					ruleConfig.gamematch.closeOneRuleDemonRoutine(rule.Id)
					ruleConfig.delOne(rule.Id)
					//mymetrics.FastLog("Rule",zlib.METRICS_OPT_DIM,0)
				}
			}
		}
	}
	//zlib.ExitPrint(watchChann)
}
func (ruleConfig *RuleConfig)Shutdown(){
	ruleConfig.WatcherCancelFunc()
}
//删除一条rule,执行此方法前，要先停掉当前rule的各种操作守护协程
func (ruleConfig *RuleConfig)delOne(ruleId int){
	_, ok  := ruleConfig.GetById(ruleId)
	if !ok {
		mylog.Error("ruleConfig.GetByI is empty~")
		return
	}
	//从内存池中删除该rule info
	delete(ruleConfig.Data,ruleId)
	//删除报名池信息
	queueSign := ruleConfig.gamematch.GetContainerSignByRuleId(ruleId)
	queueSign.delOneRule()
	//删除推送池信息
	push := ruleConfig.gamematch.getContainerPushByRuleId(ruleId)
	push.delOneRule()
	//删除匹配成功池信息
	queueSuccess := ruleConfig.gamematch.getContainerSuccessByRuleId(ruleId)
	queueSuccess.delOneRule()
	//清空该池的玩家信息
	playerIds := playerStatus.getOneRuleAllPlayer(ruleId)

	redisConnFD := myredis.GetNewConnFromPool()
	defer redisConnFD.Close()
	myredis.Send(redisConnFD,"Multi")
	for _,playerId := range playerIds{
		playerStatus.delOneById(redisConnFD,util.Atoi(playerId))
	}
	myredis.ConnDo(redisConnFD,"exec")
}
//解析etcd 里的value 字符串 => json => struct
func (ruleConfig *RuleConfig)parseOneConfigByEtcd(k string ,v string)(rule Rule,err error){
	if k == ""{
		return rule,myerr.New(620)
	}
	if v == ""{
		return rule,myerr.New(621)
	}
	//zlib.MyPrint(v)
	gamesMatchConfig := GamesMatchConfig{}
	err = json.Unmarshal( []byte(v), & gamesMatchConfig)
	if err != nil{
		//zlib.MyPrint("parseOneConfigByEtcd",k,v)
		msg := myerr.MakeOneStringReplace("ruleCategory : " +rule.CategoryKey  + err.Error())
		myerr.NewReplace(622,msg)
		return rule,err
	}

	if gamesMatchConfig.Status != RuleStatusOnline{
		return rule,myerr.New(624)
	}
	//匹配规则 - 用户权重
	gamesMatchConfigRuleStruct := PlayerWeight{}
	if gamesMatchConfig.Rule != "" {
		err := json.Unmarshal([]byte(gamesMatchConfig.Rule),&gamesMatchConfigRuleStruct)
		if err != nil{
			//zlib.MyPrint(err)
			//mylog.Error("parseOneConfigByEtcd: PlayerWeight" + k)
			mylog.Error("parseOneConfigByEtcd PlayerWeight gamesMatchConfig.rule failed...." + gamesMatchConfig.Rule + err.Error())
		}
	}
	//playerWeightRow := PlayerWeight{}
	rule  = Rule{
		Id: gamesMatchConfig.Id,
		AppId: gamesMatchConfig.Id,
		CategoryKey : gamesMatchConfig.MatchCode,
		MatchTimeout: gamesMatchConfig.Timeout,
		SuccessTimeout: gamesMatchConfig.SuccessTimeout,
		Flag:gamesMatchConfig.TeamType,
		TeamVSPerson:gamesMatchConfig.MaxPlayers / 2,
		PersonCondition: gamesMatchConfig.MaxPlayers,
		GroupPersonMax : gamesMatchConfig.MaxPlayers / 2,
		PlayerWeight: gamesMatchConfigRuleStruct,
		//IsSupportGroup: 1,
	}
	//zlib.MyPrint("parseOneConfigByEtcd:",gamesMatchConfig)
	return rule,err
}

func (ruleConfig *RuleConfig)GetDataByEtcd()  (map[int]Rule,error){
	etcdRuleList,err := myetcd.GetListByPrefix(RuleEtcdConfigPrefix)
	if err !=nil{
		return nil,errors.New("ruleConfig getByEtcd err :" + err.Error())
	}
	ruleList := make(map[int]Rule)
	if len(etcdRuleList) == 0{
		return ruleList,nil
	}
	//i := 1
	for k,v := range etcdRuleList{
		ruleRow,err := ruleConfig.parseOneConfigByEtcd(k,v)
		if err != nil{
			//msg := myerr.MakeOneStringReplace(err.Error())
			//myerr.NewErrorCodeReplace(622,msg)
			continue
		}
		_,err =  ruleConfig.CheckRuleByElement(ruleRow)
		if err != nil{
			mylog.Warn("CheckRuleByElement err :" + err.Error())
			//myerr.NewErrorCode(621)
			continue
		}
		//zlib.MyPrint("ruleRow",ruleRow)
		ruleList[ruleRow.Id] = ruleRow
		//i++
	}
	return ruleList,nil

}

func (ruleConfig *RuleConfig)strToStruct(redisStr string)Rule{

	strArr := strings.Split(redisStr,separation)
	element := Rule{
		Id 				:	util.Atoi(strArr[0]),
		AppId 			:	util.Atoi(strArr[1]),
		CategoryKey 	:	strArr[2],
		MatchTimeout 	:	util.Atoi(strArr[3]),
		SuccessTimeout 	:	util.Atoi(strArr[4]),
		PersonCondition :	util.Atoi(strArr[5]),
		//IsSupportGroup 	:	zlib.Atoi(strArr[6]),
		Flag 			:	util.Atoi(strArr[6]),
		TeamVSPerson 		:	util.Atoi(strArr[7]),
		GroupPersonMax : util.Atoi(strArr[8]),
		//WeightRule 		:	strArr[0],
	}
	return element
}

func (ruleConfig *RuleConfig)structToStr(rule Rule)string{
	//groupPersonMax	int		//玩家以组为单位，一个小组最大报名人数,暂定最大为：5

	str := strconv.Itoa(rule.Id) + separation +
		strconv.Itoa(rule.AppId) + separation +
		rule.CategoryKey + separation +
		strconv.Itoa(rule.MatchTimeout) + separation +
		strconv.Itoa(rule.SuccessTimeout) + separation +
		strconv.Itoa(rule.PersonCondition) + separation +
		//strconv.Itoa(rule.IsSupportGroup) + separation +
		strconv.Itoa(rule.Flag) + separation +
		strconv.Itoa(rule.TeamVSPerson) + separation +
		strconv.Itoa(rule.GroupPersonMax)

		//rule.WeightRule + separation
	return str
}


func (ruleConfig *RuleConfig) GetById(id int ) (Rule,bool){
	if id == 0{
		return Rule{},false
	}
	rule,ok := ruleConfig.Data[id]
	return rule,ok
}

func (ruleConfig *RuleConfig) getAll()map[int]Rule{
	return ruleConfig.Data
}

func  (ruleConfig *RuleConfig) getByCategory(category string) (rule Rule ,err error){
	if category == ""{
		return rule,myerr.New(450)
	}

	for _,rule := range ruleConfig.getAll(){
		if  rule.CategoryKey == category{
			return rule,nil
		}
	}
	return rule,myerr.New(451)
}

func (ruleConfig *RuleConfig) getIncId( ) (int){
	key := ruleConfig.getRedisIncKey()
	res,_ := redis.Int(myredis.RedisDo("INCR",key))
	return res
}

func (ruleConfig *RuleConfig) AddOne(rule Rule)(bool,error){
	checkRs,errs := ruleConfig.CheckRuleByElement(rule)
	if !checkRs{
		return false,errs
	}
	key := ruleConfig.getRedisKey()
	//id := ruleConfig.getIncId()
	ruleStr := ruleConfig.structToStr(rule)
	util.MyPrint("ruleStr : ",ruleStr, " rule struct : ",rule)
	_ ,errs = redis.Int( myredis.RedisDo("hset",redis.Args{}.Add(key).Add(rule.Id).Add(ruleStr)...))
	if errs != nil{
		return false,myerr.New(603)
	}
	return true,nil
}

func (ruleConfig *RuleConfig) CheckRuleByElement(rule Rule)(bool,error){
	if rule.Id <= 0{
		return false,myerr.New(604)
	}
	if rule.AppId <= 0{
		return false,myerr.New(605)
	}
	if rule.CategoryKey == ""{
		return false,myerr.New(616)
	}
	if rule.Flag <= 0{
		return false,myerr.New(606)
	}
	if rule.Flag == RuleFlagTeamVS{
		if rule.TeamVSPerson <= 0{
			return false,myerr.New(608)
		}

		if rule.TeamVSPerson > RuleTeamVSPersonMax{
			return false,myerr.NewReplace(609,myerr.MakeOneStringReplace(strconv.Itoa(RuleTeamVSPersonMax)))
		}
		//TeamVSPerson	int		//如果是N VS N 的类型，得确定 每个队伍几个人,必须上面的flag = 1 ，该变量才有用,目前最大是5
	}else if rule.Flag == RuleFlagCollectPerson{
		if rule.PersonCondition <= 0{
			return false,myerr.New(610)
		}

		if rule.PersonCondition > RulePersonConditionMax{
			return false,myerr.NewReplace(611,myerr.MakeOneStringReplace(strconv.Itoa(RuleTeamVSPersonMax)))
		}
	}else{
		return false,myerr.New(607)
	}
	if rule.MatchTimeout < RuleMatchTimeoutMin || rule.MatchTimeout > RuleMatchTimeoutMax{
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(RuleMatchTimeoutMin)
		msg[1] = strconv.Itoa(RuleMatchTimeoutMax)
		return false,myerr.NewReplace(612,msg)
	}

	if rule.SuccessTimeout < RuleSuccessTimeoutMin || rule.SuccessTimeout > RuleSuccessTimeoutMax{
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(RuleSuccessTimeoutMin)
		msg[1] = strconv.Itoa(RuleSuccessTimeoutMax)
		return false,myerr.NewReplace(613,msg)
	}

	if rule.GroupPersonMax <= 0{
		util.MyPrint(rule.GroupPersonMax )
		util.ExitPrint(rule)
		return false,myerr.New(614)
	}

	if rule.GroupPersonMax > RuleGroupPersonMax{
		return false,myerr.NewReplace(615,myerr.MakeOneStringReplace(strconv.Itoa(RuleGroupPersonMax)))
	}

	if rule.PlayerWeight.Formula != ""{
		if rule.PlayerWeight.ScoreMin > rule.PlayerWeight.ScoreMax{
			return false,myerr.New(617)
		}
	}

	//PlayerWeight	PlayerWeight	//权重，目前是以最小单位：玩家属性，如果是小组/团队，是计算平均值
	return true,nil
}

func (ruleConfig *RuleConfig) CheckRuleById(ruleId int)(bool,error){
	rule , ok := ruleConfig.GetById(ruleId)
	if !ok {
		msg := make(map[int]string)
		msg[0] = strconv.Itoa(ruleId)
		return false,myerr.NewReplace(602,msg)
	}

	return ruleConfig.CheckRuleByElement(rule)


}