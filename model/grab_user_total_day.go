package model

// @description 用户-抢单每日汇总表
type GrabUserTotalDay struct {
	MODEL
	Uid                 int    `json:"uid" db:"define:int;comment:uid;defaultValue:0"`
	Date                string `json:"date" db:"define:varchar(255);comment:日期;defaultValue:''"`
	OrderCnt            int    `json:"order_cnt" db:"define:int;comment:成功抢单-总次数;defaultValue:0"`
	AmountTotal         int    `json:"amount_total" db:"define:int;comment:成功失意-总金额;defaultValue:0"`
	SuccessTime         int    `json:"success_time" db:"define:int;comment:抢单成功次数;defaultValue:0"`
	FailedTime          int    `json:"failed_time" db:"define:int;comment:抢单失败次数;defaultValue:0"`
	LastGrabSuccessTime int    `json:"last_grab_success_time" db:"define:int;comment:最后成功抢单时间;defaultValue:0"`
}

func (grabUserTotalDay *GrabUserTotalDay) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户-抢单每日汇总表"

	return m
}
