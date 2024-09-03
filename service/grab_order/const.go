package grab_order

const (
	PAY_CATEGORY_ALI       = 1 //阿里
	PAY_CATEGORY_BANK_CARD = 2 //银行卡
	PAY_CATEGORY_WECHAT    = 3 //微信

	DEFAULT_PRIORITY = 1 //抢单，用户默认权重值

	EVENT_TYPE_USER_COLOSE_GRAB       = 1 //用户关闭了抢单
	EVENT_TYPE_USER_WS_CLOSE          = 2 //用户ws断了
	EVENT_TYPE_USER_FREEZE            = 3 //用户账户被冻结
	EVENT_TYPE_USER_PAY_CHANNEL_CLOSE = 4 //用户支付渠道已关闭
	EVENT_TYPE_PAY_CATEGORY_CLOSE     = 5 //系统支付分类，已关闭
	EVENT_TYPE_SET_CHANGE             = 6 //配置文件修改
)

func PayCategoryList() map[int]string {
	list := make(map[int]string)
	list[PAY_CATEGORY_ALI] = "阿里"
	list[PAY_CATEGORY_BANK_CARD] = "银行卡"
	list[PAY_CATEGORY_WECHAT] = "微信"

	return list
}
