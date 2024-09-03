package grab_order

type UserElement struct {
	Uid                        int
	WsStatus                   int
	GrabStatus                 int
	PayCategoryStatus          int
	PayChannelStatus           int
	GrabDayTotalAmountProgress int //进行中的金额，防止多个池子同时抢单
	//GrabTotalOrderCnt          int   //已抢订单数量
	GrabDayTotalOrderCnt int //今天已抢订单数量
	//GrabTotalAmount            int   //总抢金额
	GrabDayTotalAmount  int   //今天已抢总金额数
	LastGrabSuccessTime int64 //最后抢单成功的时间，用于计算间隔
	SuccessTime         int   //总成功次数
	FailedTime          int   //总失败次数
	CreateTime          int64
	UpdateTime          int64
}

type UserTotal struct {
	UserElementList map[int]UserElement
}

func NewUserTotal() *UserTotal {
	userTotal := new(UserTotal)
	return userTotal
}

func (userTotal *UserTotal) AddOne() {

}

func (userTotal *UserTotal) GetOne(uid int) UserElement {
	return userTotal.UserElementList[uid]
}

func (userTotal *UserTotal) GetUserOrderCntDay() {

}

func (userTotal *UserTotal) GetUserOrderAmountDay() {

}

func (userTotal *UserTotal) UpdateGrabDayTotalAmountProgress() {

}

func (userTotal *UserTotal) UpdateGrabDayTotalOrderCnt() {

}

func (userTotal *UserTotal) UpdateGrabDayTotalAmount() {

}

func (userTotal *UserTotal) UpdateLastGrabSuccessTime() {

}
