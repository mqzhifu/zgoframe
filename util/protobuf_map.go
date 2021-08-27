package util
import (
	"fmt"
	"go.uber.org/zap"
	"path"
	"runtime"
	"strings"
)

type ProtobufMap struct {
	//并发安全的，因为无写操作
	ActionMaps map[int]ActionMap
	Log *zap.Logger
	ConfigFileDir string
}

type ActionMap struct {
	Id 			int		`json:"id"`
	Action		string	`json:"action"`
	Desc 		string	`json:"desc"`
	Request		string	`json:"demo"`
	Response 	string 	`json:"response"`
}

//var actionMap  	map[string]map[int]ActionMap
func NewProtobufMap(log *zap.Logger,configFileDir string)*ProtobufMap {
	log.Info("NewProtobufMap")

	log.Info("NewProtocolActions instance:")
	protobufMap := new(ProtobufMap)
	protobufMap.Log = log
	protobufMap.ConfigFileDir = configFileDir
	protobufMap.initProtocolActionMap()
	return protobufMap
}

func (protobufMap *ProtobufMap)initProtocolActionMap(){
	//netway.mylog.Info("initActionMap")
	protobufMap.ActionMaps = make( 	map[int]ActionMap)
	protobufMap.ActionMaps = protobufMap.loadingActionMapConfigFile("map.txt")
	//actionMap["server"] = protobufMap.loadingActionMapConfigFile("serverActionMap.txt")
}
func (protobufMap *ProtobufMap)loadingActionMapConfigFile(fileName string)map[int]ActionMap {
	//_, _,_,dir  := getInfo(1)
	//ExitPrint(protobufMap.ConfigFileDir,fileName)
	fileContentArr,err := ReadLine(protobufMap.ConfigFileDir +"/"+fileName)
	if err != nil{
		protobufMap.Log.Error("initActionMap ReadLine err :" + err.Error())
		protobufMap.Log.Panic("initActionMap ReadLine err :" + err.Error())
	}
	am := make(map[int]ActionMap)
	for _,v:= range fileContentArr{
		contentArr := strings.Split(v,"|")
		//MyPrint(contentArr[0],contentArr[1])
		if len(contentArr)  <  5{
			protobufMap.Log.Error("read line len < 5")
			continue
		}
		id :=  Atoi(contentArr[0])
		//1000|Login|RequestLogin|ResponseLoginRes|登陆
		actionMap := ActionMap{
			Id: id,
			Action: contentArr[1],
			Request: contentArr[2],
			Response: contentArr[3],
			Desc: contentArr[4],

		}
		am[id] = actionMap
	}
	if len(am) <= 0{
		protobufMap.Log.Error("protocolActions len(am) <= 0")
		PanicPrint("protocolActions len(am) <= 0")
	}
	return am
}
//获取上层调用者的文件位置
func getInfo(depth int) (funcName, fileName string, lineNo int ,dir string) {
	pc, file, lineNo, ok := runtime.Caller(depth)
	if !ok {
		fmt.Println("runtime.Caller() failed")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file) // Base函数返回路径的最后一个元素

	i := strings.LastIndex(file, "/")
	//if i < 0 {
	//	i = strings.LastIndex(path, "\\")
	//}
	//if i < 0 {
	//	return "", errors.New(`error: Can't find "/" or "\".`)
	//}
	dir = string(file[0 : i+1])
	return
}

func(protobufMap *ProtobufMap)GetActionMap() map[int]ActionMap {
	return protobufMap.ActionMaps
}

func(protobufMap *ProtobufMap)GetActionName(id int)(actionMapT ActionMap,empty bool){
	am , ok := protobufMap.ActionMaps[id]
	if ok {
		return am,false
	}else{
		return am,true
	}
}

func (protobufMap *ProtobufMap)GetActionId(action string )(actionMapT ActionMap,empty bool){
	//netway.mylog.Info("GetActionId ",action , " ",category)
	am := protobufMap.ActionMaps
	for _,v:=range am{
		if v.Action == action {
			return v,false
		}
	}
	return  actionMapT,true
}
