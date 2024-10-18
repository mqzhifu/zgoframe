package grab_order

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
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

	grabOrder.UserTotal = NewUserTotal(grabOrder.Redis, grabOrder.Gorm)

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
		GrabTimeout:        5,
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
		grabOrder.OrderBucketList[categoryId] = NewOrderBucket(categoryId, grabOrder.Redis, grabOrder.Gorm)
		mapUserBucket := make(map[string]*UserBucket)
		for _, v := range grabOrder.AmountRange.Range {
			key := GetRangeKey(v.MinAmount, v.MaxAmount)
			mapUserBucket[key] = NewUserBucket(grabOrder.Redis, categoryId, v.MinAmount, v.MaxAmount)
		}
		grabOrder.UserBucketAmountRangeList[categoryId] = mapUserBucket
	}
}

// 用户开启 - 自动抢单功能
func (grabOrder *GrabOrder) UserOpenGrab(uid int, userOpenGrabSet []request.GrabOrderUserOpen) error {
	if uid <= 0 || len(userOpenGrabSet) == 0 {
		return errors.New("uid <= 0 or len UserOpenGrabSet== 0")
	}
	fmt.Println("userOpenGrabSet:", userOpenGrabSet)

	userTotalInfo, exist := grabOrder.UserTotal.GetOne(uid)
	if exist {
		err := grabOrder.CheckGrabLimit(uid, userTotalInfo)
		if err != nil {
			return err
		}
		if userTotalInfo.GrabStatus == 1 {
			return errors.New("用户已经开启了抢单，不要重复操作")
		}
	}

	maxAmount := userOpenGrabSet[0].AmountMax //找出用户配置的渠道中，最大的那个值
	//200 - 2000
	//100 - 500 | 501 - 1000  | 1000 - 5000
	userInAmtFlag := 0
	for _, userChannel := range userOpenGrabSet {
		fmt.Println(userChannel.AmountMin, grabOrder.AmountRange.MinAmount)
		if userChannel.AmountMax > maxAmount {
			maxAmount = userChannel.AmountMax
		}
		if userChannel.AmountMin < grabOrder.AmountRange.MinAmount {
			fmt.Println("用户支付渠道(" + strconv.Itoa(userChannel.PayCategoryId) + ")不满足 userChannel.AmountMin < grabOrder.AmountRange.MinAmount")
			continue
		}

		if userChannel.AmountMax > grabOrder.AmountRange.MaxAmount {
			fmt.Println("此用户支付渠道(" + strconv.Itoa(userChannel.PayCategoryId) + ")不满足 userChannel.AmountMax > grabOrder.AmountRange.MaxAmount")
			continue
		}

		userPayAccount := model.UserPayAccount{}
		grabOrder.Gorm.Where("uid = ? and category_id = ? ", uid, userChannel.PayCategoryId).First(&userPayAccount)
		if userPayAccount.Id <= 0 {
			return errors.New("该用户并没有：支付分类账户(" + strconv.Itoa(userChannel.PayCategoryId) + ")")
		}
		if userPayAccount.AmountMin > userChannel.AmountMin || userPayAccount.AmountMax < userChannel.AmountMax {
			return errors.New("用户支付账号：金额范围错误")
		}
		if userPayAccount.Status == 2 {
			return errors.New("用户支付账号:状态已关闭")
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
	if userInAmtFlag == 0 {
		errMsg := "用户的所有渠道的所有金额区间都没有匹配到，算是异常"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	userTotalDB := model.UserTotal{}
	grabOrder.Gorm.Where("uid = ?", uid).First(&userTotalDB)
	if userTotalDB.Id <= 0 {
		return errors.New("uid not in user_total table.")
	}
	//接单最大额度 ，不能大于 自己的账户额度，不然没法扣款了
	if maxAmount > userTotalDB.Cash {
		return errors.New("接单最大额度 ，不能大于 自己的账户额度，不然没法扣款了")
	}

	//UID会在多个渠道、金额区间的池子里，还得有个用户的总控，在这里添加/更新
	err, optType := grabOrder.UserTotal.AddOrUpdateOne(uid, userOpenGrabSet)
	fmt.Println("AddOrUpdateOne rs , err:", err, " , optType:", optType)

	return nil
}

// 用户关闭抢单
func (grabOrder *GrabOrder) UserCloseGrab(uid int) error {
	userTotal, exist := grabOrder.UserTotal.GetOne(uid)
	if !exist {
		return errors.New("user total does not exist")
	}

	userTotal.GrabStatus = USER_GRAB_STATUS_CLOSE
	for categoryId, _ := range PayCategoryList() { //每个支付分类 - 有一个 订单桶
		grabOrder.OrderBucketList[categoryId] = NewOrderBucket(categoryId, grabOrder.Redis, grabOrder.Gorm)
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
	order := Order{}
	order.InId = req.OrderId
	order.Amount = req.Amount
	order.CategoryId = req.CategoryId
	order.Uid = req.Uid
	order.Status = ORDER_MATCH_STATUS_ING
	order.StartTime = int(time.Now().Unix())
	order.PayStatus = 1

	order.Timeout = order.StartTime + grabOrder.Settings.GrabTimeout
	userBucketAmountRangeKey := ""
	//判断 - 当前订单的金额 属性哪个 金额区间
	for _, v := range grabOrder.AmountRange.Range {
		if order.Amount >= v.MinAmount && order.Amount <= v.MaxAmount {
			userBucketAmountRangeKey = GetRangeKey(v.MinAmount, v.MaxAmount)
		}
	}
	if userBucketAmountRangeKey == "" {
		fmt.Println("err1 key empty")
		return errors.New("key empty"), 0
	}

	fmt.Println("CreateOrder info id:", order.Id, ",amount:", order.Amount, ",uid:", order.Uid, "categoryId:", order.CategoryId, ",timeout:", order.Timeout, ",userBucketAmountRangeKey:", userBucketAmountRangeKey)

	//订单 - 添加到公共的桶中
	grabOrder.OrderBucketList[order.CategoryId].AddOne(order)
	//fmt.Println("grabOrder.OrderBucketList[order.CategoryId]:", grabOrder.OrderBucketList[order.CategoryId])
	timer := time.NewTimer(time.Second * time.Duration(grabOrder.Settings.GrabTimeout))
	//根据：支付渠道类型、金额区间，找到那个用户池子，从池子找一个用户接单
	userBucket := grabOrder.UserBucketAmountRangeList[order.CategoryId][userBucketAmountRangeKey]
	matchStatus := 0 //1超时2检查失败3成功
	popQueueUserList := []QueueItem{}
	successUser := 0
	order.MatchQueueUserCnt = userBucket.QueueRedis.Len() //记录下，开始匹配时，队列里的用户总数
	matchTimes := 0                                       //匹配次数
	for {
		matchTimes++
		fmt.Println("select match times:", matchTimes)
		select {
		case msg := <-grabOrder.EventMsgCh:
			fmt.Println(msg)
			matchStatus = ORDER_MATCH_STATUS_EVENT_STOP
		case <-timer.C:
			matchStatus = ORDER_MATCH_STATUS_TIMEOUT
			break
		default:
			fmt.Println("queue len:", userBucket.QueueRedis.Len())
			//为空也需要等待：超时，不排队其它协程拿走了数据
			if userBucket.QueueRedis.Len() == 0 {
				break
			}
			queueItem, err := userBucket.QueueRedis.Pop()
			if err != nil {
				fmt.Println("userBucket.PopOne err:", err.Error())
				continue
			}
			fmt.Println("pop one uid:", queueItem.Uid, ",score:", queueItem.Score)
			popQueueUserList = append(popQueueUserList, queueItem)
			err = grabOrder.CreateOrderCheckGrabLimit(order.CategoryId, order.InId, queueItem.Uid)
			if err != nil {
				grabOrder.UserTotal.UpGrabFailedTime(queueItem.Uid)
				fmt.Println("CheckGrabLimit err:", err)
				matchStatus = ORDER_MATCH_STATUS_FAILED
				break
			}
			matchStatus = ORDER_MATCH_STATUS_SUCCESS
			successUser = queueItem.Uid

		}
		fmt.Println("order match once status:"+strconv.Itoa(matchStatus), ", successUser:", successUser, ",matchTimes:", matchTimes)
		//只有匹配失败的一种情况，才需要，一直循环
		if matchStatus != ORDER_MATCH_STATUS_FAILED {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	//把之前弹出的用户，再放回去
	grabOrder.QueuePushBack(popQueueUserList, userBucket)

	order.Status = matchStatus
	order.GrabUid = successUser
	order.EndTime = int(time.Now().Unix())
	order.MatchTimes = matchTimes
	grabOrder.OrderBucketList[order.CategoryId].UpOneRecord(order)

	switch matchStatus {
	case ORDER_MATCH_STATUS_SUCCESS:
		grabOrder.GrabSuccessUpData(order, userBucket)
		return nil, successUser
	case ORDER_MATCH_STATUS_EVENT_STOP:
		return errors.New("EVENT_STOP"), 0
	case ORDER_MATCH_STATUS_TIMEOUT:
		return errors.New("order timeout"), 0
	case ORDER_MATCH_STATUS_FAILED:
		return errors.New("checkout user info failed"), 0
	}

	return errors.New("unknow"), 0

}

func (grabOrder *GrabOrder) QueuePushBack(popQueueUserList []QueueItem, userBucket *UserBucket) {
	if len(popQueueUserList) > 0 {
		for _, queueItem := range popQueueUserList {
			userBucket.QueueRedis.Push(queueItem)
		}
	}
}

func (grabOrder *GrabOrder) GrabSuccessUpData(order Order, userBucket *UserBucket) {
	//更新用户：余额、冻结金额
	userTotal := model.UserTotal{}
	upData := make(map[string]interface{})
	upData["cash"] = gorm.Expr("cash - ?", order.Amount)
	upData["FreezeCash"] = gorm.Expr("freezeCash + ?", order.Amount)
	grabOrder.Gorm.Model(&userTotal).Where("uid = ?", order.Uid).Updates(upData)
	//创建账变记录:冻结金额

	//抢单成功后，权重要下降
	userBucket.QueueRedis.IncScore(order.GrabUid, -1)

	grabOrder.UserTotal.UpGrabDayTotalOrderCnt(order.GrabUid)
	grabOrder.UserTotal.UpGrabDayTotalAmount(order.GrabUid, order.Amount)
	grabOrder.UserTotal.UpGrabSuccessTime(order.GrabUid)
	grabOrder.UserTotal.UpLastGrabSuccessTime(order.GrabUid)
}
func (grabOrder *GrabOrder) CheckUserTotal(uid int) {
	userTotalDB := model.UserTotal{}
	grabOrder.Gorm.Where("uid = ?", uid).First(&userTotalDB)
	if userTotalDB.Id <= 0 {

	}
	//userTotalDB.
}

func (grabOrder *GrabOrder) CreateOrderCheckGrabLimit(categoryId int, oid string, uid int) (err error) {
	userTotal, exist := grabOrder.UserTotal.GetOne(uid)
	if exist == false {
		return errors.New("uid not in UserTotal")
	}

	err = grabOrder.CheckGrabLimit(uid, userTotal)
	if err != nil {
		return err
	}

	order := grabOrder.OrderBucketList[categoryId].GelOne(oid)
	userTotalInfo := model.UserTotal{}
	grabOrder.Gorm.First(&userTotalInfo, uid)
	if userTotalInfo.Id <= 0 {
		return errors.New("uid not in userTotal table")
	}
	//可抢额度=账户余额-押金-冻结金额
	if order.Amount > userTotalInfo.Cash { //余额不足
		return errors.New("余额不足")
	}

	if userTotal.GrabStatus == USER_GRAB_STATUS_CLOSE {
		return errors.New("用户抢单状态：关闭")
	}

	return nil
}

func (grabOrder *GrabOrder) CheckGrabLimit(uid int, userTotal UserElement) (err error) {
	//当日 可抢数量
	if userTotal.UserDayTotal.GrabAmount > grabOrder.Settings.GrabDayTotalAmount {

	}
	//当日 可抢额度
	if userTotal.UserDayTotal.GrabCnt > grabOrder.Settings.GrabDayOrderCnt {

	}
	//下单间隔
	if int(time.Now().Unix()-userTotal.LastGrabSuccessTime) < grabOrder.Settings.GrabIntervalTime {
		return errors.New("下单间隔")
	}

	if userTotal.WsStatus == USER_WS_STATUS_OFFLINE {
		return errors.New("用户长连接状态为：关闭")
	}

	userInfo := model.User{}
	grabOrder.Gorm.First(&userInfo, uid)
	if userInfo.Id <= 0 {
		return errors.New("uid not in user table")
	}
	if userInfo.Status == 2 { //用户被禁用
		return errors.New("用户被禁用")
	}

	return nil
}

func (grabOrder *GrabOrder) ReceiveEventMsg(msg EventMsg) {
	switch msg.TypeId {
	case EVENT_TYPE_USER_WS_CHANGE:
		grabOrder.UserTotal.UpWs(msg.Uid, 1)
		break
	case EVENT_TYPE_USER_FREEZE:
		grabOrder.UserCloseGrab(msg.Uid)
		break
	case EVENT_TYPE_USER_PAY_CHANNEL_CLOSE:
		grabOrder.UserCloseGrab(msg.Uid)
		break
	case EVENT_TYPE_PAY_CATEGORY_CLOSE:
		grabOrder.UserCloseGrab(msg.Uid)
		break
	case EVENT_TYPE_SET_CHANGE:
		grabOrder.UserCloseGrab(msg.Uid)
		break
	default:
		fmt.Println("event type id err")
	}

	grabOrder.EventMsgCh <- msg
}
