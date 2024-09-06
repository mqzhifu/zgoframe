package request

// @description RBAC,角色权限控制，具体的权限(资源)
type CasbinInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// @description RBAC,角色权限控制，主要是给后台使用
type CasbinInReceive struct {
	AuthorityId string       `json:"authorityId"`
	CasbinInfos []CasbinInfo `json:"casbinInfos"`
}

// @description 上传文件
type UploadFile struct {
	File    string `json:"file" form:"file"`         //input file 控件的name
	Stream  string `json:"stream" form:"stream"`     //文件流,例：data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHgAAAB4CAMAAAAOus .......
	Module  string `json:"module" form:"module"`     //模块/业务名，可用于给文件名加前缀目录
	SyncOss int    `json:"sync_oss" form:"sync_oss"` //是否同步到云oss 1是2否
	HashDir int    `json:"hash_dir" form:"hash_dir"` //自动创建前缀目录 1不使用2月3天4小时
}

type FileDelete struct {
	SyncOss      int    `json:"sync_oss" form:"sync_oss"`           //是否同步到云oss 1是2否
	RelativePath string `json:"relative_path" form:"relative_path"` //相对路径 + 文件名
}

type FileCopy struct {
	SyncOss         int    `json:"sync_oss" form:"sync_oss"`                   //是否同步到云oss 1是2否
	SrcRelativePath string `json:"src_relative_path" form:"src_relative_path"` //源：相对路径 + 文件名
	TarRelativePath string `json:"tar_relative_path" form:"tar_relative_path"` //目标：相对路径 + 文件名

}

// @description 查看 SuperVisor 状态
type CicdSuperVisor struct {
	CicdDeploy
	Command string `json:"command"`
}

// @description 部署一个服务
type CicdDeploy struct {
	ServerId  int `json:"server_id" form:"server_id"`   //服务器ID
	ServiceId int `json:"service_id" form:"service_id"` //服务ID
	Flag      int `json:"flag" form:"flag"`             //1本地2远程
}

// @description 同步一个服务
type CicdSync struct {
	ServerId   int    `json:"server_id"`   //服务器ID
	ServiceId  int    `json:"service_id"`  //服务ID
	VersionDir string `json:"version_dir"` //当前代码(版本)的目录
}

// @description 获取声网 token
type TwinAgoraToken struct {
	Username string `json:"username"` //用户名 or 用户ID
	Channel  string `json:"channel"`  //频道名称，给rtc使用,RTM可为空
}

// @description 获取声网的一个 record_id
type TwinAgoraReq struct {
	RecordId int `json:"record_id"`
}

// //能用HTTP请求结构体，按说应该拆开，但是有些公共的参数，拆开不能统一check
//
//	type HttpReqGameMatchPlayerSign struct {
//		RuleId     int                      `json:"rule_id" desc:"ruleId，后台录入的时候，自动生成"`
//		GroupId    int                      `json:"group_id" desc:"小组ID，注：请输入唯一值，不要重复"`
//		PlayerList []HttpReqGameMatchPlayer `json:"player_list" desc:"玩家列表,ex:[{\"uid\":2,\"matchAttr\":{\"age\":1,\"sex\":2}}]"`
//		Addition   string                   `json:"addition" desc:"附加值，请求方传什么值，返回就会随带该值"`
//	}
//
//	type HttpReqGameMatchPlayer struct {
//		Uid        int            `json:"uid"`
//		WeightAttr map[string]int `json:"weight_attr"`
//	}
//
//	type HttpReqGameMatchPlayerCancel struct {
//		RuleId  int `json:"rule_id" desc:"ruleId，后台录入的时候，自动生成"`
//		GroupId int `json:"group_id" desc:"报名时候的，小组ID"`
//	}
type FrameSyncRoomHistory struct {
	RoomId              string `json:"roomId"`
	SequenceNumberStart int32  `json:"sequenceNumberStart"`
	SequenceNumberEnd   int32  `json:"sequenceNumberEnd"`
	SourceUid           int32  `json:"sourceUid"`
	PlayerId            int32  `json:"playerId"`
}

type GrabOrder struct {
	Id         int `json:"id"`
	Amount     int `json:"amount"`
	Uid        int `json:"uid"`
	CategoryId int `json:"category_id"`
	Timeout    int `json:"timeout"`
}

type GrabOrderUserOpen struct {
	PayCategoryId int
	AmountMin     int
	AmountMax     int
}

type Goods struct {
	Title         string `json:"title" db:"define:varchar(50);comment:标题;defaultValue:''"`                   //标题
	Desc          string `json:"desc" db:"define:varchar(255);comment:描述;defaultValue:''"`                   //描述
	Type          int    `json:"type" db:"define:tinyint(1);comment:分类;defaultValue:0"`                      //分类
	Price         int    `json:"price" db:"define:int;comment:价格(单位分);defaultValue:0"`                       //价格
	Amount        int    `json:"period" db:"define:int;comment:买一个商品给多少个数量;defaultValue:0"`                  //买一个商品给多少个数量
	Stock         int    `json:"stock" db:"define:int;comment:库存数量;defaultValue:0"`                          //库存数量
	AdminId       int    `json:"admin_id" db:"define:int;comment:管理员ID;defaultValue:0" `                     //管理员ID
	AllowCoupon   int    `json:"allow_coupon" db:"define:int;comment:允许使用优惠券;defaultValue:0"`                //允许使用优惠券，0不允许
	AllowGoldCoin int    `json:"allow_gold_coin" db:"define:int;comment:允许支付类型1微信2支付宝;defaultValue:0"`       //允许使用金币，0不允许
	PayAllowType  string `json:"pay_allow_type" db:"define:tinyint(1);comment:允许支付类型1微信2支付宝;defaultValue:0"` //允许支付类型1微信2支付宝
	Status        int    `json:"status" db:"define:tinyint(1);comment:状态1:正常 2:下架;defaultValue:0"`           //状态 1:正常 2:下架
	Memo          string `json:"memo" db:"define:varchar(255);comment:备注;defaultValue:''"`                   //备注
}

type Orders struct {
	Type      int    `json:"type" db:"define:tinyint(1);comment:分类;defaultValue:0"`             //分类
	GoodsId   int    `json:"goods_id" db:"define:int;comment:商品ID;defaultValue:0"`              //商品ID
	Amount    int    `json:"amount" db:"define:int;comment:购买数量;defaultValue:0"`                //购买数量
	CouponId  int    `json:"coupon_id" db:"define:int;comment:优惠卷ID;defaultValue:0" `           //优惠卷ID
	Uid       int    `json:"uid" db:"define:int;comment:用户ID;defaultValue:0" `                  //用户ID
	GoldCoin  int    `json:"gold_coin" db:"define:int;comment:金币数;defaultValue:0" `             //金币数
	Source    int    `json:"source" db:"define:tinyint(1);comment:来源;defaultValue:0"`           //来源
	UserMemo  string `json:"user_memo" db:"define:varchar(255);comment:用户备注;defaultValue:''"`   //备注
	AdminMemo string `json:"admin_memo" db:"define:varchar(255);comment:管理员备注;defaultValue:''"` //备注
}

type Payment struct {
	OrdersId int `json:"orders_id"`
	Type     int `json:"type"`
}
