package service

/*
网关->调用 后端服务 或 后端服务调用 -> 网关

这种调用的方式分成了两大类
	1. 纯内部进程函数调用，其实也就是后端服务与网关未分离，都是一个程序
	2. 网络调用
		(1) http
		(2) grpc
*/

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"zgoframe/util"
)

type GatewayMsg struct {
	Uid        int32
	ActionName string
	Data       interface{}
}

type ServiceMsg struct {
	ServiceName string
	FuncName    string
	RequestData interface{}
}

//向3方服务发送请求的适配器，属于快捷方法，也可以在自己的服务中封装
type RequestServiceAdapter struct {
	ProjectId        int
	ServiceDiscovery *util.ServiceDiscovery
	GrpcManager      *util.GrpcManager
	Flag             int
	Log              *zap.Logger
	QueueGatewayMsg  chan GatewayMsg
	QueueServiceMsg  chan ServiceMsg
}

func NewRequestServiceAdapter(ServiceDiscovery *util.ServiceDiscovery, grpcManager *util.GrpcManager, flag int, projectId int, log *zap.Logger) *RequestServiceAdapter {
	requestService := new(RequestServiceAdapter)
	requestService.ServiceDiscovery = ServiceDiscovery
	requestService.GrpcManager = grpcManager
	requestService.Flag = flag
	requestService.Log = log
	requestService.QueueGatewayMsg = make(chan GatewayMsg, 100)
	requestService.QueueServiceMsg = make(chan ServiceMsg, 100)
	//requestService.Gateway = gateway
	return requestService
}

//网关内部调用服务
func (requestService *RequestServiceAdapter) GatewaySendMsgByUids(uids string, funcName string, requestData interface{}) {
	requestService.Log.Debug("RequestServiceAdapter GatewaySendMsgByUids :" + uids + " funcName:" + funcName)
	uidsArr := strings.Split(uids, ",")
	for _, uidStr := range uidsArr {
		uid, _ := strconv.Atoi(uidStr)
		n := GatewayMsg{
			Uid:        int32(uid),
			ActionName: funcName,
			Data:       requestData,
		}
		requestService.QueueGatewayMsg <- n
	}
}

//网关内部调用服务
func (requestService *RequestServiceAdapter) GatewaySendMsgByUid(uid int32, funcName string, requestData interface{}) {
	n := GatewayMsg{
		Uid:        uid,
		ActionName: funcName,
		Data:       requestData,
	}
	requestService.QueueGatewayMsg <- n
}

//网络远程调用
func (requestService *RequestServiceAdapter) RemoteCall(serviceName string, funcName string, balanceFactor string, requestData interface{}, httpUri string) (interface{}, error) {
	switch requestService.Flag {
	case REQ_SERVICE_METHOD_HTTP:
		myService, _ := requestService.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName, "")
		http := util.NewServiceHttp(requestService.ProjectId, myService.ServiceName, myService.Ip, myService.Port, myService.ServiceId)
		http.Post(httpUri, requestData)
	case REQ_SERVICE_METHOD_GRPC:
		requestService.GrpcManager.CallGrpc(serviceName, funcName, balanceFactor, requestData.([]byte))
	case REQ_SERVICE_METHOD_NATIVE:
		serviceMsg := ServiceMsg{
			ServiceName: serviceName,
			FuncName:    funcName,
			RequestData: requestData,
		}
		requestService.QueueServiceMsg <- serviceMsg
	default:
		return "", errors.New("not found type id.")
	}
	return "", nil
}
