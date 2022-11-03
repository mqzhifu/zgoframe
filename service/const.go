package service

const (
	CTX_DONE_PRE           = "ctx.done() " //字符串标识，用于打印输出信息时的前缀
	MAIL_ADMIN_USER_UID    = 9999          //管理员默认UID，主要用于：确定发送者的UID
	PlayerMatchingMaxTimes = 3             //一个玩家，参与匹配机制的最大次数，超过这个次数，证明不用再匹配了，目前没用上，目前使用的还是绝对的超时时间为准

	//rule规格配置表
	RULE_TYPE_TEAM_VS         = 1 //moba 类 5V5 对战类型
	RULE_TYPE_TEAM_EACH_OTHER = 2 //类似吃鸡 多个队伍互相撕杀

	RuleEtcdConfigPrefix = "/v1/conf/matches/" //etcd中  ， 存放 rule  集合的前缀

	//Separation        = "#" //redis 内容-字符串分隔符
	//PayloadSeparation = "%" //push时的内容，缓存进redis时
	////RedisSeparation     = "_"           //redis key 分隔符
	//IdsSeparation = "," //多个ID 分隔符

	//微服务
	//SERVICE_MSG_SERVER = "msgServer"
	//SERVICE_MATCH_NAME		="gamematch"
	//SERVICE_PREFIX = "/v1/service"		//微服务前缀
	//SIGNAL_GOROUTINE_EXEC_EXIT   = 1 //通知协程，执行结束操作
	//SIGNAL_GOROUTINE_EXIT_FINISH = 2 //协程，通知父协程，已结束
	//SIGNAL_GOROUTINE_EXEC_ING    = 6 //协程，通知父协程，已执行
	//
	//SIGNAL_EXIT = 3 //结束所有后台守护协程，退出程序
	//
	//SIGNAL_QUIT_SOURCE            = 4
	//SIGNAL_QUIT_SOURCE_RULE_WATCH = 5
)

//@parse 配置中心-数据持久化类型
const (
	PERSISTENCE_TYPE_OFF     = 0 //关闭
	PERSISTENCE_TYPE_MYSQL   = 1 //mysql数据库
	PERSISTENCE_TYPE_REDIS   = 2 //redis缓存
	PERSISTENCE_TYPE_FILE    = 3 //文件
	PERSISTENCE_TYPE_ETCD    = 4 //etcd
	PERSISTENCE_TYPE_CONSULE = 5 //consul
)

//@parse 帧同步-锁模式
const (
	LOCK_MODE_PESSIMISTIC = 1 //囚徒
	LOCK_MODE_OPTIMISTIC  = 2 //乐观
)

//@parse 帧同步，一个副本的，一条消息的，同步状态
const (
	PLAYERS_ACK_STATUS_INIT = 1 //初始化
	PLAYERS_ACK_STATUS_WAIT = 2 //等待玩家确认
	PLAYERS_ACK_STATUS_OK   = 3 //所有玩家均已确认
)

//(帧同步 游戏匹配 好像都在用)
//@parse 玩家状态
const (
	PLAYER_STATUS_ONLINE  = 1 //在线
	PLAYER_STATUS_OFFLINE = 2 //离线
)

//@parse 帧同步-房间内，玩家准备状态
const (
	PLAYER_NO_READY  = 1 //玩家未准备
	PLAYER_HAS_READY = 2 //玩家已准备
)

//@parse 站内邮件-是否失效
const (
	MAIL_EXPIRE_TRUE  = 1 //是
	MAIL_EXPIRE_FALSE = 1 //否
)

//INIT -> READY -> EXECING -> PAUSE -> END
//@parse 帧同步-房间状态
const (
	ROOM_STATUS_INIT    = 1 //新房间，刚刚初始化，等待其它操作
	ROOM_STATUS_EXECING = 2 //已开始游戏
	ROOM_STATUS_END     = 3 //已结束
	ROOM_STATUS_READY   = 4 //准备中
	ROOM_STATUS_PAUSE   = 5 //有玩家掉线，暂停中
)

//@parse 站内邮件-发送类型
const (
	MAIL_PEOPLE_PERSON = 1 //单发
	MAIL_PEOPLE_ALL    = 2 //群发
	MAIL_PEOPLE_GROUP  = 3 //指定group
	MAIL_PEOPLE_TAG    = 4 //指定tag
	MAIL_PEOPLE_UIDS   = 5 //指定UIDS
)

//@parse 站内邮件-消息收件箱
const (
	MAIL_IN_BOX     = 1 //收件箱
	MAIL_OUT_BOX    = 2 //发件箱
	MAIL_IN_DEL_BOX = 3 //已删除箱子
	MAIL_ALL_BOX    = 4 //全部
)

//@parse 站内邮件-消息接收是否已读
const (
	RECEIVER_READ_TRUE  = 1 //是
	RECEIVER_READ_FALSE = 2 //否
)

//@parse 站内邮件-消息是否删除
const (
	RECEIVER_DEL_TRUE  = 1 //是
	RECEIVER_DEL_FALSE = 2 //否
)

//@parse http状态
const (
	HTTPD_RULE_STATE_INIT  = 1 //初始化中
	HTTPD_RULE_STATE_OK    = 2 //正常运行中
	HTTPD_RULE_STATE_CLOSE = 3 //已关闭
	HTTPD_RULE_STATE_UKNOW = 4 //未知
)

//@parse 游戏匹配-一条rule的状态
const (
	RuleStatusOnline  = 1
	RuleStatusOffline = 2
	RuleStatusDelete  = 3
)

//@parse 游戏匹配-筛选策略(如何从匹配池里拿用户)
const (
	FilterFlagAll      = 1 //全匹配，无差别匹配，rule 没有配置权重公式时，使用
	FilterFlagBlock    = 2 //权重公式，块-匹配
	FilterFlagBlockInc = 3 //权重公式，递增块匹配
	FilterFlagDIY      = 4 //权重公式，自定义块匹配
)

//@parse 游戏匹配-rule数据来源
const (
	GAME_MATCH_DATA_SOURCE_TYPE_ETCD    = 1
	GAME_MATCH_DATA_SOURCE_TYPE_DB      = 2
	GAME_MATCH_DATA_SOURCE_TYPE_SERVICE = 3
)

//@parse 游戏匹配-玩家状态
const (
	GAME_MATCH_PLAYER_STATUS_NOT_EXIST = 1 //redis中还没有该玩家信息
	GAME_MATCH_PLAYER_STATUS_SIGN      = 2 //已报名，等待匹配
	GAME_MATCH_PLAYER_STATUS_SUCCESS   = 3 //匹配成功，等待拿走
	GAME_MATCH_PLAYER_STATUS_INIT      = 4 //初始化阶段

)

//@parse 一条RULE的状态
const (
	GAME_MATCH_RULE_STATUS_INIT  = 1 //初始化
	GAME_MATCH_RULE_STATUS_EXEC  = 2 //运行中
	GAME_MATCH_RULE_STATUS_CLOSE = 3 //关闭
)

//@parse 游戏匹配-小组类型
const (
	GAME_MATCH_GROUP_TYPE_SIGN    = 1 //报名
	GAME_MATCH_GROUP_TYPE_SUCCESS = 2 //报名成功

)

//@parse 游戏匹配-HTTP推送状态
const (
	PushCategorySignTimeout    = 1 //报名超时
	PushCategorySuccess        = 2 //匹配成功
	PushCategorySuccessTimeout = 3 //匹配成功超时
)

//@parse 游戏匹配-推送状态
const (
	PUSH_STATUS_WAIT   = 1 //等待推送
	PUSH_STATUS_RETRY  = 2 //已推送过，但失败了，等待重试
	PUSH_STATUS_OK     = 3 //推送成功
	PUSH_STATUS_FAILED = 4 //推送失败
)

//@parse 呼叫类型
const (
	CALL_USER_PEOPLE_ALL     = 1 //呼叫所有人
	CALL_USER_PEOPLE_GROUP   = 2 //按照<小组>呼叫
	CALL_USER_PEOPLE_PROVIDE = 3 //用户自己指定呼叫的人
)

//@parse 房间状态2
const (
	RTC_ROOM_STATUS_CALLING = 1 //房间状态：呼叫中
	RTC_ROOM_STATUS_EXECING = 2 //房间状态：运行中
	RTC_ROOM_STATUS_END     = 3 //房间状态：已结束
)

//@parse 房间结束类型
const (
	RTC_ROOM_END_STATUS_TIMEOUT_CALLING = 10 //房间结束状态标识：呼叫超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_TIMEOUT_EXEC    = 11 //房间结束状态标识：运行超时(也可能是连接断了)
	RTC_ROOM_END_STATUS_CONN_CLOSE      = 2  //房间结束状态标识：用户退出
	RTC_ROOM_END_STATUS_DENY            = 3  //房间结束状态标识：被呼叫人拒绝
	RTC_ROOM_END_STATUS_CANCEL          = 4  //房间结束状态标识：呼叫者取消
	RTC_ROOM_END_STATUS_USER_LEAVE      = 4  //房间结束状态标识：用户离开
)

//@parse RTC推送消息类型
const (
	RTC_PUSH_MSG_EVENT_FD_CREATE_REPEAT = 400
	RTC_PUSH_MSG_EVENT_UID_NOT_IN_MAP   = 401
)
