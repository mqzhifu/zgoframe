package gamematch
//能用HTTP请求结构体，按说应该拆开，但是有些公共的参数，拆开不能统一check
type HttpReqBusiness struct {
	//公共参数
	MatchCode	string	`json:"matchCode" desc:"一条匹配规则的关键码,后台录入时确定"`
	RuleId 		int 	`json:"ruleId" desc:"同上，这个是整型，由系统自动生成"`
	//报名
	GroupId		int		`json:"groupId" desc:"小组ID，注：请输入唯一值，不要重复"`
	CustomProp	string	`json:"customProp" desc:"暂未使用"`
	PlayerList	[]HttpReqPlayer	`json:"playerList" desc:"玩家列表,ex:[{\"uid\":2,\"matchAttr\":{\"age\":1,\"sex\":2}}]"`
	Addition	string	`json:"addition" desc:"附加值，请求方传什么值，返回就会随带该值"`
	Rule_ver	int		`json:"rule_ver" desc:"一条匹配规则的公式格式-版本"`
	//取消报名
	//GroupId
	//删除一条 匹配成功记录
	SuccessId	int		`json:"success_id"`
}

type HttpReqPlayer struct {
	Uid 		int		`json:"uid" `
	MatchAttr	map[string]int	`json:"matchAttr"`
}

//type HttpReqSign struct {
//	MatchCode	string	`json:"matchCode"`
//	GroupId		int		`json:"groupId"`
//	CustomProp	string	`json:"customProp"`
//	PlayerList	[]HttpReqPlayer	`json:"playerList"`
//	Addition	string	`json:"addition"`
//	RuleId 		int 	`json:"ruleId"`
//	Rule_ver	int		`json:"rule_ver"`
//}
//type HttpClearRule struct {
//	MatchCode	string	`json:"matchCode"`
//}
//type HttpReqPlayerAttr struct {
//
//}
//type HttpReqSignCancel struct {
//	PlayerId	int
//	GroupId 	int
//	RuleId 		int
//	MatchCode	string
//}
//type HttpReqSuccessDel struct {
//	MatchCode	string
//	RuleId 		int
//	Id 	int
//}