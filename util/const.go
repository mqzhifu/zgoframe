package util

//声网 access token 自带
type Privileges uint16

const (
	//公共头信息的  key
	SERVICE_HEADER_KEY = "service_diy"

	VERSION_LENGTH = 3
	APP_ID_LENGTH  = 32

	KJoinChannel        = 1
	KPublishAudioStream = 2
	KPublishVideoStream = 3
	KPublishDataStream  = 4

	KLoginRtm = 1000

	READ_QUEUE_ORDER_FONT = "0"
	READ_QUEUE_ORDER_TAIL = "$"

	//error_msg
	CODE_NOT_EXIST = 5555
	ERR_separate   = "-_-"

	TCP_MSG_SEPARATOR = "-|"
	CTX_DONE_PRE      = "ctx.done() "

	DIR_SEPARATOR = "/"
	STR_SEPARATOR = "#"

	SUPER_VISOR_PROCESS_NAME_SERVICE_PREFIX = "service_" //进程启动时，启程名称的前缀，方便统一管理
	//log
	CONTENT_TYPE_STRING = 0
	//CONTENT_TYPE_JSON = 1

	// ping
	TimeSliceLength  = 8
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
	// ping3
	ECHO_REQUEST_HEAD_LEN = 8
	ECHO_REPLY_HEAD_LEN   = 20
)

//声网 access token 自研
// Role Type
type RTCRole uint16
type RTMRole uint16

// Role consts
const (
	RoleAttendee   = 0
	RolePublisher  = 1
	RoleSubscriber = 2
	RoleAdmin      = 101

	RoleRtmUser = 1
)

//@parse 环境变量-整形
const (
	ENV_LOCAL_INT  = 1 //开发环境
	ENV_DEV_INT    = 2 //开发环境
	ENV_TEST_INT   = 3 //测试环境
	ENV_PRE_INT    = 4 //预发布环境
	ENV_ONLINE_INT = 5 //线上环境
)

// 环境变量-字符串
const (
	ENV_LOCAL_STR  = "local"  //开发环境
	ENV_DEV_STR    = "dev"    //开发环境
	ENV_TEST_STR   = "test"   //测试环境
	ENV_PRE_STR    = "pre"    //预发布环境
	ENV_ONLINE_STR = "online" //线上环境

)

//@parse error
const (
	LOG_LEVEL_DEBUG = 1 //调试
	LOG_LEVEL_INFO  = 2 //信息
	LOG_LEVEL_OFF   = 4 //关闭
)

//@parse 日志等级
const (
	LEVEL_INFO      = 1 << iota
	LEVEL_DEBUG     = 2   //2
	LEVEL_ERROR     = 4   //4
	LEVEL_PANIC     = 8   //8
	LEVEL_EMERGENCY = 16  //16
	LEVEL_ALERT     = 32  //32
	LEVEL_CRITICAL  = 64  //64
	LEVEL_WARNING   = 128 //128
	LEVEL_NOTICE    = 256 //256
	LEVEL_TRACE     = 512 //512
	LEVEL_ALL       = LEVEL_INFO | LEVEL_DEBUG | LEVEL_ERROR | LEVEL_PANIC | LEVEL_EMERGENCY | LEVEL_ALERT | LEVEL_CRITICAL | LEVEL_WARNING | LEVEL_NOTICE | LEVEL_TRACE
	LEVEL_DEV       = LEVEL_INFO | LEVEL_DEBUG | LEVEL_ERROR | LEVEL_PANIC | LEVEL_TRACE
	LEVEL_ONLINE    = LEVEL_INFO | LEVEL_ERROR | LEVEL_PANIC
)

//(status 冲突 暂放弃,不能删 frame_sync game_match 还在使用，后期优化吧)
//@parse 玩家当着在线状态
const (
	PLAYER_STATUS_ONLINE  = 1 //在线
	PLAYER_STATUS_OFFLINE = 2 //离线
)

//@parse 一个副本的，一条消息的，同步状态
const (
	PLAYERS_ACK_STATUS_INIT = 1 //初始化
	PLAYERS_ACK_STATUS_WAIT = 2 //等待玩家确认
	PLAYERS_ACK_STATUS_OK   = 3 //所有玩家均已确认
)

//@parse 房间状态
const (
	ROOM_STATUS_INIT    = 1 //新房间，刚刚初始化，等待其它操作
	ROOM_STATUS_EXECING = 2 //已开始游戏
	ROOM_STATUS_END     = 3 //已结束
	ROOM_STATUS_READY   = 4 //准备中
	ROOM_STATUS_PAUSE   = 5 //有玩家掉线，暂停中
)

//@parse 协议类型
const (
	PROTOCOL_TCP       = 1 //传输协议 TCP
	PROTOCOL_UDP       = 3 //传输协议 UDP
	PROTOCOL_WEBSOCKET = 2 //传输协议 WEB-SOCKET
)

//@parse 长连接connFD的状态
const (
	CONN_STATUS_INIT      = 1 //初始化
	CONN_STATUS_EXECING   = 2 //运行中
	CONN_STATUS_CLOSE     = 3 //已关闭
	CONN_STATUS_CLOSE_ING = 4 //关闭中，防止重复关闭，不能用锁，因为：并发变串行后，还能重复关闭
)

//@parse 传输类型
const (
	CONTENT_TYPE_JSON     = 1 //内容类型 json
	CONTENT_TYPE_PROTOBUF = 2 //proto_buf
)

//@parse metricsc操作类型
const (
	METRICS_OPT_PLUS = 1 //1累加
	METRICS_OPT_INC  = 2 //2加加
	METRICS_OPT_LESS = 3 //3累减
	METRICS_OPT_DIM  = 4 //4减减

)

//@parse NETWAY类状态
const (
	NETWAY_STATUS_INIT  = 1 //网关状态 初始化中
	NETWAY_STATUS_START = 2 //网关状态 开始初始化
	NETWAY_STATUS_CLOSE = 3 //网关状态 已关闭
)

//@parse xxxx
const (
	TRAN_MESSAGE_TYPE_CHAR   = 1 //网络传输数据格式：字符流
	TRAN_MESSAGE_TYPE_BINARY = 2 //网络传输数据格式：二进制
)

//@parse 是否上传文件同时存储OSS
const (
	UPLOAD_STORE_OSS_OFF = 0 //关闭
	UPLOAD_STORE_OSS_ALI = 1 //阿里
)

//@parse 是否上传文件同时存储本地
const (
	UPLOAD_STORE_LOCAL_OFF  = 1 //关闭
	UPLOAD_STORE_LOCAL_OPEN = 2 //打开
)

//@parse 文件类型
const (
	FILE_TYPE_ALL   = 1 //全部
	FILE_TYPE_IMG   = 2 //图片
	FILE_TYPE_DOC   = 3 //文档
	FILE_TYPE_VIDEO = 4 //视频
	FILE_TYPE_AUDIO = 5 //音频
)

//@parse 长连接FD关闭类型
const (
	CLOSE_SOURCE_CREATE                = 3  //初始化 连接类失败，可能是连接数过大
	CLOSE_SOURCE_SERVER_HAS_CLOSE      = 11 //服务端状态已关闭
	CLOSE_SOURCE_CLIENT                = 1  //客户端-主动断开连接
	CLOSE_SOURCE_AUTH_FAILED           = 21 //客户端首次连接，登陆动作,服务端验证失败
	CLOSE_SOURCE_FD_READ_EMPTY         = 22 //客户端首次连接，登陆动作,服务端read信息为空
	CLOSE_SOURCE_FD_PARSE_CONTENT      = 23 //客户端首次连接，登陆动作,解析内容时出错
	CLOSE_SOURCE_FIRST_NO_LOGIN        = 24 //客户端首次连接，登陆动作,内容解出来了，但是action!=login
	CLOSE_SOURCE_FIRST_PARSER_LOGIN    = 25 //login  登陆不出结构体内容
	CLOSE_SOURCE_OPEN_PANIC            = 31 //初始化 新连接创建成功后，上层要再重新做一次连接，结果未知panic
	CLOSE_SOURCE_MAX_CLIENT            = 32 //当前连接数过大
	CLOSE_SOURCE_OVERRIDE              = 4  //创建新连接时，发现，该用户还有一个未关闭的连接,kickoff模式下，这条就没意义了
	CLOSE_SOURCE_TIMEOUT               = 5  //最后更新时间 ，超时.后台守护协程触发
	CLOSE_SOURCE_SIGNAL_QUIT           = 6  //接收到关闭信号，netWay.Quit触发
	CLOSE_SOURCE_CLIENT_WS_FD_GONE     = 7  //S端读取连接消息时，异常了~可能是：客户端关闭了连接
	CLOSE_SOURCE_SEND_MESSAGE          = 8  //S端给某个连接发消息，结果失败了，这里概率是连接已经断了
	CLOSE_SOURCE_CONN_RESET_BY_PEER    = 81 //对端，如果直接关闭网络，或者崩溃之类的，类库捕捉不到这个事件
	CLOSE_SOURCE_CONN_SHUTDOWN         = 12 //conn 已关闭
	CLOSE_SOURCE_CONN_LOGIN_ROUTER_ERR = 13 //登陆，路由一个方法时，未找到该方法
	CLOSE_SOURCE_RTT_TIMEOUT           = 91 //S端已收到了RTT的响应，但已超时
	CLOSE_SOURCE_RTT_TIMER_OUT         = 92 //RTT超时，定时器触发

)

//@parse http-curl类，数据传输类型
const (
	HTTP_DATA_CONTENT_TYPE_JSON   = 1 //JSON
	HTTP_DATA_CONTENT_TYPE_Nornal = 2 //普通
)

//@parse 文件存储hash类型
const (
	FILE_HASH_NONE  = 0 // 没有
	FILE_HASH_MONTH = 1 // 月
	FILE_HASH_DAY   = 2 //天
	FILE_HASH_HOUR  = 3 //小时
)

//@parse
const (
	OUT_TARGET_SC   = 1 << iota
	OUT_TARGET_FILE //文件
	OUT_TARGET_NET  //网络传输

	OUT_TARGET_ALL = OUT_TARGET_SC | OUT_TARGET_FILE | OUT_TARGET_NET

	OUT_TARGET_NET_TCP = 1 //网络协议为TCP
	OUT_TARGET_NET_UDP = 2 //网络协议为UDP
)

//@parse SERVER_STATUS
const (
	SERVER_STATUS_NORMAL = 1 //正常
	SERVER_STATUS_CLOSE  = 2 //已关闭
)

//@parse SERVER_PING服务器的状态
const (
	SERVER_PING_OK   = 1 //正常：PING 成功
	SERVER_PING_FAIL = 2 //异常：PING 失败了
)

//@parse service协议类型
const (
	SERVICE_PROTOCOL_HTTP      = 1 //HTTP
	SERVICE_PROTOCOL_GRPC      = 2 //GRPC
	SERVICE_PROTOCOL_WEBSOCKET = 3 //WS
	SERVICE_PROTOCOL_TCP       = 4 //TCP
)

//@parse 服务发现的类型，分布式DB
const (
	SERVICE_DISCOVERY_ETCD   = 1 //ETCD
	SERVICE_DISCOVERY_CONSUL = 2 //CONSULE
)

//@parse 服务发现负载类型
const (
	LOAD_BALANCE_ROBIN = 1 //轮询
	LOAD_BALANCE_HASH  = 2 //固定分子hash
)

//(主要是给CICD部署时使用，最终给前端使用)
//@parse super_visor错误类型
const (
	SV_ERROR_NONE      = 0 //无
	SV_ERROR_INIT      = 1 //初始化
	SV_ERROR_CONN      = 2 //连接中
	SV_ERROR_NOT_FOUND = 3 //未找到
)
