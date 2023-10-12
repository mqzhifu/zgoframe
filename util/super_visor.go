package util

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/abrander/go-supervisord"
	"os"
	"strconv"
	"strings"
)

type SuperVisorReplace struct {
	ScriptName           string
	StartupScriptCommand string
	ScriptWorkDir        string
	StdoutLogfile        string
	StderrLogfile        string
	ProcessName          string
}

// ==============superVisor===========

type SuperVisorOption struct {
	Ip                string
	RpcPort           string
	Username          string
	Password          string
	ConfTemplateFile  string // 每个服务的superVisor 配置文件模型(需要后期替换占位符)
	ServiceName       string // 服务名称
	ConfDir           string // 本机superVisor 的配置文件基目录(所有服务的superVisor配置文件均放在这个目录下面)
	ServiceNamePrefix string // 进程启动时，启程名称的前缀，方便统一管理
	Separator         string
	// Port 				string
}

type SuperVisor struct {
	Option                  SuperVisorOption
	ConfTemplateFileContent string
	Cli                     *supervisord.Client
}

func NewSuperVisor(superVisorOption SuperVisorOption) (*SuperVisor, error) {
	superVisor := new(SuperVisor)
	superVisorOption.Separator = STR_SEPARATOR
	superVisorOption.ServiceNamePrefix = SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX
	superVisor.Option = superVisorOption

	return superVisor, nil
}

// 通过XML Rpc 控制远程 superVisor 服务进程
func (superVisor *SuperVisor) InitXMLRpc() error {
	if superVisor.Option.Ip == "" || superVisor.Option.RpcPort == "" {
		return errors.New("ip or port empty")
	}

	dns := "http://" + superVisor.Option.Ip + ":" + superVisor.Option.RpcPort + "/RPC2"
	MyPrint("InitXMLRpc: " + dns)
	var err error
	var c *supervisord.Client
	if superVisor.Option.Username != "" && superVisor.Option.Password != "" {
		c, err = supervisord.NewClient(dns, supervisord.WithAuthentication(superVisor.Option.Username, superVisor.Option.Password))
	} else {
		c, err = supervisord.NewClient(dns)
	}
	MyPrint()

	if err != nil {
		MyPrint("superVisor init err:", err)
		return err
	}
	_, err = c.GetState()
	// state ,err := c.GetState()
	// MyPrint("superVisor: XMLRpc state:",state," err:",err)
	if err != nil {
		MyPrint("superVisor err:" + err.Error())
		return errors.New("superVisor err:" + err.Error())
	}

	superVisor.Cli = c
	return nil
}

func (superVisor *SuperVisor) StartProcess(serviceName string, wait bool) error {

	return superVisor.Cli.StartProcess(serviceName, wait)

}

func (superVisor *SuperVisor) StopProcess(serviceName string, wait bool) error {
	processName := superVisor.Option.ServiceNamePrefix + serviceName
	return superVisor.Cli.StopProcess(processName, wait)
}

func (superVisor *SuperVisor) ReloadProcess() {

}
func (superVisor *SuperVisor) SetConfTemplateFile(ConfTemplateFile string) error {
	superVisorConfTemplateFileContent, err := ReadString(ConfTemplateFile)
	if err != nil {
		return errors.New("read superVisorConfTemplateFileContent err." + err.Error() + " , " + ConfTemplateFile)
	}
	superVisor.ConfTemplateFileContent = superVisorConfTemplateFileContent
	return nil
}

func (superVisor *SuperVisor) ReplaceConfTemplate(replaceSource SuperVisorReplace) (string, error) {
	if superVisor.ConfTemplateFileContent == "" {
		return "", errors.New("ConfTemplateFileContent empty")
	}
	content := superVisor.ConfTemplateFileContent
	// MyPrint(content)
	key := superVisor.Option.Separator + "script_name" + superVisor.Option.Separator
	// MyPrint(key)
	content = strings.Replace(content, key, replaceSource.ScriptName, -1)

	key = superVisor.Option.Separator + "startup_script_command" + superVisor.Option.Separator
	content = strings.Replace(content, key, replaceSource.StartupScriptCommand, -1)

	key = superVisor.Option.Separator + "script_work_dir" + superVisor.Option.Separator
	content = strings.Replace(content, key, replaceSource.ScriptWorkDir, -1)

	key = superVisor.Option.Separator + "stdout_logfile" + superVisor.Option.Separator
	content = strings.Replace(content, key, replaceSource.StdoutLogfile, -1)

	key = superVisor.Option.Separator + "stderr_logfile" + superVisor.Option.Separator
	content = strings.Replace(content, key, replaceSource.StderrLogfile, -1)

	key = superVisor.Option.Separator + "process_name" + superVisor.Option.Separator
	content = strings.Replace(content, key, superVisor.Option.ServiceNamePrefix+replaceSource.ProcessName, -1)

	// ExitPrint(content)

	return content, nil
}

// ========监听相关
func (superVisor *SuperVisor) CreateServiceConfFile(content string, newGitCodeDir string) error {
	// 本机部署时：直接将配置文件转到superVisor目录下，立即生效
	// fileName := superVisor.Option.ConfDir +DIR_SEPARATOR +  superVisor.Option.ServiceName + ".ini"
	// 远程部署：是先在本地部署，再推送到远端，所以，是是先将配置文件生成到代码目录下，最后再同步过去
	fileName := newGitCodeDir + "/" + superVisor.Option.ServiceName + ".ini"
	// MyPrint(fileName)
	file, err := os.Create(fileName)
	MyPrint("os.Create:", fileName)
	if err != nil {
		MyPrint("os.Create :", fileName, " err:", err)
		return err
	}
	contentByte := bytes.Trim([]byte(content), "\x00") // NUL
	file.Write(contentByte)
	file.Close()

	// ExitPrint(-1)

	return nil
}

const RESP_OK = "RESULT 2\nOK"
const RESP_FAIL = "RESULT 4\nFAIL"

func (superVisor *SuperVisor) ListenerEvent() {
	stdin := bufio.NewReader(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	for {
		// 发送后等待接收event
		_, _ = stdout.WriteString("READY\n")
		_ = stdout.Flush()
		// 接收header
		line, _, _ := stdin.ReadLine()
		stderr.WriteString("stdin ReadLine: " + string(line))
		stderr.Flush()

		header, payloadSize := parseHeader(line)

		// 接收payload
		payload := make([]byte, payloadSize)
		stdin.Read(payload)
		stderr.WriteString(" , stdin Read , payload : " + string(payload) + "\n")
		stderr.Flush()

		result := alarm(header, payload)

		if result { // 发送处理结果
			stdout.WriteString(RESP_OK)
		} else {
			stdout.WriteString(RESP_FAIL)
		}
		stdout.Flush()
	}
}

func parseHeader(data []byte) (header map[string]string, payloadSize int) {
	pairs := strings.Split(string(data), " ")
	header = make(map[string]string, len(pairs))

	for _, pair := range pairs {
		token := strings.Split(pair, ":")
		header[token[0]] = token[1]
	}

	payloadSize, _ = strconv.Atoi(header["len"])
	return header, payloadSize
}

// 这里设置报警即可
func alarm(header map[string]string, payload []byte) bool {
	// send mail
	return true
}
