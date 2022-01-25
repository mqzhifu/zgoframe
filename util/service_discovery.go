package util

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

var CONSUL_NOT_OPEN = errors.New("consule not open.")


type ServiceChange struct{
	Name 	string
	Ip 		string
	Port 	string
	NewIp 	string
	NewPort string
	Action 	string
}

//服务发现-管理器
type ServiceDiscovery struct {
	//list->service.list->serviceNode
	list      		map[string]*Service      //所有-服务列表,ps:包括自己注册,包括client server
	option    		ServiceDiscoveryOption
	WatchMsg 		chan ServiceChange	//监听的数据发生变化时，发送通知给调用者,暂关闭
	CancelCxt 		context.Context
	CancelFunc 		context.CancelFunc
	OPLock 			sync.Mutex		// 注册|删除：互斥，锁的粒度有点大，后期优化
	//Status			string		//开启 关闭中  ，用于监听接收消息
	WatchClose		chan bool
	//consul 		int //待处理
}
//实例化 ServiceManager 参数
type ServiceDiscoveryOption struct {
	Etcd 			*MyEtcd
	Log 			*zap.Logger
	DiscoveryType 	int			//服务发现 存储类型
	Prefix			string		//服务发现 DB，存储的路径 前缀 ，PS:不要以 反斜杠(/) 结尾
	ServiceManager	*ServiceManager
	AutoCreateServiceLBType int	//自动创建一个服务时，默认的 负载均衡 策略
}
//创建一个 ServiceDiscovery
func NewServiceDiscovery(serviceDiscoveryOption ServiceDiscoveryOption)(sd *ServiceDiscovery,err error) {
	serviceDiscoveryOption.Log.Info("NewServiceDiscovery ， prefix:"+serviceDiscoveryOption.Prefix + " DiscoveryType:"+strconv.Itoa(serviceDiscoveryOption.DiscoveryType))

	serviceDiscovery := new(ServiceDiscovery)

	if ! CheckServiceDiscoveryExist(serviceDiscoveryOption.DiscoveryType){
		return sd,errors.New(" discoveryType err.")
	}

	if serviceDiscoveryOption.Prefix != ""{
		lastStr := serviceDiscoveryOption.Prefix[len(serviceDiscoveryOption.Prefix)-1:]
		if lastStr == "/"{
			return sd,errors.New(" last char = backslash.")
		}
	}

	if serviceDiscoveryOption.AutoCreateServiceLBType <=0 {
		serviceDiscoveryOption.AutoCreateServiceLBType = LOAD_BALANCE_ROBIN
	}

	ctx,cancelFunc := context.WithCancel(context.Background())
	serviceDiscovery.option 		= serviceDiscoveryOption
	serviceDiscovery.CancelFunc = cancelFunc
	serviceDiscovery.CancelCxt = ctx
	//serviceDiscovery.Status = "open"
	serviceDiscovery.WatchClose = make(chan bool)
	serviceDiscovery.list 			= make(map[string]*Service)
	serviceDiscovery.WatchMsg 		= make(chan ServiceChange,100)
	//初始化
	serviceDiscovery.Init()

	return serviceDiscovery,nil
}
//初始化
func (serviceDiscovery *ServiceDiscovery)Init( )error{
	//先把DB中的所有service 创建空的结构，添加到service list 中，初衷：监听etcd目录名得先确定了，如果有 创建时服务不存在 ，后期动态有人创建的情况，能够监听到
	for _,v :=range serviceDiscovery.option.ServiceManager.Pool{
		emptyService  := Service{
			Id:v.Id,
			Name: v.Name,
			DBKey:  serviceDiscovery.GetServiceDbKey( v.Name ),
			LBType: serviceDiscovery.option.AutoCreateServiceLBType,
			Log : serviceDiscovery.option.Log,
			CreateTime: GetNowTimeSecondToInt(),
		}
		serviceDiscovery.AddEmptyServiceManagerList(&emptyService)
	}
	//从etcd 中读取，已注册的服务
	_, err := serviceDiscovery.Discovery()
	if err != nil{
		return err
	}
	//监听3方服务变化
	serviceDiscovery.WatchThirdService()
	return nil
}

type EtcdKeyInfo struct {
	ServiceName string
	Ip string
	Port string
}

//服务发现 - 从分布式DB 中读取 并注册进容器中
func (serviceDiscovery *ServiceDiscovery)Discovery()( empty bool,err error){
	serviceDiscovery.option.Log.Info("Discovery , DiscoveryType:" + strconv.Itoa(serviceDiscovery.option.DiscoveryType))

	if serviceDiscovery.option.DiscoveryType == SERVICE_DISCOVERY_ETCD{
		//根据前缀 去ETCD 中读取所有<服务注册>数据
		allServiceEtcdList,err := serviceDiscovery.option.Etcd.GetListByPrefix(serviceDiscovery.option.Prefix)
		if err != nil{
			serviceDiscovery.option.Log.Error("etcd GetListByPrefix err:" +err.Error())
			return false,err
		}

		if len(allServiceEtcdList) == 0{
			serviceDiscovery.option.Log.Warn( " allServiceList is empty !")
			return true,nil
		}
		serviceDiscovery.option.Log.Info("allServiceEtcdList len:"+strconv.Itoa(len(allServiceEtcdList)))
		for k,_ := range allServiceEtcdList{
			//先把KEY 转换成 结构体
			etcdKeyInfo,err := serviceDiscovery.EtcdKeyCovertStruct(k)
			if err != nil{
				serviceDiscovery.option.Log.Error("failed : "+ k + " ,err:" + err.Error())
				continue
			}
			oriService , empty := serviceDiscovery.option.ServiceManager.GetByName(etcdKeyInfo.ServiceName)
			if empty{
				serviceDiscovery.option.Log.Error( " serviceName not in oriService pool :" + etcdKeyInfo.ServiceName)
				continue
			}
			//创建一个新的服务节点
			newServiceNode 	:= ServiceNode{
				ServiceName: etcdKeyInfo.ServiceName,
				ServiceId	: oriService.Id,
				Ip			: etcdKeyInfo.Ip,
				Port		: etcdKeyInfo.Port,
				Protocol	: SERVICE_PROTOCOL_GRPC,//仅支持GRPC,暂未实现其它
				IsSelfReg	: false,//非本服务注册，属于3方服务注册
				DBKey		: k,
			}
			serviceDiscovery.Register(newServiceNode)
		}
		return false,nil
	}else{
		//暂未实现
		return false,CONSUL_NOT_OPEN
	}
}

//给一个新的服务添加到管理列表中
func (serviceDiscovery *ServiceDiscovery)AddServiceManagerList(service *Service){
	serviceDiscovery.option.Log.Info("insert service to list :"+service.ToString())
	serviceDiscovery.list[service.Name] = service
}
//给一个新的服务添加到管理列表中
func (serviceDiscovery *ServiceDiscovery)AddEmptyServiceManagerList(service *Service){
	serviceDiscovery.option.Log.Info("insert empty service to list :"+service.ToString())
	serviceDiscovery.list[service.Name] = service
}

//设定一个监听器，用于监听；3方服务，一但出现变化通知上方
func (serviceManager *ServiceDiscovery)WatchThirdService(){
	if len(serviceManager.list) == 0{
		msg := "serviceManager.list ==0 , no need :WatchThirdService"
		serviceManager.option.Log.Info(msg)
		return
	}

	for serviceName,service:=range serviceManager.list{
		serviceManager.option.Log.Info("WatchThirdService serviceName:" + serviceName)
		go serviceManager.WatchOneService(service)
		//for _,serviceNode := range service.List{
		//	go serviceManager.WatchOneServiceNode(serviceNode)
		//}
	}
}
//监听一个服务
func (serviceDiscovery *ServiceDiscovery)WatchOneService(service *Service)error{
	prefix := "third service watching , "
	if serviceDiscovery.option.DiscoveryType == SERVICE_DISCOVERY_ETCD{
		//创建一个监听，如：/prefix/serviceName
		ctx,cancelFunc := context.WithCancel(context.Background())
		watchChan := serviceDiscovery.option.Etcd.Watch(ctx,service.DBKey)
		service.watchCancel = cancelFunc

		serviceDiscovery.option.Log.Info(prefix + " create :"+service.DBKey + " block.......")
		//进入阻塞模式
		isBreak := false
		for{
			select {
				case <-serviceDiscovery.WatchClose:
					isBreak = true
				case wresp := <- watchChan:
					for _, ev := range wresp.Events{
						action := ev.Type.String()
						key := string(ev.Kv.Key)
						val := string(ev.Kv.Value)

						msg := prefix  + " chan has event : " + action + ", key : " + key +  " val : " +val
						serviceDiscovery.option.Log.Warn(msg)

						serviceDiscovery.ServiceHasChange(service,action,key,val)

						//MyPrint(ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
						//fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

						//matchCode := strings.Replace(key,RuleEtcdConfigPrefix,"",-1)
						//matchCode = strings.Trim(matchCode," ")
						//matchCode = strings.Trim(matchCode,"/")
						//mylog.Warning(prefix , " matchCode : ",matchCode)
					}
				default:

			}
			if isBreak{
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		//for wresp := range watchChan{
		//	for _, ev := range wresp.Events{
		//		action := ev.Type.String()
		//		key := string(ev.Kv.Key)
		//		val := string(ev.Kv.Value)
		//
		//		msg := prefix  + " chan has event : " + action + ", key : " + key +  " val : " +val
		//		serviceDiscovery.option.Log.Warn(msg)
		//
		//		serviceDiscovery.ServiceHasChange(service,action,key,val)
		//
		//		//MyPrint(ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
		//		//fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		//
		//		//matchCode := strings.Replace(key,RuleEtcdConfigPrefix,"",-1)
		//		//matchCode = strings.Trim(matchCode," ")
		//		//matchCode = strings.Trim(matchCode,"/")
		//		//mylog.Warning(prefix , " matchCode : ",matchCode)
		//	}
		//}
	}else{
		return CONSUL_NOT_OPEN
	}

	return nil
}

func (serviceDiscovery *ServiceDiscovery)ServiceHasChange(service *Service,action string ,key string ,val string){
	prefix := "ServiceHasChange "
	etcdKeyInfo,err  := serviceDiscovery.EtcdKeyCovertStruct(key)
	if err != nil{
		return
	}
	hasService , ok := serviceDiscovery.list[etcdKeyInfo.ServiceName]
	if !ok {
		serviceDiscovery.option.Log.Error("ServiceName  not in serviceDiscovery list")
		return
	}
	serviceDiscovery.option.Log.Info(prefix + "now service list len:"+strconv.Itoa(len(service.List)))
	//etcd event : DELETE PUT
	if action == "PUT"{//添加|编辑
		for _,node := range hasService.List{
			if node.Ip == etcdKeyInfo.Ip && node.Port == etcdKeyInfo.Port{
				serviceDiscovery.option.Log.Error("has register,don't repeat")
				return
			}
		}
		newServiceNode := ServiceNode{
			//ServiceId: oriService.Id,
			ServiceName	: etcdKeyInfo.ServiceName,
			Ip			:etcdKeyInfo.Ip,
			Port		: etcdKeyInfo.Port,
			Protocol	: SERVICE_PROTOCOL_GRPC,//暂未实现
			IsSelfReg	: false,//非本服务注册，属于3方服务注册
			DBKey		: key,
		}
		serviceDiscovery.Register(newServiceNode)
	}else if action == "DELETE"{
		serviceDiscovery.option.Log.Info(prefix + "action is delete .  ")
		hasSearch := 0
		for _,node :=range service.List{
			if node.DBKey == key{
				hasSearch = 1
				serviceDiscovery.DelOneServiceNode(node)
				break
			}
		}
		if hasSearch == 0{
			serviceDiscovery.option.Log.Error(prefix + " delete failed ,  no search key:" + key)
		}
	}else{
		serviceDiscovery.option.Log.Error("WatchOneService event err:" + action)
		return
	}
	//serviceDiscovery.option.Log.Info(prefix + "up after ,node len:" + strconv.Itoa(len(service.List)))

	//给使用者发送信息，第一时间通知,暂时先注释掉，如调用者不接收容易出问题
	serviceChange := ServiceChange{}
	serviceChange.Action = action
	serviceChange.Name 	= etcdKeyInfo.ServiceName
	serviceChange.Ip 	= etcdKeyInfo.Ip
	serviceChange.Port 	= etcdKeyInfo.Port

	serviceDiscovery.WatchMsg <- serviceChange

}
//根据服务名，获取该服务下的一个节点，需要负载均衡
func (serviceDiscovery *ServiceDiscovery)GetLoadBalanceServiceNodeByServiceName(serviceName string,factor string)(serviceNode *ServiceNode,err error){
	serviceDiscovery.option.Log.Info("GetLoadBalanceServiceNodeByServiceName:" + serviceName)

	service ,ok := serviceDiscovery.list[serviceName]
	if !ok {
		return serviceNode,errors.New(serviceName + "serviceName 不存在 map 中 ")
	}
	if len(service.List) <=0 {
		return serviceNode,errors.New(serviceName + "serviceNode List is empty~")
	}

	node := serviceDiscovery.balanceHost(service,factor)
	return node,nil
}
//通过ETCD可以获取到服务的IP list ，这里是负载，决定 用哪 个IP
func (serviceManager *ServiceDiscovery)balanceHost(service *Service ,factor string)*ServiceNode{
	if service.LBType == LOAD_BALANCE_ROBIN{
		r := GetRandIntNumRange(0,len(service.List))
		node := service.List[r]
		return node
	}else if service.LBType == LOAD_BALANCE_HASH{
		lastStr := []byte(factor)[len(factor)-1:]
		lastByte := byte(lastStr[0])
		i := int(lastByte) % len(service.List)
		node := service.List[i]
		return node
	}else{
		serviceManager.option.Log.Error("err: service.LBType")
	}

	return nil
}

//
func (serviceDiscovery *ServiceDiscovery)DelOneServiceNode(serviceNode *ServiceNode){
	serviceDiscovery.option.Log.Info("DelOneServiceNode:" + serviceNode.ServiceName)

	serviceDiscovery.OPLock.Lock()
	defer serviceDiscovery.OPLock.Unlock()

	service ,ok := serviceDiscovery.list[serviceNode.ServiceName]
	if !ok {
		serviceDiscovery.option.Log.Error("DelOneServiceNode err: ServiceName not in map")
		return
	}
	if serviceDiscovery.option.DiscoveryType == SERVICE_DISCOVERY_ETCD{
		hasDel := false
		for k,n := range service.List{
			if n.DBKey == serviceNode.DBKey{
				service.List = append( service.List[:k],service.List[k+1:]... )
				hasDel = true
				break
			}
		}

		if !hasDel{
			serviceDiscovery.option.Log.Warn("DelOneServiceNode err:no search node")
			return
		}
		serviceDiscovery.option.Log.Info("delete one ok.")
		if serviceNode.IsSelfReg{
			err := serviceDiscovery.option.Etcd.DelOne(serviceNode.DBKey)
			if err != nil{
				serviceDiscovery.option.Log.Error("service.etcd.DelOne err " + err.Error())
				return
			}
			serviceDiscovery.option.Log.Info("delete etcd one ok.")
		}
		//if len(service.List) == 0{
		//	service.Lease.Revoke(service.LeaseCancelCtx,service.LeaseGrantId)
		//	delete(serviceDiscovery.list,service.Name)
		//}
	}else{
		str := CONSUL_NOT_OPEN.Error()
		serviceDiscovery.option.Log.Error(str)
	}
	//serviceManager.option.Log.Info("etcd DelOne :" + err.Error())
}
//动态（租约）注册一个服务，一但服务停止该服务自动取消
func (serviceDiscovery *ServiceDiscovery)Register( serviceNode ServiceNode )error {
	debugInfo := "RegisterService nodeInfo:"+ serviceNode.ToString()
	serviceDiscovery.option.Log.Info(debugInfo)

	serviceDiscovery.OPLock.Lock()
	defer serviceDiscovery.OPLock.Unlock()
	//从服务发现中，找寻服务是否存在(新服务必须申请且注册到DB中)
	oriService ,empty := serviceDiscovery.option.ServiceManager.GetByName(serviceNode.ServiceName)
	if empty{
		msg := "Register serviceName err:" + serviceNode.ServiceName
		serviceDiscovery.option.Log.Error(msg)
		return errors.New(msg)
	}

	serviceNode.ServiceId = oriService.Id
	//找到公用 DB KEY
	serviceNode.DBKey = serviceDiscovery.GetServiceNodeDbKey( serviceNode.ServiceName,serviceNode.Ip , serviceNode.Port)
	serviceNode.Log = serviceDiscovery.option.Log

	service,ok := serviceDiscovery.list[serviceNode.ServiceName]
	if !ok{
		//自动创建新服务
		serviceDiscovery.option.Log.Info("auto create service:"+ serviceNode.ServiceName)
		dbKey := serviceDiscovery.GetServiceDbKey( serviceNode.ServiceName)
		newService := Service{
			Id:serviceNode.ServiceId,
			Name: serviceNode.ServiceName,
			DBKey: dbKey,
			LBType: serviceDiscovery.option.AutoCreateServiceLBType,
			Log : serviceDiscovery.option.Log,
			CreateTime: GetNowTimeSecondToInt(),
		}

		err := newService.NewServiceNode(serviceNode)
		if err != nil{
			serviceDiscovery.option.Log.Error(err.Error())
			return err
		}
		serviceDiscovery.AddServiceManagerList(&newService)

		service = &newService
	}else{
		err := service.NewServiceNode(serviceNode)
		if err != nil{
			serviceDiscovery.option.Log.Error(err.Error())
			return err
		}
	}

	if serviceNode.IsSelfReg{//如果是自己注册的服务，需要再申请个租约
		if serviceDiscovery.option.DiscoveryType == SERVICE_DISCOVERY_ETCD {

			ctx, cancelFunc := context.WithCancel(context.Background())
			lease, leaseGrantId, err := serviceDiscovery.option.Etcd.NewLeaseGrand(ctx, 60, 1)
			if err != nil {
				serviceDiscovery.option.Log.Error(" apply etcd leaseGrand failed,err:"+err.Error())
				cancelFunc()
				return err
			}
			now := GetNowTimeSecondToInt()
			//key := serviceDiscovery.GetServiceNodeDbKey(serviceNode.ServiceName,serviceNode.Ip,serviceNode.Port)
			serviceDiscovery.option.Log.Info("create service("+serviceNode.ServiceName+") lease... key:"+serviceNode.DBKey)
			val := strconv.Itoa(now)
			_, err = serviceDiscovery.option.Etcd.putLease(ctx, leaseGrantId, serviceNode.DBKey, val)
			if err != nil {
				serviceDiscovery.option.Log.Error(" put etcd lease failed,err:"+err.Error())
				cancelFunc()
				return errors.New("service.etcd.PutOne err :" + err.Error())
			}

			service.Lease = lease
			service.LeaseGrantId = leaseGrantId
			service.LeaseCancelCtx = ctx
		}
	}

	serviceDiscovery.option.Log.Info("register one service success :"+ serviceNode.ToString())
	return nil
}
func (serviceManager *ServiceDiscovery)ShowJsonByService()string{
	if len(serviceManager.list) <=0 {
		return ""
	}

	jsonByte,err := json.Marshal(serviceManager.list)
	jsonStr := string(jsonByte)
	MyPrint(jsonStr , "err:",err)

	return jsonStr
}
func (serviceManager *ServiceDiscovery)ShowJsonByNodeServer()string{
	if len(serviceManager.list) <=0 {
		return ""
	}

	list := make(map[string][]string)
	for _,service:=range serviceManager.list{
		if len(service.List) <= 0{
			//list[service.Name] = []*ServiceNode{}
			continue
		}

		for _,node :=range service.List{
			list[node.Ip] = append(list[node.Ip],service.Name)
		}

	}

	jsonByte,err := json.Marshal(list)
	jsonStr := string(jsonByte)
	MyPrint(jsonStr , "err:",err)

	return jsonStr
}

//整体关闭
func (serviceManager *ServiceDiscovery)Shutdown( ){
	serviceManager.option.Log.Warn("service Shutdown:")
	//serviceManager.Status = "shutdown"
	//先把监听协程给关闭了
	serviceManager.WatchClose <- true
	for _,service:=range serviceManager.list{
		for _,serviceNode:=range service.List{
			//if serviceNode.IsSelfReg {
				serviceManager.DelOneServiceNode(serviceNode)
			//}else{
			//	//取消 所有 watch
			//	service.watchCancel()
			//}
		}
	}
	serviceManager.CancelFunc()
	close(serviceManager.WatchMsg)
	close(serviceManager.WatchClose)
}

func (serviceDiscovery *ServiceDiscovery)EtcdKeyCovertStruct(key string)(etcdKeyInfo EtcdKeyInfo,err error){
	eg := "/prefix/serviceName/127.0.0.1:9999"
	totalLen :=  strings.Split(key,"/")
	if len(totalLen) != 4{
		errMsg := "totalLen != 4"
		return etcdKeyInfo,errors.New(errMsg)
	}
	//把前缀去掉~
	keyUriStr := strings.Replace(key,serviceDiscovery.option.Prefix,"",-1)
	//切割
	keyUriArr := strings.Split(keyUriStr,"/")
	if len(keyUriArr)!= 3{
		msg := "parser etcdKey err: backslash != 3 , eg:"+eg
		serviceDiscovery.option.Log.Error(msg)
		return etcdKeyInfo,errors.New(msg)
	}
	etcdKeyInfo.ServiceName = keyUriArr[1]
	ipPort :=strings.Split( keyUriArr[2],":")
	if len(ipPort) != 2{
		msg := "parser etcdKey err: ipPort != 2 , eg:"+eg
		serviceDiscovery.option.Log.Error(msg)
		return etcdKeyInfo,errors.New(msg)
	}

	etcdKeyInfo.Ip = ipPort[0]
	etcdKeyInfo.Port = ipPort[1]

	return etcdKeyInfo,nil
}

func (serviceDiscovery *ServiceDiscovery)GetServiceDbKey(serviceName string)string{
	return serviceDiscovery.option.Prefix + "/" + serviceName
}

func (serviceDiscovery *ServiceDiscovery)GetServiceNodeDbKey(serviceName string,ip string ,port string)string{
	return serviceDiscovery.GetServiceDbKey(serviceName) + "/" + ip + ":" + port
}