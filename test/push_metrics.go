package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func PushMetrics(){
	testPushCounterName := "testPushCounter"

	global.V.Metric.CreateCounter(testPushCounterName,"im_test_counter")
	global.V.Metric.CounterInc(testPushCounterName)

	testPushGaugeName := "testPushGauge"

	global.V.Metric.CreateGauge(testPushGaugeName,"im_test_gauge")
	global.V.Metric.GaugeSet(testPushGaugeName,0.001)


	push_err := global.V.Metric.PushMetrics()
	util.MyPrint("test pusher err:",push_err)
	return
}
