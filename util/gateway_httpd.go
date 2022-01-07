package util
//
//import (
//	"frame_sync/myproto"
//	"frame_sync/myprotocol"
//	"github.com/gorilla/websocket"
//	"go.uber.org/zap"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	"strconv"
//	"strings"
//	"context"
//	"zlib"
//)
//
//type MyServer struct {
//	Host 			string
//	Port 			string
//	MapSize 		int32
//	RoomPeople 		int32
//	Uri				string
//	OffLineWaitTime int32
//	ActionMap		map[string]map[int32]myprotocol.ActionMap
//	ContentType		int32
//	LoginAuthType	string
//	FPS 			int32
//}
//
////var httpUpGrader = websocket.Upgrader{
////	ReadBufferSize:  1024,
////	WriteBufferSize: 1024,
////	// 允许所有的CORS 跨域请求，正式环境可以关闭
////	CheckOrigin: func(r *http.Request) bool {
////		return true
////	},
////}
//
////type MyMetrics struct {
////	Rooms		int `json:"room"`
////	Players		int `json:"players"`
////	Conns 		int `json:"conns"`
////	InputNum 	int `json:"inputNum"`
////	InputSize 	int `json:"inputSize"`
////	OutputNum 	int `json:"outputNum"`
////	OutputSize 	int `json:"outputSize"`
////	InputErrNum int `json:"inputErrNum"`
////}
//
//type Httpd struct{
//	Log 			*zap.Logger 		`json:"-"`
//	httpServer		*http.Server	`json:"-"`
//	Option HttpdOption
//}
//
//type HttpdOption struct{
//	RootPath 		string	//监听的相对路径
//	BaseRootPath 	string	//监听的物理路径 = 当前项目物理路径 + 监听的相对路径
//	Ip 				string
//	Port 			string
//	WsUri 			string
//	ParentCtx  		context.Context `json:"-"`
//	LogOption 		*zap.Logger
//	WsNewFDBack		func(*websocket.Conn)
//}
//
////var RoomList	map[string]Room
//
//func NewHttpd(option HttpdOption)*Httpd{
//	//mylog.Info("NewHttpd instance")
//	//实例化一个log类，单独给httpd使用
//	//option.LogOption.OutFileFileName = "httpd"
//	//option.LogOption.ModuleId = 2
//	//newLog,errs  := zlib.NewLog(option.LogOption)
//	//if errs != nil{
//	//	zlib.PanicPrint("new log err",errs.Error())
//	//}
//
//	httpd := new(Httpd)
//	//httpd.Log = newLog
//	httpd.Option = option
//	//httpd.RootPath = option.RootPath
//	baseDir, _ := os.Getwd()
//	httpd.Option.BaseRootPath = baseDir + option.RootPath
//
//	return httpd
//}
//
//func  (httpd *Httpd)start(outCtx context.Context){
//	dns := httpd.Option.Ip + ":" + httpd.Option.Port
//	//mylog.Alert("httpd start:",dns)
//	logger := log.New(httpd.Log.Option.OutFileFileFd ,"h_s_err",log.Ldate)
//	httpd.httpServer = & http.Server{
//		Addr:dns,
//		ErrorLog: logger,
//	}
//	http.HandleFunc(httpd.Option.RootPath, myHttpd.wwwHandler)
//	//这里开始阻塞，直到接收到停止信号
//	err := httpd.httpServer.ListenAndServe()
//	if err != nil {
//		if strings.Index(err.Error(),"Server closed") == -1{
//			zlib.PanicPrint("httpd:"+err.Error())
//		}
//		mylog.Error(" httpd ListenAndServe err:", err.Error())
//	}
//}
////临时方法
//func  (httpd *Httpd)startWs(outCtx context.Context){
//	dns := httpd.Option.Ip + ":" + httpd.Option.Port
//	mylog.Alert("httpd start:",dns)
//	logger := log.New(httpd.Log.Option.OutFileFileFd ,"h_s_err",log.Ldate)
//	httpd.httpServer = & http.Server{
//		Addr:dns,
//		ErrorLog: logger,
//	}
//	http.HandleFunc(httpd.Option.WsUri ,httpd.wsHandler)
//	//这里开始阻塞，直到接收到停止信号
//	err := httpd.httpServer.ListenAndServe()
//	if err != nil {
//		if strings.Index(err.Error(),"Server closed") == -1{
//			zlib.PanicPrint("httpd:"+err.Error())
//		}
//		mylog.Error(" httpd ListenAndServe err:", err.Error())
//	}
//}
//
//
//func  (httpd *Httpd) shutdown(){
//	mylog.Alert(" shutdown httpd")
//	httpd.httpServer.Shutdown(httpd.Option.ParentCtx)
//
//	httpd.Log.CloseChan <- 1
//}
//func  (httpd *Httpd)wsHandler(w http.ResponseWriter, r *http.Request){
//	mylog.Info("wsHandler: have a new client http request")
//	//http 升级 ws
//	wsConnFD, err := httpUpGrader.Upgrade(w, r, nil)
//	mylog.Info("Upgrade this http req to websocket")
//	if err != nil {
//		mylog.Error("Upgrade websocket failed: ", err.Error())
//		return
//	}
//	httpd.Option.WsNewFDBack(wsConnFD)
//}
////所有HTTP 请求的入口
//func  (httpd *Httpd)wwwHandler(w http.ResponseWriter, r *http.Request){
//	defer func() {
//		if err := recover(); err != nil {
//			httpd.Log.Panic("wwwHandler:",err)
//		}
//	}()
//	uri := r.URL.RequestURI()
//	httpd.Log.Info("uri:",uri)
//	if uri == "" || uri == "/" {
//		httpd.ResponseStatusCode(w,500,"RequestURI is null or uir is :  '/'")
//		return
//	}
//	//zlib.MyPrint(r.Header)
//	uri = zlib.UriTurnPath(uri)
//	httpd.Log.Info("base uri:",uri)
//	query := r.URL.Query()//GET 方式URL 中的参数 转 结构体
//	httpd.Log.Info("query:",query)
//	var jsonStr []byte
//	if uri == "/www/getServer"{
//		options := myNetWay.Option
//		//options.Host = "39.106.65.76"
//
//		format := query.Get("format")
//		format = strings.Trim(format," ")
//		if format == ""{
//			jsonStr,_ = json.Marshal(&options)
//		}else if format == "proto"{
//			cfgServer := myproto.CfgServer{
//				ListenIp			:options.ListenIp,
//				OutIp				:options.OutIp,
//				WsPort				:options.WsPort,
//				TcpPort				:options.TcpPort,
//				UdpPort				:options.UdpPort,
//				ContentType			:options.ContentType,
//				LoginAuthType		:options.LoginAuthType,
//				LoginAuthSecretKey	:options.LoginAuthSecretKey,
//				IOTimeout			:options.IOTimeout,
//				ConnTimeout			:options.ConnTimeout,
//				Protocol			:options.Protocol,
//				WsUri				:options.WsUri,
//				MaxClientConnNum	:options.MaxClientConnNum,
//				RoomPeople			:options.RoomPeople,
//				RoomReadyTimeout 	:options.RoomReadyTimeout,
//				OffLineWaitTime		:options.OffLineWaitTime,//玩家掉线后，等待多久
//				MapSize				:options.MapSize,
//				LockMode			:options.LockMode,
//				FPS					:options.FPS,
//				Store				:options.Store,
//				HttpdRootPath		:options.HttpdRootPath,
//				HttpPort			:options.HttpPort,
//			}
//			jsonStr,_ = proto.Marshal(&cfgServer)
//		}
//
//	}else if uri == "/www/apilist"{
//		format := query.Get("format")
//		info := myProtocolActions.GetActionMap()
//		if format == ""{
//			jsonStr,_ = json.Marshal(&info)
//		}else if format == "proto"{
//			cfgServer := myproto.CfgProtocolActions{}
//			client := make(map[int32]*myproto.CfgActions)
//			server := make(map[int32]*myproto.CfgActions)
//			for k,v := range info["client"]{
//				client[k] = &myproto.CfgActions{
//					Id: v.Id,
//					Action: v.Action,
//					Desc: v.Desc,
//					Demo: v.Demo,
//				}
//			}
//			for k,v := range info["server"]{
//				client[k] = &myproto.CfgActions{
//					Id: v.Id,
//					Action: v.Action,
//					Desc: v.Desc,
//					Demo: v.Demo,
//				}
//			}
//
//			cfgServer.Client = client
//			cfgServer.Server = server
//			jsonStr,_ = proto.Marshal(&cfgServer)
//		}
//
//	}else if uri == "/www/getMetrics"{
//		pool:= myMetrics.Pool
//		pool["execTime"] = int(int(zlib.GetNowMillisecond()) - pool["starupTime"])
//		jsonStr,_ = json.Marshal(&pool)
//		zlib.MyPrint(string(jsonStr))
//	}else if uri == "/www/getFD"{
//
//	}else if uri == "/www/startUpDesc"{
//		cmdArgs := CmdArgs{}
//		typeOfCmsArgs := reflect.TypeOf(cmdArgs)
//		cmsArgAfter := make(map[string]string)
//		for i:=0;i<typeOfCmsArgs.NumField();i++{
//			memVar := typeOfCmsArgs.Field(i)
//			desc := memVar.Tag.Get("desc")
//			cmsArgAfter[memVar.Name] = desc
//		}
//
//		jsonStr,_ = json.Marshal(cmsArgAfter)
//
//	}else if uri == "/www/getRoomList"{
//		type RoomList struct {
//			Rooms map[string]Room              `json:"rooms"`
//			Metrics map[string]RoomSyncMetrics `json:"metrics"`
//		}
//
//		myroomList := make(map[string]Room)
//		roomListPoint := MySyncRoomPool
//		myRoomMetrics := make(map[string]RoomSyncMetrics)
//		//var emptyArr  []*ResponseRoomHistory
//		if len(roomListPoint) > 0 {
//			for k,v := range roomListPoint{
//				tt := *v
//				tt.LogicFrameHistory = nil
//				myroomList[k] = tt
//				myRoomMetrics[k] = RoomSyncMetricsPool[k]
//			}
//		}
//
//		roomList := RoomList{
//			Rooms: myroomList,
//			Metrics: myRoomMetrics,
//		}
//
//		jsonStr,_ = json.Marshal(&roomList)
//		//mylog.Debug("jsonStr:",jsonStr,err)
//	} else if uri == "/www/actionMap"{
//
//	}else if uri == "/www/getRoomOne"{
//		roomId := query.Get("id")
//		room := MySyncRoomPool[roomId]
//		//history := room.LogicFrameHistory
//		//type RoomOne struct {
//		//	Info Room 	`json:"info"`
//		//	HistoryList []*myproto.ResponseRoomHistory `json:"historyList"`
//		//}
//		//roonOne := RoomOne{}
//		//roonOne.Info = *room
//		//roonOne.HistoryList = history
//		//zlib.MyPrint(roonOne)
//		jsonStr,_ = json.Marshal(&room)
//	} else if uri == "/www/createJwtToken"{
//		randUid := query.Get("id")
//		randUid = strings.Trim(randUid," ")
//		if randUid == ""{
//			jsonStr = []byte( "id 为空")
//		}else{
//			uidStrConvInt32,_ := strconv.ParseInt(randUid,10,32)
//			payload := zlib.JwtDataPayload{
//				Uid:int32(uidStrConvInt32),
//				ATime:int32(zlib.GetNowMillisecond()),
//				AppId:2,
//			}
//			token := zlib.CreateJwtToken(myNetWay.Option.LoginAuthSecretKey,payload)
//			type CreateJwtNewToken struct {
//				Uid 	int32
//				Token 	string
//			}
//			createJwtNewToken := CreateJwtNewToken{
//				Uid: payload.Uid,
//				Token:token,
//			}
//			jsonStr,_ = json.Marshal(&createJwtNewToken)
//		}
//
//	} else if uri == "/www/testCreateJwtToken"{
//		//info := mynetWay.testCreateJwtToken()
//		//jsonStr,_ = json.Marshal(&info)
//	}else if uri == "/www/getProtoFile"{
//		filePath := "/myproto/api.proto"
//		fileContent, err := getStaticFileContent(filePath)
//		if err != nil{
//			httpd.Log.Error("/www/getProtoFile:",err.Error())
//		}
//		jsonStr = []byte(fileContent)
//	}else{
//		err := httpd.routeStatic(w,r,uri)
//		if err != nil{
//			mylog.Error("httpd.routeStatic err:",err.Error())
//		}
//		return
//	}
//
//	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
//	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
//
//	w.Header().Set("Content-Length",strconv.Itoa( len(jsonStr) ) )
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	w.Write(jsonStr)
//
//}
//
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
//
//func  (httpd *Httpd) routeStatic(w http.ResponseWriter,r *http.Request,uri string)error{
//	fileList := zlib.GetFileListByDir(httpd.Option.BaseRootPath)
//	if len(fileList) <=0 {
//		httpd.Log.Error("GetFileListByDir is empty")
//		return nil
//	}
//
//	for _,fileName := range fileList{
//		tmpFileName := httpd.Option.RootPath + fileName
//		if uri == tmpFileName{
//			fileContent, err := getStaticFileContent(uri)
//			if err != nil {
//				httpd.ResponseStatusCode(w, 404, err.Error())
//				return errors.New("routeStatic 404")
//			}
//			//踦域处理
//			w.Header().Set("Access-Control-Allow-Origin","*")
//			w.Header().Add("Access-Control-Allow-Headers","Content-Type")
//			//w.Header().Set("content-type","text/plain")
//
//			w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))
//			w.Write([]byte(fileContent))
//			return nil
//		}
//	}
//	httpd.Log.Error("routeStatic no match :",uri)
//	//if  uri == "/www/ws.html" ||
//	//	uri == "/www/sync_frame_client_server.jpg" ||
//	//	uri == "/www/jquery.min.js"||
//	//	uri == "/www/sync.js"||
//	//	uri == "/www/api_web_pb.js"||
//	//	uri == "/www/roomlist.html"||
//	//	uri == "/www/metrics.html"||
//	//	uri == "/www/serverUpVersionMemo.html"||
//	//	uri == "/www/sync_frame_client_server.jpg" ||
//	//	uri == "/www/rsync_frame_lock_step.jpg" ||
//	//	uri == "/www/index.html" ||
//	//	uri == "/www/config.html" ||
//	//	uri == "/www/roomDetail.html" ||
//	//	uri == "/www/apilist.html"{ //静态文件
//	//
//	//	fileContent, err := getStaticFileContent(uri)
//	//	if err != nil {
//	//		httpd.ResponseStatusCode(w, 404, err.Error())
//	//		return errors.New("routeStatic 404")
//	//	}
//	//	//踦域处理
//	//	w.Header().Set("Access-Control-Allow-Origin","*")
//	//	w.Header().Add("Access-Control-Allow-Headers","Content-Type")
//	//	//w.Header().Set("content-type","text/plain")
//	//
//	//	w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))
//	//	w.Write([]byte(fileContent))
//	//}
//	return nil
//}
//
////http 响应状态码
//func    (httpd *Httpd)ResponseStatusCode(w http.ResponseWriter,code int ,responseInfo string){
//	httpd.Log.Info("ResponseStatusCode",code,responseInfo)
//
//	w.Header().Set("Content-Length",strconv.Itoa( len(responseInfo) ) )
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(403)
//	w.Write([]byte(responseInfo))
//}
