package util

import (
	"errors"
	"github.com/abrander/go-supervisord"
	"os"
	"strings"
)

const (
	SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX = "service_"//进程启动时，启程名称的前缀，方便统一管理
)

type SuperVisorOption struct {
	Ip					string
	RpcPort 			string
	ConfTemplateFile 	string	//每个服务的superVisor 配置文件模型(需要后期替换占位符)
	ServiceName			string	//服务名称
	ConfDir 			string	//本机superVisor 的配置文件基目录(所有服务的superVisor配置文件均放在这个目录下面)
	ServiceNamePrefix 	string	//进程启动时，启程名称的前缀，方便统一管理
	Separator 			string
	//Port 				string
}


type SuperVisor struct {
	//Ip						string
	//RpcPort 				string
	//ConfTemplateFile 		string
	//ServiceName 			string
	//ConfDir 				string
	//Separator 				string
	//ServiceNamePrefix 		string
	Option 					SuperVisorOption
	ConfTemplateFileContent string
	Cli *supervisord.Client
}

func NewSuperVisor(superVisorOption SuperVisorOption )(*SuperVisor,error){
	superVisor := new(SuperVisor)
	//superVisor.Ip 			= superVisorOption.Ip
	//superVisor.RpcPort 		= superVisorOption.Port
	//superVisor.ConfDir 		= superVisorOption.ServiceName
	//superVisor.ServiceName 	= superVisorOption.ConfDir
	superVisorOption.Separator	= STR_SEPARATOR
	superVisorOption.ServiceNamePrefix = SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX
	superVisor.Option = superVisorOption


	superVisorConfTemplateFileContent ,err := ReadString(superVisorOption.ConfTemplateFile)
	if err != nil{
		return superVisor,errors.New("read superVisorConfTemplateFileContent err."+err.Error() + " , "+superVisorOption.ConfTemplateFile)
	}

	//superVisor.ConfTemplateFile = superVisorOption.ConfTemplateFile
	superVisor.ConfTemplateFileContent = superVisorConfTemplateFileContent

	return superVisor,nil
}
//通过XML Rpc 控制远程 superVisor 服务进程
func(superVisor *SuperVisor) InitXMLRpc()error{
	dns := "http://" + superVisor.Option.Ip + ":" + superVisor.Option.RpcPort + "/RPC2"
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
	key := superVisor.Option.Separator+"script_name"+superVisor.Option.Separator
	MyPrint(key)
	content = strings.Replace(content,key,replaceSource.script_name,-1)

	key = superVisor.Option.Separator+"startup_script_command"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.startup_script_command,-1)

	key = superVisor.Option.Separator+"script_work_dir"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.script_work_dir,-1)

	key = superVisor.Option.Separator+"stdout_logfile"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.stdout_logfile,-1)

	key = superVisor.Option.Separator+"stderr_logfile"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.stderr_logfile,-1)

	key = superVisor.Option.Separator+"process_name"+superVisor.Option.ServiceNamePrefix + superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.process_name,-1)

	return content
}

func(superVisor *SuperVisor)CreateServiceConfFile(content string)error{
	fileName := superVisor.Option.ConfDir +STR_SEPARATOR +  superVisor.Option.ServiceName + ".ini"
	file ,err := os.Create(fileName)
	MyPrint("os.Create:" ,fileName)
	if err!= nil{
		MyPrint("os.Create :",fileName , " err:",err)
		return err
	}

	file.Write([]byte(content))
	return nil
}