package model

// @description 用户各种数据汇总
type UserTotal struct {
	MODEL
	Uid        int `json:"uid" db:"define:int;comment:uid;defaultValue:0"`
	GoldCoin   int `json:"gold_coin" db:"define:int;comment:金币;defaultValue:0"`
	Cash       int `json:"cash" db:"define:int;comment:现金;defaultValue:0"`
	FreezeCash int `json:"freeze_cash" db:"define:int;comment:冻结现金;defaultValue:0"`
	Deposit    int `json:"deposit" db:"define:int;comment:押金;defaultValue:0"`
	Experience int `json:"experience" db:"define:int;comment:经验;defaultValue:0"`
}

func (userTotal *UserTotal) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户需要统计汇总的信息"

	return m
}
