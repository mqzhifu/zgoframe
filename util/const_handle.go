package util

//注：此文件是由PHP动态生成，不要做任何修改，均会被覆盖

type ConstItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Desc  string      `json:"desc"`
}

type EnumConst struct {
	CommonPrefix string      `json:"common_prefix"`
	Desc         string      `json:"desc"`
	ConstList    []ConstItem `json:"const_list"`
	Type         string      `json:"type"`
}

type ConstHandle struct {
	EnumConstPool map[string]EnumConst
}

func NewConstHandle() *ConstHandle {
	constHandle := new(ConstHandle)
	constHandle.EnumConstPool = make(map[string]EnumConst)
	constHandle.Init()
	return constHandle
}

func (constHandle *ConstHandle) Init() {
	var constItemList []ConstItem
	var constItem ConstItem
	var enumConst EnumConst
        	constItem = ConstItem{
		Key:   "HTTP_RES_COMM_ERROR",
		Value: 4,
		Desc:  "失败",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "HTTP_RES_COMM_SUCCESS",
		Value: 200,
		Desc:  "成功",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "HTTP_RES_COMM_",
		Desc:         "HTTP 公共响应：自定义 状态码",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "ALERT_LEVEL_ALL",
		Value: -1,
		Desc:  "全部",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ALERT_LEVEL_SMS",
		Value: 1,
		Desc:  "短信",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ALERT_LEVEL_EMAIL",
		Value: 2,
		Desc:  "邮件",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ALERT_LEVEL_FEISHU",
		Value: 4,
		Desc:  "飞书",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ALERT_LEVEL_WEIXIN",
		Value: 8,
		Desc:  "微信",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ALERT_LEVEL_DINGDING",
		Value: 16,
		Desc:  "钉钉",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "ALERT_LEVEL_",
		Desc:         "报警发送渠道类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "ALERT_SEND_SYNC",
		Value: 1,
		Desc:  "同步",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ALERT_SEND_ASYNC",
		Value: 2,
		Desc:  "异步",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "ALERT_SEND_",
		Desc:         "报警发送方式类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PERSISTENCE_TYPE_OFF",
		Value: 0,
		Desc:  "关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PERSISTENCE_TYPE_MYSQL",
		Value: 1,
		Desc:  "mysql数据库",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PERSISTENCE_TYPE_REDIS",
		Value: 2,
		Desc:  "redis缓存",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PERSISTENCE_TYPE_FILE",
		Value: 3,
		Desc:  "文件",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PERSISTENCE_TYPE_ETCD",
		Value: 4,
		Desc:  "etcd",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PERSISTENCE_TYPE_CONSULE",
		Value: 5,
		Desc:  "consul",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PERSISTENCE_TYPE_",
		Desc:         "配置中心-数据持久化类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "LOCK_MODE_PESSIMISTIC",
		Value: 1,
		Desc:  "囚徒",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LOCK_MODE_OPTIMISTIC",
		Value: 2,
		Desc:  "乐观",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "LOCK_MODE_",
		Desc:         "帧同步-锁模式",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PLAYERS_ACK_STATUS_INIT",
		Value: 1,
		Desc:  "初始化",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYERS_ACK_STATUS_WAIT",
		Value: 2,
		Desc:  "等待玩家确认",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYERS_ACK_STATUS_OK",
		Value: 3,
		Desc:  "所有玩家均已确认",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PLAYERS_ACK_STATUS_",
		Desc:         "帧同步，一个副本的，一条消息的，同步状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PLAYER_STATUS_ONLINE",
		Value: 1,
		Desc:  "在线",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYER_STATUS_OFFLINE",
		Value: 2,
		Desc:  "离线",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PLAYER_STATUS_O",
		Desc:         "玩家状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PLAYER_NO_READY",
		Value: 1,
		Desc:  "玩家未准备",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYER_HAS_READY",
		Value: 2,
		Desc:  "玩家已准备",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PLAYER_",
		Desc:         "帧同步-房间内，玩家准备状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "MAIL_EXPIRE_TRUE",
		Value: 1,
		Desc:  "是",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_EXPIRE_FALSE",
		Value: 1,
		Desc:  "否",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "MAIL_EXPIRE_",
		Desc:         "站内邮件-是否失效",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "ROOM_STATUS_INIT",
		Value: 1,
		Desc:  "新房间，刚刚初始化，等待其它操作",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_EXECING",
		Value: 2,
		Desc:  "已开始游戏",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_END",
		Value: 3,
		Desc:  "已结束",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_READY",
		Value: 4,
		Desc:  "准备中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_PAUSE",
		Value: 5,
		Desc:  "有玩家掉线，暂停中",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "ROOM_STATUS_",
		Desc:         "帧同步-房间状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "MAIL_PEOPLE_PERSON",
		Value: 1,
		Desc:  "单发",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_PEOPLE_ALL",
		Value: 2,
		Desc:  "群发",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_PEOPLE_GROUP",
		Value: 3,
		Desc:  "指定group",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_PEOPLE_TAG",
		Value: 4,
		Desc:  "指定tag",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_PEOPLE_UIDS",
		Value: 5,
		Desc:  "指定UIDS",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "MAIL_PEOPLE_",
		Desc:         "站内邮件-发送类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "MAIL_IN_BOX",
		Value: 1,
		Desc:  "收件箱",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_OUT_BOX",
		Value: 2,
		Desc:  "发件箱",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_IN_DEL_BOX",
		Value: 3,
		Desc:  "已删除箱子",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "MAIL_ALL_BOX",
		Value: 4,
		Desc:  "全部",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "MAIL_",
		Desc:         "站内邮件-消息收件箱",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "RECEIVER_READ_TRUE",
		Value: 1,
		Desc:  "是",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RECEIVER_READ_FALSE",
		Value: 2,
		Desc:  "否",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "RECEIVER_READ_",
		Desc:         "站内邮件-消息接收是否已读",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "RECEIVER_DEL_TRUE",
		Value: 1,
		Desc:  "是",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RECEIVER_DEL_FALSE",
		Value: 2,
		Desc:  "否",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "RECEIVER_DEL_",
		Desc:         "站内邮件-消息是否删除",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "HTTPD_RULE_STATE_INIT",
		Value: 1,
		Desc:  "初始化中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "HTTPD_RULE_STATE_OK",
		Value: 2,
		Desc:  "正常运行中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "HTTPD_RULE_STATE_CLOSE",
		Value: 3,
		Desc:  "已关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "HTTPD_RULE_STATE_UKNOW",
		Value: 4,
		Desc:  "未知",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "HTTPD_RULE_STATE_",
		Desc:         "http状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "RuleStatusOnline",
		Value: 1,
		Desc:  "游戏匹配规则字段，状态，在线",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RuleStatusOffline",
		Value: 2,
		Desc:  "游戏匹配规则字段，状态，下线",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RuleStatusDelete",
		Value: 3,
		Desc:  "游戏匹配规则字段，状态，已删除",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "RuleStatus",
		Desc:         "游戏匹配-一条rule的状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "FilterFlagAll",
		Value: 1,
		Desc:  "全匹配，无差别匹配，rule 没有配置权重公式时，使用",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FilterFlagBlock",
		Value: 2,
		Desc:  "权重公式，块-匹配",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FilterFlagBlockInc",
		Value: 3,
		Desc:  "权重公式，递增块匹配",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FilterFlagDIY",
		Value: 4,
		Desc:  "权重公式，自定义块匹配",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "FilterFlag",
		Desc:         "游戏匹配-筛选策略(如何从匹配池里拿用户)",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "GAME_MATCH_DATA_SOURCE_TYPE_ETCD",
		Value: 1,
		Desc:  "游戏匹配 rule 数据来源：ETCD",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_DATA_SOURCE_TYPE_DB",
		Value: 2,
		Desc:  "游戏匹配 rule 数据来源：DB",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_DATA_SOURCE_TYPE_SERVICE",
		Value: 3,
		Desc:  "游戏匹配 rule 数据来源：服务",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "GAME_MATCH_DATA_SOURCE_TYPE_",
		Desc:         "游戏匹配-rule 数据来源",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "GAME_MATCH_PLAYER_STATUS_NOT_EXIST",
		Value: 1,
		Desc:  "redis中还没有该玩家信息",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_PLAYER_STATUS_SIGN",
		Value: 2,
		Desc:  "已报名，等待匹配",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_PLAYER_STATUS_SUCCESS",
		Value: 3,
		Desc:  "匹配成功，等待拿走",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_PLAYER_STATUS_INIT",
		Value: 4,
		Desc:  "初始化阶段",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "GAME_MATCH_PLAYER_STATUS_",
		Desc:         "游戏匹配-玩家状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "GAME_MATCH_RULE_STATUS_INIT",
		Value: 1,
		Desc:  "初始化",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_RULE_STATUS_EXEC",
		Value: 2,
		Desc:  "运行中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_RULE_STATUS_CLOSE",
		Value: 3,
		Desc:  "关闭",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "GAME_MATCH_RULE_STATUS_",
		Desc:         "一条RULE的状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "GAME_MATCH_GROUP_TYPE_SIGN",
		Value: 1,
		Desc:  "报名",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GAME_MATCH_GROUP_TYPE_SUCCESS",
		Value: 2,
		Desc:  "报名成功",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "GAME_MATCH_GROUP_TYPE_S",
		Desc:         "游戏匹配-小组类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PushCategorySignTimeout",
		Value: 1,
		Desc:  "报名超时",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PushCategorySuccess",
		Value: 2,
		Desc:  "匹配成功",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PushCategorySuccessTimeout",
		Value: 3,
		Desc:  "匹配成功超时",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PushCategoryS",
		Desc:         "游戏匹配-HTTP推送状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PUSH_STATUS_WAIT",
		Value: 1,
		Desc:  "等待推送",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PUSH_STATUS_RETRY",
		Value: 2,
		Desc:  "已推送过，但失败了，等待重试",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PUSH_STATUS_OK",
		Value: 3,
		Desc:  "推送成功",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PUSH_STATUS_FAILED",
		Value: 4,
		Desc:  "推送失败",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PUSH_STATUS_",
		Desc:         "游戏匹配-推送状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CALL_USER_PEOPLE_ALL",
		Value: 1,
		Desc:  "呼叫所有人",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CALL_USER_PEOPLE_GROUP",
		Value: 2,
		Desc:  "按照<小组>呼叫",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CALL_USER_PEOPLE_PROVIDE",
		Value: 3,
		Desc:  "用户自己指定呼叫的人",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CALL_USER_PEOPLE_",
		Desc:         "呼叫类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "RTC_ROOM_STATUS_CALLING",
		Value: 1,
		Desc:  "房间状态：呼叫中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RTC_ROOM_STATUS_EXECING",
		Value: 2,
		Desc:  "房间状态：运行中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RTC_ROOM_STATUS_END",
		Value: 3,
		Desc:  "房间状态：已结束",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "RTC_ROOM_STATUS_",
		Desc:         "房间状态2",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "RTC_PUSH_MSG_EVENT_FD_CREATE_REPEAT",
		Value: 400,
		Desc:  "FD 重复",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RTC_PUSH_MSG_EVENT_UID_NOT_IN_MAP",
		Value: 401,
		Desc:  " uid 不在MAP中",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "RTC_PUSH_MSG_EVENT_",
		Desc:         "RTC推送消息类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "REQ_SERVICE_METHOD_HTTP",
		Value: 1,
		Desc:  "http",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "REQ_SERVICE_METHOD_GRPC",
		Value: 2,
		Desc:  "grpc",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "REQ_SERVICE_METHOD_NATIVE",
		Value: 3,
		Desc:  "本地",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "REQ_SERVICE_METHOD_",
		Desc:         "请求3方服务的协议方法",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "DATA_ENCRYPT_BASE64",
		Value: 1,
		Desc:  "",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "DATA_ENCRYPT_AES_CBC_BASE64",
		Value: 2,
		Desc:  "",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "DATA_ENCRYPT_",
		Desc:         "加密方式",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_ROLE_NORMAL",
		Value: 1,
		Desc:  "普通用户",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_ROLE_DOCTOR",
		Value: 2,
		Desc:  "专家",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_ROLE_",
		Desc:         "用户角色",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PURPOSE_REGISTER",
		Value: 11,
		Desc:  "注册",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PURPOSE_LOGIN",
		Value: 12,
		Desc:  "登陆",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PURPOSE_FIND_BACK_PASSWORD",
		Value: 13,
		Desc:  "找回密码",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PURPOSE_SET_PASSWORD",
		Value: 21,
		Desc:  "设置密码",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PURPOSE_SET_MOBILE",
		Value: 22,
		Desc:  "设置手机号",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PURPOSE_SET_EMAIL",
		Value: 23,
		Desc:  "设置邮件",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PURPOSE_SET_PAY_PASSWORD",
		Value: 24,
		Desc:  "设置支付密码",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PURPOSE_",
		Desc:         "短信发送目的",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PLATFORM_MAC_PC_BROWSER",
		Value: 11,
		Desc:  "MAC台式浏览器",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_MAC_APP",
		Value: 12,
		Desc:  "MAC台式APP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_WIN_PC_BROWSER",
		Value: 22,
		Desc:  "WIN台式浏览器",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_WIN_APP",
		Value: 23,
		Desc:  "WIN台式APP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_ANDROID_H5_BROWSER",
		Value: 31,
		Desc:  "安卓手机浏览器",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_ANDROID_APP",
		Value: 32,
		Desc:  "安卓手机APP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_IOS_H5_BROWSER",
		Value: 41,
		Desc:  "IOS手机浏览器",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_IOS_APP",
		Value: 42,
		Desc:  "IOS手机APP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_AR",
		Value: 51,
		Desc:  "ar眼镜",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLATFORM_UNKNOW",
		Value: 99,
		Desc:  "未知",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PLATFORM_",
		Desc:         "平台类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "AUTH_CODE_STATUS_NORMAL",
		Value: 1,
		Desc:  "正常",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AUTH_CODE_STATUS_EXPIRE",
		Value: 3,
		Desc:  "失效",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AUTH_CODE_STATUS_OK",
		Value: 2,
		Desc:  "已使用",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "AUTH_CODE_STATUS_",
		Desc:         "验证码状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SEND_MSG_STATUS_OK",
		Value: 1,
		Desc:  "成功",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SEND_MSG_STATUS_FAIL",
		Value: 2,
		Desc:  "失败",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SEND_MSG_STATUS_ING",
		Value: 3,
		Desc:  "发送中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SEND_MSG_STATUS_WAIT",
		Value: 4,
		Desc:  "等待发送",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SEND_MSG_STATUS_",
		Desc:         "发送消息(sms 邮件 站内信)状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_FACEBOOK",
		Value: 21,
		Desc:  "facebook",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_GOOGLE",
		Value: 22,
		Desc:  "google",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_TWITTER",
		Value: 23,
		Desc:  "twitter",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_YOUTOBE",
		Value: 24,
		Desc:  "youtobe",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_WEIBO",
		Value: 11,
		Desc:  "微博",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_WECHAT",
		Value: 12,
		Desc:  "微信",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TYPE_THIRD_QQ",
		Value: 13,
		Desc:  "QQ",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_TYPE_THIRD_",
		Desc:         "用户注册类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_REG_TYPE_EMAIL",
		Value: 1,
		Desc:  "邮件",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_REG_TYPE_NAME",
		Value: 2,
		Desc:  "用户名",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_REG_TYPE_MOBILE",
		Value: 3,
		Desc:  "手机",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_REG_TYPE_THIRD",
		Value: 4,
		Desc:  "3方平台",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_REG_TYPE_GUEST",
		Value: 5,
		Desc:  "游客",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_REG_TYPE_",
		Desc:         "用户注册类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SEX_MALE",
		Value: 1,
		Desc:  "男",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SEX_FEMALE",
		Value: 2,
		Desc:  "女",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SEX_",
		Desc:         "性别",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_STATUS_NOMAL",
		Value: 1,
		Desc:  "正常",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_STATUS_DENY",
		Value: 2,
		Desc:  "禁止",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_STATUS_",
		Desc:         "用户状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_GUEST_TRUE",
		Value: 1,
		Desc:  "是",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_GUEST_FALSE",
		Value: 2,
		Desc:  "否",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_GUEST_",
		Desc:         "用户是否为游客",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_ROBOT_TRUE",
		Value: 1,
		Desc:  "是",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_ROBOT_FALSE",
		Value: 2,
		Desc:  "否",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_ROBOT_",
		Desc:         "用户是否为机器人",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "USER_TEST_TRUE",
		Value: 1,
		Desc:  "是",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "USER_TEST_FALSE",
		Value: 2,
		Desc:  "否",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "USER_TEST_",
		Desc:         "用户是否为测试账号",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "RULE_TYPE_AUTH_CODE",
		Value: 1,
		Desc:  "验证码",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RULE_TYPE_NOTIFY",
		Value: 2,
		Desc:  "通知",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "RULE_TYPE_MAKE",
		Value: 3,
		Desc:  "市场营销",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "RULE_TYPE_",
		Desc:         "消息rule类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SMS_CHANNEL_ALI",
		Value: 1,
		Desc:  "阿里",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SMS_CHANNEL_TENCENT",
		Value: 2,
		Desc:  "腾讯",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SMS_CHANNEL_",
		Desc:         "短信3方平台",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SERVER_PLATFORM_SELF",
		Value: 1,
		Desc:  "阿里",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVER_PLATFORM_TENGCENT",
		Value: 2,
		Desc:  "腾讯",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVER_PLATFORM_ALI",
		Value: 3,
		Desc:  "阿里",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVER_PLATFORM_HUAWEI",
		Value: 4,
		Desc:  "华为",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SERVER_PLATFORM_",
		Desc:         "服务器3方平台",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CICD_PUBLISH_DEPLOY_STATUS_ING",
		Value: 1,
		Desc:  "发布中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CICD_PUBLISH_DEPLOY_STATUS_FAIL",
		Value: 2,
		Desc:  "发布失败",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CICD_PUBLISH_DEPLOY_STATUS_FINISH",
		Value: 3,
		Desc:  "发布结束/完",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CICD_PUBLISH_DEPLOY_STATUS_",
		Desc:         "CICD部署状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CICD_PUBLISH_STATUS_WAIT_DEPLOY",
		Value: 1,
		Desc:  "待部署",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CICD_PUBLISH_STATUS_WAIT_PUB",
		Value: 2,
		Desc:  "待发布",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CICD_PUBLISH_DEPLOY_OK",
		Value: 3,
		Desc:  "成功",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CICD_PUBLISH_DEPLOY_FAIL",
		Value: 4,
		Desc:  "失败",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CICD_PUBLISH_",
		Desc:         "CICD发布状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PROJECT_STATUS_OPEN",
		Value: 1,
		Desc:  "打开",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_STATUS_CLOSE",
		Value: 2,
		Desc:  "关闭",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PROJECT_STATUS_",
		Desc:         "项目状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PROJECT_TYPE_SERVICE",
		Value: 1,
		Desc:  "服务",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_TYPE_FE",
		Value: 2,
		Desc:  "前端",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_TYPE_APP",
		Value: 3,
		Desc:  "APP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_TYPE_BE",
		Value: 4,
		Desc:  "后端",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PROJECT_TYPE_",
		Desc:         "项目大类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PROJECT_LANG_PHP",
		Value: 1,
		Desc:  "PHP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_LANG_GO",
		Value: 2,
		Desc:  "GO",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_LANG_JAVA",
		Value: 3,
		Desc:  "JAVA",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_LANG_JS",
		Value: 4,
		Desc:  "JS",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_LANG_C_PLUS",
		Value: 5,
		Desc:  "C++",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_LANG_C",
		Value: 6,
		Desc:  "C",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROJECT_LANG_C_SHARP",
		Value: 7,
		Desc:  "C#",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PROJECT_LANG_",
		Desc:         "项目开发语言类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STATUS_RESOURCE",
		Value: 1,
		Desc:  "已获取资源ID",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STATUS_START",
		Value: 2,
		Desc:  "已开始",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STATUS_END",
		Value: 3,
		Desc:  "已结束",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "AGORA_CLOUD_RECORD_STATUS_",
		Desc:         "声网录制屏幕-获取资源的状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STOP_ACTION_UNKNOW",
		Value: 0,
		Desc:  "未知",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STOP_ACTION_NORMAL",
		Value: 1,
		Desc:  "正常停止",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STOP_ACTION_RELOAD",
		Value: 2,
		Desc:  "页面刷新时拦截",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STOP_ACTION_REENTER",
		Value: 3,
		Desc:  "重新加载页面触发",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_STOP_ACTION_CALLBACK",
		Value: 4,
		Desc:  "声网回调触发",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "AGORA_CLOUD_RECORD_STOP_ACTION_",
		Desc:         "声网录制屏幕-停止的类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_SERVER_STATUS_UNDO",
		Value: 1,
		Desc:  "未处理",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_SERVER_STATUS_ING",
		Value: 2,
		Desc:  "处理中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_SERVER_STATUS_OK",
		Value: 3,
		Desc:  "处理成功",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "AGORA_CLOUD_RECORD_SERVER_STATUS_ERR",
		Value: 4,
		Desc:  "处理异常",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "AGORA_CLOUD_RECORD_SERVER_STATUS_",
		Desc:         "声网录制屏幕-服务状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CallbackEventAllUploaded",
		Value: 31,
		Desc:  " 所有录制文件已上传至指定的第三方云存储",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CallbackEventRecordExit",
		Value: 41,
		Desc:  " 录制服务已退出",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CallbackEvent",
		Desc:         "声网录制屏幕-回调状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PAY_TYPE_ALI",
		Value: 1,
		Desc:  "阿里",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PAY_TYPE_WECHAT",
		Value: 2,
		Desc:  "微信",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PAY_TYPE_",
		Desc:         "支付类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PAY_SUB_TYPE_APP",
		Value: 1,
		Desc:  "APP内",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PAY_SUB_TYPE_PC",
		Value: 2,
		Desc:  "浏览器网页端",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PAY_SUB_TYPE_H5",
		Value: 3,
		Desc:  "浏览器H5",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PAY_SUB_TYPE_QR",
		Value: 4,
		Desc:  "二维码",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PAY_SUB_TYPE_",
		Desc:         "支付子类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "GOODS_STATUS_NORMAL",
		Value: 1,
		Desc:  "上线",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "GOODS_STATUS_OFFLINE",
		Value: 2,
		Desc:  "下线",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "GOODS_STATUS_",
		Desc:         "商品状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "ORDERS_STATUS_NORMAL",
		Value: 1,
		Desc:  "正常，未支付",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ORDERS_STATUS_PAID",
		Value: 2,
		Desc:  "已支付",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ORDERS_STATUS_REFUND",
		Value: 3,
		Desc:  "退款",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ORDERS_STATUS_TIMEOUT",
		Value: 4,
		Desc:  "超时",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "ORDERS_STATUS_",
		Desc:         "订单状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "ENV_LOCAL_INT",
		Value: 1,
		Desc:  "开发环境",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ENV_DEV_INT",
		Value: 2,
		Desc:  "开发环境",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ENV_TEST_INT",
		Value: 3,
		Desc:  "测试环境",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ENV_PRE_INT",
		Value: 4,
		Desc:  "预发布环境",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ENV_ONLINE_INT",
		Value: 5,
		Desc:  "线上环境",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "ENV_",
		Desc:         "环境变量-整形",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "LOG_LEVEL_DEBUG",
		Value: 1,
		Desc:  "调试",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LOG_LEVEL_INFO",
		Value: 2,
		Desc:  "信息",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LOG_LEVEL_OFF",
		Value: 4,
		Desc:  "关闭",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "LOG_LEVEL_",
		Desc:         "error",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "LEVEL_INFO",
		Value: 1,
		Desc:  "",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_DEBUG",
		Value: 2,
		Desc:  "2",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_ERROR",
		Value: 4,
		Desc:  "4",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_PANIC",
		Value: 8,
		Desc:  "8",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_EMERGENCY",
		Value: 16,
		Desc:  "16",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_ALERT",
		Value: 32,
		Desc:  "32",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_CRITICAL",
		Value: 64,
		Desc:  "64",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_WARNING",
		Value: 128,
		Desc:  "128",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_NOTICE",
		Value: 256,
		Desc:  "256",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_TRACE",
		Value: 512,
		Desc:  "512",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_ALL",
		Value: 0,
		Desc:  "",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_DEV",
		Value: 0,
		Desc:  "",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LEVEL_ONLINE",
		Value: 0,
		Desc:  "",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "LEVEL_",
		Desc:         "日志等级",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PLAYER_STATUS_ONLINE",
		Value: 1,
		Desc:  "在线",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYER_STATUS_OFFLINE",
		Value: 2,
		Desc:  "离线",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PLAYER_STATUS_O",
		Desc:         "玩家当着在线状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PLAYERS_ACK_STATUS_INIT",
		Value: 1,
		Desc:  "初始化",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYERS_ACK_STATUS_WAIT",
		Value: 2,
		Desc:  "等待玩家确认",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PLAYERS_ACK_STATUS_OK",
		Value: 3,
		Desc:  "所有玩家均已确认",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PLAYERS_ACK_STATUS_",
		Desc:         "一个副本的，一条消息的，同步状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "ROOM_STATUS_INIT",
		Value: 1,
		Desc:  "新房间，刚刚初始化，等待其它操作",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_EXECING",
		Value: 2,
		Desc:  "已开始游戏",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_END",
		Value: 3,
		Desc:  "已结束",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_READY",
		Value: 4,
		Desc:  "准备中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "ROOM_STATUS_PAUSE",
		Value: 5,
		Desc:  "有玩家掉线，暂停中",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "ROOM_STATUS_",
		Desc:         "房间状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "PROTOCOL_TCP",
		Value: 1,
		Desc:  "传输协议 TCP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROTOCOL_UDP",
		Value: 3,
		Desc:  "传输协议 UDP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "PROTOCOL_WEBSOCKET",
		Value: 2,
		Desc:  "传输协议 WEB-SOCKET",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "PROTOCOL_",
		Desc:         "协议类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CONN_STATUS_INIT",
		Value: 1,
		Desc:  "初始化",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CONN_STATUS_EXECING",
		Value: 2,
		Desc:  "运行中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CONN_STATUS_CLOSE",
		Value: 3,
		Desc:  "已关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CONN_STATUS_CLOSE_ING",
		Value: 4,
		Desc:  "关闭中，防止重复关闭，不能用锁，因为：并发变串行后，还能重复关闭",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CONN_STATUS_",
		Desc:         "长连接connFD的状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CONTENT_TYPE_JSON",
		Value: 1,
		Desc:  "内容类型 json",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CONTENT_TYPE_PROTOBUF",
		Value: 2,
		Desc:  "proto_buf",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CONTENT_TYPE_",
		Desc:         "传输类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "METRICS_OPT_PLUS",
		Value: 1,
		Desc:  "1累加",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "METRICS_OPT_INC",
		Value: 2,
		Desc:  "2加加",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "METRICS_OPT_LESS",
		Value: 3,
		Desc:  "3累减",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "METRICS_OPT_DIM",
		Value: 4,
		Desc:  "4减减",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "METRICS_OPT_",
		Desc:         "metricsc操作类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "NETWAY_STATUS_INIT",
		Value: 1,
		Desc:  "网关状态 初始化中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "NETWAY_STATUS_START",
		Value: 2,
		Desc:  "网关状态 开始初始化",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "NETWAY_STATUS_CLOSE",
		Value: 3,
		Desc:  "网关状态 已关闭",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "NETWAY_STATUS_",
		Desc:         "NETWAY类状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "TRAN_MESSAGE_TYPE_CHAR",
		Value: 1,
		Desc:  "网络传输数据格式：字符流",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "TRAN_MESSAGE_TYPE_BINARY",
		Value: 2,
		Desc:  "网络传输数据格式：二进制",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "TRAN_MESSAGE_TYPE_",
		Desc:         "xxxx",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "UPLOAD_STORE_OSS_OFF",
		Value: 0,
		Desc:  "关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "UPLOAD_STORE_OSS_ALI",
		Value: 1,
		Desc:  "阿里",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "UPLOAD_STORE_OSS_",
		Desc:         "是否上传文件同时存储OSS",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "UPLOAD_STORE_LOCAL_OFF",
		Value: 1,
		Desc:  "关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "UPLOAD_STORE_LOCAL_OPEN",
		Value: 2,
		Desc:  "打开",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "UPLOAD_STORE_LOCAL_O",
		Desc:         "是否上传文件同时存储本地",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "FILE_TYPE_ALL",
		Value: 1,
		Desc:  "全部",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_TYPE_IMG",
		Value: 2,
		Desc:  "图片",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_TYPE_DOC",
		Value: 3,
		Desc:  "文档",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_TYPE_VIDEO",
		Value: 4,
		Desc:  "视频",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_TYPE_AUDIO",
		Value: 5,
		Desc:  "音频",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_TYPE_PACKAGES",
		Value: 6,
		Desc:  "压缩包",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "FILE_TYPE_",
		Desc:         "文件类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "FILE_SYNC_TRUE",
		Value: 1,
		Desc:  "全部",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_SYNC_FALSE",
		Value: 2,
		Desc:  "图片",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "FILE_SYNC_",
		Desc:         "文件操作是否同步到OSS",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_CREATE",
		Value: 3,
		Desc:  "初始化 连接类失败，可能是连接数过大",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_SERVER_HAS_CLOSE",
		Value: 11,
		Desc:  "服务端状态已关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_CLIENT",
		Value: 1,
		Desc:  "客户端-主动断开连接",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_AUTH_FAILED",
		Value: 21,
		Desc:  "客户端首次连接，登陆动作,服务端验证失败",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_FD_READ_EMPTY",
		Value: 22,
		Desc:  "客户端首次连接，登陆动作,服务端read信息为空",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_FD_PARSE_CONTENT",
		Value: 23,
		Desc:  "客户端首次连接，登陆动作,解析内容时出错",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_FIRST_NO_LOGIN",
		Value: 24,
		Desc:  "客户端首次连接，登陆动作,内容解出来了，但是action!=login",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_FIRST_PARSER_LOGIN",
		Value: 25,
		Desc:  "login  登陆不出结构体内容",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_OPEN_PANIC",
		Value: 31,
		Desc:  "初始化 新连接创建成功后，上层要再重新做一次连接，结果未知panic",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_MAX_CLIENT",
		Value: 32,
		Desc:  "当前连接数过大",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_OVERRIDE",
		Value: 4,
		Desc:  "创建新连接时，发现，该用户还有一个未关闭的连接,kickoff模式下，这条就没意义了",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_TIMEOUT",
		Value: 5,
		Desc:  "最后更新时间 ，超时.后台守护协程触发",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_SIGNAL_QUIT",
		Value: 6,
		Desc:  "接收到关闭信号，netWay.Quit触发",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_CLIENT_WS_FD_GONE",
		Value: 7,
		Desc:  "S端读取连接消息时，异常了~可能是：客户端关闭了连接",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_SEND_MESSAGE",
		Value: 8,
		Desc:  "S端给某个连接发消息，结果失败了，这里概率是连接已经断了",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_CONN_RESET_BY_PEER",
		Value: 81,
		Desc:  "对端，如果直接关闭网络，或者崩溃之类的，类库捕捉不到这个事件",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_CONN_SHUTDOWN",
		Value: 12,
		Desc:  "conn 已关闭",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_CONN_LOGIN_ROUTER_ERR",
		Value: 13,
		Desc:  "登陆，路由一个方法时，未找到该方法",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_RTT_TIMEOUT",
		Value: 91,
		Desc:  "S端已收到了RTT的响应，但已超时",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "CLOSE_SOURCE_RTT_TIMER_OUT",
		Value: 92,
		Desc:  "RTT超时，定时器触发",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "CLOSE_SOURCE_",
		Desc:         "长连接FD关闭类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "HTTP_DATA_CONTENT_TYPE_JSON",
		Value: 1,
		Desc:  "JSON",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "HTTP_DATA_CONTENT_TYPE_Nornal",
		Value: 2,
		Desc:  "普通",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "HTTP_DATA_CONTENT_TYPE_",
		Desc:         "http-curl类，数据传输类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "FILE_HASH_NONE",
		Value: 0,
		Desc:  " 没有",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_HASH_MONTH",
		Value: 1,
		Desc:  " 月",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_HASH_DAY",
		Value: 2,
		Desc:  "天",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "FILE_HASH_HOUR",
		Value: 3,
		Desc:  "小时",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "FILE_HASH_",
		Desc:         "文件存储hash类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SERVER_STATUS_NORMAL",
		Value: 1,
		Desc:  "正常",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVER_STATUS_CLOSE",
		Value: 2,
		Desc:  "已关闭",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SERVER_STATUS_",
		Desc:         "SERVER_STATUS",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SERVER_PING_OK",
		Value: 1,
		Desc:  "正常：PING 成功",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVER_PING_FAIL",
		Value: 2,
		Desc:  "异常：PING 失败了",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SERVER_PING_",
		Desc:         "SERVER_PING服务器的状态",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SERVICE_PROTOCOL_HTTP",
		Value: 1,
		Desc:  "HTTP",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVICE_PROTOCOL_GRPC",
		Value: 2,
		Desc:  "GRPC",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVICE_PROTOCOL_WEBSOCKET",
		Value: 3,
		Desc:  "WS",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVICE_PROTOCOL_TCP",
		Value: 4,
		Desc:  "TCP",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SERVICE_PROTOCOL_",
		Desc:         "service协议类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SERVICE_DISCOVERY_ETCD",
		Value: 1,
		Desc:  "ETCD",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SERVICE_DISCOVERY_CONSUL",
		Value: 2,
		Desc:  "CONSULE",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SERVICE_DISCOVERY_",
		Desc:         "服务发现的类型，分布式DB",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "LOAD_BALANCE_ROBIN",
		Value: 1,
		Desc:  "轮询",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "LOAD_BALANCE_HASH",
		Value: 2,
		Desc:  "固定分子hash",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "LOAD_BALANCE_",
		Desc:         "服务发现负载类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}

    	constItem = ConstItem{
		Key:   "SV_ERROR_NONE",
		Value: 0,
		Desc:  "无",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SV_ERROR_INIT",
		Value: 1,
		Desc:  "初始化",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SV_ERROR_CONN",
		Value: 2,
		Desc:  "连接中",
	}
	constItemList = append(constItemList, constItem)

	constItem = ConstItem{
		Key:   "SV_ERROR_NOT_FOUND",
		Value: 3,
		Desc:  "未找到",
	}
	constItemList = append(constItemList, constItem)


	enumConst = EnumConst{
		CommonPrefix: "SV_ERROR_",
		Desc:         "super_visor错误类型",
		ConstList:    constItemList,
		Type:          "int",
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}


}
