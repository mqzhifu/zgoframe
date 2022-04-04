package util

import (
	"context"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"zgoframe/model"
)

//一个服务下面的一个节点,给服务发现使用
type ServiceNode struct {
	ProjectId   int         `json:"project_id"`
	ServiceId   int         `json:"service_id"`
	ServiceName string      `json:"service_name"`
	ListenIp    string      `json:"listen_ip"` //监听IP
	Ip          string      `json:"ip"`        //对外显示IP
	Port        string      `json:"port"`
	Protocol    int         `json:"protocol"` //暂未实现
	Desc        string      `json:"desc"`
	Status      int         `json:"status"`
	IsSelfReg   bool        `json:"is_self_reg"` //是否为当前服务自己注册的服务
	RegTime     int         `json:"reg_time"`    //注册时间
	DBKey       string      `json:"db_key"`
	CreateTime  int         `json:"create_time"`
	Log         *zap.Logger `json:"-"`
}

//服务，这里两个地方用，
//1. 从DB里读出来，后台录入的情况
//2. 服务发现也会创建这个节点体
type Service struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	//Key 		string			`json:"key"`
	DBKey      string         `json:"db_key"`
	Status     int            `json:"status"`
	Desc       string         `json:"desc"`
	Type       int            `json:"type"`
	CreateTime int            `json:"create_time"`
	SecretKey  string         `json:"secret_key"`
	Git        string         `json:"git"`
	LBType     int            `json:"lb_type"`
	Deploy     int            `json:"deploy"`
	List       []*ServiceNode `json:"-"`

	Log *zap.Logger `json:"-"`

	watchCancel context.CancelFunc `json:"-"` //监听 分布式DB 取消函数

	LeaseGrantId   clientv3.LeaseID `json:"-"` //etcd
	Lease          clientv3.Lease   `json:"-"` //etcd
	LeaseCancelCtx context.Context  `json:"-"` //etcd
}

func (service *Service) ToString() string {
	str := "id:" + strconv.Itoa(service.Id) + " name:" + service.Name + " DBKey:" + service.DBKey + " CreateTime:" + strconv.Itoa(service.CreateTime) + " LBType:" + strconv.Itoa(service.LBType)
	return str
}

//给一个新的节点添加到一个服务中
func (service *Service) AddServiceList(node *ServiceNode) error {
	service.Log.Info("insert node to serviceList:" + node.ToString())
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

	service.List = append(service.List, node)
	return nil
	//serviceManager.list[service.Name] = service
}
func (serviceNode *ServiceNode) GetDns() string {
	return serviceNode.Ip + ":" + serviceNode.Port
}

func (serviceNode *ServiceNode) ToString() string {
	isSelf := "false"
	if serviceNode.IsSelfReg {
		isSelf = "true"
	}
	//ServiceId 	int 			`json:"service_id"`
	//RegTime 	int 			`json:"reg_time"`		//注册时间
	debugInfo := " ServiceName:" + serviceNode.ServiceName + " dns : " + serviceNode.GetDns() + ", protocol:" + strconv.Itoa(serviceNode.Protocol) + ", isSelf:" + isSelf + " ,ListenIp:" + serviceNode.ListenIp + " DBKey:" + serviceNode.DBKey + " createTime:" + strconv.Itoa(serviceNode.CreateTime)
	return debugInfo
}

//创建一个新的服务的节点
func (service *Service) NewServiceNode(node ServiceNode) error {
	prefix := "NewServiceNode"
	msgPrefix := prefix + " err: info check err:"
	if !CheckServiceProtocolExist(node.Protocol) {
		errMsg := msgPrefix + " err: protocol empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	if node.Ip == "" {
		errMsg := msgPrefix + " err:ip empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	if node.Port == "" {
		errMsg := msgPrefix + " err:port empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	//if node.DBKey == ""{
	//	errMsg := msgPrefix + " DBKey empty"
	//	MyPrint(errMsg)
	//	return errors.New(errMsg)
	//}

	node.CreateTime = GetNowTimeSecondToInt()

	if node.ServiceName == "" {
		errMsg := msgPrefix + " err:serviceName empty"
		//MyPrint(errMsg)
		return errors.New(errMsg)
	}

	err := service.AddServiceList(&node)
	return err
}

//===========================
type ServiceManager struct {
	Pool map[int]Service
	Gorm *gorm.DB
}

func NewServiceManager(gorm *gorm.DB) (*ServiceManager, error) {
	serviceManager := new(ServiceManager)
	serviceManager.Pool = make(map[int]Service)
	serviceManager.Gorm = gorm
	err := serviceManager.initAppPool()

	return serviceManager, err
}

func (serviceManager *ServiceManager) initAppPool() error {
	//appManager.GetTestData()
	return serviceManager.GetFromDb()
}

func (serviceManager *ServiceManager) GetFromDb() error {
	db := serviceManager.Gorm.Model(&model.Project{})
	var serviceList []model.Service
	err := db.Where(" status = ? and type = ? ", 1, model.PROJECT_TYPE_SERVICE).Find(&serviceList).Error
	if err != nil {
		return err
	}
	if len(serviceList) == 0 {
		return errors.New("app list empty!!!")
	}

	for _, v := range serviceList {
		//fmt.Println(v)
		n := Service{
			Id:     int(v.Id),
			Status: v.Status,
			Name:   v.Name,
			Desc:   v.Desc,
			//Key: v.Key,
			Type:      v.Type,
			SecretKey: v.SecretKey,
			Git:       v.Git,
		}
		serviceManager.AddOne(n)
	}
	return nil
}

func (serviceManager *ServiceManager) AddOne(app Service) {
	serviceManager.Pool[app.Id] = app
}

func (serviceManager *ServiceManager) GetById(id int) (Service, bool) {
	one, ok := serviceManager.Pool[id]
	if ok {
		return one, false
	}
	return one, true
}

func (serviceManager *ServiceManager) GetByName(name string) (service Service, isEmpty bool) {
	//fmt.Println("serviceManager GetByName:"+name)
	if len(serviceManager.Pool) <= 0 {
		return service, isEmpty
	}

	for _, v := range serviceManager.Pool {
		if v.Name == name {
			return v, false
		}
	}

	return service, true
}

//func  (serviceManager *ServiceManager)GetTypeName(typeValue int)string{
//	v ,_ := PROJECT_TYPE_MAP[typeValue]
//	return v
//}
