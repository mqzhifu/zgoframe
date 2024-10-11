package grab_order

import (
	"gorm.io/gorm"
	"time"
	"zgoframe/model"
	"zgoframe/util"
)

//type Order struct {
//	Amount            int `json:"amount"`               //金额
//	Uid               int `json:"uid"`                  //推单用户id
//	CategoryId        int `json:"category_id"`          //分类ID
//	Timeout           int `json:"timeout"`              //超时
//	StartTime         int `json:"start_time"`           //开始时间
//	SuccessTime       int `json:"success_time"`         //成功时间
//	MatchTimes        int `json:"match_times"` //匹配了多少次
//	GrabUid           int `json:"grab_uid"`             //抢到单用户
//	MatchQueueUserCnt int `json:"match_queue_user_cnt"` //开始匹配时，池里有多少个用户
//}

// 支付-类型桶
type OrderBucket struct {
	CategoryId int              `json:"category_id"`
	List       map[string]Order `json:"list"`
	Redis      *util.MyRedis    `json:"-"`
	Gorm       *gorm.DB         `json:"-"`
}

type Order struct {
	model.PayOrderMatch
	//CloseOrderTimeoutDemon chan int
}

func NewOrderBucket(categoryId int, redis *util.MyRedis, gorm *gorm.DB) *OrderBucket {
	orderBucket := new(OrderBucket)
	orderBucket.CategoryId = categoryId
	orderBucket.List = make(map[string]Order)
	orderBucket.Redis = redis
	orderBucket.Gorm = gorm

	orderBucket.LoadDataFromRedis(categoryId)

	return orderBucket
}
func (ob OrderBucket) LoadDataFromRedis(categoryId int) {
	list := []model.PayOrderMatch{}
	ob.Gorm.Where("category_id = ?", categoryId).Find(&list)

	if len(list) <= 0 {
		return
	}

	for _, dbOrder := range list {
		oo := Order{}
		ob.List[dbOrder.InId] = oo
	}

	return
}
func (ob OrderBucket) AddOne(order Order) error {
	//order.CloseOrderTimeoutDemon = make(chan int)
	ob.List[order.InId] = order
	ob.Gorm.Create(&order)
	go ob.CheckOrderTimeout(order)
	return nil
}

// 更新一整条记录
func (ob OrderBucket) UpOneRecord(order Order) {
	ob.List[order.InId] = order
	ob.Gorm.Save(&order)
}

func (ob OrderBucket) GelOne(oid string) Order {
	return ob.List[oid]
}

func (ob OrderBucket) DelOne(oid string) error {
	delete(ob.List, oid)
	return nil
}

// 监听：订单超时
func (ob OrderBucket) CheckOrderTimeout(order Order) {
	for {
		if order.Status != ORDER_MATCH_STATUS_ING {
			break
		}
		//这里多一秒，是前面的方法已经做了超时算是，这里做个保底操作
		if int(time.Now().Unix()) > order.Timeout+1 {

		}
		time.Sleep(time.Second * 1)
	}
}
