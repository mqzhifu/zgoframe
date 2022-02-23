package v1


/*
	需要实现的一些功能接口，列表
*/

type Project interface {
	GetInfo()
	GetAppInfo()
}

type TaskInterface interface {
	GetListByUid(uid int)
}

type UserCenter interface{
	Reg()//注册
	Login()//用户名|邮箱|手机 登陆
	ThirdLogin()//3方登陆
	QRCodeLogin()//二维码/扫码登陆
	Logout()//登出
	FindBackPSAuthByPhone()//通过手机号找回密码，验证
	FindBackPSResetPS()//找回密码，重置密码
	RestPS()//修改密码
	SetPayPs()//修改支付密码
	BindEmail()//绑定邮件
	BindPhone()//绑定手机
	UpInfo()//修改个人基础信息
	UploadAvatar()//上传头像
	UpIdCard()//更新 身份证信息
	GetInfo()//获取一个用户的基础信息
}

type Friend interface {
	Apply()//申请添加好友
	GetList()//好友列表

	AppliedList()//别人申请添加你的列表
	AppliedDeny()//拒绝别人的申请
	AppliedOk()//通过别人的申请

	Remove()//解除好友关系

	AddBlack()
	GetBlackList()
	CancelBlack()
}

type Room interface{
	GetRoomCategory()//房间有不同类型
	Entry()//进入房间
	Ready()//准备
}

type GameMatch interface{
	Add()
	Cancel()
}

type FrameSync interface {

}

type IMInterface interface {
	SendOneMsg()
	EntrySession()
	DeleteSession()
}

type Game interface{
	GetList()//获取所有游戏列表
	GetHistoryByUid()//获取一个用户玩过的所有游戏列表
	GetMatchHistoryByUid()//获取一个用户玩过的所有游戏的战绩列表
}

type Pay interface{
	GetCategory()
	Doing()
}

type Goods interface {
	GetList()
	UpInfo()
	GetOne()
}
//货币
type Currency interface {
	AddCoin()
	LessCoin()
	GetCoinBalance()
	GetCoinHistory()
}

type Order interface{
	Create()
	PayCallback()
	UpInfo()
}

type MsgCenter interface{
	//微信 企业微信/公众号
	//app msg
	//站内信 | 内部邮件
	//短信
	//邮件
	//IM
	//task
	//game
	//pay
	//friend

	//长连接
	Connect()//接收到C端建立长连接
	Login()//连接建立后，做验证
	Heartbeat()//心跳
	Close()//C端主动关闭连接
	PushLong()//接收到C端推送消息

	//短连接
	Push()
}