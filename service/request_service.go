package service

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

func (requestService *RequestServiceAdapter) GatewaySendMsgByUid(uid int32, funcName string, requestData interface{}) {
	n := GatewayMsg{
		Uid:        uid,
		ActionName: funcName,
		Data:       requestData,
	}
	requestService.QueueGatewayMsg <- n
}

func (requestService *RequestServiceAdapter) RemoteCall(serviceName string, funcName string, balanceFactor string, requestData interface{}, httpUri string) (interface{}, error) {
	switch requestService.Flag {
	case REQ_SERVICE_METHOD_HTTP:
		myService, _ := requestService.ServiceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName, "")
		http := util.NewServiceHttp(requestService.ProjectId, myService.ServiceName, myService.Ip, myService.Port, myService.ServiceId)
		http.Post(httpUri, requestData)
	case REQ_SERVICE_METHOD_GRPC:
		requestService.GrpcManager.CallGrpc(serviceName, funcName, balanceFactor, requestData.([]byte))
	case REQ_SERVICE_METHOD_INNER:
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
