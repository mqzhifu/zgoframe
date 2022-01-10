package util

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

type Gateway struct {
	GrpcManager *GrpcManager
	Log 		*zap.Logger
}

func NewGateway(grpcManager *GrpcManager,log *zap.Logger )*Gateway{
	gateway := new(Gateway)
	gateway.GrpcManager = grpcManager
	gateway.Log = log
	return gateway
}
func  (gateway *Gateway)HttpCallGrpc(serviceName string,funcName string,balanceFactor string,requestData []byte)( resJsonStr string,err error){
	callGrpcResData ,err := gateway.CallGrpc(serviceName,funcName,balanceFactor,requestData)
	if err != nil{

	}
	resJsonStrByte ,err  := json.Marshal(callGrpcResData )
	if err != nil{

	}

	return string(resJsonStrByte),err
}

func  (gateway *Gateway)CallGrpc(serviceName string,funcName string,balanceFactor string,requestData []byte)( resData interface{},err error){
	//先确定service,根据service再确定执行哪个函数
	switch serviceName {
	case "Zgoframe":
		resData , err = gateway.GrpcManager.CallServiceFuncZgoframe(funcName,balanceFactor,requestData)
	}

	return resData,err
}

func (gateway *Gateway)StartSocket(netWayOption NetWayOption){
	NewNetWay(netWayOption)
	for  {
		time.Sleep(1)
	}
	//
	//roomId := "aabbccdd"
	//ZgoframeClient ,err := gateway.GrpcManager.GetZgoframeClient(roomId)
}
