package model

// @description 电子邮件 - 发送规则
type Orders struct {
	MODEL
	Status      int    `json:"status" db:"define:tinyint(1);comment:状态1:正常 2:下架;defaultValue:0"`                            //状态 1:未支付 2:已支付 3超时
	Type        int    `json:"type" db:"define:tinyint(1);comment:分类;defaultValue:0"`                                       //分类
	GoodsId     int    `json:"goods_id" db:"define:int;comment:商品ID;defaultValue:0"`                                        //商品ID
	Price       int    `json:"price" db:"define:int;comment:价格(单位分);defaultValue:0"`                                        //价格(单位分)
	Amount      int    `json:"amount" db:"define:int;comment:购买数量;defaultValue:0"`                                          //购买数量
	RealPrice   int    `json:"real_price" db:"define:int;comment:实付价格(单位分);defaultValue:0" `                                //实付价格(单位分)
	CouponId    int    `json:"coupon_id" db:"define:int;comment:优惠卷ID;defaultValue:0" `                                     //优惠卷ID
	Uid         int    `json:"uid" db:"define:int;comment:用户ID;defaultValue:0" `                                            //用户ID
	GoldCoin    int    `json:"gold_coin" db:"define:int;comment:金币数;defaultValue:0" `                                       //金币数
	Source      int    `json:"source" db:"define:tinyint(1);comment:来源;defaultValue:0"`                                     //来源
	PayOutId    string `json:"pay_out_id" db:"define:varchar(100);comment:3方返回的订单号;defaultValue:''"`                        //3方返回的订单号
	UserMemo    string `json:"user_memo" db:"define:varchar(255);comment:用户备注;defaultValue:''"`                             //备注
	AdminMemo   string `json:"admin_memo" db:"define:varchar(255);comment:管理员备注;defaultValue:''"`                           //备注
	PayType     int    `json:"pay_type" db:"define:tinyint(1);comment:支付类型1微信2支付宝;defaultValue:0"`                          //支付类型1微信2支付宝
	PaySubType  int    `json:"pay_sub_type" db:"define:tinyint(1);comment:支付子类型1.APP 2.浏览器PC 3.浏览器H5 4.二维码;defaultValue:0"` //
	PayBackTime int64  `json:"pay_back_time" db:"define:bigint(20);comment:支付回调时间;defaultValue:0"`                          //支付回调时间
	PayInId     string `json:"pay_in_id" db:"define:varchar(100);comment:给3方支付的订单号;defaultValue:''"`                        //给3方支付的订单号
}

func (orders *Orders) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "订单"

	return m
}
