package grab_order

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"zgoframe/util"
)

type UserElement struct {
	Uid                 int
	WsStatus            int   //web-socket 在线状态
	GrabStatus          int   //抢单开关
	LastGrabSuccessTime int64 //最后抢单成功的时间，用于计算间隔
	CreateTime          int64
	UpdateTime          int64
	UserDayTotal        UserDayTotal
	//GrabTotalCnt        int //已抢订单总数量
	//GrabTotalAmount     int //已抢订单总金额
	//PayCategoryStatus   int   //支付分类的状态
	//PayChannelStatus    int   //用户支付通道的状态
}

type UserDayTotal struct {
	Date                string `json:"date"`                   //天
	GrabAmountProgress  int    `json:"grab_amount_progress"`   //进行中的金额，防止多个池子同时抢单
	GrabCnt             int    `json:"grab_cnt"`               //今天已抢订单数量
	GrabAmount          int    `json:"grab_amount"`            //今天已抢总金额数
	SuccessTime         int    `json:"success_time"`           //今天总成功次数
	FailedTime          int    `json:"failed_time"`            //今天总失败次数
	LastGrabSuccessTime int64  `json:"last_grab_success_time"` //最后成功抢单时间
}

type UserTotal struct {
	UserElementList map[int]UserElement
	Redis           *util.MyRedis
}

func NewUserTotal(redis *util.MyRedis) *UserTotal {
	userTotal := new(UserTotal)
	userTotal.Redis = redis
	return userTotal
}

func (userTotal *UserTotal) AddOne(uid int) error {
	userElement, exist := userTotal.GetOne(uid)

	ymd := time.Now().Format("2006") + time.Now().Format("01") + time.Now().Format("02")
	key := "grab_order_day_total_" + ymd + strconv.Itoa(uid)

	if exist {
		userElement.UpdateTime = time.Now().Unix()
		if userElement.UserDayTotal.Date == "" {
			fmt.Println("err userElement.UserDayTotal.Date empty")
		}

		if userElement.UserDayTotal.Date != ymd {
			fmt.Println("err UserDayTotal.Date != today")
			userElement.UserDayTotal = userTotal.GetStructUserDayTotal(key)
		}

		return nil
	}

	userElement = UserElement{
		Uid:        uid,
		WsStatus:   0,
		GrabStatus: 0,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	keyExist := userTotal.Redis.Redis.Exists(context.TODO(), key).Val()
	if keyExist <= 0 {
		userDayTotal := UserDayTotal{
			Date:                ymd,
			GrabAmountProgress:  0,
			GrabCnt:             0,
			GrabAmount:          0,
			SuccessTime:         0,
			FailedTime:          0,
			LastGrabSuccessTime: 0,
		}
		bytes, _ := json.Marshal(userDayTotal)
		userTotal.Redis.Redis.Set(context.Background(), key, string(bytes), 0)
		userElement.UserDayTotal = userDayTotal
	} else {
		userElement.UserDayTotal = userTotal.GetStructUserDayTotal(key)
	}

	userTotal.UserElementList[uid] = userElement

	return nil
}
func (userTotal *UserTotal) GetStructUserDayTotal(key string) (userDayTotal UserDayTotal) {
	res := userTotal.Redis.Redis.Get(context.Background(), key)
	if res.Err() != nil {
		fmt.Println("GetStructUserDayTotal: " + res.Err().Error())
		return
	}
	err := json.Unmarshal([]byte(res.Val()), &userDayTotal)
	if err != nil {
		fmt.Println("GetStructUserDayTotal :", err)
		return
	}
	//fmt.Println(userElement.UserDayTotal)
	return userDayTotal
}

func (userTotal *UserTotal) GetOne(uid int) (e UserElement, exist bool) {
	element, ok := userTotal.UserElementList[uid]
	if !ok {
		return element, false
	}
	return element, true
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
