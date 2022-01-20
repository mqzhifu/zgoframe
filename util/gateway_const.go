package util

const (
	CLOSE_SOURCE_SERVER_HAS_CLOSE   = 11 //服务端状态已关闭
	CLOSE_SOURCE_CLIENT 			= 1	//客户端-主动断开连接
	CLOSE_SOURCE_AUTH_FAILED 		= 21//客户端首次连接，登陆动作,服务端验证失败
	CLOSE_SOURCE_FD_READ_EMPTY 		= 22//客户端首次连接，登陆动作,服务端read信息为空
	CLOSE_SOURCE_FD_PARSE_CONTENT 	= 23//客户端首次连接，登陆动作,解析内容时出错
	CLOSE_SOURCE_FIRST_NO_LOGIN 	= 24//客户端首次连接，登陆动作,内容解出来了，但是action!=login
	//CLOSE_SOURCE_CREATE 			= 3	//初始化 连接类失败，可能是连接数过大
	CLOSE_SOURCE_OPEN_PANIC			= 31//初始化 新连接创建成功后，上层要再重新做一次连接，结果未知panic
	CLOSE_SOURCE_MAX_CLIENT			= 32//当前连接数过大
	CLOSE_SOURCE_OVERRIDE 			= 4	//创建新连接时，发现，该用户还有一个未关闭的连接,kickoff模式下，这条就没意义了
	CLOSE_SOURCE_TIMEOUT 			= 5	//最后更新时间 ，超时.后台守护协程触发
	CLOSE_SOURCE_SIGNAL_QUIT 		= 6 //接收到关闭信号，netWay.Quit触发
	CLOSE_SOURCE_CLIENT_WS_FD_GONE 	= 7	//S端读取连接消息时，异常了~可能是：客户端关闭了连接
	CLOSE_SOURCE_SEND_MESSAGE 		= 8 //S端给某个连接发消息，结果失败了，这里概率是连接已经断了
	CLOSE_SOURCE_CONN_SHUTDOWN 		= 11

	CLOSE_SOURCE_RTT_TIMEOUT 		= 91//S端已收到了RTT的响应，但已超时
	CLOSE_SOURCE_RTT_TIMER_OUT 		= 92//RTT超时，定时器触发

	TCP_MSG_SEPARATOR   = "-|"

	CTX_DONE_PRE = "ctx.done() "

	CONTENT_TYPE_JSON 		= 1		//内容类型 json
	CONTENT_TYPE_PROTOBUF 	= 2		//proto_buf

	CONN_STATUS_INIT 	= 1	//初始化
	CONN_STATUS_EXECING = 2	//运行中
	CONN_STATUS_CLOSE 	= 3	//已关闭

	PROTOCOL_TCP 		= 1
	PROTOCOL_UDP 		= 3
	PROTOCOL_WEBSOCKET 	= 2

	ROOM_STATUS_INIT 	= 1		//新房间，刚刚初始化，等待其它操作
	ROOM_STATUS_EXECING = 2		//已开始游戏
	ROOM_STATUS_END 	= 3		//已结束
	ROOM_STATUS_READY 	= 4		//准备中
	ROOM_STATUS_PAUSE 	= 5		//有玩家掉线，暂停中

	//一个副本的，一条消息的，同步状态
	PLAYERS_ACK_STATUS_INIT = 1	//初始化
	PLAYERS_ACK_STATUS_WAIT = 2	//等待玩家确认
	PLAYERS_ACK_STATUS_OK 	= 3	//所有玩家均已确认

	PLAYER_STATUS_ONLINE 	= 1	//在线
	PLAYER_STATUS_OFFLINE 	= 2	//离线

	LOCK_MODE_PESSIMISTIC 	= 1	//囚徒
	LOCK_MODE_OPTIMISTIC 	= 2	//乐观

	METRICS_OPT_PLUS 	= 1  	//1累加
	METRICS_OPT_INC 	= 2		//2加加
	METRICS_OPT_LESS 	= 3		//3累减
	METRICS_OPT_DIM 	= 4		//4减减

	NETWAY_STATUS_INIT = 1
	NETWAY_STATUS_START = 2
	NETWAY_STATUS_CLOSE = 3

	TRAN_MESSAGE_TYPE_CHAR 	 = 1 //网络传输数据格式：字符流
	TRAN_MESSAGE_TYPE_BINARY = 2 //网络传输数据格式：二进制
)


type CmdArgs struct {
	Env 			string	`seq:"1" err:"env=local" desc:"环境变量： local test dev pre online"`
	Ip 				string	`seq:"2" err:"ip=127.0.0.1" desc:"监听的IP地址"`
	HttpPort 		string	`seq:"3" err:"HttpPort=2222" desc:"短连接监听端口号"`
	WsPort 			string	`seq:"4" err:"WsPort=2223" desc:"websocket监听端口号"`
	TcpPort 		string	`seq:"5" err:"TcpPort=2224" desc:"tcp协议监听端口号"`
	LogBasePath 	string	`seq:"6" err:"log_base_path=/golang/logs" desc:"日志文件保存位置"`
	//ClientServer 	string 	`seq:"7" err:"cs=serve"`
}