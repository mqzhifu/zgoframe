package util

import (
	"errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

type Instance struct {
	Id     int
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	Env    int    `json:"env"`
	User   string `json:"user"`
	Ps     string `json:"ps"`
	Ext    string `json:"ext"`
	Status int    `json:"status"`
}

type InstanceManager struct {
	Pool map[int]Instance
	Gorm *gorm.DB
}

// 3方实例管理：mysql redis etcd
func NewInstanceManager(gorm *gorm.DB) (*InstanceManager, error) {
	instanceManager := new(InstanceManager)

	instanceManager.Pool = make(map[int]Instance)
	instanceManager.Gorm = gorm

	err := instanceManager.initInstancePool()

	return instanceManager, err
}

func (instanceManager *InstanceManager) initInstancePool() error {
	// appManager.GetTestData()
	return instanceManager.GetFromDb()
}

func (instanceManager *InstanceManager) GetFromDb() error {
	db := instanceManager.Gorm.Model(&model.Instance{})
	var instanceList []model.Instance
	err := db.Where(" status = ?  ", 1).Find(&instanceList).Error
	if err != nil {
		return err
	}
	if len(instanceList) == 0 {
		return errors.New("app list empty!!!")
	}

	for _, v := range instanceList {
		n := Instance{
			Id:   v.Id,
			Name: v.Name,
			Host: v.Host,
			Port: v.Port,
			Env:  v.Env,
			User: v.User,
			Ps:   v.Ps,
			Ext:  v.Ext,
		}
		instanceManager.AddOne(n)
	}
	return nil
}

func (instanceManager *InstanceManager) AddOne(instance Instance) {
	instanceManager.Pool[instance.Id] = instance
}

func (instanceManager *InstanceManager) GetById(id int) (Instance, bool) {
	one, ok := instanceManager.Pool[id]
	if ok {
		return one, false
	}
	return one, true
}

func (instanceManager *InstanceManager) GetByEnv(env int) (in []Instance) {

	for _, v := range instanceManager.Pool {
		if v.Env == env {
			in = append(in, v)
		}
	}
	return in
}

func (instanceManager *InstanceManager) GetByEnvName(env int, name string) (in Instance, empty bool) {
	// MyPrint("GetByEnvName:",env , " name:",name)

	for _, v := range instanceManager.Pool {
		if v.Env == env && v.Name == name {
			// MyPrint(v)
			return v, false
		}
	}
	MyPrint("GetByEnvName empty", " ", env, " ", name)
	return in, true
}

// func BasePathPlusTypeStr(basePath string,typeStr string)string{
//	return basePath + "/" + typeStr + "/"
// }
// func (appManager *AppManager)GetTestData(){
//	app := App{
//		Id:        1,
//		Name:      "gamematch",
//		Type:      APP_TYPE_SERVICE,
//		Desc:      "游戏匹配",
//		Key:       "gamematch",
//		SecretKey: "123456",
//	}
//	appManager.AddOne(app)
//	app = App{
//		Id:        2,
//		Name:      "frame_sync",
//		Type:      APP_TYPE_SERVICE,
//		Desc:      "帧同步",
//		Key:       "frame_sync",
//		SecretKey: "123456",
//	}
//	appManager.AddOne(app)
//	app = App{
//		Id:        3,
//		Name:      "logslave",
//		Type:      APP_TYPE_SERVICE,
//		Desc:      "日志收集器",
//		Key:       "logslave",
//		SecretKey: "123456",
//	}
//	appManager.AddOne(app)
//	app = App{
//		Id:        4,
//		Name:      "frame_sync_fe",
//		Type:      APP_TYPE_FE,
//		Desc:      "帧同步-前端",
//		Key:       "frame_sync_fe",
//		SecretKey: "123456",
//	}
//	appManager.AddOne(app)
//	app = App{
//		Id:        5,
//		Name:      "zgoframe",
//		Type:      APP_TYPE_SERVICE,
//		Desc:      "测试-框架端",
//		Key:       "test_frame",
//		SecretKey: "123456",
//	}
//	appManager.AddOne(app)
// }
