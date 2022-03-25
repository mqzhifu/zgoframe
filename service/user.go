package service

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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
	Gorm *gorm.DB
}

func NewUser(gorm *gorm.DB) *User {
	user := new(User)
	user.Gorm = gorm
	return user
}

//注册，用户名/密码
func (user *User) RegisterByUsername(R request.Register, h request.Header) (err error, userInter model.User) {
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

//最终 - 注册
func (user *User) Register(formUser model.User, h request.Header, userRegInfo UserRegInfo) (err error, userInter model.User) {
	var userRegType int

	formUser.Status = model.USER_STATUS_NOMAL

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
			_, exist, _ := user.FindUserByMobile(formUser.Mobile)
			if exist {
				return errors.New("mobile 已注册:" + formUser.Username), userInter
			}
		} else if userRegType == model.USER_REG_TYPE_EMAIL {
			_, exist, _ := user.FindUserByEmail(formUser.Username)
			if exist {
				return errors.New("email 已注册:" + formUser.Username), userInter
			}
		} else {
			_, exist, _ := user.FindUserByUsername(formUser.Username)
			if exist {
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
	userInter, exist, _ := user.FindUserByMobile(mobile)
	if exist {
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
func (user *User) LoginThird(rLoginThird request.RLoginThird, h request.Header) (userInfo model.User, isNewReg bool, err error) {
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
		return userInfo, false, nil
	}
	if err != nil {
		return userInfo, false, err
	}
	return userInfo, true, nil
}

func (user *User) FindUserByEmail(email string) (userInfo model.User, empty bool, err error) {
	err = user.Gorm.Where("email = ? ", email).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userInfo, false, nil
	}
	if err != nil {
		return userInfo, false, err
	}
	return userInfo, true, nil
}

func (user *User) FindUserByMobile(mobile string) (userInfo model.User, empty bool, err error) {
	err = user.Gorm.Where("mobile = ? ", mobile).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userInfo, false, nil
	}
	if err != nil {
		return userInfo, false, err
	}
	return userInfo, true, nil
}

//根据uuid查找一个用户信息
func (user *User) FindUserByUuid(uuid string) (err error, userInfo *model.User) {
	if err = user.Gorm.Where("`uuid` = ?", uuid).First(&userInfo).Error; err != nil {
		return errors.New("用户不存在"), userInfo
	}
	return nil, userInfo
}

func (user *User) DeleteUser(id float64) (err error) {
	var userInfo model.User
	err = user.Gorm.Where("id = ?", id).Delete(&userInfo).Error
	return err
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
	var userInfo model.User
	userInfo.Id = uid
	userInfo.Mobile = mobile
	userInfo.Guest = model.USER_GUEST_FALSE

	err = user.Gorm.Updates(&userInfo).Error
	return err
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
