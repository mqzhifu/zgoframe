package model

// @description 用户基础信息
type User struct {
	MODEL
	Uuid      string `json:"uuid" db:"define:varchar(50);comment:UID字条串化;unique:uuid;index:uuid;defaultValue:''"` //UID字条串化
	ProjectId int    `json:"project_id" db:"define:tinyint(1);comment:项目ID;defaultValue:0"`                       //项目ID
	Sex       int    `json:"sex" db:"define:tinyint(1);comment:性别1男2女;defaultValue:0"`                            //性别1男2女
	Birthday  int    `json:"birthday" db:"define:int;comment:出生日期,unix时间戳;defaultValue:0"`                        //出生日期,unix时间戳
	Username  string `json:"username" db:"define:varchar(50);comment:用户登录名;defaultValue:''"`                      //用户登录名
	Password  string `json:"-" db:"define:varchar(50);comment:用户登录密码;defaultValue:''"`                            //登陆密码
	PayPs     string `json:"-" db:"define:varchar(50);comment:用户支付密码;defaultValue:''"`                            //支付密码
	NickName  string `json:"nick_name" db:"define:varchar(50);comment:用户昵称;defaultValue:''" `                     //昵称
	Mobile    string `json:"mobile" db:"define:varchar(50);comment:手机号;defaultValue:''"`                          //手机号
	Email     string `json:"email" db:"define:varchar(50);comment:邮箱;defaultValue:''"`                            //邮箱
	Robot     int    `json:"robot" db:"define:tinyint(1);comment:机器人;defaultValue:0"`                             //机器人
	Status    int    `json:"status" db:"define:tinyint(1);comment:状态1正常2禁用;defaultValue:0"`                       //状态1正常2禁用
	Guest     int    `json:"type" db:"define:tinyint(1);comment:是否游客,1是2否;defaultValue:0"`                        //是否游客,1是2否
	Test      int    `json:"test" db:"define:tinyint(1);comment:是否测试,1是2否;defaultValue:0"`                        //是否测试,1是2否
	Recommend string `json:"recommend" db:"define:varchar(50);comment:推荐人;defaultValue:''"`                       //推荐人
	HeaderImg string `json:"headerImg" gorm:"" db:"define:varchar(50);comment:头像url地址;defaultValue:''"`           //头像url地址
	//ThirdId   string `json:"third_id" db:"define:varchar(50);comment:三方平台(登陆)用户ID;defaultValue:''"`
	//AuthorityId string       `json:"authority_id" db:"define:varchar(50);comment:用户角色ID(后台使用);defaultValue:''"`
	//Authority   SysAuthority `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;" db:"define:varchar(50);comment:用户角色(后台使用);defaultValue:''"`
}

func (user *User) Db() {

}

func (user *User) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "用户表"

	return m
}

func (user *User) Count() {

}

////
////根据主键ID查找一条记录
//func (user *User) GetRowById(Id int)(*User, error){
//	userRow := &User{}
//	//这里没必要用first last ，因为它两会排序，既然ID是主键不可能出现重复|多条记录，也就不需要排序
//	txDb := Db.Take( userRow,Id)
//	//gorm.ErrRecordNotFound
//	return userRow, txDb.Error
//	//err := global.V.Gorm.Where("id = ?", Id).Take(userRow).Error
//	//return userRow, err
//}
////根据主键ID查找一组记录
//func (user *User) GetRowByIds(Ids []int)([]*User, error){
//	var users []*User
//	txDb := Db.Find( &users,Ids)
//	return users, txDb.Error
//}
//
//
//func (user *User) GetRow(where string)(*User, error){
//	userRow := &User{}
//	err := Db.Where(where).First(userRow).Error
//	return userRow, err
//}
//
//type GetListPara struct {
//	Where	string
//	Fields 	string
//	Order 	string
//	Group   stringΩ
//	Limit 	int
//	Offset	int
//}
//
//func (user *User) GetList(para GetListPara )([]*User, error){
//	var users []*User
//	db := global.V.Gorm.Find(para.Where)
//	if para.Fields != ""{
//		db.Select(para.Fields)
//	}
//
//	if para.Group != ""{
//		db.Order(para.Group)
//	}
//
//	if para.Order != ""{
//		db.Order(para.Order)
//	}
//
//	if para.Limit >= 0 {
//		db.Limit(para.Offset)
//	}
//
//	if para.Limit  >= 0{
//		db.Limit(para.Limit )
//	}
//
//	err := db.Find(users).Error
//	return users, err
//}
//
//func (user *User) UpRowById(Id int)(*User, error){
//
//}
//
//func (user *User) UpRowByIds(Id int)(*User, error){
//
//}
//
//func (user *User) UpRow(){
//
//}
//
//func (user *User) UpRows(){
//
//}
//
//func (user *User) AddOne(){
//
//}
//
