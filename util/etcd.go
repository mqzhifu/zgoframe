package util

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ResponseMsgST struct {
	Code	int
	Msg 	interface{}
	//Code 	int `json:"code"`
	//Data 	interface{} `json:"data"`
}

type MyEtcd struct {
	cli         *clientv3.Client
	option      EtcdOption
	AppConflist map[string]string
}

type EtcdOption struct {
	AppName string
	AppKey string
	AppENV	string
	FindEtcdUrl		string
	LinkAddressList	[]string
	Username string
	Password string
	Ip 		string
	Port 	string
	Log *zap.Logger
}
//通过http 请求配置中心，获取返回结果
type EtcdHttpResp struct {
	Code int        `json:"code"`
	Data Etcdconfig `json:"data"`
}
type Etcdconfig struct {
	Username	string `json:"username"`
	Password	string	`json:"password"`
	Hosts		[]string `json:"hosts"`
}
func NewMyEtcdSdk(etcdOption EtcdOption)(myEtcd *MyEtcd,errs error){
	myEtcd = new (MyEtcd)
	var clientv3Config  clientv3.Config
	dns := etcdOption.Ip + ":" + etcdOption.Port
	etcdOption.Log.Info("etcd connect:"+dns)
	//建立连接，优先使用 FindEtcdUrl ，走网络发现，这种 扩展更好
	if etcdOption.FindEtcdUrl != ""{
		etcdOption.Log.Info("use FindEtcdUrl node")
		//获取etcd 服务器配置信息
		jsonStruct ,errs := getEtcdHostPort(etcdOption)
		if errs != nil {
			return nil,errors.New("http request err :" + errs.Error())
		}
		//etcdOption.Log.Info("etcdConfig ip list : ", json.Marshal(jsonStruct.Data.Hosts))
		etcdOption.LinkAddressList = jsonStruct.Data.Hosts
		//开启建立连接
		clientv3Config  = clientv3.Config{
			Endpoints:  jsonStruct.Data.Hosts,
			DialTimeout: 5 * time.Second,
			Username: jsonStruct.Data.Username,
			Password: jsonStruct.Data.Password,
		}
	}else{
		etcdOption.Log.Info("use configFile node")
		dns := etcdOption.Ip + ":" + etcdOption.Port

		etcdOption.Log.Info("etcd confg: " +  dns+ " , username:"+etcdOption.Username + " ps:"+etcdOption.Password)

		etcdOption.LinkAddressList = append(etcdOption.LinkAddressList,dns)

		clientv3Config  = clientv3.Config{
			//#http://114.116.212.202/account/dev/v1/sys/etcd
			Endpoints:  etcdOption.LinkAddressList,
			DialTimeout: 5 * time.Second,
			Username: etcdOption.Username,
			Password: etcdOption.Password,
		}
	}

	cli, errs := clientv3.New(clientv3Config)
	if errs != nil {
		return nil,errors.New("clientv3.New error :  " + errs.Error())
	}
	myEtcd.cli = cli
	myEtcd.option = etcdOption
	//获取自己项目想着的配置信息
	myEtcd.iniAppConf()
	return myEtcd,nil
}
func (myEtcd *MyEtcd)Shutdown(){
	myEtcd.cli.Close()
	myEtcd.option.Log.Warn("etcd shutdown.")
}
//寻找etcd host ip 列表
func getEtcdHostPort(etcdOption EtcdOption)( etcdHttpResp EtcdHttpResp,err error){
	etcdOption.Log.Info("http.get remote etcd host:port  : " +etcdOption.FindEtcdUrl)
	resp, errs := http.Get(etcdOption.FindEtcdUrl)
	if errs != nil{
		return etcdHttpResp,errs
	}
	htmlContentJson,_ := ioutil.ReadAll(resp.Body)
	//解析请求回来的配置信息
	if len(htmlContentJson) == 0{
		return etcdHttpResp,errors.New("http request content empty! :" + errs.Error())
	}
	//jsonStruct :=  EtcdHttpResp{}
	errs = json.Unmarshal(htmlContentJson,&etcdHttpResp)
	if errs != nil {
		return etcdHttpResp,errors.New("http request err : Unmarshal " + errs.Error())
	}
	//etcdConfig := strings.Split(jsonStruct.Msg.(string),",")
	if len(etcdHttpResp.Data.Hosts) == 0 {
		return etcdHttpResp,errors.New("http request err : etcdConfig is empty ")
	}
	return etcdHttpResp,errs
}
//申请一个X秒TTL的租约
//autoKeepAlive:一个租约到时候后，是否自动继续续租
func (myEtcd *MyEtcd)NewLeaseGrand(ctx context.Context ,ttl int64,autoKeepAlive int)(l clientv3.Lease,leaseGrantId clientv3.LeaseID,e error){
	//创建一个租约实体
	lease :=  clientv3.NewLease(myEtcd.cli)
	//授权实体：申请一个60秒的 租约 实体
	leaseGrant, err := lease.Grant(ctx, ttl)
	if  err != nil {
		myEtcd.option.Log.Error("lease.Grant err :" + err.Error())
		return l,leaseGrantId,err
	}
	if autoKeepAlive == 1{
		_,err :=lease.KeepAlive(ctx,leaseGrant.ID)
		if err !=nil{
			myEtcd.option.Log.Error("lease.KeepAlive err :" + err.Error())
			return l,leaseGrantId,err
		}
	}
	//myEtcd.option.Log.Info("create New Lease and Grand ,  ttl :",ttl, " id : ",leaseGrant.ID)
	return lease,leaseGrant.ID,nil
}

//往一个租约里写入内容，跟NewLeaseGrand联合使用的
func (myEtcd *MyEtcd)putLease(ctx context.Context,leaseId clientv3.LeaseID,k string,v string)(putResponse *clientv3.PutResponse,err error){
	//创建一个KV 容器
	kv := clientv3.KV(myEtcd.cli)
	myEtcd.option.Log.Info("putLease k:" + k + " v:" + v)
	putResponse, err = kv.Put(ctx, k, v, clientv3.WithLease(leaseId))
	//myEtcd.option.Log.Info("putLease (",leaseId,"): ",putResponse, err)
	if err != nil{
		return putResponse,err
	}

	return putResponse,nil
}
//根据前缀，获取该前缀下面的所有路径信息
func (myEtcd *MyEtcd)GetListByPrefix(key string)(list map[string]string,err error){
	myEtcd.option.Log.Info(" etcd GetListByPrefix key: " + key)
	rootContext := context.Background()
	kvc := clientv3.NewKV(myEtcd.cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	response, err := kvc.Get(ctx, key,clientv3.WithPrefix())
	defer cancelFunc()
	//myEtcd.option.Log.Debug(" ",response, err)
	if err != nil {
		myEtcd.option.Log.Warn("client Get err : " + err.Error())
		return list,errors.New("client Get err : "+err.Error())
	}

	if response.Count == 0{
		return list,nil
	}

	kvs := response.Kvs
	list = make(map[string]string)
	for _,v := range kvs{
		//MyPrint(string(v.Key),string(v.Value))
		list[string(v.Key)] =  string(v.Value)
	}
	//MyPrint(list)
	return list,nil
}

//func (myEtcd *MyEtcd)GetListValue(key string)(list []string){
//	myEtcd.option.Log.Info(" etcd GetOne , ",key ," : ")
//	rootContext := context.Background()
//	kvc := clientv3.NewKV(myEtcd.cli)
//	//获取值
//	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
//	response, err := kvc.Get(ctx, key)
//	myEtcd.option.Log.Debug(" ",response, err)
//	if err != nil {
//		myEtcd.option.Log.Error("Get",err)
//	}
//	cancelFunc()
//
//	if response.Count == 0{
//		return nil
//	}
//
//	kvs := response.Kvs
//
//	for _,v := range kvs{
//		list = append(list,string(v.Value))
//	}
//	return list
//}
//
//func (myEtcd *MyEtcd)GetOneValue(key string)string{
//	myEtcd.option.Log.Info(" etcd GetOne , ",key ," : ")
//	rootContext := context.Background()
//	kvc := clientv3.NewKV(myEtcd.cli)
//	//获取值
//	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
//	response, err := kvc.Get(ctx, key)
//	myEtcd.option.Log.Debug(" ",response, err)
//	if err != nil {
//		myEtcd.option.Log.Error("Get",err)
//	}
//	cancelFunc()
//
//	if response.Count == 0{
//		return ""
//	}
//
//	kvs := response.Kvs
//	value := string( kvs[0].Value )
//	return value
//}
//
//func (myEtcd *MyEtcd)SetLog(log *Log){
//	myEtcd.option.Log = log
//}

func (myEtcd *MyEtcd) PutOne(k string, v string)(putResponse *clientv3.PutResponse,errs error){
	myEtcd.option.Log.Info(" etcd PutOne: " + k  + v)
	rootContext := context.Background()
	kvc := clientv3.NewKV(myEtcd.cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	defer cancelFunc()
	putResponse, errs = kvc.Put(ctx, k,v)

	if errs != nil {
		myEtcd.option.Log.Error("RegOneService : " + errs.Error())
		switch errs {
		case context.Canceled:
			myEtcd.option.Log.Error("ctx is canceled by another routine: " + errs.Error())
		case context.DeadlineExceeded:
			myEtcd.option.Log.Error("ctx is attached with a deadline is exceeded: " + errs.Error())
		//case rpctypes.ErrEmptyKey:
		//	log.Error("client-side error: %v", err)
		default:
			myEtcd.option.Log.Error("bad cluster endpoints, which are not etcd servers: %v" +  errs.Error())
		}
	}
	myEtcd.option.Log.Info("RegOneService success" + putResponse.Header.String() + putResponse.PrevKv.String())
	return putResponse, errs
}

func  (myEtcd *MyEtcd)Watch(ctx context.Context,key string) <-chan clientv3.WatchResponse {
	myEtcd.option.Log.Warn("etcd create new watch :" + key)
	watchChan  := myEtcd.cli.Watch(ctx,key,clientv3.WithPrefix())
	//MyPrint("return watchChan")
	return watchChan
	//rch := cli.Watch(context.Background(), "/xi")
}

func  (myEtcd *MyEtcd)getConfRootPrefix()string{
	rootPath := "/config/"+myEtcd.option.AppKey + "/"+  myEtcd.option.AppENV + "/"
	return rootPath
}
//初始化，一个项目下的，所有配置文件（ 路径：/项目名/环境名/）
//减少请求etcd次数
func  (myEtcd *MyEtcd)iniAppConf() {
	myEtcd.option.Log.Info("etcd iniAppConf : ")
	confListEtcd,err := myEtcd.GetListByPrefix(myEtcd.getConfRootPrefix())
	if err != nil{
		myEtcd.option.Log.Error("iniAppConf err:" + err.Error())
		return
	}
	if len(confListEtcd) == 0{
		myEtcd.option.Log.Warn("iniAppConf confListEtcd is empty!")
		return
	}
	confList := make(map[string]string)
	for k,v := range confListEtcd{
		str := strings.Replace(k,myEtcd.getConfRootPrefix(),"",-1)
		//serviceArr := strings.Split(str,"/")
		myEtcd.option.Log.Info("conf "  + str + v )
		confList[str] = v
	}
	myEtcd.AppConflist = confList
}
func  (myEtcd *MyEtcd)DelOne(key string)error{
	myEtcd.option.Log.Info("myEtcd DelOne:" + key)
	_, err := myEtcd.cli.Delete(context.TODO(),key)
	if err != nil{
		myEtcd.option.Log.Error(" etcd del one err:" + err.Error())
	}
	return err
}

func  (myEtcd *MyEtcd)GetLinkAddressList()([]string){
	val := myEtcd.option.LinkAddressList
	return val
}

func  (myEtcd *MyEtcd)GetAppConf()(map[string]string){
	val := myEtcd.AppConflist
	return val
}

func  (myEtcd *MyEtcd)GetAppConfByKey(key string)(str string){
	val := myEtcd.AppConflist[key]
	return val
}
