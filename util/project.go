package util

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

//type Project struct {
//	Id        int    `json:"id"`
//	Name      string `json:"name"`
//	Desc      string `json:"desc"`
//	Type      int    `json:"type"`
//	SecretKey string `json:"secretKey"`
//	Status    int    `json:"status"`
//	Git       string `json:"git"`
//	Access    string `json:"access"`
//}

type ProjectManager struct {
	Pool map[int]model.Project
	Gorm *gorm.DB
}

func NewProjectManager(gorm *gorm.DB) (*ProjectManager, error) {
	projectManager := new(ProjectManager)
	projectManager.Pool = make(map[int]model.Project)
	projectManager.Gorm = gorm

	err := projectManager.initAppPool()

	return projectManager, err
}

// 初始化，会从MYSQL中读取数据，而没有监听MYSQL数据的变化，可重新再加载一次
func (projectManager *ProjectManager) DataReload() error {
	return projectManager.initAppPool()
}

func (projectManager *ProjectManager) initAppPool() error {
	return projectManager.GetDataFromDb()
}

// 清空 - 内存缓存
func (projectManager *ProjectManager) cleanPool() {
	projectManager.Pool = make(map[int]model.Project)
}

func (projectManager *ProjectManager) GetDataFromDb() error {
	projectManager.cleanPool() //清空原数据，重新从DB中读取
	db := projectManager.Gorm.Model(&model.Project{})
	var projectList []model.Project
	//读取DB中，所有状态为正常的 project 记录
	err := db.Where(" status = ? ", model.PROJECT_STATUS_OPEN).Find(&projectList).Error
	if err != nil {
		return err
	}
	if len(projectList) == 0 {
		return errors.New("app list empty!!!")
	}
	//将DB数据  追回到  内存中
	for _, v := range projectList {
		projectManager.AddOne(v)
	}
	return nil
}

func (projectManager *ProjectManager) AddOne(project model.Project) {
	projectManager.Pool[project.Id] = project
}

func (projectManager *ProjectManager) GetById(id int) (model.Project, bool) {
	one, ok := projectManager.Pool[id]
	if ok {
		return one, false
	}
	return one, true
}

func (projectManager *ProjectManager) GetByName(name string) (project model.Project, empty bool) {
	for _, v := range projectManager.Pool {
		if v.Name == name {
			return v, false
		}
	}
	return project, true
}
