package model

// @description 用户-自动抢单开启
type GrabUserAutoOpen struct {
	MODEL
	Uid        int `json:"uid" db:"define:int;comment:uid;defaultValue:0"`
	CategoryId int `json:"category_id" db:"define:int;comment:支付类型;defaultValue:0"` //分类，区分业务
	AmountMin  int `json:"amount_min"  db:"define:int;comment:最小值;defaultValue:0"`
	AmountMax  int `json:"amount_max"  db:"define:int;comment:最大值;defaultValue:0"`
	BatchId    int `json:"batch_id"  db:"define:int;comment:批次号;defaultValue:0"`
}

func (grabUserAutoOpen *GrabUserAutoOpen) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户-自动抢单开启"

	return m
}
