package model

// @description 电子邮件 - 发送规则
type Goods struct {
	MODEL
	Title         string `json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`                     //标题
	Desc          string `json:"desc" db:"define:varchar(255);comment:描述;defaultValue:''"`                     //描述
	Type          int    `json:"type" db:"define:tinyint(1);comment:分类;defaultValue:0"`                        //分类
	Price         int    `json:"price" db:"define:int;comment:价格(单位分);defaultValue:0"`                         //价格
	Amount        int    `json:"period" db:"define:int;comment:买一个商品给多少个数量;defaultValue:0"`                    //买一个商品给多少个数量
	Stock         int    `json:"stock" db:"define:int;comment:库存数量;defaultValue:0"`                            //库存数量
	AdminId       int    `json:"admin_id" db:"define:int;comment:管理员ID;defaultValue:0" `                       //管理员ID
	AllowCoupon   int    `json:"allow_coupon" db:"define:int;comment:允许使用优惠券;defaultValue:0"`                  //允许使用优惠券，0不允许
	AllowGoldCoin int    `json:"allow_gold_coin" db:"define:int;comment:允许支付类型1微信2支付宝;defaultValue:0"`         //允许使用金币，0不允许
	PayAllowType  string `json:"pay_allow_type" db:"define:varchar(255);comment:允许支付类型1微信2支付宝;defaultValue:0"` //允许支付类型1微信2支付宝
	Status        int    `json:"status" db:"define:tinyint(1);comment:状态1:正常 2:下架;defaultValue:0"`             //状态 1:正常 2:下架
	Memo          string `json:"memo" db:"define:varchar(255);comment:备注;defaultValue:''"`                     //备注
}

func (goods *Goods) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "商品"

	return m
}
