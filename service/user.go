package service

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

type UserRegInfo struct {
	Channel   int    `json:"channel"`     //来源渠道
	ThirdType int    `json:"third_type" ` //三方平台类型
	ThirdId   string `json:"third_id"`    //三方平台ID
}

type User struct {
	Gorm  *gorm.DB
	Redis *util.MyRedis
}

func NewUser(gorm *gorm.DB, redis *util.MyRedis) *User {
	user := new(User)
	user.Gorm = gorm
	user.Redis = redis
	return user
}

//注册，用户名/密码
func (user *User) RegisterByUsername(R request.Register, h request.HeaderRequest) (err error, userInter model.User) {
	u := model.User{
		Username:  R.Username,
		NickName:  R.NickName,
		Password:  R.Password,
		HeaderImg: R.HeaderImg,
		ProjectId: R.ProjectId,
		Sex:       R.Sex,
		Recommend: R.Recommend,
		Guest:     R.Guest,
		Robot:     model.USER_ROBOT_FALSE,
	}

	if u.Guest != model.USER_GUEST_TRUE && u.Guest != model.USER_GUEST_FALSE {
		return errors.New("Guest value err."), userInter
	}

	userRegInfo := UserRegInfo{
		ThirdType: R.ThirdType,
		ThirdId:   R.ThirdId,
		Channel:   R.Channel,
	}

	return user.Register(u, h, userRegInfo)
}

func (user *User) Delete(uid int) (map[string]int, error) {
	rsMap := make(map[string]int)
	err, userInfo := user.FindUserById(uid)
	if err != nil {
		return rsMap, err
	}
	if userInfo.Mobile != "" {
		var smsLog model.SmsLog
		obj := user.Gorm.Unscoped().Where("Receiver = ?", userInfo.Mobile).Delete(&smsLog)
		rsMap["SmsLog"] = int(obj.RowsAffected)
	}

	if userInfo.Email != "" {
		var emailLog model.EmailLog
		obj := user.Gorm.Unscoped().Where("Receiver = ?", userInfo.Email).Delete(&emailLog)
		rsMap["EmailLog"] = int(obj.RowsAffected)
	}
	var userLogin model.UserLogin
	obj := user.Gorm.Unscoped().Where("uid = ?", uid).Delete(&userLogin)
	rsMap["UserLogin"] = int(obj.RowsAffected)

	listPlatform := model.GetConstListPlatform()
	for _, v := range listPlatform {
		redisElement, _ := user.Redis.GetElementByIndex("jwt", strconv.Itoa(v), strconv.Itoa(uid))
		user.Redis.Del(redisElement)
	}

	var userReg model.UserReg
	obj = user.Gorm.Unscoped().Where("uid = ?", uid).Delete(&userReg)
	rsMap["UserReg"] = int(obj.RowsAffected)

	var u model.User
	user.Gorm.Unscoped().Delete(&u, uid)
	rsMap["User"] = int(obj.RowsAffected)

	util.MyPrint(rsMap)

	return rsMap, nil
}

//最终 - 注册
func (user *User) Register(formUser model.User, h request.HeaderRequest, userRegInfo UserRegInfo) (err error, userInter model.User) {
	var userRegType int

	formUser.Status = model.USER_STATUS_NOMAL

	if formUser.Test <= 0 {
		formUser.Test = model.USER_TEST_FALSE
	}

	if formUser.Guest == model.USER_REG_TYPE_GUEST {
		util.MyPrint("reg in GUEST")
		//deviceId = username
		if formUser.Username == "" {
			formUser.Username = MakeGuestUsername()
		}
		_, exist, _ := user.FindUserByUsername(formUser.Username)
		if exist {
			return errors.New("username 已注册:" + formUser.Username), userInter
		}
		userRegType = model.USER_REG_TYPE_NAME
	} else if userRegInfo.ThirdType > 0 && userRegInfo.ThirdId != "" {
		util.MyPrint("reg in Third")
		userRegType = model.USER_REG_TYPE_THIRD
		var userThird model.UserThird
		if !errors.Is(user.Gorm.Where("third_id = ? and platform_type = ?  ", userRegInfo.ThirdId, userRegInfo.ThirdType).First(&userThird).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
			return errors.New("third_id 已注册"), userInter
		}
	} else {
		userRegType = user.TurnRegByUsername(formUser.Username)
		util.MyPrint("reg in USERNAME , type:", userRegType)
		if userRegType == model.USER_REG_TYPE_MOBILE {
			_, empty, _ := user.FindUserByMobile(formUser.Mobile)
			if !empty {
				return errors.New("mobile 已注册:" + formUser.Username), userInter
			}
		} else if userRegType == model.USER_REG_TYPE_EMAIL {
			_, empty, _ := user.FindUserByEmail(formUser.Username)
			if !empty {
				return errors.New("email 已注册:" + formUser.Username), userInter
			}
		} else {
			_, empty, _ := user.FindUserByUsername(formUser.Username)
			if !empty {
				return errors.New("username 已注册:" + formUser.Username), userInter
			}
		}

		userRegType = user.TurnRegByUsername(formUser.Username)
	}

	if formUser.NickName == "" {
		formUser.NickName = MakeNickname()
	}

	if formUser.Password != "" {
		formUser.Password = util.MD5V([]byte(formUser.Password))
	}

	formUser.Uuid = uuid.NewV4().String()

	formatHeader := fmt.Sprintf("%+v", formUser)
	util.MyPrint("create user Info", formatHeader)

	err = user.Gorm.Create(&formUser).Error
	if err != nil {
		return err, userInter
	}
	//util.MyPrint("u.id:",u.Id)
	channel := model.CHANNEL_DEFAULT
	if userRegInfo.Channel > 0 {
		channel = userRegInfo.Channel
	}

	userReg := model.UserReg{
		ProjectId: formUser.ProjectId,
		Uid:       formUser.Id,
		ThirdType: userRegInfo.ThirdType,
		Type:      userRegType,
		Channel:   channel,
		AutoIp:    h.AutoIp,

		Ip:            h.BaseInfo.Ip,
		AppVersion:    h.BaseInfo.AppVersion,
		SourceType:    h.SourceType,
		Os:            h.BaseInfo.OS,
		OsVersion:     h.BaseInfo.OSVersion,
		Device:        h.BaseInfo.Device,
		DeviceVersion: h.BaseInfo.DeviceVersion,
		Lat:           h.BaseInfo.Lat,
		Lon:           h.BaseInfo.Lon,
		DeviceId:      h.BaseInfo.DeviceId,
		Dpi:           h.BaseInfo.DPI,
	}
	util.PrintStruct(userReg, ":")
	//fmt.Sprintf("aaaaf:%+v", &userReg)
	//util.MyPrint("userReg:",userReg)
	err = user.Gorm.Create(&userReg).Error
	if err != nil {
		return errors.New("create user_Reg err:" + err.Error()), userInter
	}

	if userRegType == model.USER_REG_TYPE_THIRD {
		userThird := model.UserThird{
			Uid:          formUser.Id,
			ThirdId:      userRegInfo.ThirdId,
			PlatformType: userRegInfo.ThirdType,
		}

		err = user.Gorm.Create(&userThird).Error
		if err != nil {
			return errors.New("create user_third err:" + err.Error()), userInter
		}
	}

	return nil, formUser
}

//手机号登陆 - 上一步需要 ： 短信验证没问题
func (user *User) LoginSms(mobile string) (userInter model.User, err error) {
	userInter, empty, _ := user.FindUserByMobile(mobile)
	if !empty {
		return userInter, nil
	}

	return userInter, errors.New("手机号不存在DB")
}

//用户名/密码 登陆
func (user *User) Login(u *model.User) (err error, userInter model.User) {
	var userInfo model.User
	if u.Username == "" || u.Password == "" {
		return errors.New("username || Password empty"), userInter
	}

	u.Password = util.MD5V([]byte(u.Password))
	regType := user.TurnRegByUsername(u.Username)
	if regType == model.USER_REG_TYPE_MOBILE {
		err = user.Gorm.Where("mobile = ? AND password = ?   ", u.Mobile, u.Password).First(&userInfo).Error
	} else if regType == model.USER_REG_TYPE_EMAIL {
		err = user.Gorm.Where("email = ? AND password = ?   ", u.Email, u.Password).First(&userInfo).Error
	} else {
		err = user.Gorm.Where("username = ? AND password = ?   ", u.Username, u.Password).First(&userInfo).Error
	}

	if err == nil {
		if userInfo.Status != model.USER_STATUS_NOMAL {
			return errors.New("status err"), userInfo
		}
	}

	return err, userInfo
}

//3方平台登陆 - 不需要密码
func (user *User) LoginThird(rLoginThird request.RLoginThird, h request.HeaderRequest) (userInfo model.User, isNewReg bool, err error) {
	var userThird model.UserThird
	//var userInfo model.User
	if errors.Is(user.Gorm.Where("third_id = ? and platform_type = ?  ", rLoginThird.ThirdId, rLoginThird.PlatformType).First(&userThird).Error, gorm.ErrRecordNotFound) {
		rLoginThird.Guest = model.USER_GUEST_FALSE

		regUserInfo := model.User{
			ProjectId: rLoginThird.ProjectId,
			Username:  rLoginThird.Username,
			NickName:  rLoginThird.NickName,

			HeaderImg: rLoginThird.NickName,
			Sex:       rLoginThird.Sex,
			Birthday:  rLoginThird.Birthday,
			Recommend: rLoginThird.Recommend,
			Guest:     model.USER_GUEST_FALSE,
			Robot:     model.USER_ROBOT_FALSE,
		}

		userRegInfo := UserRegInfo{
			ThirdType: rLoginThird.ThirdType,
			ThirdId:   rLoginThird.ThirdId,
		}

		err, userInfo = user.Register(regUserInfo, h, userRegInfo)
		if err != nil {
			return userInfo, false, err
		}
		return userInfo, true, nil
	} else {
		err = user.Gorm.Where("id = ?   ", userThird.Uid).First(&userInfo).Error
		return userInfo, false, err
	}
}

//根据用户名 判断  ：手机号 用户名 邮箱
func (user *User) TurnRegByUsername(username string) int {
	isEmail := util.CheckEmailRule(username)
	isMobile := util.CheckMobileRule(username)
	userRegType := model.USER_REG_TYPE_NAME
	if isEmail {
		userRegType = model.USER_REG_TYPE_EMAIL
	} else if isMobile {
		userRegType = model.USER_REG_TYPE_MOBILE
	}
	return userRegType
}

//编辑用户基础信息
func (user *User) SetUserInfo(reqUser model.User) (err error, userInfo model.User) {
	err = user.Gorm.Updates(&reqUser).Error
	return err, reqUser
}

//根据ID查找一个用户信息
func (user *User) FindUserById(id int) (err error, userInfo *model.User) {
	var u model.User
	err = user.Gorm.Where("`id` = ?", id).First(&u).Error
	return err, &u
}

func (user *User) FindUserByUsername(username string) (userInfo model.User, empty bool, err error) {
	err = user.Gorm.Where("username = ? ", username).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userInfo, true, nil
	}
	if err != nil {
		return userInfo, false, err
	}
	return userInfo, false, nil
}

func (user *User) FindUserByEmail(email string) (userInfo model.User, empty bool, err error) {
	err = user.Gorm.Where("email = ? ", email).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userInfo, true, nil
	}
	if err != nil {
		return userInfo, false, err
	}
	return userInfo, false, nil
}

func (user *User) FindUserByMobile(mobile string) (userInfo model.User, empty bool, err error) {
	err = user.Gorm.Where("mobile = ? ", mobile).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userInfo, true, nil
	}
	if err != nil {
		return userInfo, false, err
	}
	return userInfo, false, nil
}

//根据uuid查找一个用户信息
func (user *User) FindUserByUuid(uuid string) (err error, userInfo *model.User) {
	if err = user.Gorm.Where("`uuid` = ?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("用户不存在"), userInfo
	}
	return nil, userInfo
}

//随机生成一个用户 - 游客
func MakeGuestUsername() string {
	return uuid.NewV4().String()
}

//随机生成一个昵称 - 游客
func MakeNickname() string {
	return uuid.NewV4().String()
}

//修改密码
func (user *User) ChangePassword(uid int, newPassword string) (err error) {
	var userInfo model.User
	userInfo.Id = uid
	userInfo.Password = util.MD5V([]byte(newPassword))
	err = user.Gorm.Updates(&userInfo).Error
	return err
}

//绑定手机号
func (user *User) BindMobile(uid int, mobile string) (err error) {

	_, empty, err := user.FindUserByMobile(mobile)
	if !empty {
		return errors.New("手机号已绑定过了，请不要重复操作")
	}

	//err, _ = user.FindUserById(uid)
	//if err != nil {
	//	return errors.New("uid err:" + err.Error())
	//}
	//
	//if userInfo.Mobile != "" {
	//	return errors.New("该用户已经绑定过手机号，请不要重复操作:" + mobile)
	//}

	var userInfoEdit model.User
	userInfoEdit.Id = uid
	userInfoEdit.Mobile = mobile
	userInfoEdit.Guest = model.USER_GUEST_FALSE

	err = user.Gorm.Updates(&userInfoEdit).Error
	return err
}

//绑定邮箱
func (user *User) BindEmail(uid int, email string) (err error) {

	_, empty, err := user.FindUserByEmail(email)
	if !empty {
		return errors.New("邮箱已绑定过了，请不要重复操作")
	}

	//err, _ = user.FindUserById(uid)
	//if err != nil {
	//	return errors.New("uid err:" + err.Error())
	//}
	//
	//if userInfo.Email != "" {
	//	return errors.New("该用户已经绑定过邮箱，请不要重复操作:" + email)
	//}

	var userInfoEdit model.User
	userInfoEdit.Id = uid
	userInfoEdit.Email = email
	userInfoEdit.Guest = model.USER_GUEST_FALSE

	err = user.Gorm.Updates(&userInfoEdit).Error
	return err
}

type ConstInfo struct {
	List map[string]int
	Key  string
	Name string
}

func (user *User) GetConstList() []ConstInfo {
	ConstList := []ConstInfo{}

	ConstList = append(ConstList, ConstInfo{
		List: util.GetConstListEnv(),
		Name: "env-环境",
		Key:  "ENV",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListProjectType(),
		Name: "项目类型",
		Key:  "PROJECT_TYPE",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListProjectStatus(),
		Name: "项目状态",
		Key:  "PROJECT_STATUS",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListPlatform(),
		Name: "平台类型",
		Key:  "PLATFORM",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserTypeThird(),
		Name: "用户类型3方",
		Key:  "USER_TYPE_THIRD",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserRegType(),
		Name: "用户注册类型",
		Key:  "USER_REG_TYPE",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserSex(),
		Name: "用户性别",
		Key:  "USER_SEX",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserStatus(),
		Name: "用户状态",
		Key:  "USER_STATUS",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserTypeThirdCN(),
		Name: "用户类型-中国",
		Key:  "USER_TYPE_THIRD_CN",
	})
	//
	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserTypeThirdNotCN(),
		Name: "用户类型-外国",
		Key:  "USER_TYPE_THIRD_NOT_CN",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserGuest(),
		Name: "游客分类",
		Key:  "USER_GUEST",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserRobot(),
		Name: "机器人",
		Key:  "USER_ROBOT",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListUserTest(),
		Name: "测试账号",
		Key:  "USER_TEST",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListPurpose(),
		Name: "目的",
		Key:  "PURPOSE",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListAuthCodeStatus(),
		Name: "验证码状态",
		Key:  "AUTH_CODE_STATUS",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListRuleType(),
		Name: "配置规则类型",
		Key:  "RULE_TYPE",
	})

	//
	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListSmsChannel(),
		Name: "短信渠道",
		Key:  "SMS_CHANNEL",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListServerPlatform(),
		Name: "服务器平台",
		Key:  "SERVER_PLATFORM",
	})

	ConstList = append(ConstList, ConstInfo{
		List: model.GetConstListCicdPublishStatus(),
		Name: "CICD发布状态",
		Key:  "CICD_PUBLISH_STATUS",
	})

	ConstList = append(ConstList, ConstInfo{
		List: GetConstListMailBoxType(),
		Name: "站内信,信件箱类型",
		Key:  "MAIL_BOX",
	})

	ConstList = append(ConstList, ConstInfo{
		List: GetConstListMailPeople(),
		Name: "站内信,接收人群类型",
		Key:  "MAIL_PEOPLE",
	})


	ConstList = append(ConstList, ConstInfo{
		List: GetConstListConfigPersistenceType(),
		Name: "配置中心持久化类型",
		Key:  "CONFIG_PERSISTENCE_TYPE",
	})









	//ConstList = append(ConstList, ConstInfo{
	//	List: model.GetConstListProjectType(),
	//	Name: "项目类型",
	//	Key:  "PROJECT_TYPE",
	//})

	return ConstList
}

//
////批量获取用户信息
//func GetUserInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
//	limit := info.PageSize
//	offset := info.PageSize * (info.Page - 1)
//	db := user.Gorm.Model(&model.User{})
//	var userList []model.User
//	err = db.Count(&total).Error
//	err = db.Limit(limit).Offset(offset).Preload("Authority").Find(&userList).Error
//	return err, userList, total
//}

////后台使用，权限控制
//func SetUserAuthority(uuid uuid.UUID, authorityId string) (err error) {
//	err = user.Gorm.Where("uuid = ?", uuid).First(&model.User{}).Update("authority_id", authorityId).Error
//	return err
//}
//
//func CheckUserIsCpByUserId(userId int) (res bool) {
//	_, user := FindUserById(userId)
//	auid, _ := strconv.Atoi(user.AuthorityId)
//	if auid == 9528 {
//		return true
//	}
//	return false
//}
//
//func CheckIsSuperAdmin(userId int) (res bool) {
//	_, user := FindUserById(userId)
//	auid, _ := strconv.Atoi(user.AuthorityId)
//	if auid == 888 {
//		return true
//	}
//	return false
//}
