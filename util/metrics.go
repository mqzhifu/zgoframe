package util

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
)

type MyMetrics struct {
	Groups map[string]interface{}
}

func NewMyMetrics()*MyMetrics{
	myMetrics := new(MyMetrics)
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
	if ok {
		return true
	}else{
		return false
	}
}
func (myMetrics *MyMetrics)CreateGauge(name string )error{
	if myMetrics.GroupNameHasExist(name) {
		return errors.New("GroupNameHasExist:"+name)
	}

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      name,
		//Help:      "the temperature of CPU",
	})

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
		return errors.New("GroupNameHasExist:"+name)
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
		return errors.New("GroupNameHasExist:"+name)
	}
	var AccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
		},
	)

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

//Counter end

