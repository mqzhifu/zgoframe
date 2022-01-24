package util

import (
	"encoding/json"
	"fmt"
	"github.com/abrander/go-supervisord"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"zgoframe/model"
)
/*
自动化部署，从DB中读取出所有信息基础信息，GIT CLONE 配置super visor 监听进程
依赖
	supervisor 依赖 python 、 xmlrpc
	代码依赖：git
*/
const (
	DIR_SEPARATOR = "/"
	STR_SEPARATOR = "#"
)

//===============配置 结构体 开始===================
type ConfigCicdSystem struct {
	Env []string
	LogDir string
	ServiceDir string
}

type ConfigServiceCICDSystem struct{
	Startup	string
	ListeningPorts string
	TestUnit string
	Build string
	Command string
}

type ConfigServiceCICDDepend struct{
	Go string
	Node string
	Mysql string
	Redis string
}

type ConfigServiceCICD struct {
	System	ConfigServiceCICDSystem
	Depend	ConfigServiceCICDDepend
}

type ConfigCicdSuperVisor struct {
	RpcPort	string
	ConfTemplateFile string
	ConfDir string
}

type ConfigCicd struct {
	System ConfigCicdSystem
	SuperVisor ConfigCicdSuperVisor
}
//===============配置 结构体 结束===================



type SuperVisorReplace struct{
	script_name	string
	startup_script_command string
	script_work_dir string
	stdout_logfile string
	stderr_logfile string
	process_name string
}

//==============superVisor===========

type CicdPublish struct {
	Id				int
	RegTime 		int
	Status 			int
	ServiceName 	string
	Logs			[]string
	TotalExecTime 	int
	Server 			Server
}


type CicdManager struct {
	Option CicdManagerOption
}

type CicdManagerOption struct{
	ServerList 		map[int]Server	//所有服务器
	ServiceList 	map[int]Service	//所有项目/服务

	InstanceManager *InstanceManager
	Config 			ConfigCicd
	PublicManager 	*CICDPublicManager
}

func NewCicdManager(cicdManagerOption CicdManagerOption)*CicdManager{
	cicdManager := new(CicdManager)
	cicdManager.Option = cicdManagerOption

	return cicdManager
}
func(cicdManager *CicdManager)ReplaceInstance(content string,serviceName string ,env string)string{
	category := []string{"mysql","redis","etcd","rabbitmq","kafka","log","alert"}
	//attr := []string{"ip","port","user","ps"}
	separator := STR_SEPARATOR
	content = strings.Replace(content,separator + "env" + separator,env,-1)
	content = strings.Replace(content,separator + "log_dir" + separator,cicdManager.Option.Config.System.LogDir,-1)
	for _,v := range category{
		//for _,attrOne := range attr{
			instance,empty :=  cicdManager.Option.InstanceManager.GetByEnvName(env,v)
			if empty{
				MyPrint("cicdManager.Option.InstanceManager.GetByEnvName is empty,",env,v)
				continue
			}
			key := separator+ v  +"_" + "ip"  +separator
			content = strings.Replace(content,key,instance.Host,-1)

		key = separator  + v  +"_" + "port"  +separator
			content = strings.Replace(content,key,instance.Port,-1)

		key = separator  + v  +"_" + "user"  +separator
			content = strings.Replace(content,key,instance.User,-1)

		key = separator  + v  +"_" + "ps"  +separator
			content = strings.Replace(content,key,instance.Ps,-1)

		//}
	}

	return content
}

type ServiceDeployConfig struct {
	Name 				string 	//服务名称
	BaseDir 			string	//所有service项目统一放在一个目录下，由host.toml 中配置
	FullPath 			string 	//最终一个服务的目录名
	MasterDirName 		string	//一个服务的线上使用版本-软连目录名称
	CICDConfFileName 	string	//一个服务自己的需要执行的cicd脚本
	ConfigTmpFileName 	string	//一个服务的配置文件的模板文件名
	ConfigFileName 		string	//一个服务的配置文件名,由上面CP
	GitCloneTmpDirName 	string	//git clone 一个服务的项目代码时，临时存在的目录名
	CICDShellFileName 	string	//有一一些操作需要借用于shell 执行，如：git clone . 该变量就是shell 文件名
}
//一次部署全部服务项目
func(cicdManager *CicdManager)DeployAllService(){
	serviceDeployConfig := ServiceDeployConfig{
		BaseDir 			: cicdManager.Option.Config.System.ServiceDir,
		MasterDirName 		: "master",
		CICDConfFileName 	: "cicd.toml",
		ConfigTmpFileName	: "config.toml.tmp",
		ConfigFileName 		: "config.toml",
		GitCloneTmpDirName 	: "clone",
		CICDShellFileName 	: "./cicd.sh",
	}
	//先遍历所有服务器，然后，把所有已知服务部署到每台服务器上(每台机器都可以部署任何服务)
	for _,server :=range cicdManager.Option.ServerList{
		fmt.Println("for each service:" + server.OutIp + " " + server.Env)
		//遍历所有服务
		for _,service :=range cicdManager.Option.ServiceList{
			cicdManager.DeployOneService(server,serviceDeployConfig,service)
		}
	}
}
//部署一个服务
func(cicdManager *CicdManager)DeployOneService(server Server , serviceDeployConfig ServiceDeployConfig ,  service Service){
	serviceDeployConfig.Name = service.Name
	//创建发布记录
	publish := cicdManager.Option.PublicManager.InsertOne(service,server)
	MyPrint("create publish:",publish.Id)

	//一个服务的根目录，大部分操作都在这个目录下，除了superVisor
	servicePath := serviceDeployConfig.BaseDir + DIR_SEPARATOR +  service.Name
	serviceDeployConfig.FullPath = servicePath
	MyPrint("servicePath:",servicePath)
	pathNotExistCreate(servicePath)

	serviceMasterPath := servicePath + DIR_SEPARATOR + serviceDeployConfig.MasterDirName
	MyPrint("serviceMasterPath:"+serviceMasterPath)

	//git clone 目录
	serviceGitClonePath := servicePath + DIR_SEPARATOR + serviceDeployConfig.GitCloneTmpDirName
	pathNotExistCreate(serviceGitClonePath)
	//通过shell 执行git clone ，同时获取当前clone master 的版本号
	//gitLastCommitId :=GitCloneAndGetLastCommitIdByShell(serviceGitClonePath,service.Name,service.Git)
	//构建 shell 执行时所需 参数
	shellArgc := service.Git + " " + serviceGitClonePath + " " +  service.Name
	//执行shell 脚本 后：service项目代码已被clone, git 版本号已知了
	gitLastCommitId := ExecShellFile(serviceDeployConfig.CICDShellFileName,shellArgc)

	MyPrint("gitLastCommitId:",gitLastCommitId)
	//刚刚clone完后，项目的目录
	serviceCodeGitClonePath := serviceGitClonePath + DIR_SEPARATOR + service.Name
	//新刚刚克隆好的项目目录，移动一个新目录下，新目录名：git_master_versionId + 当前时间
	newGitCodeDir := servicePath + DIR_SEPARATOR + strconv.Itoa(GetNowTimeSecondToInt())  + "_" + gitLastCommitId
	MyPrint("service code move :",serviceCodeGitClonePath +" to "+ newGitCodeDir)
	//执行 移动操作
	os.Rename(serviceCodeGitClonePath,newGitCodeDir)

	//项目自带的CICD配置文件，这里有 服务启动脚本 和 依赖的环境
	serviceSelfCICDConf := newGitCodeDir + DIR_SEPARATOR + serviceDeployConfig.CICDConfFileName
	MyPrint("read file:"+serviceSelfCICDConf)
	serviceCICDConfig := ConfigServiceCICD{}
	//读取项目自己的cicd配置文件，并映射到结构体中
	ReadConfFile(serviceSelfCICDConf,&serviceCICDConfig)
	PrintStruct(serviceCICDConfig,":")

	//生成该服务的，superVisor 配置文件
	superVisorOption := SuperVisorOption{
		ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
		ServiceName		: service.Name,
		ConfDir			: cicdManager.Option.Config.SuperVisor.ConfDir,
	}
	serviceSuperVisor := NewSuperVisor(superVisorOption)
	//superVisor 配置文件中 动态的占位符，需要替换掉
	superVisorReplace := SuperVisorReplace{
		script_name				: service.Name,
		startup_script_command	: serviceCICDConfig.System.Startup,
		script_work_dir 		: serviceMasterPath,
		stdout_logfile 			: serviceDeployConfig.BaseDir + DIR_SEPARATOR + "super_visor_stdout.log",
		stderr_logfile 			: serviceDeployConfig.BaseDir + DIR_SEPARATOR + "super_visor_stderr.log",
		process_name : service.Name,
	}
	//替换配置文件中的动态值，并生成配置文件
	serviceConfFileContent := serviceSuperVisor.ReplaceConfTemplate(superVisorReplace)
	//将已替换好的文件，生成一个新的配置文件
	serviceSuperVisor.CreateServiceConfFile(serviceConfFileContent)
	//读取该服务自己的配置文件
	serviceSelfConfigTmpFileDir := newGitCodeDir + DIR_SEPARATOR + serviceDeployConfig.GitCloneTmpDirName
	MyPrint("read file:"+serviceSelfConfigTmpFileDir)
	//读取模板文件内容
	serviceSelfConfigTmpFileContent,err := ReadString(serviceSelfConfigTmpFileDir)
	if err != nil{
		ExitPrint("read file err ,"+err.Error())
	}
	//开始替换 服务自己配置文件中的，实例信息，如：IP PORT
	serviceSelfConfigTmpFileContentNew := cicdManager.ReplaceInstance(serviceSelfConfigTmpFileContent,service.Name,server.Env)
	//生成新的配置文件
	newConfig := newGitCodeDir + DIR_SEPARATOR + serviceDeployConfig.ConfigFileName
	newConfigFile ,_:= os.Create(newConfig)
	newConfigFile.Write([]byte(serviceSelfConfigTmpFileContentNew))

	//先执行 服务自带的 shell 预处理
	//if serviceCICDConfig.System.Command != ""{
	//	ExecShellCommand(serviceCICDConfig.System.Command,"")
	//}
	//
	//if serviceCICDConfig.System.Build != ""{
	//	ExecShellCommand(serviceCICDConfig.System.Build,"")
	//}
	//
	//if serviceCICDConfig.System.TestUnit != ""{
	//	ExecShellCommand(serviceCICDConfig.System.TestUnit,"")
	//}

	//将master软链 指向 上面刚刚clone下的最新代码上
	MyPrint("os.Symlink:",newGitCodeDir , " to ",serviceMasterPath)
	pathExist ,_ := PathExists(serviceMasterPath)
	if pathExist{
		MyPrint("master path exist , so need del .",serviceMasterPath)
		err = os.Remove(serviceMasterPath)
		if err != nil{
			MyPrint("os.Remove ",serviceMasterPath, " err:",err)
		}
	}

	err = os.Symlink(newGitCodeDir, serviceMasterPath)
	if err != nil{
		ExitPrint("os.Symlink err :",err)
	}
	cicdManager.Option.PublicManager.UpStatus(publish,2)
	return
}


func (cicdManager *CicdManager)StartHttp(staticDir string){
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

	ginRouter.StaticFS("/static",http.Dir(staticDir))


	ginRouter.GET("/getServerList", cicdManager.GetServerList)
	ginRouter.GET("/getInstanceList", cicdManager.GetInstanceList)
	ginRouter.GET("/getPublishList", cicdManager.GetPublishList)
	ginRouter.GET("/getSuperVisorList", cicdManager.GetSuperVisorList)
	ginRouter.GET("/getServiceList", cicdManager.GetServiceList)

	ginRouter.Run("127.0.0.1:1111")

	//404
	//ginRouter.NoMethod(HandleNotFound)
}

func (cicdManager *CicdManager)GetServiceList(c *gin.Context){

	str,_ := json.Marshal(cicdManager.Option.ServiceList)
	c.String(200,string(str))

}

func (cicdManager *CicdManager)GetServerList(c *gin.Context){

	str,_ := json.Marshal(cicdManager.Option.ServerList)
	c.String(200,string(str))

}

func (cicdManager *CicdManager)GetInstanceList(c *gin.Context){

	str,_ := json.Marshal(cicdManager.Option.InstanceManager.Pool)
	c.String(200,string(str))

}

func (cicdManager *CicdManager)GetPublishList(c *gin.Context){
	listArr,_ := cicdManager.Option.PublicManager.GetList()
	listMap := make(map[int]model.CICDPublish)
	for _,v:= range listArr{
		listMap[v.Id] = v
	}
	str,_ := json.Marshal(listArr)
	c.String(200,string(str))

}
//这里有一条简单的操作，80端口基本上都得用，测试服务器状态，用ping curl 也可以.
func (cicdManager *CicdManager)TestServerStateHttp(hostUri string)int{
	httpClient := http.Client{
		Timeout: time.Second * 1,
	}

	resp,err := httpClient.Get(hostUri)
	if err != nil{
		MyPrint("http get err:",err)
		return 0
	}

	MyPrint("http get status:",resp.Status)

	return 1

}

func (cicdManager *CicdManager)GetSuperVisorList(c *gin.Context){
	serviceBaseDir := cicdManager.Option.Config.System.ServiceDir

	if len(cicdManager.Option.ServerList) == 0{
		MyPrint("GetSuperVisorList err:ServerList is empty")
		return
	}

	if len(cicdManager.Option.ServiceList) == 0{
		MyPrint("GetSuperVisorList err:ServiceList is empty")
		return
	}

	MyPrint("serverList len:",len(cicdManager.Option.ServerList), " ServiceList len:",len(cicdManager.Option.ServiceList))

	//serverId=>superVisorList
	serverServiceSuperVisor := make(map[int][]supervisord.ProcessInfo)
	serverStatus := make(map[int]int)
	for _,server :=range cicdManager.Option.ServerList{
		fmt.Println("for each servce:" + server.OutIp + " " + server.Env)

		testServerRs := cicdManager.TestServerStateHttp("http://" + server.OutIp)
		if testServerRs == 0{
			MyPrint("")
			serverStatus[server.Id] = 3
			continue
		}
		//c.String(200,"ok")
		//return
		//生成该服务的，superVisor 配置文件
		superVisorOption := SuperVisorOption{
			Ip:	server.OutIp,
			Port: cicdManager.Option.Config.SuperVisor.RpcPort,
			ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
			ServiceName : "" ,
			ConfDir: cicdManager.Option.Config.SuperVisor.ConfDir,
		}
		serviceSuperVisor := NewSuperVisor(superVisorOption)
		err := serviceSuperVisor.InitXMLRpc()
		if err != nil {
			MyPrint("serviceSuperVisor InitXMLRpc err:",err)
			serverStatus[server.Id] = 4
			continue
		}

		serverStatus[server.Id] = server.Status

		processInfoList,_ := serviceSuperVisor.Cli.GetAllProcessInfo()
		for _,service :=range cicdManager.Option.ServiceList{
			servicePath := serviceBaseDir + DIR_SEPARATOR +  service.Name
			MyPrint("servicePath:",servicePath)

			superVisorProcessInfo := supervisord.ProcessInfo{
				Name: SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX + service.Name,
				State: 999,//项目未部署过
			}

			for _,process :=range processInfoList{
				if process.Name == SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX + service.Name{
					superVisorProcessInfo = process
					break
				}
			}

			_ ,ok := serverServiceSuperVisor[server.Id]
			if !ok {
				MyPrint(22222)
				serverServiceSuperVisor[server.Id] = []supervisord.ProcessInfo{superVisorProcessInfo}
			}else{
				MyPrint(3333)
				serverServiceSuperVisor[server.Id] = append(serverServiceSuperVisor[server.Id],superVisorProcessInfo)
			}
		}

	}

	type response struct{
		ServerStatus 			map[int]int							`json:"server_status"`
		ServerServiceSuperVisor	map[int][]supervisord.ProcessInfo	`json:"server_service_super_visor"`
	}

	myresponse := response{
		ServerStatus : serverStatus,
		ServerServiceSuperVisor: serverServiceSuperVisor,
	}
	str ,err := json.Marshal(myresponse)
	MyPrint("json err:",err)
	c.String(200 , string(str) )
}

func pathNotExistCreate(path string){
	pathExist ,_ := PathExists(path)
	//fmt.Print(err)
	if !pathExist {
		//创建一个目录
		err := os.Mkdir(path, 0777)
		fmt.Println("create path:",path)
		if err != nil {
			fmt.Println("create path failed , err:",err)
		}
	}else{
		fmt.Println("path :" + path + " exist , no need create.")
	}
}

func ExecShellCommand(command string ,argc string)string{
	//shellCommand := command + " " + argc
	c := exec.Command(command, argc)

	output, err := c.CombinedOutput()
	if err != nil{
		fmt.Println("exec.Command err:",err)
		return ""
	}
	outStr := string(output)
	outArr := strings.Split(outStr,"\n")

	return outArr[1]
}

func ExecShellFile(shellFile string ,argc string)string{
	MyPrint("ExecShellFile:",shellFile , " ", argc)
	shellCommand := shellFile + " " + argc
	c := exec.Command("sh", "-c", shellCommand)

	output, err := c.CombinedOutput()
	if err != nil{
		fmt.Println("exec.Command err:",err)
	}
	outStr := string(output)
	outArr := strings.Split(outStr,"\n")

	return outArr[1]
}

//func GitCloneAndGetLastCommitIdByShell(serviceGitClonePath string,serviceName string,gitCloneUrl string)string{
//	argc := gitCloneUrl + " " + serviceGitClonePath + " " +  serviceName
//
//	shellFileName := "./cicd.sh" + " " + argc
//	println(shellFileName)
//	c := exec.Command("sh", "-c", shellFileName)
//
//	output, err := c.CombinedOutput()
//	if err != nil{
//		fmt.Println("exec.Command err:",err)
//	}
//	outStr := string(output)
//	outArr := strings.Split(outStr,"\n")
//
//	return outArr[1]
//	//fmt.Println(string(output), " err :",err)
//	//var shellCommands []string
//	//shellCommands = append(shellCommands,"./.sh")
//	//return shellCommands
//}
