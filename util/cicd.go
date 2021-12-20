package util

import (
	"fmt"
	"github.com/abrander/go-supervisord"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//已知：依赖
//supervisor 依赖 python 、 xmlrpc
//代码依赖：git

const (
	DIR_SEPARATOR = "/"
	STR_SEPARATOR
)

type CicdPublish struct {
	Id				int
	RegTime 		int
	Status 			int
	ServiceName 	string
	Logs			[]string
	TotalExecTime 	int
	Server 			HostServer
}

type CicdConfigSystem struct {
	Env []string
	LogDir string
	ServiceDir string
}

type ServiceCICDConfigSystem struct{
	Startup	string
	ListeningPorts string
	TestUnit string
	Build string
	Command string
}

type ServiceCICDConfigDepend struct{
	Go string
	Node string
	Mysql string
	Redis string
}

type ServiceCICDConfig struct {
	System	ServiceCICDConfigSystem
	Depend	ServiceCICDConfigDepend
}

type CicdConfigSuperVisor struct {
	RpcPort	string
	ConfTemplateFile string
	ConfDir string
}

type CicdConfig struct {
	System CicdConfigSystem
	SuperVisor CicdConfigSuperVisor
}

type SuperVisorReplace struct{
	script_name	string
	startup_script_command string
	script_work_dir string
	stdout_logfile string
	stderr_logfile string
	process_name string
}

type SuperVisor struct {
	Ip	string
	RpcPort string
	ConfTemplateFile string
	ConfTemplateFileContent string
	ServiceName string
	ConfDir string
	Separator string
	Cli *supervisord.Client
}

type CicdManager struct {
	Option CicdManagerOption
}

func NewSuperVisor(ip string,port string,ConfTemplateFile string,serviceName string,confDir string)*SuperVisor{
	superVisor := new(SuperVisor)
	superVisor.Ip = ip
	superVisor.RpcPort = port

	superVisor.ConfDir = confDir
	superVisor.ServiceName = serviceName
	superVisor.Separator = STR_SEPARATOR

	superVisorConfTemplateFileContent ,err := ReadString(ConfTemplateFile)
	if err != nil{
		ExitPrint("read superVisorConfTemplateFileContent err.")
	}

	superVisor.ConfTemplateFile = ConfTemplateFile
	superVisor.ConfTemplateFileContent = superVisorConfTemplateFileContent

	//ExitPrint(superVisorConfTemplateFileContent)

	return superVisor
}

func(superVisor *SuperVisor) Init(){
	dns := "http://" + superVisor.Ip + ":" + superVisor.RpcPort + "/RPC2"
	c, err := supervisord.NewClient(dns)
	if err != nil{

	}
	superVisor.Cli = c

}

func(superVisor *SuperVisor)CreateServiceConfFile(content string)error{
	fileName := superVisor.ConfDir +STR_SEPARATOR +  superVisor.ServiceName + ".ini"
	file ,err := os.Create(fileName)
	MyPrint("os.Create:" ,fileName)
	if err!= nil{
		MyPrint("os.Create :",fileName , " err:",err)
		return err
	}

	file.Write([]byte(content))
	return nil
}

func(superVisor *SuperVisor)ReplaceConfTemplate(replaceSource SuperVisorReplace)string{

	content := superVisor.ConfTemplateFileContent
	key := superVisor.Separator+"script_name"+superVisor.Separator
	MyPrint(key)
	content = strings.Replace(content,key,replaceSource.script_name,-1)

	key = superVisor.Separator+"startup_script_command"+superVisor.Separator
	content = strings.Replace(content,key,replaceSource.startup_script_command,-1)

	key = superVisor.Separator+"script_work_dir"+superVisor.Separator
	content = strings.Replace(content,key,replaceSource.script_work_dir,-1)

	key = superVisor.Separator+"stdout_logfile"+superVisor.Separator
	content = strings.Replace(content,key,replaceSource.stdout_logfile,-1)

	key = superVisor.Separator+"stderr_logfile"+superVisor.Separator
	content = strings.Replace(content,key,replaceSource.stderr_logfile,-1)

	key = superVisor.Separator+"process_name"+superVisor.Separator
	content = strings.Replace(content,key,replaceSource.process_name,-1)

	return content
}

type CicdManagerOption struct{
	ServerList map[int]HostServer
	AppList map[int]App
	ServiceList map[int]Service
	InstanceManager *InstanceManager
	Config CicdConfig
	PublicManager *CICDPublicManager
}

func NewCicdManager(cicdManagerOption CicdManagerOption)*CicdManager{
	cicdManager := new(CicdManager)
	cicdManager.Option = cicdManagerOption

	//cicdManager.Option.AppList = make(map[int]App)
	//cicdManager.ServerList = hostServerList
	//cicdManager.AppList[1] = App{
	//	Name: "zgoframe",
	//	Git: "git://github.com/mqzhifu/zgoframe.git",
	//}

	//cicdManager.ServerList["127.0.0.1"] = Server{
	//	Ip: "127.0.0.1",
	//}

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
			instance,empty :=  cicdManager.Option.InstanceManager.GetByEnvName(env,serviceName)
			if empty{
				ExitPrint("cicdManager.Option.InstanceManager.GetByEnvName is empty,",env,serviceName)
			}
			key := separator+v +"_" + "ip" + "_" + instance.Name  +separator
			content = strings.Replace(content,key,strconv.Itoa(instance.Id),-1)

			key = separator+v +"_" + "port" + "_" + instance.Name  +separator
			content = strings.Replace(content,key,instance.Port,-1)

			key = separator+v +"_" + "user" + "_" + instance.Name  +separator
			content = strings.Replace(content,key,instance.User,-1)

			key = separator+v +"_" + "ps" + "_" + instance.Name  +separator
			content = strings.Replace(content,key,instance.Ps,-1)

		//}
	}

	return content
}
func(cicdManager *CicdManager)Init(){
	serviceBaseDir := cicdManager.Option.Config.System.ServiceDir
	serviceMasterPathName := "master"
	serviceSelfCICDConfFile := "cicd.toml"
	serviceSelfConfigTmpFile := "config.toml.tmp"
	//fmt.Println("superVisorPort:",superVisorPort)
	//fmt.Println("serviceBaseDir:",serviceBaseDir)
	//fmt.Println("serviceMasterPathName:",serviceMasterPathName)




	for _,server :=range cicdManager.Option.ServerList{
		fmt.Println("for each servce:" + server.OutIp + " " + server.Env)
		for _,service :=range cicdManager.Option.ServiceList{
			//创建发布记录
			publish := cicdManager.Option.PublicManager.InsertOne(service,server)
			MyPrint("create publish:",publish.Id)

			//一个服务的根目录，大部分操作都在这个目录下，除了superVisor
			servicePath := serviceBaseDir + DIR_SEPARATOR +  service.Name
			MyPrint("servicePath:",servicePath)
			pathNotExistCreate(servicePath)

			serviceMasterPath := servicePath + DIR_SEPARATOR + serviceMasterPathName
			MyPrint("serviceMasterPath:"+serviceMasterPath)

			//git clone 目录
			serviceGitClonePath := servicePath + DIR_SEPARATOR + "clone"
			pathNotExistCreate(serviceGitClonePath)
			//通过shell 执行git clone ，同时获取当前clone master 的版本号
			//gitLastCommitId :=GitCloneAndGetLastCommitIdByShell(serviceGitClonePath,service.Name,service.Git)
			shellArgc := service.Git + " " + serviceGitClonePath + " " +  service.Name
			gitLastCommitId := ExecShellFile("./cicd.sh",shellArgc)

			MyPrint("gitLastCommitId:",gitLastCommitId)
			//刚刚clone完后，项目的目录
			serviceCodeGitClonePath := serviceGitClonePath + DIR_SEPARATOR + service.Name
			//新刚刚克隆好的项目目录，移动一个新目录下，新目录名：git_master_versionId + 当前时间
			newGitCodeDir := servicePath + "/" + strconv.Itoa(GetNowTimeSecondToInt())  + "_" + gitLastCommitId
			MyPrint("service code move :",serviceCodeGitClonePath +" to "+ newGitCodeDir)
			//执行 移动操作
			os.Rename(serviceCodeGitClonePath,newGitCodeDir)

			//项目自带的CICD配置文件，这里有 服务启动脚本 和 依赖的环境
			serviceSelfCICDConf := newGitCodeDir + DIR_SEPARATOR + serviceSelfCICDConfFile
			MyPrint("read file:"+serviceSelfCICDConf)
			serviceCICDConfig := ServiceCICDConfig{}
			ReadConfFile(serviceSelfCICDConf,&serviceCICDConfig)
			PrintStruct(serviceCICDConfig,":")

			//生成该服务的，superVisor 配置文件
			serviceSuperVisor := NewSuperVisor(server.OutIp,cicdManager.Option.Config.SuperVisor.RpcPort, cicdManager.Option.Config.SuperVisor.ConfTemplateFile,service.Name,cicdManager.Option.Config.SuperVisor.ConfDir)
			superVisorReplace := SuperVisorReplace{
				script_name: service.Name,
				startup_script_command:serviceCICDConfig.System.Startup,
				script_work_dir :serviceMasterPath,
				stdout_logfile :serviceBaseDir + DIR_SEPARATOR + "super_visor_stdout.log",
				stderr_logfile :serviceBaseDir + DIR_SEPARATOR + "super_visor_stderr.log",
				process_name :"service_"+service.Name,
			}
			//替换配置文件中的动态值，并生成配置文件
			serviceConfFileContent := serviceSuperVisor.ReplaceConfTemplate(superVisorReplace)
			serviceSuperVisor.CreateServiceConfFile(serviceConfFileContent)

			//读取该服务自己的配置文件
			serviceSelfConfigTmpFileDir := newGitCodeDir + DIR_SEPARATOR + serviceSelfConfigTmpFile
			MyPrint("read file:"+serviceSelfConfigTmpFileDir)
			serviceSelfConfigTmpFileContent,err := ReadString(serviceSelfConfigTmpFileDir)
			if err != nil{
				ExitPrint("read file err ,"+err.Error())
			}

			//开始替换 服务自己配置文件中的，实例信息，如：IP PORT
			serviceSelfConfigTmpFileContentNew := cicdManager.ReplaceInstance(serviceSelfConfigTmpFileContent,service.Name,server.Env)

			ExitPrint(serviceSelfConfigTmpFileContentNew)
			//先执行 服务自带的 shell 预处理
			if serviceCICDConfig.System.Command != ""{
				ExecShellCommand(serviceCICDConfig.System.Command,"")
			}

			if serviceCICDConfig.System.Build != ""{
				ExecShellCommand(serviceCICDConfig.System.Build,"")
			}

			if serviceCICDConfig.System.TestUnit != ""{
				ExecShellCommand(serviceCICDConfig.System.TestUnit,"")
			}

			//将master软链 指向 上面刚刚clone下的最新代码上
			MyPrint("os.Symlink:",newGitCodeDir , " to ",serviceMasterPath)
			err = os.Symlink(newGitCodeDir, serviceMasterPath)
			if err != nil{
				ExitPrint("os.Symlink err :",err)
			}
			return
		}
	}
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

func GitCloneAndGetLastCommitIdByShell(serviceGitClonePath string,serviceName string,gitCloneUrl string)string{
	argc := gitCloneUrl + " " + serviceGitClonePath + " " +  serviceName

	shellFileName := "./cicd.sh" + " " + argc
	println(shellFileName)
	c := exec.Command("sh", "-c", shellFileName)

	output, err := c.CombinedOutput()
	if err != nil{
		fmt.Println("exec.Command err:",err)
	}
	outStr := string(output)
	outArr := strings.Split(outStr,"\n")

	return outArr[1]
	//fmt.Println(string(output), " err :",err)
	//var shellCommands []string
	//shellCommands = append(shellCommands,"./.sh")
	//return shellCommands
}
