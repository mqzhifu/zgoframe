package gamematch

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"zgoframe/http/request"
	"zgoframe/service"
	"zgoframe/util"
)

/*
	报名 - 加入匹配队列
	所有玩家是以组为单位的，校验是要对：单个玩家+组ID
	大概步骤：
	1. 校验基础参数
	2. 校验每个玩家的状态
	3. 处理权重
	4. 增加到匹配池
*/
func (gameMatch *GameMatch) PlayerJoin(form request.HttpReqGameMatchPlayerSign) (group Group, err error) {
	formBytes, err := json.Marshal(&form)
	gameMatch.Option.Log.Debug("PlayerJoin " + string(formBytes))
	ruleId := form.RuleId
	outGroupId := form.GroupId
	//这里只做最基础的验证，前置条件是已经在HTTP层作了验证
	rule, err := gameMatch.RuleManager.GetById(ruleId)
	if err != nil {
		return group, gameMatch.Err.New(400)
	}

	if rule.Status != service.GAME_MATCH_RULE_STATUS_EXEC {
		return group, gameMatch.Err.New(400)
	}

	lenPlayers := len(form.PlayerList)
	if lenPlayers <= 0 {
		return group, gameMatch.Err.New(401)
	}
	if form.GroupId <= 0 {
		return group, gameMatch.Err.New(452)
	}
	//groupsTotal := queueSign.getAllGroupsWeightCnt()	//报名 小组总数
	//playersTotal := queueSign.getAllPlayersCnt()	//报名 玩家总数
	//mylog.Info(" action :  Sign , players : " + strconv.Itoa(lenPlayers) +" ,queue cnt : groupsTotal",groupsTotal ," , playersTotal",playersTotal)
	//queueSign := gamematch.GetContainerSignByRuleId(ruleId)
	//mylog.Info("new sign :[ruleId : " + strconv.Itoa(ruleId) + "(" + rule.CategoryKey + ") , outGroupId : " + strconv.Itoa(outGroupId) + " , playersCount : " + strconv.Itoa(lenPlayers) + "] ")
	//mylog.Info("new sign :[ruleId : ",ruleId,"(",rule.CategoryKey,") , outGroupId : ",outGroupId," , playersCount : ",lenPlayers,"] ")
	//queueSign.Log.Info("new sign :[ruleId : " ,  ruleId   ,"(",rule.CategoryKey,") , outGroupId : ",outGroupId," , playersCount : ",lenPlayers,"] ")
	processStartTime := util.GetNowTimeSecondToInt()
	//util.PrintStruct(queueSign.Rule, ":")
	if lenPlayers > rule.TeamMaxPeople {
		errMsgMap := gameMatch.Err.MakeOneStringReplace(" rule.TeamMaxPeople:" + strconv.Itoa(rule.TeamMaxPeople))
		return group, gameMatch.Err.NewReplace(408, errMsgMap)
	}
	if lenPlayers > rule.ConditionPeople {
		errMsgMap := gameMatch.Err.MakeOneStringReplace(strconv.Itoa(rule.ConditionPeople) + " > " + strconv.Itoa(lenPlayers))
		return group, gameMatch.Err.NewReplace(410, errMsgMap)
	}
	//这里有个小问题，目前仅支持：每个组最多5人，回头我再优化
	if lenPlayers > gameMatch.RuleTeamMaxPeople {
		return group, gameMatch.Err.New(411)
	}

	//验证3方传过来的groupId 是否重复
	allGroupIds := rule.QueueSign.GetGroupSignTimeoutAll()
	for _, hasGroupId := range allGroupIds {
		if outGroupId == hasGroupId {
			//util.MyPrint(allGroupIds, outGroupId, hasGroupId, httpReqBusiness)
			errMsg := gameMatch.Err.MakeOneStringReplace(strconv.Itoa(outGroupId))
			return group, gameMatch.Err.NewReplace(409, errMsg)
		}
	}

	gameMatch.Option.Log.Info("check base info finish , start check player status :")
	//检查，所有玩家的状态
	var playerList []Player
	for _, httpPlayer := range form.PlayerList {
		if httpPlayer.Uid <= 0 {
			return group, gameMatch.Err.New(412)
		}
		//player := Player{Id: httpPlayer.Uid, MatchAttr: httpPlayer.MatchAttr}
		player, isEmpty := rule.PlayerManager.GetById(httpPlayer.Uid)
		//queueSign.Log.Info("player(" + strconv.Itoa(player.Id) + ") GetById :  status = " + strconv.Itoa(playerStatusElement.Status) + " isEmpty:" + strconv.Itoa(isEmpty))
		if isEmpty == 1 {
			//这是正常
		} else if player.Status == service.GAME_MATCH_PLAYER_STATUS_SUCCESS { //玩家已经匹配成功，并等待开始游戏
			//queueSign.Log.Error(" player status = PlayerStatusSuccess ,demon not clean.")
			errMsg := gameMatch.Err.MakeOneStringReplace("strconv.Itoa(player.Id)")
			return group, gameMatch.Err.NewReplace(403, errMsg)
		} else if player.Status == service.GAME_MATCH_PLAYER_STATUS_SIGN { //报名成功，等待匹配
			isTimeout := rule.PlayerManager.checkSignTimeout(player)
			if !isTimeout { //未超时
				//queueSign.Log.Error(" player status = matching...  not timeout")
				errMsg := gameMatch.Err.MakeOneStringReplace(strconv.Itoa(player.Id))
				return group, gameMatch.Err.NewReplace(402, errMsg)
			} else { //报名已超时，等待后台守护协程处理
				//这里其实也可以先一步处理，但是怕与后台协程冲突
				//queueSign.Log.Error(" player status is timeout ,but not clear , wait a moment!!!")
				errMsg := gameMatch.Err.MakeOneStringReplace(strconv.Itoa(player.Id))
				return group, gameMatch.Err.NewReplace(407, errMsg)
			}
		}
		//player.WeightAttrs = httpPlayer.WeightAttr
		playerList = append(playerList, player)
		//playerStatusElementMap[player.Id] = playerStatusElement
	}
	gameMatch.Option.Log.Info("finish check player status ，start calculate group/player weight:")

	//先计算一下权重平均值
	var groupWeightTotal float32
	groupWeightTotal = 0.00

	if rule.Formula != "" {
		//util.MyPrint(queueSign, "rule weight , Formula : ", rule.PlayerWeight.Formula)
		var weight float32
		weight = 0.00
		var playerWeightValue []float32
		for k, p := range playerList {
			onePlayerWeight := gameMatch.getPlayerWeightByFormula(rule.Formula, p.WeightAttrs)
			util.MyPrint("onePlayerWeight : ", onePlayerWeight)
			if onePlayerWeight > float32(gameMatch.WeightMaxValue) {
				onePlayerWeight = float32(gameMatch.WeightMaxValue)
			}
			weight += onePlayerWeight
			playerWeightValue = append(playerWeightValue, onePlayerWeight)
			playerList[k].Weight = onePlayerWeight
		}
		switch rule.WeightTeamAggregation {
		case "sum":
			groupWeightTotal = weight
		case "min":
			groupWeightTotal = util.FindMinNumInArrFloat32(playerWeightValue)
		case "max":
			groupWeightTotal = util.FindMaxNumInArrFloat32(playerWeightValue)
		case "average":
			groupWeightTotal = weight / float32(len(playerList))
		default:
			groupWeightTotal = weight / float32(len(playerList))
		}
		//保留2位小数
		tmp := util.FloatToString(groupWeightTotal, 2)
		groupWeightTotal = util.StringToFloat(tmp)
	} else {
		gameMatch.Option.Log.Info("rule.Formula empty , no need calculate group/player weight.")
	}

	//这里有个特殊的情况:报名(人数)即满足条件
	//按说也不应该进到我这儿，直接调用同步服务即可，可能为了统一走一个服务，且调用方不想麻烦再写代码了
	if rule.Type == service.RULE_TYPE_TEAM_EACH_OTHER {
		if lenPlayers == rule.ConditionPeople {
			go rule.Push.RuntimeSuccess(form, ruleId)
			return
		}
	}
	//这里再做一次检查，防止，此时某个 rule 关闭了
	if rule.Status != service.GAME_MATCH_RULE_STATUS_EXEC {
		return group, gameMatch.Err.New(400)
	}

	//下面两行必须是原子操作，如果pushOne执行成功，但是 upInfo 没成功会导致报名队列里，同一个用户能再报名一次
	redisConnFD := gameMatch.Option.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()
	gameMatch.Option.Redis.Multi(redisConnFD) //开始多指令缓存模式
	//匹配超时时间
	expire := processStartTime + rule.MatchTimeout
	//创建一个新的小组
	group = gameMatch.NewGroupStruct(rule)
	//这里有偷个懒，还是用外部的groupId , 不想再给redis加 groupId映射outGroupId了
	gameMatch.Option.Log.Warn(" outGroupId replace groupId :" + strconv.Itoa(outGroupId) + " " + strconv.Itoa(group.Id))
	group.Id = outGroupId
	group.Players = playerList
	group.Type = service.GAME_MATCH_GROUP_TYPE_SIGN
	group.SignTimeout = expire
	group.Person = len(playerList)
	group.Weight = groupWeightTotal
	group.OutGroupId = outGroupId
	group.Addition = form.Addition
	//group.CustomProp = httpReqBusiness.CustomProp
	//group.MatchCode = rule.CategoryKey
	//util.MyPrint(queueSign, "newGroupId : ", group.Id, "player/group weight : ", groupWeightTotal, " now : ", now, " expire : ", expire)
	//mylog.Info("newGroupId : ",group.Id , "player/group weight : " ,groupWeightTotal ," now : ",now ," expire : ",expire )
	//queueSign.Log.Info("newGroupId : ",group.Id , "player/group weight : " ,groupWeightTotal ," now : ",now ," expire : ",expire)
	rule.QueueSign.AddOne(group, redisConnFD)
	playerIds := ""
	for _, player := range playerList {

		newPlayerStatusElement := rule.PlayerManager.createEmptyPlayer()
		newPlayerStatusElement.Id = player.Id
		newPlayerStatusElement.Status = service.GAME_MATCH_PLAYER_STATUS_SIGN
		newPlayerStatusElement.Weight = player.Weight
		newPlayerStatusElement.SignTimeout = expire
		newPlayerStatusElement.GroupId = group.Id
		//newPlayerStatusElement.RuleId = ruleId
		//queueSign.Log.Info("playerStatus.upInfo:" + strconv.Itoa(service.PlayerStatusSign))
		rule.PlayerManager.Create(newPlayerStatusElement, redisConnFD)

		playerIds += strconv.Itoa(player.Id) + ","
	}
	//提交缓存中的指令
	_, err = gameMatch.Option.Redis.Exec(redisConnFD)
	//if err != nil {
	//	queueSign.Log.Error("transaction failed : " + err.Error())
	//}
	//queueSign.Log.Info(" sign finish ,total : newGroupId " + strconv.Itoa(group.Id) + " success players : " + strconv.Itoa(len(players)))
	gameMatch.Option.Log.Info(" sign finish ,total : newGroupId " + strconv.Itoa(group.Id) + " success players : " + strconv.Itoa(len(playerList)))

	//signSuccessReturnData = SignSuccessReturnData{
	//	RuleId: ruleId,
	//	GroupId: outGroupId,
	//	PlayerIds: playerIds,
	//
	//}

	return group, nil
}

func (gameMatch *GameMatch) getPlayerWeightByFormula(formula string, MatchAttr map[string]int) float32 {
	//mylog.Debug("getPlayerWeightByFormula , formula:",formula)
	grep := gameMatch.FormulaFirst + "([\\s\\S]*?)" + gameMatch.FormulaEnd
	var imgRE = regexp.MustCompile(grep)
	findRs := imgRE.FindAllStringSubmatch(formula, -1)
	//util.MyPrint(sign, "parse PlayerWeightByFormula : ", findRs)
	if len(findRs) == 0 {
		return 0
	}
	for _, v := range findRs {
		val, ok := MatchAttr[v[1]]
		if !ok {
			val = 0
		}
		formula = strings.Replace(formula, v[0], strconv.Itoa(val), -1)

	}
	//util.MyPrint(sign, "final formula replaced str :", formula)
	rs, err := util.Eval(formula)
	if err != nil {
		return 0
	}
	f, _ := rs.Float32()
	return f
}
