package model

// @description 支付订单
type PayOrderMatch struct {
	MODEL
	InId              string `json:"in_id" db:"define:varchar(100);comment:对内ID自己生成;defaultValue:''"`         //对内ID，自己生成
	Status            int    `json:"status" db:"define:tinyint(1);comment:状态1匹配中2成功3失败;defaultValue:0"`       //状态1:匹配中 2:成功 3失败
	Amount            int    `json:"amount" db:"define:int;comment:价格(单位分);defaultValue:0"`                   //价格(单位分)
	Uid               int    `json:"uid" db:"define:int;comment:用户ID;defaultValue:0" `                        //用户ID
	CategoryId        int    `json:"category_id" db:"define:tinyint(1);comment:支付类型1微信2支付宝;defaultValue:0"`   //支付类型1微信2支付宝
	GrabUid           int    `json:"grab_uid" db:"define:int;comment:抢单成功用户;defaultValue:0"`                  //支付子类型，未使用
	Timeout           int    `json:"timeout" db:"define:int;comment:超时时间;defaultValue:0"`                     //超时时间
	StartTime         int    `json:"start_time" db:"define:int;comment:开始匹配时间;defaultValue:0"`                //超时时间
	EndTime           int    `json:"end_time" db:"define:int;comment:匹配结束时间;defaultValue:0"`                  //超时时间
	MatchTimes        int    `json:"match_times" db:"define:int;comment:匹配次数;defaultValue:0"`                 //超时时间
	MatchQueueUserCnt int    `json:"match_queue_user_cnt" db:"define:int;comment:开始时，池里有多少用户;defaultValue:0"` //超时时间
}

func (payOrderMatch *PayOrderMatch) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "抢单匹配"

	return m
}
