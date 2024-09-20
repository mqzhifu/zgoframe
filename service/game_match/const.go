package gamematch

// @parse 配置中心-数据持久化类型
const (
	PERSISTENCE_TYPE_OFF     = 0 //关闭
	PERSISTENCE_TYPE_MYSQL   = 1 //mysql数据库
	PERSISTENCE_TYPE_REDIS   = 2 //redis缓存
	PERSISTENCE_TYPE_FILE    = 3 //文件
	PERSISTENCE_TYPE_ETCD    = 4 //etcd
	PERSISTENCE_TYPE_CONSULE = 5 //consul
)

// @parse 游戏匹配-筛选策略(如何从匹配池里拿用户)
const (
	FilterFlagAll      = 1 //全匹配，无差别匹配，rule 没有配置权重公式时，使用
	FilterFlagBlock    = 2 //权重公式，块-匹配
	FilterFlagBlockInc = 3 //权重公式，递增块匹配
	FilterFlagDIY      = 4 //权重公式，自定义块匹配
)

const (
	CTX_DONE_PRE           = "ctx.done() " //字符串标识，用于打印输出信息时的前缀
	PlayerMatchingMaxTimes = 3             //一个玩家，参与匹配机制的最大次数，超过这个次数，证明不用再匹配了，目前没用上，目前使用的还是绝对的超时时间为准

	//rule规格配置表
	RULE_TYPE_TEAM_VS         = 1 //moba 类 5V5 对战类型
	RULE_TYPE_TEAM_EACH_OTHER = 2 //类似吃鸡 多个队伍互相撕杀

	RuleEtcdConfigPrefix = "/v1/conf/matches/" //etcd中  ， 存放 rule  集合的前缀

	//Separation        = "#" //redis 内容-字符串分隔符
	//PayloadSeparation = "%" //push时的内容，缓存进redis时
	//RedisSeparation     = "_"           //redis key 分隔符
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

// @parse 一条RULE的状态
const (
	GAME_MATCH_RULE_STATUS_INIT  = 1 //初始化
	GAME_MATCH_RULE_STATUS_EXEC  = 2 //运行中
	GAME_MATCH_RULE_STATUS_CLOSE = 3 //关闭
)

// @parse 游戏匹配-小组类型
const (
	GAME_MATCH_GROUP_TYPE_SIGN    = 1 //报名
	GAME_MATCH_GROUP_TYPE_SUCCESS = 2 //报名成功

)

// @parse 游戏匹配-HTTP推送状态
const (
	PushCategorySignTimeout    = 1 //报名超时
	PushCategorySuccess        = 2 //匹配成功
	PushCategorySuccessTimeout = 3 //匹配成功超时
)

// @parse 游戏匹配-推送状态
const (
	PUSH_STATUS_WAIT   = 1 //等待推送
	PUSH_STATUS_RETRY  = 2 //已推送过，但失败了，等待重试
	PUSH_STATUS_OK     = 3 //推送成功
	PUSH_STATUS_FAILED = 4 //推送失败
)

// @parse 游戏匹配-一条rule的状态
const (
	RuleStatusOnline  = 1 //游戏匹配规则字段，状态，在线
	RuleStatusOffline = 2 //游戏匹配规则字段，状态，下线
	RuleStatusDelete  = 3 //游戏匹配规则字段，状态，已删除
)

// @parse 游戏匹配-rule 数据来源
const (
	GAME_MATCH_DATA_SOURCE_TYPE_ETCD    = 1 //游戏匹配 rule 数据来源：ETCD
	GAME_MATCH_DATA_SOURCE_TYPE_DB      = 2 //游戏匹配 rule 数据来源：DB
	GAME_MATCH_DATA_SOURCE_TYPE_SERVICE = 3 //游戏匹配 rule 数据来源：服务
)

// @parse 游戏匹配-玩家状态
const (
	GAME_MATCH_PLAYER_STATUS_NOT_EXIST = 1 //redis中还没有该玩家信息
	GAME_MATCH_PLAYER_STATUS_SIGN      = 2 //已报名，等待匹配
	GAME_MATCH_PLAYER_STATUS_SUCCESS   = 3 //匹配成功，等待拿走
	GAME_MATCH_PLAYER_STATUS_INIT      = 4 //初始化阶段

)
