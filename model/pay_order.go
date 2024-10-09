package model

// @description 支付订单
type PayOrder struct {
	MODEL
	OutId       string `json:"pay_out_id" db:"define:varchar(100);comment:3方返回的订单号;defaultValue:''"`                        //对外ID，3方返回的订单号
	InId        string `json:"pay_in_id" db:"define:varchar(100);comment:给3方支付的订单号;defaultValue:''"`                        //对内ID，给3方支付的订单号
	ProjectId   int    `json:"project_id" db:"define:int;comment:项目ID;defaultValue:0"`                                      //分类，区分业务
	Status      int    `json:"status" db:"define:tinyint(1);comment:状态1:正常 2:下架;defaultValue:0"`                            //状态 1:未支付 2:已支付 3超时
	Price       int    `json:"price" db:"define:int;comment:价格(单位分);defaultValue:0"`                                        //价格(单位分)
	Uid         int    `json:"uid" db:"define:int;comment:用户ID;defaultValue:0" `                                            //用户ID
	PayType     int    `json:"pay_type" db:"define:tinyint(1);comment:支付类型1微信2支付宝;defaultValue:0"`                          //支付类型1微信2支付宝
	PaySubType  int    `json:"pay_sub_type" db:"define:tinyint(1);comment:支付子类型1.APP 2.浏览器PC 3.浏览器H5 4.二维码;defaultValue:0"` //支付子类型，未使用
	PayBackTime int64  `json:"pay_back_time" db:"define:bigint(20);comment:支付回调时间;defaultValue:0"`                          //支付回调时间
	PayBackInfo string `json:"pay_back_info" db:"define:varchar(255);comment:支付回调数据;defaultValue:''"`                       //支付回调数据
	Timeout     int    `json:"timeout" db:"define:int;comment:超时时间;defaultValue:0"`                                         //超时时间
	UserMemo    string `json:"user_memo" db:"define:varchar(255);comment:用户备注;defaultValue:''"`                             //备注
	AdminMemo   string `json:"admin_memo" db:"define:varchar(255);comment:管理员备注;defaultValue:''"`                           //备注
}

func (payOrder *PayOrder) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "支付订单"

	return m
}
