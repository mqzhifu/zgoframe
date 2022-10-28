package gamematch
//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"flag"
//	"io/ioutil"
//	"net/http"
//	"os"
//	"reflect"
//	"strconv"
//	"strings"
//	"zlib"
//)
//
////content_type 类型 枚举
//const (
//	CT_EMPTY 		= ""
//	CT_JSON 		= "application/json"
//	CT_MULTIPART 	= "multipart/form-data"
//	CT_URLENCODED 	= "application/x-www-form-urlencoded"
//	CT_STREAM 		= "application/octet-stream"
//	CT_PLAIN 		= "text/plain"
//	CT_HTML 		= "text/html"
//	CT_JS 			= "application/javascript"
//	CT_XML 			= "application/xml"
//)
////请求-ContentType
//type ContentType struct {
//	Name		string
//	Char 		string
//	Addition	string
//}
////响应结构体
//type ResponseMsgST struct {
//	Code 	int
//	Msg 	interface{}
//}
////httpd 类
//type Httpd struct {
//	Option HttpdOption
//	Gamematch *Gamematch
//	Log *zlib.Log
//	HttpServer *http.Server
//}
////实例化 初始值
//type HttpdOption struct {
//	Host			string
//	Port 			string
//	Log 			*zlib.Log
//	BaseRootPath 	string
//	RootPath 		string
//}
////实例化
//func NewHttpd(httpdOption HttpdOption,gamematch *Gamematch)(httpd  *Httpd,err error){
//
//	httpdOption.Log.Info("NewHttpd : ",httpdOption)
//	httpd = new (Httpd)
//	httpd.Option = httpdOption
//	httpd.Option.RootPath = "/www/"
//
//	baseDir, _ := os.Getwd()
//	httpd.Option.BaseRootPath = baseDir + httpd.Option.RootPath
//
//	httpd.Gamematch = gamematch
//	newLog ,err  := getModuleLogInc("httpd")//初始化日志文件
//	if err != nil{
//		return httpd,err
//	}
//	httpd.Log = newLog
//	return httpd,nil
//}
////开启HTTP监听
//func (httpd *Httpd)Start()error{
//	//mymetrics.CreateOneNode("httpSignReq")//请求数
//	//mymetrics.CreateOneNode("httpSignReqSuccess")//请求成功数
//	//mymetrics.CreateOneNode("httpSignReqFiled")//请求失败数
//
//	//开启HTTP监听进程
//	httpd.StartHttpLisnten()
//	//将所有rule的HTTPD请求状态都打开，允许外面访问
//	for k,_:=range httpd.Gamematch.HttpdRuleState{
//		httpd.Gamematch.HttpdRuleState[k] = HTTPD_RULE_STATE_OK
//	}
//	//获取一个管道，用于接收停止信号
//	myChan := httpd.Gamematch.NewSignalChan(0,"httpd")
//	//阻塞，等待接收信号
//	httpd.Gamematch.signReceive(myChan,"httpd")
//	//一但执行到这里，就证明接收到了信号(关闭)
//	ctx ,_ := context.WithCancel(context.Background())
//	httpd.HttpServer.Shutdown(ctx)
//	httpd.Log.CloseChan <- 1
//	//发送信号给管道，告知，HTTPD守护协程已结束成功
//	httpd.Gamematch.signSend(myChan,SIGNAL_GOROUTINE_EXIT_FINISH,"httpd")
//
//	mylog.Warning("httpd goRoutune end")
//	httpd.Log.Warning("httpd goRoutune end")
//	return nil
//}
////开启http 守护监听 协程
//func (httpd *Httpd)StartHttpLisnten(){
//	dns := httpd.Option.Host + ":" + httpd.Option.Port
//	var addr = flag.String("server addr", dns, "server address")
//	httpServer := & http.Server{
//		Addr:*addr,
//	}
//	//监听目录：/，设置统一回调函数
//	http.HandleFunc("/", httpd.RouterHandler)
//
//	httpd.HttpServer = httpServer
//	httpd.Option.Log.Info("httpd start loop:",dns , " Listen : /")
//	go httpd.HttpListen()
//	//httpd.Gamematch.Option.Goroutine.CreateExec(httpd,"HttpListen")
//
//}
//func  (httpd *Httpd)HttpListen( ){
//	err := httpd.HttpServer.ListenAndServe()
//	if err != nil{
//		mylog.Error("http.ListenAndServe() err :  " + err.Error())
//		httpd.Log.Error("http.ListenAndServe() err :  " + err.Error())
//	}
//}
////启动成功后，接收所有HTTP 并路由
//func (httpd *Httpd)RouterHandler(w http.ResponseWriter, r *http.Request){
//
//	parameter := r.URL.Query()//GET 方式URL 中的参数 转 结构体
//	uri := r.URL.RequestURI()
//
//	contentType := httpd.GetContentType(r)
//	httpd.Option.Log.Info("receiver :  uri :",uri," , url.query : ",parameter, " method : ",r.Method , " content_type : ",contentType)
//	//httpd.Log.Info("receiver :  uri :",uri," , url.query : ",parameter, " method : ",r.Method , " content_type : ",contentType)
//	//httpd.Log.Debug(r.Header)
//	var postJsonStr string
//	if strings.ToUpper(r.Method)  == "POST"{
//		GetPostDataMap ,myJsonStr,errs := httpd.GetPostData(r,contentType.Name)
//		if errs != nil{
//			httpd.ResponseStatusCode(w,500,"httpd.GetPostDat" + errs.Error() )
//			return
//		}
//		httpd.Option.Log.Info("PostData",GetPostDataMap,errs)
//		httpd.Log.Info("PostData",GetPostDataMap,errs)
//		postJsonStr = myJsonStr
//	}else{
//		//this is GET method
//	}
//	if r.URL.RequestURI() == "/favicon.ico" {//浏览器的ICON
//		httpd.ResponseStatusCode(w,403,"no power")
//		return
//	}
//	if uri == "" || uri == "/" {
//		httpd.ResponseStatusCode(w,500,"RequestURI is null or uri is :  '/'")
//		return
//	}
//	//去掉 URI 中最后一个 /
//	uriLen := len(uri)
//	if string([]byte(uri)[uriLen-1:uriLen]) == "/"{
//		uri = string([]byte(uri)[0:uriLen - 1])
//	}
//	httpd.Log.Info("final uri : ",uri , " start routing ...")
//	//*********: 还没有加  v1  v2 版本号
//	code := 200
//	var msg interface{}
//	if uri == "/sign" {//报名
//		mymetrics.FastLog("HttpSign",zlib.METRICS_OPT_INC,0)
//		code,msg = httpd.signHandler(postJsonStr)
//		if code == 200{
//			mymetrics.FastLog("HttpSignSuccess",zlib.METRICS_OPT_INC,0)
//		}else{
//			mymetrics.FastLog("httpSignFiled",zlib.METRICS_OPT_INC,0)
//		}
//	}else if uri == "/sign/cancel"{//取消报名
//		mymetrics.FastLog("HttpCancel",zlib.METRICS_OPT_INC,0)
//		code,msg = httpd.signCancelHandler(postJsonStr)
//		if code == 200{
//			mymetrics.FastLog("HttpCancelSuccess",zlib.METRICS_OPT_INC,0)
//		}else{
//			mymetrics.FastLog("HttpCancelFiled",zlib.METRICS_OPT_INC,0)
//		}
//	}else if uri == "/success/del"{//匹配成功记录，不想要了，删除一掉
//		code,msg = httpd.successDelHandler(postJsonStr)
//	}else if uri == "/config"{//
//		code,msg = httpd.ConfigHandler(postJsonStr)
//	}else if uri == "/rule/add" {//添加一条rule
//		//code,msg = httpd.ruleAddOne(postDataMap)
//	}else if uri == "/tools/getErrorInfo" {//所有错误码列表
//		code,msg = httpd.getErrorInfoHandler()
//	}else if uri == "/tools/clearRuleByCode"{//清空一条rule的所有数组，用于测试
//		code,msg = httpd.clearRuleByCodeHandler(postJsonStr)
//	}else if uri == "/tools/getNormalMetrics"{//html api
//		code,msg = httpd.normalMetrics()
//	}else if uri == "/tools/getRedisMetrics"{//html api
//		code,msg = httpd.redisMetrics()
//	}else if uri == "/tools/RedisStoreDb"{//html api
//		code,msg = httpd.RedisStoreDb()
//	}else if uri == "/tools/getHttpReqBusiness"{//html api
//		httpReqBusinessStruct := HttpReqBusiness{}
//		//httpReqBusinessJson,_ := json.Marshal(httpReqBusinessStruct)
//		httpReqBusinessStructDesc := make(map[string]string)
//		types := reflect.TypeOf(&httpReqBusinessStruct)
//		for i:=0 ; i < types.Elem().NumField() ; i++{
//			field := types.Elem().Field(i)
//			tagName1 := field.Tag.Get("desc")
//			mapFieldName := zlib.StrFirstToLower(field.Name)
//			httpReqBusinessStructDesc[mapFieldName] = tagName1
//		}
//		//httpReqBusinessStructDescJson,_ := json.Marshal(httpReqBusinessStructDesc)
//
//		type MyResult struct {
//			MyStrcut  HttpReqBusiness
//			MyStrcutDesc map[string]string
//		}
//
//		rs := MyResult{}
//		rs.MyStrcut = httpReqBusinessStruct
//		rs.MyStrcutDesc = httpReqBusinessStructDesc
//		finalRs ,_ := json.Marshal(&rs)
//		code = 200
//		msg = string(finalRs)
//	}else{
//		//先匹配下静态资源
//		hasMatch, err := httpd.routeStatic(w,r,uri)
//		if err != nil{
//			code = 500
//			msg = " uri router failed."
//		}else{
//			if hasMatch{
//				//routeStatic 内部已处理
//				return
//			}else{
//				code = 500
//				msg = " uri router failed."
//			}
//		}
//	}
//
//	httpd.ResponseMsg(w,code,msg)
//
//}
//func  getStaticFileContent(fileSuffix string)(content string ,err error){
//	//code,msg = httpd.redisMetrics()
//	baseDir ,_ := os.Getwd()
//	//path := baseDir + "/../gamematch/www"
//	path := baseDir
//	mylog.Debug("final path:",path)
//	filePath := path+fileSuffix
//	mylog.Info("getStaticFileContent File path:",filePath)
//	b, err := ioutil.ReadFile(filePath) // just pass the file name
//	return string(b),err
//}
//func  (httpd *Httpd) routeStatic(w http.ResponseWriter,r *http.Request,uri string)(hasMatch bool,err error){
//	//zlib.ExitPrint(httpd.Option.BaseRootPath)
//	fileList := zlib.GetFileListByDir(httpd.Option.BaseRootPath)
//	if len(fileList) <=0 {
//		httpd.Log.Error("GetFileListByDir is empty")
//		return false,nil
//	}
//
//	for _,fileName := range fileList{
//		tmpFileName := httpd.Option.RootPath + fileName
//		if uri == tmpFileName{
//			fileContent, err := getStaticFileContent(uri)
//			if err != nil {
//				return false,errors.New("getStaticFileContent err:" + err.Error())
//			}
//			//踦域处理
//			w.Header().Set("Access-Control-Allow-Origin","*")
//			w.Header().Add("Access-Control-Allow-Headers","Content-Type")
//			//w.Header().Set("content-type","text/plain")
//
//			w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))
//			w.Write([]byte(fileContent))
//			return true,nil
//		}
//	}
//	httpd.Log.Error("routeStatic no match :",uri)
//	return false,nil
//}
//
////func  (httpd *Httpd)routeStatic(w http.ResponseWriter,r *http.Request,uri string){
//	//uriSplit := strings.Split(uri,"?")
//	//if uriSplit[0] == "/apireq.html" {
//	//	uri = uriSplit[0]
//	//}
//	//if uri == "/jquery.min.js" ||
//	//	uri == "/metrics.html" ||
//	//	uri == "/errorlist.html" ||
//	//	uri == "/apireq.html" ||
//	//	uri == "/index.html" ||
//	//	uri == "/apilist.html" ||
//	//	uri == "/apireq.html" ||
//	//	uri == "/favicon.ico" ||
//	//	uri == "/flow.jpg" { //静态文件
//	//	fileContent, err := httpd.getStaticFileContent(uri)
//	//	if err != nil {
//	//		httpd.ResponseStatusCode(w, 404, err.Error())
//	//		return
//	//	}
//	//	//踦域处理
//	//	w.Header().Set("Access-Control-Allow-Origin","*")
//	//	w.Header().Add("Access-Control-Allow-Headers","Content-Type")
//	//	//w.Header().Set("content-type","text/plain")
//	//
//	//	w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))
//	//	w.Write([]byte(fileContent))
//	//}
////}
//func  (httpd *Httpd)getStaticFileContent(fileSuffix string)(content string ,err error){
//	//code,msg = httpd.redisMetrics()
//	baseDir ,_ := os.Getwd()
//	path := baseDir + "/../gamematch/www"
//	mylog.Debug("final path:",path)
//	filePath := path+fileSuffix
//	httpd.Option.Log.Info("getStaticFileContent File path:",filePath)
//	b, err := ioutil.ReadFile(filePath) // just pass the file name
//	return string(b),err
//}
////响应的具体内容
//func  (httpd *Httpd)ResponseMsg(w http.ResponseWriter,code int ,msg interface{} ){
//	responseMsgST := ResponseMsgST{Code:code,Msg:msg}
//	//msg = msg[1:len(msg)-1]
//	//这里有个无奈的地方，为了兼容非网络请求，正常使用时，返回的就是json,现在HTTP套一层，还得再一层JSON，冲突了
//	jsonResponseMsg , err := json.Marshal(responseMsgST)
//	//jsonResponseMsgNew := strings.Replace(string(jsonResponseMsg),"#msg#",msg,-1)
//	if code == 200{
//		httpd.Option.Log.Info("ResponseMsg rs",err, string(jsonResponseMsg))
//		httpd.Log.Info("ResponseMsg",string(jsonResponseMsg))
//	}else{
//		httpd.Option.Log.Error("ResponseMsg rs",err, string(jsonResponseMsg))
//		httpd.Log.Error("ResponseMsg",string(jsonResponseMsg))
//	}
//
//	w.Header().Set("Content-Length",strconv.Itoa( len(jsonResponseMsg) ) )
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	w.Write([]byte(jsonResponseMsg))
//}
////http 响应状态码
//func (httpd *Httpd)ResponseStatusCode(w http.ResponseWriter,code int ,responseInfo string){
//	httpd.Option.Log.Info("ResponseStatusCode",code,responseInfo)
//	httpd.Log.Info("ResponseStatusCode",code,responseInfo)
//
//	w.Header().Set("Content-Length",strconv.Itoa( len(responseInfo) ) )
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(403)
//	w.Write([]byte(responseInfo))
//}
//
//
//func (httpd *Httpd) GetContentType( r *http.Request)ContentType{
//	//r.Header.Get("Content-Type")
//	contentTypeArr ,ok := r.Header["Content-Type"]
//	//httpd.Option.Log.Debug(contentTypeArr)
//
//	//正常的请求基本上没这个值，除了 FORM，因为只有传输内容的时候才有意义
//	contentType := ContentType{}
//	if ok {
//		contentType.Name = contentTypeArr[0]
//		//httpd.Option.Log.Debug(contentType.Name)
//		if strings.Index(contentType.Name,"multipart/form-data") != -1{
//			tmpArr := strings.Split(contentType.Name,";")
//			contentType.Addition = strings.TrimSpace(tmpArr[1])
//			contentType.Name = CT_MULTIPART
//		}else{
//			tmpArr := strings.Split(contentType.Name,";")
//			if len(tmpArr) >= 2{
//				contentType.Char = strings.TrimSpace(tmpArr[1])
//			}
//		}
//		elementIndex := zlib.ElementStrInArrIndex(GetContentTypeList(),contentType.Name)
//		if elementIndex == -1{
//			httpd.Option.Log.Notice("content type is unknow ")
//		}
//	}else{
//		contentType.Name =CT_EMPTY
//	}
//	return contentType
//}
//
//
//func GetContentTypeList()[]string{
//	list := []string{
//		CT_JSON,CT_MULTIPART,CT_URLENCODED,CT_STREAM,CT_EMPTY,CT_PLAIN,CT_JS,CT_HTML,CT_XML,
//	}
//	return list
//}
//
//func (httpd *Httpd)GetPostData(r *http.Request,contentType string)( data  map[string]interface{},jsonStr string, err error){
//	//httpd.Option.Log.Debug(" getPostData ")
//	if r.ContentLength == 0{//获取主体数据的长度
//		return data,jsonStr,nil
//	}
//	switch contentType {
//		case CT_JSON:
//			body := make([]byte, r.ContentLength)
//			r.Body.Read(body)
//			//mylog.Debug("body : ",string(body))
//
//			jsonDataMap := make(map[string]interface{})
//			errs := json.Unmarshal(body,&jsonDataMap)
//			if errs != nil{
//				httpd.Log.Error("json.Unmarshal failed , ",body)
//			}
//			return jsonDataMap,string(body),nil
//		case CT_MULTIPART:
//			data = make( map[string]interface{})
//			r.ParseMultipartForm(r.ContentLength)
//			for k,v:= range r.Form{
//				data[k] = v
//			}
//			return data,jsonStr,nil
//		case CT_URLENCODED:
//			err := r.ParseForm()
//			if err != nil{
//				return data,jsonStr,err
//			}
//			data = make( map[string]interface{})
//			for k,v:= range r.Form{
//				if len(v) == 1{
//					data[k] = v[0]
//				}else{
//					zlib.ExitPrint(" bug !!!")
//				}
//			}
//			httpReqBusinessStruct := HttpReqBusiness{}
//			types := reflect.TypeOf(&httpReqBusinessStruct)
//			for i:=0 ; i < types.Elem().NumField() ; i++{
//				thisFiledType := types.Elem().Field(i).Type.String()
//				fileName :=  zlib.StrFirstToLower(types.Elem().Field(i).Name)
//				if  thisFiledType == "int"{
//					switch  data[fileName].(type) {
//					case nil:
//						data[fileName] = 0
//					default:
//						//if data[types.Elem().Field(i).Name].(string) == ""  {
//						//	data[types.Elem().Field(i).Name] = 0
//						//}else{
//							data[fileName] = zlib.Atoi( data[fileName].(string))
//						//}
//					}
//				}
//			}
//			//这里先做个特殊处理，回头想下如何动态化
//			//if r.URL.RequestURI() == "/sign"{
//				playerListStr ,ok := data["playerList"]
//				if ok {
//					//[{"uid":2,"matchAttr":{"age":1,"sex":2}}]
//					//[{\"uid\":2,\"matchAttr\":{\"age\":1,\"sex\":2}}]
//					playerListArr := []HttpReqPlayer{}
//					json.Unmarshal([]byte(playerListStr.(string)),&playerListArr)
//					data["playerList"] = playerListArr
//				}else{
//					data["playerList"] = nil
//				}
//			//}
//			jsonStr ,_ := json.Marshal(data)
//
//			return data,string(jsonStr),nil
//		default:
//			httpd.Option.Log.Error("contentType no support : ",contentType , " ,no data")
//	}
//
//	return data,jsonStr,nil
//}