package util
import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type ProtobufMap struct {
	//并发安全的，因为无写操作
	ActionMaps 		map[int]ActionMap
	ServiceAction 	map[int][]*ActionMap
	Log 			*zap.Logger
	ConfigFileDir 	string
	MapFileName 	string 	`json:"map_file_name"`
	ProjectManager *ProjectManager
}

type ActionMap struct {
	ServiceName string 	`json:"service_name"`
	ServiceId 	int		`json:"service_id"`
	Id 			int		`json:"id"`
	Action		string	`json:"action"`
	Desc 		string	`json:"desc"`
	Request		string	`json:"demo"`
	Response 	string 	`json:"response"`

}

//var actionMap  	map[string]map[int]ActionMap
func NewProtobufMap(log *zap.Logger,configFileDir string,MapFileName string, projectManager *ProjectManager)(*ProtobufMap,error) {
	log.Info("NewProtobufMap:"+configFileDir)

	protobufMap := new(ProtobufMap)
	protobufMap.Log = log
	protobufMap.ConfigFileDir = configFileDir
	protobufMap.MapFileName = MapFileName
	protobufMap.ProjectManager = projectManager

	err := protobufMap.initProtocolActionMap()
	return protobufMap,err
}

func (protobufMap *ProtobufMap)initProtocolActionMap()error{
	//netway.mylog.Info("initActionMap")
	protobufMap.ActionMaps = make( 	map[int]ActionMap)
	mapList ,err := protobufMap.loadingActionMapConfigFile(protobufMap.MapFileName)
	if err != nil{
		return err
	}
	protobufMap.ServiceAction =  make( 	map[int][]*ActionMap)
	for _,v:=range mapList{
		_,ok := protobufMap.ServiceAction[v.ServiceId]
		if ok {
			protobufMap.ServiceAction[v.ServiceId] = append(protobufMap.ServiceAction[v.ServiceId] , &v)
		}else{
			protobufMap.ServiceAction[v.ServiceId] = []*ActionMap{&v}
		}
	}

	protobufMap.ActionMaps = mapList

	return err
}
func (protobufMap *ProtobufMap)loadingActionMapConfigFile(fileName string)(map[int]ActionMap,error) {
	//_, _,_,dir  := getInfo(1)
	//ExitPrint(protobufMap.ConfigFileDir,fileName)
	fileContentArr,err := ReadLine(protobufMap.ConfigFileDir +"/"+fileName)
	if err != nil{
		protobufMap.Log.Error("initActionMap ReadLine err :" + err.Error())
		return nil,err
		//protobufMap.Log.Panic("initActionMap ReadLine err :" + err.Error())
	}
	am := make(map[int]ActionMap)
	for _,v:= range fileContentArr{
		contentArr := strings.Split(v,"|")
		//MyPrint(contentArr[0],contentArr[1])
		if len(contentArr)  <  5{
			protobufMap.Log.Error("read line len < 5")
			continue
		}
		serviceIdStr := contentArr[0][0:3]
		serviceId ,_ := strconv.Atoi(serviceIdStr)
		funcIdStr  := contentArr[0][3:]
		funcId :=  Atoi(funcIdStr)
		//id :=  Atoi(contentArr[0])
		//1000|Login|RequestLogin|ResponseLoginRes|登陆
		serviceName := contentArr[1]

		//map txt 里都是首字节大写，这里转成小写
		lowServiceName :=StrFirstToLower(serviceName)
		_, empty := protobufMap.ProjectManager.GetByKey(lowServiceName)
		if empty{
			return nil,errors.New("serviceName not in project list :" + serviceName)
		}

		actionMap := ActionMap{
			ServiceId : serviceId,
			ServiceName: serviceName,
			Id: funcId,
			Action: contentArr[2],
			Request: contentArr[3],
			Response: contentArr[4],
			Desc: contentArr[5],
		}

		am[funcId] = actionMap
	}
	if len(am) <= 0{
		protobufMap.Log.Error("protocolActions len(am) <= 0")
		panic("protocolActions len(am) <= 0")
	}
	return am,nil
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
