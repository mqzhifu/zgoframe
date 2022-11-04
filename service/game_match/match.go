package gamematch

import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
	"zgoframe/service"
	"zgoframe/util"
)

//匹配 (这个就是整个包的核心)
type Match struct {
	Mutex      sync.Mutex
	Rule       *Rule
	rangeStart string
	rangeEnd   string
	Log        *zap.Logger //log 实例
	Redis      *util.MyRedisGo
	Err        *util.ErrMsg
	CloseChan  chan int
	prefix     string
}

func NewMatch(rule *Rule) *Match {
	match := new(Match)
	match.Rule = rule
	match.rangeStart = ""
	match.rangeEnd = ""
	match.Redis = rule.RuleManager.Option.GameMatch.Option.Redis
	match.Log = rule.RuleManager.Option.GameMatch.Option.Log
	match.CloseChan = make(chan int)
	match.prefix = "match"
	return match
}
func (match *Match) Close() {
	match.CloseChan <- 1
}

//守护 协程
func (match *Match) Demon() {
	match.Log.Info(match.prefix + " Demon.")
	for {
		select {
		case signal := <-match.CloseChan:
			match.Log.Warn(match.prefix + "Demon CloseChan receive :" + strconv.Itoa(signal))
			goto forEnd
		default:
			match.matching()
			time.Sleep(time.Millisecond * time.Duration(match.Rule.RuleManager.Option.GameMatch.LoopSleepTime))
		}
	}
forEnd:
	match.Log.Warn("match Demon end .")
}

/*
一次完整的匹配大流程
1. 先检查是否有报名超时的情况，并做超时处理
1. 获取一下当前池子里：共有多少个组，共有多少个玩家，如果太少，直接结束
2. 根据 rule 是否有公式(权重)，分成了两类处理方式：
	1. 没有公式(权重)，这种比较简单，直接暴力去池子里捞用户，就一个维度：谁先报名的，谁最快匹配
	2. 有公式(权重)，这种略复杂。先计算出一个步长，然后根据此步长循环，每次有一个权限范围值，根据这个范围，去池子里面捞用户
	>步长的计算方式：rule 配置表里的 权重最大值 - 权重最小值
3. 将捞出来的用户，做进一步细致的匹配，这一步才是略核心的。这里分为两类吧，但最终会合并
	1. N VS N ，这种是略复杂一点的，它多出了步，即：计算互补数，如：5人一队，优先匹配对策也是5人，3+2人一队，优先匹配3+2的组队。
	2. 吃鸡模式，这个比较粗暴，略简单，就正常匹配即可。
4. 做排队组合：计算出  哪个队的人数 + 哪个组的人数 = 符合条件的一个结果
5. 将捞出来的用户数据 配合 排队组合的结合 ，计算出，最终：匹配成功的结果
6. 格式化数据，将匹配成功的结果，创建一条记录：success result,再添加子组数据：group ，再添加超时数据，再添加push数据



对于 如何捞用户的计算：
1. 先，从最大的范围搜索，扫整个集合，如果元素过少，直接在这个维度就结束了
2. 缩小搜索范围，把整个集合，划分成10个维度，每个维度单纯计算，如果成功，那就结束了，如果人数过多，还会再划分一次最细粒度的匹配
3. 这种是介于上面两者中间~即不是全部集合，也不是单独计算一个维度，而是逐步，放大到:最大90%集合，1-1，1-2....1-9

总结：以上的算法，其实就是不断切换搜索范围（由最大到中小，再到中大），加速匹配时间

*/
func (match *Match) matching() {
	//每次匹配之前，要先检查一下数据是否超时，后面就不再检查了
	match.Rule.QueueSign.CheckTimeout()

	playersTotal := match.Rule.QueueSign.getAllPlayersCnt()
	groupsTotal := match.Rule.QueueSign.getAllGroupsWeightCnt()
	match.clearMemberRange()

	now := util.GetNowTimeSecondToInt()
	if playersTotal == 0 || groupsTotal == 0 {
		//match.Log.Debug(" first total is empty ")
		if now%match.Rule.DemonDebugTime == 0 {
			match.Log.Info("new times <matching> ,ruleId: " + strconv.Itoa(match.Rule.Id))
			match.Log.Info(match.prefix + " matching total is empty ")
		}

		return
	}
	match.Log.Info("new once matching func , playersTotal total:" + strconv.Itoa(playersTotal) + " groupsTotal : " + strconv.Itoa(groupsTotal))
	//设置了匹配权重公式
	if match.Rule.Formula != "" && match.Rule.WeightScoreMax > match.Rule.WeightScoreMin {
		match.Log.Info("if case in :PlayerWeight.Formula")
		dead := 0
		//步长，两个<极值>的距离，用于循环
		distance := match.Rule.WeightScoreMax - match.Rule.WeightScoreMin
		match.Log.Info("case in Formula , distance:" + strconv.Itoa(distance))
		//这里，比较好的循环值应该是0-100，保证每走向前走一步，上下值的范围都能获取全了
		//但是，这样有点浪费，步长设置成  上面的距离值，更快一些
		for i := 0; i < match.Rule.RuleManager.Option.GameMatch.WeightMaxValue; i = i + distance {
			start := i - match.Rule.WeightScoreMin
			if start < 0 {
				start = 0
			}

			end := i + match.Rule.WeightScoreMax
			if end > match.Rule.RuleManager.Option.GameMatch.WeightMaxValue {
				end = match.Rule.RuleManager.Option.GameMatch.WeightMaxValue
				dead = 1
			}

			if dead == 1 {
				break
			}
			tmpMax := float32(end)
			rangeStart := strconv.Itoa(start)
			rangeEnd := util.FloatToString(tmpMax, 3)

			match.setMemberRange(rangeStart, rangeEnd)
			//successGroupIds: map[int]map[int]int   ,外层的map 是证明一次匹配成功的结果，内层map 是 groupIds
			successGroupIds, isEmpty := match.searchByRange(service.FilterFlagDIY)
			match.Log.Info("searchByRange rs : " + rangeStart + " ~ " + rangeStart + " , rs , cnt:" + strconv.Itoa(len(successGroupIds)) + " isEmpty : " + strconv.Itoa(isEmpty))

			util.MyPrint("searchByRange rs : ", successGroupIds)
			if isEmpty == 1 {
				continue
			}
			//match.Log.Debug("successGroupIds",successGroupIds)
			util.MyPrint("searchByRange rs : ", successGroupIds)
			//match.Log.Info("successGroupIds",successGroupIds)
			//将计算好的组ID，团队，插入到 成功 队列中
			if len(successGroupIds) > 0 {
				match.successConditions(successGroupIds)
			}
		}
		//zlib.ExitPrint(11111)
	} else {
		match.Log.Info("case in  no Formula:")
		//matchingRange := []int{service.FilterFlagAll, service.FilterFlagBlock, service.FilterFlagBlockInc}
		//for i := 0; i < len(matchingRange); i++ {
		//	groupIds, isEmpty := match.matchingRange(matchingRange[i])
		//	//groupIds, isEmpty := match.matchingRange(service.FilterFlagAll)
		//	if isEmpty == 1 {
		//		match.Log.Warn("signed players is empty ")
		//		return
		//	}
		//	if matchingRange[i] == service.FilterFlagAll {
		//		if len(groupIds) > 0 {
		//			match.Log.Info("FilterFlagAll hit...")
		//			return
		//		}
		//	}
		//}

		//rule 配置表中：没有设置公式，即:用户都没有权重属性，那么匹配的时候，间接等于就一个维度了：报名的时间，报名早的可能就被优先匹配了
		groupIds, isEmpty := match.matchingRange(service.FilterFlagAll)
		if isEmpty == 1 {
			match.Log.Warn("signed players is empty ")
			return
		}
		//if matchingRange[i] == service.FilterFlagAll {
		if len(groupIds) > 0 {
			match.Log.Info("FilterFlagAll hit...")
			return
		}

	}

	finalPlayersTotal := match.Rule.QueueSign.getAllPlayersCnt()
	finalGroupsTotal := match.Rule.QueueSign.getAllGroupsWeightCnt()

	match.Log.Info("once matching end, playersTotal total:" + strconv.Itoa(playersTotal) + " groupsTotal : " + strconv.Itoa(groupsTotal) +
		"finalPlayersTotal total:" + strconv.Itoa(finalPlayersTotal) + " finalGroupsTotal : " + strconv.Itoa(finalGroupsTotal))

	match.clearMemberRange()
}

//用于输出测试
func (match *Match) SuccessGroupIdsToStr(successGroupIds map[int]map[int]int) string {
	if len(successGroupIds) <= 0 {
		return " SuccessGroupIdsToStr empty"
	}
	str := ""
	for k1, groupList := range successGroupIds {
		rowStr := strconv.Itoa(k1) + "\n"
		for k2, v := range groupList {
			rowStr += util.GetSpaceStr(4) + "k:" + strconv.Itoa(k2) + " groupId:" + strconv.Itoa(v)
		}
		str += rowStr
	}

	return str
}

//用于输出测试
func (match *Match) GroupByPersonCntToStr(list map[int]int) string {
	if len(list) <= 0 {
		return " GroupByPersonCntToStr empty"
	}
	str := ""
	for k1, groupNum := range list {
		str += "person:" + strconv.Itoa(k1) + " groupNum:" + strconv.Itoa(groupNum) + "\n"
	}
	return str
}

/*
flag
	1:平均10等份
	2：逐渐 递增份
	3：全部
PS：这个方法好像没太大用处，权限匹配，走的是自己的测试（根据步长DIY筛选），而没有权重公式的，进了些方法，也只能使用全匹配模式，因为这个方法的匹配方式是基于权限的，没权限，你就是用 递增 它也没效果啊
*/
//范围性匹配：搜索一定<权重>范围内的数据，一共是3次，由最粗到最细粒度
func (match *Match) matchingRange(flag int) (successGroupIds map[int]map[int]int, isEmpty int) {
	var tmpMin int
	var tmpMax float32
	rangeStart := "" //权限范围的起始值
	rangeEnd := ""   //权限范围的结束值
	forStart := 0    //循环的开始值
	//循环的结束值
	forEnd := match.Rule.RuleManager.Option.GameMatch.WeightMaxValue
	if flag == service.FilterFlagAll { //全匹配模式下就不需要多次循环了，一次即可
		forEnd = 1
	}

	match.Log.Info("matchingRange , flag : " + strconv.Itoa(flag) + " forEnd:" + strconv.Itoa(forEnd))

	for i := forStart; i < forEnd; i++ {
		if flag == service.FilterFlagBlock { //区域性的
			tmpMin = i
			tmpMax = float32(i)*10 + 0.9999

			rangeStart = strconv.Itoa(tmpMin)
			rangeEnd = util.FloatToString(tmpMax, 3)
		} else if flag == service.FilterFlagBlockInc { //递增，范围更大
			if i == 0 {
				continue
			}
			if i == 1 {
				//因为0-1 区间上面已经计算过了，所以 从1开始
				continue
			}
			if i == 9 {
				//0-10 最开始全范围搜索已民经计算过了
				continue
			}
			tmpMin = 0
			tmpMax = float32(i)*10 + 0.9999

			rangeStart = strconv.Itoa(tmpMin)
			rangeEnd = util.FloatToString(tmpMax, 3)
		} else { //全部
			rangeStart = "-inf"
			rangeEnd = "+inf"
		}
		match.setMemberRange(rangeStart, rangeEnd)
		successGroupIds, isEmpty = match.searchByRange(flag)
		//match.Log.Info("searchByRange rs : ",rangeStart , " ~ ",rangeStart , " , rs , cnt:",len(successGroupIds) , "isEmpty : ",isEmpty)
		match.Log.Info("searchByRange rs : " + rangeStart + " ~ " + rangeStart + " , rs , cnt:" + strconv.Itoa(len(successGroupIds)) + "isEmpty : " + strconv.Itoa(isEmpty))
		match.Log.Debug("searchByRange rs successGroupIds : " + match.SuccessGroupIdsToStr(successGroupIds))
		if isEmpty == 1 && flag == service.FilterFlagAll {
			//最大范围的查找没有数据，并且从redis读不出数据，证明就是数据太少，不可能成团
			match.Log.Warn(" FilterFlagAll  isEmpty  = 1 ,break this foreach")
			match.Log.Info(" FilterFlagAll  isEmpty  = 1 ,break this foreach")
			return successGroupIds, isEmpty
		}
		//match.Log.Debug("successGroupIds",successGroupIds)
		//util.MyPrint("successGroupIds", successGroupIds)
		//将计算好的组ID，团队，插入到 成功 队列中
		if len(successGroupIds) > 0 {
			match.Log.Info("has success len : " + strconv.Itoa(len(successGroupIds)) + "  and start insert : ")
			match.successConditions(successGroupIds)
			//if flag == service.FilterFlagAll {
			//	//如果最大范围的匹配命中了数据，就证明，只有在最大范围才有数据，再细化的范围搜索已无用
			//	return successGroupIds, isEmpty
			//}
		}
	}

	return successGroupIds, isEmpty
}

//根据具体的<权重>值，去池子里面捞用户
//这里主要还是：查询 redis 报名池里的 总量数据，之后对总量数据进行基础检验，并没有实际的去捞用户与计算
func (match *Match) searchByRange(flag int) (successGroupIds map[int]map[int]int, isEmpty int) {
	playersTotal := 0                  //当前池子里共有多少个玩家
	if flag == service.FilterFlagAll { //全匹配
		playersTotal = match.Rule.QueueSign.getAllPlayersCnt()
	} else { //块匹配
		playersTotal = match.Rule.QueueSign.getPlayersCntTotalByWeight(match.rangeStart, match.rangeEnd)
	}
	//当前池子里的小组数
	groupsTotal := match.Rule.QueueSign.getGroupsWeightCnt(match.rangeStart, match.rangeEnd)
	match.Log.Info("searchByRange : " + match.rangeStart + " , " + match.rangeEnd + "  , flag " + strconv.Itoa(flag) + " playersTotal:" + strconv.Itoa(playersTotal) + " groupsTotal : " + strconv.Itoa(groupsTotal))
	//玩家数 或 小组数 <=0 没必要再算了，直接返回
	if playersTotal <= 0 || groupsTotal <= 0 {
		match.Log.Warn("searchByRange : playersTotal or groupsTotal  is 0")
		return successGroupIds, 1
	}
	//组成一个队的：最低人数要求 (好像是错的)
	//一次匹配成功：需要的总人数
	personCondition := match.Rule.ConditionPeople
	//if match.Rule.Type == service.RULE_TYPE_TEAM_VS { //组队/对战类
	//	personCondition = match.Rule.TeamMaxPeople * 2
	//}
	match.Log.Info("searchByRange success condition when person =" + strconv.Itoa(personCondition))
	if playersTotal < personCondition {
		match.Log.Warn("searchByRange : playersTotal(" + strconv.Itoa(playersTotal) + ") < personCondition (" + strconv.Itoa(personCondition) + ") ")
		return successGroupIds, 1
	}
	//因为：最先做的就是 <全匹配模式> ，但略有些粗暴，如果此时待匹配的玩家过多，最好还是留给后面的 细匹配来处理
	//这里以 满足条件人数的 4倍 做为一个定值
	//if flag == service.FilterFlagAll && playersTotal > personCondition*4 {
	//	match.Log.Warn(" FilterFlagAll playersTotal > personCondition , end this searchByRange ")
	//	return successGroupIds, 0
	//}
	//上面仅仅是从 redis 中查看了一下数据，并没有实际操作，更偏向从总数来简单分析下数据情况，没有问题后~,开始详情的计算
	successGroupIds = match.searchFilterDetail()
	return successGroupIds, 0
}

//前面是基于各种<总数>汇总，的各种验证，都没有问题了，这里才算是正式的：细节匹配
func (match *Match) searchFilterDetail() (defaultSuccessGroupIds map[int]map[int]int) {
	match.Log.Info("searchFilterDetail , rule flag :" + strconv.Itoa(service.RULE_TYPE_TEAM_VS))
	//使用地址引用，后面的子函数，不用返回值了
	successGroupIds := make(map[int]map[int]int)       //保留最终值：匹配成功的情况。注：N VS N 保留的结果维护：是组队成5人组，而吃鸡模式的一个维度的值就代表是完全匹配成功一次
	TeamVSSuccessGroupIds := make(map[int]map[int]int) //用于计算N VS N 的匹配成功结果，它是先计算出一部分的值，最终会与 successGroupIds 合并
	//N vs N ，会多一步，公平性匹配：5人组 匹配的对手 也是5人组 ，3+2 匹配的最好也是 3+2
	//但最终两种方法，还是会统一再处理一次 groupPersonCalculateNumberCombination，N VS N处理的一步只是公平性匹配，可能还剩下不少不公平的情况
	if match.Rule.Type == service.RULE_TYPE_TEAM_VS {
		/*
			共2轮过滤筛选
			1、只求互补数，为了公平：5人组~匹配，对手也是5人组队的，3人组队匹配到的对手里也应该有3人组队的
			   	如：5没有互补，4的互补是1，3的互补是2（2的互补是3，1互补是4，但前两步已经计算过了，这两步忽略）
			2、接上面处理过后的数据，正常排序组合
			   	上面一轮筛选过后，5人的：最多只剩下一个
				4人的，可能匹配光了，也可能剩下N个，因为取决于1人有多少，如果1人被匹配光了，4人组剩下多少个都没意义了，因为无法再组成一个团队了，所以4人组可直接忽略
				剩下 3 2 1 ，没啥快捷的办法了，只能单纯排列组合=5
		*/
		match.logarithmic(TeamVSSuccessGroupIds)
	}
	//注： 上面一但执行，虽然还是查询操作且主要是查询索引，但真实索引值已经出队列了，如果后面操作失败，要把索引补全，不然数据缺失
	match.groupPersonCalculateNumberCombination(successGroupIds)
	util.MyPrint("TeamVSSuccessGroupIds:", TeamVSSuccessGroupIds, " successGroupIds:", successGroupIds)
	//合并  TeamVSSuccessGroupIds  successGroupIds
	if len(TeamVSSuccessGroupIds) > 0 {
		if len(successGroupIds) > 0 {
			match.Log.Info("merge TeamVSSuccessGroupIds successGroupIds in case 1")
			index := len(successGroupIds)
			for _, v := range TeamVSSuccessGroupIds {
				successGroupIds[index] = v
				index++
			}
		} else {
			match.Log.Info("merge TeamVSSuccessGroupIds successGroupIds in case 2")
			successGroupIds = TeamVSSuccessGroupIds
		}
	}
	//因用的是map ,预先make 下，地址就已经分配了，可能有时，其实一个值都没有用到，但是len 依然是 > 0
	if util.CheckMap2IntIsEmpty(successGroupIds) {
		match.Log.Warn("searchFilterDetail is empty")
		return defaultSuccessGroupIds
	}

	return successGroupIds
}

//小组人数计算  数字组合方式
func (match *Match) groupPersonCalculateNumberCombination(successGroupIds map[int]map[int]int) {
	//这里吧，按说取，一条rule最大的值就行，没必要取全局最大的5，但是吧，后面的算法有点LOW，外层循环数就是5 ，除了矩阵，太麻烦，回头我再想想
	groupPersonNum := match.Rule.QueueSign.getPlayersCntByWeight(match.rangeStart, match.rangeEnd)
	match.Log.Info("groupPersonCalculateNumberCombination every group person total : " + match.GroupByPersonCntToStr(groupPersonNum))
	//根据组人数，做排列组合，得到最终自然数和
	condition := match.Rule.ConditionPeople
	if match.Rule.Type == service.RULE_TYPE_TEAM_VS { //吃鸡模式就是只要满足设定的条件即可，而N VS N 模式并不需要那么多的人匹配，它只需要匹配出一个队伍即可
		condition = match.Rule.ConditionPeople / 2
	}
	calculateNumberTotalRs := match.calculateNumberTotal(condition, groupPersonNum)
	//上面的函数，虽然计算完了，总人数是够了，但是也可能不会成团，比如：全是4人为一组报名，成团总人数却是奇数，偶数和是不可能出现奇数的
	if len(calculateNumberTotalRs) == 0 {
		//match.Log.Notice(zlib.GetSpaceStr(4)+"calculateNumberTotal is empty")
		match.Log.Warn("calculateNumberTotal is empty")
		return
	}
	/*
		calculateNumberTotal 函数是排列组合，它关心的是 N个数的总和，一个数可以使用N次，重复使用。而匹配：如果一个数使用了一次，下次组合就得少一个
		所以排列组合的结果，会有很多，而实际结果却会少很多
		接着分析，对两种匹配机制
		1、只要够多少人即成功，那么排列一次，随机，取其中一个算法，下一次，再排列一下
		2、N VS N PK，

	*/

	//这里是倒序，因为最大的数在最后，如5人组成立的条件，使用5最多
	inc := len(successGroupIds)
	for i := len(calculateNumberTotalRs) - 1; i >= 0; i-- {

		//zlib.MyPrint(zlib.GetSpaceStr(4)+"condition " , i )
		oneConditionGroupIds := match.oneConditionConvertGroup(calculateNumberTotalRs[i])
		//successGroupIds[inc] = oneConditionGroupIds
		//inc++
		if len(oneConditionGroupIds) > 0 {
			successGroupIds[inc] = oneConditionGroupIds
			inc++
		} else {
			match.Log.Warn("this calculateNumberTotal condition not found ~")
		}
	}
	match.Log.Debug("groupPersonCalculateNumberCombination rs :" + match.SuccessGroupIdsToStr(successGroupIds))
}
func (match *Match) oneConditionConvertGroup(oneOkPersonCondition [5]int) map[int]int {
	oneConditionGroupIds := make(map[int]int)
	inc := 0
	//match.Log.Debug(oneOkPersonCondition)
	someOneEmpty := 0 //redis 里取不出来数据了，或者取出的数据 小于 应取数据个数
	for index, num := range oneOkPersonCondition {
		person := index + 1
		//zlib.MyPrint( zlib.GetSpaceStr(4)+" groupPerson : ",person , " get num : ",num )
		if num <= 0 {
			util.MyPrint(util.GetSpaceStr(4) + " get num <= 0  ,continue")
			continue
		}
		//从redis的GroupPersonIndex索引中 到内存中（redis的数据已被删除，这里不能断了，如果断了得把取出来的数据再塞回去，不然对不上)
		groupIds := match.Rule.QueueSign.getGroupPersonIndexList(person, match.rangeStart, match.rangeEnd, 0, num, true)
		if len(groupIds) == 0 {
			someOneEmpty = 1
			match.Log.Error("getGroupPersonIndexList empty")
			break
		}
		for i := 0; i < len(groupIds); i++ {
			oneConditionGroupIds[inc] = groupIds[i]
			inc++
		}

		if len(groupIds) != num {
			someOneEmpty = 1
			util.MyPrint("getGroupPersonIndexList != num "+strconv.Itoa(len(groupIds))+" "+strconv.Itoa(num)+" ", groupIds)
			break
		}
		//zlib.MyPrint("groupIds : ",groupIds)
	}
	if someOneEmpty == 1 {
		//match.Log.Error("oneConditionConvertGroup someOneEmpty  empty")
		//match.Log.Error("oneConditionConvertGroup someOneEmpty  empty")
		msg := "oneConditionConvertGroup someOneEmpty  empty"
		if len(oneConditionGroupIds) <= 0 {
			msg += " bug oneConditionGroupIds len = 0 ,no need pushBack redis info"
			match.Log.Error(msg)
		} else {
			msg += " oneConditionGroupIds len = " + strconv.Itoa(len(oneConditionGroupIds))
			match.Log.Error(msg)
			redisConnFD := match.Redis.GetNewConnFromPool()
			match.Redis.Multi(redisConnFD)
			match.groupPushBackCondition(redisConnFD, oneConditionGroupIds)
			match.Redis.Exec(redisConnFD)
			redisConnFD.Close()
			//重置：该变量
			oneConditionGroupIds = make(map[int]int)
		}

	}

	return oneConditionGroupIds
}

//N V N  ,求:互补数/对数
//这里只是做公平匹配，如：5V5 优先匹配出来，然后是4+1 VS 4+1 ，依此类推，保证5人组最好是直接匹配成5人组
func (match *Match) logarithmic(successGroupIds map[int]map[int]int) {
	//match.Log.Info(zlib.GetSpaceStr(3)+ " action RuleFlagTeamVS logarithmic :")
	prefix := "logarithmic "
	groupPersonNum := match.Rule.QueueSign.getPlayersCntByWeight(match.rangeStart, match.rangeEnd)
	match.Log.Debug(prefix + " rangeStart:" + match.rangeStart + "  ~ rangeEnd" + match.rangeEnd + " groupPersonNum , " + match.GroupByPersonCntToStr(groupPersonNum) + " ")

	successGroupIdsInc := 0   //最终成功的团队数
	var processedNumber []int //已处理过的互补数，如：4计算完了，1就不用算了，3计算完了，2其实也不用算了
	wishSuccessGroups := 0    //预计/希望 应该成功的  团队~  用于统计debug
	for personUnit, groupTotal := range groupPersonNum {
		match.Log.Info(prefix + "foreach , personUnit " + strconv.Itoa(personUnit) + " ,  groupTotal  " + strconv.Itoa(groupTotal) + " successGroupIdsInc" + strconv.Itoa(successGroupIdsInc))
		//判断是否已经处理过了
		elementInArrIndex := util.ElementInArrIndex(processedNumber, personUnit)
		if elementInArrIndex != -1 {
			continue
		}
		//5人 或 设置的最大值，已是最大的，直接就满足了，不需要处理互补数了，但得给它找到一同样为 5人的 小组
		if personUnit == match.Rule.TeamMaxPeople {
			match.Log.Info(util.GetSpaceStr(4) + "in max TeamVSPerson , no need remainder number")
			if groupTotal <= 1 { //<5人组>如果只有一个的情况，满足不了条件
				match.Log.Warn("in max TeamVSPerson 1 , but personTotal <= 1 , continue")
				continue
			}
			//当前值 与 需要的互补数，找出最小的一个
			maxNumber := match.getMinPersonNum(groupTotal, groupTotal)
			match.Log.Info("maxNumber " + strconv.Itoa(maxNumber))
			if maxNumber <= 0 {
				match.Log.Warn("in max TeamVSPerson 2 , but personTotal <= 1 , continue")
				continue
			}
			wishSuccessGroups += maxNumber / 2
			//取出(这里才是真的取，删除redis里的数据)集合中，所有人数为5的组group_ids.
			groupIds := match.Rule.QueueSign.getGroupPersonIndexList(match.Rule.TeamMaxPeople, match.rangeStart, match.rangeEnd, 0, maxNumber, true)
			//util.ExitPrint("groupIds:", groupIds)
			j := 0
			for i := 0; i < maxNumber/2; i++ {
				//tmp := make(map[int]int)
				//tmp[0] = groupIds[j]
				//j++
				//tmp[1] = groupIds[j]
				//j++
				//successGroupIds[successGroupIdsInc] = tmp
				//successGroupIdsInc++
				tmp := make(map[int]int)
				tmp[0] = groupIds[j]
				j++
				tmp2 := make(map[int]int)
				tmp2[0] = groupIds[j]
				j++
				successGroupIds[successGroupIdsInc] = tmp
				successGroupIdsInc++

				successGroupIds[successGroupIdsInc] = tmp2
				successGroupIdsInc++
			}
			//zlib.MyPrint("person 5 final groupIds",groupIds)
			continue
		}
		//团队最大值 - 当前人数 = 需要补哪个<组人数>   补数
		needRemainderNum := match.Rule.TeamMaxPeople - personUnit
		if groupPersonNum[needRemainderNum] <= 0 {
			//互补值 不存在 ，或者 互补值 人数 为 0
			continue
		}
		maxNumber := match.getMinPersonNum(groupTotal, groupPersonNum[needRemainderNum])
		match.Log.Debug(util.GetSpaceStr(4) + "needNumber : " + strconv.Itoa(needRemainderNum) + "needNumberPersonTotal" + strconv.Itoa(groupPersonNum[needRemainderNum]) + " maxNumber : " + strconv.Itoa(maxNumber))
		if maxNumber <= 0 {
			continue
		}
		//取出两个组的 redis 数据，这里是真的取数据了
		setA := match.Rule.QueueSign.getGroupPersonIndexList(needRemainderNum, match.rangeStart, match.rangeEnd, 0, maxNumber, true)
		setB := match.Rule.QueueSign.getGroupPersonIndexList(personUnit, match.rangeStart, match.rangeEnd, 0, maxNumber, true)
		//逐条合并 setA setB
		for k, _ := range setA {
			tmp := make(map[int]int)
			tmp[0] = setA[k]
			tmp[1] = setB[k]
			successGroupIds[successGroupIdsInc] = tmp
			successGroupIdsInc++
		}

		wishSuccessGroups += maxNumber

		processedNumber = append(processedNumber, needRemainderNum)
	}
	match.Log.Debug("logarithmic rs : " + match.SuccessGroupIdsToStr(successGroupIds))
}

func (match *Match) setMemberRange(rangeStart string, rangeEnd string) {
	match.rangeStart = rangeStart
	match.rangeEnd = rangeEnd
	//match.Log.Info(" set MemberVar Range ", rangeStart , " ", rangeEnd)
	//match.Log.Info(" set MemberVar Range ", rangeStart , " ", rangeEnd)
}

func (match *Match) clearMemberRange() {
	match.rangeStart = ""
	match.rangeEnd = ""
	//match.Log.Info(" clear MemberVar : rangeStart  rangeEnd" )
	//match.Log.Info(" clear MemberVar : rangeStart  rangeEnd")
}

//互补数中，有一步是，取：两个互补数，人数的，最小的那个
//如：团队满足人数是5人
//1. 当前： 3人组 的数量是10，它需要 2人组的数量至少也要10个，这样才能合成 10个组(每组5人)
//2. 上面是理想状态，非理想状态，3人组 的数量是10， 2人组 却只有 8个，那么 8就是此函数的结果
func (match *Match) getMinPersonNum(personTotal int, needRemainderNumPerson int) int {
	maxNumber := personTotal
	if needRemainderNumPerson < maxNumber {
		maxNumber = needRemainderNumPerson
	}
	divider := personTotal % 2
	if divider > 0 { //证明是奇数个
		maxNumber--
	}
	return maxNumber
}

//走到这里就证明，有匹配成功的玩家了
//当匹配成功最终，成功筛选出匹配的组/玩家后，该函数开始执行后续插入操作
//注：successGroupIds ，这里只是组ID，是从索引里拿出来的，只是单一删除了索引值，并没有删除真正的组信息
func (match *Match) successConditions(successGroupIds map[int]map[int]int) {
	length := len(successGroupIds)
	match.Log.Info("successConditions  ...   len :   " + strconv.Itoa(length))
	//match.Log.Info("successConditions  ...   len :   " + strconv.Itoa(length))
	redisConnFD := match.Redis.GetNewConnFromPool()
	defer redisConnFD.Close()

	if match.Rule.Type == service.RULE_TYPE_TEAM_EACH_OTHER { //满足人数即开团
		match.Log.Debug("case : RuleFlagCollectPerson")
		//zlib.MyPrint("successGroupIds",successGroupIds)
		for _, oneCondition := range successGroupIds {
			match.Redis.Multi(redisConnFD)
			util.MyPrint("oneCondition : ", oneCondition)
			resultElement := match.Rule.QueueSuccess.NewResult()
			match.Log.Info("new ResultElement struct")
			//zlib.MyPrint("new resultElement : ",resultElement)
			teamId := 1
			groupIdsArr := make(map[int]int)
			playerIdsArr := make(map[int]int)
			for _, groupId := range oneCondition {
				match.successConditionAddOneGroup(redisConnFD, resultElement.Id, groupId, teamId, groupIdsArr, playerIdsArr)
			}

			//zlib.MyPrint("groupIdsArr",groupIdsArr,"playerIdsArr",playerIdsArr)
			resultElement.GroupIds = util.MapCovertArr(groupIdsArr)
			resultElement.PlayerIds = util.MapCovertArr(playerIdsArr)
			resultElement.Teams = []int{teamId}
			//zlib.MyPrint("resultElement",resultElement)
			util.MyPrint("QueueSuccess.addOne", resultElement)
			pushElement := match.Rule.QueueSuccess.addOne(redisConnFD, resultElement)
			match.Redis.Exec(redisConnFD)
			match.Rule.RuleManager.Option.GameMatch.PersistenceRecordSuccessResult(resultElement, match.Rule.Id)
			match.Rule.RuleManager.Option.GameMatch.PersistenceRecordSuccessPush(pushElement, match.Rule.Id)
		}
	} else { //组队互相PK
		match.Log.Debug("case : RuleFlagTeamVS")
		if length == 1 {
			match.groupPushBackCondition(redisConnFD, successGroupIds[0])
			match.Log.Warn("successGroupIds length = 1 , break")
			match.Log.Warn("successGroupIds length = 1 , break")
			return
		}
		if length%2 > 0 {

			//组队PK，肯定是至少有2个组，如果出现奇数，证明肯定最后一个不能用了
			//把最后一个数，塞回到redis里，再清空这个数
			index := length - 1
			match.groupPushBackCondition(redisConnFD, successGroupIds[index])
			successGroupIds[index] = nil
			length--

			match.Log.Warn("have single group, index:" + strconv.Itoa(index) + " new length:" + strconv.Itoa(length))
		}
		util.MyPrint("final success cnt : ", successGroupIds, " length : "+strconv.Itoa(length))
		var teamId int
		var resultElement Result

		var groupIdsArr map[int]int
		var playerIdsArr map[int]int
		for i := 0; i < length; i++ {
			match.Log.Info("i:" + strconv.Itoa(i))
			//zlib.MyPrint(successGroupIds[i],i)
			//if len(successGroupIds[i]) == 1 {
			//	match.Log.Info(" has a single")
			//	match.groupPushBackCondition(redisConnFD, successGroupIds[i])
			//	continue
			//}
			//一个成功的结果需要：A队(N个小组) B队(N个小组)
			//第一次是创建结果集，同时，把A队里的小组插入进该结果集中，第二次就不创建结果集了
			if i%2 == 0 {
				match.Redis.Multi(redisConnFD)
				resultElement = match.Rule.QueueSuccess.NewResult()
				teamId = 1
				groupIdsArr = make(map[int]int)
				playerIdsArr = make(map[int]int)

				match.Log.Warn("groupId : " + strconv.Itoa(resultElement.Id))
			}
			//将小组信息依次：插入到结果集中
			for _, groupId := range successGroupIds[i] {
				match.successConditionAddOneGroup(redisConnFD, resultElement.Id, groupId, teamId, groupIdsArr, playerIdsArr)
			}

			resultElement.GroupIds = util.MapCovertArr(groupIdsArr)
			resultElement.PlayerIds = util.MapCovertArr(playerIdsArr)
			//resultElement.Teams = []int{teamId}
			teamId = 2
			if i%2 == 1 {
				teamIds := []int{1, 2}
				resultElement.Teams = teamIds
				util.MyPrint("QueueSuccess.addOne", resultElement)
				pushElement := match.Rule.QueueSuccess.addOne(redisConnFD, resultElement)
				match.Redis.Exec(redisConnFD)
				match.Rule.RuleManager.Option.GameMatch.PersistenceRecordSuccessResult(resultElement, match.Rule.Id)
				match.Rule.RuleManager.Option.GameMatch.PersistenceRecordSuccessPush(pushElement, match.Rule.Id)
			}
		}
		//zlib.ExitPrint(123123)
	}
	match.Log.Info("finish successConditions  ...")
	//util.ExitPrint(11)
	//match.Log.Info("finish successConditions  ...")
}

//取出来的groupIds 可能某些原因 最终并没有用上，但是得给塞回到redis里
//这里其实只是将index数据补充上即可，因为计算的时候，删的也只是索引值
func (match *Match) groupPushBackCondition(redisConn redis.Conn, oneCondition map[int]int) {
	util.MyPrint("groupPushBackCondition:", oneCondition)
	//match.Log.Info("groupPushBackCondition", oneCondition)
	for _, groupId := range oneCondition {
		group := match.Rule.QueueSign.getGroupElementById(groupId)
		match.Rule.QueueSign.addOneGroupIndex(redisConn, groupId, group.Person, group.Weight)
	}
}

//添加一个组
func (match *Match) successConditionAddOneGroup(redisConnFD redis.Conn, resultId int, groupId int, teamId int, groupIdsArr map[int]int, playerIdsArr map[int]int) Group {
	match.Log.Info("successConditionAddOneGroup , resultId:" + strconv.Itoa(resultId) + " ,groupId:" + strconv.Itoa(groupId) + " ,teamId:" + strconv.Itoa(teamId))
	match.Log.Info("successConditionAddOneGroup")
	//先以出之前报名的组信息
	group := match.Rule.QueueSign.getGroupElementById(groupId)
	//match.Log.Debug("getGroupElementById group ",group)
	//groupIdsArr = append(  (*groupIdsArr),groupId)
	groupIdsArr[len(groupIdsArr)] = groupId
	playerIdsArrInc := len(playerIdsArr)
	for _, player := range group.Players {
		playerIdsArr[playerIdsArrInc] = player.Id
		playerIdsArrInc++
	}
	//将之前<报名小组>信息复制，并更新相关值
	SuccessGroup := group
	SuccessGroup.Type = service.GAME_MATCH_GROUP_TYPE_SUCCESS
	SuccessGroup.SuccessTimeout = util.GetNowTimeSecondToInt() + match.Rule.SuccessTimeout
	//SuccessGroup.LinkId = resultId
	SuccessGroup.SuccessTime = util.GetNowTimeSecondToInt()
	SuccessGroup.TeamId = teamId
	//fmt.Printf("%+v",SuccessGroup)
	//zlib.ExitPrint(222)
	//添加一条新的小组
	util.MyPrint("addOneGroup", SuccessGroup)
	match.Rule.QueueSuccess.addOneGroup(redisConnFD, SuccessGroup)
	//开始删除，旧的<报名小组>
	match.Log.Warn("delSingOldGroup " + strconv.Itoa(groupId))
	match.Rule.QueueSign.delOneRuleOneGroup(redisConnFD, groupId, 0)
	match.Rule.RuleManager.Option.GameMatch.PersistenceRecordGroup(group, match.Rule.Id) //持久化数据
	//更新玩家状态值，上面其实已经把原玩家状态给清空了
	for _, player := range group.Players {
		_, isEmpty := match.Rule.PlayerManager.GetById(player.Id)
		//var newPlayerStatusElement Player
		if isEmpty == 1 {
			match.Log.Error("GetById empty " + strconv.Itoa(player.Id))
			continue
			//newPlayerStatusElement = match.Rule.PlayerManager.create()
			//} else {
			//newPlayerStatusElement = playerStatusElement
		}
		match.Rule.PlayerManager.UpStatus(player.Id, service.GAME_MATCH_PLAYER_STATUS_SUCCESS, group.SuccessTimeout, redisConnFD)

		//newPlayerStatusElement.Status = service.GAME_MATCH_PLAYER_STATUS_SUCCESS
		//newPlayerStatusElement.SuccessTimeout = group.SuccessTimeout
		//newPlayerStatusElement.GroupId = group.Id

		//queueSign.Log.Info("playerStatus.upInfo:" ,PlayerStatusSign)
		//match.Rule.PlayerManager.setInfo(newPlayerStatusElement, redisConnFD)
		//match.Log.Info("playerStatus.upInfo ", "oldStatus : ",PlayerStatusElement.Status,"newStatus : ",newPlayerStatusElement.Status)
	}
	//zlib.MyPrint( "add one group : ")
	//fmt.Printf("%+v",SuccessGroup)
	return group
}
