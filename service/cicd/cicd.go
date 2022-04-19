package cicd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abrander/go-supervisord"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"zgoframe/model"
	"zgoframe/util"
)

/*
自动化部署，从DB中读取出所有信息基础信息，GIT CLONE 配置super visor 监听进程
依赖
	supervisor 依赖 python 、 xmlrpc
	代码依赖：git
	系统脚本：依赖shell
*/


//===============配置 结构体 开始===================
type ConfigCicdSystem struct {
	Env        			[]string
	LogDir     			string
	ServiceDir 			string
	MasterDirName		string
	GitCloneTmpDirName	string
	HttpPort 			string
}

type ConfigServiceCICDSystem struct {
	Startup        string
	ListeningPorts string
	TestUnit       string
	Build          string
	Command        string
	ConfigTmpFileName string
	ConfigFileName string
}

type ConfigServiceCICDDepend struct {
	Go    string
	Node  string
	Mysql string
	Redis string
}

type ConfigServiceCICD struct {
	System ConfigServiceCICDSystem
	Depend ConfigServiceCICDDepend
}

type ConfigCicdSuperVisor struct {
	RpcPort          string
	ConfTemplateFile string
	ConfDir          string
}

type ConfigCicd struct {
	System     ConfigCicdSystem
	SuperVisor ConfigCicdSuperVisor
}

//===============配置 结构体 结束===================



type CicdPublish struct {
	Id            int
	RegTime       int
	Status        int
	ServiceName   string
	Logs          []string
	TotalExecTime int
	Server        util.Server
}

type CicdManager struct {
	Option CicdManagerOption
}

type CicdManagerOption struct {
	ServerList  map[int]util.Server  //所有服务器
	ServiceList map[int]util.Service //所有项目/服务

	HttpPort        string
	InstanceManager *util.InstanceManager
	Config          ConfigCicd
	PublicManager   *CICDPublicManager
	Log             *zap.Logger
	OpDirName       string
	TestServerList 	[]string
}


func NewCicdManager(cicdManagerOption CicdManagerOption) (*CicdManager, error) {
	cicdManager := new(CicdManager)

	cicdManagerOption.TestServerList = []string{"127.0.0.1","8.142.177.235"}

	cicdManager.Option = cicdManagerOption

	_, err := util.PathExists(cicdManagerOption.Config.System.ServiceDir) //service 根目录
	if err != nil {
		return cicdManager, cicdManager.MakeError("Option.Config.System.ServiceDir :" + err.Error())
	}
	//SuperVisor 模板文件
	_, err = util.FileExist(cicdManager.Option.Config.SuperVisor.ConfTemplateFile) //superVisor 模板文件
	if err != nil {
		return cicdManager, cicdManager.MakeError("SuperVisor.ConfTemplateFile :" + err.Error())
	}
	//SuperVisor 配置文件统一放置目录
	_, err = util.PathExists(cicdManager.Option.Config.SuperVisor.ConfDir) //superVisor 配置文件统一放置目录
	if err != nil {
		return cicdManager, cicdManager.MakeError("SuperVisor.ConfDir :" + err.Error())
	}
	//日志统一放置目录
	_, err = util.PathExists(cicdManager.Option.Config.System.LogDir)
	if err != nil {
		return cicdManager, cicdManager.MakeError("System.LogDir :" + err.Error())
	}

	return cicdManager, nil
}

func (cicdManager *CicdManager) MakeError(errMsg string) error {
	cicdManager.Option.Log.Error(errMsg)
	return errors.New(errMsg)
}

//开始HTTP监听，供管理员UI可视化管理
func (cicdManager *CicdManager) StartHttp(staticDir string) {
	//HttpZapLog = zapLog
	//ginRouter := gin.Default()
	//单独的日志记录，GIN默认的日志不会持久化的
	//ginRouter.Use(ZapLog())
	//加载静态目录
	//	Router.Static("/form-generator", "./resource/page")
	//加载swagger api 工具
	//ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//设置跨域
	//ginRouter.Use(httpmiddleware.Cors())

	//ginRouter.StaticFS("/static", http.Dir(staticDir))

	//ginRouter.GET("/ping", cicdManager.Ping)
	//ginRouter.GET("/getServerList", cicdManager.GetServerList)
	//ginRouter.GET("/getInstanceList", cicdManager.GetInstanceList)
	//ginRouter.GET("/getPublishList", cicdManager.GetPublishList)
	//ginRouter.GET("/getSuperVisorList", cicdManager.GetSuperVisorList)
	//ginRouter.GET("/getServiceList", cicdManager.GetServiceList)
	//ginRouter.POST("/publish", cicdManager.Publish)

	//ginRouter.Run("0.0.0.0:" + cicdManager.Option.HttpPort)

	//404
	//ginRouter.NoMethod(HandleNotFound)
}

//发布|部署 一次 服务
func (cicdManager *CicdManager) Ping(c *gin.Context) {

	//str,_ := json.Marshal(cicdManager.Option.ServiceList)
	//c.String(200,string(str))
}

////发布|部署 一次 服务
//func (cicdManager *CicdManager) Publish(c *gin.Context) {
//
//	//str,_ := json.Marshal(cicdManager.Option.ServiceList)
//	//c.String(200,string(str))
//}

//在当前服务器上，从<部署目录>中检索出每个服务（目录名），分析出：哪些服务~已经部署
func (cicdManager *CicdManager) GetServiceList( ) map[int]util.Service {
	list := make(map[int]util.Service)

	dirList := util.ForeachDir(cicdManager.Option.Config.System.ServiceDir)
	for k, service := range cicdManager.Option.ServiceList {
		s := service
		for _, dirInfo := range dirList {
			if s.Name == dirInfo.Name {
				s.Deploy = 1
				break
			}
		}
		list[k] = s
	}
	return list
	//str, _ := json.Marshal(cicdManager.Option.ServiceList)
}

//获取所有服务器列表，并做ping，确定状态
func (cicdManager *CicdManager) GetServerList( ) map[int]util.Server{
	list := make(map[int]util.Server)
	for k, server := range cicdManager.Option.ServerList {
		//这里是测试代码，不然PING太慢
		if !cicdManager.CheckInTestServer(server.OutIp){
			server.Status = 3
		}else{
			//status := util.PingByShell(server.OutIp, "2")
			//if !status {
			//	server.Status = 3
			//}
			server.Status = 3
		}

		list[k] = server
	}

	return list

}

//获取所有3方服务列表
func (cicdManager *CicdManager) GetInstanceList(c *gin.Context) {
	for _, instance := range cicdManager.Option.InstanceManager.Pool {
		status := util.CheckIpPort(instance.Host, instance.Port, 2)
		if !status {
			instance.Status = 3
		}
	}

	str, _ := json.Marshal(cicdManager.Option.InstanceManager.Pool)
	c.String(200, string(str))

}

//获取所有 部署发布 记录列表，ps:未加分页
func (cicdManager *CicdManager) GetPublishList(c *gin.Context) {
	listArr, _ := cicdManager.Option.PublicManager.GetList()
	listMap := make(map[int]model.CicdPublish)
	for _, v := range listArr {
		listMap[v.Id] = v
	}
	str, _ := json.Marshal(listArr)
	c.String(200, string(str))
}

//创建一个新的结构体,主要是给前端返回结果使用
type ServerServiceSuperVisorList struct {
	ServerPingStatus            map[int]int	`json:"server_ping_status"`
	SuperVisorStatus            map[int]int `json:"super_visor_status"`
	ServerServiceSuperVisor map[int][]supervisord.ProcessInfo `json:"server_service_super_visor"`
}
//每台服务器上 都会启动一个superVisor进程
//列出每台机器上的：superVisor进程 的所有服务进程的状态信息
func (cicdManager *CicdManager) GetSuperVisorList( )(list ServerServiceSuperVisorList , err error) {
	//serviceBaseDir := cicdManager.Option.Config.System.ServiceDir

	if len(cicdManager.Option.ServerList) == 0 {
		//服务器 为空
		errMsg := "GetSuperVisorList err:ServerList is empty"
		//util.MyPrint()
		return list,errors.New(errMsg)
	}

	if len(cicdManager.Option.ServiceList) == 0 {
		//服务为空
		errMsg := "GetSuperVisorList err:ServiceList is empty"
		//util.MyPrint(errMsg)
		return list,errors.New(errMsg)
	}

	util.MyPrint("serverList len:", len(cicdManager.Option.ServerList), " ServiceList len:", len(cicdManager.Option.ServiceList))
	//服务器 上面:已经开启的 superVisor	map[serverId]=>superVisorList
	serverServiceSuperVisor := make(map[int][]supervisord.ProcessInfo)
	//服务器 状态
	serverStatus := make(map[int]int)
	superVisorStatus := make(map[int]int)
	for _, server := range cicdManager.Option.ServerList {
		fmt.Println("for each service:" + server.OutIp + " " + strconv.Itoa(server.Env))

		dns := "http://" + server.OutIp + ":" + cicdManager.Option.HttpPort
		//dns := "http://" + server.OutIp + ":9001"
		if cicdManager.CheckInTestServer(server.OutIp){
			//ping 测试一下 对端服务器：是否开启了sdk HTTP
			testServerRs := cicdManager.TestServerStateHttp(dns + "/cicd/ping")
			if testServerRs == 0 {
				util.MyPrint("")
				serverStatus[server.Id] = util.SERVER_PING_FAIL
				superVisorStatus[server.Id ] = util.SV_ERROR_INIT
				continue
			}else{
				serverStatus[server.Id] = util.SERVER_PING_OK
			}
		}else{
			superVisorStatus[server.Id ] = util.SV_ERROR_INIT
			serverStatus[server.Id] = util.SERVER_PING_FAIL
			continue
		}
		//创建实例
		superVisorOption := util.SuperVisorOption{
			Ip:               server.OutIp,
			RpcPort:          cicdManager.Option.Config.SuperVisor.RpcPort,
			ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
			ServiceName:      "",
			ConfDir:          cicdManager.Option.Config.SuperVisor.ConfDir,
		}
		serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
		if err != nil {
			util.MyPrint("NewSuperVisor err:", err)
			superVisorStatus[server.Id ] = util.SV_ERROR_INIT
			continue
		}
		//建立 XMLRpc
		err = serviceSuperVisor.InitXMLRpc()
		if err != nil {
			util.MyPrint("serviceSuperVisor InitXMLRpc err:", err)
			superVisorStatus[server.Id ] =  util.SV_ERROR_CONN
			continue
		}
		//获取当前机器上的superVisor 的所有 服务进程 状态
		processInfoList, err := serviceSuperVisor.Cli.GetAllProcessInfo()
		if err != nil{
			superVisorStatus[server.Id ] =  util.SV_ERROR_CONN
			util.MyPrint("SuperVisor GetAllProcessInfo err:"+err.Error())
			continue
		}else{
			//jsonStr,_  := json.Marshal(processInfoList)
			//util.MyPrint(jsonStr)
			//util.ExitPrint(string(jsonStr))
		}
		superVisorStatus[server.Id ] = util.SV_ERROR_NONE

		for _, service := range cicdManager.Option.ServiceList {
			//获取当前服务器，已部署的服务目录
			//servicePath := serviceBaseDir + util.DIR_SEPARATOR + service.Name
			//util.MyPrint("servicePath:", servicePath)

			superVisorProcessInfo := supervisord.ProcessInfo{
				Name:  util.SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX + service.Name,
				State: 999, //项目未部署过
			}
			//search := 0
			//var superVisorProcessInfo supervisord.ProcessInfo
			//先筛选一下，看看该服务有没有被 添加到superVisor中
			for _, process := range processInfoList {
				//rs := fmt.Sprintf("%+v",process)
				//util.MyPrint(rs)
				if process.Name == util.SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX+service.Name {
					util.MyPrint(fmt.Sprintf("%+v",process))
					superVisorProcessInfo = process
					//superVisorStatus[service.Id ] = util.SV_ERROR_NONE
					//search = 1
					break
				}
			}

			//if search == 0{
			//	superVisorStatus[server.Id ] = util.SV_ERROR_NOT_FOUND
			//	continue
			//}

			_, ok := serverServiceSuperVisor[server.Id]
			if !ok {
				//util.MyPrint(22222)
				serverServiceSuperVisor[server.Id] = []supervisord.ProcessInfo{superVisorProcessInfo}
			} else {
				//util.MyPrint(3333)
				serverServiceSuperVisor[server.Id] = append(serverServiceSuperVisor[server.Id], superVisorProcessInfo)
			}
		}

	}

	list = ServerServiceSuperVisorList{
		ServerPingStatus		: serverStatus,
		SuperVisorStatus		: superVisorStatus,
		ServerServiceSuperVisor	: serverServiceSuperVisor,
	}
	return list,nil
	//str, err := json.Marshal(myresponse)
	//util.MyPrint("json err:", err)
	//c.String(200, string(str))
}

func (cicdManager *CicdManager) GenerateAllFilebeat() {
	pwd, _ := os.Getwd() //当前路径
	opDirFull := pwd + "/" + cicdManager.Option.OpDirName

	for _, server := range cicdManager.Option.ServerList {
		cicdManager.GenerateFilebeat(server, opDirFull)
		util.MyPrint("finish one ...........")
	}
	util.ExitPrint(33)
}
func (cicdManager *CicdManager) GenerateFilebeat(server util.Server, opDir string) {

	instance, empty := cicdManager.Option.InstanceManager.GetByEnvName(server.Env, "es")
	if empty {
		util.ExitPrint("ProcessFilebeat GetByEnvName es empty :" + strconv.Itoa(server.Env))
	}

	esDns := instance.Host + ":" + instance.Port
	filebeatConfigFile := opDir + "/" + "filebeat.yaml"
	filebeatConfigFileContent, _ := util.ReadString(filebeatConfigFile)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#elasticsearch_output_hosts#", esDns, -1)

	filebeatInput := ""
	for _, service := range cicdManager.Option.ServiceList {
		filebeat_input_file := opDir + "/" + "filebeat_input.yaml"
		filebeat_input_content, _ := util.ReadString(filebeat_input_file)
		serviceLogDir := cicdManager.Option.Config.System.LogDir + "/" + service.Name + "/*.log"
		//ExitPrint(serviceLogDir)
		filebeat_input_content = strings.Replace(filebeat_input_content, "#paths#", serviceLogDir, -1)
		filebeat_input_content = strings.Replace(filebeat_input_content, "#source#", service.Name, -1)

		filebeatInput += filebeat_input_content + "\n"
	}
	esOutput := ""
	for _, service := range cicdManager.Option.ServiceList {
		esOutputFile := opDir + "/" + "filebeat_es_output.yaml"
		esOutputFileContent, _ := util.ReadString(esOutputFile)
		esOutputFileContent = strings.Replace(esOutputFileContent, "#index#", service.Name, -1)

		esOutput += esOutputFileContent + "\n"
	}
	//MyPrint(filebeatInput)
	//MyPrint(esOutput)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#filebeat_inputs#", filebeatInput, -1)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#elasticsearch_output_index#", esOutput, -1)

	util.MyPrint(filebeatConfigFileContent)
}

//如果一个路径不存在
func pathNotExistCreate(path string) error {
	_, err := util.PathExists(path)
	if err == nil { //目录存在
		return nil
	}
	if os.IsNotExist(err) { //目录不存在
		//创建一个目录
		err = os.Mkdir(path, 0777)
		fmt.Println("create path:", path)
		if err != nil {
			fmt.Println("create path failed , err:", err)
		}
		//return err
	} else { //其它错误
		//	fmt.Println("path :" + path + " exist , no need create.")
		//	return err
	}
	return err
}

//执行shell文件
func ExecShellFile(shellFile string, argc string) (string, error) {
	util.MyPrint("ExecShellFile:", shellFile, " ", argc)
	shellCommand := shellFile + " " + argc
	c := exec.Command("sh", "-c", shellCommand)

	output, err := c.CombinedOutput()
	if err != nil {
		util.MyPrint("exec.Command err:", err)
		return "", err
	}
	outStr := string(output)
	outArr := strings.Split(outStr, "\n")

	return outArr[1], nil
}

//执行shell 指令
func ExecShellCommand(command string, argc string) (string, error) {
	util.MyPrint("ExecShellCommand:", command, argc)
	//shellCommand := command + " " + argc
	c := exec.Command("bash", "-c", command)

	output, err := c.CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command err:", err)
		return "", err
	}
	outStr := string(output)
	//outArr := strings.Split(outStr,"\n")
	//return outArr[1],nil
	return outStr, nil
}

//这里有一条简单的操作，80端口基本上都得用，测试服务器状态，用ping curl 也可以.
func (cicdManager *CicdManager) TestServerStateHttp(hostUri string) int {
	client := &http.Client{
		Timeout:  1 * time.Second,
	}
	//提交请求
	reqest, err := http.NewRequest("GET", hostUri, nil)
	if err != nil {
		return 0
	}
	util.MyPrint("TestServerStateHttp uri:"+hostUri)
	//增加header选项
	reqest.Header.Add("X-Source-Type", "11")
	reqest.Header.Add("X-Project-Id", "6")
	reqest.Header.Add("X-Access", "imzgoframe")
	//处理返回结果
	response, err := client.Do(reqest)
	//defer response.Body.Close()
	if  err != nil{
		return 0
	}

	util.MyPrint("TestServerStateHttp http status:", response.Status)



	return 1

}

func (cicdManager *CicdManager) CheckInTestServer(ip string)bool{
	for _ , v := range cicdManager.Option.TestServerList{
		if ip == v {
			return true
		}
	}
	return false
}
