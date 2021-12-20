package util

import (
	"encoding/json"
	"gorm.io/gorm"
	"zgoframe/model"
)

type CICDPublicManager struct{
	Db *gorm.DB
}

func NewCICDPublicManager(gorm *gorm.DB)*CICDPublicManager{
	cICDPublicManager := new(CICDPublicManager)
	cICDPublicManager.Db = gorm
	return cICDPublicManager
}

func(CICDPublicManager *CICDPublicManager) InsertOne(service Service,server HostServer)model.CICDPublish{
	serviceInfo ,_ := json.Marshal(service)
	serverInfo ,_ := json.Marshal(server)
	data := model.CICDPublish{
		Status: 0,
		ServiceId: service.Id,
		HostId: server.Id,
		ServiceInfo: string(serviceInfo),
		ServerInfo: string(serverInfo),
	}
	CICDPublicManager.Db.Create(&data)
	return data
}
