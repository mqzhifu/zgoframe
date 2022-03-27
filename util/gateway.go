package util

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type Gateway struct {
	GrpcManager  *GrpcManager
	Log          *zap.Logger
	NetWayOption NetWayOption
}

func NewGateway(grpcManager *GrpcManager, log *zap.Logger) *Gateway {
	gateway := new(Gateway)
	gateway.GrpcManager = grpcManager
	gateway.Log = log
	return gateway
}

func (gateway *Gateway) HttpCallGrpc(serviceName string, funcName string, balanceFactor string, requestData []byte) (resJsonStr string, err error) {
	fmt.Print("HttpCallGrpc :", serviceName, funcName, balanceFactor, requestData)
	//gateway.Log.Info("HttpCallGrpc:")
	callGrpcResData, err := gateway.GrpcManager.CallGrpc(serviceName, funcName, balanceFactor, requestData)
	if err != nil {
		return resJsonStr, err
	}
	resJsonStrByte, err := json.Marshal(callGrpcResData)
	if err != nil {
		return resJsonStr, err
	}
	return string(resJsonStrByte), err
	//return resJsonStr,err
}

func (gateway *Gateway) StartSocket(netWayOption NetWayOption) {
	gateway.NetWayOption = netWayOption
	netWay, err := NewNetWay(netWayOption)
	if err != nil {
		errMsg := "NewNetWay err:" + err.Error()
		ExitPrint(errMsg)
	}

	for {
		time.Sleep(time.Second * 1)
	}
	netWay.Shutdown()
	//
	//roomId := "aabbccdd"
	//ZgoframeClient ,err := gateway.GrpcManager.GetZgoframeClient(roomId)

}
