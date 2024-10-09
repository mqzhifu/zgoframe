package grab_order

import "strconv"

const (
	PAY_CATEGORY_ALI       = 1 //阿里
	PAY_CATEGORY_BANK_CARD = 2 //银行卡
	PAY_CATEGORY_WECHAT    = 3 //微信

	DEFAULT_PRIORITY = 100 //抢单，用户默认权重值

	EVENT_TYPE_USER_COLOSE_GRAB       = 1 //用户关闭了抢单
	EVENT_TYPE_USER_WS_CLOSE          = 2 //用户ws断了
	EVENT_TYPE_USER_FREEZE            = 3 //用户账户被冻结
	EVENT_TYPE_USER_PAY_CHANNEL_CLOSE = 4 //用户支付渠道已关闭
	EVENT_TYPE_PAY_CATEGORY_CLOSE     = 5 //系统支付分类，已关闭
	EVENT_TYPE_SET_CHANGE             = 6 //配置文件修改

	USER_TOTAL_OPT_TYPE_ADD = 1
	USER_TOTAL_OPT_TYPE_UP  = 2

	USER_GRAP_STATUS_OPEN  = 1
	USER_GRAP_STATUS_CLOSE = 2

	USER_WS_STATUS_ONLINE  = 1
	USER_WS_STATUS_OFFLINE = 2

	//LOOP_SELECT_USER_QUEUE_TIMEOUT    = 1
	//LOOP_SELECT_USER_QUEUE_FAILED     = 2
	//LOOP_SELECT_USER_QUEUE_SUCCESS    = 3
	//LOOP_SELECT_USER_QUEUE_EVENT_STOP = 4

	ORDER_MATCH_STATUS_ING        = 1
	ORDER_MATCH_STATUS_SUCCESS    = 2
	ORDER_MATCH_STATUS_FAILED     = 3
	ORDER_MATCH_STATUS_TIMEOUT    = 4
	ORDER_MATCH_STATUS_EVENT_STOP = 5

	//LOOP_SELECT_USER_QUEUE_EMPTY      = 5
)

// 用户开启自动抢单 - 设置的金额(支付渠道)列表
type UserOpenGrabSet struct {
	PayCategoryId int
	AmountMin     int
	AmountMax     int
}

type EventMsg struct {
	OrderId int
	TypeId  int
	Content string
	Uid     int
}

// 存储 AmountRangeElement 的信息，增加了：最大值 和 最小值字段，方便运算
type AmountRange struct {
	MinAmount int                  `json:"min_amount"`
	MaxAmount int                  `json:"max_amount"`
	Range     []AmountRangeElement `json:"range"`
}

// 订单金额 - 区间段 - 划分
type AmountRangeElement struct {
	MinAmount int `json:"min_amount"`
	MaxAmount int `json:"max_amount"`
}

// 抢单 - 配置信息 - 各种限制
type Settings struct {
	GrabTimeout        int `json:"grab_timeout"`          //抢单超时时间
	GrabDayTotalAmount int `json:"grab_day_total_amount"` //每天可抢总额度
	GrabDayOrderCnt    int `json:"grab_day_order_cnt"`    //每天可抢总订单数量
	GrabIntervalTime   int `json:"grab_interval_time"`    //抢单间隔
}

func PayCategoryList() map[int]string {
	list := make(map[int]string)
	list[PAY_CATEGORY_ALI] = "阿里"
	list[PAY_CATEGORY_BANK_CARD] = "银行卡"
	list[PAY_CATEGORY_WECHAT] = "微信"

	return list
}

func GetRangeKey(amountMin int, amountMax int) string {
	return strconv.Itoa(amountMin) + "_" + strconv.Itoa(amountMax)
}
