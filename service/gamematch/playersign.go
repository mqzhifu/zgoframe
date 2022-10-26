package gamematch

import (
	"regexp"
	"strconv"
	"strings"
	"zgoframe/service"
	"zgoframe/util"
)

//报名 - 加入匹配队列
//此方法，得有前置条件：验证所有参数是否正确，因为使用者为http请求，数据的验证交由HTTP层做处理，如果是非HTTP，要验证一下
func (gamematch *Gamematch) Sign(httpReqBusiness HttpReqBusiness) (group Group, err error) {
	ruleId := httpReqBusiness.RuleId
	outGroupId := httpReqBusiness.GroupId
	//这里只做最基础的验证，前置条件是已经在HTTP层作了验证
	rule, ok := gamematch.RuleConfig.GetById(ruleId)
	if !ok {
		return group, myerr.New(400)
	}
	lenPlayers := len(httpReqBusiness.PlayerList)
	if lenPlayers == 0 {
		return group, myerr.New(401)
	}
	//groupsTotal := queueSign.getAllGroupsWeightCnt()	//报名 小组总数
	//playersTotal := queueSign.getAllPlayersCnt()	//报名 玩家总数
	//mylog.Info(" action :  Sign , players : " + strconv.Itoa(lenPlayers) +" ,queue cnt : groupsTotal",groupsTotal ," , playersTotal",playersTotal)
	queueSign := gamematch.GetContainerSignByRuleId(ruleId)
	mylog.Info("new sign :[ruleId : " + strconv.Itoa(ruleId) + "(" + rule.CategoryKey + ") , outGroupId : " + strconv.Itoa(outGroupId) + " , playersCount : " + strconv.Itoa(lenPlayers) + "] ")
	//mylog.Info("new sign :[ruleId : ",ruleId,"(",rule.CategoryKey,") , outGroupId : ",outGroupId," , playersCount : ",lenPlayers,"] ")
	//queueSign.Log.Info("new sign :[ruleId : " ,  ruleId   ,"(",rule.CategoryKey,") , outGroupId : ",outGroupId," , playersCount : ",lenPlayers,"] ")
	now := util.GetNowTimeSecondToInt()

	util.PrintStruct(queueSign.Rule, ":")
	if rule.Flag == service.RuleFlagTeamVS {
		if lenPlayers > queueSign.Rule.GroupPersonMax {
			msg := make(map[int]string)
			msg[0] = " rule.RuleFlagTeamVS " + strconv.Itoa(queueSign.Rule.GroupPersonMax)
			return group, myerr.NewReplace(408, msg)
		}
	} else {
		if lenPlayers > queueSign.Rule.PersonCondition {
			msg := make(map[int]string)
			msg[0] = strconv.Itoa(queueSign.Rule.PersonCondition) + " > " + strconv.Itoa(lenPlayers)
			return group, myerr.NewReplace(410, msg)
		}

		if lenPlayers > 5 && lenPlayers != rule.PersonCondition {
			return group, myerr.New(411)
		}

		//if lenPlayers > queueSign.Rule.GroupPersonMax &&  lenPlayers != rule.PersonCondition {
		//	msg := make(map[int]string)
		//	msg[0] = " rule.PersonCondition " + strconv.Itoa( queueSign.Rule.GroupPersonMax)
		//	return group,myerr.NewErrorCodeReplace(408,msg)
		//}
	}

	queueSign.Log.Info("start check player status :")
	//检查，所有玩家的状态
	var players []Player
	for _, httpPlayer := range httpReqBusiness.PlayerList {
		player := Player{Id: httpPlayer.Uid, MatchAttr: httpPlayer.MatchAttr}
		playerStatusElement, isEmpty := playerStatus.GetById(player.Id)
		queueSign.Log.Info("player(" + strconv.Itoa(player.Id) + ") GetById :  status = " + strconv.Itoa(playerStatusElement.Status) + " isEmpty:" + strconv.Itoa(isEmpty))
		if isEmpty == 1 {
			//这是正常
		} else if playerStatusElement.Status == service.PlayerStatusSuccess { //玩家已经匹配成功，并等待开始游戏
			queueSign.Log.Error(" player status = PlayerStatusSuccess ,demon not clean.")
			msg := make(map[int]string)
			msg[0] = strconv.Itoa(player.Id)
			return group, myerr.NewReplace(403, msg)
		} else if playerStatusElement.Status == service.PlayerStatusSign { //报名成功，等待匹配
			isTimeout := playerStatus.checkSignTimeout(rule, playerStatusElement)
			if !isTimeout { //未超时
				//queueSign.Log.Error(" player status = matching...  not timeout")
				msg := make(map[int]string)
				msg[0] = strconv.Itoa(player.Id)
				return group, myerr.NewReplace(402, msg)
			} else { //报名已超时，等待后台守护协程处理
				//这里其实也可以先一步处理，但是怕与后台协程冲突
				//queueSign.Log.Error(" player status is timeout ,but not clear , wait a moment!!!")
				msg := make(map[int]string)
				msg[0] = strconv.Itoa(player.Id)
				return group, myerr.NewReplace(407, msg)
			}
		}
		players = append(players, player)
		//playerStatusElementMap[player.Id] = playerStatusElement
	}
	mylog.Info("finish check player status.")
	//验证3方传过来的groupId 是否重复
	allGroupIds := queueSign.GetGroupSignTimeoutAll()
	for _, hasGroupId := range allGroupIds {
		if outGroupId == hasGroupId {
			util.MyPrint(allGroupIds, outGroupId, hasGroupId, httpReqBusiness)
			msg := make(map[int]string)
			msg[0] = strconv.Itoa(outGroupId)
			return group, myerr.NewReplace(409, msg)
		}
	}
	//这里有个特殊的情况 ，报名的人数即满足条件，具体啥需求不知道，按说也不应该进到我这儿，既然知道已然满足了还要匹配干毛？勉强给做了吧。
	if rule.Flag == service.RuleFlagCollectPerson {
		if lenPlayers == rule.PersonCondition {
			push := gamematch.getContainerPushByRuleId(ruleId)
			go push.RuntimeSuccess(httpReqBusiness, ruleId, gamematch)
			return
		}
	}
	////按说这里是重复的，但是为了兼容上面这个功能，无奈再判断一次吧
	//if lenPlayers > queueSign.Rule.GroupPersonMax{
	//	msg := make(map[int]string)
	//	msg[0] = strconv.Itoa( queueSign.Rule.GroupPersonMax)
	//	return group,myerr.NewErrorCodeReplace(408,msg)
	//}

	//zlib.ExitPrint(allGroupIds)
	//先计算一下权重平均值
	var groupWeightTotal float32
	groupWeightTotal = 0.00

	if rule.PlayerWeight.Formula != "" {
		util.MyPrint(queueSign, "rule weight , Formula : ", rule.PlayerWeight.Formula)
		var weight float32
		weight = 0.00
		var playerWeightValue []float32
		for k, p := range players {
			onePlayerWeight := getPlayerWeightByFormula(rule.PlayerWeight.Formula, p.MatchAttr, queueSign)
			util.MyPrint("onePlayerWeight : ", onePlayerWeight)
			if onePlayerWeight > service.WeightMaxValue {
				onePlayerWeight = service.WeightMaxValue
			}
			weight += onePlayerWeight
			playerWeightValue = append(playerWeightValue, onePlayerWeight)
			players[k].Weight = onePlayerWeight
		}
		switch rule.PlayerWeight.Aggregation {
		case "sum":
			groupWeightTotal = weight
		case "min":
			groupWeightTotal = util.FindMinNumInArrFloat32(playerWeightValue)
		case "max":
			groupWeightTotal = util.FindMaxNumInArrFloat32(playerWeightValue)
		case "average":
			groupWeightTotal = weight / float32(len(players))
		default:
			groupWeightTotal = weight / float32(len(players))
		}
		//保留2位小数
		tmp := util.FloatToString(groupWeightTotal, 2)
		groupWeightTotal = util.StringToFloat(tmp)
	} else {
		util.MyPrint(queueSign, "rule weight , Formula is empty!!!")
	}
	//下面两行必须是原子操作，如果pushOne执行成功，但是upInfo没成功会导致报名队列里，同一个用户能再报名一次
	redisConnFD := myredis.GetNewConnFromPool()
	defer redisConnFD.Close()
	//开始多指令缓存模式
	myredis.Multi(redisConnFD)

	//超时时间
	expire := now + rule.MatchTimeout
	//创建一个新的小组
	group = gamematch.NewGroupStruct(rule)
	//这里有偷个懒，还是用外部的groupId , 不想再给redis加 groupId映射outGroupId了
	mylog.Warn(" outGroupId replace groupId :" + strconv.Itoa(outGroupId) + " " + strconv.Itoa(group.Id))
	group.Id = outGroupId
	group.Players = players
	group.SignTimeout = expire
	group.Person = len(players)
	group.Weight = groupWeightTotal
	group.OutGroupId = outGroupId
	group.Addition = httpReqBusiness.Addition
	group.CustomProp = httpReqBusiness.CustomProp
	group.MatchCode = rule.CategoryKey
	util.MyPrint(queueSign, "newGroupId : ", group.Id, "player/group weight : ", groupWeightTotal, " now : ", now, " expire : ", expire)
	//mylog.Info("newGroupId : ",group.Id , "player/group weight : " ,groupWeightTotal ," now : ",now ," expire : ",expire )
	//queueSign.Log.Info("newGroupId : ",group.Id , "player/group weight : " ,groupWeightTotal ," now : ",now ," expire : ",expire)
	queueSign.AddOne(group, redisConnFD)
	playerIds := ""
	for _, player := range players {

		newPlayerStatusElement := playerStatus.newPlayerStatusElement()
		newPlayerStatusElement.PlayerId = player.Id
		newPlayerStatusElement.Status = service.PlayerStatusSign
		newPlayerStatusElement.RuleId = ruleId
		newPlayerStatusElement.Weight = player.Weight
		newPlayerStatusElement.SignTimeout = expire
		newPlayerStatusElement.GroupId = group.Id

		queueSign.Log.Info("playerStatus.upInfo:" + strconv.Itoa(service.PlayerStatusSign))
		playerStatus.upInfo(newPlayerStatusElement, redisConnFD)

		playerIds += strconv.Itoa(player.Id) + ","
	}
	//提交缓存中的指令
	_, err = myredis.Exec(redisConnFD)
	if err != nil {
		queueSign.Log.Error("transaction failed : " + err.Error())
	}
	queueSign.Log.Info(" sign finish ,total : newGroupId " + strconv.Itoa(group.Id) + " success players : " + strconv.Itoa(len(players)))
	mylog.Info(" sign finish ,total : newGroupId " + strconv.Itoa(group.Id) + " success players : " + strconv.Itoa(len(players)))

	//signSuccessReturnData = SignSuccessReturnData{
	//	RuleId: ruleId,
	//	GroupId: outGroupId,
	//	PlayerIds: playerIds,
	//
	//}

	return group, nil
}

func getPlayerWeightByFormula(formula string, MatchAttr map[string]int, sign *QueueSign) float32 {
	//mylog.Debug("getPlayerWeightByFormula , formula:",formula)
	grep := service.FormulaFirst + "([\\s\\S]*?)" + service.FormulaEnd
	var imgRE = regexp.MustCompile(grep)
	findRs := imgRE.FindAllStringSubmatch(formula, -1)
	util.MyPrint(sign, "parse PlayerWeightByFormula : ", findRs)
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
	util.MyPrint(sign, "final formula replaced str :", formula)
	rs, err := util.Eval(formula)
	if err != nil {
		return 0
	}
	f, _ := rs.Float32()
	return f
}
