package util

import (
	"fmt"
	"github.com/abrander/go-supervisord"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//已知：依赖
//supervisor 依赖 python 依赖 xmlrpc
//git

type CicdPublish struct {
	Id				int
	RegTime 		int
	Status 			int
	ServiceName 	string
	Logs			[]string
	TotalExecTime 	int
	Server 			Server
}

type Server struct {
	Ip		string
	RegTime int
	From 	string //ali huwei tencent
	InnerIp string
	Env  	string
	Name 	string
}

type SuperVisor struct {
	Ip	string
	Port string
	Cli *supervisord.Client
}

type CicdManager struct {
	ServerList map[string]Server
	AppList map[int]App
}

func NewSuperVisor(ip string,port string)*SuperVisor{
	superVisor := new(SuperVisor)
	superVisor.Ip = ip
	superVisor.Port = port

	dns := "http://" + ip + ":" + port + "/RPC2"
	c, err := supervisord.NewClient(dns)
	if err != nil{

	}
	superVisor.Cli = c
	return superVisor
}

func NewCicdManager()*CicdManager{
	cicdManager := new(CicdManager)
	cicdManager.AppList = make(map[int]App)
	cicdManager.ServerList = make(map[string]Server)

	cicdManager.AppList[1] = App{
		Name: "zgoframe",
		Git: "git://github.com/mqzhifu/zgoframe.git",
	}

	cicdManager.ServerList["127.0.0.1"] = Server{
		Ip: "127.0.0.1",
	}

	return cicdManager
}

func(cicdManager *CicdManager)Init(){
	superVisorPort := "9001"
	serviceBaseDir := "/data/www/golang/testcicd/"
	serviceMasterPathName := "master"
	fmt.Println("superVisorPort:",superVisorPort)
	fmt.Println("serviceBaseDir:",serviceBaseDir)
	fmt.Println("serviceMasterPathName:",serviceMasterPathName)

	for _,server :=range cicdManager.ServerList{
		fmt.Println(server.Ip)
		//superVisor := NewSuperVisor(server.Ip,superVisorPort)
		for _,app :=range cicdManager.AppList{
			servicePath := serviceBaseDir + app.Name
			fmt.Println("servicePath:",servicePath)
			pathNotExistCreate(servicePath)

			//fmt.Println("serviceMasterPath:",serviceMasterPath)
			//serviceMasterPathExist ,_ := PathExists(serviceMasterPath)
			//fmt.Print(err)
			//if !serviceMasterPathExist {
			//	//创建一个目录
			//}
			//
			serviceGitClonePath := servicePath + "/" + "clone"
			pathNotExistCreate(serviceGitClonePath)
			gitLastCommitId :=GitCloneAndGetLastCommitIdByShell(serviceGitClonePath,app.Name,app.Git)
			fmt.Println("gitLastCommitId:",gitLastCommitId)
			hasGitClonePath := serviceGitClonePath + "/" + app.Name
			newGitCodeDir := servicePath + "/" + gitLastCommitId + "_" + strconv.Itoa(GetNowTimeSecondToInt())
			os.Rename(hasGitClonePath,newGitCodeDir)

			//serviceMasterPath := servicePath + "/" + serviceMasterPathName
			//err := os.Symlink(newGitCodeDir, serviceMasterPath)
			return
		}
	}
}

func pathNotExistCreate(path string){
	pathExist ,_ := PathExists(path)
	//fmt.Print(err)
	if !pathExist {
		//创建一个目录
		err := os.Mkdir(path, 0777)
		fmt.Println("create path:",path)
		if err != nil {
			fmt.Println("create path failed , err:",err)
		}
	}else{
		fmt.Println("path exist,",path)
	}
}

func GitCloneAndGetLastCommitIdByShell(serviceGitClonePath string,serviceName string,gitCloneUrl string)string{
	argc := gitCloneUrl + " " + serviceGitClonePath + " " +  serviceName

	shellFileName := "./cicd.sh" + " " + argc
	println(shellFileName)
	c := exec.Command("sh", "-c", shellFileName)

	output, err := c.CombinedOutput()
	if err != nil{
		fmt.Println("exec.Command err:",err)
	}
	outStr := string(output)
	outArr := strings.Split(outStr,"\n")

	return outArr[1]
	//fmt.Println(string(output), " err :",err)
	//var shellCommands []string
	//shellCommands = append(shellCommands,"./.sh")
	//return shellCommands
}
