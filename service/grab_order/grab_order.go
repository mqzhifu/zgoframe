package grab_order

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

type GrabOrder struct {
	OrderBucketList           map[int]*OrderBucket           `json:"order_bucket_list"`             //[category_id]OrderBucket 存放所有订单
	UserBucketAmountRangeList map[int]map[string]*UserBucket `json:"user_bucket_amount_range_list"` //[category_id][amount_range]UserBucket 存放抢单的人
	UserTotal                 *UserTotal                     `json:"user_total"`                    //每个用户的详细信息
	AmountRange               AmountRange                    `json:"amount_range"`                  //存储 settings 的配置信息
	Settings                  Settings                       `json:"settings"`                      //配置信息 - 各种限制
	Gorm                      *gorm.DB                       `json:"-"`                             //DB
	Redis                     *util.MyRedis                  `json:"-"`                             //redis
	EventMsgCh                chan EventMsg                  `json:"-"`                             //接收 eventMsg 中断抢单死循环
	CloseOrderTimeoutDemon    chan int                       `json:"-"`
}

func NewGrabOrder(db *gorm.DB, redis *util.MyRedis) *GrabOrder {
	grabOrder := new(GrabOrder)
	grabOrder.Redis = redis
	grabOrder.Gorm = db
	grabOrder.EventMsgCh = make(chan EventMsg)
	grabOrder.CloseOrderTimeoutDemon = make(chan int)

	grabOrder.Init()

	return grabOrder
}

func (grabOrder *GrabOrder) Init() {
	grabOrder.InitAmountRange()
	grabOrder.InitSettings()
	grabOrder.InitUserBucketOrderBucket()

	grabOrder.UserTotal = NewUserTotal(grabOrder.Redis)

}

// 初始化 - 装载 抢单金额区间 池子
func (grabOrder *GrabOrder) InitAmountRange() {
	amountRange := AmountRange{
		MinAmount: 100,
		MaxAmount: 5000,
	}

	amountRange.Range = append(amountRange.Range, AmountRangeElement{MinAmount: 100, MaxAmount: 500})
	amountRange.Range = append(amountRange.Range, AmountRangeElement{MinAmount: 501, MaxAmount: 1000})
	amountRange.Range = append(amountRange.Range, AmountRangeElement{MinAmount: 1001, MaxAmount: 5000})

	grabOrder.AmountRange = amountRange
}

// 初始化 - 配置信息，各种限制值
func (grabOrder *GrabOrder) InitSettings() {
	settings := Settings{
		GrabTimeout:        30,
		GrabDayTotalAmount: 10000,
		GrabDayOrderCnt:    50,
		GrabIntervalTime:   60,
	}
	grabOrder.Settings = settings
}

// 初始化 - 抢单的用户池子 - [支付渠道类型][金额区间][用户池]
func (grabOrder *GrabOrder) InitUserBucketOrderBucket() {
	grabOrder.UserBucketAmountRangeList = make(map[int]map[string]*UserBucket)
	grabOrder.OrderBucketList = make(map[int]*OrderBucket)
	for categoryId, _ := range PayCategoryList() { //每个支付分类 - 有一个 订单桶
		grabOrder.OrderBucketList[categoryId] = NewOrderBucket(categoryId)
		mapUserBucket := make(map[string]*UserBucket)
		for _, v := range grabOrder.AmountRange.Range {
			key := GetRangeKey(v.MinAmount, v.MaxAmount)
			mapUserBucket[key] = NewUserBucket(grabOrder.Redis, categoryId, v.MinAmount, v.MaxAmount)
		}
		grabOrder.UserBucketAmountRangeList[categoryId] = mapUserBucket
	}
}

// 监听：订单超时
func (grabOrder *GrabOrder) CheckOrderTimeout() {
	stop := 0
	for {
		select {
		case <-grabOrder.CloseOrderTimeoutDemon:
			stop = 1
			break
		default:
			time.Sleep(time.Second * 1)
		}
		if stop == 1 {
			break
		}
	}
}

// 用户开启 - 自动抢单功能
func (grabOrder *GrabOrder) UserOpenGrab(uid int, userOpenGrabSet []request.GrabOrderUserOpen) error {
	if uid <= 0 || len(userOpenGrabSet) == 0 {
		return errors.New("uid <= 0 or len UserOpenGrabSet== 0")
	}
	fmt.Println("userOpenGrabSet:", userOpenGrabSet)
	//200 - 2000
	//100 - 500 | 501 - 1000  | 1000 - 5000
	userInAmtFlag := 0
	for _, userChannel := range userOpenGrabSet {
		fmt.Println(userChannel.AmountMin, grabOrder.AmountRange.MinAmount)
		if userChannel.AmountMin < grabOrder.AmountRange.MinAmount {
			fmt.Println("此用户支付渠道不满足 userChannel.AmountMin < grabOrder.AmountRange.MinAmount")
			continue
		}

		if userChannel.AmountMax > grabOrder.AmountRange.MaxAmount {
			fmt.Println("此用户支付渠道不满足 userChannel.AmountMax > grabOrder.AmountRange.MaxAmount")
			continue
		}
		//每个渠道，根据每个金额范围，都会有一个用户桶
		//用户：每个渠道，会设置一个 金额范围
		//一个渠道  => 多个金额区间的用户池 ，用户可能同时在多个用户池中
		//将用户设置的金额区间，计算出应该在哪些用户池中添加该UID
		userAmountSection := []string{}                 //存储，用户金额范围内，属于哪些区间
		for _, v := range grabOrder.AmountRange.Range { //100 - 500 | 501 - 1000  | 1000 - 5000
			if userChannel.AmountMin >= v.MaxAmount {
				continue
			}
			userAmountSection = append(userAmountSection, GetRangeKey(v.MinAmount, v.MaxAmount))
			if userChannel.AmountMax <= v.MaxAmount {
				break
			}
		}

		fmt.Println("userAmountSection:", userAmountSection)

		if len(userAmountSection) == 0 {
			continue
		}
		userInAmtFlag = 1

		for _, v := range userAmountSection {
			//找到该金额区间的用户池，然后添加进去
			userBucket := grabOrder.UserBucketAmountRangeList[userChannel.PayCategoryId][v]
			queueItem := QueueItem{
				Score: DEFAULT_PRIORITY,
				Uid:   uid,
			}
			//可能用户池的UID会重复，但用序队列会自动覆盖
			userBucket.QueueRedis.Push(queueItem)
		}
	}
	fmt.Println("=======111111")
	if userInAmtFlag == 0 {
		errMsg := "用户的所有渠道的所有金额区间都没有匹配到，算是异常"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}
	//UID会在多个渠道、金额区间的池子里，还得有个用户的总控，在这里添加/更新
	err, optType := grabOrder.UserTotal.AddOrUpdateOne(uid)
	fmt.Println("AddOrUpdateOne rs , err:", err, " , optType:", optType)

	return nil
}

// 用户关闭抢单
func (grabOrder *GrabOrder) UserCloseGrab(uid int) error {
	userTotal, exist := grabOrder.UserTotal.GetOne(uid)
	if !exist {
		return errors.New("user total does not exist")
	}

	userTotal.GrabStatus = USER_GRAP_STATUS_CLOSE
	for categoryId, _ := range PayCategoryList() { //每个支付分类 - 有一个 订单桶
		grabOrder.OrderBucketList[categoryId] = NewOrderBucket(categoryId)
		for _, v := range grabOrder.AmountRange.Range {
			key := GetRangeKey(v.MinAmount, v.MaxAmount)
			userBucket := grabOrder.UserBucketAmountRangeList[categoryId][key]
			userBucket.QueueRedis.DelOneByUid(uid)
		}
	}
	return nil
}

// 其它服务，推单过来
func (grabOrder *GrabOrder) CreateOrder(req request.GrabOrder) (error, int) {
	order := Order{
		Id:         req.Id,
		Amount:     req.Amount,
		CategoryId: req.CategoryId,
		Uid:        req.Uid,
	}
	//先给订单 设置 超时时间
	order.StartTime = int(time.Now().Unix())
	order.Timeout = order.StartTime + grabOrder.Settings.GrabTimeout
	key := ""
	//判断 - 当前订单的金额 属性哪个 金额区间
	for _, v := range grabOrder.AmountRange.Range {
		if order.Amount >= v.MinAmount && order.Amount <= v.MaxAmount {
			key = GetRangeKey(v.MinAmount, v.MaxAmount)
		}
	}
	if key == "" {
		fmt.Println("err1 key empty")
		return errors.New("key empty"), 0
	}
	//订单 - 添加到公共的桶中
	grabOrder.OrderBucketList[order.CategoryId].AddOne(order)
	fmt.Println("grabOrder.OrderBucketList[order.CategoryId]:", grabOrder.OrderBucketList[order.CategoryId])
	timer := time.NewTimer(time.Second * time.Duration(order.StartTime+grabOrder.Settings.GrabTimeout))
	//根据：支付渠道类型、金额区间，找到那个用户池子，从池子找一个用户接单
	userBucket := grabOrder.UserBucketAmountRangeList[order.CategoryId][key]
	selectStatus := 0 //1超时2检查失败3成功
	popQueueList := []QueueItem{}
	successUser := 0
	for {
		select {
		case msg := <-grabOrder.EventMsgCh:
			fmt.Println(msg)
			selectStatus = LOOP_SELECT_USER_QUEUE_EVENT_STOP
			//exceptUid = append(exceptUid, msg.Uid)
		case <-timer.C:
			selectStatus = LOOP_SELECT_USER_QUEUE_TIMEOUT
			break
		default:
			if userBucket.QueueRedis.Len() == 0 {
				selectStatus = LOOP_SELECT_USER_QUEUE_EMPTY
				break
			}
			queueItem, err := userBucket.QueueRedis.Pop()
			if err != nil {
				fmt.Println("userBucket.PopOne err:", err.Error())
				continue
			}
			popQueueList = append(popQueueList, queueItem)
			fmt.Println("popQueueList:", popQueueList)
			err = grabOrder.GrabDoing(queueItem.Uid, order.Id, order.CategoryId)
			if err != nil {
				selectStatus = LOOP_SELECT_USER_QUEUE_FAILED
				break
			}
			selectStatus = LOOP_SELECT_USER_QUEUE_SUCCESS
			successUser = queueItem.Uid
		}
		//只有匹配失败的一种情况，才需要，一直循环
		if selectStatus != LOOP_SELECT_USER_QUEUE_FAILED {
			break
		}
	}
	//把之前弹出的用户，再放回去
	if len(popQueueList) > 0 {

	}

	switch selectStatus {
	case LOOP_SELECT_USER_QUEUE_SUCCESS:
		return nil, successUser
	case LOOP_SELECT_USER_QUEUE_EVENT_STOP:
	case LOOP_SELECT_USER_QUEUE_TIMEOUT:
	case LOOP_SELECT_USER_QUEUE_EMPTY:
	case LOOP_SELECT_USER_QUEUE_FAILED:

	}

	return errors.New("unknow"), 0

}

// 正式开始抢单
func (grabOrder *GrabOrder) GrabDoing(uid int, oid string, categoryId int) error {
	fmt.Println("GrabDoing uid:", uid, " categoryId:", categoryId, " oid:", oid)
	err := grabOrder.CheckGrabLimit(categoryId, oid, uid)
	grabOrder.UserTotal.UpdateLastGrabFailedTime()
	if err != nil {
		return err
	}
	//order record 落盘
	//更新用户：余额、冻结金额
	//合建账变记录

	grabOrder.UserTotal.UpdateGrabDayTotalOrderCnt()
	//grabOrder.UserTotal.UpdateGrabDayTotalAmountProgress()
	grabOrder.UserTotal.UpdateGrabDayTotalAmount()
	grabOrder.UserTotal.UpdateLastGrabSuccessTime()

	return nil
}

func (grabOrder *GrabOrder) CheckGrabLimit(category int, oid string, uid int) (err error) {
	order := grabOrder.OrderBucketList[category].GelOne(oid)
	if order.Timeout > int(time.Now().Unix()) {
		return errors.New("order timeout")
	}

	userTotal, exist := grabOrder.UserTotal.GetOne(uid)
	if exist == false {
		return errors.New("uid not in UserTotal")
	}
	if userTotal.UserDayTotal.GrabAmount > grabOrder.Settings.GrabDayTotalAmount {

	}

	if userTotal.UserDayTotal.GrabCnt > grabOrder.Settings.GrabDayOrderCnt {

	}

	if int(time.Now().Unix()-userTotal.LastGrabSuccessTime) < grabOrder.Settings.GrabIntervalTime {
		return errors.New("下单间隔")
	}

	if userTotal.WsStatus == USER_WS_STATUS_OFFLINE {
		return errors.New("用户长连接状态为：关闭")
	}

	if userTotal.GrabStatus == USER_GRAP_STATUS_CLOSE {
		return errors.New("uid not in UserTotal")
	}

	userInfo := model.User{}
	grabOrder.Gorm.First(&userInfo, uid)
	if userInfo.Id <= 0 {
		return errors.New("uid not in user table")
	}
	if userInfo.Status == 2 { //用户被禁用
		return errors.New("用户被禁用")
	}
	userTotalInfo := model.UserTotal{}
	grabOrder.Gorm.First(&userTotalInfo, uid)
	if userTotalInfo.Id <= 0 {
		return errors.New("uid not in userTotal table")
	}
	//可抢额度=账户余额-押金-冻结金额
	if order.Amount > userTotalInfo.Cash { //余额不足
		return errors.New("余额不足")
	}
	//if userTotal.PayCategoryStatus == 0 {
	//
	//}
	//
	//if userTotal.PayChannelStatus == 0 {
	//
	//}
	//

	return nil
}

func (grabOrder *GrabOrder) ReceiveEventMsg(msg EventMsg) {
	switch msg.TypeId {
	case EVENT_TYPE_USER_COLOSE_GRAB:
		break
	case EVENT_TYPE_USER_WS_CLOSE:
		break
	case EVENT_TYPE_USER_FREEZE:
		break
	case EVENT_TYPE_USER_PAY_CHANNEL_CLOSE:
		break
	case EVENT_TYPE_PAY_CATEGORY_CLOSE:
		break
	case EVENT_TYPE_SET_CHANGE:
		break
	default:
		fmt.Println("event type id err")
	}
}

// 给前端API
func (grabOrder *GrabOrder) GetPayCategory() ([]model.PayCategory, error) {
	list := []model.PayCategory{}
	grabOrder.Gorm.Find(&list)
	return list, nil
}

// 给前端API
func (grabOrder *GrabOrder) GetData() (*GrabOrder, error) {
	return grabOrder, nil
}
