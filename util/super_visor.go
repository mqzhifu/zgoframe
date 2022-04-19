package util

import (
	"bytes"
	"errors"
	"github.com/abrander/go-supervisord"
	"os"
	"strings"
)

const (
	DIR_SEPARATOR = "/"
	STR_SEPARATOR = "#"

	SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX = "service_"//进程启动时，启程名称的前缀，方便统一管理

	//下面是错误标识，主要是给CICD部署时使用，最终给前端使用
	SV_ERROR_NONE = 0
	SV_ERROR_INIT = 1
	SV_ERROR_CONN = 2
	SV_ERROR_NOT_FOUND = 3
)

type SuperVisorReplace struct {
	Script_name            string
	Startup_script_command string
	Script_work_dir        string
	Stdout_logfile         string
	Stderr_logfile         string
	Process_name           string
}

//==============superVisor===========

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
	MyPrint(dns)
	c, err := supervisord.NewClient(dns)
	if err != nil{
		MyPrint("superVisor init err:",err)
		return err
	}
	state ,err := c.GetState()
	MyPrint("superVisor: XMLRpc state:",state," err:",err)
	if err != nil{
		MyPrint("superVisor err:"+err.Error())
	}

	superVisor.Cli = c
	return nil
}

func(superVisor *SuperVisor)StartProcess(serviceName string,wait bool)error{
	//processListInfo  ,err := superVisor.Cli.GetAllProcessInfo()
	//if err != nil{
	//
	//}
	//
	//var serviceSuperVisorInfo supervisord.ProcessInfo
	//hasSearch := false
	//for _,v := range processListInfo{
	//	if serviceName == v.Name{
	//		serviceSuperVisorInfo = v
	//		break
	//	}
	//}
	//if !hasSearch{
	//
	//}
	//
	//if serviceSuperVisorInfo.State !=  supervisord.StateStopped && serviceSuperVisorInfo.State == supervisord.StateExited {
	//
	//}
	return superVisor.Cli.StartProcess(serviceName,wait)

}

func(superVisor *SuperVisor)StopProcess(){

}

func(superVisor *SuperVisor)ReloadProcess(){

}



func(superVisor *SuperVisor)ReplaceConfTemplate(replaceSource SuperVisorReplace)string{
	content := superVisor.ConfTemplateFileContent
	key := superVisor.Option.Separator+"script_name"+superVisor.Option.Separator
	MyPrint(key)
	content = strings.Replace(content,key,replaceSource.Script_name,-1)

	key = superVisor.Option.Separator+"startup_script_command"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.Startup_script_command,-1)

	key = superVisor.Option.Separator+"script_work_dir"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.Script_work_dir,-1)

	key = superVisor.Option.Separator+"stdout_logfile"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.Stdout_logfile,-1)

	key = superVisor.Option.Separator+"stderr_logfile"+superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.Stderr_logfile,-1)

	key = superVisor.Option.Separator+"process_name"+superVisor.Option.ServiceNamePrefix + superVisor.Option.Separator
	content = strings.Replace(content,key,replaceSource.Process_name,-1)

	//ExitPrint(content)

	return content
}

func(superVisor *SuperVisor)CreateServiceConfFile(content string)error{
	fileName := superVisor.Option.ConfDir +DIR_SEPARATOR +  superVisor.Option.ServiceName + ".ini"
	file ,err := os.Create(fileName)
	MyPrint("os.Create:" ,fileName)
	if err!= nil{
		MyPrint("os.Create :",fileName , " err:",err)
		return err
	}
	contentByte := bytes.Trim([]byte(content),"\x00")//NUL
	file.Write(contentByte)
	file.Close()
	return nil
}