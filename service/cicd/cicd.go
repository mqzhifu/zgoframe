// 自动化：持续集成/持续交付/持续部署
package cicd

import (
	"errors"
	"fmt"
	"github.com/abrander/go-supervisord"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"zgoframe/http/request"
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

// ===============配置 结构体 开始===================
type ConfigCicdSystem struct {
	Env               []string
	LogDir            string
	WorkBaseDir       string
	RemoteBaseDir     string
	RemoteUploadDir   string
	RemoteDownloadDir string
	// UploadPath      string
	// DownloadPath string
	RootDir string
	// ServiceDir 			string	//远程部署目录
	// LocalServiceDir		string 	//本地部署目录
	MasterDirName      string
	GitCloneTmpDirName string
	HttpPort           string
}

type ConfigServiceCICDSystem struct {
	Startup           string
	ListeningPorts    string
	TestUnit          string
	Build             string
	Command           string
	ConfigTmpFileName string
	ConfigFileName    string
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
	RpcPort              string
	ConfTemplateFile     string
	ConfTemplateFileName string
	ConfDir              string
}

type ConfigCicd struct {
	System     ConfigCicdSystem
	SuperVisor ConfigCicdSuperVisor
}

// ===============配置 结构体 结束===================

// 创建一个新的结构体,主要是给前端返回结果使用
type ServerServiceSuperVisorList struct {
	// ServerPingStatus        map[int]int             `json:"server_ping_status"`
	SuperVisorStatus        map[int]int             `json:"super_visor_status"`
	ServerServiceSuperVisor map[int][]MyProcessInfo `json:"server_service_super_visor"`
}

type LocalServerServiceList struct {
	ServerList  map[int]util.Server   `json:"server_list"`
	ServiceList map[int]model.Project `json:"service_list"`
}

type MyProcessInfo struct {
	ServiceId int    `json:"service_id"`
	MasterSrc string `json:"master_src"`
	supervisord.ProcessInfo
}

// ============================

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
	Deploy *Deploy
}

type CicdManagerOption struct {
	ServerList       map[int]util.Server  // 所有服务器
	ServiceList      map[int]util.Service // 所有项目/服务
	ProjectList      map[int]model.Project
	HttpPort         string
	InstanceManager  *util.InstanceManager
	Config           ConfigCicd
	PublicManager    *CICDPublicManager
	Log              *zap.Logger
	OpDirName        string
	TestServerList   []string
	UploadDiskPath   string
	DownloadDiskPath string
}

func NewCicdManager(cicdManagerOption CicdManagerOption) (*CicdManager, error) {
	cicdManager := new(CicdManager)

	cicdManagerOption.TestServerList = []string{"127.0.0.1", "8.142.177.235"}

	cicdManager.Option = cicdManagerOption
	cicdManager.Deploy = NewDeploy(cicdManagerOption)
	// _, err := util.PathExists(cicdManagerOption.Config.System.ServiceDir) //service 根目录
	_, err := util.PathExists(cicdManagerOption.Config.System.WorkBaseDir)
	if err != nil {
		return cicdManager, cicdManager.Deploy.MakeError("Option.Config.System.ServiceDir :" + err.Error())
	}
	// SuperVisor 模板文件
	_, err = util.FileExist(cicdManager.Option.Config.SuperVisor.ConfTemplateFile) // superVisor 模板文件
	if err != nil {
		return cicdManager, cicdManager.Deploy.MakeError("SuperVisor.ConfTemplateFile :" + err.Error())
	}
	// SuperVisor 配置文件统一放置目录
	_, err = util.PathExists(cicdManager.Option.Config.SuperVisor.ConfDir) // superVisor 配置文件统一放置目录
	if err != nil {
		return cicdManager, cicdManager.Deploy.MakeError("SuperVisor.ConfDir :" + err.Error())
	}
	// 日志统一放置目录
	_, err = util.PathExists(cicdManager.Option.Config.System.LogDir)
	if err != nil {
		return cicdManager, cicdManager.Deploy.MakeError("System.LogDir :" + err.Error())
	}

	return cicdManager, nil
}

// 获取所有 部署发布 记录列表，ps:未加分页
func (cicdManager *CicdManager) GetPublishList(limit int) map[int]model.CicdPublish {
	listArr, _ := cicdManager.Option.PublicManager.GetList(limit)
	listMap := make(map[int]model.CicdPublish)
	for _, v := range listArr {
		listMap[v.Id] = v
	}
	return listMap
	// str, _ := json.Marshal(listArr)
	// c.String(200, string(str))
}
func (cicdManager *CicdManager) SuperVisorProcess(form request.CicdSuperVisor) (err error) {
	server, service, err := cicdManager.Deploy.CheckCicdRequestForm(form.CicdDeploy)
	if err != nil {
		return err
	}
	if form.Command == "" {
		return errors.New("command is empty")
	}
	// 创建实例
	superVisorOption := util.SuperVisorOption{
		Ip:      server.OutIp,
		RpcPort: cicdManager.Option.Config.SuperVisor.RpcPort,
	}

	serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
	// 建立 XMLRpc
	err = serviceSuperVisor.InitXMLRpc()
	if err != nil {
		return errors.New("serviceSuperVisor InitXMLRpc err:" + err.Error())
	}

	if form.Command == "startup" {
		util.MyPrint("process service name :", service.Name)
		err = serviceSuperVisor.Cli.StartProcess(service.Name, true)
		// err = serviceSuperVisor.StopProcess(service.Name,true)
		util.ExitPrint(err)
		err = serviceSuperVisor.Cli.StartProcess(service.Name, true)
	} else if form.Command == "stop" {
		name := service.Name + ":service_" + service.Name
		util.MyPrint("name:", name)
		err = serviceSuperVisor.Cli.StopProcess(name, true)
		// err = serviceSuperVisor.StopProcess(service.Name,true)
		util.ExitPrint(err)
		// err = serviceSuperVisor.Cli.StopProcess(service.Name,true)
	} else if form.Command == "restart" {
	} else {
		return errors.New("command err")
	}

	return err
}

// 浏览器会请求此函数，拿到所有服务器和服务列表，用于部署
// 此方法较慢，因为要 ping 服务器。还要 连接 superVisor
func (cicdManager *CicdManager) LocalAllServerServiceList() (list LocalServerServiceList, err error) {
	list.ServiceList = make(map[int]model.Project)
	list.ServerList = make(map[int]util.Server)

	// util.MyPrint(cicdManager.Option.ServerList)

	if len(cicdManager.Option.ServerList) == 0 {
		// 服务器 为空
		errMsg := "GetSuperVisorList err:ServerList is empty"
		// util.MyPrint()
		return list, errors.New(errMsg)
	}

	if len(cicdManager.Option.ServiceList) == 0 {
		// 服务为空
		errMsg := "GetSuperVisorList err:ServiceList is empty"
		// util.MyPrint(errMsg)
		return list, errors.New(errMsg)
	}
	// list := LocalServerServiceList{}
	instanceManager := cicdManager.Option.InstanceManager

	for _, server := range cicdManager.Option.ServerList {
		// 先ping 一下，确定该服务器网络正常
		argsmap := map[string]interface{}{}
		p := util.NewPingOption()
		err = p.Ping3(server.OutIp, argsmap)
		// util.MyPrint("Ping3 rs:",err)
		if err != nil {
			server.PingStatus = util.SERVER_PING_FAIL
			list.ServerList[server.Id] = server
			continue
		} else {
			server.PingStatus = util.SERVER_PING_OK
		}

		instance, empty := instanceManager.GetByEnvName(server.Env, "super_visor")
		if empty {
			return list, errors.New("not found super_visor instance")
		}
		// 再测试下无端的superVisor是否正常
		// superVisorStatus := make(map[int]int)
		// 创建实例
		superVisorOption := util.SuperVisorOption{
			Ip: server.OutIp,
			// RpcPort:          cicdManager.Option.Config.SuperVisor.RpcPort
			RpcPort:  instance.Port,
			Username: instance.User,
			Password: instance.Ps,
		}

		serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
		if err != nil {
			util.MyPrint("NewSuperVisor err:", err)
			// superVisorStatus[server.Id ] = util.SV_ERROR_INIT
			server.SuperVisorStatus = util.SERVER_PING_FAIL
			list.ServerList[server.Id] = server
			continue
		}
		if server.Env == 1 {
			server.SuperVisorStatus = util.SERVER_PING_FAIL
			list.ServerList[server.Id] = server
			continue
		} else {
			// 建立 XMLRpc
			err = serviceSuperVisor.InitXMLRpc()
			if err != nil {
				util.MyPrint("serviceSuperVisor InitXMLRpc err:", err)
				server.SuperVisorStatus = util.SERVER_PING_FAIL
				list.ServerList[server.Id] = server
				continue
			}
		}

		server.SuperVisorStatus = util.SERVER_PING_OK
		list.ServerList[server.Id] = server

		for _, service := range cicdManager.Option.ProjectList {
			list.ServiceList[service.Id] = service
		}
	}

	return list, nil
}

// 每台服务器上 都会启动一个 superVisor 进程
// 列出每台机器上的：superVisor 进程 的所有服务进程的状态信息
func (cicdManager *CicdManager) GetSuperVisorList() (list ServerServiceSuperVisorList, err error) {
	if len(cicdManager.Option.ServerList) == 0 { // 服务器 为空
		errMsg := "GetSuperVisorList err:ServerList is empty"
		return list, errors.New(errMsg)
	}

	if len(cicdManager.Option.ServiceList) == 0 { // 服务为空
		errMsg := "GetSuperVisorList err:ServiceList is empty"
		return list, errors.New(errMsg)
	}

	util.MyPrint("GetSuperVisorList serverList len:", len(cicdManager.Option.ServerList), " ServiceList len:", len(cicdManager.Option.ServiceList))
	// 服务器 上面:已经开启的 superVisor	map[serverId]=>superVisorList
	serverServiceSuperVisor := make(map[int][]MyProcessInfo)
	// 服务器 状态
	// serverStatus := make(map[int]int)
	superVisorStatus := make(map[int]int)
	// 遍历服务器列表
	for _, server := range cicdManager.Option.ServerList {
		fmt.Println("for each service , outIp:" + server.OutIp + " env:" + strconv.Itoa(server.Env))
		instance, empty := cicdManager.Option.InstanceManager.GetByEnvName(server.Env, "super_visor")
		if empty {
			return list, errors.New("not found super_visor instance")
		}

		// 创建实例
		// superVisorOption := util.SuperVisorOption{
		// 	Ip:      server.OutIp,
		// 	RpcPort: cicdManager.Option.Config.SuperVisor.RpcPort,
		// 	Username:
		// 	Username: "ckadmin",
		// 	Password: "ckckarar",
		// }
		superVisorOption := util.SuperVisorOption{
			Ip:       server.OutIp,
			RpcPort:  instance.Port,
			Username: instance.User,
			Password: instance.Ps,
		}
		util.MyPrint("=====", superVisorOption)
		serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
		if err != nil {
			util.MyPrint("NewSuperVisor err:", err)
			superVisorStatus[server.Id] = util.SV_ERROR_INIT
			continue
		}
		// 建立 XMLRpc
		err = serviceSuperVisor.InitXMLRpc()
		if err != nil {
			util.MyPrint("serviceSuperVisor InitXMLRpc err:", err)
			superVisorStatus[server.Id] = util.SV_ERROR_CONN
			continue
		}
		// 获取当前机器上的superVisor 的所有 服务进程 状态
		processInfoList, err := serviceSuperVisor.Cli.GetAllProcessInfo()
		if err != nil {
			superVisorStatus[server.Id] = util.SV_ERROR_CONN
			util.MyPrint("SuperVisor GetAllProcessInfo err:" + err.Error())
			continue
		} else {
			// jsonStr,_  := json.Marshal(processInfoList)
			// util.MyPrint(jsonStr)
			// util.ExitPrint(string(jsonStr))
		}
		superVisorStatus[server.Id] = util.SV_ERROR_NONE

		for _, service := range cicdManager.Option.ProjectList {
			// 获取当前服务器，已部署的服务目录
			// servicePath := serviceBaseDir + util.DIR_SEPARATOR + service.Name
			// util.MyPrint("servicePath:", servicePath)

			defaultProcessInfo := supervisord.ProcessInfo{
				Name:  util.SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX + service.Name,
				State: 999, // 项目未部署过
			}

			superVisorProcessInfo := MyProcessInfo{
				ServiceId:   service.Id,
				ProcessInfo: defaultProcessInfo,
			}
			// search := 0
			// var superVisorProcessInfo supervisord.ProcessInfo
			// 先筛选一下，看看该服务有没有被 添加到superVisor中
			for _, process := range processInfoList {
				// rs := fmt.Sprintf("%+v",process)
				// util.MyPrint(rs)
				if process.Name == util.SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX+service.Name {
					util.MyPrint(fmt.Sprintf("%+v", process))
					superVisorProcessInfo.ProcessInfo = process
					// superVisorStatus[service.Id ] = util.SV_ERROR_NONE
					// search = 1

					serviceDeployConfig := cicdManager.Deploy.GetDeployConfig(DEPLOY_TARGET_TYPE_LOCAL)
					serviceDeployConfig, _ = cicdManager.Deploy.DeployServiceCheck(serviceDeployConfig, service, server)
					// path := serviceDeployConfig.MasterPath + "/" + serviceDeployConfig.OpDirName
					// masterSrc,_ := ExecShellFile2(path + "/" + "get_soft_link_src.sh",serviceDeployConfig.MasterPath)
					masterSrcPath, _ := filepath.EvalSymlinks(serviceDeployConfig.MasterPath)
					masterSrcPathArr := strings.Split(masterSrcPath, "/")
					superVisorProcessInfo.MasterSrc = masterSrcPathArr[len(masterSrcPathArr)-1]
					break
				}
			}

			_, ok := serverServiceSuperVisor[server.Id]
			if !ok {
				serverServiceSuperVisor[server.Id] = []MyProcessInfo{superVisorProcessInfo}
			} else {
				serverServiceSuperVisor[server.Id] = append(serverServiceSuperVisor[server.Id], superVisorProcessInfo)
			}
		}

	}

	list = ServerServiceSuperVisorList{
		// ServerPingStatus:        serverStatus,
		SuperVisorStatus:        superVisorStatus,
		ServerServiceSuperVisor: serverServiceSuperVisor,
	}
	return list, nil
}

// 如果一个路径不存在
func pathNotExistCreate(path string) error {
	_, err := util.PathExists(path)
	if err == nil { // 目录存在
		return nil
	}
	if os.IsNotExist(err) { // 目录不存在
		// 创建一个目录
		err = os.Mkdir(path, 0777)
		fmt.Println("create path:", path)
		if err != nil {
			fmt.Println("create path failed , err:", err)
		}
		// return err
	} else { // 其它错误
		//	fmt.Println("path :" + path + " exist , no need create.")
		//	return err
	}
	return err
}

// 执行shell文件
func ExecShellFile(shellFile string, argc string) (string, error) {
	util.MyPrint("ExecShellFile:", shellFile, " ", argc)
	shellCommand := shellFile + " " + argc
	c := exec.Command("sh", "-c", shellCommand)

	output, err := c.CombinedOutput()
	// util.MyPrint(string(output),err)
	if err != nil {
		util.MyPrint("exec.Command err:", err)
		return "", err
	}
	outStr := string(output)
	outArr := strings.Split(outStr, "\n")

	return outArr[1], nil
}

// 执行shell文件
func ExecShellFile2(shellFile string, argc string) (string, error) {
	util.MyPrint("ExecShellFile:", shellFile, " ", argc)
	shellCommand := shellFile + " " + argc
	c := exec.Command("sh", "-c", shellCommand)

	output, err := c.CombinedOutput()
	// util.MyPrint(string(output),err)
	if err != nil {
		util.MyPrint("exec.Command err:", err)
		return "", err
	}
	outStr := string(output)
	// outArr := strings.Split(outStr, "\n")

	return outStr, nil
}

// 执行shell 指令
func ExecShellCommand(command string, argc string) (string, error) {
	// util.MyPrint("ExecShellCommand:", command, argc)
	// shellCommand := command + " " + argc
	c := exec.Command("bash", "-c", command)

	output, err := c.CombinedOutput()
	strOutput := string(output)
	if err != nil {
		util.MyPrint("ExecShellCommand : <"+command+"> ,  has error , output:", strOutput, err.Error())
	} else {
		util.MyPrint("ExecShellCommand : <"+command+"> ,  success , output:", strOutput)
	}
	// if err != nil {
	//	fmt.Println("exec.Command err:", err)
	//	return "", err
	// }
	// outStr := string(output)
	// outArr := strings.Split(outStr,"\n")
	// return outArr[1],nil
	return strOutput, err
}

// 这里有一条简单的操作，80端口基本上都得用，测试服务器状态，用ping curl 也可以.
func (cicdManager *CicdManager) TestServerStateHttp(hostUri string) int {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	// 提交请求
	reqest, err := http.NewRequest("GET", hostUri, nil)
	if err != nil {
		return 0
	}
	util.MyPrint("TestServerStateHttp uri:" + hostUri)
	// 增加header选项
	reqest.Header.Add("X-Source-Type", "11")
	reqest.Header.Add("X-Project-Id", "6")
	reqest.Header.Add("X-Access", "imzgoframe")
	// 处理返回结果
	response, err := client.Do(reqest)
	// defer response.Body.Close()
	if err != nil {
		return 0
	}

	util.MyPrint("TestServerStateHttp http status:", response.Status)

	return 1

}

func (cicdManager *CicdManager) CheckInTestServer(ip string) bool {
	for _, v := range cicdManager.Option.TestServerList {
		if ip == v {
			return true
		}
	}
	return false
}

//
// //发布|部署 一次 服务
// func (cicdManager *CicdManager) Ping(c *gin.Context) {
//
//	//str,_ := json.Marshal(cicdManager.Option.ServiceList)
//	//c.String(200,string(str))
// }

// 获取所有3方服务列表
func (cicdManager *CicdManager) GetInstanceList(c *gin.Context) {
	// for _, instance := range cicdManager.Option.InstanceManager.Pool {
	//	status := util.CheckIpPort(instance.Host, instance.Port, 2)
	//	if !status {
	//		instance.Status = 3
	//	}
	// }
	//
	// str, _ := json.Marshal(cicdManager.Option.InstanceManager.Pool)
	// c.String(200, string(str))
}

// 在当前服务器上，从<部署目录>中检索出每个服务（目录名），分析出：哪些服务~已经部署
func (cicdManager *CicdManager) GetServiceList() map[int]model.Project {
	list := make(map[int]model.Project)

	for k, service := range cicdManager.Option.ProjectList {
		// server := cicdManager.Option.ServerList[service.Id]
		// localServiceDeployConfig := cicdManager.Deploy.GetDeployConfig(DEPLOY_TARGET_TYPE_LOCAL)
		// localServiceDeployConfig, _ = cicdManager.Deploy.DeployServiceCheck(localServiceDeployConfig, service, server)

		// dirList := util.ForeachDir(localServiceDeployConfig.BaseDir)
		// s := service
		// for _, dirInfo := range dirList {
		//	if s.Name == dirInfo.Name {
		//		s.Deploy = 1
		//		break
		//	}
		// }
		// list[k] = s
		list[k] = service
	}
	return list
	// str, _ := json.Marshal(cicdManager.Option.ServiceList)
}

// 获取所有服务器列表，并做ping，确定状态
func (cicdManager *CicdManager) GetServerList() map[int]util.Server {
	list := make(map[int]util.Server)
	for k, server := range cicdManager.Option.ServerList {
		// 这里是测试代码，不然PING太慢
		// if !cicdManager.CheckInTestServer(server.OutIp){
		//	server.PingStatus = 2
		// }else{
		arg_smap := map[string]interface{}{}
		p := util.NewPingOption()
		// host := "111.1.34.56"
		err := p.Ping3(server.OutIp, arg_smap)
		if err != nil {
			server.PingStatus = util.SERVER_PING_FAIL
		} else {
			server.PingStatus = util.SERVER_PING_OK
		}
		// }
		list[k] = server
	}

	return list

}

// func (cicdManager *CicdManager) GetHasDeployService() map[int]map[int][]string {
//	list := make(map[int]map[int][]string)
//	for _, server := range cicdManager.Option.ServerList {
//		serverDirList := make(map[int][]string)
//		for _, service := range cicdManager.Option.ServiceList {
//			form := request.CicdDeploy{
//				ServiceId: service.Id,
//				ServerId:  server.Id,
//				Flag:      DEPLOY_TARGET_TYPE_REMOTE,
//			}
//			dirList, _ := cicdManager.GetHasDeployServiceDirList(form)
//			serverDirList[service.Id] = dirList
//		}
//
//		list[server.Id] = serverDirList
//	}
//	//util.ExitPrint(list)
//	return list
// }

// //获取当前服务器上的，已部署过的，服务的，目录列表
// func (cicdManager *CicdManager) GetHasDeployServiceDirList(form request.CicdDeploy) ([]string, error) {
//	list := []string{}
//
//	server, service, err := cicdManager.Deploy.CheckCicdRequestForm(form)
//	if err != nil {
//		return nil, err
//	}
//	serviceDeployConfig := cicdManager.Deploy.GetDeployConfig(DEPLOY_TARGET_TYPE_REMOTE)
//	serviceDeployConfig, err = cicdManager.Deploy.DeployServiceCheck(serviceDeployConfig, service, server)
//	if err != nil {
//		return list, err
//	}
//
//	_, err = util.PathExists(serviceDeployConfig.FullPath)
//	if err != nil {
//		return list, err
//	}
//
//	dirList := util.ForeachDir(serviceDeployConfig.FullPath)
//
//	//util.MyPrint("lis len:",len(list) , " FullPath:",serviceDeployConfig.FullPath," list:",dirList)
//	for _, v := range dirList {
//		if v.Cate == "file" {
//			continue
//		}
//
//		if util.CheckServiceDeployDirName(v.Name) {
//			//util.MyPrint(111111,"===========")
//			list = append(list, v.Name)
//		}
//	}
//	//util.MyPrint(list)
//	return list, nil
//
// }

// func (cicdManager *CicdManager) LocalSyncTarget(form request.CicdSync) error {
//	sFrom := request.CicdDeploy{
//		ServerId:  form.ServerId,
//		ServiceId: form.ServiceId,
//		Flag:      DEPLOY_TARGET_TYPE_LOCAL,
//	}
//	server, service, err := cicdManager.Deploy.CheckCicdRequestForm(sFrom)
//	if err != nil {
//		return err
//	}
//
//	if form.VersionDir == "" {
//		return errors.New("VersionDir empty")
//	}
//
//	targetServiceDeployConfig := cicdManager.Deploy.GetDeployConfig(DEPLOY_TARGET_TYPE_LOCAL)
//	targetServiceDeployConfig, _ = cicdManager.Deploy.DeployServiceCheck(targetServiceDeployConfig, service, server)
//	targetDir := targetServiceDeployConfig.FullPath
//
//	localServiceDeployConfig := cicdManager.Deploy.GetDeployConfig(DEPLOY_TARGET_TYPE_REMOTE)
//	localServiceDeployConfig, _ = cicdManager.Deploy.DeployServiceCheck(localServiceDeployConfig, service, server)
//	localDir := localServiceDeployConfig.FullPath + "/" + form.VersionDir
//
//	//scp local_file remote_username@remote_ip:remote_folder
//	//shellArgc := " -r " + localDir + " root@"+server.OutIp + ":" + targetDir
//	//util.ExitPrint("scp "+ shellArgc)
//	//ExecShellCommand("scp",shellArgc)
//
//	shellArgc := "scp  -r " + localDir + " root@" + server.OutIp + ":" + targetDir
//	util.MyPrint(shellArgc)
//	ExecShellCommand(shellArgc, "")
//	return nil
// }
