package util

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
)

//.proto 描述文件中：service 、函数名、结构体的定义，编译后，可以直接调用
//但是反向使用，就得用反射,如：长连接中，有接口的ID，根据ID找具体的一个服务下的一个方法，反射就有点麻烦
//长连接的通信时，为什么要用接口ID号？
//1. 压缩传输内容提升性能
//2. 加密
//所以，长连接中的内容传输是只有接口ID号，并没有具体的服务名及函数名的，这里就需要一个ID映射关系
//具体实现：
//1. 先用PHP 解析：.proto 描述文件，分析出：service_name 、func_name、请求参数、返回参数，并保存成.txt文件
//2. 读取上面生成的.txt文件内容，读入到GO的内存(结构体)中

//待解决：
//1. 借助PHP生成txt，可考虑用GO来写
//2. 不要中转.txt 直接用反射直接读到内存中

type ProtoMap struct {
	//并发安全的，因为无写操作
	ServiceFuncMap map[int]ProtoServiceFunc
	//ServiceAction  map[int][]*ActionMap
	Log            *zap.Logger
	ConfigFileDir  string
	MapFileName    string `json:"map_file_name"`
	ProjectManager *ProjectManager
}

type ProtoServiceFunc struct {
	ServiceName string `json:"service_name"`
	ServiceId   int    `json:"service_id"`
	Id          int    `json:"id"`
	FuncId      int    `json:"func_id"`
	FuncName    string `json:"func_name"`
	Request     string `json:"request"`
	Response    string `json:"response"`
	Desc        string `json:"desc"`
}

//var actionMap  	map[string]map[int]ActionMap
func NewProtoMap(log *zap.Logger, configFileDir string, MapFileName string, projectManager *ProjectManager) (*ProtoMap, error) {
	log.Info("NewProtobufMap:" + configFileDir)

	protoMap := new(ProtoMap)
	protoMap.Log = log
	protoMap.ConfigFileDir = configFileDir
	protoMap.MapFileName = MapFileName
	protoMap.ProjectManager = projectManager

	err := protoMap.initProtocolActionMap()
	log.Info("protoMap.ServiceFuncMap len:" + strconv.Itoa(len(protoMap.ServiceFuncMap)))

	return protoMap, err
}

func (protoMap *ProtoMap) initProtocolActionMap() error {
	//netway.mylog.Info("initActionMap")
	protoMap.ServiceFuncMap = make(map[int]ProtoServiceFunc)
	mapList, err := protoMap.loadingActionMapConfigFile(protoMap.MapFileName)
	if err != nil {
		return err
	}
	//protobufMap.ServiceAction = make(map[int][]*ActionMap)
	//for _, v := range mapList {
	//	_, ok := protobufMap.ServiceAction[v.ServiceId]
	//	if ok {
	//		protobufMap.ServiceAction[v.ServiceId] = append(protobufMap.ServiceAction[v.ServiceId], &v)
	//	} else {
	//		protobufMap.ServiceAction[v.ServiceId] = []*ActionMap{&v}
	//	}
	//}
	protoMap.ServiceFuncMap = mapList

	return err
}

func (protoMap *ProtoMap) loadingActionMapConfigFile(fileName string) (map[int]ProtoServiceFunc, error) {
	//_, _,_,dir  := getInfo(1)
	//ExitPrint(protobufMap.ConfigFileDir,fileName)
	pathFile := protoMap.ConfigFileDir + "/" + fileName
	protoMap.Log.Info("protobufMap loadingActionMapConfigFile:" + pathFile)
	fileContentArr, err := ReadLine(pathFile)
	if err != nil {
		protoMap.Log.Error("initActionMap ReadLine err :" + err.Error())
		return nil, err
		//protobufMap.Log.Panic("initActionMap ReadLine err :" + err.Error())
	}
	am := make(map[int]ProtoServiceFunc)
	for _, v := range fileContentArr {
		contentArr := strings.Split(v, "|")
		//MyPrint(contentArr[0],contentArr[1])
		if len(contentArr) < 5 {
			protoMap.Log.Error("read line len < 5")
			continue
		}
		serviceIdStr := contentArr[0][0:2]
		serviceId, _ := strconv.Atoi(serviceIdStr)

		id := Atoi(contentArr[0])
		serviceName := contentArr[1]

		funcIdStr := contentArr[0][2:]
		funcId, _ := strconv.Atoi(funcIdStr)

		//这里有BUG，回头处理
		//map txt 里都是首字节大写，这里转成小写
		//lowServiceName :=StrFirstToLower(serviceName)
		//_, empty := protoMap.ProjectManager.GetByName(serviceName)

		//if empty {
		//	return nil, errors.New("serviceName not in project list :" + serviceName)
		//}

		//id ,_:= strconv.Atoi(contentArr[0])
		actionMap := ProtoServiceFunc{
			Id:          id,
			ServiceId:   serviceId,
			ServiceName: serviceName,
			FuncId:      funcId,
			FuncName:    contentArr[2],
			Request:     contentArr[3],
			Response:    contentArr[4],
			Desc:        contentArr[5],
		}
		//PrintStruct(actionMap,":")
		//ExitPrint(111)
		am[id] = actionMap
	}
	if len(am) <= 0 {
		protoMap.Log.Error("protocolActions len(am) <= 0")
		panic("protocolActions len(am) <= 0")
	}
	return am, nil
}

////获取上层调用者的文件位置
//func getInfo(depth int) (funcName, fileName string, lineNo int, dir string) {
//	pc, file, lineNo, ok := runtime.Caller(depth)
//	if !ok {
//		fmt.Println("runtime.Caller() failed")
//		return
//	}
//	funcName = runtime.FuncForPC(pc).Name()
//	fileName = path.Base(file) // Base函数返回路径的最后一个元素
//
//	i := strings.LastIndex(file, "/")
//	//if i < 0 {
//	//	i = strings.LastIndex(path, "\\")
//	//}
//	//if i < 0 {
//	//	return "", errors.New(`error: Can't find "/" or "\".`)
//	//}
//	dir = string(file[0 : i+1])
//	return
//}

//获取全部列表数据
func (protoMap *ProtoMap) GetServiceFuncMap() map[int]ProtoServiceFunc {
	return protoMap.ServiceFuncMap
}

//根据函数名，获取一条记录，这种方法不太严谨，按说应该 服务名+函数名，保证唯一 ，但代码已经写了，先这样，后期优化
func (protoMap *ProtoMap) GetServiceFuncByFuncName(funcName string) (protoServiceFunc ProtoServiceFunc, empty bool) {
	for _, v := range protoMap.ServiceFuncMap {
		if v.FuncName == funcName {
			return v, false
		}
	}
	return protoServiceFunc, true
}

func (protoMap *ProtoMap) GetServiceFuncById(id int) (protoServiceFunc ProtoServiceFunc, empty bool) {
	am, ok := protoMap.ServiceFuncMap[id]
	if ok {
		return am, false
	} else {
		return am, true
	}
}

func (protoMap *ProtoMap) GetIdByMergeSidFid(serviceId int, funcId int) int {
	id, _ := strconv.Atoi(strconv.Itoa(serviceId) + strconv.Itoa(funcId))
	return id
}

func (protoMap *ProtoMap) GetIdByMergeStringSidFid(serviceId string, funcId string) int {
	id, _ := strconv.Atoi(serviceId + funcId)
	return id
}

//func (protoMap *ProtoMap) GetServiceFuncById(id int) (actionMapT ActionMap, empty bool) {
//	//protobufMap.Log.Info("GetActionId " + action)
//	am := protoMap.ServiceFuncMap
//	for id,  := range am {
//		if id  == action {
//			return v, false
//		}
//	}
//	return actionMapT, true
//}
