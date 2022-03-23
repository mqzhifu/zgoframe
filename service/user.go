package service

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"zgoframe/core/global"
	"zgoframe/http/request"
	"zgoframe/model"
	"zgoframe/util"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册-仅限：邮件、用户名、手机
//@param: u model.User
//@return: err error, userInter model.User
func Register(R request.Register, h request.Header) (err error, userInter model.User) {
	var user model.User
	var userRegType int

	u := model.User{
		Username:  R.Username,
		NickName:  R.NickName,
		Password:  R.Password,
		HeaderImg: R.HeaderImg,
		ProjectId: R.ProjectId,
		Sex:       R.Sex,
		Recommend: R.Recommend,
		Guest:     R.Guest,
		ThirdId:   R.ThirdId,
		Robot:     model.USER_ROBOT_FALSE,
		Status:    model.USER_STATUS_NOMAL,
	}

	if u.Guest != model.USER_GUEST_TRUE && u.Guest != model.USER_GUEST_FALSE {
		return errors.New("Guest value err."), userInter
	}

	if u.Guest == model.USER_REG_TYPE_GUEST {
		//deviceId = username
		if u.Username == "" {
			u.Username = MakeGuestUsername()
		}

		if !errors.Is(global.V.Gorm.Where("username = ? ", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
			return errors.New("用户名已注册"), userInter
		}
		userRegType = model.USER_REG_TYPE_NAME
	} else if R.ThirdType > 0 && R.ThirdId != "" {
		userRegType = model.USER_REG_TYPE_THIRD

		if !errors.Is(global.V.Gorm.Where("third_id = ?  ", u.ThirdId).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
			return errors.New("用户名已注册"), userInter
		}
	} else {
		userRegType = TurnRegByUsername(u.Username)
		if userRegType == model.USER_REG_TYPE_MOBILE {
			if !errors.Is(global.V.Gorm.Where("username = ? ", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
				return errors.New("mobile 已注册"), userInter
			}
		} else if userRegType == model.USER_REG_TYPE_EMAIL {
			if !errors.Is(global.V.Gorm.Where("email = ? ", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
				return errors.New("email 已注册"), userInter
			}
		} else {
			if !errors.Is(global.V.Gorm.Where("username = ? ", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
				return errors.New("用户名已注册"), userInter
			}
		}

		userRegType = TurnRegByUsername(u.Username)
	}

	if u.NickName == "" {
		u.NickName = MakeNickname()
	}

	if u.Password != "" {
		u.Password = util.MD5V([]byte(u.Password))
	}

	u.Uuid = uuid.NewV4().String()

	formatHeader := fmt.Sprintf("%+v", u)
	util.MyPrint("create user Info", formatHeader)

	err = global.V.Gorm.Create(&u).Error
	if err != nil {
		return err, userInter
	}
	//util.MyPrint("u.id:",u.Id)
	channel := model.CHANNEL_DEFAULT
	if R.Channel > 0 {
		channel = R.Channel
	}
	if err == nil {
		userReg := model.UserReg{
			ProjectId: u.ProjectId,
			Uid:       u.Id,
			ThirdType: R.ThirdType,
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
		err = global.V.Gorm.Create(&userReg).Error
		if err != nil {
			return errors.New("create user_Reg err:" + err.Error()), userInter
		}
	}
	return err, u
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Login
//@description: 用户登录
//@param: u *model.User
//@return: err error, userInter *model.User

func Login(u *model.User) (err error, userInter *model.User) {
	var user model.User
	if u.Username == "" {
		return errors.New("username empty"), userInter
	}

	if u.Password == "" {
		return errors.New("password empty"), userInter
	}

	u.Password = util.MD5V([]byte(u.Password))
	regType := TurnRegByUsername(u.Username)
	if regType == model.USER_REG_TYPE_MOBILE {
		err = global.V.Gorm.Where("mobile = ? AND password = ?   ", u.Mobile, u.Password).First(&user).Error
	} else if regType == model.USER_REG_TYPE_EMAIL {
		err = global.V.Gorm.Where("email = ? AND password = ?   ", u.Email, u.Password).First(&user).Error
	} else {
		err = global.V.Gorm.Where("username = ? AND password = ?   ", u.Username, u.Password).First(&user).Error
	}

	if err == nil {
		if user.Status != model.USER_STATUS_NOMAL {
			return errors.New("status err"), &user
		}
	}

	return err, &user
}

func LoginThird(user *model.User) (err error, userInter *model.User) {
	if errors.Is(global.V.Gorm.Where("third_id = ?  ", user.ThirdId).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册"), userInter
	}
	return nil, user
}

func MakeGuestUsername() string {
	return uuid.NewV4().String()
}

func MakeNickname() string {
	return uuid.NewV4().String()
}

func TurnRegByUsername(username string) int {
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

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.User, newPassword string
//@return: err error, userInter *model.User

func ChangePassword(u *model.User, newPassword string) (err error, userInter *model.User) {
	var user model.User
	u.Password = util.MD5V([]byte(u.Password))
	err = global.V.Gorm.Where("username = ? AND password = ?", u.Username, u.Password).First(&user).Update("password", util.MD5V([]byte(newPassword))).Error
	return err, u
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func GetUserInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.V.Gorm.Model(&model.User{})
	var userList []model.User
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Preload("Authority").Find(&userList).Error
	return err, userList, total
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func SetUserAuthority(uuid uuid.UUID, authorityId string) (err error) {
	err = global.V.Gorm.Where("uuid = ?", uuid).First(&model.User{}).Update("authority_id", authorityId).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func DeleteUser(id float64) (err error) {
	var user model.User
	err = global.V.Gorm.Where("id = ?", id).Delete(&user).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.User
//@return: err error, user model.User

func SetUserInfo(reqUser model.User) (err error, user model.User) {
	err = global.V.Gorm.Updates(&reqUser).Error
	return err, reqUser
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.User

func FindUserById(id int) (err error, user *model.User) {
	var u model.User
	err = global.V.Gorm.Where("`id` = ?", id).First(&u).Error
	return err, &u
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.User

func FindUserByUuid(uuid string) (err error, user *model.User) {
	var u model.User
	if err = global.V.Gorm.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}

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
