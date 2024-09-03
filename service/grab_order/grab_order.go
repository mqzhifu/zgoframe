package grab_order

import (
	"fmt"
	"strconv"
	"time"
)

type GrabOrder struct {
	OrderBucketList           map[int]*OrderBucket           //[category_id]OrderBucket
	UserBucketAmountRangeList map[int]map[string]*UserBucket //[category_id][amount_range]UserBucket
	UserTotal                 *UserTotal
	AmountRange               AmountRange
	Settings                  Settings
	//GrabEventInterrupt        chan int
}

type AmountRange struct {
	MinAmount int
	MaxAmount int
	Range     []AmountRangeElement
}

type AmountRangeElement struct {
	MinAmount int
	MaxAmount int
}

type Settings struct {
	GrabTimeout        int //抢单超时时间
	GrabDayTotalAmount int //每天可抢总额度
	GrabDayOrderCnt    int //每天可抢总订单数量
	GrabIntervalTime   int //抢单间隔
}

func NewGrabOrder() *GrabOrder {
	grabOrder := new(GrabOrder)

	grabOrder.InitAmountRange()
	grabOrder.InitSettings()
	grabOrder.InitUserBucketOrderBucket()

	grabOrder.UserTotal = NewUserTotal()
	//grabOrder.GrabEventInterrupt = make(chan int)

	return grabOrder
}

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

func (grabOrder *GrabOrder) InitSettings() {
	settings := Settings{
		GrabTimeout:        30,
		GrabDayTotalAmount: 10000,
		GrabDayOrderCnt:    50,
		GrabIntervalTime:   60,
	}
	grabOrder.Settings = settings
}

func (grabOrder *GrabOrder) InitUserBucketOrderBucket() {
	grabOrder.OrderBucketList = make(map[int]*OrderBucket)
	for categoryId, _ := range PayCategoryList() { //每个支付分类 - 有一个 订单桶
		grabOrder.OrderBucketList[categoryId] = NewOrderBucket(categoryId)
		for _, v := range grabOrder.AmountRange.Range {
			key := GetRangeKey(v.MinAmount, v.MaxAmount)
			grabOrder.UserBucketAmountRangeList[categoryId][key] = NewUserBucket(categoryId, v.MinAmount, v.MaxAmount)
		}
	}
}

func (grabOrder *GrabOrder) CreateOrder(order Order) {
	order.Timeout = grabOrder.Settings.GrabTimeout
	key := ""
	for _, v := range grabOrder.AmountRange.Range {
		if order.Amount >= v.MinAmount && order.Amount <= v.MaxAmount {
			key = GetRangeKey(v.MinAmount, v.MaxAmount)
		}
	}
	if key == "" {
		fmt.Println("err1")
	}

	grabOrder.OrderBucketList[order.CategoryId].AddOne(order)
	timer := time.NewTimer(time.Second * time.Duration(grabOrder.Settings.GrabTimeout))

	userBucket := grabOrder.UserBucketAmountRangeList[order.CategoryId][key]
	times := 1
	for {
		if times >= 3 {
			break
		}

		userInfo, err := userBucket.PopOne()
		if err != nil {
			fmt.Println("err2")
		}

		//rs := grabOrder.GrabStart(userInfo.Uid, order.Id)
		grabDoingExecRsChan := make(chan int)
		grabDoingExecRs := 0
		doneFlag := 0
		timeoutFlag := 0
		go grabOrder.GrabDoing(order.Uid, order.Id, grabDoingExecRsChan)

		select {
		case <-timer.C:
			fmt.Println("timeout")
			timeoutFlag = 1
		case grabDoingExecRs = <-grabDoingExecRsChan:
			fmt.Println("抢单返回结果")
			doneFlag = 1
		}

		if timeoutFlag == 1 {

		}

		if doneFlag == 1 {
			if grabDoingExecRs == 1 { //抢单成功
				userInfo.Weight++
				userInfo.SuccessTime++
			} else {
				userInfo.Weight--
				userInfo.FailedTime++
			}
		}

		userBucket.PushOne(userInfo)
		times++
	}
}

type EventMsg struct {
	OrderId int
	TypeId  int
	Content string
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

type UserGrabInfo struct {
	PayCategoryId int
	AmountMin     int
	AmountMax     int
}

func GetRangeKey(amountMin int, amountMax int) string {
	return strconv.Itoa(amountMin) + "_" + strconv.Itoa(amountMax)
}

func (grabOrder *GrabOrder) UserOpenGrab(uid int, userGrabInfo []UserGrabInfo) bool {
	if uid <= 0 || len(userGrabInfo) == 0 {
		fmt.Println("err4")
	}

	//200 - 2000
	//100 - 500 | 501 - 1000  | 1000 - 5000

	for _, userChannel := range userGrabInfo {
		if userChannel.AmountMin < grabOrder.AmountRange.MinAmount {
			fmt.Println("userChannel.AmountMin < grabOrder.AmountRange.MinAmount")
		}

		if userChannel.AmountMax > grabOrder.AmountRange.MaxAmount {
			fmt.Println("userChannel.AmountMax > grabOrder.AmountRange.MaxAmount")
		}
		for _, v := range grabOrder.AmountRange.Range { //100 - 500 | 501 - 1000  | 1000 - 5000
			if userChannel.AmountMin >= v.MinAmount && userChannel.AmountMax <= v.MaxAmount {
				userBucketAmountRange := grabOrder.UserBucketAmountRangeList[userChannel.PayCategoryId]
				key := GetRangeKey(userChannel.AmountMin, userChannel.AmountMax)
				userBucketAmountRange[key].PushNewOne(uid)
				break
			}
		}
	}

	return true
}

//func (grabOrder *GrabOrder) GrabStart(uid int, oid int) bool {
//
//}

func (grabOrder *GrabOrder) GrabDoing(uid int, oid int, grabCh chan int) {
	order := grabOrder.OrderBucketList[oid].GelOne(oid)
	if order.Timeout > int(time.Now().Unix()) {

	}

	userTotal := grabOrder.UserTotal.GetOne(uid)
	if userTotal.GrabDayTotalAmount > grabOrder.Settings.GrabDayTotalAmount {

	}

	if userTotal.GrabDayTotalOrderCnt > grabOrder.Settings.GrabDayOrderCnt {

	}

	if int(time.Now().Unix()-userTotal.LastGrabSuccessTime) < grabOrder.Settings.GrabIntervalTime {

	}

	if userTotal.WsStatus == 0 {

	}

	if userTotal.GrabStatus == 0 {

	}

	if userTotal.PayCategoryStatus == 0 {

	}

	if userTotal.PayChannelStatus == 0 {

	}

	//可抢额度=账户余额-押金-冻结金额

	//order record 落盘
	//更新用户：余额、冻结金额
	//合建账变记录

	grabOrder.UserTotal.UpdateGrabDayTotalOrderCnt()
	grabOrder.UserTotal.UpdateGrabDayTotalAmountProgress()
	grabOrder.UserTotal.UpdateGrabDayTotalAmount()
	grabOrder.UserTotal.UpdateLastGrabSuccessTime()

	grabCh <- 1
}
