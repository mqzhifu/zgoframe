package grab_order

type Order struct {
	Id                string `json:"id"`                   //订单ID
	Amount            int    `json:"amount"`               //金额
	Uid               int    `json:"uid"`                  //推单用户id
	CategoryId        int    `json:"category_id"`          //分类ID
	Timeout           int    `json:"timeout"`              //超时
	StartTime         int    `json:"start_time"`           //开始时间
	SuccessTime       int    `json:"success_time"`         //成功时间
	GrabUid           int    `json:"grab_uid"`             //抢到单用户
	MatchQueueUserCnt int    `json:"match_queue_user_cnt"` //开始匹配时，池里有多少个用户
}

// 支付-类型桶
type OrderBucket struct {
	CategoryId int              `json:"category_id"`
	List       map[string]Order `json:"list"`
}

func NewOrderBucket(categoryId int) *OrderBucket {
	orderBucket := new(OrderBucket)
	orderBucket.CategoryId = categoryId
	orderBucket.List = make(map[string]Order)

	return orderBucket
}

func (ob OrderBucket) AddOne(order Order) error {
	ob.List[order.Id] = order
	return nil
}

func (ob OrderBucket) GelOne(oid string) Order {
	return ob.List[oid]
}

func (ob OrderBucket) DelOne(oid string) error {
	delete(ob.List, oid)
	return nil
}
