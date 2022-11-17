package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ServiceHttp struct {
	//AppId 				int
	//ServiceId 			int
	ProjectId         int
	ServiceName       string
	Ip                string
	Port              string
	TargetServiceId   int
	TargetServiceName string
}

func NewServiceHttp(projectId int, targetServiceName string, ip string, port string, targetServiceId int) *ServiceHttp {
	serviceHttp := new(ServiceHttp)
	serviceHttp.ProjectId = projectId
	serviceHttp.Ip = ip
	serviceHttp.Port = port
	serviceHttp.TargetServiceName = targetServiceName
	serviceHttp.TargetServiceId = targetServiceId

	return serviceHttp
}

func (serviceHttp *ServiceHttp) GetDns() string {
	return "http://" + serviceHttp.Ip + ":" + serviceHttp.Port
}

func (serviceHttp *ServiceHttp) PostGateway(serviceId int, funcId int, targetUids string, data interface{}) (responseMsgST ResponseMsgST, errs error) {
	uri := "/gateway/send/msg"
	return serviceHttp.Post(uri, data)
}

//给一个服务，发送一条http消息
func (serviceHttp *ServiceHttp) Post(uri string, data interface{}) (responseMsgST ResponseMsgST, errs error) {
	//node,err := serviceDiscovery.GetLoadBalanceServiceNodeByServiceName(serviceName)
	//if err != nil{
	//	return responseMsgST,err
	//}
	//serviceHost := node.Ip + ":" + node.Port
	//url := "http://"+serviceHost + "/" + uri
	url := serviceHttp.GetDns() + "/" + uri
	//serviceDiscovery.option.Log.Info("HttpPost" + serviceName + serviceHost + uri + url)
	jsonStr, _ := json.Marshal(data)
	//serviceDiscovery.option.Log.Info("jsonStr:" + string(jsonStr))
	//ExitPrint(1111)
	req, errs := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if errs != nil {
		return responseMsgST, errors.New("NewRequest err")
	}

	req.Header.Add("content-type", "application/json")

	clientHeader := NewServiceClientHeader()
	//clientHeader.ProjectId = strconv.Itoa(serviceHttp.AppId)
	//clientHeader.ServiceId = strconv.Itoa(serviceHttp.ServiceId)
	clientHeader.ProjectId = strconv.Itoa(serviceHttp.ProjectId)
	clientHeader.TargetServiceName = serviceHttp.TargetServiceName

	clientHeaderStr, _ := json.Marshal(clientHeader)
	//clientHeader.AppId
	req.Header.Add(SERVICE_HEADER_KEY, string(clientHeaderStr))

	defer req.Body.Close()
	//5秒超时
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	//service.option.Log.Debug(resp,error)
	if error != nil {
		return responseMsgST, errors.New("client.Do  err" + error.Error())
	}

	if resp.StatusCode != 200 {
		return responseMsgST, errors.New("http response code != 200")
	}

	if resp.ContentLength == 0 {
		return responseMsgST, errors.New("http response content = 0")
	}
	contentJsonStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseMsgST, errors.New("ioutil.ReadAll err : " + err.Error())
	}

	errs = json.Unmarshal(contentJsonStr, &responseMsgST)
	if errs != nil {
		return responseMsgST, errors.New(" json.Unmarshal html content err : " + err.Error())
	}

	//service.option.Log.Debug("responseMsgST : ",responseMsgST)
	return responseMsgST, nil
}
