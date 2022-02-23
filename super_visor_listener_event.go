package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const RESP_OK = "RESULT 2\nOK"
const RESP_FAIL = "RESULT 4\nFAIL"

func main() {

	superVisorListenerEvent := NewSuperVisorListenerEvent()
	superVisorListenerEvent.Start()

}

type SuperVisorListenerEvent struct {

}

func NewSuperVisorListenerEvent()*SuperVisorListenerEvent{
	superVisorListenerEvent := new(SuperVisorListenerEvent)
	return superVisorListenerEvent
}

func(superVisorListenerEvent SuperVisorListenerEvent) Start(){
	//filename := "/tmp/test_super.log"
	//fd, err := os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
	//fmt.Println("OpenFile ",filename," err:",err)
	//
	//fmt.Println("start fro loop...")
	//for {
	//	// 发送后等待接收event
	//	num, err := stdout.WriteString("READY\n")
	//	fmt.Println("stdout.WriteString READY , num:",num,",err:",err)
	//	_ = stdout.Flush()
	//	// 接收header
	//	line, _, err := stdin.ReadLine()
	//	fmt.Println("stdin.ReadLin ,line:",string(line), ", err:",err)
	//	io.WriteString(fd,string(line))
	//
	//	stderr.WriteString("read" + string(line))
	//	stderr.Flush()
	//
	//	header, payloadSize := parseHeader(line)
	//	fmt.Println("parseHeader header:",header, " , payloadSize:",payloadSize)
	//	// 接收payload
	//	payload := make([]byte, payloadSize)
	//	stdin.Read(payload)
	//	stderr.WriteString("read : " + string(payload))
	//	stderr.Flush()
	//
	//	result := alarm(header, payload)
	//
	//	if result {   // 发送处理结果
	//		stdout.WriteString(RESP_OK)
	//	} else {
	//		stdout.WriteString(RESP_FAIL)
	//	}
	//	stdout.Flush()
	//}

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

		if result {   // 发送处理结果
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