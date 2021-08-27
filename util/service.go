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
	POROTOCOL_HTTP = 1
	POROTOCOL_GRPC = 2
)
type DynamicService struct {
	Name 	string
	Value 	string
	cancelFunc 	context.CancelFunc
	cancelCtx context.Context
	LeaseGrantId clientv3.LeaseID
	Lease clientv3.Lease
}

type Service struct {
	list      map[string][]string       //可用服务列表
	selfList  map[string]DynamicService //自己注册的服务列表
	watchList []context.CancelFunc
	etcd      *MyEtcd
	option    ServiceOption
	Porotocol	int 		//1http 2grpc
}

type ServiceOption struct {
	Etcd 	*MyEtcd
	Log 	*zap.Logger
	Prefix	string
	TestHttpGamematchPushReceiveHsot string
}

func NewService(serviceOption ServiceOption)*Service {
	service := new(Service)

	service.etcd = serviceOption.Etcd
	service.option = serviceOption
	service.selfList = make(map[string]DynamicService)
	service.ReadAdnRegThird()

	serviceOption.Log.Info("NewService")

	return service
}
//从配置中心读取：3方可用服务列表,注册到内存中
func (service *Service)ReadAdnRegThird( ){
	//从etcd 中读取，已注册的服务
	allServiceList,err := service.etcd.GetListByPrefix(service.option.Prefix)
	if err != nil{
		service.option.Log.Error("ReadAdnRegThird err:" +err.Error())
		return
	}
	if len(allServiceList) == 0{
		service.option.Log.Warn( " allServiceList is empty !")
		return
	}

	serviceListMap := make(map[string][]string)
	for k,_ := range allServiceList{
		str := strings.Replace(k,service.option.Prefix,"",-1)
		//MyPrint(str,k)
		serviceArr := strings.Split(str,"/")
		serviceListMap[serviceArr[1]] = append(serviceListMap[serviceArr[1]], serviceArr[2])
	}
	//service.option.Log.Info("RegThird:",serviceListMap)
	service.list = serviceListMap
	//service.option.Log.Debug(serviceListMap)
	//AddRoutineList("WatchThridService")
	//监听3方服务变化
	go service.WatchThridService()
	//service.option.Goroutine.CreateExec(service,"WatchThridService")
}
//设定一个监听器，用于监听；3方服务，一但出现变化通知上方
func (service *Service)WatchThridService(){
	ctx,cancelFunc := context.WithCancel(context.Background())
	watchChann := service.etcd.Watch(ctx,service.option.Prefix)
	service.watchList = append(service.watchList,cancelFunc)
	prefix := "third service watching receive , "
	//service.option.Log.Notice(prefix , " , new key : ",service.option.Prefix)
	//watchChann := myetcd.Watch("/testmatch")
	for wresp := range watchChann{
		for _, ev := range wresp.Events{
			action := ev.Type.String()
			key := string(ev.Kv.Key)
			val := string(ev.Kv.Value)

			service.option.Log.Warn(prefix  + " chan has event : " + action)
			service.option.Log.Warn(prefix +" key : " + key)
			service.option.Log.Warn(prefix + " val : " +val)
			//MyPrint(ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
			//fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			//matchCode := strings.Replace(key,RuleEtcdConfigPrefix,"",-1)
			//matchCode = strings.Trim(matchCode," ")
			//matchCode = strings.Trim(matchCode,"/")
			//mylog.Warning(prefix , " matchCode : ",matchCode)
		}
	}
}
//注册自己的服务
//func (service *Service)RegOne(serviceName string,ipPort string){
//	now := GetNowTimeSecondToInt()
//	putResponse,err := service.etcd.PutOne( service.option.Prefix +"/"+serviceName +"/"+ipPort , strconv.Itoa(now))
//	if err != nil{
//		ExitPrint("service.etcd.PutOne err ",err.Error())
//	}
//	service.selfList[serviceName] = ipPort
//	service.option.Log.Info("etcd put one ",putResponse.Header)
//}
//删除自己的服务
func (service *Service)DelOne(serviceName string,ipPort string){
	//now := GetNowTimeSecondToInt()
	//putResponse,err := service.etcd.PutOne( service.option.Prefix +"/"+serviceName +"/"+ipPort , strconv.Itoa(now))
	name := service.option.Prefix +"/"+serviceName
	err := service.etcd.DelOne(name)
	if err != nil{
		ExitPrint("service.etcd.PutOne err ",err.Error())
	}
	service.option.Log.Info("etcd DelOne :" + err.Error())
}
//通过ETCD可以获取到服务的IP list ，这里是负载，决定 用哪 个IP
func (service *Service)balanceHost(list []string)string{
	return list[0]
}
//动态（租约）注册一个服务，一但服务停止该服务自动取消
func (service *Service)RegOneDynamic(serviceName string,ipPort string)error{
	ctx,cancelFunc  := context.WithCancel(context.Background())
	lease,leaseGrantId ,err := service.etcd.NewLeaseGrand(ctx,60,1)
	//service.option.Log.Debug("leaseGrantId : ",leaseGrantId ,err)
	if err != nil{
		cancelFunc()
		return err
	}
	now := GetNowTimeSecondToInt()
	//putResponse,err := service.etcd.PutOne( service.option.Prefix +"/"+serviceName +"/"+ipPort , strconv.Itoa(now))
	key := service.option.Prefix +"/"+serviceName +"/"+ipPort
	val := strconv.Itoa(now)
	_,err = service.etcd.putLease(ctx,leaseGrantId,key,val)
	if err != nil{
		cancelFunc()
		return errors.New("service.etcd.PutOne err :" + err.Error())
	}

	dynamicService := DynamicService{
		 cancelFunc:cancelFunc,
		 cancelCtx :ctx,
		 Name: serviceName,
		 Value: ipPort,
		 LeaseGrantId : leaseGrantId,
		 Lease:lease,
	}
	service.selfList[serviceName] = dynamicService
	return nil
}
func (service *Service)Shutdown( ){
	service.option.Log.Warn("service Shutdown:")
	if len(service.selfList ) >=0  {
		for _,dynamicService :=range service.selfList{
			service.option.Log.Info("service cancelFunc :" +dynamicService.Name)
			//dynamicService.cancelFunc()
			dynamicService.Lease.Revoke(dynamicService.cancelCtx,dynamicService.LeaseGrantId)
		}
	}else{
		service.option.Log.Warn( "service.selfList  <= 0")
	}

	if len(service.watchList ) >=0  {
		for _,cancelFunc :=range service.watchList{
			cancelFunc()
		}
	}else{
		service.option.Log.Warn( "service.watchList  <= 0")
	}
	service.option.Log.Warn("service shutdown.")

}
//给一个服务，发送一条http消息
func (service *Service)HttpPost(serviceName string,uri string,data interface{}) (responseMsgST ResponseMsgST,errs error){
	//先从池中找到该服务
	serviceIpList ,ok := service.list[serviceName]
	if !ok {
		return responseMsgST,errors.New(serviceName + " 不存在 map 中 ")
	}
	//找一个该服务下的一个IP地址
	serviceHost := service.balanceHost(serviceIpList)
	serviceHost = service.option.TestHttpGamematchPushReceiveHsot
	url := "http://"+serviceHost + "/" + uri
	service.option.Log.Debug("HttpPost" + serviceName + serviceHost + uri + url)
	jsonStr, _ := json.Marshal(data)
	service.option.Log.Debug("jsonStr:" + string(jsonStr))
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