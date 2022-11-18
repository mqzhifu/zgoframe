package gamematch

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"zgoframe/protobuf/pb"
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
func (gameMatch *GameMatch) PlayerJoin(form pb.GameMatchSign) (group Group, err error) {
	//func (gameMatch *GameMatch) PlayerJoin(form request.HttpReqGameMatchPlayerSign) (group Group, err error) {
	formBytes, err := json.Marshal(&form)
	gameMatch.Option.Log.Debug("PlayerJoin " + string(formBytes))
	ruleId := form.RuleId

	//这里只做最基础的验证，前置条件是已经在HTTP层作了验证
	rule, err := gameMatch.RuleManager.GetById(int(ruleId))
	if err != nil {
		return group, gameMatch.Err.New(400)
	}

	if rule.Status != service.GAME_MATCH_RULE_STATUS_EXEC {
		return group, gameMatch.Err.New(624)
	}

	lenPlayers := len(form.PlayerList)
	if lenPlayers <= 0 {
		return group, gameMatch.Err.New(401)
	}
	if form.GroupId <= 0 {
		//这里是做一下兼容，按说应该直接返回错误。但有时候就是单用户直接报名
		form.GroupId = int32(gameMatch.GetGroupIncId(int(form.RuleId)))
		//return group, gameMatch.Err.New(452)
	}
	outGroupId := form.GroupId

	gameMatch.DebugShowQueueInfo(rule, form)
	processStartTime := util.GetNowTimeSecondToInt()
	if lenPlayers > rule.TeamMaxPeople {
		errMsgMap := gameMatch.Err.MakeOneStringReplace(" rule.TeamMaxPeople:" + strconv.Itoa(rule.TeamMaxPeople))
		return group, gameMatch.Err.NewReplace(408, errMsgMap)
	}
	if lenPlayers > rule.ConditionPeople {
		errMsgMap := gameMatch.Err.MakeOneStringReplace(strconv.Itoa(rule.ConditionPeople) + " > " + strconv.Itoa(lenPlayers))
		return group, gameMatch.Err.NewReplace(410, errMsgMap)
	}
	//这里有个小问题，目前仅支持：每个组最多5人，回头我再优化
	if lenPlayers > gameMatch.Option.RuleTeamMaxPeople {
		return group, gameMatch.Err.New(411)
	}

	//验证3方传过来的groupId 是否重复
	allGroupIds := rule.QueueSign.GetGroupSignTimeoutAll()
	for _, hasGroupId := range allGroupIds {
		if int(outGroupId) == hasGroupId {
			errMsg := gameMatch.Err.MakeOneStringReplace(strconv.Itoa(int(outGroupId)))
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
		player, isEmpty := rule.PlayerManager.GetById(int(httpPlayer.Uid))
		//queueSign.Log.Info("player(" + strconv.Itoa(player.Id) + ") GetById :  status = " + strconv.Itoa(playerStatusElement.Status) + " isEmpty:" + strconv.Itoa(isEmpty))
		if isEmpty == 1 {
			//这是正常，用户之前没有登陆过，或 玩过一次，结算的时候把该数据清了
			newPlayer := rule.PlayerManager.createEmptyPlayer()
			newPlayer.Id = int(httpPlayer.Uid)
			player = newPlayer
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
		player.WeightAttrs = httpPlayer.WeightAttr
		playerList = append(playerList, player)
	}
	gameMatch.Option.Log.Info("finish check player status ，start calculate group/player weight:")

	//先计算一下权重平均值
	var groupWeightTotal float32
	groupWeightTotal = 0.00

	if rule.Formula != "" {
		//util.MyPrint(queueSign, "rule weight , Formula : ", rule.PlayerWeight.Formula)
		var weight float32
		weight = 0.00 //一个组的总权重值
		var playerWeightValue []float32
		for k, p := range playerList {
			onePlayerWeight := gameMatch.getPlayerWeightByFormula(rule.Formula, p.WeightAttrs)
			if onePlayerWeight > float32(gameMatch.Option.WeightMaxValue) {
				onePlayerWeight = float32(gameMatch.Option.WeightMaxValue)
			}
			weight += onePlayerWeight
			playerWeightValue = append(playerWeightValue, onePlayerWeight)
			playerList[k].Weight = onePlayerWeight
		}
		//util.MyPrint("rule.WeightTeamAggregation:", rule.WeightTeamAggregation)
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
			//util.MyPrint("len(playerList):", len(playerList), " weight:", weight)
			groupWeightTotal = weight / float32(len(playerList))
		}
		//保留2位小数
		tmp := util.FloatToString(groupWeightTotal, 2)
		groupWeightTotal = util.StringToFloat(tmp)
		//util.MyPrint("groupWeightTotal:", groupWeightTotal, " tmp ", tmp)
	} else {
		gameMatch.Option.Log.Info("rule.Formula empty , no need calculate group/player weight.")
	}

	//这里有个特殊的情况:报名(人数)即满足条件
	//按说也不应该进到我这儿，直接调用同步服务即可，可能为了统一走一个服务，且调用方不想麻烦再写代码了
	//此方法放弃，如果单独写，得弄两套代码，但逻辑差不多，而一边改了，另一边忘改就出BUG，还是正常走匹配逻辑，正常的产生数据，正常的推送
	//if rule.Type == service.RULE_TYPE_TEAM_EACH_OTHER {
	//	if lenPlayers == rule.ConditionPeople {
	//		go rule.Push.RuntimeSuccess(form, ruleId)
	//		return
	//	}
	//}

	//这里再做一次检查，防止，此时某个 rule 关闭了
	if rule.Status != service.GAME_MATCH_RULE_STATUS_EXEC {
		return group, gameMatch.Err.New(400)
	}

	//验证都成功了，下面开始处理具体的添加组、添加用户的操作，因为是强依赖REDIS，用redis当DB
	//所以必须是原子操作。否则数据一但更新不全，某些用户可能就永远不能再匹配了......

	//这里有偷个懒，还是用外部的groupId , 不想再给redis加 groupId映射outGroupId了
	//gameMatch.Option.Log.Warn(" outGroupId replace groupId :" + strconv.Itoa(outGroupId) + " " + strconv.Itoa(group.Id))

	redisConnFD := gameMatch.Option.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()
	gameMatch.Option.Redis.Multi(redisConnFD) //开始多指令缓存模式
	//匹配超时时间
	expire := processStartTime + rule.MatchTimeout
	//创建一个新的小组
	group = gameMatch.NewGroupStruct(rule)
	group.Id = int(outGroupId)
	group.Players = playerList
	group.Type = service.GAME_MATCH_GROUP_TYPE_SIGN
	group.SignTimeout = expire
	group.SignTime = processStartTime
	group.Person = len(playerList)
	group.Weight = groupWeightTotal
	group.OutGroupId = int(outGroupId)
	group.Addition = form.Addition
	//group.CustomProp = httpReqBusiness.CustomProp
	//group.MatchCode = rule.CategoryKey
	rule.QueueSign.AddOne(group, redisConnFD)
	groupBytes, _ := json.Marshal(&group)
	gameMatch.Option.Log.Info("add one group:" + string(groupBytes))
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
		rule.PlayerManager.delOneById(redisConnFD, player.Id)
		rule.PlayerManager.Create(newPlayerStatusElement, redisConnFD)

		playerIds += strconv.Itoa(player.Id) + ","
	}

	//持久化,报名记录
	gameMatch.PersistenceRecordGroup(group, rule.Id)

	//提交事务(缓存中的redis指令)
	_, err = gameMatch.Option.Redis.Exec(redisConnFD)
	//if err != nil {
	//	queueSign.Log.Error("transaction failed : " + err.Error())
	//}
	gameMatch.Option.Log.Info(" sign finish ,total : newGroupId " + strconv.Itoa(group.Id) + " , players len : " + strconv.Itoa(len(playerList)))

	return group, nil
}

func (gameMatch *GameMatch) Cancel(form pb.GameMatchPlayerCancel) error {
	if form.RuleId <= 0 {
		return errors.New("rule id empty")
	}

	if form.GroupId <= 0 {
		return errors.New("GroupId empty")
	}

	rule, err := gameMatch.RuleManager.GetById(int(form.RuleId))
	if err != nil {
		return err
	}
	return rule.QueueSign.CancelByGroupId(int(form.GroupId))
}

//注： formula 不支持小数点，变量用尖括号：( <age> * 20 ) + ( <level> * 50)
func (gameMatch *GameMatch) getPlayerWeightByFormula(formula string, MatchAttr map[string]int32) float32 {
	//mylog.Debug("getPlayerWeightByFormula , formula:",formula)
	grep := gameMatch.Option.FormulaFirst + "([\\s\\S]*?)" + gameMatch.Option.FormulaEnd
	var imgRE = regexp.MustCompile(grep)
	findRs := imgRE.FindAllStringSubmatch(formula, -1)
	//util.MyPrint("parse PlayerWeightByFormula : ", findRs)
	if len(findRs) == 0 {
		return 0
	}
	for _, v := range findRs {
		val, ok := MatchAttr[v[1]]
		if !ok {
			val = 0
		}
		formula = strings.Replace(formula, v[0], strconv.Itoa(int(val)), -1)

	}
	//util.MyPrint("final formula replaced str :", formula)
	rs, err := util.Eval(formula)
	//util.MyPrint(rs, err)
	if err != nil {
		return 0
	}
	f, _ := rs.Float32()
	return f
}

func (gameMatch *GameMatch) DebugShowQueueInfo(rule *Rule, form pb.GameMatchSign) {

	groupsTotal := rule.QueueSign.getAllGroupsWeightCnt() //报名 小组总数
	playersTotal := rule.QueueSign.getAllPlayersCnt()     //报名 玩家总数
	gameMatch.Option.Log.Debug("ShowQueueInfo , groupId:" + strconv.Itoa(int(form.GroupId)) + " , playersLen : " + strconv.Itoa(len(form.PlayerList)) + " ,queue cnt : groupsTotal" + strconv.Itoa(groupsTotal) + " , playersTotal" + strconv.Itoa(playersTotal))
}
