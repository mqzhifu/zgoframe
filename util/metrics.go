package util

/*
	自实现服务的exporter
	依赖：client_golang/prometheus 库，其核心：
	1. 收集器
	2. 指定定义器
	3. 推送/接收器 (http)

	目前metric类型只实现两种：counter Gauge ，未实现：Histogram Summary
	数据多维度label 未实现
 */


import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"strconv"

	//"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

type PushGateway struct {
	Status 	string
	Ip 		string
	Port 	string
	JobName string
}

type MyMetricsOption struct{
	Log *zap.Logger
	PushGateway PushGateway
	NameSpace string
	Env int
}

type MyMetrics struct {
	Groups map[string]interface{}
	Pusher 	*push.Pusher
	Option MyMetricsOption
}

func NewMyMetrics(option MyMetricsOption)*MyMetrics{
	myMetrics := new(MyMetrics)
	myMetrics.Groups = make(map[string]interface{})
	//myMetrics.Log = log
	//myMetrics.PushGateway = pushGateway
	//myMetrics.NameSpace = nameSpace
	myMetrics.Option = option
	option.Log.Info("NewMyMetrics")

	if myMetrics.Option.PushGateway.Status == "open"{
		//dns := "http://"+pushGateway.Ip + ":" + pushGateway.Port + "/metrics"
		dns := "http://"+myMetrics.Option.PushGateway.Ip + ":" + myMetrics.Option.PushGateway.Port
		pusher := push.New(dns,myMetrics.Option.PushGateway.JobName)
		myMetrics.Pusher = pusher
		//testPushGateway()
	}


	return myMetrics
}

func (myMetrics *MyMetrics)GroupNameHasExist(name string)bool{
	_,ok := myMetrics.Groups[name]
	rs := false
	if ok {
		rs = true
	}
	//fmt.Println("GroupNameHasExist "+ name + " rs:",rs)
	return rs
}
func (myMetrics *MyMetrics)PushMetrics()error{
	if myMetrics.Option.PushGateway.Status != "open"{
		return errors.New("PushGateway.Status != open")
	}

	myMetrics.Pusher.Grouping("instance", myMetrics.Option.PushGateway.Ip ).Grouping("env",strconv.Itoa(myMetrics.Option.Env)).Push()

	return nil
}
func (myMetrics *MyMetrics)CreateGauge(name string,help string )error{
	if myMetrics.GroupNameHasExist(name) {
		return errors.New("CreateGauge GroupNameHasExist:"+name)
	}

	//labels :=  make(map[string]string)
	//labels["label_create_type"] = "CreateGauge"
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      name,
		//Namespace: myMetrics.Option.NameSpace,
		Help:     help,
		//ConstLabels:labels,
	})


	//myMetrics.Option.Log.Info("metrics: CreateGauge "+ name)

	prometheus.MustRegister(gauge)

	if myMetrics.Option.PushGateway.Status == "open"{
		myMetrics.Pusher.Collector(gauge)
	}

	myMetrics.Groups[name] = gauge

	return nil
}
//Gauge  start
func (myMetrics *MyMetrics)GaugeSet(name string,value float64 )error{
	if !myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupNameHasExist:"+name)
	}

	gauge := myMetrics.Groups[name].(prometheus.Gauge)
	gauge.Set(value)

	return nil
}

func (myMetrics *MyMetrics)GaugeInc(name string )error{
	if !myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupName not Exist:"+name)
	}

	gauge := myMetrics.Groups[name].(prometheus.Gauge)
	gauge.Inc()

	return nil
}

func (myMetrics *MyMetrics)GaugeDec(name string )error{
	if !myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupNameHasExist:"+name)
	}

	gauge := myMetrics.Groups[name].(prometheus.Gauge)
	gauge.Dec()

	return nil
}

func (myMetrics *MyMetrics)GaugeAdd(name string,value float64 )error{
	if !myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupNameHasExist:"+name)
	}

	gauge := myMetrics.Groups[name].(prometheus.Gauge)
	gauge.Add(value)

	return nil
}
//Gauge end

//Counter start
func (myMetrics *MyMetrics)CreateCounter(name string,help string )error{
	if myMetrics.GroupNameHasExist(name) {
		return errors.New("CreateCounter GroupNameHasExist:"+name)
	}
	var AccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			//Namespace: myMetrics.Option.NameSpace,
			Help: help,
		},
	)

	//myMetrics.Option.Log.Info("metrics: CreateCounter "+name )


	if myMetrics.Option.PushGateway.Status == "open"{
		myMetrics.Pusher.Collector(AccessCounter)
	}

	prometheus.MustRegister(AccessCounter)
	myMetrics.Groups[name] = AccessCounter

	return nil
}
func (myMetrics *MyMetrics)CounterInc(name string )error{
	if !myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupNameHasExist:"+name)
	}

	counter := myMetrics.Groups[name].(prometheus.Counter)
	counter.Inc()

	return nil
}

func (myMetrics *MyMetrics)CounterDec(name string,value float64 )error{
	if !myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupNameHasExist:"+name)
	}

	counter := myMetrics.Groups[name].(prometheus.Counter)
	counter.Add(value)

	return nil
}

func (myMetrics *MyMetrics)Shutdown(){

}

//Counter end


//func testPushGateway(pusher *push.Pusher,dns string){
//
//	myTimer := time.NewTimer(time.Second * 2)
//	MyPrint("start push metrics :"+ dns)
//	cnt := prometheus.NewCounter(prometheus.CounterOpts{
//		Name:      "pushName",
//		Namespace: "testNamespace",
//		Help:     "test golang pusher",
//	})
//	cnt.Inc()
//	<-myTimer.C
//	err := pusher.Collector(cnt).Grouping("instance", "1.1.1.1").Push()
//	MyPrint("push metrics err :",err)
//}
//
//func (myMetrics *MyMetrics)Test(){
//	myMetrics.CreateCounter("paySuccess")
//	myMetrics.CounterInc("paySuccess")
//	myMetrics.CounterInc("paySuccess")
//
//	myMetrics.CreateGauge("payUser")
//	myMetrics.GaugeSet("payUser",100)
//}