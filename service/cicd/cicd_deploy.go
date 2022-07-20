package cicd

import (
	"errors"
	"go.uber.org/zap"
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
	DEPLOY_TARGET_TYPE_LOCAL  = 1 //本地部署
	DEPLOY_TARGET_TYPE_REMOTE = 2 //远程部署，并同步到了本机

	DEPLOY_TARGET_TYPE_LOCAL_NAME  = "local"
	DEPLOY_TARGET_TYPE_REMOTE_NAME = "remote"
)

func (cicdManager *CicdManager) GetConstListCicdDeployTargetType() map[string]int {
	list := make(map[string]int)

	list["本要部署再同步到远端"] = DEPLOY_TARGET_TYPE_LOCAL
	list["直接在远端部署"] = DEPLOY_TARGET_TYPE_REMOTE

	return list
}

type ServiceDeployConfig struct {
	Name             string //服务名称
	BaseDir          string //所有service项目统一放在一个目录下，由host.toml 中配置
	RemoteBaseDir    string //无端机器的部署代码的基础路径
	FullPath         string //最终一个服务的目录名,BaseDir + serviceName
	MasterDirName    string //一个服务的线上使用版本-软连目录名称
	MasterPath       string //full path + MasterDirName
	CICDConfFileName string //一个服务自己的需要执行的cicd脚本
	OpDirName        string //存放所有：运维工具脚本的目录
	FullOpDirName    string //当前正在执行的脚本，运维目录
	//ConfigTmpFileName  string //一个服务的配置文件的模板文件名
	//ConfigFileName     string //一个服务的配置文件名,由上面CP
	GitCloneTmpDirName    string //git clone 一个服务的项目代码时，临时存在的目录名
	ClonePath             string //service dir + GitCloneTmpDirName ，先把代码clone 到这个目录下面,后续再转移
	CodeGitClonePath      string // ClonePath + service name ,，之后再重合名(文件名：unixTime + gitCommitId)，转移到service目录下
	SuperVisorConfDir     string //superVisor配置文件存放目录
	SuperConfTemplateFile string //superVisor原配置文件,模板文件，用这个文件再生成每个项目的superVisor配置文件
	CICDShellFileName     string //有一一些操作需要借用于shell 执行，如：git clone . 该变量就是shell 文件名
	DeployTargetType      int    //1本地部署2远端部署
}

//公共变量
type DeployOption struct {
	ServerList      map[int]util.Server  //所有服务器
	ServiceList     map[int]util.Service //所有项目/服务
	ProjectList     map[int]util.Project
	InstanceManager *util.InstanceManager
	Config          ConfigCicd
	PublicManager   *CICDPublicManager
	Log             *zap.Logger
	OpDirName       string
	TestServerList  []string
	UploadDiskPath  string
}

type Deploy struct {
	Option CicdManagerOption
}

func NewDeploy(option CicdManagerOption) *Deploy {
	deploy := new(Deploy)
	deploy.Option = option
	return deploy
}

func (deploy *Deploy) MakeError(errMsg string) error {
	deploy.Option.Log.Error(errMsg)
	return errors.New(errMsg)
}

func (deploy *Deploy) ApiDeployOneService(form request.CicdDeploy) error {
	server, project, err := deploy.CheckCicdRequestForm(form)
	if err != nil {
		return err
	}
	serviceDeployConfig := deploy.GetDeployConfig(form.Flag)
	_, _, err = deploy.OneService(server, serviceDeployConfig, project)
	return err
}

func (deploy *Deploy) CheckCicdRequestForm(form request.CicdDeploy) (server util.Server, service util.Project, err error) {
	server, ok := deploy.Option.ServerList[form.ServerId]
	if !ok {
		return server, service, errors.New("serviceId not in list")
	}
	service, ok = deploy.Option.ProjectList[form.ServiceId]
	if !ok {
		return server, service, errors.New("serviceId not in list")
	}

	if form.Flag <= 0 {
		return server, service, errors.New("Flag <= 0")
	}

	if form.Flag != DEPLOY_TARGET_TYPE_LOCAL && form.Flag != DEPLOY_TARGET_TYPE_REMOTE {
		return server, service, errors.New("Flag err")
	}

	return server, service, nil
}

//一次部署: 所有服务器的全部服务项目
func (deploy *Deploy) AllService(deployTargetType int) {
	util.MyPrint("DeployAllService:")
	serviceDeployConfig := deploy.GetDeployConfig(deployTargetType)
	//先遍历所有服务器，然后，把所有已知服务部署到每台服务器上(每台机器都可以部署任何服务)
	for _, server := range deploy.Option.ServerList {
		//if server.OutIp != "8.142.161.156" {
		//	continue
		//}
		//遍历所有服务
		for _, service := range deploy.Option.ProjectList {
			publishId, _, err := deploy.OneService(server, serviceDeployConfig, service)
			util.MyPrint("DeployOneService err:", err, " publishId:", publishId)
			if err == nil {
				err = deploy.Publish(publishId, deployTargetType)
				util.MyPrint("DeployOneService err:", err)
			}

			//if err != nil {
			//	util.ExitPrint(err)
			//}
		}
	}
}

func (deploy *Deploy) GetDeployConfig(deployTargetType int) ServiceDeployConfig {
	serviceBaseDir := ""
	if deployTargetType == DEPLOY_TARGET_TYPE_REMOTE {
		serviceBaseDir = deploy.Option.Config.System.WorkBaseDir + "/" + DEPLOY_TARGET_TYPE_REMOTE_NAME + "/"
	} else if deployTargetType == DEPLOY_TARGET_TYPE_LOCAL {
		serviceBaseDir = deploy.Option.Config.System.WorkBaseDir + "/" + DEPLOY_TARGET_TYPE_LOCAL_NAME + "/"
	} else {
		util.ExitPrint("deployTargetType err:", deployTargetType)
	}

	//cicdManager.Option.Log.Info("DeployAllService:")
	serviceDeployConfig := ServiceDeployConfig{
		DeployTargetType:   deployTargetType,
		BaseDir:            serviceBaseDir,
		RemoteBaseDir:      deploy.Option.Config.System.RemoteBaseDir,
		OpDirName:          deploy.Option.OpDirName,
		MasterDirName:      deploy.Option.Config.System.MasterDirName,
		GitCloneTmpDirName: deploy.Option.Config.System.GitCloneTmpDirName,

		CICDConfFileName:      "cicd.toml", //本项目的相关 脚本及配置,这个是写死的，与程序员约定好，且里面的内容由程序DIY
		CICDShellFileName:     "cicd.sh",   //执行：git clone 代码，并获取当前git最新版本号
		SuperVisorConfDir:     deploy.Option.Config.SuperVisor.ConfDir,
		SuperConfTemplateFile: deploy.Option.Config.SuperVisor.ConfTemplateFile,
	}
	return serviceDeployConfig
}

//部署时，如果是测试，指定一些参数即可，不用全部署
func (deploy *Deploy) CheckTest(server util.Server, serviceDeployConfig ServiceDeployConfig, service util.Project) error {
	//if server.OutIp != ""{
	//	return errors.New("CheckTest is err outIp != ''")
	//}
	if server.Env != 5 { //测试代码,只部署：正式
		errMsg := " server.Env != 5 "
		util.MyPrint(errMsg)
		return errors.New(errMsg)
	}
	//目前仅允许这3个项目部署，3个全开放，是给HTTP使用，指令行测试，把下面两个打开即可
	test_allow_project_name := []string{"Zgoframe", "Zwebuivue", "Zwebuivgo"}
	//test_allow_project_name := []string{"Zwebuivue"}
	//test_allow_project_name := []string{"Zwebuivgo"}
	search := 0
	for _, v := range test_allow_project_name {
		if service.Name == v { //测试代码,只部署：选择的项目
			search = 1

		}
	}

	if search == 0 {
		errMsg := "test_allow_project_name service name no search : " + service.Name
		util.MyPrint(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

//部署一个服务
func (deploy *Deploy) OneService(server util.Server, serviceDeployConfig ServiceDeployConfig, service util.Project) (publishId int, deployOneServiceFlowRecord DeployOneServiceFlowRecord, err error) {
	startTime := util.GetNowTimeSecondToInt()
	checkTestRs := deploy.CheckTest(server, serviceDeployConfig, service)
	if checkTestRs != nil {
		return publishId, deployOneServiceFlowRecord, checkTestRs
	}
	//记录一个服务：部署的整个执行过程
	deployOneServiceFlowRecord = DeployOneServiceFlowRecord{
		StartTime:           startTime,
		Step:                0,
		ServiceDeployConfig: serviceDeployConfig,
		Server:              server,
		service:             service,
	}
	deploy.Option.Log.Info("DeployOneService:" + server.OutIp + " Env:" + strconv.Itoa(server.Env) + " serviceName: " + service.Name)
	//创建一条发布记录
	publish, err := deploy.Option.PublicManager.InsertOne(service, server, serviceDeployConfig.DeployTargetType)
	if err != nil {
		return 0, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, "PublicManager.InsertOne err:"+err.Error())
	}
	deployOneServiceFlowRecord.PublishId = publish.Id
	deploy.Option.Log.Info("create one publish record , id: " + strconv.Itoa(publish.Id))
	//step 1 : 预处理，检查各种路径是否正确
	deployOneServiceFlowRecord.Step = 1
	newServiceDeployConfig, err := deploy.DeployServiceCheck(serviceDeployConfig, service, server)
	deployOneServiceFlowRecord.NewServiceDeployConfig = newServiceDeployConfig
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, err.Error())
	}
	deploy.Option.Log.Info("DeployServiceCheck ok~")
	serviceDeployConfig = newServiceDeployConfig
	//step 2 : 项目代码及目录(git)相关
	deployOneServiceFlowRecord.Step = 2
	newGitCodeDir, projectDirName, gitLastCommitId, err := deploy.DeployOneServiceGitCode(serviceDeployConfig, service)

	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, err.Error())
	}

	p := model.CicdPublish{}
	p.Id = publish.Id
	p.CodeDir = projectDirName //保留目录名，使用浏览器操作的时候，切换master使用
	deploy.Option.PublicManager.UpInfo(p)
	//step 3 : 读取service项目代码里自带的cicd.toml ,供:后面使用
	deployOneServiceFlowRecord.Step = 3
	util.MyPrint("newGitCodeDir:", newGitCodeDir)
	deployOneServiceFlowRecord.NewGitCodeDir = newGitCodeDir

	serviceCICDConfig, serviceSelfCICDConf, err := deploy.DeployOneServiceCICIConfig(newGitCodeDir, serviceDeployConfig, server, gitLastCommitId)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, err.Error())
	}
	//deployOneServiceFlowRecord.ShowDeployOneServiceRecord()
	//util.ExitPrint(22)

	deployOneServiceFlowRecord.ServiceSelfCicdFile = serviceSelfCICDConf
	deployOneServiceFlowRecord.Step = 4
	//step 4: 生成该服务的，superVisor 配置文件
	err = deploy.DeployOneServiceSuperVisor(serviceDeployConfig, serviceCICDConfig, newGitCodeDir)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, err.Error())
	}
	//step 5: 处理项目自带的主配置文件
	deployOneServiceFlowRecord.Step = 5
	serviceSelfConfigTmpFileDir, serviceSelfConfigFileDir, err := deploy.DeployOneServiceProjectConfig(newGitCodeDir, server, serviceDeployConfig, serviceCICDConfig, service)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, err.Error())
	}
	deployOneServiceFlowRecord.ServiceSelfConfigTmpFileDir = serviceSelfConfigTmpFileDir
	deployOneServiceFlowRecord.ServiceSelfCicdFile = serviceSelfConfigFileDir
	//step 6 : 先执行 服务自带的 shell 预处理
	deployOneServiceFlowRecord.Step = 6
	shellCICDCommand, shellCICDBuild, _, err := deploy.DeployOneServiceCommand(newGitCodeDir, serviceDeployConfig, serviceCICDConfig)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, deploy.DeployOneServiceFailed(publish, err.Error())
	}
	deployOneServiceFlowRecord.ShellCICDCommand = shellCICDCommand
	deployOneServiceFlowRecord.ShellCICDBuild = shellCICDBuild
	//step 7 : 目前均是在本机部署的代码，现在要将代码同步到服务器上
	deployOneServiceFlowRecord.Step = 7
	syncCodeShellCommand, syncSuperVisorShellCommand, err := deploy.SyncOneServiceToRemote(serviceDeployConfig, server, newGitCodeDir, service)
	deployOneServiceFlowRecord.SyncCodeShellCommand = syncCodeShellCommand
	deployOneServiceFlowRecord.SyncSuperVisorShellCommand = syncSuperVisorShellCommand
	//更新部署的状态
	deploy.Option.PublicManager.UpDeployStatus(publish, model.CICD_PUBLISH_DEPLOY_STATUS_FINISH)
	//更新本条发布记录的状态
	deploy.Option.PublicManager.UpStatus(publish, model.CICD_PUBLISH_STATUS_WAIT_PUB)
	//更新本条发布记录的基础信息
	endTime := util.GetNowTimeSecondToInt()
	execTime := endTime - startTime
	e := model.CicdPublish{}
	e.Id = publish.Id
	e.ExecTime = execTime
	e.Log = deployOneServiceFlowRecord.ToString()
	e.Step = deployOneServiceFlowRecord.Step
	deploy.Option.PublicManager.UpInfo(e)

	deployOneServiceFlowRecord.EndTime = endTime
	deployOneServiceFlowRecord.ShowDeployOneServiceRecord()

	return publish.Id, deployOneServiceFlowRecord, nil
}

func GetRsyncCommandPre() string {
	return "rsync -avz --progress --port=8877 "
}

func (deploy *Deploy) Publish(id int, deployTargetType int) error {
	serviceDeployConfig := deploy.GetDeployConfig(deployTargetType)
	publishRecord, err := deploy.Option.PublicManager.GetById(id)
	if err != nil {
		return err
	}

	server := deploy.Option.ServerList[publishRecord.ServerId]
	service := deploy.Option.ProjectList[publishRecord.ServiceId]
	serviceDeployConfig, _ = deploy.DeployServiceCheck(serviceDeployConfig, service, server)
	//step 8 将master软链 指向 上面刚刚clone下的最新代码上
	err = deploy.DeployOneServiceLinkMaster(publishRecord.CodeDir, serviceDeployConfig)
	if err != nil {
		deploy.Option.PublicManager.UpStatus(publishRecord, model.CICD_PUBLISH_DEPLOY_FAIL)
		return err
		//return cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	if service.Type == util.PROJECT_TYPE_SERVICE {
		//1 同步代码
		syncCodeShellCommand := GetRsyncCommandPre() + serviceDeployConfig.FullPath + "/" + serviceDeployConfig.MasterDirName + " rsync@" + server.OutIp + "::www/" + serviceDeployConfig.Name
		_, err = ExecShellCommand(syncCodeShellCommand, "")
		util.MyPrint("SyncOneServiceToRemote:", syncCodeShellCommand, " err:", err)
	} else if service.Type == util.PROJECT_TYPE_FE {
		syncCodeShellCommand := GetRsyncCommandPre() + serviceDeployConfig.FullPath + "/" + serviceDeployConfig.MasterDirName + " rsync@" + server.OutIp + "::www/" + serviceDeployConfig.Name
		_, err = ExecShellCommand(syncCodeShellCommand, "")
		util.MyPrint("SyncOneServiceToRemote:", syncCodeShellCommand, " err:", err)
	} else {
		return errors.New("SyncOneServiceToRemote :project type err.")
	}

	deploy.Option.PublicManager.UpStatus(publishRecord, model.CICD_PUBLISH_DEPLOY_OK)
	return nil
}

//部署一个服务失败，统一处理接口
func (deploy *Deploy) DeployOneServiceFailed(publish model.CicdPublish, errMsg string) error {
	deploy.Option.PublicManager.UpDeployStatus(publish, model.CICD_PUBLISH_DEPLOY_FAIL)
	return deploy.MakeError(errMsg)
}

var ThirdInstance = []string{"mysql", "redis", "log", "email", "etcd", "rabbitmq", "kafka", "alert", "cdn", "consul", "sms", "prometheus", "es", "kibana", "grafana", "push_gateway", "http", "domain", "oss", "grpc", "gateway", "agora", "super_visor"}

func (deploy *Deploy) ReplaceInstance(content string, serviceName string, env int, serviceId int) string {
	category := ThirdInstance
	//attr := []string{"ip","port","user","ps"}
	separator := util.STR_SEPARATOR
	content = strings.Replace(content, separator+"env"+separator, strconv.Itoa(env), -1)
	projectLogDir := deploy.Option.Config.System.LogDir + util.DIR_SEPARATOR + serviceName

	pathNotExistCreate(projectLogDir)

	content = strings.Replace(content, separator+"log_dir"+separator, projectLogDir, -1)
	content = strings.Replace(content, separator+"projectId"+separator, strconv.Itoa(serviceId), -1)
	for _, v := range category {
		//for _,attrOne := range attr{
		instance, empty := deploy.Option.InstanceManager.GetByEnvName(env, v)
		if empty {
			//MyPrint("cicdManager.Option.InstanceManager.GetByEnvName is empty,",env,v)
			continue
		}

		if v == "gateway" {
			host := strings.Split(instance.Host, ",")

			key := separator + v + "_" + "listen_ip" + separator
			content = strings.Replace(content, key, host[0], -1)

			key = separator + v + "_" + "out_ip" + separator
			content = strings.Replace(content, key, host[1], -1)

			ports := strings.Split(instance.Port, ",")

			key = separator + v + "_" + "ws_port" + separator
			content = strings.Replace(content, key, ports[0], -1)

			key = separator + v + "_" + "tcp_port" + separator
			content = strings.Replace(content, key, ports[1], -1)

			key = separator + v + "_" + "ws_uri" + separator
			content = strings.Replace(content, key, instance.User, -1)

			continue
		}

		key := separator + v + "_" + "ip" + separator
		content = strings.Replace(content, key, instance.Host, -1)

		key = separator + v + "_" + "host" + separator
		content = strings.Replace(content, key, instance.Host, -1)

		key = separator + v + "_" + "port" + separator
		content = strings.Replace(content, key, instance.Port, -1)

		key = separator + v + "_" + "user" + separator
		content = strings.Replace(content, key, instance.User, -1)

		key = separator + v + "_" + "ps" + separator
		content = strings.Replace(content, key, instance.Ps, -1)

		key = separator + v + "_" + "password" + separator
		content = strings.Replace(content, key, instance.Ps, -1)

		key = separator + v + "_" + "from" + separator
		content = strings.Replace(content, key, instance.Host, -1)

		key = separator + v + "_" + "auth_code" + separator
		content = strings.Replace(content, key, instance.Host, -1)

		//}
	}

	return content
}
