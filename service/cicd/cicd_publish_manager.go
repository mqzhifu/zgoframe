package cicd

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
	"zgoframe/util"
)

//部署项目时，触发按钮及一次部署的数据记录

type CICDPublicManager struct {
	Db *gorm.DB
}

func NewCICDPublicManager(gorm *gorm.DB) *CICDPublicManager {
	cICDPublicManager := new(CICDPublicManager)
	cICDPublicManager.Db = gorm
	return cICDPublicManager
}

func (CICDPublicManager *CICDPublicManager) InsertOne(service util.Service, server util.Server) model.CicdPublish {
	serviceInfo, _ := json.Marshal(service)
	serverInfo, _ := json.Marshal(server)
	data := model.CicdPublish{
		Status:      model.CICD_PUBLISH_STATUS_WAIT_DEPLOY,
		DeployStatus: model.CICD_PUBLISH_DEPLOY_STATUS_ING,
		ServiceId:   service.Id,
		ServerId:    server.Id,
		ServiceInfo: string(serviceInfo),
		ServerInfo:  string(serverInfo),
	}
	CICDPublicManager.Db.Create(&data)
	return data
}

func (CICDPublicManager *CICDPublicManager) UpDeployStatus(m model.CicdPublish, status int) {
	util.MyPrint("CICDPublicManager DeployStatus publishId:", m.Id, " new status:"+strconv.Itoa(status))
	CICDPublicManager.Db.Model(&m).Update("deploy_status", status)
	//db.Model(&Food{}).Update("price", 25)
}

func (CICDPublicManager *CICDPublicManager) UpStatus(m model.CicdPublish, status int) {
	util.MyPrint("CICDPublicManager UpStatus publishId:", m.Id, " new status:"+strconv.Itoa(status))
	CICDPublicManager.Db.Model(&m).Update("status", status)
	//db.Model(&Food{}).Update("price", 25)
}

func (CICDPublicManager *CICDPublicManager) UpInfo(m model.CicdPublish ) {
	util.MyPrint("CICDPublicManager UpInfo publishId:", m.Id)
	err := CICDPublicManager.Db.Updates(m).Error
	if err != nil{
		util.MyPrint("UpInfo err:"+err.Error())
	}
	//db.Model(&Food{}).Update("price", 25)
}

func (CICDPublicManager *CICDPublicManager) GetById(id int)(model.CicdPublish , error) {
	var m model.CicdPublish
	err  := CICDPublicManager.Db.First(&m,id).Error
	return m,err
}

func (CICDPublicManager *CICDPublicManager) GetList() ([]model.CicdPublish, error) {
	db := CICDPublicManager.Db.Model(&model.CicdPublish{})
	var cicdPublishList []model.CicdPublish
	err := db.Limit(10).Order("id desc").Find(&cicdPublishList).Error
	if err != nil {
		return cicdPublishList, err
	}
	if len(cicdPublishList) == 0 {
		return cicdPublishList, errors.New("app list empty!!!")
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
	return cicdPublishList, nil
}
