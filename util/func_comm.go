package util

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

// 这个函数只是懒......
func MyPrint(a ...interface{}) (n int, err error) {
	return fmt.Println(a)
}

// 这个函数只是懒......debug 调试使用
func ExitPrint(a ...interface{}) {
	fmt.Println(a)
	os.Exit(999)
}

// 输出复杂类型的数据，如：结构体
func MyComplexPrint(a ...interface{}) (n int, err error) {
	return fmt.Printf("%+v", a)
}

// 四舍五入
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// 一次获取N个空格，用于测试时 输出时 加些空格格式化内容
func GetSpaceStr(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += " "
	}
	return str
}

// 获取一个随机数：int
func GetRandIntNum(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

// 获取一个随机数：int32
func GetRandInt32Num(max int32) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(int(max))
}

// 获取一个随机整数:可设置范围
func GetRandIntNumRange(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

// 自带的<strconv.Atoi>函数返回的是两个参数，很麻烦，这里简化，只返回一个参数
func Atoi(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

// 类型转换：浮点转字符串
func FloatToString(number float32, little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(number), 'f', little, 64)
}

// 类型转换：float64 -> 字符串
func Float64ToString(number float64, little int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(number, 'f', little, 64)
}

// 类型转换：string -> float32
func StringToFloat(str string) float32 {
	v1, _ := strconv.ParseFloat(str, 32)
	number := float32(v1)
	return number
}

// 指令行映射，根据使用者提供的map从指令行读取取，有查错，并映射进去
// 练习了两个知识点：1给定一个struct，和一堆string，反射struct成员值，把string映射进map里 2从一个struct的tag中解析数据
func CmsArgs(data interface{}) (argMap map[string]string, err error) {
	//读取 data 类型 反射
	typeOfCmsArgs := reflect.TypeOf(data)
	if len(os.Args) < typeOfCmsArgs.NumField()+1 {
		errInfo := "os.Args len < " + strconv.Itoa(typeOfCmsArgs.NumField()) + " , eg:"
		for i := 0; i < typeOfCmsArgs.NumField(); i++ {
			memVar := typeOfCmsArgs.Field(i)
			errInfo += memVar.Tag.Get("err") + " ,"
		}
		return argMap, errors.New(errInfo)
	}
	cmsArg := make(map[string]string)
	for i := 0; i < typeOfCmsArgs.NumField(); i++ {
		memVar := typeOfCmsArgs.Field(i)   //获取结构体中的一个成员对象
		sqeNum := memVar.Tag.Get("seq")    //读取出该成员的tag
		num, _ := strconv.Atoi(sqeNum)     //转换成字符串
		cmsArg[memVar.Name] = os.Args[num] //根据成员名，写入map中
	}
	return cmsArg, nil
}

// 获取本机的Ip地址
func GetLocalIp() (ip string, err error) {
	netInterfaces, err := net.Interfaces()
	//MyPrint(netInterfaces, err)
	if err != nil {
		return ip, errors.New("net.Interfaces failed, err:" + err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}

	return ip, nil
}

// GetRemoteClientIp 获取远程客户端IP
func GetRemoteClientIp(r *http.Request) string {
	remoteIp := r.RemoteAddr

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteIp = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}

	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp

}
func PingByShell(host string, timeOut string) bool {
	os := runtime.GOOS
	//sendPackageNum := 4
	var cmd *exec.Cmd
	str := "PingByShell host:" + host + " timeOut:" + timeOut + " os:" + os
	cmdStr := ""
	if os == "darwin" {
		cmdStr = host + " -c 4 -t " + timeOut
	} else if os == "windows" {
		cmdStr = host + " -n 4 -w " + timeOut
		//cmd = exec.Command("ping", host +  "-n 4 -w" +  timeOut)
	} else if os == "linux" {
		cmdStr = host + " -c 4 -w " + timeOut
		//cmd = exec.Command("ping", host +  "-c 4 -w" +  timeOut)
	} else {
		MyPrint("get os err:", os)
		return false
	}
	MyPrint(str + " ping " + cmdStr)
	cmd = exec.Command("ping", cmdStr)
	//MyPrint("ping " +  host + "-c " +timeOut)
	//fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	err := cmd.Run()
	fmt.Println("PingByShell finish ,  :", time.Now().Unix())
	if err != nil {
		fmt.Println("PingByShell err:", err.Error())
		return false
	} else {
		fmt.Println("PingByShell ok~")
	}
	return true

}

func CheckIpPort(ip string, port string, timeout int64) bool {
	timeoutDuration := time.Duration(timeout) * time.Second
	_, err := net.DialTimeout("tcp", ip+":"+port, timeoutDuration)
	if err != nil {
		//fmt.Println("Site unreachable, error: ", err)
		return false
	}
	return true
}
