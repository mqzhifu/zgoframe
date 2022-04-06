package util

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
	"strings"
	"time"
	"zgoframe/model"
)

/*
自动化部署，从DB中读取出所有信息基础信息，GIT CLONE 配置super visor 监听进程
依赖
	supervisor 依赖 python 、 xmlrpc
	代码依赖：git
	系统脚本：依赖shell
*/
const (
	DIR_SEPARATOR = "/"
	STR_SEPARATOR = "#"
)

//===============配置 结构体 开始===================
type ConfigCicdSystem struct {
	Env        []string
	LogDir     string
	ServiceDir string
}

type ConfigServiceCICDSystem struct {
	Startup        string
	ListeningPorts string
	TestUnit       string
	Build          string
	Command        string
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

type SuperVisorReplace struct {
	script_name            string
	startup_script_command string
	script_work_dir        string
	stdout_logfile         string
	stderr_logfile         string
	process_name           string
}

//==============superVisor===========

type CicdPublish struct {
	Id            int
	RegTime       int
	Status        int
	ServiceName   string
	Logs          []string
	TotalExecTime int
	Server        Server
}

type CicdManager struct {
	Option CicdManagerOption
}

type CicdManagerOption struct {
	ServerList  map[int]Server  //所有服务器
	ServiceList map[int]Service //所有项目/服务

	HttpPort        string
	InstanceManager *InstanceManager
	Config          ConfigCicd
	PublicManager   *CICDPublicManager
	Log             *zap.Logger
	OpDirName       string
}

func NewCicdManager(cicdManagerOption CicdManagerOption) (*CicdManager, error) {
	cicdManager := new(CicdManager)
	cicdManager.Option = cicdManagerOption

	_, err := PathExists(cicdManagerOption.Config.System.ServiceDir) //service 根目录
	if err != nil {
		return cicdManager, cicdManager.MakeError("Option.Config.System.ServiceDir :" + err.Error())
	}
	//SuperVisor 模板文件
	_, err = FileExist(cicdManager.Option.Config.SuperVisor.ConfTemplateFile) //superVisor 模板文件
	if err != nil {
		return cicdManager, cicdManager.MakeError("SuperVisor.ConfTemplateFile :" + err.Error())
	}
	//SuperVisor 配置文件统一放置目录
	_, err = PathExists(cicdManager.Option.Config.SuperVisor.ConfDir) //superVisor 配置文件统一放置目录
	if err != nil {
		return cicdManager, cicdManager.MakeError("SuperVisor.ConfDir :" + err.Error())
	}
	//日志统一放置目录
	_, err = PathExists(cicdManager.Option.Config.System.LogDir)
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
	ginRouter := gin.Default()
	//单独的日志记录，GIN默认的日志不会持久化的
	//ginRouter.Use(ZapLog())
	//加载静态目录
	//	Router.Static("/form-generator", "./resource/page")
	//加载swagger api 工具
	//ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//设置跨域
	//ginRouter.Use(httpmiddleware.Cors())

	ginRouter.StaticFS("/static", http.Dir(staticDir))

	ginRouter.GET("/ping", cicdManager.Ping)
	ginRouter.GET("/getServerList", cicdManager.GetServerList)
	ginRouter.GET("/getInstanceList", cicdManager.GetInstanceList)
	ginRouter.GET("/getPublishList", cicdManager.GetPublishList)
	ginRouter.GET("/getSuperVisorList", cicdManager.GetSuperVisorList)
	ginRouter.GET("/getServiceList", cicdManager.GetServiceList)

	ginRouter.POST("/publish", cicdManager.Publish)

	ginRouter.Run("0.0.0.0:" + cicdManager.Option.HttpPort)

	//404
	//ginRouter.NoMethod(HandleNotFound)
}

//发布|部署 一次 服务
func (cicdManager *CicdManager) Ping(c *gin.Context) {

	//str,_ := json.Marshal(cicdManager.Option.ServiceList)
	//c.String(200,string(str))
}

//发布|部署 一次 服务
func (cicdManager *CicdManager) Publish(c *gin.Context) {

	//str,_ := json.Marshal(cicdManager.Option.ServiceList)
	//c.String(200,string(str))
}

//获取所有服务列表
func (cicdManager *CicdManager) GetServiceList(c *gin.Context) {
	dirList := ForeachDir(cicdManager.Option.Config.System.ServiceDir)
	for _, service := range cicdManager.Option.ServiceList {
		for _, dirInfo := range dirList {
			if service.Name == dirInfo.Name {
				service.Deploy = 1
				break
			}
		}
	}
	str, _ := json.Marshal(cicdManager.Option.ServiceList)
	c.String(200, string(str))
}

//获取所有服务器列表
func (cicdManager *CicdManager) GetServerList(c *gin.Context) {
	for _, server := range cicdManager.Option.ServerList {
		status := PingByShell(server.OutIp, "2")
		if !status {
			server.Status = 3
		}
	}
	str, _ := json.Marshal(cicdManager.Option.ServerList)
	c.String(200, string(str))

}

//获取所有3方服务列表
func (cicdManager *CicdManager) GetInstanceList(c *gin.Context) {
	for _, instance := range cicdManager.Option.InstanceManager.Pool {
		status := CheckIpPort(instance.Host, instance.Port, 2)
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

//每台服务器上 都会启动一个superVisor进程
//列出每台机器上的：superVisor进程 的所有服务进程的状态信息
func (cicdManager *CicdManager) GetSuperVisorList(c *gin.Context) {
	serviceBaseDir := cicdManager.Option.Config.System.ServiceDir

	if len(cicdManager.Option.ServerList) == 0 {
		MyPrint("GetSuperVisorList err:ServerList is empty")
		return
	}

	if len(cicdManager.Option.ServiceList) == 0 {
		MyPrint("GetSuperVisorList err:ServiceList is empty")
		return
	}

	MyPrint("serverList len:", len(cicdManager.Option.ServerList), " ServiceList len:", len(cicdManager.Option.ServiceList))

	//serverId=>superVisorList
	serverServiceSuperVisor := make(map[int][]supervisord.ProcessInfo)
	serverStatus := make(map[int]int)
	for _, server := range cicdManager.Option.ServerList {
		fmt.Println("for each service:" + server.OutIp + " " + server.Env)

		dns := "http://" + server.OutIp + ":" + cicdManager.Option.HttpPort
		//ping 测试一下 其它机器是否开启了sdk HTTP
		testServerRs := cicdManager.TestServerStateHttp(dns + "/ping")
		if testServerRs == 0 {
			MyPrint("")
			serverStatus[server.Id] = 3
			continue
		}
		//创建实例
		superVisorOption := SuperVisorOption{
			Ip:               server.OutIp,
			RpcPort:          cicdManager.Option.Config.SuperVisor.RpcPort,
			ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
			ServiceName:      "",
			ConfDir:          cicdManager.Option.Config.SuperVisor.ConfDir,
		}
		serviceSuperVisor, err := NewSuperVisor(superVisorOption)
		if err != nil {
			MyPrint("NewSuperVisor err:", err)
			serverStatus[server.Id] = 4
			continue
		}
		//建立 XMLRpc
		err = serviceSuperVisor.InitXMLRpc()
		if err != nil {
			MyPrint("serviceSuperVisor InitXMLRpc err:", err)
			serverStatus[server.Id] = 4
			continue
		}

		serverStatus[server.Id] = server.Status

		processInfoList, _ := serviceSuperVisor.Cli.GetAllProcessInfo()
		for _, service := range cicdManager.Option.ServiceList {
			servicePath := serviceBaseDir + DIR_SEPARATOR + service.Name
			MyPrint("servicePath:", servicePath)

			superVisorProcessInfo := supervisord.ProcessInfo{
				Name:  SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX + service.Name,
				State: 999, //项目未部署过
			}

			for _, process := range processInfoList {
				if process.Name == SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX+service.Name {
					superVisorProcessInfo = process
					break
				}
			}

			_, ok := serverServiceSuperVisor[server.Id]
			if !ok {
				MyPrint(22222)
				serverServiceSuperVisor[server.Id] = []supervisord.ProcessInfo{superVisorProcessInfo}
			} else {
				MyPrint(3333)
				serverServiceSuperVisor[server.Id] = append(serverServiceSuperVisor[server.Id], superVisorProcessInfo)
			}
		}

	}

	type response struct {
		ServerStatus            map[int]int                       `json:"server_status"`
		ServerServiceSuperVisor map[int][]supervisord.ProcessInfo `json:"server_service_super_visor"`
	}

	myresponse := response{
		ServerStatus:            serverStatus,
		ServerServiceSuperVisor: serverServiceSuperVisor,
	}
	str, err := json.Marshal(myresponse)
	MyPrint("json err:", err)
	c.String(200, string(str))
}

func (cicdManager *CicdManager) GenerateAllFilebeat() {
	pwd, _ := os.Getwd() //当前路径
	opDirFull := pwd + "/" + cicdManager.Option.OpDirName

	for _, server := range cicdManager.Option.ServerList {
		cicdManager.GenerateFilebeat(server, opDirFull)
		MyPrint("finish one ...........")
	}
	ExitPrint(33)
}
func (cicdManager *CicdManager) GenerateFilebeat(server Server, opDir string) {

	instance, empty := cicdManager.Option.InstanceManager.GetByEnvName(server.Env, "es")
	if empty {
		ExitPrint("ProcessFilebeat GetByEnvName es empty :" + server.Env)
	}

	esDns := instance.Host + ":" + instance.Port
	filebeatConfigFile := opDir + "/" + "filebeat.yaml"
	filebeatConfigFileContent, _ := ReadString(filebeatConfigFile)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#elasticsearch_output_hosts#", esDns, -1)

	filebeatInput := ""
	for _, service := range cicdManager.Option.ServiceList {
		filebeat_input_file := opDir + "/" + "filebeat_input.yaml"
		filebeat_input_content, _ := ReadString(filebeat_input_file)
		serviceLogDir := cicdManager.Option.Config.System.LogDir + "/" + service.Name + "/*.log"
		//ExitPrint(serviceLogDir)
		filebeat_input_content = strings.Replace(filebeat_input_content, "#paths#", serviceLogDir, -1)
		filebeat_input_content = strings.Replace(filebeat_input_content, "#source#", service.Name, -1)

		filebeatInput += filebeat_input_content + "\n"
	}
	esOutput := ""
	for _, service := range cicdManager.Option.ServiceList {
		esOutputFile := opDir + "/" + "filebeat_es_output.yaml"
		esOutputFileContent, _ := ReadString(esOutputFile)
		esOutputFileContent = strings.Replace(esOutputFileContent, "#index#", service.Name, -1)

		esOutput += esOutputFileContent + "\n"
	}
	//MyPrint(filebeatInput)
	//MyPrint(esOutput)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#filebeat_inputs#", filebeatInput, -1)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#elasticsearch_output_index#", esOutput, -1)

	MyPrint(filebeatConfigFileContent)
}

//如果一个路径不存在
func pathNotExistCreate(path string) error {
	_, err := PathExists(path)
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
	MyPrint("ExecShellFile:", shellFile, " ", argc)
	shellCommand := shellFile + " " + argc
	c := exec.Command("sh", "-c", shellCommand)

	output, err := c.CombinedOutput()
	if err != nil {
		MyPrint("exec.Command err:", err)
		return "", err
	}
	outStr := string(output)
	outArr := strings.Split(outStr, "\n")

	return outArr[1], nil
}

//执行shell 指令
func ExecShellCommand(command string, argc string) (string, error) {
	MyPrint("ExecShellCommand:", command, argc)
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
	httpClient := http.Client{
		Timeout: time.Second * 1,
	}

	resp, err := httpClient.Get(hostUri)
	if err != nil {
		MyPrint("http get err:", err)
		return 0
	}

	MyPrint("http get status:", resp.Status)

	return 1

}
