package grab_order

type Order struct {
	Id         int `json:"id"`
	Amount     int `json:"amount"`
	Uid        int `json:"uid"`
	CategoryId int `json:"category_id"`
	Timeout    int `json:"timeout"`
	StartTime  int `json:"start_time"`
}

// 支付-类型桶
type OrderBucket struct {
	CategoryId int           `json:"category_id"`
	List       map[int]Order `json:"list"`
}

func NewOrderBucket(categoryId int) *OrderBucket {
	payCategoryBucket := new(OrderBucket)
	payCategoryBucket.CategoryId = categoryId

	return payCategoryBucket
}

func (ob OrderBucket) AddOne(order Order) error {
	ob.List[order.Id] = order
	return nil
}

func (ob OrderBucket) GelOne(oid int) Order {
	return ob.List[oid]
}

func (ob OrderBucket) DelOne(oid int) error {
	delete(ob.List, oid)
	return nil
}
