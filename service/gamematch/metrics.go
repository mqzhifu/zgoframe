package gamematch

import (
	"zgoframe/service"
	"zgoframe/util"
)

func (gamematch *Gamematch) RedisMetrics() (rulelist map[int]Rule, list map[int]map[string]int, playerCnt map[string]int, rulePersonNum map[int]map[int]int) {
	rulelist = gamematch.RuleConfig.getAll()

	playerList, _ := playerStatus.getAllPlayers()
	playerCnt = make(map[string]int)
	playerCnt["total"] = 0
	playerCnt["signTimeout"] = 0
	playerCnt["successTimeout"] = 0
	//下面是状态
	playerCnt["sign"] = 0
	playerCnt["success"] = 0
	playerCnt["int"] = 0
	playerCnt["unknow"] = 0

	if len(playerList) > 0 {
		now := util.GetNowTimeSecondToInt()
		for _, playerStatusElement := range playerList {
			if playerStatusElement.Status == service.PlayerStatusSign {
				playerCnt["sign"]++
			} else if playerStatusElement.Status == service.PlayerStatusSuccess {
				playerCnt["success"]++
			} else if playerStatusElement.Status == service.PlayerStatusInit {
				playerCnt["int"]++
			} else {
				playerCnt["unknow"]++
			}

			if now-playerStatusElement.SignTimeout > rulelist[playerStatusElement.RuleId].MatchTimeout {
				playerCnt["signTimeout"]++
			}

			if playerStatusElement.SuccessTimeout > 0 {
				if now-playerStatusElement.SuccessTimeout > rulelist[playerStatusElement.RuleId].SuccessTimeout {
					playerCnt["successTimeout"]++
				}
			}

			playerCnt["total"]++
		}
	}
	list = make(map[int]map[string]int)

	rulePersonNum = make(map[int]map[int]int)
	for ruleId, _ := range rulelist {
		//prefix := "rule("+strconv.Itoa(ruleId) + ")"

		row := make(map[string]int)
		row["player"] = playerStatus.getOneRuleAllPlayerCnt(ruleId)
		//playerCnt := playerStatus.getOneRuleAllPlayerCnt(ruleId)
		push := gamematch.getContainerPushByRuleId(ruleId)
		row["push"] = push.getAllCnt()
		row["pushWaitStatut"] = push.getStatusCnt(service.PushStatusWait)
		row["pushRetryStatut"] = push.getStatusCnt(service.PushStatusRetry)

		sign := gamematch.GetContainerSignByRuleId(ruleId)
		row["signGroup"] = sign.getAllGroupsWeightCnt()
		//row["groupPersonCnt"] =
		//map[int]int
		playersByPerson := sign.getPlayersCntByWeight("0", "100")
		rulePersonNum[ruleId] = playersByPerson

		success := gamematch.getContainerSuccessByRuleId(ruleId)
		row["success"] = success.GetAllTimeoutCnt()
		//
		//	//match := httpd.Gamematch.getContainerMatchByRuleId(ruleId)
		list[ruleId] = row
	}
	return rulelist, list, playerCnt, rulePersonNum
}
