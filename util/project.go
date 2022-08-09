package util

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

const (
	PROJECT_TYPE_SERVICE = 1
	PROJECT_TYPE_FE      = 2
	PROJECT_TYPE_APP     = 4
	PROJECT_TYPE_BE      = 3

	PROJECT_STATUS_OPEN  = 1
	PROJECT_STATUS_CLOSE = 2
)

//此函数在model const 里重新定义了
//func GetConstListProjectType() map[string]int {
//	list := make(map[string]int)
//	list["微服务"] = PROJECT_TYPE_SERVICE
//	list["前端"] = PROJECT_TYPE_FE
//	list["APP"] = PROJECT_TYPE_APP
//	list["后端"] = PROJECT_TYPE_BE
//
//	return list
//}

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	//Key 	string		`json:"key"`
	Type      int    `json:"type"`
	SecretKey string `json:"secretKey"`
	Status    int    `json:"status"`
	Git       string `json:"git"`
	Access    string `json:"access"`
}

type ProjectManager struct {
	Pool map[int]Project
	Gorm *gorm.DB
}

func NewProjectManager(gorm *gorm.DB) (*ProjectManager, error) {
	projectManager := new(ProjectManager)
	projectManager.Pool = make(map[int]Project)
	projectManager.Gorm = gorm

	err := projectManager.initAppPool()

	return projectManager, err
}

//初始化，会从MYSQL中读取数据，而没有监听MYSQL数据的变化，可重新再加载一次
func (projectManager *ProjectManager) DataReload() {
	projectManager.initAppPool()
}

func (projectManager *ProjectManager) initAppPool() error {
	//appManager.GetTestData()
	return projectManager.GetDataFromDb()
}

func (projectManager *ProjectManager) cleanPool() {
	projectManager.Pool = make(map[int]Project)
}

func (projectManager *ProjectManager) GetDataFromDb() error {
	projectManager.cleanPool() //清空原数据，重新从DB中读取
	db := projectManager.Gorm.Model(&model.Project{})
	var projectList []model.Project
	err := db.Where(" status = ? ", PROJECT_STATUS_OPEN).Find(&projectList).Error
	if err != nil {
		return err
	}
	if len(projectList) == 0 {
		return errors.New("app list empty!!!")
	}

	for _, v := range projectList {
		n := Project{
			Id:     v.Id,
			Status: v.Status,
			Name:   v.Name,
			Desc:   v.Desc,
			//Key		: v.Key,
			Type:      v.Type,
			Git:       v.Git,
			SecretKey: v.SecretKey,
			Access:    v.Access,
		}
		projectManager.AddOne(n)
	}
	return nil
}

func (projectManager *ProjectManager) AddOne(project Project) {
	projectManager.Pool[project.Id] = project
}

func (projectManager *ProjectManager) GetById(id int) (Project, bool) {
	one, ok := projectManager.Pool[id]
	//MyPrint("GetById : ",one ,ok)
	if ok {
		return one, false
	}
	return one, true
}

func (projectManager *ProjectManager) GetByName(name string) (project Project, empty bool) {
	for _, v := range projectManager.Pool {
		if v.Name == name {
			return v, false
		}
	}
	return project, true
}
