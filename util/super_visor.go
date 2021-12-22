package util

import (
	"github.com/abrander/go-supervisord"
	"strings"
)

const (
	SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX = "service_"
)

type SuperVisorOption struct {
	Ip	string
	Port string
	ConfTemplateFile string
	ServiceName	string
	ConfDir string
	ServiceNamePrefix string
}


type SuperVisor struct {
	Ip	string
	RpcPort string
	ConfTemplateFile string
	ConfTemplateFileContent string
	ServiceName string
	ConfDir string
	Separator string
	ServiceNamePrefix string
	Cli *supervisord.Client
}

func NewSuperVisor(superVisorOption SuperVisorOption )*SuperVisor{
	superVisor := new(SuperVisor)
	superVisor.Ip 		= superVisorOption.Ip
	superVisor.RpcPort 	= superVisorOption.Port

	superVisor.ConfDir 	= superVisorOption.ServiceName
	superVisor.ServiceName = superVisorOption.ConfDir
	superVisor.Separator= STR_SEPARATOR
	superVisor.ServiceNamePrefix = SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX
	
	superVisorConfTemplateFileContent ,err := ReadString(superVisorOption.ConfTemplateFile)
	if err != nil{
		ExitPrint("read superVisorConfTemplateFileContent err.")
	}

	superVisor.ConfTemplateFile = superVisorOption.ConfTemplateFile
	superVisor.ConfTemplateFileContent = superVisorConfTemplateFileContent

	return superVisor
}

func(superVisor *SuperVisor) InitXMLRpc()error{
	dns := "http://" + superVisor.Ip + ":" + superVisor.RpcPort + "/RPC2"
	c, err := supervisord.NewClient(dns)
	if err != nil{
		MyPrint("superVisor init err:",err)
		return err
	}
	superVisor.Cli = c
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

	key = superVisor.Separator+"process_name"+superVisor.ServiceNamePrefix + superVisor.Separator
	content = strings.Replace(content,key,replaceSource.process_name,-1)

	return content
}