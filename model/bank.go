package model

// @description 银行
type Bank struct {
	MODEL
	Name    string `json:"name" db:"define:varchar(50);comment:名称;defaultValue:''"`               //标题
	Code    string `json:"code" db:"define:varchar(255);comment:sn码;defaultValue:''"`             //模板内容,可变量替换
	Status  int    `json:"status" db:"define:tinyint(1);comment:0未知1开启2关闭;defaultValue:0"`        //分类,1验证码2通知3营销4报警
	Address string `json:"address" db:"define:varchar(255);comment:描述，主要是给3方审核用;defaultValue:''"` //备注
}

func (bank *Bank) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "银行"

	return m
}
