package util

import (
	"errors"
	"net"
	"time"
	"fmt"
	"strconv"
	"os"
)

type PingOption struct{
	Count int
	Size int
	Timeout int64
	Nerverstop bool
}

func NewPingOption()*PingOption{
	return &PingOption{
		Count:4,
		Size:32,
		Timeout:1000,
		Nerverstop:false,
	}
}

//ping连接用的协议是ICMP，原理：
//Ping的基本原理是发送和接受ICMP请求回显报文。接收方将报文原封不动的返回发送方，发送方校验报文，校验成功则表示ping通。
//一台主机向一个节点发送一个类型字段值为8的ICMP报文，如果途中没有异常（如果没有被路由丢弃，目标不回应ICMP或者传输失败），
//则目标返回类型字段值为0的ICMP报文，说明这台主机可达
func (p *PingOption)Ping3(host string, args map[string]interface{})error {
	//要发送的回显请求数
	var count int = 1
	//要发送缓冲区大小,单位：字节
	var size int = 32
	//等待每次回复的超时时间(毫秒)
	var timeout int64 = 1000
	//Ping 指定的主机，直到停止
	var neverstop bool = false
	fmt.Println(args,"args")
	if len(args)!=0{
		count = args["n"].(int)
		size = args["l"].(int)
		timeout = args["w"].(int64)
		neverstop = args["t"].(bool)
	}

	//查找规范的dns主机名字  eg.www.baidu.com->www.a.shifen.com
	cname, _ := net.LookupCNAME(host)
	starttime := time.Now()
	//此处的链接conn只是为了获得ip := conn.RemoteAddr(),显示出来，因为后面每次连接都会重新获取conn,todo 但是每次重新获取的conn,其连接的ip保证一致么？
	conn, err := net.DialTimeout("ip4:icmp", host, time.Duration(timeout*1000*1000))
	//每个域名可能对应多个ip，但实际连接时，请求只会转发到某一个上，故需要获取实际连接的远程ip，才能知道实际ping的机器是哪台
	//  ip := conn.RemoteAddr()
	//  fmt.Println("正在 Ping " + cname + " [" + ip.String() + "] 具有 32 字节的数据:")

	var seq int16 = 1
	id0, id1 := genidentifier3(host)
	//ICMP报头的长度至少8字节，如果报文包含数据部分则大于8字节。
	//ping命令包含"请求"（Echo Request，报头类型是8）和"应答"（Echo Reply，类型是0）2个部分，由ICMP报头的类型决定
	const ECHO_REQUEST_HEAD_LEN = 8

	//记录发送次数
	sendN := 0
	//成功应答次数
	recvN := 0
	//记录失败请求数
	lostN := 0
	//所有请求中应答时间最短的一个
	shortT := -1
	//所有请求中应答时间最长的一个
	longT := -1
	//所有请求的应答时间和
	sumT := 0

	for count > 0 || neverstop {
		sendN++
		//ICMP报文长度，报头8字节，数据部分32字节
		var msg []byte = make([]byte, size+ECHO_REQUEST_HEAD_LEN)
		//第一个字节表示报文类型，8表示回显请求
		msg[0] = 8                        // echo
		//ping的请求和应答，该code都为0
		msg[1] = 0                        // code 0
		//校验码占2字节
		msg[2] = 0                        // checksum
		msg[3] = 0                        // checksum
		//ID标识符 占2字节
		msg[4], msg[5] = id0, id1         //identifier[0] identifier[1]
		//序号占2字节
		msg[6], msg[7] = gensequence3(seq) //sequence[0], sequence[1]

		length := size + ECHO_REQUEST_HEAD_LEN
		//计算检验和。
		check := checkSum3(msg[0:length])
		//左乘右除，把二进制位向右移动位
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 255)

		conn, err = net.DialTimeout("ip:icmp", host, time.Duration(timeout*1000*1000))

		//todo test
		//ip := conn.RemoteAddr()
		fmt.Println("remote ip:",host)

		checkError3(err)

		starttime = time.Now()
		//conn.SetReadDeadline可以在未收到数据的指定时间内停止Read等待，并返回错误err，然后判定请求超时
		conn.SetDeadline(starttime.Add(time.Duration(timeout * 1000 * 1000)))
		//onn.Write方法执行之后也就发送了一条ICMP请求，同时进行计时和计次
		_, err = conn.Write(msg[0:length])

		//在使用Go语言的net.Dial函数时，发送echo request报文时，不用考虑i前20个字节的ip头；
		// 但是在接收到echo response消息时，前20字节是ip头。后面的内容才是icmp的内容，应该与echo request的内容一致
		const ECHO_REPLY_HEAD_LEN = 20

		var receive []byte = make([]byte, ECHO_REPLY_HEAD_LEN+length)
		n, err := conn.Read(receive)
		_ = n

		var endduration int = int(int64(time.Since(starttime)) / (1000 * 1000))

		sumT += endduration

		time.Sleep(1000 * 1000 * 1000)

		//除了判断err!=nil，还有判断请求和应答的ID标识符，sequence序列码是否一致，以及ICMP是否超时（receive[ECHO_REPLY_HEAD_LEN] == 11，即ICMP报头的类型为11时表示ICMP超时）
		if err != nil || receive[ECHO_REPLY_HEAD_LEN+4] != msg[4] || receive[ECHO_REPLY_HEAD_LEN+5] != msg[5] || receive[ECHO_REPLY_HEAD_LEN+6] != msg[6] || receive[ECHO_REPLY_HEAD_LEN+7] != msg[7] || endduration >= int(timeout) || receive[ECHO_REPLY_HEAD_LEN] == 11 {
			lostN++
			//todo
			//fmt.Println("对 " + cname + "[" + ip.String() + "]" + " 的请求超时。")
			fmt.Println("对 " + cname + "[" + host + "]" + " 的请求超时。")
			return errors.New("timeout")
		} else {
			if shortT == -1 {
				shortT = endduration
			} else if shortT > endduration {
				shortT = endduration
			}
			if longT == -1 {
				longT = endduration
			} else if longT < endduration {
				longT = endduration
			}
			recvN++
			ttl := int(receive[8])
			//          fmt.Println(ttl)
			//todo
			//fmt.Println("来自 " + cname + "[" + ip.String() + "]" + " 的回复: 字节=32 时间=" + strconv.Itoa(endduration) + "ms TTL=" + strconv.Itoa(ttl))
			fmt.Println("来自 " + cname + "[" + host + "]" + " 的回复: 字节=32 时间=" + strconv.Itoa(endduration) + "ms TTL=" + strconv.Itoa(ttl))
		}

		seq++
		count--
	}
	//todo 先注释，用下一行测试
	//stat3(host, sendN, lostN, recvN, shortT, longT, sumT)
	return nil
}

func checkSum3(msg []byte) uint16 {
	sum := 0

	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256 // notice here, why *256?
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func checkError3(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func gensequence3(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genidentifier3(host string) (byte, byte) {
	return host[0], host[1]
}

func stat3(ip string, sendN int, lostN int, recvN int, shortT int, longT int, sumT int) {
	fmt.Println()
	fmt.Println(ip, " 的 Ping 统计信息:")
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%d%% 丢失)，\n", sendN, recvN, lostN, int(lostN*100/sendN))
	fmt.Println("往返行程的估计时间(以毫秒为单位):")
	if recvN != 0 {
		fmt.Printf("    最短 = %dms，最长 = %dms，平均 = %dms\n", shortT, longT, sumT/sendN)
	}
}