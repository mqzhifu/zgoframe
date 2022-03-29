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
	list["PURPOSE_REGISTER"] = PURPOSE_REGISTER
	list["PURPOSE_FIND_PASSWORD"] = PURPOSE_FIND_BACK_PASSWORD
	list["PURPOSE_SET_PASSWORD"] = PURPOSE_SET_PASSWORD
	list["PURPOSE_LOGIN"] = PURPOSE_LOGIN
	list["PURPOSE_SET_MOBILE"] = PURPOSE_SET_MOBILE
	list["PURPOSE_SET_EMAIL"] = PURPOSE_SET_EMAIL
	list["PURPOSE_SET_PAY_PASSWORD"] = PURPOSE_SET_PAY_PASSWORD

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
	list["PLATFORM_MAC_PC_BROWSER"] = PLATFORM_MAC_PC_BROWSER
	list["PLATFORM_MAC_APP"] = PLATFORM_MAC_APP

	list["PLATFORM_WIN_PC_BROWSER"] = PLATFORM_WIN_PC_BROWSER
	list["PLATFORM_WIN_APP"] = PLATFORM_WIN_APP

	list["PLATFORM_ANDROID_H5_BROWSER"] = PLATFORM_ANDROID_H5_BROWSER
	list["PLATFORM_ANDROID_APP"] = PLATFORM_ANDROID_APP

	list["PLATFORM_IOS_H5_BROWSER"] = PLATFORM_IOS_H5_BROWSER
	list["PLATFORM_IOS_APP"] = PLATFORM_IOS_APP

	list["PLATFORM_UNKNOW"] = PLATFORM_UNKNOW

	return list
}

const (
	RULE_PERIORD_MIN = 30

	AUTH_CODE_STATUS_NORMAL = 1
	AUTH_CODE_STATUS_EXPIRE = 3
	AUTH_CODE_STATUS_OK     = 2
)

func GetConstListAuthCode() map[string]int {
	list := make(map[string]int)
	list["AUTH_CODE_STATUS_NORMAL"] = AUTH_CODE_STATUS_NORMAL
	list["AUTH_CODE_STATUS_EXPIRE"] = AUTH_CODE_STATUS_EXPIRE
	list["AUTH_CODE_STATUS_OK"] = AUTH_CODE_STATUS_OK

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
	list["USER_REG_TYPE_EMAIL"] = USER_REG_TYPE_EMAIL
	list["USER_REG_TYPE_NAME"] = USER_REG_TYPE_NAME
	list["USER_REG_TYPE_MOBILE"] = USER_REG_TYPE_MOBILE
	list["USER_REG_TYPE_GUEST"] = USER_REG_TYPE_GUEST

	return list
}

func GetConstListUserTypeThird() map[string]int {
	list := make(map[string]int)
	list["USER_TYPE_THIRD_WEIBO"] = USER_TYPE_THIRD_WEIBO
	list["USER_TYPE_THIRD_WECHAT"] = USER_TYPE_THIRD_WECHAT
	list["USER_TYPE_THIRD_FACEBOOK"] = USER_TYPE_THIRD_FACEBOOK
	list["USER_TYPE_THIRD_GOOGLE"] = USER_TYPE_THIRD_GOOGLE
	list["USER_TYPE_THIRD_TWITTER"] = USER_TYPE_THIRD_TWITTER
	list["USER_TYPE_THIRD_YOUTOBE"] = USER_TYPE_THIRD_YOUTOBE
	list["USER_TYPE_THIRD_QQ"] = USER_TYPE_THIRD_QQ

	return list
}

func GetConstListUserTypeThirdCN() map[string]int {
	list := make(map[string]int)
	list["USER_TYPE_THIRD_WEIBO"] = USER_TYPE_THIRD_WEIBO
	list["USER_TYPE_THIRD_WECHAT"] = USER_TYPE_THIRD_WECHAT
	list["USER_TYPE_THIRD_QQ"] = USER_TYPE_THIRD_QQ

	return list
}

func GetConstListUserTypeThirdNotCN() map[string]int {
	list := make(map[string]int)
	list["USER_TYPE_THIRD_FACEBOOK"] = USER_TYPE_THIRD_FACEBOOK
	list["USER_TYPE_THIRD_GOOGLE"] = USER_TYPE_THIRD_GOOGLE
	list["USER_TYPE_THIRD_TWITTER"] = USER_TYPE_THIRD_TWITTER
	list["USER_TYPE_THIRD_YOUTOBE"] = USER_TYPE_THIRD_YOUTOBE

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
	list["SEX_MALE"] = SEX_MALE
	list["SEX_FEMALE"] = SEX_FEMALE

	return list
}

func GetConstListUserStatus() map[string]int {
	list := make(map[string]int)
	list["USER_STATUS_NOMAL"] = USER_STATUS_NOMAL
	list["USER_STATUS_DENY"] = USER_STATUS_DENY

	return list
}

func GetConstListUserGuest() map[string]int {
	list := make(map[string]int)
	list["USER_GUEST_TRUE"] = USER_GUEST_TRUE
	list["USER_GUEST_FALSE"] = USER_GUEST_FALSE

	return list
}

func GetConstListUserRobot() map[string]int {
	list := make(map[string]int)
	list["USER_ROBOT_TRUE"] = USER_ROBOT_TRUE
	list["USER_ROBOT_FALSE"] = USER_ROBOT_FALSE

	return list
}

func GetConstListUserTest() map[string]int {
	list := make(map[string]int)
	list["USER_TEST_TRUE"] = USER_TEST_TRUE
	list["USER_TEST_FALSE"] = USER_TEST_FALSE

	return list
}

const (
	RULE_TYHP_AUTH_CODE = 1 //验证码
	RULE_TYHP_NOTIFY    = 2 //通知
	RULE_TYHP_MAKE      = 3 //市场营销
)

func GetConstListRuleType() map[string]int {
	list := make(map[string]int)

	list["RULE_TYHP_AUTH_CODE"] = RULE_TYHP_AUTH_CODE
	list["RULE_TYHP_NOTIFY"] = RULE_TYHP_NOTIFY
	list["RULE_TYHP_MAKE"] = RULE_TYHP_MAKE

	return list
}

const (
	SMS_CHANNEL_ALI     = 1 //阿里
	SMS_CHANNEL_TENCENT = 2 //腾讯
)

func GetConstListSmsChannel() map[string]int {
	list := make(map[string]int)

	list["SMS_CHANNEL_ALI"] = SMS_CHANNEL_ALI
	list["SMS_CHANNEL_TENCENT"] = SMS_CHANNEL_TENCENT

	return list
}
