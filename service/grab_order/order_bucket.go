package grab_order

type Order struct {
	Id         int
	Amount     int
	Uid        int
	CategoryId int
	Timeout    int
}

// 支付-类型桶
type OrderBucket struct {
	CategoryId     int
	QueueOrderInfo map[int]Order
}

func NewOrderBucket(categoryId int) *OrderBucket {
	payCategoryBucket := new(OrderBucket)
	payCategoryBucket.CategoryId = categoryId

	return payCategoryBucket
}

func (ob OrderBucket) AddOne(order Order) error {
	ob.QueueOrderInfo[order.Id] = order
	return nil
}

func (ob OrderBucket) GelOne(oid int) Order {
	return ob.QueueOrderInfo[oid]
}

func (ob OrderBucket) DelOne(oid int) error {
	delete(ob.QueueOrderInfo, oid)
	return nil
}
