package initialize

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
	"zgoframe/core/global"
	"zgoframe/util"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"context"
)

type ViperOption struct {
	SourceType string
	ConfigFileType string
	ConfigFileName string
	EtcdUrl string
	ENV string
}
//读取配置文件：目前权支持文件，ETCD只写了一半
func GetNewViper(viperOption  ViperOption )(myViper *viper.Viper,config global.Config,err error){
	//util.MyPrint("SourceType:",viperOption.SourceType, " ConfigFileType:",viperOption.ConfigFileType ," , ConfigFileName:",viperOption.ConfigFileName)
	myViper = viper.New()


	if viperOption.SourceType == "file"{
		myViper.SetConfigType(viperOption.ConfigFileType)
		//myViper.SetConfigName(ConfigName + "." + ConfigType)
		configFile := viperOption.ConfigFileName + "." + viperOption.ConfigFileType
		//util.MyPrint(configFile)
		myViper.SetConfigFile(configFile)
		myViper.AddConfigPath(".")
		err = myViper.ReadInConfig()
		if err != nil{
			util.MyPrint("myViper.ReadInConfig() err :",err)
			return myViper,config,err
		}
		//config := Config{}
		err = myViper.Unmarshal(&config)
		if err != nil{
			util.MyPrint(" myViper.Unmarshal err:",err)
			return myViper,config,err
		}
	}else{
		util.MyPrint("get etcd config url:",viperOption.EtcdUrl)
		if viperOption.SourceType != "etcd"{
			util.MyPrint("configSourceType err: etcd or file")
			return  myViper,config,err
		}

		if viperOption.EtcdUrl == ""{
			util.MyPrint("viperOption.EtcdUrl == empty")
			return  myViper,config,err
		}

		jsonStruct ,errs := getEtcdHostPort(viperOption.EtcdUrl)
		util.MyPrint("getEtcdHostPort:",jsonStruct,errs)
		if errs != nil {
			return myViper,config,errors.New("http request err :" + errs.Error())
		}
		//etcdOption.Log.Info("etcdConfig ip list : ", jsonStruct.Data.Hosts)
		//linkAddressList := jsonStruct.Data.Hosts
		//开启建立连接
		clientv3Config  := clientv3.Config{
			Endpoints:  jsonStruct.Data.Hosts,
			DialTimeout: 5 * time.Second,
			Username: jsonStruct.Data.Username,
			Password: jsonStruct.Data.Password,
		}

		cli, errs := clientv3.New(clientv3Config)
		if errs != nil {
			return myViper,config,errors.New("clientv3.New error :  " + errs.Error())
		}
		prefix := "/gamematch/"+viperOption.ENV
		etcdList,err := GetListByPrefix(cli,prefix)
		util.ExitPrint(err,etcdList)
	}

	if config.Viper.Watch == global.CONFIG_STATUS_OPEN{
		util.MyPrint("viper watch open")
		myViper.WatchConfig()
		handleFunc := func(in fsnotify.Event) {
			util.MyPrint("myViper.WatchConfig onChange:",in.Name ,in.String())

			//if err := viper.Unmarshal(Conf); err != nil {
			//	panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
			//}
		}
		myViper.OnConfigChange(handleFunc)
		viper.OnConfigChange(handleFunc)
	}

	return myViper,config,nil
}

type Etcdconfig struct {
	Username	string `json:"username"`
	Password	string	`json:"password"`
	Hosts		[]string `json:"hosts"`
}

type EtcdHttpResp struct {
	Code int        `json:"code"`
	Data Etcdconfig `json:"data"`
}

func getEtcdHostPort(etcdUrl string)( etcdHttpResp EtcdHttpResp,err error){
	resp, errs := http.Get(etcdUrl)
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

//根据前缀，获取该前缀下面的所有路径信息
func GetListByPrefix(cli *clientv3.Client,key string)(list map[string]string,err error){
	util.MyPrint(" etcd GetListByPrefix , key "+ ":" + key)
	rootContext := context.Background()
	kvc := clientv3.NewKV(cli)
	//获取值
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	response, err := kvc.Get(ctx, key,clientv3.WithPrefix())
	defer cancelFunc()
	//myEtcd.option.Log.Debug(" ",response, err)
	if err != nil {
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