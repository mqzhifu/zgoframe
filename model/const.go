package model

func CheckConstInList(list map[string]int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

const (
	PURPOSE_REGISTER           = 11
	PURPOSE_LOGIN              = 12
	PURPOSE_FIND_BACK_PASSWORD = 13
	PURPOSE_SET_PASSWORD       = 21
	PURPOSE_SET_MOBILE         = 22
	PURPOSE_SET_EMAIL          = 23
	PURPOSE_SET_PAY_PASSWORD   = 24
)

func GetConstListPurpose() map[string]int {
	list := make(map[string]int)
	list["注册"] = PURPOSE_REGISTER
	list["找回密码"] = PURPOSE_FIND_BACK_PASSWORD
	list["设置密码"] = PURPOSE_SET_PASSWORD
	list["登陆"] = PURPOSE_LOGIN
	list["设置手机号"] = PURPOSE_SET_MOBILE
	list["设置邮件"] = PURPOSE_SET_EMAIL
	list["设置支付密码"] = PURPOSE_SET_PAY_PASSWORD

	return list
}

const (
	PLATFORM_MAC_PC_BROWSER = 11
	PLATFORM_MAC_APP        = 12

	PLATFORM_WIN_PC_BROWSER = 22
	PLATFORM_WIN_APP        = 23

	PLATFORM_ANDROID_H5_BROWSER = 31
	PLATFORM_ANDROID_APP        = 32

	PLATFORM_IOS_H5_BROWSER = 41
	PLATFORM_IOS_APP        = 42

	PLATFORM_UNKNOW = 99
)

func GetConstListPlatform() map[string]int {
	list := make(map[string]int)
	list["MAC台式浏览器"] = PLATFORM_MAC_PC_BROWSER
	list["MAC台式APP"] = PLATFORM_MAC_APP

	list["WIN台式浏览器"] = PLATFORM_WIN_PC_BROWSER
	list["WIN台式APP"] = PLATFORM_WIN_APP

	list["安卓手机浏览器"] = PLATFORM_ANDROID_H5_BROWSER
	list["安卓手机APP"] = PLATFORM_ANDROID_APP

	list["IOS手机浏览器"] = PLATFORM_IOS_H5_BROWSER
	list["IOS手机APP"] = PLATFORM_IOS_APP

	list["未知"] = PLATFORM_UNKNOW

	return list
}

const (
	RULE_PERIORD_MIN = 30

	AUTH_CODE_STATUS_NORMAL = 1
	AUTH_CODE_STATUS_EXPIRE = 3
	AUTH_CODE_STATUS_OK     = 2
)

func GetConstListAuthCodeStatus() map[string]int {
	list := make(map[string]int)
	list["未使用"] = AUTH_CODE_STATUS_NORMAL
	list["已失效"] = AUTH_CODE_STATUS_EXPIRE
	list["已使用"] = AUTH_CODE_STATUS_OK

	return list
}

const (
	USER_REG_TYPE_EMAIL  = 1
	USER_REG_TYPE_NAME   = 2
	USER_REG_TYPE_MOBILE = 3
	USER_REG_TYPE_THIRD  = 4
	USER_REG_TYPE_GUEST  = 5

	USER_TYPE_THIRD_WEIBO  = 11
	USER_TYPE_THIRD_WECHAT = 12
	USER_TYPE_THIRD_QQ     = 13

	USER_TYPE_THIRD_FACEBOOK = 21
	USER_TYPE_THIRD_GOOGLE   = 22
	USER_TYPE_THIRD_TWITTER  = 23
	USER_TYPE_THIRD_YOUTOBE  = 24

	CHANNEL_DEFAULT = 1
)

func GetConstListUserRegType() map[string]int {
	list := make(map[string]int)
	list["邮件"] = USER_REG_TYPE_EMAIL
	list["用户名"] = USER_REG_TYPE_NAME
	list["手机"] = USER_REG_TYPE_MOBILE
	list["游客"] = USER_REG_TYPE_GUEST

	return list
}

func GetConstListUserTypeThird() map[string]int {
	list := make(map[string]int)
	list["微博"] = USER_TYPE_THIRD_WEIBO
	list["微信"] = USER_TYPE_THIRD_WECHAT
	list["facebook"] = USER_TYPE_THIRD_FACEBOOK
	list["google"] = USER_TYPE_THIRD_GOOGLE
	list["twitter"] = USER_TYPE_THIRD_TWITTER
	list["youtobe"] = USER_TYPE_THIRD_YOUTOBE
	list["qq"] = USER_TYPE_THIRD_QQ

	return list
}

func GetConstListUserTypeThirdCN() map[string]int {
	list := make(map[string]int)
	list["微博"] = USER_TYPE_THIRD_WEIBO
	list["微信"] = USER_TYPE_THIRD_WECHAT
	list["qq"] = USER_TYPE_THIRD_QQ

	return list
}

func GetConstListUserTypeThirdNotCN() map[string]int {
	list := make(map[string]int)
	list["facebook"] = USER_TYPE_THIRD_FACEBOOK
	list["google"] = USER_TYPE_THIRD_GOOGLE
	list["twitter"] = USER_TYPE_THIRD_TWITTER
	list["youtobe"] = USER_TYPE_THIRD_YOUTOBE

	return list
}

//func GetUserThirdTypeList() []int {
//	UserThirdType := []int{USER_TYPE_THIRD_WEIBO, USER_TYPE_THIRD_WECHAT, USER_TYPE_THIRD_FACEBOOK, USER_TYPE_THIRD_GOOGLE, USER_TYPE_THIRD_TWITTER, USER_TYPE_THIRD_YOUTOBE, USER_TYPE_THIRD_QQ}
//	return UserThirdType
//}
//
//func GetUserRegTypeList() []int {
//	UserThirdType := []int{USER_REG_TYPE_EMAIL, USER_REG_TYPE_NAME, USER_REG_TYPE_MOBILE, USER_REG_TYPE_THIRD, USER_REG_TYPE_THIRD}
//	return UserThirdType
//}

const (
	SEX_MALE   = 1
	SEX_FEMALE = 2

	USER_STATUS_NOMAL = 1
	USER_STATUS_DENY  = 2

	USER_GUEST_TRUE  = 1
	USER_GUEST_FALSE = 2

	USER_ROBOT_TRUE  = 1
	USER_ROBOT_FALSE = 2

	USER_TEST_TRUE  = 1
	USER_TEST_FALSE = 2
)

func GetConstListUserSex() map[string]int {
	list := make(map[string]int)
	list["男"] = SEX_MALE
	list["女"] = SEX_FEMALE

	return list
}

func GetConstListUserStatus() map[string]int {
	list := make(map[string]int)
	list["正常"] = USER_STATUS_NOMAL
	list["禁止"] = USER_STATUS_DENY

	return list
}

func GetConstListUserGuest() map[string]int {
	list := make(map[string]int)
	list["是"] = USER_GUEST_TRUE
	list["否"] = USER_GUEST_FALSE

	return list
}

func GetConstListUserRobot() map[string]int {
	list := make(map[string]int)
	list["是"] = USER_ROBOT_TRUE
	list["否"] = USER_ROBOT_FALSE

	return list
}

func GetConstListUserTest() map[string]int {
	list := make(map[string]int)
	list["是\""] = USER_TEST_TRUE
	list["否"] = USER_TEST_FALSE

	return list
}

const (
	RULE_TYPE_AUTH_CODE = 1 //验证码
	RULE_TYPE_NOTIFY    = 2 //通知
	RULE_TYPE_MAKE      = 3 //市场营销
)

func GetConstListRuleType() map[string]int {
	list := make(map[string]int)

	list["验证码"] = RULE_TYPE_AUTH_CODE
	list["通知"] = RULE_TYPE_NOTIFY
	list["营销"] = RULE_TYPE_MAKE

	return list
}

const (
	SMS_CHANNEL_ALI     = 1 //阿里
	SMS_CHANNEL_TENCENT = 2 //腾讯
)

func GetConstListSmsChannel() map[string]int {
	list := make(map[string]int)

	list["阿里"] = SMS_CHANNEL_ALI
	list["腾讯"] = SMS_CHANNEL_TENCENT

	return list
}

const (
	SERVER_PLATFORM_SELF     = 1 //阿里
	SERVER_PLATFORM_TENGCENT = 2 //腾讯
	SERVER_PLATFORM_ALI      = 3
	SERVER_PLATFORM_HUAWEI   = 4
)

func GetConstListServerPlatform() map[string]int {
	list := make(map[string]int)

	list["自家"] = SERVER_PLATFORM_SELF
	list["腾讯"] = SERVER_PLATFORM_TENGCENT
	list["阿里"] = SERVER_PLATFORM_ALI
	list["华为"] = SERVER_PLATFORM_HUAWEI

	return list
}

const (
	CICD_PUBLISH_STATUS_ING  = 1
	CICD_PUBLISH_STATUS_FAIL = 2
	CICD_PUBLISH_STATUS_OK   = 3
)

func GetConstListCicdPublishStatus() map[string]int {
	list := make(map[string]int)

	list["发送中"] = CICD_PUBLISH_STATUS_ING
	list["失败"] = CICD_PUBLISH_STATUS_FAIL
	list["成功"] = CICD_PUBLISH_STATUS_OK

	return list
}

const (
	PROJECT_STATUS_OPEN  = 1
	PROJECT_STATUS_CLOSE = 2
)

func GetConstListProjectStatus() map[string]int {
	list := make(map[string]int)

	list["打开"] = PROJECT_STATUS_OPEN
	list["关闭"] = PROJECT_STATUS_CLOSE

	return list
}

const (
	PROJECT_TYPE_SERVICE = 1
	PROJECT_TYPE_FE      = 2
	PROJECT_TYPE_APP     = 3
	PROJECT_TYPE_BE      = 4
)

func GetConstListProjectType() map[string]int {
	list := make(map[string]int)
	list["服务"] = PROJECT_TYPE_SERVICE
	list["前端"] = PROJECT_TYPE_FE
	list["APP"] = PROJECT_TYPE_APP
	list["后端"] = PROJECT_TYPE_BE
	return list
}

//var PROJECT_TYPE_MAP = map[int]string{
//	PROJECT_TYPE_SERVICE: "service",
//	PROJECT_TYPE_FE:      "frontend",
//	PROJECT_TYPE_APP:     "app",
//	PROJECT_TYPE_BE:      "backend",
//}
