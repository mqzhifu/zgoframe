package gamematch

import (
	"zgoframe/util"
)

//func TestAddRuleData(gameMatch *GameMatch) {
//	rule1 := Rule{
//		Id:              1,
//		AppId:           3,
//		CategoryKey:     "RuleFlagCollectPerson",
//		MatchTimeout:    7,
//		SuccessTimeout:  60,
//		IsSupportGroup:  1,
//		Flag:            2,
//		PersonCondition: 4,
//		GroupPersonMax:  5,
//		//TeamPerson: 5,
//		PlayerWeight: PlayerWeight{
//			ScoreMin:   -1,
//			ScoreMax:   -1,
//			AutoAssign: true,
//			Formula:    "",
//			//Flag:1,
//		},
//	}
//
//	rule2 := Rule{
//		Id:              2,
//		AppId:           4,
//		CategoryKey:     "RuleFlagTeamVS",
//		MatchTimeout:    10,
//		SuccessTimeout:  70,
//		IsSupportGroup:  1,
//		Flag:            1,
//		PersonCondition: 5,
//		TeamVSPerson:    5,
//		GroupPersonMax:  5,
//		//TeamPerson: 5,
//		PlayerWeight: PlayerWeight{
//			ScoreMin:   -1,
//			ScoreMax:   -1,
//			AutoAssign: true,
//			Formula:    "",
//			//Flag:1,
//		},
//	}
//
func getOneRandomPlayerUid() int {
	return util.GetRandIntNumRange(1000, 9999)
}

func getOneRandomGroupId() int {
	return util.GetRandIntNumRange(10, 99)
}

//func TestSign() {
//	//{"groupId":100001,"customProp":"","addition":"","matchCode":"test_vs","playerList":[{"uid":2,"matchAttr":{"age":1,"sex":2}}]}
//	localIp, err := util.GetLocalIp()
//	if err != nil {
//		util.ExitPrint("GetLocalIp err : ", err)
//	}
//	myHost := localIp
//
//	//每次报名，生成几个人player
//	signGroupPersonArr := []int{1, 5, 4, 5, 3, 3, 3, 2, 2, 1, 1, 1, 1, 5, 4, 4}
//	var signSuccessGroup []HttpReqGameMatchPlayerSign
//
//	//matchCode := "test_vs"
//	for _, playerNumMax := range signGroupPersonArr {
//		addition := "map_" + strconv.Itoa(util.GetRandIntNumRange(1, 10))
//		var playerStructArr []HttpReqGameMatchPlayer
//		for i := 0; i < playerNumMax; i++ {
//			matchAttr := make(map[string]int)
//			matchAttr["age"] = util.GetRandIntNumRange(1, 2)
//			matchAttr["sex"] = util.GetRandIntNumRange(1, 10)
//			playerUid := getOneRandomPlayerUid()
//			player := HttpReqGameMatchPlayer{Uid: playerUid, WeightAttr: matchAttr}
//			playerStructArr = append(playerStructArr, player)
//		}
//		httpReqSign := HttpReqGameMatchPlayerSign{
//			GroupId:    getOneRandomGroupId(),
//			PlayerList: playerStructArr,
//			Addition:   addition,
//		}
//		reqData, _ := json.Marshal(httpReqSign)
//		url := "http://" + myHost + ":5678/sign"
//
//		client := &http.Client{Timeout: 5 * time.Second}
//		resp, errs := client.Post(url, "application/json", bytes.NewBuffer(reqData))
//		//etcdOption.Log.Info(" get etcd config ip:port list : ",etcdOption.FindEtcdUrl,errs)
//		if errs != nil {
//			util.ExitPrint("http.Post err : ", errs)
//		}
//		htmlContentJson, _ := ioutil.ReadAll(resp.Body)
//		util.MyPrint("post rs : ", string(htmlContentJson))
//		htmlContent := httpresponse.Response{}
//		json.Unmarshal(htmlContentJson, &htmlContent)
//		if htmlContent.Code == 200 {
//			signSuccessGroup = append(signSuccessGroup, httpReqSign)
//		}
//		//zlib.ExitPrint("post rs : ",string(htmlContentJson))
//		//myGamematch.Sign(ruleId,9999,customProp,playerStructArr , "im_addition")
//	}
//	util.MyPrint("signSuccessGroup len : ", len(signSuccessGroup), " hope len: ", len(signGroupPersonArr), "filed :", len(signGroupPersonArr)-len(signSuccessGroup))
//	for k, v := range signSuccessGroup {
//		util.MyPrint(k, v)
//	}
//	util.ExitPrint("im end")
//}
