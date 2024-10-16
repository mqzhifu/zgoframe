package grab_order

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

type UserElement struct {
	Uid                 int                         `json:"uid"`
	WsStatus            int                         `json:"ws_status"`              //web-socket 在线状态
	GrabStatus          int                         `json:"grab_status"`            //抢单开关
	LastGrabSuccessTime int64                       `json:"last_grab_success_time"` //最后抢单成功的时间，用于计算间隔
	CreateTime          int64                       `json:"create_time"`
	UpdateTime          int64                       `json:"update_time"`
	UserDayTotal        UserDayTotal                `json:"user_day_total"`
	UserOpenGrab        []request.GrabOrderUserOpen `json:"user_open_grab"` //用户开始自动抢单时的配置数据
	//GrabTotalCnt        int //已抢订单总数量
	//GrabTotalAmount     int //已抢订单总金额
	//PayCategoryStatus   int   //支付分类的状态
	//PayChannelStatus    int   //用户支付通道的状态
}

type UserDayTotal struct {
	Date                string `json:"date"`                   //天
	GrabCnt             int    `json:"grab_cnt"`               //今天已抢订单数量
	GrabAmount          int    `json:"grab_amount"`            //今天已抢总金额数
	SuccessTime         int    `json:"success_time"`           //今天总成功次数
	FailedTime          int    `json:"failed_time"`            //今天总失败次数
	LastGrabSuccessTime int64  `json:"last_grab_success_time"` //最后成功抢单时间
}

type UserTotal struct {
	UserElementList map[int]UserElement `json:"user_element_list"`
	Redis           *util.MyRedis       `json:"-"`
	Gorm            *gorm.DB            `json:"-"`
}

func NewUserTotal(redis *util.MyRedis, gorm *gorm.DB) *UserTotal {
	userTotal := new(UserTotal)
	userTotal.Redis = redis
	userTotal.Gorm = gorm
	userTotal.UserElementList = make(map[int]UserElement)
	return userTotal
}

func (userTotal *UserTotal) InsertGrabOrderUserOpenRecord(list []request.GrabOrderUserOpen) {
	rowMax := model.GrabUserAutoOpen{}
	userTotal.Gorm.Select("max(id) as id ").First(&rowMax)

	for _, v := range list {
		row := model.GrabUserAutoOpen{}
		row.Uid = v.Uid
		row.CategoryId = v.PayCategoryId
		row.AmountMax = v.AmountMax
		row.AmountMin = v.AmountMin
		row.BatchId = rowMax.BatchId + 1
		userTotal.Gorm.Create(&row)
	}
}

func (userTotal *UserTotal) AddOrUpdateOne(uid int, grabOrderUserOpen []request.GrabOrderUserOpen) (err error, optType int) {
	userElement, exist := userTotal.GetOne(uid)

	ymd := time.Now().Format("2006") + time.Now().Format("01") + time.Now().Format("02")
	key := "grab_order_day_total_" + ymd + "_" + strconv.Itoa(uid)

	if exist { //如果已经存在，做更新处理
		userElement.UpdateTime = time.Now().Unix()
		userElement.GrabStatus = USER_GRAB_STATUS_OPEN
		userElement.WsStatus = userTotal.GetUserWsStatus(uid)
		userElement.UserOpenGrab = grabOrderUserOpen
		if userElement.UserDayTotal.Date == "" {
			fmt.Println("err userElement.UserDayTotal.Date empty")
			//重新获取一下，当日的金额相关的统计数据
			userElement.UserDayTotal = userTotal.GetStructUserDayTotal(key)
		}

		if userElement.UserDayTotal.Date != ymd {
			fmt.Println("err UserDayTotal.Date != today")
			//重新获取一下，当日的金额相关的统计数据
			userElement.UserDayTotal = userTotal.GetStructUserDayTotal(key)
		}

		userTotal.InsertGrabOrderUserOpenRecord(grabOrderUserOpen)

		return nil, USER_TOTAL_OPT_TYPE_UP
	}
	//走到这里，证明，用户数据不存在进程中，需要重新创建一下
	userElement = UserElement{
		Uid:        uid,
		WsStatus:   userTotal.GetUserWsStatus(uid),
		GrabStatus: USER_GRAB_STATUS_OPEN,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	keyExist := userTotal.Redis.Redis.Exists(context.TODO(), key).Val()
	if keyExist <= 0 {
		userDayTotal := UserDayTotal{
			Date:                ymd,
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
	//添加到集合中
	userTotal.UserElementList[uid] = userElement
	userTotal.InsertGrabOrderUserOpenRecord(grabOrderUserOpen)

	return nil, USER_TOTAL_OPT_TYPE_ADD
}

func (userTotal *UserTotal) GetUserWsStatus(uid int) int {
	info, _ := userTotal.GetOne(uid)
	return info.WsStatus
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

// 获取 - 用户当天抢单总数
func (userTotal *UserTotal) GetUserOrderCntDay() {

}

// 获取 - 用户当天抢单总金额
func (userTotal *UserTotal) GetUserOrderAmountDay() {

}

func (userTotal *UserTotal) UpWs(uid int, status int) {
	info, exist := userTotal.GetOne(uid)
	if exist { // 获取 - 用户当天抢单-进行中的-订单数量
		info.WsStatus = status
	} //
} //}

// 更新 - 用户当天抢单总数
func (userTotal *UserTotal) UpGrabDayTotalOrderCnt(uid int) {
	info, exist := userTotal.GetOne(uid)
	if exist {
		info.UserDayTotal.GrabCnt++
	}
}

// 更新 - 用户当天抢单总金额
func (userTotal *UserTotal) UpGrabDayTotalAmount(uid int, amount int) {
	info, exist := userTotal.GetOne(uid)
	if exist {
		info.UserDayTotal.GrabAmount += amount
	}
}

// 更新 - 用户当天抢单-失败次数
func (userTotal *UserTotal) UpGrabFailedTime(uid int) {
	info, exist := userTotal.GetOne(uid)
	if exist {
		info.UserDayTotal.FailedTime++
	}
}

// 更新 - 用户当天抢单-进行中的-订单数量
func (userTotal *UserTotal) UpGrabSuccessTime(uid int) {
	info, exist := userTotal.GetOne(uid)
	if exist {
		info.UserDayTotal.SuccessTime++
	}
}

func (userTotal *UserTotal) UpLastGrabSuccessTime(uid int) {
	info, exist := userTotal.GetOne(uid)
	if exist {
		info.UserDayTotal.LastGrabSuccessTime = time.Now().Unix()
	}
}
