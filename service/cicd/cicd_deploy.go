package cicd

import (
	"bytes"
	"encoding/json"
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
	DEPLOY_TARGET_TYPE_LOCAL  = 1 //本要部署
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

func (cicdManager *CicdManager) ApiDeployOneService(form request.CicdDeploy) error {
	server, project, err := cicdManager.CheckCicdRequestForm(form)
	if err != nil {
		return err
	}
	serviceDeployConfig := cicdManager.GetDeployConfig(form.Flag)
	_, _, err = cicdManager.DeployOneService(server, serviceDeployConfig, project)
	return err
}

func (cicdManager *CicdManager) CheckCicdRequestForm(form request.CicdDeploy) (server util.Server, service util.Project, err error) {
	server, ok := cicdManager.Option.ServerList[form.ServerId]
	if !ok {
		return server, service, errors.New("serviceId not in list")
	}
	service, ok = cicdManager.Option.ProjectList[form.ServiceId]
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

//一次部署全部服务项目
func (cicdManager *CicdManager) DeployAllService(deployTargetType int) {
	util.MyPrint("DeployAllService:")
	serviceDeployConfig := cicdManager.GetDeployConfig(deployTargetType)
	//先遍历所有服务器，然后，把所有已知服务部署到每台服务器上(每台机器都可以部署任何服务)
	for _, server := range cicdManager.Option.ServerList {
		//if server.OutIp != "8.142.161.156" {
		//	continue
		//}
		//遍历所有服务
		for _, service := range cicdManager.Option.ProjectList {
			publishId, _, err := cicdManager.DeployOneService(server, serviceDeployConfig, service)
			util.MyPrint("DeployOneService err:", err, " publishId:", publishId)
			if err == nil {
				err = cicdManager.Publish(publishId, deployTargetType)
				util.MyPrint("DeployOneService err:", err)
			}

			//if err != nil {
			//	util.ExitPrint(err)
			//}
		}
	}
}

func (cicdManager *CicdManager) GetDeployConfig(deployTargetType int) ServiceDeployConfig {
	serviceBaseDir := ""
	if deployTargetType == DEPLOY_TARGET_TYPE_REMOTE {
		serviceBaseDir = cicdManager.Option.Config.System.WorkBaseDir + "/" + DEPLOY_TARGET_TYPE_REMOTE_NAME + "/"
	} else if deployTargetType == DEPLOY_TARGET_TYPE_LOCAL {
		serviceBaseDir = cicdManager.Option.Config.System.WorkBaseDir + "/" + DEPLOY_TARGET_TYPE_LOCAL_NAME + "/"
	} else {
		util.ExitPrint("deployTargetType err:", deployTargetType)
	}

	//cicdManager.Option.Log.Info("DeployAllService:")
	serviceDeployConfig := ServiceDeployConfig{
		DeployTargetType:   deployTargetType,
		BaseDir:            serviceBaseDir,
		RemoteBaseDir:      cicdManager.Option.Config.System.RemoteBaseDir,
		OpDirName:          cicdManager.Option.OpDirName,
		MasterDirName:      cicdManager.Option.Config.System.MasterDirName,
		GitCloneTmpDirName: cicdManager.Option.Config.System.GitCloneTmpDirName,

		CICDConfFileName:      "cicd.toml", //本项目的相关 脚本及配置,这个是写死的，与程序员约定好，且里面的内容由程序DIY
		CICDShellFileName:     "cicd.sh",   //执行：git clone 代码，并获取当前git最新版本号
		SuperVisorConfDir:     cicdManager.Option.Config.SuperVisor.ConfDir,
		SuperConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
	}
	return serviceDeployConfig
}

func (cicdManager *CicdManager) DeployServiceCheck(serviceDeployConfig ServiceDeployConfig, service util.Project, server util.Server) (ServiceDeployConfig, error) {
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
	if err != nil {
		if os.IsNotExist(err) {
			util.MyPrint("DeployServiceCheck create dir:", serviceDeployConfig.BaseDir)
			err = os.Mkdir(serviceDeployConfig.BaseDir, 0777)
			if err != nil {
				util.ExitPrint("os.Mkdir :", serviceDeployConfig.BaseDir, " err:", err.Error())
			}
		}
	}
	//本机部分编译，要把远程部署多出一层： 服务器IP目录->服务目录
	if serviceDeployConfig.DeployTargetType == DEPLOY_TARGET_TYPE_REMOTE {
		newBaseDir := serviceDeployConfig.BaseDir + "/" + server.OutIp
		_, err = util.PathExists(newBaseDir)
		if err != nil {
			if os.IsNotExist(err) {
				util.MyPrint("DEPLOY_TARGET_TYPE_LOCAL create dir:", newBaseDir)
				err = os.Mkdir(newBaseDir, 0777)
				if err != nil {
					util.ExitPrint("os.Mkdir :", newBaseDir, " err:", err.Error())
				}
			} else {
				util.ExitPrint("util.PathExists err:", err.Error())
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
	serviceDeployConfig.FullOpDirName = cicdManager.Option.Config.System.RootDir + "/" + serviceDeployConfig.OpDirName

	//serviceDeployConfig.RemoteBaseDir = serviceDeployConfig.RemoteBaseDir
	newServiceDeployConfig := serviceDeployConfig

	//util.PrintStruct(newServiceDeployConfig, ":")

	return newServiceDeployConfig, nil
}

//部署时，如果是测试，指定一些参数即可，不用全部署
func (cicdManager *CicdManager) CheckTest(server util.Server, serviceDeployConfig ServiceDeployConfig, service util.Project) error {
	//if server.OutIp != ""{
	//	return errors.New("CheckTest is err outIp != ''")
	//}
	//if server.Env != 1 &&  server.Env != 4 { //测试代码,只部署：local Zgoframe
	//	errMsg := " server.Env != 1 "
	//	util.MyPrint(errMsg)
	//	return errors.New(errMsg)
	//}

	test_allow_project_name := []string{"Zgoframe"}
	//test_allow_project_name := []string{"Zwebuivue"}
	//test_allow_project_name := []string{"Zwebuivgo"}
	for _, v := range test_allow_project_name {
		if service.Name != v { //测试代码,只部署：选择的项目
			errMsg := "test_allow_project_name service name != " + v
			util.MyPrint(errMsg)
			return errors.New(errMsg)
		}
	}

	return nil
}

//记录一次部署的全过程
type DeployOneServiceFlowRecord struct {
	ServiceDeployConfig    ServiceDeployConfig `json:"service_deploy_config"`
	NewServiceDeployConfig ServiceDeployConfig `json:"new_service_deploy_config"`
	Server                 util.Server         `json:"-"` //mysql里有存这个字段
	service                util.Project        `json:"-"` //mysql里有存这个字段

	StartTime                   int    `json:"start_time"`
	EndTime                     int    `json:"end_time"`
	Step                        int    `json:"step"`
	PublishId                   int    `json:"publish_id"`
	ServiceSelfCicdFile         string `json:"service_self_cicd_file"`           //clone 出来的项目代码中，自带的cicd.config 文件
	NewGitCodeDir               string `json:"new_git_code_dir"`                 //最终的代码clone 目录
	ServiceSelfConfigTmpFileDir string `json:"service_self_config_tmp_file_dir"` //项目自带的配置文件模板
	serviceSelfConfigFileDir    string `json:"service_self_config_file_dir"`     //项目自带的配置文件，由上面替换变量后，最终文件
	ShellCICDCommand            string `json:"shell_cicd_command"`               //项目自带的cicd文件，执行预处理shell脚本
	ShellCICDBuild              string `json:"shell_cicd_build"`                 //项目自带的cicd文件，对代码进行编译
	SyncCodeShellCommand        string `json:"sync_code_shell_command"`          //代码同步到远端指令
	SyncSuperVisorShellCommand  string `json:"sync_super_visor_shell_command"`   //supervisor 配置文件 - 同步到远端指令
}

func (deployOneServiceFlowRecord *DeployOneServiceFlowRecord) ShowDeployOneServiceRecord() {
	util.MyPrint("ShowDeployOneServiceRecord:")
	util.MyPrint("ServiceSelfCicdFile:" + deployOneServiceFlowRecord.ServiceSelfCicdFile)
	util.MyPrint("NewGitCodeDir:" + deployOneServiceFlowRecord.NewGitCodeDir)
	util.MyPrint("serviceSelfConfigFileDir:" + deployOneServiceFlowRecord.serviceSelfConfigFileDir)
	util.MyPrint("ServiceSelfConfigTmpFileDir:" + deployOneServiceFlowRecord.ServiceSelfConfigTmpFileDir)
	util.MyPrint("ShellCICDBuild:" + deployOneServiceFlowRecord.ShellCICDBuild)
	util.MyPrint("ShellCICDCommand:" + deployOneServiceFlowRecord.ShellCICDCommand)
	util.PrintStruct(deployOneServiceFlowRecord.NewServiceDeployConfig, ":")

}
func (deployOneServiceFlowRecord *DeployOneServiceFlowRecord) ToString() string {
	str, err := json.Marshal(deployOneServiceFlowRecord)
	if err != nil {
		util.MyPrint("err deployOneServiceFlowRecord json.Marshal:", err.Error())
	}
	util.MyPrint("toString:", string(str))
	return string(str)
}

//部署一个服务
func (cicdManager *CicdManager) DeployOneService(server util.Server, serviceDeployConfig ServiceDeployConfig, service util.Project) (publishId int, deployOneServiceFlowRecord DeployOneServiceFlowRecord, err error) {
	startTime := util.GetNowTimeSecondToInt()
	checkTestRs := cicdManager.CheckTest(server, serviceDeployConfig, service)
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
	cicdManager.Option.Log.Info("DeployOneService:" + server.OutIp + " Env:" + strconv.Itoa(server.Env) + " serviceName: " + service.Name)
	//创建一条发布记录
	publish, err := cicdManager.Option.PublicManager.InsertOne(service, server, serviceDeployConfig.DeployTargetType)
	if err != nil {
		return 0, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, "PublicManager.InsertOne err:"+err.Error())
	}
	deployOneServiceFlowRecord.PublishId = publish.Id
	cicdManager.Option.Log.Info("create one publish record , id: " + strconv.Itoa(publish.Id))
	//step 1 : 预处理，检查各种路径是否正确
	deployOneServiceFlowRecord.Step = 1
	newServiceDeployConfig, err := cicdManager.DeployServiceCheck(serviceDeployConfig, service, server)
	deployOneServiceFlowRecord.NewServiceDeployConfig = newServiceDeployConfig
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	cicdManager.Option.Log.Info("DeployServiceCheck ok~")
	serviceDeployConfig = newServiceDeployConfig
	//step 2 : 项目代码及目录(git)相关
	deployOneServiceFlowRecord.Step = 2
	newGitCodeDir, projectDirName, gitLastCommitId, err := cicdManager.DeployOneServiceGitCode(serviceDeployConfig, service)

	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, err.Error())
	}

	p := model.CicdPublish{}
	p.Id = publish.Id
	p.CodeDir = projectDirName //保留目录名，使用浏览器操作的时候，切换master使用
	cicdManager.Option.PublicManager.UpInfo(p)
	//step 3 : 读取service项目代码里自带的cicd.toml ,供:后面使用
	deployOneServiceFlowRecord.Step = 3
	util.MyPrint("newGitCodeDir:", newGitCodeDir)
	deployOneServiceFlowRecord.NewGitCodeDir = newGitCodeDir

	serviceCICDConfig, serviceSelfCICDConf, err := cicdManager.DeployOneServiceCICIConfig(newGitCodeDir, serviceDeployConfig, server, gitLastCommitId)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, err.Error())
	}

	deployOneServiceFlowRecord.ServiceSelfCicdFile = serviceSelfCICDConf
	deployOneServiceFlowRecord.Step = 4
	//step 4: 生成该服务的，superVisor 配置文件
	err = cicdManager.DeployOneServiceSuperVisor(serviceDeployConfig, serviceCICDConfig, newGitCodeDir)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	//step 5: 处理项目自带的主配置文件
	deployOneServiceFlowRecord.Step = 5
	serviceSelfConfigTmpFileDir, serviceSelfConfigFileDir, err := cicdManager.DeployOneServiceProjectConfig(newGitCodeDir, server, serviceDeployConfig, serviceCICDConfig, service)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	deployOneServiceFlowRecord.ServiceSelfConfigTmpFileDir = serviceSelfConfigTmpFileDir
	deployOneServiceFlowRecord.ServiceSelfCicdFile = serviceSelfConfigFileDir
	//step 6 : 先执行 服务自带的 shell 预处理
	deployOneServiceFlowRecord.Step = 6
	shellCICDCommand, shellCICDBuild, _, err := cicdManager.DeployOneServiceCommand(newGitCodeDir, serviceDeployConfig, serviceCICDConfig)
	if err != nil {
		return publish.Id, deployOneServiceFlowRecord, cicdManager.DeployOneServiceFailed(publish, err.Error())
	}
	deployOneServiceFlowRecord.ShellCICDCommand = shellCICDCommand
	deployOneServiceFlowRecord.ShellCICDBuild = shellCICDBuild
	//step 7 : 目前均是在本机部署的代码，现在要将代码同步到服务器上
	deployOneServiceFlowRecord.Step = 7
	syncCodeShellCommand, syncSuperVisorShellCommand, err := cicdManager.SyncOneServiceToRemote(serviceDeployConfig, server, newGitCodeDir, service)
	deployOneServiceFlowRecord.SyncCodeShellCommand = syncCodeShellCommand
	deployOneServiceFlowRecord.SyncSuperVisorShellCommand = syncSuperVisorShellCommand
	//更新部署的状态
	cicdManager.Option.PublicManager.UpDeployStatus(publish, model.CICD_PUBLISH_DEPLOY_STATUS_FINISH)
	//更新本条发布记录的状态
	cicdManager.Option.PublicManager.UpStatus(publish, model.CICD_PUBLISH_STATUS_WAIT_PUB)
	//更新本条发布记录的基础信息
	endTime := util.GetNowTimeSecondToInt()
	execTime := endTime - startTime
	e := model.CicdPublish{}
	e.Id = publish.Id
	e.ExecTime = execTime
	e.Log = deployOneServiceFlowRecord.ToString()
	cicdManager.Option.PublicManager.UpInfo(e)

	deployOneServiceFlowRecord.EndTime = endTime
	deployOneServiceFlowRecord.ShowDeployOneServiceRecord()

	return publish.Id, deployOneServiceFlowRecord, nil
}

func GetRsyncCommandPre() string {
	return "rsync -avz --progress --port=8877 "
}

//本机部署均已完成，需要将本地代码同步到远端
func (cicdManager *CicdManager) SyncOneServiceToRemote(serviceDeployConfig ServiceDeployConfig, server util.Server, newGitCodeDir string, project util.Project) (syncCodeShellCommand string, syncSuperVisorShellCommand string, err error) {
	if project.Type == util.PROJECT_TYPE_SERVICE {
		//1 同步代码
		syncCodeShellCommand = GetRsyncCommandPre() + " --exclude=master " + serviceDeployConfig.FullPath + " rsync@" + server.OutIp + "::www"
		_, err := ExecShellCommand(syncCodeShellCommand, "")
		util.MyPrint("SyncOneServiceToRemote:", syncCodeShellCommand, " err:", err)
		//2 同步superVisor
		syncSuperVisorShellCommand = GetRsyncCommandPre() + newGitCodeDir + "/" + serviceDeployConfig.Name + ".ini" + " rsync@" + server.OutIp + "::super_visor"
		_, err = ExecShellCommand(syncSuperVisorShellCommand, "")
		util.MyPrint("syncSuperVisorShellCommand:", syncSuperVisorShellCommand, " err:", err)
	} else if project.Type == util.PROJECT_TYPE_FE {
		//util.MyPrint(serviceDeployConfig)
		syncCodeShellCommand = GetRsyncCommandPre() + " --exclude=node_modules " + newGitCodeDir + " rsync@" + server.OutIp + "::www/" + serviceDeployConfig.Name
		//util.ExitPrint(syncCodeShellCommand)
		_, err := ExecShellCommand(syncCodeShellCommand, "")
		util.MyPrint("SyncOneServiceToRemote:", syncCodeShellCommand, " err:", err)
	} else {
		return "", "", errors.New("SyncOneServiceToRemote :project type err.")
	}

	return syncCodeShellCommand, syncSuperVisorShellCommand, nil
}

func (cicdManager *CicdManager) Publish(id int, deployTargetType int) error {
	serviceDeployConfig := cicdManager.GetDeployConfig(deployTargetType)
	publishRecord, err := cicdManager.Option.PublicManager.GetById(id)
	if err != nil {
		return err
	}

	server := cicdManager.Option.ServerList[publishRecord.ServerId]
	service := cicdManager.Option.ProjectList[publishRecord.ServiceId]
	serviceDeployConfig, _ = cicdManager.DeployServiceCheck(serviceDeployConfig, service, server)
	//将master软链 指向 上面刚刚clone下的最新代码上
	err = cicdManager.DeployOneServiceLinkMaster(publishRecord.CodeDir, serviceDeployConfig)
	if err != nil {
		cicdManager.Option.PublicManager.UpStatus(publishRecord, model.CICD_PUBLISH_DEPLOY_FAIL)
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

	cicdManager.Option.PublicManager.UpStatus(publishRecord, model.CICD_PUBLISH_DEPLOY_OK)
	return nil
}

//step 1
func (cicdManager *CicdManager) DeployOneServiceGitCode(serviceDeployConfig ServiceDeployConfig, service util.Project) (string, string, string, error) {
	cicdManager.Option.Log.Info("step 1 : git clone project code and get git commit id.")
	//FullPath 一个服务的根目录，大部分操作都在这个目录下(除了superVisor)
	//查看服务根目录是否存在，不存在即新创建
	pathNotExistCreate(serviceDeployConfig.FullPath)
	//查看git clone 目录是否存在，不存在即新创建
	pathNotExistCreate(serviceDeployConfig.ClonePath)
	//构建 shell 执行时所需 参数
	shellArgc := service.Git + " " + serviceDeployConfig.ClonePath + " " + service.Name + " " + cicdManager.Option.Config.System.RemoteUploadDir + " " + cicdManager.Option.Config.System.UploadPath
	//执行shell 脚本 后：service项目代码已被clone, git 版本号已知了

	gitLastCommitId, err := ExecShellFile(serviceDeployConfig.FullOpDirName+"/"+serviceDeployConfig.CICDShellFileName, shellArgc)
	if err != nil {
		return "", "", "", errors.New("ExecShellFile err:" + err.Error())
	}
	//新刚刚克隆好的项目目录，移动一个新目录下，新目录名：git_master_versionId + 当前时间
	projectDirName := strconv.Itoa(util.GetNowTimeSecondToInt()) + "_" + gitLastCommitId
	newGitCodeDir := serviceDeployConfig.FullPath + util.DIR_SEPARATOR + projectDirName
	cicdManager.Option.Log.Info(" service code move :" + serviceDeployConfig.CodeGitClonePath + " to " + newGitCodeDir)
	//执行 移动操作
	err = os.Rename(serviceDeployConfig.CodeGitClonePath, newGitCodeDir)
	if err != nil {
		return newGitCodeDir, "", "", errors.New("serviceCodeGitClonePath os.Rename err:" + err.Error())
	}
	cicdManager.Option.Log.Info("step 1 finish , newGitCodeDir :  " + newGitCodeDir + " , gitLastCommitId:" + gitLastCommitId)

	//处理图片目录 的软件 连接
	//_, err := util.FileExist(cicdManager.Option.UploadDiskPath)
	//cicdManager.Option.Log.Info("ln -s " + cicdManager.Option.Config.System.RemoteUploadDir + " " + cicdManager.Option.UploadDiskPath)
	//err = os.Symlink(cicdManager.Option.Config.System.RemoteUploadDir,cicdManager.Option.UploadDiskPath)
	//if err != nil{
	//	return newGitCodeDir , projectDirName ,gitLastCommitId, errors.New("link file upload err:" + err.Error())
	//}
	//util.ExitPrint(33)

	return newGitCodeDir, projectDirName, gitLastCommitId, nil
}

//step 2
func (cicdManager *CicdManager) DeployOneServiceCICIConfig(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig, server util.Server, gitLastCommitId string) (ConfigServiceCICD, string, error) {
	cicdManager.Option.Log.Info("step 2:load service CICD config ")
	//项目自带的CICD配置文件，这里有 服务启动脚本 和 依赖的环境
	serviceSelfCICDConf := newGitCodeDir + util.DIR_SEPARATOR + serviceDeployConfig.CICDConfFileName
	cicdManager.Option.Log.Info("read file:" + serviceSelfCICDConf)
	serviceCICDConfig := ConfigServiceCICD{}
	//读取项目自己的cicd配置文件，并映射到结构体中
	err := util.ReadConfFileAutoExt(serviceSelfCICDConf, &serviceCICDConfig)
	if err != nil {
		return serviceCICDConfig, serviceSelfCICDConf, errors.New(err.Error())
	}
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#service_name#", serviceDeployConfig.Name, -1)
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#datetime#", strconv.Itoa(util.GetNowTimeSecondToInt()), -1)
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#git_version#", gitLastCommitId, -1)
	//util.MyPrint(serviceCICDConfig.System.Build)
	//util.ExitPrint(33)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#env#", strconv.Itoa(server.Env), -1)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#master_path#", serviceDeployConfig.RemoteBaseDir+"/"+serviceDeployConfig.Name+"/"+serviceDeployConfig.MasterDirName, -1)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#service_name#", serviceDeployConfig.Name, -1)

	//util.ExitPrint(serviceCICDConfig.System.Startup)
	//util.PrintStruct(serviceCICDConfig, ":")

	return serviceCICDConfig, serviceSelfCICDConf, nil
}

//step 3 生成该服务的，superVisor 配置文件
func (cicdManager *CicdManager) DeployOneServiceSuperVisor(serviceDeployConfig ServiceDeployConfig, configServiceCICD ConfigServiceCICD, newGitCodeDir string) error {
	cicdManager.Option.Log.Info("step 3 : create superVisor conf file.")
	superVisorOption := util.SuperVisorOption{
		ConfDir:     cicdManager.Option.Config.SuperVisor.ConfDir,
		ServiceName: serviceDeployConfig.Name,
		//ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
	}

	serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
	if err != nil {
		return err
	}
	serviceSuperVisor.SetConfTemplateFile(cicdManager.Option.Config.SuperVisor.ConfTemplateFile)
	//superVisor 配置文件中 动态的占位符，需要替换掉
	superVisorReplace := util.SuperVisorReplace{}
	if serviceDeployConfig.DeployTargetType == DEPLOY_TARGET_TYPE_REMOTE {
		superVisorReplace = util.SuperVisorReplace{
			ScriptName:           serviceDeployConfig.Name,
			StartupScriptCommand: configServiceCICD.System.Startup,
			ScriptWorkDir:        serviceDeployConfig.RemoteBaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name + "/" + serviceDeployConfig.MasterDirName,
			StdoutLogfile:        serviceDeployConfig.RemoteBaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name + "/super_visor_stdout.log",
			StderrLogfile:        serviceDeployConfig.RemoteBaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name + "/super_visor_stderr.log",
			ProcessName:          serviceDeployConfig.Name,
		}
	} else {
		superVisorReplace = util.SuperVisorReplace{
			ScriptName:           serviceDeployConfig.Name,
			StartupScriptCommand: configServiceCICD.System.Startup,
			ScriptWorkDir:        serviceDeployConfig.MasterPath,
			StdoutLogfile:        serviceDeployConfig.BaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name + "/super_visor_stdout.log",
			StderrLogfile:        serviceDeployConfig.BaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name + "/super_visor_stderr.log",
			ProcessName:          serviceDeployConfig.Name,
		}
	}

	//util.PrintStruct(superVisorReplace,":")
	//替换配置文件中的动态值，并生成配置文件
	serviceConfFileContent, _ := serviceSuperVisor.ReplaceConfTemplate(superVisorReplace)
	//util.ExitPrint(serviceConfFileContent)
	//将已替换好的文件，生成一个新的配置文件
	err = serviceSuperVisor.CreateServiceConfFile(serviceConfFileContent, newGitCodeDir)
	if err != nil {
		return err
	}

	return nil
}

//step 4
func (cicdManager *CicdManager) DeployOneServiceProjectConfig(newGitCodeDir string, server util.Server, serviceDeployConfig ServiceDeployConfig, configServiceCICD ConfigServiceCICD, service util.Project) (string, string, error) {
	cicdManager.Option.Log.Info("step 4 : create project self conf file.")
	//读取该服务自己的配置文件 config.toml
	serviceSelfConfigTmpFileDir := newGitCodeDir + util.DIR_SEPARATOR + configServiceCICD.System.ConfigTmpFileName
	_, err := util.FileExist(serviceSelfConfigTmpFileDir)
	if err != nil {
		return "", "", errors.New("serviceSelfConfigTmpFileDir CheckFileIsExist err:" + err.Error())
	}
	cicdManager.Option.Log.Info("read file:" + serviceSelfConfigTmpFileDir)
	//读取模板文件内容
	serviceSelfConfigTmpFileContent, err := util.ReadString(serviceSelfConfigTmpFileDir)
	if err != nil {
		return "", "", errors.New(err.Error())
	}
	//开始替换 服务自己配置文件中的，实例信息，如：IP PORT
	serviceSelfConfigTmpFileContentNew := cicdManager.ReplaceInstance(serviceSelfConfigTmpFileContent, serviceDeployConfig.Name, server.Env, service.Id)

	key := util.STR_SEPARATOR + "projectId" + util.STR_SEPARATOR
	serviceSelfConfigTmpFileContentNew = strings.Replace(serviceSelfConfigTmpFileContentNew, key, strconv.Itoa(service.Id), -1)

	//生成新的配置文件
	newConfig := newGitCodeDir + util.DIR_SEPARATOR + configServiceCICD.System.ConfigFileName
	newConfigFile, _ := os.Create(newConfig)
	contentByte := bytes.Trim([]byte(serviceSelfConfigTmpFileContentNew), "\x00") //NUL
	newConfigFile.Write(contentByte)

	return serviceSelfConfigTmpFileDir, newConfig, nil
}

//step 5
func (cicdManager *CicdManager) DeployOneServiceCommand(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig, serviceCICDConfig ConfigServiceCICD) (command string, build string, output string, err error) {
	cicdManager.Option.Log.Info("step 5 : DeployOneServiceCommand.")
	ExecShellCommandPre := "cd " + newGitCodeDir + "  ; pwd ; "
	output1 := ""
	output2 := ""
	if serviceCICDConfig.System.Command != "" {
		command = ExecShellCommandPre + serviceCICDConfig.System.Command
		cicdManager.Option.Log.Info("step 5.1 : System.Command: " + command)
		output1, err = ExecShellCommand(ExecShellCommandPre+serviceCICDConfig.System.Command, "")
		if err != nil {
			return command, build, output, errors.New("ExecShellCommand " + command + " err " + err.Error())
		}
	}
	//编译项目代码
	if serviceCICDConfig.System.Build != "" {
		build = ExecShellCommandPre + serviceCICDConfig.System.Build
		cicdManager.Option.Log.Info("step 5.2 : project build command :" + build)
		output2, err = ExecShellCommand(build, "")
		if err != nil {
			return command, build, output, errors.New("ExecShellCommand " + command + "  err " + err.Error())
		}
	}

	return command, build, output1 + " <br/> " + output2, nil
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

var ThirdInstance = []string{"mysql", "redis", "log", "email", "etcd", "rabbitmq", "kafka", "alert", "cdn", "consul", "sms", "prometheus", "es", "kibana", "grafana", "push_gateway", "http", "domain", "oss", "grpc", "gateway", "agora", "super_visor"}

func (cicdManager *CicdManager) ReplaceInstance(content string, serviceName string, env int, serviceId int) string {
	category := ThirdInstance
	//attr := []string{"ip","port","user","ps"}
	separator := util.STR_SEPARATOR
	content = strings.Replace(content, separator+"env"+separator, strconv.Itoa(env), -1)
	projectLogDir := cicdManager.Option.Config.System.LogDir + util.DIR_SEPARATOR + serviceName

	pathNotExistCreate(projectLogDir)

	content = strings.Replace(content, separator+"log_dir"+separator, projectLogDir, -1)
	content = strings.Replace(content, separator+"projectId"+separator, strconv.Itoa(serviceId), -1)
	for _, v := range category {
		//for _,attrOne := range attr{
		instance, empty := cicdManager.Option.InstanceManager.GetByEnvName(env, v)
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
