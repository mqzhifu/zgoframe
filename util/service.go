package util

import (
	"errors"
	"gorm.io/gorm"
	"zgoframe/model"
	"context"
	"go.etcd.io/etcd/clientv3"
)
//一个服务下面的一个节点,给服务发现使用
type ServiceNode struct {
	ServiceId 	int 			`json:"service_id"`
	ServiceName string			`json:"service_name"`
	ListenIp 	string			`json:"listen_ip"`
	Ip			string			`json:"ip"`
	Port		string			`json:"port"`
	Protocol 	int				`json:"protocol"`		//暂未实现
	Desc 		string			`json:"desc"`
	Status 		int				`json:"status"`
	IsSelfReg 	bool			`json:"is_self_reg"`	//是否为当前服务自己注册的服务
	RegTime 	int 			`json:"reg_time"`		//注册时间
	DBKey		string 			`json:"db_key"`
}
//服务，这里两个地方用，
//1. 从DB里读出来，后台录入的情况
//2. 服务发现也会创建这个节点体
type Service struct {
	Id 			int				`json:"id"`
	Name 		string			`json:"name"`
	Key 		string
	DBKey 		string
	Status 		int
	Desc 		string
	Type 		int
	CreateTime 	int
	SecretKey	string
	Git 		string
	LBType 		int
	List 		[]*ServiceNode	`json:"list"`

	watchCancel	context.CancelFunc	`json:"-"`			//监听 分布式DB 取消函数

	LeaseGrantId clientv3.LeaseID	`json:"-"`			//etcd
	Lease clientv3.Lease			`json:"-"`			//etcd
	LeaseCancelCtx context.Context	`json:"-"`			//etcd
}
//给一个新的节点添加到一个服务中
func (service *Service)AddServiceList(node *ServiceNode)error{
	MyPrint("insert node to serviceList.")
	//msgPrefix := "AddServiceList "
	//if node.ServiceName == ""{
	//	errMsg := msgPrefix + " serviceName empty"
	//	MyPrint(errMsg)
	//	return errors.New(errMsg)
	//}
	//
	//if node.Ip == ""{
	//	errMsg := msgPrefix + " IP empty"
	//	MyPrint(errMsg)
	//	return errors.New(errMsg)
	//}
	//
	//if node.Port == ""{
	//	errMsg := msgPrefix + " port empty"
	//	MyPrint(errMsg)
	//	return errors.New(errMsg)
	//}

	service.List = append(service.List,node)
	return nil
	//serviceManager.list[service.Name] = service
}
func (serviceNode *ServiceNode)GetDns()string{
	return serviceNode.Ip + ":" + serviceNode.Port
}
//创建一个新的服务的节点
func (service *Service)NewServiceNode(node ServiceNode)error{
	prefix := "NewServiceNode"
	msgPrefix := prefix + " err: info check err:"
	if !CheckServiceProtocolExist(node.Protocol){
		errMsg := msgPrefix + " err: protocol empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	if node.Ip == ""{
		errMsg := msgPrefix + " err:ip empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	if node.Port == ""{
		errMsg := msgPrefix + " err:port empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	//if node.DBKey == ""{
	//	errMsg := msgPrefix + " DBKey empty"
	//	MyPrint(errMsg)
	//	return errors.New(errMsg)
	//}

	if node.ServiceName == ""{
		errMsg := msgPrefix + " err:serviceName empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	MyPrint("NewServiceNode success:" , node)

	err := service.AddServiceList(&node)
	return err
}

//===========================
type ServiceManager struct {
	Pool map[int]Service
	Gorm 	*gorm.DB
}

func NewServiceManager (gorm *gorm.DB)(*ServiceManager,error) {
	serviceManager := new(ServiceManager)
	serviceManager.Pool = make(map[int]Service)
	serviceManager.Gorm = gorm
	err := serviceManager.initAppPool()

	return serviceManager,err
}

func (serviceManager *ServiceManager)initAppPool()error{
	//appManager.GetTestData()
	return serviceManager.GetFromDb()
}

func (serviceManager *ServiceManager)GetFromDb()error{
	db := serviceManager.Gorm.Model(&model.Service{})
	var serviceList []model.Service
	err := db.Where(" status = ? ", 1).Find(&serviceList).Error
	if err != nil{
		return err
	}
	if len(serviceList) == 0{
		return errors.New("app list empty!!!")
	}

	for _,v:=range serviceList{
		n := Service{
			Id : int(v.Id),
			Status: v.Status,
			Name: v.Name,
			Desc: v.Desc,
			Key: v.Key,
			Type: v.Type,
			SecretKey: v.SecretKey ,
			Git:v.Git,
		}
		serviceManager.AddOne(n)
	}
	return nil
}

func (serviceManager *ServiceManager) AddOne(app Service){
	serviceManager.Pool[app.Id] = app
}

func (serviceManager *ServiceManager) GetById(id int)(Service,bool){
	one ,ok := serviceManager.Pool[id]
	if ok {
		return one,false
	}
	return one,true
}

func (serviceManager *ServiceManager) GetByName(name string)(service Service,isEmpty bool){
	if len(serviceManager.Pool) <= 0{
		return service,isEmpty
	}

	for _,v:= range serviceManager.Pool{
		if v.Name == name{
			return v,false
		}
	}

	return service,true
}

func  (serviceManager *ServiceManager)GetTypeName(typeValue int)string{
	v ,_ := APP_TYPE_MAP[typeValue]
	return v
}
