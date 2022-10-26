package model

const (
	RULE_PERIORD_MIN = 30 //
	CHANNEL_DEFAULT  = 1  //
)

//@parse 用户角色
const (
	USER_ROLE_NORMAL = 1 //普通用户
	USER_ROLE_DOCTOR = 2 //专家
)

func CheckConstInList(list map[string]int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//@parse 短信发送目的
const (
	PURPOSE_REGISTER           = 11 //注册
	PURPOSE_LOGIN              = 12 //登陆
	PURPOSE_FIND_BACK_PASSWORD = 13 //找回密码
	PURPOSE_SET_PASSWORD       = 21 //设置密码
	PURPOSE_SET_MOBILE         = 22 //设置手机号
	PURPOSE_SET_EMAIL          = 23 //设置邮件
	PURPOSE_SET_PAY_PASSWORD   = 24 //设置支付密码
)

//func GetConstListPurpose() map[string]int {
//	list := make(map[string]int)
//	list["注册"] = PURPOSE_REGISTER
//	list["找回密码"] = PURPOSE_FIND_BACK_PASSWORD
//	list["设置密码"] = PURPOSE_SET_PASSWORD
//	list["登陆"] = PURPOSE_LOGIN
//	list["设置手机号"] = PURPOSE_SET_MOBILE
//	list["设置邮件"] = PURPOSE_SET_EMAIL
//	list["设置支付密码"] = PURPOSE_SET_PAY_PASSWORD
//
//	return list
//}

//@parse 平台类型
const (
	PLATFORM_MAC_PC_BROWSER     = 11 //MAC台式浏览器
	PLATFORM_MAC_APP            = 12 //MAC台式APP
	PLATFORM_WIN_PC_BROWSER     = 22 //WIN台式浏览器
	PLATFORM_WIN_APP            = 23 //WIN台式APP
	PLATFORM_ANDROID_H5_BROWSER = 31 //安卓手机浏览器
	PLATFORM_ANDROID_APP        = 32 //安卓手机APP
	PLATFORM_IOS_H5_BROWSER     = 41 //IOS手机浏览器
	PLATFORM_IOS_APP            = 42 //IOS手机APP
	PLATFORM_UNKNOW             = 99 //未知
)

//func GetConstListPlatform() map[string]int {
//	list := make(map[string]int)
//	list["MAC台式浏览器"] = PLATFORM_MAC_PC_BROWSER
//	list["MAC台式APP"] = PLATFORM_MAC_APP
//
//	list["WIN台式浏览器"] = PLATFORM_WIN_PC_BROWSER
//	list["WIN台式APP"] = PLATFORM_WIN_APP
//
//	list["安卓手机浏览器"] = PLATFORM_ANDROID_H5_BROWSER
//	list["安卓手机APP"] = PLATFORM_ANDROID_APP
//
//	list["IOS手机浏览器"] = PLATFORM_IOS_H5_BROWSER
//	list["IOS手机APP"] = PLATFORM_IOS_APP
//
//	list["未知"] = PLATFORM_UNKNOW
//
//	return list
//}

//@parse 验证码状态
const (
	AUTH_CODE_STATUS_NORMAL = 1 //正常
	AUTH_CODE_STATUS_EXPIRE = 3 //失效
	AUTH_CODE_STATUS_OK     = 2 //已使用
)

//func GetConstListAuthCodeStatus() map[string]int {
//	list := make(map[string]int)
//	list["未使用"] = AUTH_CODE_STATUS_NORMAL
//	list["已失效"] = AUTH_CODE_STATUS_EXPIRE
//	list["已使用"] = AUTH_CODE_STATUS_OK
//
//	return list
//}

//@parse 发送消息(sms 邮件 站内信)状态
const (
	SEND_MSG_STATUS_OK   = 1 //成功
	SEND_MSG_STATUS_FAIL = 2 //失败
	SEND_MSG_STATUS_ING  = 3 //发送中
	SEND_MSG_STATUS_WAIT = 4 //等待发送
)

//func GetConstListSendMsgStatus() map[string]int {
//	list := make(map[string]int)
//	list["成功"] = SEND_MSG_STATUS_OK
//	list["失败"] = SEND_MSG_STATUS_FAIL
//	list["发送中"] = SEND_MSG_STATUS_ING
//	list["等待发送"] = SEND_MSG_STATUS_WAIT
//
//	return list
//}

//@parse 用户注册类型
const (
	USER_TYPE_THIRD_FACEBOOK = 21 //facebook
	USER_TYPE_THIRD_GOOGLE   = 22 //google
	USER_TYPE_THIRD_TWITTER  = 23 //twitter
	USER_TYPE_THIRD_YOUTOBE  = 24 //youtobe

	USER_TYPE_THIRD_WEIBO  = 11 //微博
	USER_TYPE_THIRD_WECHAT = 12 //微信
	USER_TYPE_THIRD_QQ     = 13 //QQ
)

//@parse 用户注册类型
const (
	USER_REG_TYPE_EMAIL  = 1 //邮件
	USER_REG_TYPE_NAME   = 2 //用户名
	USER_REG_TYPE_MOBILE = 3 //手机
	USER_REG_TYPE_THIRD  = 4 //3方平台
	USER_REG_TYPE_GUEST  = 5 //游客
)

//
//func GetConstListUserRegType() map[string]int {
//	list := make(map[string]int)
//	list["邮件"] = USER_REG_TYPE_EMAIL
//	list["用户名"] = USER_REG_TYPE_NAME
//	list["手机"] = USER_REG_TYPE_MOBILE
//	list["游客"] = USER_REG_TYPE_GUEST
//
//	return list
//}
//
//func GetConstListUserTypeThird() map[string]int {
//	list := make(map[string]int)
//	list["微博"] = USER_TYPE_THIRD_WEIBO
//	list["微信"] = USER_TYPE_THIRD_WECHAT
//	list["facebook"] = USER_TYPE_THIRD_FACEBOOK
//	list["google"] = USER_TYPE_THIRD_GOOGLE
//	list["twitter"] = USER_TYPE_THIRD_TWITTER
//	list["youtobe"] = USER_TYPE_THIRD_YOUTOBE
//	list["qq"] = USER_TYPE_THIRD_QQ
//
//	return list
//}
//
//func GetConstListUserTypeThirdCN() map[string]int {
//	list := make(map[string]int)
//	list["微博"] = USER_TYPE_THIRD_WEIBO
//	list["微信"] = USER_TYPE_THIRD_WECHAT
//	list["qq"] = USER_TYPE_THIRD_QQ
//
//	return list
//}
//
//func GetConstListUserTypeThirdNotCN() map[string]int {
//	list := make(map[string]int)
//	list["facebook"] = USER_TYPE_THIRD_FACEBOOK
//	list["google"] = USER_TYPE_THIRD_GOOGLE
//	list["twitter"] = USER_TYPE_THIRD_TWITTER
//	list["youtobe"] = USER_TYPE_THIRD_YOUTOBE
//
//	return list
//}
////
////func GetUserThirdTypeList() []int {
////	UserThirdType := []int{USER_TYPE_THIRD_WEIBO, USER_TYPE_THIRD_WECHAT, USER_TYPE_THIRD_FACEBOOK, USER_TYPE_THIRD_GOOGLE, USER_TYPE_THIRD_TWITTER, USER_TYPE_THIRD_YOUTOBE, USER_TYPE_THIRD_QQ}
////	return UserThirdType
////}
////
////func GetUserRegTypeList() []int {
////	UserThirdType := []int{USER_REG_TYPE_EMAIL, USER_REG_TYPE_NAME, USER_REG_TYPE_MOBILE, USER_REG_TYPE_THIRD, USER_REG_TYPE_THIRD}
////	return UserThirdType
////}

//@parse 性别
const (
	SEX_MALE   = 1 //男
	SEX_FEMALE = 2 //女
)

//@parse 用户状态
const (
	USER_STATUS_NOMAL = 1 //正常
	USER_STATUS_DENY  = 2 //禁止
)

//@parse 用户是否为游客
const (
	USER_GUEST_TRUE  = 1 //是
	USER_GUEST_FALSE = 2 //否
)

//@parse 用户是否为机器人
const (
	USER_ROBOT_TRUE  = 1 //是
	USER_ROBOT_FALSE = 2 //否
)

//@parse 用户是否为测试账号
const (
	USER_TEST_TRUE  = 1 //是
	USER_TEST_FALSE = 2 //否
)

//func GetConstListUserSex() map[string]int {
//	list := make(map[string]int)
//	list["男"] = SEX_MALE
//	list["女"] = SEX_FEMALE
//
//	return list
//}
//
//func GetConstListUserStatus() map[string]int {
//	list := make(map[string]int)
//	list["正常"] = USER_STATUS_NOMAL
//	list["禁止"] = USER_STATUS_DENY
//
//	return list
//}
//
//func GetConstListUserGuest() map[string]int {
//	list := make(map[string]int)
//	list["是"] = USER_GUEST_TRUE
//	list["否"] = USER_GUEST_FALSE
//
//	return list
//}
//
//func GetConstListUserRobot() map[string]int {
//	list := make(map[string]int)
//	list["是"] = USER_ROBOT_TRUE
//	list["否"] = USER_ROBOT_FALSE
//
//	return list
//}
//
//func GetConstListUserTest() map[string]int {
//	list := make(map[string]int)
//	list["是"] = USER_TEST_TRUE
//	list["否"] = USER_TEST_FALSE
//
//	return list
//}

//@parse 消息rule类型
const (
	RULE_TYPE_AUTH_CODE = 1 //验证码
	RULE_TYPE_NOTIFY    = 2 //通知
	RULE_TYPE_MAKE      = 3 //市场营销
)

//func GetConstListRuleType() map[string]int {
//	list := make(map[string]int)
//
//	list["验证码"] = RULE_TYPE_AUTH_CODE
//	list["通知"] = RULE_TYPE_NOTIFY
//	list["营销"] = RULE_TYPE_MAKE
//
//	return list
//}

//@parse 短信3方平台
const (
	SMS_CHANNEL_ALI     = 1 //阿里
	SMS_CHANNEL_TENCENT = 2 //腾讯
)

//func GetConstListSmsChannel() map[string]int {
//	list := make(map[string]int)
//
//	list["阿里"] = SMS_CHANNEL_ALI
//	list["腾讯"] = SMS_CHANNEL_TENCENT
//
//	return list
//}

//@parse 服务器3方平台
const (
	SERVER_PLATFORM_SELF     = 1 //阿里
	SERVER_PLATFORM_TENGCENT = 2 //腾讯
	SERVER_PLATFORM_ALI      = 3 //阿里
	SERVER_PLATFORM_HUAWEI   = 4 //华为
)

//func GetConstListServerPlatform() map[string]int {
//	list := make(map[string]int)
//
//	list["自家"] = SERVER_PLATFORM_SELF
//	list["腾讯"] = SERVER_PLATFORM_TENGCENT
//	list["阿里"] = SERVER_PLATFORM_ALI
//	list["华为"] = SERVER_PLATFORM_HUAWEI
//
//	return list
//}

//@parse CICD部署状态
const (
	CICD_PUBLISH_DEPLOY_STATUS_ING    = 1 //发布中
	CICD_PUBLISH_DEPLOY_STATUS_FAIL   = 2 //发布失败
	CICD_PUBLISH_DEPLOY_STATUS_FINISH = 3 //发布结束/完
)

//
//func GetConstListCicdPublishDeployStatus() map[string]int {
//	list := make(map[string]int)
//
//	list["部署中"] = CICD_PUBLISH_DEPLOY_STATUS_ING
//	list["失败"] = CICD_PUBLISH_DEPLOY_STATUS_FAIL
//	list["完成"] = CICD_PUBLISH_DEPLOY_STATUS_FINISH
//
//	return list
//}

//@parse CICD发布状态
const (
	CICD_PUBLISH_STATUS_WAIT_DEPLOY = 1 //待部署
	CICD_PUBLISH_STATUS_WAIT_PUB    = 2 //待发布
	CICD_PUBLISH_DEPLOY_OK          = 3 //成功
	CICD_PUBLISH_DEPLOY_FAIL        = 4 //失败
)

//func GetConstListCicdPublishStatus() map[string]int {
//	list := make(map[string]int)
//
//	list["待部署"] = CICD_PUBLISH_STATUS_WAIT_DEPLOY
//	list["待发布"] = CICD_PUBLISH_STATUS_WAIT_PUB
//	list["成功"] = CICD_PUBLISH_DEPLOY_OK
//	list["失败"] = CICD_PUBLISH_DEPLOY_OK
//
//	return list
//}

//@parse 项目状态
const (
	PROJECT_STATUS_OPEN  = 1 //打开
	PROJECT_STATUS_CLOSE = 2 //关闭
)

//func GetConstListProjectStatus() map[string]int {
//	list := make(map[string]int)
//
//	list["打开"] = PROJECT_STATUS_OPEN
//	list["关闭"] = PROJECT_STATUS_CLOSE
//
//	return list
//}

//@parse 项目大类型
const (
	PROJECT_TYPE_SERVICE = 1 //服务
	PROJECT_TYPE_FE      = 2 //前端
	PROJECT_TYPE_APP     = 3 //APP
	PROJECT_TYPE_BE      = 4 //后端
)

//func GetConstListProjectType() map[string]int {
//	list := make(map[string]int)
//	list["服务"] = PROJECT_TYPE_SERVICE
//	list["前端"] = PROJECT_TYPE_FE
//	list["APP"] = PROJECT_TYPE_APP
//	list["后端"] = PROJECT_TYPE_BE
//	return list
//}

//@parse 项目开发语言类型
const (
	PROJECT_LANG_PHP     = 1 //PHP
	PROJECT_LANG_GO      = 2 //GO
	PROJECT_LANG_JAVA    = 3 //JAVA
	PROJECT_LANG_JS      = 4 //JS
	PROJECT_LANG_C_PLUS  = 5 //C++
	PROJECT_LANG_C       = 6 //C
	PROJECT_LANG_C_SHARP = 7 //C#
)

//func GetConstListProjectLanguage() map[string]int {
//	list := make(map[string]int)
//	list["php"] = PROJECT_LANG_PHP
//	list["golang"] = PROJECT_LANG_GO
//	list["java"] = PROJECT_LANG_JAVA
//	list["js"] = PROJECT_LANG_JS
//	list["c++"] = PROJECT_LANG_C_PLUS
//	list["c"] = PROJECT_LANG_C
//	list["c#"] = PROJECT_LANG_C_SHARP
//	return list
//}

//@parse 声网录制屏幕-获取资源的状态
const (
	AGORA_CLOUD_RECORD_STATUS_RESOURCE = 1 //已获取资源ID
	AGORA_CLOUD_RECORD_STATUS_START    = 2 //已开始
	AGORA_CLOUD_RECORD_STATUS_END      = 3 //已结束
)

//func GetConstListAgoraCloudRecordStatus() map[string]int {
//	list := make(map[string]int)
//	list["已获取资源ID"] = AGORA_CLOUD_RECORD_STATUS_RESOURCE
//	list["已开始"] = AGORA_CLOUD_RECORD_STATUS_START
//	list["已结束"] = AGORA_CLOUD_RECORD_STATUS_END
//
//	return list
//}

//@parse 声网录制屏幕-停止的类型
const (
	AGORA_CLOUD_RECORD_STOP_ACTION_UNKNOW   = 0 //未知
	AGORA_CLOUD_RECORD_STOP_ACTION_NORMAL   = 1 //正常停止
	AGORA_CLOUD_RECORD_STOP_ACTION_RELOAD   = 2 //页面刷新时拦截
	AGORA_CLOUD_RECORD_STOP_ACTION_REENTER  = 3 //重新加载页面触发
	AGORA_CLOUD_RECORD_STOP_ACTION_CALLBACK = 4 //声网回调触发
)

//@parse 声网录制屏幕-服务状态
const (
	AGORA_CLOUD_RECORD_SERVER_STATUS_UNDO = 1 //未处理
	AGORA_CLOUD_RECORD_SERVER_STATUS_ING  = 2 //处理中
	AGORA_CLOUD_RECORD_SERVER_STATUS_OK   = 3 //处理成功
	AGORA_CLOUD_RECORD_SERVER_STATUS_ERR  = 4 //处理异常
)

//func GetConstListAgoraCloudRecordServerStatus() map[string]int {
//	list := make(map[string]int)
//	list["未处理"] = AGORA_CLOUD_RECORD_SERVER_STATUS_UNDO
//	list["处理中"] = AGORA_CLOUD_RECORD_SERVER_STATUS_ING
//	list["处理成功"] = AGORA_CLOUD_RECORD_SERVER_STATUS_OK
//	list["处理异常"] = AGORA_CLOUD_RECORD_SERVER_STATUS_ERR
//
//	return list
//}

//@parse 声网录制屏幕-回调状态
const (
	CallbackEventAllUploaded = 31 // 所有录制文件已上传至指定的第三方云存储
	CallbackEventRecordExit  = 41 // 录制服务已退出
)
