package service

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"zgoframe/util"
)

type Gateway struct {
	GrpcManager  *util.GrpcManager
	Log          *zap.Logger
	NetWayOption util.NetWayOption
}

//网关，目前主要是分为2部分
//1. http 代理 grpc
//2. 长连接代理，这里才是重点
func NewGateway(grpcManager *util.GrpcManager, log *zap.Logger) *Gateway {
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

func (gateway *Gateway) StartSocket(netWayOption util.NetWayOption) (*util.NetWay, error) {
	gateway.NetWayOption = netWayOption
	netWay, err := util.NewNetWay(netWayOption)
	return netWay, err
	//if err != nil {
	//	//errMsg := "NewNetWay err:" + err.Error()
	//	return netWay, err
	//}
	//for {
	//	time.Sleep(time.Second * 1)
	//}
	//netWay.Shutdown()
	//
	//roomId := "aabbccdd"
	//ZgoframeClient ,err := gateway.GrpcManager.GetZgoframeClient(roomId)

}
