package util

import (
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"errors"
)

//部署项目时，触发按钮及一次部署的数据记录

type CICDPublicManager struct{
	Db *gorm.DB
}

func NewCICDPublicManager(gorm *gorm.DB)*CICDPublicManager{
	cICDPublicManager := new(CICDPublicManager)
	cICDPublicManager.Db = gorm
	return cICDPublicManager
}

func(CICDPublicManager *CICDPublicManager) InsertOne(service Service,server Server)model.CICDPublish{
	serviceInfo ,_ := json.Marshal(service)
	serverInfo ,_ := json.Marshal(server)
	data := model.CICDPublish{
		Status: 1,
		ServiceId: service.Id,
		ServerId: server.Id,
		ServiceInfo: string(serviceInfo),
		ServerInfo: string(serverInfo),
	}
	CICDPublicManager.Db.Create(&data)
	return data
}

func (CICDPublicManager *CICDPublicManager) UpStatus(m model.CICDPublish,status int){
	MyPrint("CICDPublicManager UpStatus publishId:",m.Id , " new status:"+strconv.Itoa(status))
	CICDPublicManager.Db.Model(&m).Update("status",status)
	//db.Model(&Food{}).Update("price", 25)
}

func (CICDPublicManager *CICDPublicManager) GetList()([]model.CICDPublish,error){
	db := CICDPublicManager.Db.Model(&model.CICDPublish{})
	var cicdPublishList []model.CICDPublish
	err := db.Limit(10).Order("id desc").Find(&cicdPublishList).Error
	if err != nil{
		return cicdPublishList,err
	}
	if len(cicdPublishList) == 0{
		return cicdPublishList,errors.New("app list empty!!!")
	}

	//for _,v:=range cicdPublishList{
	//	n := model.CICDPublish{
	//		ServiceId:v.ServiceId,
	//		HostId: v.HostId,
	//		Status: v.Status,
	//
	//		ServiceInfo: v.ServiceInfo,
	//		ServerInfo: v.ServerInfo,
	//	}
	//}
	return cicdPublishList,nil
}