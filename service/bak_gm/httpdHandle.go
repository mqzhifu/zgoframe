package gamematch
//
//import (
//	"encoding/json"
//	"runtime"
//	"strconv"
//	"strings"
//	"zlib"
//)
//////报名 - 玩家
////func  (httpd *Httpd)signHandler( postJsonStr string)(code int ,msg interface{}){
////	httpd.Log.Info(" routing in signHandler : ")
////
////	errCode , httpReqBusiness:= httpd.businessCheckData(postJsonStr)
////	if errCode != 0{
////		errs := myerr.NewErrorCode(errCode)
////		errInfo := zlib.ErrInfo{}
////		json.Unmarshal([]byte(errs.Error()),&errInfo)
////		return errInfo.Code,errInfo.Msg
////	}
////	errs := httpd.Gamematch.CheckHttpSignData(httpReqBusiness)
////	if errs != nil{
////		errInfo := zlib.ErrInfo{}
////		json.Unmarshal([]byte(errs.Error()),&errInfo)
////
////		return errInfo.Code,errInfo.Msg
////	}
////	signRsData, errs := httpd.Gamematch.Sign(httpReqBusiness)
////	if errs != nil{
////		errInfo := zlib.ErrInfo{}
////		json.Unmarshal([]byte(errs.Error()),&errInfo)
////
////		return errInfo.Code,errInfo.Msg
////	}
////	return 200,signRsData
////}
//
//type AppAllConfig struct {
//	CmsArg				map[string]string
//	PushRetryPeriod		[]int
//	EtcdAppConf			map[string]string
//	AppInfo 			zlib.App
//	Other 				map[string]string
//}
//
//func  (httpd *Httpd)ConfigHandler(postJsonStr string)(code int ,msg interface{}){
//	//cmdArgs := CmdArgs{}
//	//typeOfCmsArgs := reflect.TypeOf(cmdArgs)
//	//cmsArgAfter := make(map[string]string)
//	//for i:=0;i<typeOfCmsArgs.NumField();i++{
//	//	memVar := typeOfCmsArgs.Field(i)
//	//	desc := memVar.Tag.Get("desc")
//	//	cmsArgAfter[memVar.Name] = desc
//	//}
//
//	appM  := zlib.NewAppManager()
//	app,empty := appM.GetById(APP_ID)
//	if !empty{
//		mylog.Error("ConfigHandler appM.GetById empty~")
//	}
//
//	other := make(map[string]string)
//	getLinkAddressList := myetcd.GetLinkAddressList()
//	other["EtcdHost"] = getLinkAddressList[0]
//	other["ServicePrefix"] = SERVICE_PREFIX
//	other["RuleEtcdConfigPrefix"] = RuleEtcdConfigPrefix
//	appAllConfig := AppAllConfig{
//		CmsArg : myGamematchOption.Option.CmsArg,
//		PushRetryPeriod : PushRetryPeriod,
//		EtcdAppConf : myetcd.GetAppConf(),
//		AppInfo: app,
//		Other: other,
//	}
//
//	return 200,appAllConfig
//}
//
////
//func  (httpd *Httpd)successDelHandler(postJsonStr string)(code int ,msg interface{}){
//	httpd.Log.Info(" routing in successDelHandler : ")
//
//	errCode , httpReqBusiness:= httpd.businessCheckData(postJsonStr)
//	if errCode != 0{
//		errs := myerr.NewErrorCode(errCode)
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//		return errInfo.Code,errInfo.Msg
//	}
//	errs := httpd.Gamematch.CheckHttpSuccessDelData(httpReqBusiness)
//	if errs != nil{
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//		return errInfo.Code,  errInfo.Msg
//	}
//	successClass := httpd.Gamematch.getContainerSuccessByRuleId(httpReqBusiness.RuleId)
//	_ ,isEmpty := successClass.GetResultById(httpReqBusiness.SuccessId,0,0)
//	if isEmpty == 1{
//		errs := myerr.NewErrorCode(807)
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//		return errInfo.Code,errInfo.Msg
//	}
//
//	redisConn := myredis.GetNewConnFromPool()
//	myredis.Multi(redisConn )
//	successClass.delOneResult(redisConn,httpReqBusiness.SuccessId,1,1,1,1)
//	myredis.Exec(redisConn )
//	return 200,"ok"
//}
////取消报名 - 删除已参与匹配的玩家信息
//func  (httpd *Httpd)signCancelHandler(postJsonStr string)(code int ,msg interface{}){
//	errCode , httpReqBusiness := httpd.businessCheckData(postJsonStr)
//	if errCode != 0{
//		errs := myerr.NewErrorCode(errCode)
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//		return errInfo.Code,errInfo.Msg
//	}
//	errs := httpd.Gamematch.CheckHttpSignCancelData(httpReqBusiness)
//	if errs != nil{
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//		return errInfo.Code,  errInfo.Msg
//	}
//
//	signClass := httpd.Gamematch.GetContainerSignByRuleId(httpReqBusiness.RuleId)
//	httpd.Log.Info("del by groupId")
//	err := signClass.cancelByGroupId(httpReqBusiness.GroupId)
//	if err != nil{
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(err.Error()),&errInfo)
//
//		return errInfo.Code,errInfo.Msg
//	}
//
//	return 200,"ok"
//}
////获取错误码
//func  (httpd *Httpd)getErrorInfoHandler()(code int ,msg interface{}){
//	httpd.Log.Info(" routing in getErrorInfoHandler : ")
//
//	container := getErrorCode()
//	var res []MyErrorCode
//	for _,v:= range  container{
//		row := strings.Split(v,",")
//		myErrorCode := MyErrorCode{
//			Code: zlib.Atoi(row[0]),
//			Msg :row[1],
//			Flag: row[2],
//			MsgCn: row[3],
//		}
//		res = append(res,myErrorCode)
//	}
//
//	//msg,_ := json.Marshal(res)
//	return 200,res
//}
////清除一条rule的所有数据，用于测试
//func  (httpd *Httpd)clearRuleByCodeHandler(postJsonStr string)(code int ,msg interface{}){
//	httpd.Log.Info(" routing in clearRuleByCodeHandler : ")
//
//	errCode , httpReqBusiness:= httpd.businessCheckData(postJsonStr)
//	if errCode != 0{
//		errs := myerr.NewErrorCode(errCode)
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//		return errInfo.Code,errInfo.Msg
//	}
//
//	checkCodeRs := false
//	ruleId := 0
//	for _,rule := range httpd.Gamematch.RuleConfig.getAll(){
//		if  rule.CategoryKey == httpReqBusiness.MatchCode{
//			ruleId = rule.Id
//			checkCodeRs = true
//			break
//		}
//	}
//	if !checkCodeRs{
//		errs := myerr.NewErrorCode(451)
//		errInfo := zlib.ErrInfo{}
//		json.Unmarshal([]byte(errs.Error()),&errInfo)
//
//		return errInfo.Code,errInfo.Msg
//	}
//
//	httpd.Gamematch.RuleConfig.delOne(ruleId)
//	return 200,"ok"
//}
//func  (httpd *Httpd)normalMetrics()(code int ,msg interface{}){
//	data := mymetrics.GetAll()
//	data["sysGoroutineNum"] = runtime.NumGoroutine()
//	return 200,data
//	//return code,msg
//}
//func  (httpd *Httpd)redisMetrics()(code int ,msg interface{}){
//	//rulelist map[int]Rule ,list map[int]map[string]int,playerCnt
//	rulelist,list,playerCnt,rulePersonNum := httpd.Gamematch.RedisMetrics()
//	data := make(map[string]interface{})
//	data["ruleList"] = rulelist
//	data["ruleTotal"] = list
//	data["playerStatus"] = playerCnt
//	data["ruleSignPerson"] = rulePersonNum
//
//	PushRetryPeriodStr := ""
//	for _,v := range PushRetryPeriod{
//		PushRetryPeriodStr += " "+strconv.Itoa(v)+ ""
//	}
//	data["pushRetryPeriodStr"] = PushRetryPeriodStr
//
//
//	return 200,data
//}
//
//func  (httpd *Httpd)RedisStoreDb()(code int ,msg interface{}){
//	//player := Player{}
//	rule := Rule{}
//	playerWeight := PlayerWeight{}
//	playerStatusElement := PlayerStatusElement{}
//	group := Group{}
//	result := Result{}
//	pushElement := PushElement{}
//
//	data := make(map[string]string)
//
//	ruleMap := zlib.StructCovertMap(rule)
//	playerWeightMap := zlib.StructCovertMap(playerWeight)
//	playerStatusElementMap := zlib.StructCovertMap(playerStatusElement)
//	groupMap := zlib.StructCovertMap(group)
//	resultMap := zlib.StructCovertMap(result)
//	pushElementMap := zlib.StructCovertMap(pushElement)
//
//
//
//	ruleMapStr,err1 := json.Marshal(ruleMap)
//	playerWeightMapStr,err2 := json.Marshal(playerWeightMap)
//	playerStatusElementMapStr,err3 := json.Marshal(playerStatusElementMap)
//	groupMapStr,err4 := json.Marshal(groupMap)
//	resultMapStr,err5 := json.Marshal(resultMap)
//	pushElementMapStr,err6 := json.Marshal(pushElementMap)
//
//	mylog.Debug("RedisStoreDb json.Marshal:",err1,err2,err3,err4,err5,err6)
//
//	data["rule"] = string(ruleMapStr)
//	data["playerWeight"] = string(playerWeightMapStr)
//	data["playerStatusElement"] = string(playerStatusElementMapStr)
//	data["group"] = string(groupMapStr)
//	data["result"] = string(resultMapStr)
//	data["pushElement"] = string(pushElementMapStr)
//
//	return 200,data
//}
//
////================以上是直接HTTP API 请求的接口，下面是内部服务方法==================================
////通用 业务型  请求 数据  检查
//func  (httpd *Httpd) businessCheckData(postJsonStr string )(errCode int,httpReqBusiness HttpReqBusiness){
//	httpd.Log.Info(" businessCheckData : ")
//	if postJsonStr == ""{
//		return 802,httpReqBusiness
//	}
//	var jsonUnmarshalErr error
//	jsonUnmarshalErr = json.Unmarshal([]byte(postJsonStr),&httpReqBusiness)
//	if jsonUnmarshalErr != nil{
//		httpd.Log.Error(jsonUnmarshalErr)
//		mylog.Error(jsonUnmarshalErr)
//		return 459,httpReqBusiness
//	}
//	if httpReqBusiness.MatchCode == ""{
//		return 450,httpReqBusiness
//	}
//	rule ,err := httpd.Gamematch.RuleConfig.getByCategory(httpReqBusiness.MatchCode)
//	if err !=nil{
//		return 806,httpReqBusiness
//	}
//	httpReqBusiness.RuleId = rule.Id
//	_,err  = httpd.checkHttpdState(httpReqBusiness.RuleId)
//	if err != nil{
//		return 804,httpReqBusiness
//	}
//
//	return 0 ,httpReqBusiness
//}
//
//func  (httpd *Httpd)checkHttpdState(ruleId int)(bool,error){
//	state ,ok := httpd.Gamematch.HttpdRuleState[ruleId]
//	if !ok {
//		return false,myerr.NewErrorCode(803)
//	}
//	if state == HTTPD_RULE_STATE_OK{
//		return true,nil
//	}
//	return false,myerr.NewErrorCode(804)
//}
//
//
////----------------------------------------------------------------------
//
////func  (httpd *Httpd)ruleAddOne(  postDataMap map[string]interface{})(code int ,msg interface{}){
////	data,errs := httpd.Gamematch.CheckHttpAddRule(jsonDataMap)
////	//zlib.MyPrint(errs)
////	if errs != nil{
////		errInfo := zlib.ErrInfo{}
////		json.Unmarshal([]byte(errs.Error()),&errInfo)
////
////		return errInfo.Code,errInfo.Msg
////	}
////
////	httpd.Gamematch.RuleConfig.AddOne()
////
////	return code,msg
////}
//
////func  (httpd *Httpd)getConstList(){
////	httpd.Log.Info(" routing in signHandler : ")
////
////	fileDir := "/data/www/golang/src/gamematch/const.go"
////	doc := zlib.NewDocRegular(fileDir)
////	doc.ParseConst()
////}