package model

// @description 支付分类
type PayCategory struct {
	MODEL
	Name   string `json:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`              //标题
	Sn     string `json:"sn" db:"define:varchar(255);comment:sn码;defaultValue:''"`              //模板内容,可变量替换
	Status int    `json:"status" db:"define:tinyint(1);comment:0未知1开启2关闭;defaultValue:0"`       //分类,1验证码2通知3营销4报警
	Sort   int    `json:"sort" db:"define:int;comment:排序权重;defaultValue:0"`                     //每天最多发送次数
	Remark string `json:"remark" db:"define:varchar(255);comment:描述，主要是给3方审核用;defaultValue:''"` //备注
	Icon   string `json:"icon" db:"define:varchar(255);comment:图标URL;defaultValue:0" `          //周期内最多可发送次数
	MinAmt int    `json:"min_amt" db:"define:int;comment:最小金额限制;defaultValue:0" `               //验证码的失效时间
	MaxAmt int    `json:"max_amt" db:"define:int;comment:最大金额限制;defaultValue:0" `               //验证码的失效时间
}

func (payCategory *PayCategory) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "支付分类"

	return m
}
