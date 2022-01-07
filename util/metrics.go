package util

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type MyMetrics struct {
	Groups map[string]interface{}
	Log *zap.Logger
}

func NewMyMetrics(log *zap.Logger)*MyMetrics{
	myMetrics := new(MyMetrics)
	myMetrics.Groups = make(map[string]interface{})
	myMetrics.Log = log

	log.Info("NewMyMetrics")

	return myMetrics
}

func (myMetrics *MyMetrics)Test(){
	myMetrics.CreateCounter("paySuccess")
	myMetrics.CounterInc("paySuccess")
	myMetrics.CounterInc("paySuccess")

	myMetrics.CreateGauge("payUser")
	myMetrics.GaugeSet("payUser",100)
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

func (myMetrics *MyMetrics)CreateGauge(name string )error{
	if myMetrics.GroupNameHasExist(name) {
		return errors.New("CreateGauge GroupNameHasExist:"+name)
	}
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      name,
		//Help:      "the temperature of CPU",
	})

	myMetrics.Log.Info("metrics: CreateGauge "+ name)

	prometheus.MustRegister(gauge)
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
func (myMetrics *MyMetrics)CreateCounter(name string )error{
	if myMetrics.GroupNameHasExist(name) {
		return errors.New("CreateCounter GroupNameHasExist:"+name)
	}
	var AccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
		},
	)

	myMetrics.Log.Info("metrics: CreateCounter "+name )

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

