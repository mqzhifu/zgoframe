package model

// @description 用户-支付通道
type UserPayAccount struct {
	MODEL
	Uid           int    `json:"uid" db:"define:int;comment:uid;defaultValue:0"`
	CategoryId    int    `json:"category_id" db:"define:int;comment:支付类型;defaultValue:0"` //分类，区分业务
	AmountMin     int    `json:"amount_min"  db:"define:int;comment:最小值;defaultValue:0"`
	AmountMax     int    `json:"amount_max"  db:"define:int;comment:最大值;defaultValue:0"`
	BankId        int    `json:"bank_id"  db:"define:int;comment:银行ID;defaultValue:0"`
	Status        int    `json:"status"  db:"define:int;comment:状态;defaultValue:0"`
	AccountName   string `json:"account_name"  db:"define:int;comment:账户名;defaultValue:''"`   //账户名
	AccountNumber string `json:"account_number"  db:"define:int;comment:账户号;defaultValue:''"` //账户号
	ReceiveQRCode string `json:"receive_qr_code"  db:"define:int;comment:收款码;defaultValue:''"`
}

func (userPayAccount *UserPayAccount) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户-支付通道"

	return m
}
