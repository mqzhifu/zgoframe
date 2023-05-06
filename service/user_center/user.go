package user_center

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

//func (user *User) ProjectAutoCreateUserDbRecord() {
//	for _, v := range projectManager.Pool {
//		var user model.User
//		//projectManager.Gorm.Model(&model.User{}).Where("project_id = ?",v.Id).FirstOrCreate(&model.User{ProjectId:v.Id})
//		err := projectManager.Gorm.Where("username = ?   ", v.Name).First(&user).Error
//		if err == nil { //证明该用户记录已经存在，不需要再创建
//			continue
//		}
//		newUser := model.User{}
//	}
//}
