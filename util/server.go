package util

import (
	"errors"
	"gorm.io/gorm"
	"zgoframe/model"
)

type Server struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Platform         string `json:"platform"`
	OutIp            string `json:"out_ip"`
	InnerIp          string `json:"inner_ip"`
	Env              int    `json:"env"`
	Status           int    `json:"status"`      //1正常2关闭
	PingStatus       int    `json:"ping_status"` //1正常2异常
	SuperVisorStatus int    `json:"super_visor_status"`
}

type ServerManager struct {
	Pool map[int]Server
	Gorm *gorm.DB
}

func NewServerManger(gorm *gorm.DB) (*ServerManager, error) {
	serverManager := new(ServerManager)
	serverManager.Pool = make(map[int]Server)
	serverManager.Gorm = gorm

	err := serverManager.initAppPool()

	return serverManager, err
}

func (serverManager *ServerManager) initAppPool() error {
	//appManager.GetTestData()
	return serverManager.GetFromDb()
}

func (serverManager *ServerManager) GetFromDb() error {
	db := serverManager.Gorm.Model(&model.Server{})
	var serverList []model.Server
	err := db.Where(" status = ?  ", SERVER_STATUS_NORMAL).Find(&serverList).Error
	if err != nil {
		return err
	}
	if len(serverList) == 0 {
		return errors.New("app list empty!!!")
	}

	for _, v := range serverList {
		n := Server{
			Id:       v.Id,
			Name:     v.Name,
			Platform: v.Platform,
			InnerIp:  v.InnerIp,
			OutIp:    v.OutIp,
			Env:      v.Env,
			Status:   v.Status,
		}
		serverManager.AddOne(n)
	}
	return nil
}

func (serverManager *ServerManager) AddOne(hostServer Server) {
	serverManager.Pool[hostServer.Id] = hostServer
}

func (serverManager *ServerManager) GetById(id int) (Server, bool) {
	one, ok := serverManager.Pool[id]
	if ok {
		return one, false
	}
	return one, true
}

//func  (serverManager *Server)GetTypeName(typeValue int)string{
//	v ,_ := PROJECT_TYPE_MAP[typeValue]
//		return v
//}
