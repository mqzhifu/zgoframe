package util

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"zgoframe/model"
)
/*
自动化部署，从DB中读取出所有信息基础信息，GIT CLONE 配置super visor 监听进程
依赖
	supervisor 依赖 python 、 xmlrpc
	代码依赖：git
*/

type ServiceDeployConfig struct {
	Name 				string 	//服务名称
	BaseDir 			string	//所有service项目统一放在一个目录下，由host.toml 中配置
	FullPath 			string 	//最终一个服务的目录名,BaseDir + serviceName
	MasterDirName 		string	//一个服务的线上使用版本-软连目录名称
	MasterPath			string  //full path + MasterDirName
	CICDConfFileName 	string	//一个服务自己的需要执行的cicd脚本
	ConfigTmpFileName 	string	//一个服务的配置文件的模板文件名
	ConfigFileName 		string	//一个服务的配置文件名,由上面CP
	GitCloneTmpDirName 	string	//git clone 一个服务的项目代码时，临时存在的目录名
	ClonePath			string 	//service dir + GitCloneTmpDirName
	CodeGitClonePath   	string  // ClonePath + service name
	CICDShellFileName 	string	//有一一些操作需要借用于shell 执行，如：git clone . 该变量就是shell 文件名
}

//一次部署全部服务项目
func(cicdManager *CicdManager)DeployAllService(){
	cicdManager.Option.Log.Info("DeployAllService:")
	serviceDeployConfig := ServiceDeployConfig{
		BaseDir 			: cicdManager.Option.Config.System.ServiceDir,
		MasterDirName 		: "master",
		CICDConfFileName 	: "cicd.toml",
		ConfigTmpFileName	: "config.toml.tmp",
		ConfigFileName 		: "config.toml",
		GitCloneTmpDirName 	: "clone",
		CICDShellFileName 	: "./cicd.sh",
	}
	PrintStruct(serviceDeployConfig , ":")

	//先遍历所有服务器，然后，把所有已知服务部署到每台服务器上(每台机器都可以部署任何服务)
	for _,server :=range cicdManager.Option.ServerList{
		//遍历所有服务
		for _,service :=range cicdManager.Option.ServiceList{
			err := cicdManager.DeployOneService(server,serviceDeployConfig,service)
			if err != nil{
				ExitPrint(err)
			}
		}
	}
}

func(cicdManager *CicdManager)DeployServiceCheck(server Server , serviceDeployConfig ServiceDeployConfig ,  service Service)(ServiceDeployConfig,error){
	if service.Git == ""{
		errMsg := "service.Git is empty~"+service.Name
		return serviceDeployConfig,errors.New(errMsg)
	}

	if service.Name == ""{
		errMsg := "service.Name is empty~"
		return serviceDeployConfig,errors.New(errMsg)
	}

	if serviceDeployConfig.MasterDirName == ""{
		errMsg := "MasterDirName is empty~"
		return serviceDeployConfig,errors.New(errMsg)
	}

	if serviceDeployConfig.GitCloneTmpDirName == ""{
		errMsg := "GitCloneTmpDirName is empty~"
		return serviceDeployConfig,errors.New(errMsg)
	}

	//baseDir 已由 构造函数做校验了

	serviceDeployConfig.Name 				= service.Name
	serviceDeployConfig.FullPath 			= serviceDeployConfig.BaseDir + DIR_SEPARATOR + serviceDeployConfig.Name
	serviceDeployConfig.MasterPath 			= serviceDeployConfig.FullPath + DIR_SEPARATOR + serviceDeployConfig.MasterDirName
	serviceDeployConfig.ClonePath 			= serviceDeployConfig.FullPath + DIR_SEPARATOR + serviceDeployConfig.GitCloneTmpDirName
	serviceDeployConfig.CodeGitClonePath 	= serviceDeployConfig.ClonePath + DIR_SEPARATOR + service.Name


	newServiceDeployConfig := serviceDeployConfig
	return newServiceDeployConfig,nil
}
//部署一个服务
func(cicdManager *CicdManager)DeployOneService(server Server , serviceDeployConfig ServiceDeployConfig ,  service Service)error{
	if service.Name != "Zgoframe"{//测试代码
		MyPrint("service name != Zgoframe")
		return nil
	}
	cicdManager.Option.Log.Info("DeployOneService:" + server.OutIp + " " + server.Env + " "+service.Name)
	//创建发布记录
	publish := cicdManager.Option.PublicManager.InsertOne(service,server)
	cicdManager.Option.Log.Info("create publish:"+ strconv.Itoa(publish.Id))
	//检查各种路径是否正确
	newServiceDeployConfig , err := cicdManager.DeployServiceCheck(server,serviceDeployConfig,service)
	if err !=nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	cicdManager.Option.Log.Info("DeployServiceCheck ok~")

	serviceDeployConfig = newServiceDeployConfig
	//step 1 : 项目代码及目录(git)相关
	newGitCodeDir , err := cicdManager.DeployOneServiceGitCode(serviceDeployConfig,service)
	if err != nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	//step 2 : 读取service项目代码里自带的cicd.toml ,供:后面使用
	serviceCICDConfig , err := cicdManager.DeployOneServiceCICIConfig(newGitCodeDir,serviceDeployConfig,server)
	if err != nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	//step 3: 生成该服务的，superVisor 配置文件
	err = cicdManager.DeployOneServiceSuperVisor(serviceDeployConfig,serviceCICDConfig)
	if err !=nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	//step 4: 处理项目自带的主配置文件
	err = cicdManager.DeployOneServiceProjectConfig(newGitCodeDir,server,serviceDeployConfig)
	if err !=nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	//step 5 : 先执行 服务自带的 shell 预处理
	_ , err = cicdManager.DeployOneServiceCommand(newGitCodeDir,serviceDeployConfig,serviceCICDConfig)
	if err != nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	//将master软链 指向 上面刚刚clone下的最新代码上
	err = cicdManager.DeployOneServiceLinkMaster(newGitCodeDir,serviceDeployConfig)
	if err != nil{
		return cicdManager.DeployOneServiceFailed(publish,err.Error())
	}
	cicdManager.Option.PublicManager.UpStatus(publish,2)
	ExitPrint("finish one.")
	return nil
}

func (cicdManager *CicdManager)DeployOneServiceGitCode(serviceDeployConfig ServiceDeployConfig,service Service)(string,error){
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

	shellArgc := service.Git + " " + serviceDeployConfig.ClonePath + " " +  service.Name
	//执行shell 脚本 后：service项目代码已被clone, git 版本号已知了
	gitLastCommitId ,err  := ExecShellFile(serviceDeployConfig.CICDShellFileName,shellArgc)
	if err != nil {
		return "",errors.New("ExecShellFile err:"+err.Error())
	}
	//cicdManager.Option.Log.Info("gitLastCommitId:" + gitLastCommitId)
	//刚刚clone完后，项目的目录
	//serviceCodeGitClonePath := serviceDeployConfig.ClonePath + DIR_SEPARATOR + service.Name
	//新刚刚克隆好的项目目录，移动一个新目录下，新目录名：git_master_versionId + 当前时间
	newGitCodeDir := serviceDeployConfig.FullPath + DIR_SEPARATOR + strconv.Itoa(GetNowTimeSecondToInt())  + "_" + gitLastCommitId
	cicdManager.Option.Log.Info(" service code move :" + serviceDeployConfig.CodeGitClonePath +" to "+ newGitCodeDir)
	//执行 移动操作
	err = os.Rename(serviceDeployConfig.CodeGitClonePath,newGitCodeDir)
	if err != nil{
		return newGitCodeDir,errors.New("serviceCodeGitClonePath os.Rename err:"+err.Error() )
	}
	cicdManager.Option.Log.Info("step 1 finish , newGitCodeDir :  "+newGitCodeDir + " , gitLastCommitId:"+gitLastCommitId)
	return newGitCodeDir, nil
}
func (cicdManager *CicdManager)DeployOneServiceCICIConfig(newGitCodeDir string,serviceDeployConfig ServiceDeployConfig,server Server)(ConfigServiceCICD,error){
	cicdManager.Option.Log.Info("step 2:load service CICD config ")
	//项目自带的CICD配置文件，这里有 服务启动脚本 和 依赖的环境
	serviceSelfCICDConf := newGitCodeDir + DIR_SEPARATOR + serviceDeployConfig.CICDConfFileName
	cicdManager.Option.Log.Info("read file:"+serviceSelfCICDConf)
	serviceCICDConfig := ConfigServiceCICD{}
	//读取项目自己的cicd配置文件，并映射到结构体中
	err := ReadConfFile(serviceSelfCICDConf,&serviceCICDConfig)
	if err != nil{
		return serviceCICDConfig,errors.New(err.Error())
	}
	serviceCICDConfig.System.Build  	= strings.Replace(serviceCICDConfig.System.Build,"#service_name#",serviceDeployConfig.Name,-1)
	serviceCICDConfig.System.Startup  	= strings.Replace(serviceCICDConfig.System.Startup,"#env#",server.Env,-1)
	PrintStruct(serviceCICDConfig,":")

	return serviceCICDConfig,nil
}
//生成该服务的，superVisor 配置文件
func (cicdManager *CicdManager)DeployOneServiceSuperVisor(serviceDeployConfig ServiceDeployConfig ,configServiceCICD  ConfigServiceCICD)error{
	cicdManager.Option.Log.Info("step 3 : create superVisor conf file.")
	superVisorOption := SuperVisorOption{
		ServiceName		: serviceDeployConfig.Name,
		ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
		ConfDir			: cicdManager.Option.Config.SuperVisor.ConfDir,
	}

	serviceSuperVisor,err := NewSuperVisor(superVisorOption)
	if err != nil{
		return err
	}
	//superVisor 配置文件中 动态的占位符，需要替换掉
	superVisorReplace := SuperVisorReplace{
		script_name				: serviceDeployConfig.Name,
		startup_script_command	: configServiceCICD.System.Startup,
		script_work_dir 		: serviceDeployConfig.MasterPath,
		stdout_logfile 			: serviceDeployConfig.BaseDir + DIR_SEPARATOR + "super_visor_stdout.log",
		stderr_logfile 			: serviceDeployConfig.BaseDir + DIR_SEPARATOR + "super_visor_stderr.log",
		process_name 			: serviceDeployConfig.Name,
	}
	//替换配置文件中的动态值，并生成配置文件
	serviceConfFileContent := serviceSuperVisor.ReplaceConfTemplate(superVisorReplace)
	//将已替换好的文件，生成一个新的配置文件
	err = serviceSuperVisor.CreateServiceConfFile(serviceConfFileContent)
	if err != nil{
		return err
	}

	return nil
}


func (cicdManager *CicdManager)DeployOneServiceProjectConfig(newGitCodeDir string,server Server, serviceDeployConfig ServiceDeployConfig)error{
	cicdManager.Option.Log.Info("step 4 : create project self conf file.")
	//读取该服务自己的配置文件 config.toml
	serviceSelfConfigTmpFileDir := newGitCodeDir + DIR_SEPARATOR + serviceDeployConfig.ConfigTmpFileName
	_ ,err  := FileExist(serviceSelfConfigTmpFileDir)
	if err != nil{
		return errors.New("serviceSelfConfigTmpFileDir CheckFileIsExist err:"+err.Error() )
	}
	cicdManager.Option.Log.Info("read file:"+serviceSelfConfigTmpFileDir)
	//读取模板文件内容
	serviceSelfConfigTmpFileContent,err := ReadString(serviceSelfConfigTmpFileDir)
	if err != nil{
		return errors.New(err.Error())
	}
	//开始替换 服务自己配置文件中的，实例信息，如：IP PORT
	serviceSelfConfigTmpFileContentNew := cicdManager.ReplaceInstance(serviceSelfConfigTmpFileContent,serviceDeployConfig.Name,server.Env)
	//生成新的配置文件
	newConfig := newGitCodeDir + DIR_SEPARATOR + serviceDeployConfig.ConfigFileName
	newConfigFile ,_:= os.Create(newConfig)
	newConfigFile.Write([]byte(serviceSelfConfigTmpFileContentNew))

	return nil
}


func (cicdManager *CicdManager)DeployOneServiceCommand(newGitCodeDir string,serviceDeployConfig ServiceDeployConfig,serviceCICDConfig ConfigServiceCICD)(output string ,err error){
	cicdManager.Option.Log.Info("step 5 : DeployOneServiceCommand.")
	cicdManager.Option.Log.Info("step 6.1 : project pre command "+serviceCICDConfig.System.Command)
	//    /usr/local/Cellar/go/1.16.5/bin/
	ExecShellCommandPre := "cd "+newGitCodeDir + " ; pwd ; "
	//ExecShellCommandPre := " ls -l "
	output1 := ""
	output2 := ""
	if serviceCICDConfig.System.Command != ""{
		output1,err = ExecShellCommand(ExecShellCommandPre + serviceCICDConfig.System.Command ,"")
		if err != nil{
			return output,errors.New("ExecShellCommand err "  +err.Error())
		}
		MyPrint(output)
	}
	//编译项目代码
	cicdManager.Option.Log.Info("step 6.2 : project build command "+serviceCICDConfig.System.Build)
	if serviceCICDConfig.System.Build != ""{
		output2,err = ExecShellCommand(ExecShellCommandPre + serviceCICDConfig.System.Build,"")
		if err != nil{
			return output,errors.New("ExecShellCommand err "  +err.Error())
		}
		MyPrint(output)
	}

	return output1 + " ssss " +output2,nil
	//cicdManager.Option.Log.Info("step 6.3 :  project testUnit command "+serviceCICDConfig.System.Command)
	//if serviceCICDConfig.System.TestUnit != ""{
	//	ExecShellCommand(serviceCICDConfig.System.TestUnit,"")
	//}
}

func (cicdManager *CicdManager)DeployOneServiceLinkMaster(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig)error{
	cicdManager.Option.Log.Info("step 6 : master dir softLink , os.Symlink:" + newGitCodeDir  +  " to " + serviceDeployConfig.MasterPath)
	_,err := PathExists(serviceDeployConfig.MasterPath)
	if err == nil{
		cicdManager.Option.Log.Info("master path exist , so need del ." + serviceDeployConfig.MasterPath)
		err = os.Remove(serviceDeployConfig.MasterPath)
		if err != nil{
			return errors.New("os.Remove " + serviceDeployConfig.MasterPath +  " err:" +err.Error())
		}
	}else if  os.IsNotExist(err){

	}else{
		//return cicdManager.DeployOneServiceFailed(publish,"unkonw err:"+err.Error())
		cicdManager.Option.Log.Info("master path exist , so need del ." + serviceDeployConfig.MasterPath)
		err = os.Remove(serviceDeployConfig.MasterPath)
		if err != nil{
			return errors.New("os.Remove " + serviceDeployConfig.MasterPath +  " err:" +err.Error())
		}
	}

	err = os.Symlink(newGitCodeDir, serviceDeployConfig.MasterPath)
	if err != nil{
		return errors.New("os.Symlink err :" + err.Error() )
	}
	return nil
}
//部署一个服务失败，统一处理接口
func (cicdManager *CicdManager)DeployOneServiceFailed(publish model.CICDPublish ,errMsg string)error{
	cicdManager.Option.PublicManager.UpStatus(publish,3)
	return cicdManager.MakeError(errMsg)
}