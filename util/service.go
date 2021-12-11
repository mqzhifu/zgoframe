package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)
const (
	SERVICE_PROTOCOL_HTTP = 1
	SERVICE_PROTOCOL_GRPC = 2
	SERVICE_PROTOCOL_WEBSOCKET = 3
	SERVICE_PROTOCOL_TCP = 4

	SERVICE_DISCOVERY_ETCD = 1
	SERVICE_DISCOVERY_CONSUL = 2
)
//type DynamicService struct {
//	Name 	string
//	Value 	string
//	cancelFunc 	context.CancelFunc
//	cancelCtx context.Context
//	LeaseGrantId clientv3.LeaseID
//	Lease clientv3.Lease
//}

//一个服务下面的一个节点
type ServiceNode struct {
	ServiceName string	`json:"service_name"`
	Ip			string	`json:"ip"`
	Port		string	`json:"port"`
	Protocol 	string	`json:"protocol"`
	Desc 		string	`json:"desc"`
	Status 		int		`json:"status"`
	IsSelfReg 	bool	`json:"is_self_reg"`
	RegTime 	int 	`json:"reg_time"`
	watch 		context.CancelFunc	`json:"-"`

	LeaseGrantId clientv3.LeaseID	`json:"-"`
	Lease clientv3.Lease			`json:"-"`
	LeaseCancelCtx context.Context	`json:"-"`
}
//一个服务
type Service struct {
	Id int	`json:"id"`
	Name string	`json:"name"`
	List []*ServiceNode	`json:"list"`
}
//服务管理器
//ServiceManager.list->service.list->serviceNode
type ServiceManager struct {

	list      map[string]*Service      //可用服务列表
	etcd      *MyEtcd

	option    ServiceOption
	DiscoveryType int
	Protocol	int 		//1http 2grpc 3websocket 4tcp

	//consul 		int //待处理,https://github.com/hashicorp/consul/
	//selfList  map[string]DynamicService //自己注册的服务列表
}

type ServiceOption struct {
	Etcd 	*MyEtcd
	Log 	*zap.Logger
	DiscoveryType int
	Prefix	string
	//TestHttpGamematchPushReceiveHsot string
}

func NewServiceManager(serviceOption ServiceOption)*ServiceManager {
	serviceOption.Log.Info("NewService")

	serviceManager := new(ServiceManager)

	serviceManager.DiscoveryType = serviceOption.DiscoveryType
	serviceManager.etcd = serviceOption.Etcd
	serviceManager.option = serviceOption
	serviceManager.list = make(map[string]*Service)
	serviceManager.ReadThirdService()



	return serviceManager
}
//根据服务名，获取该服务下的一个节点，需要负载均衡
func (serviceManager *ServiceManager)GetLoadBalanceServiceNodeByServiceName(serviceName string)(serviceNode *ServiceNode,err error){
	MyPrint("serviceManager list:",serviceManager.list)
	service ,ok := serviceManager.list[serviceName]
	if !ok {
		return serviceNode,errors.New(serviceName + "serviceName 不存在 map 中 ")
	}
	//if len(service.List) <=0 {
	//	return serviceNode,errors.New(serviceName + "serviceNode List is empty~")
	//}

	node := serviceManager.balanceHost(service)
	return node,nil
}
func (serviceManager *ServiceManager)AddServiceManagerList(service *Service){
	serviceManager.list[service.Name] = service
}
func (service *Service)AddServiceList(node *ServiceNode){
	if node.ServiceName == ""{
		MyPrint("AddServiceList err: serviceName empty")
	}

	if node.Ip == ""{
		MyPrint("AddServiceList err: IP empty")
	}

	if node.Port == ""{
		MyPrint("AddServiceList err: port empty")
	}

	service.List = append(service.List,node)
	//serviceManager.list[service.Name] = service
}
//服务发现
func (serviceManager *ServiceManager)Discovery()(list map[string][]Service,err error){
	if serviceManager.DiscoveryType == SERVICE_DISCOVERY_ETCD{
		allServiceList,err := serviceManager.etcd.GetListByPrefix(serviceManager.option.Prefix)
		if err != nil{
			serviceManager.option.Log.Error("ReadAdnRegThird err:" +err.Error())
			return list,err
		}
		if len(allServiceList) == 0{
			serviceManager.option.Log.Warn( " allServiceList is empty !")
			return list,err
		}

		for k,_ := range allServiceList{
			str := strings.Replace(k,serviceManager.option.Prefix,"",-1)
			//MyPrint(str,k)
			serviceArr := strings.Split(str,"/")
			//MyPrint(serviceArr)
			serviceName := serviceArr[1]
			hasServiceInc,ok := serviceManager.list[serviceName]
			var serviceInc *Service
			if !ok {
				serviceInc = &Service{
					Name: serviceName,
				}
			}else{
				serviceInc = hasServiceInc
			}
			ipPort :=strings.Split( serviceArr[2],":")
			newServiceNode := ServiceNode{
				ServiceName: serviceName,
				Ip:ipPort[0],
				Port: ipPort[1],
				Protocol: "",
				IsSelfReg: false,
			}
			MyPrint("newServiceNode:",newServiceNode)
			//serviceInc.List = append(serviceInc.List,newServiceNode)
			serviceManager.AddServiceManagerList(serviceInc)
			serviceInc.AddServiceList(&newServiceNode)
		}
	}else{

	}
	return list,err
}
//从配置中心读取：3方可用服务列表,注册到内存中
func (serviceManager *ServiceManager)ReadThirdService( ){
	//从etcd 中读取，已注册的服务
	serviceManager.Discovery()
	//监听3方服务变化
	go serviceManager.WatchThirdService()
}
//设定一个监听器，用于监听；3方服务，一但出现变化通知上方
func (serviceManager *ServiceManager)WatchThirdService(){
	for serviceName,service:=range serviceManager.list{
		serviceManager.option.Log.Info("WatchThirdService serviceName:" + serviceName)
		for _,serviceNode := range service.List{
			if serviceManager.DiscoveryType == SERVICE_DISCOVERY_ETCD{
				ctx,cancelFunc := context.WithCancel(context.Background())
				watchChann := serviceManager.etcd.Watch(ctx,serviceManager.option.Prefix)
				//service.watchList = append(service.watchList,cancelFunc)
				serviceNode.watch = cancelFunc
				prefix := "third service watching receive , "
				//service.option.Log.Notice(prefix , " , new key : ",service.option.Prefix)
				//watchChann := myetcd.Watch("/testmatch")
				for wresp := range watchChann{
					for _, ev := range wresp.Events{
						action := ev.Type.String()
						key := string(ev.Kv.Key)
						val := string(ev.Kv.Value)

						serviceManager.option.Log.Warn(prefix  + " chan has event : " + action)
						serviceManager.option.Log.Warn(prefix +" key : " + key)
						serviceManager.option.Log.Warn(prefix + " val : " +val)
						//MyPrint(ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
						//fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

						//matchCode := strings.Replace(key,RuleEtcdConfigPrefix,"",-1)
						//matchCode = strings.Trim(matchCode," ")
						//matchCode = strings.Trim(matchCode,"/")
						//mylog.Warning(prefix , " matchCode : ",matchCode)
					}
				}
			}
		}
	}
}
//删除自己的服务
func (serviceManager *ServiceManager)DelOneServiceNode(serviceNode *ServiceNode){
	//now := GetNowTimeSecondToInt()
	//putResponse,err := service.etcd.PutOne( service.option.Prefix +"/"+serviceName +"/"+ipPort , strconv.Itoa(now))
	service ,ok := serviceManager.list[serviceNode.ServiceName]
	if !ok {
		MyPrint("DelOneServiceNode err: ServiceName not in map")
		return
	}
	if serviceManager.DiscoveryType == SERVICE_DISCOVERY_ETCD{
		name := serviceManager.option.Prefix +"/"+serviceNode.ServiceName
		err := serviceManager.etcd.DelOne(name)
		if err != nil{
			ExitPrint("service.etcd.PutOne err ",err.Error())
		}
		for _,node := range service.List{
			node.watch()
		}
		hasDel := false
		for k,n := range service.List{
			if n == serviceNode{
				service.List = append( service.List[:k],service.List[k:]... )
				hasDel = true
				break
			}
		}
		if !hasDel{
			MyPrint("DelOneServiceNode err:no search node")
		}
		//delete(service.list,serviceName)
		serviceNode.Lease.Revoke(serviceNode.LeaseCancelCtx,serviceNode.LeaseGrantId)
	}
	//serviceManager.option.Log.Info("etcd DelOne :" + err.Error())
}
//通过ETCD可以获取到服务的IP list ，这里是负载，决定 用哪 个IP
func (serviceManager *ServiceManager)balanceHost(service *Service)*ServiceNode{
	r := GetRandIntNumRange(0,len(service.List))
	node := service.List[r]

	return node
}
//动态（租约）注册一个服务，一但服务停止该服务自动取消
func (serviceManager *ServiceManager)Register(service Service ,serviceNode ServiceNode)error {
	debugInfo := "Register serviceName:"+ service.Name +  " node :"+ serviceNode.Ip +":" +serviceNode.Port +", "+ serviceNode.Protocol
	serviceManager.option.Log.Info(debugInfo)
	_,ok := serviceManager.list[service.Name]
	if !ok{
		serviceManager.option.Log.Info("Register serviceName add new key")
		serviceManager.AddServiceManagerList(&service)
	}
	serviceNode.IsSelfReg = true

	if serviceManager.DiscoveryType == SERVICE_DISCOVERY_ETCD {

		ctx, cancelFunc := context.WithCancel(context.Background())
		lease, leaseGrantId, err := serviceManager.etcd.NewLeaseGrand(ctx, 60, 1)
		//_      , leaseGrantId, err := serviceManager.etcd.NewLeaseGrand(ctx, 60, 1)
		//service.option.Log.Debug("leaseGrantId : ",leaseGrantId ,err)
		if err != nil {
			cancelFunc()
			return err
		}
		now := GetNowTimeSecondToInt()
		//putResponse,err := service.etcd.PutOne( service.option.Prefix +"/"+serviceName +"/"+ipPort , strconv.Itoa(now))
		key := serviceManager.option.Prefix + "/" + service.Name + "/" + serviceNode.Ip + ":" + serviceNode.Port
		val := strconv.Itoa(now)
		_, err = serviceManager.etcd.putLease(ctx, leaseGrantId, key, val)
		if err != nil {
			cancelFunc()
			return errors.New("service.etcd.PutOne err :" + err.Error())
		}

		//dynamicService := DynamicService{
		//	cancelFunc:   cancelFunc,
		//	cancelCtx:    ctx,
		//	Name:         serviceName,
		//	Value:        ipPort,
		//	LeaseGrantId: leaseGrantId,
		//	Lease:        lease,
		//}
		//service.selfList[serviceName] = dynamicService

		serviceNode.Lease = lease
		serviceNode.LeaseGrantId = leaseGrantId
		serviceNode .LeaseCancelCtx = ctx
		//serviceInc := serviceManager.list[service.Name]
		service.AddServiceList(&serviceNode)
		//serviceInc.List = append(serviceInc.List,serviceNode)
		//serviceManager.list[service.Name].List =mylist
	}
	serviceManager.option.Log.Info("register one service success.")
	return nil
}
func (serviceManager *ServiceManager)ShowJsonByService()string{
	if len(serviceManager.list) <=0 {
		return ""
	}

	jsonByte,err := json.Marshal(serviceManager.list)
	jsonStr := string(jsonByte)
	MyPrint(jsonStr , "err:",err)

	return jsonStr
}
func (serviceManager *ServiceManager)ShowJsonByNodeServer()string{
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
func (serviceManager *ServiceManager)Shutdown( ){
	serviceManager.option.Log.Warn("service Shutdown:")

	for _,service:=range serviceManager.list{
		for _,serviceNode:=range service.List{
			if serviceNode.IsSelfReg {
				serviceManager.DelOneServiceNode(serviceNode)
			}
			//取消 所有 watch
			serviceNode.watch()
		}
	}

	//if len(service.selfList ) >=0  {
	//	for _,dynamicService :=range service.selfList{
	//		service.option.Log.Info("service cancelFunc :" +dynamicService.Name)
	//		//dynamicService.cancelFunc()

		//}
	//}else{
	//	service.option.Log.Warn( "service.selfList  <= 0")
	//}

	//if len(service.watchList ) >=0  {
	//	for _,cancelFunc :=range service.watchList{
	//		cancelFunc()
	//	}
	//}else{
	//	service.option.Log.Warn( "service.watchList  <= 0")
	//}
	//service.option.Log.Warn("service shutdown.")
}
//给一个服务，发送一条http消息
func (serviceManager *ServiceManager)HttpPost(serviceName string,uri string,data interface{}) (responseMsgST ResponseMsgST,errs error){
	//先从池中找到该服务
	service ,ok := serviceManager.list[serviceName]
	if !ok {
		return responseMsgST,errors.New(serviceName + "serviceName 不存在 map 中 ")
	}
	if len(service.List) <=0 {
		return responseMsgST,errors.New(serviceName + "serviceNode List is empty~")
	}
	//找一个该服务下的一个IP地址
	node := serviceManager.balanceHost(service)
	serviceHost := node.Ip + ":" + node.Port
	//serviceHost = serviceManager.option.TestHttpGamematchPushReceiveHsot
	url := "http://"+serviceHost + "/" + uri
	serviceManager.option.Log.Info("HttpPost" + serviceName + serviceHost + uri + url)
	jsonStr, _ := json.Marshal(data)
	serviceManager.option.Log.Info("jsonStr:" + string(jsonStr))
	//ExitPrint(1111)
	req, errs := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", "application/json")
	defer req.Body.Close()
	//
	if errs != nil {
		return responseMsgST,errors.New("NewRequest err")
	}
	//5秒超时
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	//service.option.Log.Debug(resp,error)
	if error != nil {
		return responseMsgST,errors.New("client.Do  err"+error.Error())
	}

	if resp.StatusCode != 200{
		return responseMsgST,errors.New("http response code != 200")
	}

	if resp.ContentLength == 0{
		return responseMsgST,errors.New("http response content = 0")
	}
	contentJsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return responseMsgST,errors.New("ioutil.ReadAll err : "+err.Error() )
	}

	errs = json.Unmarshal(contentJsonStr,&responseMsgST)
	if errs != nil{
		return responseMsgST,errors.New(" json.Unmarshal html content err : "+err.Error() )
	}

	//service.option.Log.Debug("responseMsgST : ",responseMsgST)
	return responseMsgST,nil
}