package cicd

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

/*
自动化部署，从DB中读取出所有信息基础信息，GIT CLONE 配置super visor 监听进程
依赖
	supervisor 依赖 python 、 xmlrpc
	代码依赖：git
*/

const (
	DEPLOY_TARGET_TYPE_LOCAL = 1			//本要部署
	DEPLOY_TARGET_TYPE_REMOTE = 2			//远程部署，并同步到了本机

	DEPLOY_TARGET_TYPE_LOCAL_NAME ="local"
	DEPLOY_TARGET_TYPE_REMOTE_NAME = "remote"
)


func GetConstListCicdDeployTargetType() map[string]int {
	list := make(map[string]int)

	list["本要部署再同步到远端"] = DEPLOY_TARGET_TYPE_LOCAL
	list["直接在远端部署"] = DEPLOY_TARGET_TYPE_REMOTE

	return list
}

type ServiceDeployConfig struct {
	Name               string //服务名称
	BaseDir            string //所有service项目统一放在一个目录下，由host.toml 中配置
	FullPath           string //最终一个服务的目录名,BaseDir + serviceName
	MasterDirName      string //一个服务的线上使用版本-软连目录名称
	MasterPath         string //full path + MasterDirName
	CICDConfFileName   string //一个服务自己的需要执行的cicd脚本
	OpDirName          string //存放所有：运维工具脚本的目录
	//ConfigTmpFileName  string //一个服务的配置文件的模板文件名
	//ConfigFileName     string //一个服务的配置文件名,由上面CP
	GitCloneTmpDirName string //git clone 一个服务的项目代码时，临时存在的目录名
	ClonePath          string //service dir + GitCloneTmpDirName
	CodeGitClonePath   string // ClonePath + service name
	CICDShellFileName  string //有一一些操作需要借用于shell 执行，如：git clone . 该变量就是shell 文件名
	DeployTargetType	int	//1本地部署2远端部署
}



func (cicdManager *CicdManager) ApiDeployOneService(form request.CicdDeploy)(error){
	server , service ,err := cicdManager.CheckCicdRequestForm(form)
	if err != nil{
		return err
	}
	serviceDeployConfig := cicdManager.GetDeployConfig(form.Flag)
	err = cicdManager.DeployOneService(server, serviceDeployConfig, service)
	return err
}

func (cicdManager *CicdManager)CheckCicdRequestForm(form request.CicdDeploy)(server util.Server,service util.Service, err error){
	server , ok := cicdManager.Option.ServerList[form.ServerId]
	if !ok {
		return server,service,errors.New("serviceId not in list")
	}
	service , ok = cicdManager.Option.ServiceList[form.ServiceId]
	if !ok {
		return server,service,errors.New("serviceId not in list")
	}

	if form.Flag <= 0{
		return server,service, errors.New("Flag <= 0")
	}

	if form.Flag != DEPLOY_TARGET_TYPE_LOCAL && form.Flag != DEPLOY_TARGET_TYPE_REMOTE{
		return server,service, errors.New("Flag err")
	}

	return server,service,nil
}

//一次部署全部服务项目
func (cicdManager *CicdManager) DeployAllService(deployTargetType int) {

	serviceDeployConfig := cicdManager.GetDeployConfig(deployTargetType)
	//先遍历所有服务器，然后，把所有已知服务部署到每台服务器上(每台机器都可以部署任何服务)
	for _, server := range cicdManager.Option.ServerList {
		//遍历所有服务
		for _, service := range cicdManager.Option.ServiceList {
			err := cicdManager.DeployOneService(server, serviceDeployConfig, service)
			if err != nil {
				util.ExitPrint(err)
			}
		}
	}
}

func (cicdManager *CicdManager)GetDeployConfig(deployTargetType int)ServiceDeployConfig{
	serviceBaseDir := ""
	if deployTargetType == DEPLOY_TARGET_TYPE_REMOTE{
		serviceBaseDir = cicdManager.Option.Config.System.WorkBaseDir + "/" + DEPLOY_TARGET_TYPE_REMOTE_NAME + "/"
	}else if deployTargetType == DEPLOY_TARGET_TYPE_LOCAL{
		serviceBaseDir  = cicdManager.Option.Config.System.WorkBaseDir + "/" + DEPLOY_TARGET_TYPE_LOCAL_NAME + "/"
	}else{
		util.ExitPrint("deployTargetType err:",deployTargetType)
	}

	//cicdManager.Option.Log.Info("DeployAllService:")
	serviceDeployConfig := ServiceDeployConfig{
		DeployTargetType: deployTargetType,
		BaseDir:            serviceBaseDir,
		OpDirName:          cicdManager.Option.OpDirName,
		MasterDirName:     cicdManager.Option.Config.System.MasterDirName,
		GitCloneTmpDirName: cicdManager.Option.Config.System.GitCloneTmpDirName,

		CICDConfFileName:   "cicd.toml",//本项目的相关 脚本及配置
		CICDShellFileName:  "cicd.sh",//执行：git clone 代码，并获取当前git最新版本号
	}
	return serviceDeployConfig
}

func (cicdManager *CicdManager) DeployServiceCheck( serviceDeployConfig ServiceDeployConfig, service util.Service ,server util.Server) (ServiceDeployConfig, error) {
	if service.Git == "" {
		errMsg := "service.Git is empty~" + service.Name
		return serviceDeployConfig, errors.New(errMsg)
	}

	if service.Name == "" {
		errMsg := "service.Name is empty~"
		return serviceDeployConfig, errors.New(errMsg)
	}

	if serviceDeployConfig.MasterDirName == "" {
		errMsg := "MasterDirName is empty~"
		return serviceDeployConfig, errors.New(errMsg)
	}

	if serviceDeployConfig.GitCloneTmpDirName == "" {
		errMsg := "GitCloneTmpDirName is empty~"
		return serviceDeployConfig, errors.New(errMsg)
	}

	if serviceDeployConfig.OpDirName == "" {
		errMsg := "OpDirName is empty~"
		return serviceDeployConfig, errors.New(errMsg)
	}

	_, err := util.PathExists(serviceDeployConfig.BaseDir)
	if err != nil{
		if os.IsNotExist(err) {
			util.MyPrint("DeployServiceCheck create dir:",serviceDeployConfig.BaseDir)
			err = os.Mkdir(serviceDeployConfig.BaseDir,0777)
			if err != nil{
				util.ExitPrint("os.Mkdir :",serviceDeployConfig.BaseDir, " err:",err.Error())
			}
		}
	}
	//本机部分编译，要把远程部署多出一层： 服务器IP目录->服务目录
	if serviceDeployConfig.DeployTargetType == DEPLOY_TARGET_TYPE_REMOTE{
		newBaseDir := serviceDeployConfig.BaseDir + "/" + server.OutIp
		_, err := util.PathExists(newBaseDir)
		if err != nil{
			if os.IsNotExist(err) {
				util.MyPrint("DEPLOY_TARGET_TYPE_LOCAL create dir:",newBaseDir)
				err = os.Mkdir(newBaseDir,0777)
				if err != nil{
					util.ExitPrint("os.Mkdir :",newBaseDir, " err:",err.Error())
				}
			}else{
				util.ExitPrint("util.PathExists err:",err.Error())
			}
		}
		serviceDeployConfig.BaseDir += server.OutIp
	}
	//baseDir 已由 构造函数做校验了

	serviceDeployConfig.Name = service.Name
	serviceDeployConfig.FullPath = serviceDeployConfig.BaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name
	serviceDeployConfig.MasterPath = serviceDeployConfig.FullPath + util.DIR_SEPARATOR + serviceDeployConfig.MasterDirName
	serviceDeployConfig.ClonePath = serviceDeployConfig.FullPath + util.DIR_SEPARATOR + serviceDeployConfig.GitCloneTmpDirName
	serviceDeployConfig.CodeGitClonePath = serviceDeployConfig.ClonePath + util.DIR_SEPARATOR + service.Name

	newServiceDeployConfig := serviceDeployConfig


	//util.PrintStruct(newServiceDeployConfig, ":")

	return newServiceDeployConfig, nil
}

//部署一个服务
func (cicdManager *CicdManager) DeployOneService(server util.Server, serviceDeployConfig ServiceDeployConfig, service util.Service) error {
	startTime := util.GetNowTimeSecondToInt()
	if service.Name != "Zgoframe"  { //测试代码,只部署：local Zgoframe
		errMsg := "service name != Zgoframe"
		util.MyPrint(errMsg)
		return errors.New(errMsg)
	}

	if server.Env != 1 &&  server.Env != 4 { //测试代码,只部署：local Zgoframe
		errMsg := " server.Env != 1 "
			util.MyPrint(errMsg)
		return errors.New(errMsg)
	}

	cicdManager.Option.Log.Info("DeployOneService:" + server.OutIp + " " + strconv.Itoa(server.Env) + " " + service.Name)
	//创建发布记录
	publish := cicdManager.Option.PublicManager.InsertOne(service, server,serviceDeployConfig.DeployTargetType)
	cicdManager.Option.Log.Info("create publish:" + strconv.Itoa(publish.Id))
	//检查各种路径是否正确
	newServiceDeployConfig, err := cicdManager.DeployServiceCheck( serviceDeployConfig, service , server)
	//util.PrintStruct(newServiceDeployConfig,":")
	//util.ExitPrint(11)
	if err != nil {
		return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	cicdManager.Option.Log.Info("DeployServiceCheck ok~")

	serviceDeployConfig = newServiceDeployConfig
	//step 1 : 项目代码及目录(git)相关
	newGitCodeDir, projectDirName ,err := cicdManager.DeployOneServiceGitCode(serviceDeployConfig, service)

	if err != nil {
		return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	p := model.CicdPublish{}
	p.Id = publish.Id
	p.CodeDir = projectDirName
	cicdManager.Option.PublicManager.UpInfo(p)
	//util.ExitPrint(p)
	//step 2 : 读取service项目代码里自带的cicd.toml ,供:后面使用
	serviceCICDConfig, err := cicdManager.DeployOneServiceCICIConfig(newGitCodeDir, serviceDeployConfig, server)
	if err != nil {
		return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	//step 3: 生成该服务的，superVisor 配置文件
	err = cicdManager.DeployOneServiceSuperVisor(serviceDeployConfig, serviceCICDConfig)
	if err != nil {
		return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	//step 4: 处理项目自带的主配置文件
	err = cicdManager.DeployOneServiceProjectConfig(newGitCodeDir, server, serviceDeployConfig,serviceCICDConfig,service)
	if err != nil {
		return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	//step 5 : 先执行 服务自带的 shell 预处理
	_, err = cicdManager.DeployOneServiceCommand(newGitCodeDir, serviceDeployConfig, serviceCICDConfig)
	if err != nil {
		return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}

	cicdManager.Option.PublicManager.UpDeployStatus(publish, model.CICD_PUBLISH_DEPLOY_STATUS_FINISH)
	cicdManager.Option.PublicManager.UpStatus(publish,model.CICD_PUBLISH_STATUS_WAIT_PUB)

	endTime :=util.GetNowTimeSecondToInt()
	execTime := endTime - startTime
	e := model.CicdPublish{}
	e.Id = publish.Id
	e.ExecTime = execTime
	cicdManager.Option.PublicManager.UpInfo(e)

	return nil
}

func (cicdManager *CicdManager)Publish(id int,deployTargetType int)error{
	serviceDeployConfig := cicdManager.GetDeployConfig(deployTargetType)
	publishRecord , err := cicdManager.Option.PublicManager.GetById(id)
	if err !=nil{
		return err
	}

	server := cicdManager.Option.ServerList[publishRecord.ServiceId]

	service := cicdManager.Option.ServiceList[publishRecord.ServiceId]
	serviceDeployConfig ,_ = cicdManager.DeployServiceCheck(serviceDeployConfig,service,server)
	//将master软链 指向 上面刚刚clone下的最新代码上
	err = cicdManager.DeployOneServiceLinkMaster(publishRecord.CodeDir, serviceDeployConfig)
	if err != nil {
		cicdManager.Option.PublicManager.UpStatus(publishRecord, model.CICD_PUBLISH_DEPLOY_FAIL)
		return err
		//return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	cicdManager.Option.PublicManager.UpStatus(publishRecord, model.CICD_PUBLISH_DEPLOY_OK)
	return nil
}

//step 1
func (cicdManager *CicdManager) DeployOneServiceGitCode(serviceDeployConfig ServiceDeployConfig, service util.Service) (string, string, error) {
	cicdManager.Option.Log.Info("step 1 : git clone project code and get git commit id.")

	//FullPath 一个服务的根目录，大部分操作都在这个目录下(除了superVisor)
	//servicePath := serviceDeployConfig.BaseDir + DIR_SEPARATOR +  service.Name
	//serviceDeployConfig.FullPath = servicePath
	//cicdManager.Option.Log.Info("servicePath:" + serviceDeployConfig.FullPath)
	//查看服务根目录是否存在，不存在即新创建
	pathNotExistCreate(serviceDeployConfig.FullPath)
	//serviceMasterPath := servicePath + DIR_SEPARATOR + serviceDeployConfig.MasterDirName
	//cicdManager.Option.Log.Info("serviceMasterPath:"+serviceMasterPath)

	//git clone 目录
	//serviceGitClonePath := serviceDeployConfig.FullPath + DIR_SEPARATOR + serviceDeployConfig.GitCloneTmpDirName
	//查看git clone 目录是否存在，不存在即新创建
	pathNotExistCreate(serviceDeployConfig.ClonePath)
	//通过shell 执行git clone ，同时获取当前clone master 的版本号
	//gitLastCommitId :=GitCloneAndGetLastCommitIdByShell(serviceGitClonePath,service.Name,service.Git)
	//构建 shell 执行时所需 参数

	shellArgc := service.Git + " " + serviceDeployConfig.ClonePath + " " + service.Name
	//执行shell 脚本 后：service项目代码已被clone, git 版本号已知了

	pwd, _ := os.Getwd() //当前路径]
	opDirFull := pwd + "/" + cicdManager.Option.OpDirName

	gitLastCommitId, err := ExecShellFile(opDirFull+"/"+serviceDeployConfig.CICDShellFileName, shellArgc)
	if err != nil {
		return "","", errors.New("ExecShellFile err:" + err.Error())
	}
	//cicdManager.Option.Log.Info("gitLastCommitId:" + gitLastCommitId)
	//刚刚clone完后，项目的目录
	//serviceCodeGitClonePath := serviceDeployConfig.ClonePath + DIR_SEPARATOR + service.Name
	//新刚刚克隆好的项目目录，移动一个新目录下，新目录名：git_master_versionId + 当前时间
	projectDirName := strconv.Itoa(util.GetNowTimeSecondToInt()) + "_" + gitLastCommitId
	newGitCodeDir := serviceDeployConfig.FullPath + util.DIR_SEPARATOR + projectDirName
	cicdManager.Option.Log.Info(" service code move :" + serviceDeployConfig.CodeGitClonePath + " to " + newGitCodeDir)
	//执行 移动操作
	err = os.Rename(serviceDeployConfig.CodeGitClonePath, newGitCodeDir)
	if err != nil {
		return newGitCodeDir, "", errors.New("serviceCodeGitClonePath os.Rename err:" + err.Error())
	}
	cicdManager.Option.Log.Info("step 1 finish , newGitCodeDir :  " + newGitCodeDir + " , gitLastCommitId:" + gitLastCommitId)
	return newGitCodeDir , projectDirName , nil
}

//step 2
func (cicdManager *CicdManager) DeployOneServiceCICIConfig(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig, server util.Server) (ConfigServiceCICD, error) {
	cicdManager.Option.Log.Info("step 2:load service CICD config ")
	//项目自带的CICD配置文件，这里有 服务启动脚本 和 依赖的环境
	serviceSelfCICDConf := newGitCodeDir + util.DIR_SEPARATOR + serviceDeployConfig.OpDirName + util.DIR_SEPARATOR + serviceDeployConfig.CICDConfFileName
	cicdManager.Option.Log.Info("read file:" + serviceSelfCICDConf)
	serviceCICDConfig := ConfigServiceCICD{}
	//读取项目自己的cicd配置文件，并映射到结构体中
	err := util.ReadConfFile(serviceSelfCICDConf, &serviceCICDConfig)
	if err != nil {
		return serviceCICDConfig, errors.New(err.Error())
	}
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#service_name#", serviceDeployConfig.Name, -1)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#env#",strconv.Itoa( server.Env), -1)
	util.PrintStruct(serviceCICDConfig, ":")

	return serviceCICDConfig, nil
}

//step 3 生成该服务的，superVisor 配置文件
func (cicdManager *CicdManager) DeployOneServiceSuperVisor(serviceDeployConfig ServiceDeployConfig, configServiceCICD ConfigServiceCICD) error {
	cicdManager.Option.Log.Info("step 3 : create superVisor conf file.")
	superVisorOption := util.SuperVisorOption{
		ConfDir:          cicdManager.Option.Config.SuperVisor.ConfDir,
		ServiceName:      serviceDeployConfig.Name,
		//ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
	}

	serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
	if err != nil {
		return err
	}
	serviceSuperVisor.SetConfTemplateFile(cicdManager.Option.Config.SuperVisor.ConfTemplateFile)
	//superVisor 配置文件中 动态的占位符，需要替换掉
	superVisorReplace := util.SuperVisorReplace{
		ScriptName:            serviceDeployConfig.Name,
		StartupScriptCommand: configServiceCICD.System.Startup,
		ScriptWorkDir:        serviceDeployConfig.MasterPath,
		StdoutLogfile:         serviceDeployConfig.BaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name +  "/super_visor_stdout.log",
		StderrLogfile:         serviceDeployConfig.BaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name + "/super_visor_stderr.log",
		ProcessName:           serviceDeployConfig.Name,
	}
	//util.PrintStruct(superVisorReplace,":")
	//替换配置文件中的动态值，并生成配置文件
	serviceConfFileContent,_ := serviceSuperVisor.ReplaceConfTemplate(superVisorReplace)
	//util.ExitPrint(serviceConfFileContent)
	//将已替换好的文件，生成一个新的配置文件
	err = serviceSuperVisor.CreateServiceConfFile(serviceConfFileContent)
	if err != nil {
		return err
	}

	return nil
}

//step 4
func (cicdManager *CicdManager) DeployOneServiceProjectConfig(newGitCodeDir string, server util.Server, serviceDeployConfig ServiceDeployConfig,configServiceCICD ConfigServiceCICD ,service util.Service ) error {
	cicdManager.Option.Log.Info("step 4 : create project self conf file.")
	//读取该服务自己的配置文件 config.toml
	serviceSelfConfigTmpFileDir := newGitCodeDir + util.DIR_SEPARATOR + configServiceCICD.System.ConfigTmpFileName
	_, err := util.FileExist(serviceSelfConfigTmpFileDir)
	if err != nil {
		return errors.New("serviceSelfConfigTmpFileDir CheckFileIsExist err:" + err.Error())
	}
	cicdManager.Option.Log.Info("read file:" + serviceSelfConfigTmpFileDir)
	//读取模板文件内容
	serviceSelfConfigTmpFileContent, err := util.ReadString(serviceSelfConfigTmpFileDir)
	if err != nil {
		return errors.New(err.Error())
	}
	//开始替换 服务自己配置文件中的，实例信息，如：IP PORT
	serviceSelfConfigTmpFileContentNew := cicdManager.ReplaceInstance(serviceSelfConfigTmpFileContent, serviceDeployConfig.Name, server.Env)

	key := util.STR_SEPARATOR + "projectId" + util.STR_SEPARATOR
	serviceSelfConfigTmpFileContentNew = strings.Replace(serviceSelfConfigTmpFileContentNew, key, strconv.Itoa(service.Id) , -1)

	//生成新的配置文件
	newConfig := newGitCodeDir + util.DIR_SEPARATOR + configServiceCICD.System.ConfigFileName
	newConfigFile, _ := os.Create(newConfig)
	contentByte := bytes.Trim([]byte(serviceSelfConfigTmpFileContentNew),"\x00")//NUL
	newConfigFile.Write(contentByte)

	return nil
}

//step 5
func (cicdManager *CicdManager) DeployOneServiceCommand(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig, serviceCICDConfig ConfigServiceCICD) (output string, err error) {
	cicdManager.Option.Log.Info("step 5 : DeployOneServiceCommand.")
	//cicdManager.Option.Log.Info("step 6.1 : project pre command "+serviceCICDConfig.System.Command)
	//    /usr/local/Cellar/go/1.16.5/bin/
	//workDir := newGitCodeDir + "/" + cicdManager.Option.OpDirName
	ExecShellCommandPre := "cd " + newGitCodeDir + "  ; pwd ; "
	//ExecShellCommandPre := " ls -l "
	output1 := ""
	output2 := ""
	if serviceCICDConfig.System.Command != "" {
		cicdManager.Option.Log.Info("step 5.1 : System.Command " + serviceCICDConfig.System.Command)
		output1, err = ExecShellCommand(ExecShellCommandPre+serviceCICDConfig.System.Command, "")
		if err != nil {
			return output, errors.New("ExecShellCommand err " + err.Error())
		}
		util.MyPrint(output)
	}
	//编译项目代码
	if serviceCICDConfig.System.Build != "" {
		cicdManager.Option.Log.Info("step 5.2 : project build command " + serviceCICDConfig.System.Build)
		output2, err = ExecShellCommand(ExecShellCommandPre+serviceCICDConfig.System.Build, "")
		if err != nil {
			return output, errors.New("ExecShellCommand err " + err.Error())
		}
		util.MyPrint(output)
	}

	return output1 + " <br/> " + output2, nil
	//cicdManager.Option.Log.Info("step 6.3 :  project testUnit command "+serviceCICDConfig.System.Command)
	//if serviceCICDConfig.System.TestUnit != ""{
	//	ExecShellCommand(serviceCICDConfig.System.TestUnit,"")
	//}
}

//step 6
func (cicdManager *CicdManager) DeployOneServiceLinkMaster(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig) error {
	cicdManager.Option.Log.Info("step 6 : master dir softLink , os.Symlink:" + newGitCodeDir + " to " + serviceDeployConfig.MasterPath)
	_, err := util.PathExists(serviceDeployConfig.MasterPath)
	if err == nil {
		cicdManager.Option.Log.Info("master path exist , so need del ." + serviceDeployConfig.MasterPath)
		err = os.Remove(serviceDeployConfig.MasterPath)
		if err != nil {
			return errors.New("os.Remove " + serviceDeployConfig.MasterPath + " err:" + err.Error())
		}
	} else if os.IsNotExist(err) {

	} else {
		//return cicdManager.DeployOneServiceFailed(publish,"unkonw err:"+err.Error())
		cicdManager.Option.Log.Info("master path exist , so need del ." + serviceDeployConfig.MasterPath)
		err = os.Remove(serviceDeployConfig.MasterPath)
		if err != nil {
			return errors.New("os.Remove " + serviceDeployConfig.MasterPath + " err:" + err.Error())
		}
	}

	err = os.Symlink(newGitCodeDir, serviceDeployConfig.MasterPath)
	if err != nil {
		return errors.New("os.Symlink err :" + err.Error())
	}
	return nil
}

//部署一个服务失败，统一处理接口
func (cicdManager *CicdManager) DeployOneServiceFailed(publish model.CicdPublish, errMsg string) error {
	cicdManager.Option.PublicManager.UpDeployStatus(publish, model.CICD_PUBLISH_DEPLOY_FAIL)
	return cicdManager.MakeError(errMsg)
}

var ThirdInstance = []string{"mysql", "redis", "log", "email", "etcd", "rabbitmq", "kafka", "alert", "cdn", "consul", "sms", "prometheus", "es", "kibana", "grafana", "push_gateway"}

func (cicdManager *CicdManager) ReplaceInstance(content string, serviceName string, env int) string {
	category := ThirdInstance
	//attr := []string{"ip","port","user","ps"}
	separator := util.STR_SEPARATOR
	content = strings.Replace(content, separator+"env"+separator, strconv.Itoa(env), -1)
	projectLogDir := cicdManager.Option.Config.System.LogDir + util.DIR_SEPARATOR + serviceName

	pathNotExistCreate(projectLogDir)

	content = strings.Replace(content, separator+"log_dir"+separator, projectLogDir, -1)
	for _, v := range category {
		//for _,attrOne := range attr{
		instance, empty := cicdManager.Option.InstanceManager.GetByEnvName(env, v)
		if empty {
			//MyPrint("cicdManager.Option.InstanceManager.GetByEnvName is empty,",env,v)
			continue
		}
		key := separator + v + "_" + "ip" + separator
		content = strings.Replace(content, key, instance.Host, -1)

		key = separator + v + "_" + "port" + separator
		content = strings.Replace(content, key, instance.Port, -1)

		key = separator + v + "_" + "user" + separator
		content = strings.Replace(content, key, instance.User, -1)

		key = separator + v + "_" + "ps" + separator
		content = strings.Replace(content, key, instance.Ps, -1)

		//}
	}

	return content
}
