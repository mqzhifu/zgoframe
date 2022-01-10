package util

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
)

const (
	//协议
	SERVICE_PROTOCOL_HTTP = 1
	SERVICE_PROTOCOL_GRPC = 2
	SERVICE_PROTOCOL_WEBSOCKET = 3
	SERVICE_PROTOCOL_TCP = 4

	//服务发现 - 分布式DB
	SERVICE_DISCOVERY_ETCD = 1
	SERVICE_DISCOVERY_CONSUL = 2

	//公共头信息的  key
	SERVICE_HEADER_KEY = "service_diy"

	//服务发现的 负载均衡 类型
	LOAD_BALANCE_ROBIN = 1
	LOAD_BALANCE_HASH = 2
)
//公共请求头
type ServiceClientHeader struct {//这里全用的string 方便 http header 处理
	TraceId 			string	`json:"trace_id"`
	RequestId 			string	`json:"request_id"`
	Protocol 			string	`json:"protocol"`
	ProjectId 			string	`json:"project_id"`
	ServiceId			string 	`json:"service_id"`
	RequestTime 		string	`json:"request_time"`
	TargetServiceName 	string	`json:"target_service_name"`
	ServerReceiveTime 	string	`json:"server_receive_time"`
	ServerResponseTime 	string	`json:"server_response_time"`
}
//公共响应头
type ServiceServerHeader struct {//这里全用的string 方便 http header 处理
	TraceId 		string	`json:"trace_id"`
	RequestId 		string	`json:"request_id"`
	Protocol 		string	`json:"protocol"`
	ProjectId 			string	`json:"project_id"`
	TargetProjectId 	string	`json:"target_project_id"`
	ServiceName 	string	`json:"service_name"`
	ReceiveTime 	string	`json:"receive_time"`
	ResponseTime 	string	`json:"response_time"`
}

func GetServiceProtocolList()[]int{
	list := []int{SERVICE_PROTOCOL_HTTP, SERVICE_PROTOCOL_GRPC, SERVICE_PROTOCOL_WEBSOCKET, SERVICE_PROTOCOL_TCP}
	return list
}

func CheckServiceProtocolExist(protocol int)bool{
	list := GetServiceProtocolList()
	for _,v :=range list{
		if v == protocol{
			return true
		}
	}
	return false
}

func GetServiceDiscoveryList()[]int{
	list := []int{SERVICE_DISCOVERY_ETCD,SERVICE_DISCOVERY_CONSUL}
	return list
}

func CheckServiceDiscoveryExist(protocol int)bool{
	list := GetServiceDiscoveryList()
	for _,v :=range list{
		if v == protocol{
			return true
		}
	}
	return false
}

func MakeTraceId()string{
	id :=  uuid.NewV4()
	return id.String()
}

func MakeRequestId()string{
	id :=  uuid.NewV4()
	return id.String()
}

func NewServiceClientHeader()ServiceClientHeader{
	clientHeader := ServiceClientHeader{
		TraceId: MakeTraceId(),
		RequestTime: strconv.FormatInt(GetNowTimeSecondToInt64(),10),
		Protocol: "http",
		RequestId: MakeRequestId(),
		//TargetServiceName:  myGrpcClient.ServiceName,
		//AppId: strconv.Itoa(myGrpcClient.AppId),
	}
	return clientHeader
}
