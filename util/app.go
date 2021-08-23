package util

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

const(
	APP_TYPE_SERVICE = 1
	APP_TYPE_FE = 2
	APP_TYPE_APP = 3
	APP_TYPE_BE = 4
)

type App struct {
	Id 		int			`json:"id"`
	Name	string		`json:"name"`
	Desc 	string		`json:"desc"`
	Key 	string		`json:"key"`
	Type 	int			`json:"type"`
	SecretKey string 	`json:"secretKey"`
	Status  int 		`json:"status"`
}

var APP_TYPE_MAP = map[int]string{
	APP_TYPE_SERVICE: "service",
	APP_TYPE_FE:      "frontend",
	APP_TYPE_APP:     "app",
	APP_TYPE_BE:      "backend",
}

type AppManager struct {
	Pool map[int]App
	Gorm 	*gorm.DB
}

func NewAppManager (gorm *gorm.DB)(*AppManager,error) {
	appManager := new(AppManager)
	appManager.Pool = make(map[int]App)
	appManager.Gorm = gorm
	err := appManager.initAppPool()

	return appManager,err
}

func (appManager *AppManager)initAppPool()error{
	//app := App{
	//	Id:        1,
	//	Name:      "gamematch",
	//	Type:      APP_TYPE_SERVICE,
	//	Desc:      "游戏匹配",
	//	Key:       "gamematch",
	//	SecretKey: "123456",
	//}
	//appManager.AddOne(app)
	//app = App{
	//	Id:        2,
	//	Name:      "frame_sync",
	//	Type:      APP_TYPE_SERVICE,
	//	Desc:      "帧同步",
	//	Key:       "frame_sync",
	//	SecretKey: "123456",
	//}
	//appManager.AddOne(app)
	//app = App{
	//	Id:        3,
	//	Name:      "logslave",
	//	Type:      APP_TYPE_SERVICE,
	//	Desc:      "日志收集器",
	//	Key:       "logslave",
	//	SecretKey: "123456",
	//}
	//appManager.AddOne(app)
	//app = App{
	//	Id:        4,
	//	Name:      "frame_sync_fe",
	//	Type:      APP_TYPE_FE,
	//	Desc:      "帧同步-前端",
	//	Key:       "frame_sync_fe",
	//	SecretKey: "123456",
	//}
	//appManager.AddOne(app)
	//app = App{
	//	Id:        5,
	//	Name:      "zgoframe",
	//	Type:      APP_TYPE_SERVICE,
	//	Desc:      "测试-框架端",
	//	Key:       "test_frame",
	//	SecretKey: "123456",
	//}
	//appManager.AddOne(app)


	return appManager.GetFromDb()
}

func (appManager *AppManager)GetFromDb()error{
	db := appManager.Gorm.Model(&model.App{})
	var appList []model.App
	err := db.Where(" status = ? ", 1).Find(&appList).Error
	if err != nil{
		return err
	}
	if len(appList) == 0{
		return errors.New("app list empty!!!")
	}

	for _,v:=range appList{
		n := App{
			Id : int(v.Id),
			Status: v.Status,
			Name: v.Name,
			Desc: v.Desc,
			Key: v.Key,
			Type: v.Type,
			SecretKey: v.SecretKey,

		}
		appManager.AddOne(n)
	}
	return nil
}

func (appManager *AppManager) AddOne(app App){
	appManager.Pool[app.Id] = app
}

func (appManager *AppManager) GetById(id int)(App,bool){
	one ,ok := appManager.Pool[id]
	if ok {
		return one,false
	}
	return one,true
}

func  (appManager *AppManager)GetTypeName(typeValue int)string{
	v ,_ := APP_TYPE_MAP[typeValue]
		return v
}

func BasePathPlusTypeStr(basePath string,typeStr string)string{
	return basePath + "/" + typeStr + "/"
}