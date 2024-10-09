package grab_order

import (
	"gorm.io/gorm"
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
	CategoryId int                            `json:"category_id"`
	List       map[string]model.PayOrderMatch `json:"list"`
	Redis      *util.MyRedis                  `json:"-"`
	Gorm       *gorm.DB                       `json:"-"`
}

func NewOrderBucket(categoryId int, redis *util.MyRedis, gorm *gorm.DB) *OrderBucket {
	orderBucket := new(OrderBucket)
	orderBucket.CategoryId = categoryId
	orderBucket.List = make(map[string]model.PayOrderMatch)
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
		ob.List[dbOrder.InId] = dbOrder
	}

	return
}
func (ob OrderBucket) AddOne(order model.PayOrderMatch) error {
	ob.List[order.InId] = order
	ob.Gorm.Create(&order)

	return nil
}

// 更新一整条记录
func (ob OrderBucket) UpOneRecord(order model.PayOrderMatch) {
	ob.List[order.InId] = order
	ob.Gorm.Save(&order)
}

func (ob OrderBucket) GelOne(oid string) model.PayOrderMatch {
	return ob.List[oid]
}

func (ob OrderBucket) DelOne(oid string) error {
	delete(ob.List, oid)
	return nil
}
