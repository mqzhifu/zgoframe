package service

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
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
func Register(u model.User,h request.Header) (err error, userInter model.User) {
	var user model.User
	if !errors.Is(global.V.Gorm.Where("username = ? and app_id = ?", u.Username,u.AppId).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册"), userInter
	}
	isEmail := util.CheckEmailRule(u.Username)
	isMobile := util.CheckMobileRule(u.Username)
	userRegType := model.USER_REG_TYPE_NAME
	if isEmail{
		userRegType = model.USER_REG_TYPE_EMAIL
	}else if isMobile{
		userRegType = model.USER_REG_TYPE_MOBILE
	}
	// 否则 附加uuid 密码md5简单加密 注册
	u.Password = util.MD5V([]byte(u.Password))
	u.Uuid = uuid.NewV4()
	err = global.V.Gorm.Create(&u).Error

	//util.MyPrint("u.id:",u.Id)

	if err == nil{
		userReg := model.UserReg{
			AppId: u.AppId,
			Uid: u.Id,
			Type: userRegType,
			Ip: h.Ip,
			AutoIp: h.AutoIp,
			AppVersion: h.AppVersion,
			SourceType: h.SourceType,
			Os: h.OS,
			OsVersion: h.OSVersion,
			Device: h.Device,
			DeviceVersion: h.DeviceVersion,
			Lat: h.Lat,
			Lon: h.Lon,
			DeviceId: h.DeviceId,
			Dpi: h.DPI,
			Channel: model.CHANNEL_DEFAULT,
		}
		util.PrintStruct(userReg,":")
		//fmt.Sprintf("aaaaf:%+v", &userReg)
		//util.MyPrint("userReg:",userReg)
		err = global.V.Gorm.Create(&userReg).Error
		if err != nil{
			global.V.Zap.Error("create user_Reg err:"+err.Error())
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
	u.Password = util.MD5V([]byte(u.Password))
	err = global.V.Gorm.Where("username = ? AND password = ? AND app_id = ? ", u.Username, u.Password,u.AppId).Preload("Authority").First(&user).Error
	return err, &user
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
	if err = global.V.Gorm.Where("`uuid` = ?", uuid).First(&u).Error; err != nil{
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}

func CheckUserIsCpByUserId(userId int)(res bool){
	_,user := FindUserById(userId)
	auid,_ := strconv.Atoi(user.AuthorityId)
	if auid == 9528 {
		return true
	}
	return false
}


func CheckIsSuperAdmin(userId int)(res bool){
	_,user := FindUserById(userId)
	auid,_ := strconv.Atoi(user.AuthorityId)
	if auid == 888 {
		return true
	}
	return false
}