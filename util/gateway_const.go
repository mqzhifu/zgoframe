package util

type CmdArgs struct {
	Env         string `seq:"1" err:"env=local" desc:"环境变量： local test dev pre online"`
	Ip          string `seq:"2" err:"ip=127.0.0.1" desc:"监听的IP地址"`
	HttpPort    string `seq:"3" err:"HttpPort=2222" desc:"短连接监听端口号"`
	WsPort      string `seq:"4" err:"WsPort=2223" desc:"websocket监听端口号"`
	TcpPort     string `seq:"5" err:"TcpPort=2224" desc:"tcp协议监听端口号"`
	LogBasePath string `seq:"6" err:"log_base_path=/golang/logs" desc:"日志文件保存位置"`
	//ClientServer 	string 	`seq:"7" err:"cs=serve"`
}
