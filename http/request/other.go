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
