package util

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

const(
	PROJECT_TYPE_SERVICE = 1
	PROJECT_TYPE_FE = 2
	PROJECT_TYPE_APP = 3
	PROJECT_TYPE_BE = 4
)

type Project struct {
	Id 		int			`json:"id"`
	Name	string		`json:"name"`
	Desc 	string		`json:"desc"`
	//Key 	string		`json:"key"`
	Type 	int			`json:"type"`
	SecretKey string 	`json:"secretKey"`
	Status  int 		`json:"status"`
	Git 	string 		`json:"git"`
	Access 	string		`json:"access"`
}

var PROJECT_TYPE_MAP = map[int]string{
	PROJECT_TYPE_SERVICE: "service",
	PROJECT_TYPE_FE:      "frontend",
	PROJECT_TYPE_APP:     "app",
	PROJECT_TYPE_BE:      "backend",
}

type ProjectManager struct {
	Pool 	map[int]Project
	Gorm 	*gorm.DB
}

func NewProjectManager (gorm *gorm.DB)(*ProjectManager,error) {
	projectManager 		:= new(ProjectManager)
	projectManager.Pool = make(map[int]Project)
	projectManager.Gorm = gorm

	err := projectManager.initAppPool()

	return projectManager,err
}

func (projectManager *ProjectManager)initAppPool()error{
	//appManager.GetTestData()
	return projectManager.GetFromDb()
}

func (projectManager *ProjectManager)GetFromDb()error{
	db := projectManager.Gorm.Model(&model.Project{})
	var appList []model.Project
	err := db.Where(" status = ? ", 1).Find(&appList).Error
	if err != nil{
		return err
	}
	if len(appList) == 0{
		return errors.New("app list empty!!!")
	}

	for _,v:=range appList{
		n := Project{
			Id 		: v.Id,
			Status	: v.Status,
			Name	: v.Name,
			Desc	: v.Desc,
			//Key		: v.Key,
			Type	: v.Type,
			Git		: v.Git,
			SecretKey: v.SecretKey,
			Access: v.Access,
		}
		projectManager.AddOne(n)
	}
	return nil
}

func (projectManager *ProjectManager) AddOne(project Project){
	projectManager.Pool[project.Id] = project
}

func (projectManager *ProjectManager) GetById(id int)(Project,bool){
	one ,ok := projectManager.Pool[id]
	if ok {
		return one,false
	}
	return one,true
}

func (projectManager *ProjectManager) GetByName(name string)(project Project,empty bool){
	for _,v:= range projectManager.Pool{
		if v.Name == name{
			return v,false
		}
	}
	return project,true
}


func  (projectManager *ProjectManager)GetTypeName(typeValue int)string{
	v ,_ := PROJECT_TYPE_MAP[typeValue]
		return v
}

func BasePathPlusTypeStr(basePath string,typeStr string)string{
	return basePath + "/" + typeStr + "/"
}

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