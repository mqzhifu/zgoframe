package util

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

const(
	//APP_TYPE_SERVICE = 1
	//APP_TYPE_FE = 2
	//APP_TYPE_APP = 3
	//APP_TYPE_BE = 4
)

type HostServer struct {
	Id 			int 	`json:"id"`
	Name        string	`json:"name"`
	Platform    string	`json:"platform"`
	OutIp 		string	`json:"out_ip"`
	InnerIp     string	`json:"inner_ip"`
	Env 		string	`json:"env"`
}

//var APP_TYPE_MAP = map[int]string{
//	APP_TYPE_SERVICE: "service",
//	APP_TYPE_FE:      "frontend",
//	APP_TYPE_APP:     "app",
//	APP_TYPE_BE:      "backend",
//}

type HostServerManager struct {
	Pool map[int]HostServer
	Gorm 	*gorm.DB
}

func NewHostServer (gorm *gorm.DB)(*HostServerManager,error) {
	hostServerManager 		:= new(HostServerManager)
	hostServerManager.Pool = make(map[int]HostServer)
	hostServerManager.Gorm = gorm

	err := hostServerManager.initAppPool()

	return hostServerManager,err
}

func (hostServerManager *HostServerManager)initAppPool()error{
	//appManager.GetTestData()
	return hostServerManager.GetFromDb()
}

func (hostServerManager *HostServerManager)GetFromDb()error{
	db := hostServerManager.Gorm.Model(&model.Host{})
	var hoseServerList []model.Host
	err := db.Where(" status = ?  ", 1).Find(&hoseServerList).Error
	if err != nil{
		return err
	}
	if len(hoseServerList) == 0{
		return errors.New("app list empty!!!")
	}

	for _,v:=range hoseServerList{
		n := HostServer{
			Id 		: v.Id,
			Name	: v.Name,
			Platform	: v.Platform,
			InnerIp		: v.InnerIp,
			OutIp	: v.OutIp,
			Env :v.Env,
		}
		hostServerManager.AddOne(n)
	}
	return nil
}

func (hostServerManager *HostServerManager) AddOne(hostServer HostServer ){
	hostServerManager.Pool[hostServer.Id] = hostServer
}

func (hostServerManager *HostServerManager) GetById(id int)(HostServer,bool){
	one ,ok := hostServerManager.Pool[id]
	if ok {
		return one,false
	}
	return one,true
}

func  (hostServer *HostServerManager)GetTypeName(typeValue int)string{
	v ,_ := APP_TYPE_MAP[typeValue]
		return v
}

//func BasePathPlusTypeStr(basePath string,typeStr string)string{
//	return basePath + "/" + typeStr + "/"
//}
//func (appManager *AppManager)GetTestData(){
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
//}