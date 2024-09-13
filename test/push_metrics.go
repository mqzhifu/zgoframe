package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func PushMetrics() {
	testPushCounterName := "testPushCounter"

	global.V.Util.Metric.CreateCounter(testPushCounterName, "im_test_counter")
	global.V.Util.Metric.CounterInc(testPushCounterName)

	testPushGaugeName := "testPushGauge"

	global.V.Util.Metric.CreateGauge(testPushGaugeName, "im_test_gauge")
	global.V.Util.Metric.GaugeSet(testPushGaugeName, 0.001)

	push_err := global.V.Util.Metric.PushMetrics()
	util.MyPrint("test pusher err:", push_err)
	return
}
