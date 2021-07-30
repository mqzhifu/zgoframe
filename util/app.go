package util

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
}

var APP_TYPE_MAP = map[int]string{
	APP_TYPE_SERVICE: "service",
	APP_TYPE_FE:      "frontend",
	APP_TYPE_APP:     "app",
	APP_TYPE_BE:      "backend",
}

type AppManager struct {
	Pool map[int]App
}

func NewAppManager ()*AppManager {
	appManager := new(AppManager)
	appManager.Pool = make(map[int]App)
	appManager.initAppPool()
	return appManager
}

func (appManager *AppManager)initAppPool(){
	app := App{
		Id:        1,
		Name:      "gamematch",
		Type:      APP_TYPE_SERVICE,
		Desc:      "游戏匹配",
		Key:       "gamematch",
		SecretKey: "123456",
	}
	appManager.AddOne(app)
	app = App{
		Id:        2,
		Name:      "frame_sync",
		Type:      APP_TYPE_SERVICE,
		Desc:      "帧同步",
		Key:       "frame_sync",
		SecretKey: "123456",
	}
	appManager.AddOne(app)
	app = App{
		Id:        3,
		Name:      "logslave",
		Type:      APP_TYPE_SERVICE,
		Desc:      "日志收集器",
		Key:       "logslave",
		SecretKey: "123456",
	}
	appManager.AddOne(app)
	app = App{
		Id:        4,
		Name:      "frame_sync_fe",
		Type:      APP_TYPE_FE,
		Desc:      "帧同步-前端",
		Key:       "frame_sync_fe",
		SecretKey: "123456",
	}
	app = App{
		Id:        5,
		Name:      "test_frame",
		Type:      APP_TYPE_FE,
		Desc:      "测试-前端框架端",
		Key:       "test_frame",
		SecretKey: "123456",
	}

	appManager.AddOne(app)
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