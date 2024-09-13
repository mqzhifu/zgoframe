package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

// @Tags Gateway
// @Summary 网关 - 短连接
// @Description 通过网关调取后端服务(grpc)
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param service_name path string true "服务名"
// @Param func_name path string true "函数名"
// @Param data body string true "任意，请参考.proto"
// @Success 200 {boolean} true "true:成功 false:否"
// @Router /gateway/service/{service_name}/{func_name} [post]
func GatewayService(c *gin.Context) {
	prefix := "Gateway Http , "
	serviceName := c.Param("service_name")
	funcName := c.Param("func_name")

	data, err := c.GetRawData()
	util.MyPrint("c.GetRawData data:", data, " err:", err)
	// if err != nil {
	//	errMsg := prefix + " GetRawData err:" + err.Error()
	//	//util.ExitPrint(prefix + " GetRawData err:"+err.Error())
	//	httpresponse.FailWithMessage(errMsg, c)
	//	c.Abort()
	// }

	fmt.Println(prefix+" ServiceName:"+serviceName, " funcName:"+funcName+" data:"+string(data))
	backData, err := global.V.Service.Gateway.HttpCallGrpc(serviceName, funcName, "", data)
	if err != nil {
		fmt.Println(err)
		httpresponse.FailWithMessage(err.Error(), c)
		c.Abort()
	}

	fmt.Println(backData)
}

// @Tags Gateway
// @Summary 获取网关配置信息
// @Description 主要是长连接的配置(端口|协议)
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Success 200 {object} util.NetWayOption "dddd"
// @Router /gateway/config [get]
func GatewayConfig(c *gin.Context) {
	httpresponse.OkWithAll(global.V.Service.Gateway.NetWayOption, "ok", c)
	// httpresponse.OkWithAll(global.V.Gate.Option, "ok", c)

}

// @Tags Gateway
// @Summary 获取所有服务的.proto 配置文件
// @Description proto接口及GRPC微服务函数的信息等
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Success 200 {string} bbbb " "
// @Router /gateway/proto [get]
func GatewayProto(c *gin.Context) {
	url := "http:/127.0.0.1:" + global.C.Http.Port + "/" + global.C.Http.Status + "/proto"
	msg := "去 <a target='_blakn' href='" + url + "'" + "点我</a>"
	httpresponse.OkWithMessage(msg, c)
}

// @Tags Gateway
// @Summary php解析:.proto文件，生成.txt , 再通过GO读取出来
// @Description 后期考虑替换掉PHP解析过程，直接用GO
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Success 200 {string} ssss " "
// @Router /gateway/action/map [get]
func ActionMap(c *gin.Context) {
	list := global.V.Util.ProtoMap.GetServiceFuncMap()
	// 格式化数据，方便前端使用
	rsList := make(map[string]map[int]util.ProtoServiceFunc)
	clientList := make(map[int]util.ProtoServiceFunc)
	serverList := make(map[int]util.ProtoServiceFunc)
	for _, v := range list {
		cate := string([]byte(v.FuncName)[:2])
		if cate == "CS" {
			clientList[v.Id] = v
		} else if cate == "SC" {
			serverList[v.Id] = v
		}
	}
	rsList["client"] = clientList
	rsList["server"] = serverList

	// client
	// server

	httpresponse.OkWithAll(rsList, "ok", c)
}

// @Tags Gateway
// @Summary 网关 - 长连接 - metrics
// @Description metrics
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Router /gateway/total [get]
// @Success 200 {object} util.Conn "连接结构体"
func GatewayTotal(c *gin.Context) {
	myMetrics, _ := global.V.Service.Gateway.Netway.Metrics.GetAllByPrefix()
	httpresponse.OkWithAll(myMetrics, "ok", c)
	// mm, _ := prometheus.DefaultGatherer.Gather()
	// for _, v := range mm {
	//	util.MyPrint("mm==:", v.GetType())
	//	for _, b := range v.GetMetric() {
	//		util.MyPrint(b.String())
	//	}
	// }

	// util.MyPrint("prometheus.DefaultGatherer:======", prometheus.DefaultGatherer.)
}

// @Tags Gateway
// @Summary 网关 - 长连接
// @Description 长连接列表，FD => UID
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Router /gateway/fd/list [get]
// @Success 200 {object} util.Conn "连接结构体"
func GatewayFDList(c *gin.Context) {
	connManager := global.V.Service.Gateway.Netway.ConnManager
	if len(connManager.Pool) <= 0 {
		emptyMap := make(map[int32]*util.Conn)
		httpresponse.OkWithAll(emptyMap, "ok", c)
		return
	}

	connFDStrByte, err := json.Marshal(connManager.Pool)
	if err != nil {
		util.MyPrint("json.Marshal(connManager.Pool) err:", err)
	}
	connFDStr := string(connFDStrByte)
	util.MyPrint(connManager.Pool, connFDStr, err)

	httpresponse.OkWithAll(connManager.Pool, "ok", c)
}

// @Tags Gateway
// @Summary 网关 - 长连接 - 给某个用户-发送一条消息
// @Description 给某个UID发送一条消息，主要用于测试
// @Security ApiKeyAuth
// @Param X-Source-Type header string true "来源" Enums(11,12,21,22)
// @Param X-Project-Id header string true "项目ID" default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Param data body pb.ProjectPushMsg true "基础信息"
// @Router /gateway/send/msg [post]
// @Success 200 {string} bbbb " "
func GatewaySendMsg(c *gin.Context) {
	connManager := global.V.Service.Gateway.Netway.ConnManager
	// if len(connManager.Pool) <= 0 {
	//	msg := "失败，user pool = 0"
	//	httpresponse.FailWithMessage(msg, c)
	//	return
	// }

	var form pb.ProjectPushMsg
	c.ShouldBind(&form)
	if form.SourceUid <= 0 {
		httpresponse.FailWithMessage("SourceUid empty!!!", c)
		return
	}

	if form.SourceProjectId <= 0 {
		httpresponse.FailWithMessage("SourceProjectId empty!!!", c)
		return
	}

	if form.Content == "" {
		httpresponse.FailWithMessage("Content empty!!!", c)
		return
	}

	if form.TargetUids == "" {
		httpresponse.FailWithMessage("TargetUids empty!!!", c)
		return
	}

	formBytes, _ := json.Marshal(&form)
	protoMap, _ := global.V.Util.ProtoMap.GetServiceByName("Gateway", "SC_ProjectPushMsg")
	SifFId := global.V.Util.ProtoMap.GetIdBySidFid(protoMap.ServiceId, protoMap.FuncId)
	msg := pb.Msg{
		ServiceId:   int32(protoMap.ServiceId),
		FuncId:      int32(protoMap.FuncId),
		SidFid:      int32(SifFId),
		Content:     string(formBytes),
		ContentType: util.CONTENT_TYPE_JSON,
	}

	global.V.Service.Gateway.NativeServiceFuncRouter(msg)
	httpresponse.OkWithAll(connManager.Pool, "ok", c)
}
