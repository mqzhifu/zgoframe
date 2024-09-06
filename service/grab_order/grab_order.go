package grab_order

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
	"zgoframe/model"
	"zgoframe/util"
)

type GrabOrder struct {
	OrderBucketList           map[int]*OrderBucket           `json:"order_bucket_list"`             //[category_id]OrderBucket
	UserBucketAmountRangeList map[int]map[string]*UserBucket `json:"user_bucket_amount_range_list"` //[category_id][amount_range]UserBucket
	UserTotal                 *UserTotal                     `json:"user_total"`
	AmountRange               AmountRange                    `json:"amount_range"`
	Settings                  Settings                       `json:"settings"`
	Gorm                      *gorm.DB                       `json:"-"`
	Redis                     *util.MyRedis                  `json:"-"`
	//GrabEventInterrupt        chan int
}

type AmountRange struct {
	MinAmount int                  `json:"min_amount"`
	MaxAmount int                  `json:"max_amount"`
	Range     []AmountRangeElement `json:"range"`
}

type AmountRangeElement struct {
	MinAmount int `json:"min_amount"`
	MaxAmount int `json:"max_amount"`
}

type Settings struct {
	GrabTimeout        int `json:"grab_timeout"`          //抢单超时时间
	GrabDayTotalAmount int `json:"grab_day_total_amount"` //每天可抢总额度
	GrabDayOrderCnt    int `json:"grab_day_order_cnt"`    //每天可抢总订单数量
	GrabIntervalTime   int `json:"grab_interval_time"`    //抢单间隔
}

func NewGrabOrder(db *gorm.DB, redis *util.MyRedis) *GrabOrder {
	grabOrder := new(GrabOrder)
	grabOrder.Redis = redis
	grabOrder.Gorm = db

	grabOrder.InitAmountRange()
	grabOrder.InitSettings()
	grabOrder.InitUserBucketOrderBucket()

	grabOrder.UserTotal = NewUserTotal(redis)

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
func (grabOrder *GrabOrder) GetPayCategory() ([]model.PayCategory, error) {
	list := []model.PayCategory{}
	grabOrder.Gorm.Find(&list)
	return list, nil
}
func (grabOrder *GrabOrder) GetData() (*GrabOrder, error) {
	return grabOrder, nil
}

func (grabOrder *GrabOrder) CreateOrder(order Order) error {
	order.Timeout = grabOrder.Settings.GrabTimeout
	key := ""
	for _, v := range grabOrder.AmountRange.Range {
		if order.Amount >= v.MinAmount && order.Amount <= v.MaxAmount {
			key = GetRangeKey(v.MinAmount, v.MaxAmount)
		}
	}
	if key == "" {
		fmt.Println("err1 key empty")
		return errors.New("key empty")
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
		//doneFlag := 0
		timeoutFlag := 0
		go grabOrder.GrabDoing(order.Uid, order.Id, grabDoingExecRsChan)

		select {
		case <-timer.C:
			fmt.Println("timeout")
			timeoutFlag = 1
		case grabDoingExecRs = <-grabDoingExecRsChan:
			fmt.Println("抢单返回结果")
			//doneFlag = 1
		}

		if timeoutFlag == 1 {

		}
		fmt.Println(grabDoingExecRs)

		//if doneFlag == 1 {
		//	if grabDoingExecRs == 1 { //抢单成功
		//		userInfo.Weight++
		//		userInfo.SuccessTime++
		//	} else {
		//		userInfo.Weight--
		//		userInfo.FailedTime++
		//	}
		//}

		userBucket.PushOne(userInfo)
		times++
	}
	return nil
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

func (grabOrder *GrabOrder) UserOpenGrab(uid int, userGrabInfo []UserGrabInfo) error {
	if uid <= 0 || len(userGrabInfo) == 0 {
		fmt.Println("err4")
		return errors.New("err4")
	}
	//200 - 2000
	//100 - 500 | 501 - 1000  | 1000 - 5000
	fmt.Println("3333========", grabOrder.AmountRange.MinAmount, grabOrder.AmountRange.MaxAmount)
	userAmountSection := []string{}
	for _, userChannel := range userGrabInfo {
		//fmt.Println(userChannel.AmountMin, userChannel.AmountMax)
		if userChannel.AmountMin < grabOrder.AmountRange.MinAmount {
			fmt.Println("此用户支付渠道不满足 userChannel.AmountMin < grabOrder.AmountRange.MinAmount")
			continue
		}

		if userChannel.AmountMax > grabOrder.AmountRange.MaxAmount {
			fmt.Println("此用户支付渠道不满足 userChannel.AmountMax > grabOrder.AmountRange.MaxAmount")
			continue
		}

		for _, v := range grabOrder.AmountRange.Range { //100 - 500 | 501 - 1000  | 1000 - 5000
			if userChannel.AmountMin >= v.MaxAmount {
				continue
			}
			userAmountSection = append(userAmountSection, GetRangeKey(v.MinAmount, v.MaxAmount))
			if userChannel.AmountMax <= v.MaxAmount {
				break
			}
		}

		//for _, v := range userAmountSection {
		//grabOrder.UserTotal.AddOne(uid)
		//userBucket := grabOrder.UserBucketAmountRangeList[userChannel.PayCategoryId][v]
		//userBucket.PushNewOne(uid)
		//}
	}

	return nil
}

func (grabOrder *GrabOrder) GrabDoing(uid int, oid int, grabCh chan int) {
	order := grabOrder.OrderBucketList[oid].GelOne(oid)
	if order.Timeout > int(time.Now().Unix()) {

	}

	//userTotal, _ := grabOrder.UserTotal.GetOne(uid)
	//if userTotal.GrabDayTotalAmount > grabOrder.Settings.GrabDayTotalAmount {
	//
	//}
	//
	//if userTotal.GrabDayTotalOrderCnt > grabOrder.Settings.GrabDayOrderCnt {
	//
	//}
	//
	//if int(time.Now().Unix()-userTotal.LastGrabSuccessTime) < grabOrder.Settings.GrabIntervalTime {
	//
	//}
	//
	//if userTotal.WsStatus == 0 {
	//
	//}
	//
	//if userTotal.GrabStatus == 0 {
	//
	//}
	//
	//if userTotal.PayCategoryStatus == 0 {
	//
	//}
	//
	//if userTotal.PayChannelStatus == 0 {
	//
	//}
	//
	////可抢额度=账户余额-押金-冻结金额
	//
	////order record 落盘
	////更新用户：余额、冻结金额
	////合建账变记录
	//
	//grabOrder.UserTotal.UpdateGrabDayTotalOrderCnt()
	//grabOrder.UserTotal.UpdateGrabDayTotalAmountProgress()
	//grabOrder.UserTotal.UpdateGrabDayTotalAmount()
	//grabOrder.UserTotal.UpdateLastGrabSuccessTime()
	//
	//grabCh <- 1
}
