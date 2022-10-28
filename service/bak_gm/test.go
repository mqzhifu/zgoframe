package gamematch

//func clear(myGamematch Gamematch){
//	myGamematch.DelAll()
//	zlib.ExitPrint("clear all done. ")
//}
//func TestAddRuleData(myGamematch Gamematch){
//	rule1 := Rule{
//		Id: 1,
//		AppId: 3,
//		CategoryKey: "RuleFlagCollectPerson",
//		MatchTimeout: 7,
//		SuccessTimeout: 60,
//		IsSupportGroup: 1,
//		Flag: 2,
//		PersonCondition: 4,
//		GroupPersonMax:5,
//		//TeamPerson: 5,
//		PlayerWeight: PlayerWeight{
//			ScoreMin:-1,
//			ScoreMax:-1,
//			AutoAssign:true,
//			Formula:"",
//			//Flag:1,
//		},
//	}
//
//	rule2 := Rule{
//		Id: 2,
//		AppId: 4,
//		CategoryKey: "RuleFlagTeamVS",
//		MatchTimeout: 10,
//		SuccessTimeout: 70,
//		IsSupportGroup: 1,
//		Flag: 1,
//		PersonCondition: 5,
//		TeamVSPerson:5,
//		GroupPersonMax:5,
//		//TeamPerson: 5,
//		PlayerWeight: PlayerWeight{
//			ScoreMin:-1,
//			ScoreMax:-1,
//			AutoAssign:true,
//			Formula:"",
//			//Flag:1,
//		},
//	}
//
//	myGamematch.RuleConfig.AddOne(rule1)
//	myGamematch.RuleConfig.AddOne(rule2)
//}

//func testContext(){
//	ctx ,cancel := context.WithCancel(context.Background())
//
//	go func(ctx context.Context,url string ){
//		select {
//			case <- ctx.Done():
//				zlib.MyPrint("sub func done.")
//				return
//		}
//	}(ctx,"aaaa")
//	time.Sleep(3 * time.Second)
//	cancel()
//	dddd()
//}
//
//func dddd()  {
//	for  {
//		zlib.MyPrint("dddd sleep 1")
//		time.Sleep(1 * time.Second)
//	}
//}
//
//func delOneRule(myGamematch Gamematch ,ruleId int){
//	queueSign := myGamematch.getContainerSignByRuleId(ruleId)
//	queueSign.delOneRule()
//
//	queueSuccess := myGamematch.getContainerSuccessByRuleId(ruleId)
//	queueSuccess.delOneRule()
//
//	playerStatus.delAllPlayers()
//	mylog.Notice("testSignDel finish.")
//}
//func testReidsPool(){
//tttttt
//rro := zlib.RedisOption{
//	Host: "127.0.0.1",
//	Port: "6379",
//	Log: mylog,
//}
//mmrr ,err := zlib.NewRedisConn(rro)
//inc := 0
//for{
//	inc ++
//	if inc == 5{
//		zlib.MyPrint(" close ")
//		mmrr.Conn.Close()
//	}
//	time.Sleep(1 * time.Second)
//}
//inc := 0
//redisConn := myredis.GetNewConnFromPool()
//for{
//	inc ++
//	if inc == 5{
//		zlib.MyPrint(" ping start:")
//		myredis.MyRedisDo(redisConn,"ping")
//	}
//	if inc == 10{
//		zlib.MyPrint("close start :")
//		rs := redisConn.Close()
//		zlib.MyPrint("close rs ",rs)
//	}
//	time.Sleep(1 * time.Second)
//}
//tttttt
//}

//========================


//func getOneRandomPlayerUid()int{
//	return zlib.GetRandIntNumRange(1000,9999)
//}
//
//func getOneRandomGroupId()int{
//	return zlib.GetRandIntNumRange(10,99)
//}
//
//func TestSign( ){
//	//{"groupId":100001,"customProp":"","addition":"","matchCode":"test_vs","playerList":[{"uid":2,"matchAttr":{"age":1,"sex":2}}]}
//	localIp,err := zlib.GetLocalIp()
//	if err !=nil{
//		zlib.ExitPrint("GetLocalIp err : ",err)
//	}
//	myHost := localIp
//
//	//每次报名，生成几个人player
//	signGroupPersonArr := []int{1,5,4,5,3,3,3,2,2,1,1,1,1,5,4,4}
//	var signSuccessGroup  []HttpReqBusiness
//
//	matchCode := "test_vs"
//	for _,playerNumMax := range signGroupPersonArr{
//		addition := "map_" + strconv.Itoa(zlib.GetRandIntNumRange(1,10))
//		var playerStructArr []HttpReqPlayer
//		for i:=0;i<playerNumMax;i++{
//			matchAttr := make(map[string]int)
//			matchAttr["age"] = zlib.GetRandIntNumRange(1,2)
//			matchAttr["sex"] = zlib.GetRandIntNumRange(1,10)
//			//player := Player{Id:playerIdInc}
//			playerUid := getOneRandomPlayerUid()
//			player := HttpReqPlayer{Uid:playerUid,MatchAttr:matchAttr}
//			playerStructArr = append(playerStructArr,player)
//			//playerIdInc++
//		}
//		httpReqSign := HttpReqBusiness{
//			MatchCode:  matchCode,
//			GroupId:    getOneRandomGroupId(),
//			PlayerList: playerStructArr,
//			Addition:   addition,
//		}
//		//zlib.ExitPrint(httpReqSign)
//		reqData,_ := json.Marshal(httpReqSign)
//		//zlib.ExitPrint(string(reqData))
//		url := "http://"+myHost+":5678/sign"
//
//		client := &http.Client{Timeout: 5 * time.Second}
//		resp, errs := client.Post(url,"application/json",bytes.NewBuffer(reqData))
//		//etcdOption.Log.Info(" get etcd config ip:port list : ",etcdOption.FindEtcdUrl,errs)
//		if errs != nil{
//			zlib.ExitPrint("http.Post err : ",errs)
//		}
//		htmlContentJson,_ := ioutil.ReadAll(resp.Body)
//		zlib.MyPrint("post rs : ",string(htmlContentJson))
//		htmlContent := ResponseMsgST{}
//		json.Unmarshal(htmlContentJson,&htmlContent)
//		if htmlContent.Code == 200{
//			signSuccessGroup = append(signSuccessGroup,httpReqSign)
//		}
//		//zlib.ExitPrint("post rs : ",string(htmlContentJson))
//		//myGamematch.Sign(ruleId,9999,customProp,playerStructArr , "im_addition")
//	}
//	zlib.MyPrint("signSuccessGroup len : ",len(signSuccessGroup) , " hope len: ",len(signGroupPersonArr) , "filed :",len(signGroupPersonArr)  - len(signSuccessGroup))
//	for k,v := range signSuccessGroup{
//		zlib.MyPrint(k,v)
//	}
//	zlib.ExitPrint("im end")
//}
//
