package cicd

import (
	"encoding/json"
	"zgoframe/model"
	"zgoframe/util"
)

//记录一次部署的全过程
type DeployOneServiceFlowRecord struct {
	ServiceDeployConfig    ServiceDeployConfig `json:"service_deploy_config"`
	NewServiceDeployConfig ServiceDeployConfig `json:"new_service_deploy_config"`
	Server                 util.Server         `json:"-"` //mysql里有存这个字段
	service                model.Project       `json:"-"` //mysql里有存这个字段

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
