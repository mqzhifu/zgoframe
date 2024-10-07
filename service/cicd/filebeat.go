package cicd

import (
	"os"
	"strconv"
	"strings"
	"zgoframe/util"
)

func (cicdManager *CicdManager) GenerateAllFilebeat() {
	pwd, _ := os.Getwd() //当前路径
	opDirFull := pwd + "/" + cicdManager.Option.OpDirName

	for _, server := range cicdManager.Option.ServerList {
		cicdManager.GenerateFilebeat(server, opDirFull)
		util.MyPrint("finish one ...........")
	}
	util.ExitPrint(33)
}
func (cicdManager *CicdManager) GenerateFilebeat(server util.Server, opDir string) {

	instance, empty := cicdManager.Option.InstanceManager.GetByEnvName(server.Env, "es")
	if empty {
		util.ExitPrint("ProcessFilebeat GetByEnvName es empty :" + strconv.Itoa(server.Env))
	}

	esDns := instance.Host + ":" + instance.Port
	filebeatConfigFile := opDir + "/" + "filebeat.yaml"
	filebeatConfigFileContent, _ := util.ReadString(filebeatConfigFile)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#elasticsearch_output_hosts#", esDns, -1)

	filebeatInput := ""
	for _, service := range cicdManager.Option.ServiceList {
		filebeat_input_file := opDir + "/" + "filebeat_input.yaml"
		filebeat_input_content, _ := util.ReadString(filebeat_input_file)
		serviceLogDir := cicdManager.Option.Config.System.LogDir + "/" + service.Name + "/*.log"
		filebeat_input_content = strings.Replace(filebeat_input_content, "#paths#", serviceLogDir, -1)
		filebeat_input_content = strings.Replace(filebeat_input_content, "#source#", service.Name, -1)

		filebeatInput += filebeat_input_content + "\n"
	}
	esOutput := ""
	for _, service := range cicdManager.Option.ServiceList {
		esOutputFile := opDir + "/" + "filebeat_es_output.yaml"
		esOutputFileContent, _ := util.ReadString(esOutputFile)
		esOutputFileContent = strings.Replace(esOutputFileContent, "#index#", service.Name, -1)

		esOutput += esOutputFileContent + "\n"
	}
	//MyPrint(filebeatInput)
	//MyPrint(esOutput)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#filebeat_inputs#", filebeatInput, -1)
	filebeatConfigFileContent = strings.Replace(filebeatConfigFileContent, "#elasticsearch_output_index#", esOutput, -1)

	util.MyPrint(filebeatConfigFileContent)
}
