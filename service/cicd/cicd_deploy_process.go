package cicd

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	"zgoframe/model"
	"zgoframe/util"
)

// step 1
func (deploy *Deploy) DeployServiceCheck(serviceDeployConfig ServiceDeployConfig, service model.Project, server util.Server) (ServiceDeployConfig, error) {
	deploy.Option.Log.Info("step 1 : DeployServiceCheck ")
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
	// 本机部分编译，要把远程部署多出一层： 服务器IP目录->服务目录
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
	// baseDir 已由 构造函数做校验了

	serviceDeployConfig.Name = service.Name
	serviceDeployConfig.FullPath = serviceDeployConfig.BaseDir + util.DIR_SEPARATOR + serviceDeployConfig.Name
	serviceDeployConfig.MasterPath = serviceDeployConfig.FullPath + util.DIR_SEPARATOR + serviceDeployConfig.MasterDirName
	serviceDeployConfig.ClonePath = serviceDeployConfig.FullPath + util.DIR_SEPARATOR + serviceDeployConfig.GitCloneTmpDirName
	serviceDeployConfig.CodeGitClonePath = serviceDeployConfig.ClonePath + util.DIR_SEPARATOR + service.Name
	serviceDeployConfig.FullOpDirName = deploy.Option.Config.System.RootDir + "/" + serviceDeployConfig.OpDirName

	// serviceDeployConfig.RemoteBaseDir = serviceDeployConfig.RemoteBaseDir
	newServiceDeployConfig := serviceDeployConfig

	// util.PrintStruct(newServiceDeployConfig, ":")

	return newServiceDeployConfig, nil
}

// step 2
func (deploy *Deploy) DeployOneServiceGitCode(serviceDeployConfig ServiceDeployConfig, service model.Project) (string, string, string, error) {
	deploy.Option.Log.Info("step 2 : git clone project code and get git commit id.")
	// FullPath 一个服务的根目录，大部分操作都在这个目录下(除了superVisor)
	// 查看服务根目录是否存在，不存在即新创建
	pathNotExistCreate(serviceDeployConfig.FullPath)
	// 查看git clone 目录是否存在，不存在即新创建
	pathNotExistCreate(serviceDeployConfig.ClonePath)
	// 构建 shell 执行时所需 参数
	shellArgc := service.Git + " " + serviceDeployConfig.ClonePath + " " + service.Name + " " + deploy.Option.Config.System.RemoteUploadDir + " " + deploy.Option.UploadDiskPath + " " + deploy.Option.Config.System.RemoteDownloadDir + " " + deploy.Option.DownloadDiskPath
	CICDShellFileName := ""
	// 执行shell 脚本 后：service项目代码已被clone, git 版本号已知了
	if service.Type == model.PROJECT_TYPE_FE {
		CICDShellFileName = "cicd_fe.sh"
	} else {
		CICDShellFileName = serviceDeployConfig.CICDShellFileName
	}

	gitLastCommitId, err := ExecShellFile(serviceDeployConfig.FullOpDirName+"/"+CICDShellFileName, shellArgc)
	if err != nil {
		return "", "", "", errors.New("ExecShellFile err:" + err.Error())
	}
	// 新刚刚克隆好的项目目录，移动一个新目录下，新目录名：git_master_versionId + 当前时间
	projectDirName := strconv.Itoa(util.GetNowTimeSecondToInt()) + "_" + gitLastCommitId
	newGitCodeDir := serviceDeployConfig.FullPath + util.DIR_SEPARATOR + projectDirName
	deploy.Option.Log.Info(" service code move :" + serviceDeployConfig.CodeGitClonePath + " to " + newGitCodeDir)
	// 执行 移动操作
	err = os.Rename(serviceDeployConfig.CodeGitClonePath, newGitCodeDir)
	if err != nil {
		return newGitCodeDir, "", "", errors.New("serviceCodeGitClonePath os.Rename err:" + err.Error())
	}
	deploy.Option.Log.Info("step 2 finish , newGitCodeDir :  " + newGitCodeDir + " , gitLastCommitId:" + gitLastCommitId)

	// 处理图片目录 的软件 连接
	// _, err := util.FileExist(cicdManager.Option.UploadDiskPath)
	// cicdManager.Option.Log.Info("ln -s " + cicdManager.Option.Config.System.RemoteUploadDir + " " + cicdManager.Option.UploadDiskPath)
	// err = os.Symlink(cicdManager.Option.Config.System.RemoteUploadDir,cicdManager.Option.UploadDiskPath)
	// if err != nil{
	//	return newGitCodeDir , projectDirName ,gitLastCommitId, errors.New("link file upload err:" + err.Error())
	// }
	// util.ExitPrint(33)

	return newGitCodeDir, projectDirName, gitLastCommitId, nil
}

// step 3
func (deploy *Deploy) DeployOneServiceCICIConfig(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig, server util.Server, gitLastCommitId string) (ConfigServiceCICD, string, error) {
	deploy.Option.Log.Info("step 3:load service CICD config ")
	// 项目自带的CICD配置文件，这里有 服务启动脚本 和 依赖的环境
	serviceSelfCICDConf := newGitCodeDir + util.DIR_SEPARATOR + serviceDeployConfig.CICDConfFileName
	deploy.Option.Log.Info("read file:" + serviceSelfCICDConf)
	serviceCICDConfig := ConfigServiceCICD{}
	// 读取项目自己的cicd配置文件，并映射到结构体中
	err := util.ReadConfFileAutoExt(serviceSelfCICDConf, &serviceCICDConfig)
	if err != nil {
		return serviceCICDConfig, serviceSelfCICDConf, errors.New(err.Error())
	}
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#service_name#", serviceDeployConfig.Name, -1)
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#datetime#", strconv.Itoa(util.GetNowTimeSecondToInt()), -1)
	serviceCICDConfig.System.Build = strings.Replace(serviceCICDConfig.System.Build, "#git_version#", gitLastCommitId, -1)
	// util.MyPrint(serviceCICDConfig.System.Build)
	// util.ExitPrint(33)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#env#", strconv.Itoa(server.Env), -1)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#master_path#", serviceDeployConfig.RemoteBaseDir+"/"+serviceDeployConfig.Name+"/"+serviceDeployConfig.MasterDirName, -1)
	serviceCICDConfig.System.Startup = strings.Replace(serviceCICDConfig.System.Startup, "#service_name#", serviceDeployConfig.Name, -1)

	// util.ExitPrint(serviceCICDConfig.System.Startup)
	// util.PrintStruct(serviceCICDConfig, ":")

	return serviceCICDConfig, serviceSelfCICDConf, nil
}

// step 4 生成该服务的，superVisor 配置文件
func (deploy *Deploy) DeployOneServiceSuperVisor(serviceDeployConfig ServiceDeployConfig, configServiceCICD ConfigServiceCICD, newGitCodeDir string) error {
	deploy.Option.Log.Info("step 4 : create superVisor conf file.")
	superVisorOption := util.SuperVisorOption{
		ConfDir:     deploy.Option.Config.SuperVisor.ConfDir,
		ServiceName: serviceDeployConfig.Name,
		// ConfTemplateFile: cicdManager.Option.Config.SuperVisor.ConfTemplateFile,
	}

	serviceSuperVisor, err := util.NewSuperVisor(superVisorOption)
	if err != nil {
		return err
	}
	serviceSuperVisor.SetConfTemplateFile(deploy.Option.Config.SuperVisor.ConfTemplateFile)
	// superVisor 配置文件中 动态的占位符，需要替换掉
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

	// 替换配置文件中的动态值，并生成配置文件
	serviceConfFileContent, _ := serviceSuperVisor.ReplaceConfTemplate(superVisorReplace)
	// 将已替换好的文件，生成一个新的配置文件
	err = serviceSuperVisor.CreateServiceConfFile(serviceConfFileContent, newGitCodeDir)
	if err != nil {
		return err
	}

	return nil
}

// step 5
func (deploy *Deploy) DeployOneServiceProjectConfig(newGitCodeDir string, server util.Server, serviceDeployConfig ServiceDeployConfig, configServiceCICD ConfigServiceCICD, service model.Project) (string, string, error) {
	deploy.Option.Log.Info("step 5 : create project self conf file.")
	// 读取该服务自己的配置文件 config.toml
	// serviceSelfConfigTmpFileDir := newGitCodeDir + util.DIR_SEPARATOR + configServiceCICD.System.ConfigTmpFileName
	// 原 config.toml 是放在项目根目录下，后期做docker makefile 的时候遇到把，新建了个 config 文件夹存在
	serviceSelfConfigTmpFileDir := newGitCodeDir + util.DIR_SEPARATOR + "config" + util.DIR_SEPARATOR + configServiceCICD.System.ConfigTmpFileName
	_, err := util.FileExist(serviceSelfConfigTmpFileDir)
	if err != nil {
		return "", "", errors.New("serviceSelfConfigTmpFileDir CheckFileIsExist err:" + err.Error())
	}
	deploy.Option.Log.Info("read file:" + serviceSelfConfigTmpFileDir)
	// 读取模板文件内容
	serviceSelfConfigTmpFileContent, err := util.ReadString(serviceSelfConfigTmpFileDir)
	if err != nil {
		return "", "", errors.New(err.Error())
	}
	// 开始替换 服务自己配置文件中的，实例信息，如：IP PORT
	serviceSelfConfigTmpFileContentNew := deploy.ReplaceInstance(serviceSelfConfigTmpFileContent, serviceDeployConfig.Name, server.Env, service.Id)

	key := util.STR_SEPARATOR + "projectId" + util.STR_SEPARATOR
	serviceSelfConfigTmpFileContentNew = strings.Replace(serviceSelfConfigTmpFileContentNew, key, strconv.Itoa(service.Id), -1)

	// 生成新的配置文件
	newConfig := newGitCodeDir + util.DIR_SEPARATOR + configServiceCICD.System.ConfigFileName
	newConfigFile, _ := os.Create(newConfig)
	contentByte := bytes.Trim([]byte(serviceSelfConfigTmpFileContentNew), "\x00") // NUL
	newConfigFile.Write(contentByte)

	return serviceSelfConfigTmpFileDir, newConfig, nil
}

// step 6
func (deploy *Deploy) DeployOneServiceCommand(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig, serviceCICDConfig ConfigServiceCICD) (command string, build string, output string, err error) {
	deploy.Option.Log.Info("step 6 : DeployOneServiceCommand.")
	ExecShellCommandPre := "cd " + newGitCodeDir + "  ; pwd ; "
	output1 := ""
	output2 := ""
	if serviceCICDConfig.System.Command != "" {
		command = ExecShellCommandPre + serviceCICDConfig.System.Command
		deploy.Option.Log.Info("step 6.1 : System.Command: " + command)
		output1, err = ExecShellCommand(ExecShellCommandPre+serviceCICDConfig.System.Command, "")
		if err != nil {
			return command, build, output, errors.New("ExecShellCommand " + command + " err " + err.Error())
		}
	}
	// 编译项目代码
	if serviceCICDConfig.System.Build != "" {
		build = ExecShellCommandPre + serviceCICDConfig.System.Build
		deploy.Option.Log.Info("step 6.2 : project build command :" + build)
		output2, err = ExecShellCommand(build, "")
		if err != nil {
			return command, build, output, errors.New("ExecShellCommand " + command + "  err " + err.Error())
		}
	}

	return command, build, output1 + " <br/> " + output2, nil
}

func (deploy *Deploy) GetRsyncInstance(server util.Server) (is util.Instance, err error) {
	rsyncInstance, empty := deploy.Option.InstanceManager.GetByEnvName(server.Env, "rsync")
	if empty {
		deploy.Option.Log.Error("GetByEnvName rsync empty: env=" + strconv.Itoa(server.Env) + " rsync")
		return is, errors.New("SyncOneServiceToRemote err1")
	}

	if rsyncInstance.Host == "" || rsyncInstance.Port == "" || rsyncInstance.User == "" || rsyncInstance.Ext == "" {
		deploy.Option.Log.Error("rsyncInstance someone empty : rsyncHost rsyncPort rsyncUserName rsyncModule")
		return is, errors.New("SyncOneServiceToRemote err2")
	}

	return is, nil
}

func GetRsyncCommand(remoteHost string, port string, username string, ps string, module string, exclude string, localPath string) string {
	comm := ""
	if exclude == "" {
		comm = "rsync -avz --progress --port=" + port + " " + localPath + " " + username + "@" + remoteHost + "::" + module
	} else {
		comm = "rsync -avz --progress --port=" + port + " --exclude=" + exclude + " " + localPath + " " + username + "@" + remoteHost + "::" + module
	}
	return comm
}

// 本机部署均已完成，需要将本地代码同步到远端
func (deploy *Deploy) SyncOneServiceToRemote(serviceDeployConfig ServiceDeployConfig, server util.Server, newGitCodeDir string, project model.Project) (syncCodeShellCommand string, syncSuperVisorShellCommand string, err error) {
	rsyncInstance, err := deploy.GetRsyncInstance(server)
	if err != nil {
		return "", "", err
	}
	if project.Type == model.PROJECT_TYPE_SERVICE {
		// 1 同步代码
		// syncCodeShellCommand = GetRsyncCommandPre() + " --exclude=master " + serviceDeployConfig.FullPath + " rsync@" + server.OutIp + "::www"
		syncCodeShellCommand = GetRsyncCommand(rsyncInstance.Host, rsyncInstance.Port, rsyncInstance.User, "", rsyncInstance.Ext, "master", serviceDeployConfig.FullPath)
		util.ExitPrint(syncCodeShellCommand)
		_, err := ExecShellCommand(syncCodeShellCommand, "")
		util.MyPrint("SyncOneServiceToRemote:", syncCodeShellCommand, " err:", err)
		// 2 同步superVisor
		// syncSuperVisorShellCommand = GetRsyncCommandPre() + newGitCodeDir + "/" + serviceDeployConfig.Name + ".ini" + " rsync@" + server.OutIp + "::super_visor"
		syncCodeShellCommand = GetRsyncCommand(rsyncInstance.Host, rsyncInstance.Port, rsyncInstance.User, "", "super_visor", "", newGitCodeDir+"/"+serviceDeployConfig.Name+".ini")
		_, err = ExecShellCommand(syncSuperVisorShellCommand, "")
		util.MyPrint("syncSuperVisorShellCommand:", syncSuperVisorShellCommand, " err:", err)
	} else if project.Type == model.PROJECT_TYPE_FE {
		// util.MyPrint(serviceDeployConfig)
		// syncCodeShellCommand = GetRsyncCommandPre() + " --exclude=node_modules " + newGitCodeDir + " rsync@" + server.OutIp + "::www/" + serviceDeployConfig.Name
		syncCodeShellCommand = GetRsyncCommand(rsyncInstance.Host, rsyncInstance.Port, rsyncInstance.User, "", rsyncInstance.Ext+"/"+serviceDeployConfig.Name, "node_modules", newGitCodeDir)
		// util.ExitPrint(syncCodeShellCommand)
		_, err := ExecShellCommand(syncCodeShellCommand, "")
		util.MyPrint("SyncOneServiceToRemote:", syncCodeShellCommand, " err:", err)
	} else {
		return "", "", errors.New("SyncOneServiceToRemote :project type err.")
	}

	return syncCodeShellCommand, syncSuperVisorShellCommand, nil
}

// step 8
func (deploy *Deploy) DeployOneServiceLinkMaster(newGitCodeDir string, serviceDeployConfig ServiceDeployConfig) error {
	deploy.Option.Log.Info("step 8 : master dir softLink , os.Symlink:" + newGitCodeDir + " to " + serviceDeployConfig.MasterPath)
	_, err := util.PathExists(serviceDeployConfig.MasterPath)
	if err == nil {
		deploy.Option.Log.Info("master path exist , so need del ." + serviceDeployConfig.MasterPath)
		err = os.Remove(serviceDeployConfig.MasterPath)
		if err != nil {
			return errors.New("os.Remove " + serviceDeployConfig.MasterPath + " err:" + err.Error())
		}
	} else if os.IsNotExist(err) {

	} else {
		// return cicdManager.DeployOneServiceFailed(publish,"unkonw err:"+err.Error())
		deploy.Option.Log.Info("master path exist , so need del ." + serviceDeployConfig.MasterPath)
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
