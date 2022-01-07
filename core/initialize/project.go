package initialize

import (
	"errors"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/util"
)

func InitProject()(err error){
	if global.C.System.ProjectId <=0  {
		return errors.New("ProjectId empty")
	}
	global.V.ProjectMng ,err  = util.NewProjectManager(global.V.Gorm)
	if err != nil{
		return err
	}
	empty := false
	global.V.Project,empty = global.V.ProjectMng.GetById(global.C.System.ProjectId)
	if empty {
		return errors.New("AppId not match : " + strconv.Itoa(global.C.System.ProjectId) )
	}
	return nil
}
//func InitAppService()(err error){
//	if global.C.System.AppId <=0 && global.C.System.ServiceId <=0 {
//		return errors.New("appId and serviceId both empty")
//	}
//	//一个项目要么是APP 要么是service
//	if global.C.System.AppId >0 && global.C.System.ServiceId >0 {
//		return errors.New("appId and serviceId both >= 0 ")
//	}
//
//	if global.C.System.AppId > 0{
//		global.V.AppMng ,err  = GetNewAppManager()
//		if err != nil{
//			util.MyPrint("GetNewApp err:",err)
//			return err
//		}
//		//根据APPId去DB中查找详细信息
//		app,empty := global.V.AppMng.GetById(global.C.System.AppId)
//		if empty {
//			return errors.New("AppId not match : " + strconv.Itoa(global.C.System.AppId) )
//		}
//		global.V.App = app
//		//util.MyPrint("project app info flow:")
//		//util.PrintStruct(app,":")
//	}else{
//		global.V.ServiceManager ,err  = GetNewServiceManager()
//		if err != nil{
//			util.MyPrint("GetNewServiceManager err:",err)
//			return err
//		}
//		//根据APPId去DB中查找详细信息
//		service,empty := global.V.ServiceManager.GetById(global.C.System.ServiceId)
//		if empty {
//			return errors.New("ServiceId not match : " + strconv.Itoa(global.C.System.ServiceId) )
//		}
//		global.V.Service = service
//		//util.MyPrint("service info flow:")
//		//util.PrintStruct(service,":")
//	}
//	return nil
//}
//
////初始化app管理容器
//func GetNewAppManager()(m *util.AppManager,e error){
//	appM,err := util.NewAppManager(global.V.Gorm)
//	if err != nil{
//		return m,err
//	}
//
//	return appM,nil
//}
//
//func GetNewServiceManager()(m *util.ServiceManager,e error){
//	sm,err := util.NewServiceManager(global.V.Gorm)
//	if err != nil{
//		return sm,err
//	}
//
//	return sm,nil
//}
