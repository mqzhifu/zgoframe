// 公共网关，主要是想把所有网络请求做统一的处理,尤其是：长连接
package gateway

import (
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/protobuf/pb"
	"zgoframe/service/bridge"
	"zgoframe/util"
)

////这是个快捷变量，目前所有代码均在一起，直接挂在这个变量上即可，后期所有服务分拆出去，网关没那么多附加功能此变量就没用了
//type MyServiceList struct {
//	//Match      *Match
//	GameMatch *gamematch.GameMatch
//	FrameSync *frame_sync.FrameSync
//	//RoomManage *frame_sync.RoomManager
//	TwinAgora *seed_business.TwinAgora
//}

type Gateway struct {
	GrpcManager  *util.GrpcManager //通过GRPC反射代理其它微服务
	Log          *zap.Logger       //日志
	Netway       *util.NetWay      //长连接公共类
	NetWayOption util.NetWayOption //长连接公共类的初始化参数
	//MyServiceList *MyServiceList    //快捷访问内部微服务
	//RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
	ServiceBridge *bridge.Bridge
}

/*
网关，目前主要是分为3个主要功能：
1. http 代理 grpc(中等)
2. 长连接代理(重点)
3. http 代理 http(鸡肋)
*/
func NewGateway(grpcManager *util.GrpcManager, log *zap.Logger, serviceBridge *bridge.Bridge) *Gateway {
	//func NewGateway(grpcManager *util.GrpcManager, log *zap.Logger, requestServiceAdapter *service.RequestServiceAdapter) *Gateway {
	gateway := new(Gateway)
	gateway.GrpcManager = grpcManager
	gateway.Log = log
	gateway.ServiceBridge = serviceBridge
	go gateway.ListeningBridgeMsg()
	//gateway.MyServiceList = &MyServiceList{}
	//gateway.RequestServiceAdapter = requestServiceAdapter
	//go gateway.ListeningMsg()
	return gateway
}

//func (gateway *Gateway) ListeningMsg() {
//	for {
//		select {
//		case GatewayMsg := <-gateway.RequestServiceAdapter.QueueGatewayMsg:
//			conn, exist := gateway.Netway.ConnManager.GetConnPoolById(GatewayMsg.Uid)
//			if !exist {
//				gateway.Log.Error("ListeningMsg conn empty uid:" + strconv.Itoa(int(GatewayMsg.Uid)))
//				break
//			}
//			conn.SendMsgCompressByUid(GatewayMsg.Uid, GatewayMsg.ActionName, GatewayMsg.Data)
//		case ServiceMsg := <-gateway.RequestServiceAdapter.QueueServiceMsg:
//			//工程太大不写了
//			util.MyPrint("gateway ListeningMsg:", ServiceMsg)
//		default:
//			time.Sleep(time.Millisecond * 50)
//		}
//	}
//
//}

// 开启长连接监听
func (gateway *Gateway) StartSocket(netWayOption util.NetWayOption) (*util.NetWay, error) {
	gateway.Log.Info("gateway StartSocket:")
	//netWayOption.RouterBack = gateway.Router //公共回调 路由器，用于给最底层的长连接公共类回调
	netWayOption.RouterBack = gateway.ServiceBridge.RouterBack
	//创建长连接:底层-公共类
	gateway.NetWayOption = netWayOption
	netWay, err := util.NewNetWay(netWayOption)
	gateway.Netway = netWay
	return netWay, err
}

// 广播给所有服务，如：心跳 PING PONG 关闭事件(不广播给gateway)
func (gateway *Gateway) BroadcastService(funcName string, msg pb.Msg) {
	gateway.Log.Debug("BroadcastService funcId:" + strconv.Itoa(int(msg.FuncId)))
	//gateway.RouterServiceSync(msg)
	////gateway.RouterServiceGameMatch(msg, conn)
	//gateway.RouterServiceTwinAgora(msg)

	//serviceDesc, empty := gateway.NetWayOption.ProtoMap.GetServiceByName("FrameSync", funcName)
	//if empty {
	//	util.ExitPrint("BroadcastService get service1 empty , name:" + funcName)
	//}
	//msg.ServiceId = int32(serviceDesc.ServiceId)
	//msg.FuncId = int32(serviceDesc.FuncId)
	//msg.SidFid = int32(gateway.NetWayOption.ProtoMap.GetIdBySidFid(serviceDesc.ServiceId, serviceDesc.FuncId))
	//
	//gateway.ServiceBridge.Call(service.CallMsg{Msg: msg})
	//
	//serviceDesc, _ = gateway.NetWayOption.ProtoMap.GetServiceByName("GameMatch", funcName)
	//if empty {
	//	util.ExitPrint("BroadcastService get service2 empty, name:" + funcName)
	//}
	//msg.ServiceId = int32(serviceDesc.ServiceId)
	//msg.FuncId = int32(serviceDesc.FuncId)
	//msg.SidFid = int32(gateway.NetWayOption.ProtoMap.GetIdBySidFid(serviceDesc.ServiceId, serviceDesc.FuncId))
	//gateway.ServiceBridge.Call(service.CallMsg{Msg: msg})

	serviceDesc, empty := gateway.NetWayOption.ProtoMap.GetServiceByName("TwinAgora", funcName)
	if empty {
		util.ExitPrint("BroadcastService get service3 empty, name:" + funcName)
	}
	msg.ServiceId = int32(serviceDesc.ServiceId)
	msg.FuncId = int32(serviceDesc.FuncId)
	msg.SidFid = int32(gateway.NetWayOption.ProtoMap.GetIdBySidFid(serviceDesc.ServiceId, serviceDesc.FuncId))
	gateway.ServiceBridge.Call(bridge.CallMsg{Msg: msg})

}

func (gateway *Gateway) ClientPong(requestClientPong pb.PongRes) {
	gateway.Log.Debug("ClientPong")
}

func (gateway *Gateway) heartbeat(requestClientHeartbeat pb.Heartbeat) {
	util.MyPrint("================", requestClientHeartbeat.SourceUid)
	conn, _ := gateway.Netway.ConnManager.GetConnPoolById(requestClientHeartbeat.SourceUid)

	now := util.GetNowTimeSecondToInt()
	now64 := util.GetNowMillisecond()
	//util.MyPrint("=================", now, now64)
	conn.UpTime = int32(now)
	conn.RTT = now64 - requestClientHeartbeat.Time

	gateway.Log.Debug("gateway heartbeat , now64:", zap.Int64("now", now64), zap.Int64(" client_time", requestClientHeartbeat.Time), zap.Int64(" RTT:", conn.RTT))
	responseHeartbeat := pb.Heartbeat{
		Time:              now64,
		ReqTime:           requestClientHeartbeat.ClientReqTime,
		ClientReqTime:     requestClientHeartbeat.ClientReqTime,
		ServerReceiveTime: now64,
		RequestId:         requestClientHeartbeat.RequestId,
	}

	conn.SendMsgCompressByName("Gateway", "SC_Heartbeat", &responseHeartbeat)
}

func (gateway *Gateway) clientPing(ping pb.PingReq) {
	conn, exist := gateway.Netway.ConnManager.GetConnPoolById(ping.SourceUid)
	if !exist {
		gateway.Log.Error("clientPing conn empty uid:" + strconv.Itoa(int(ping.SourceUid)))
		return
	}
	responseServerPong := pb.PongRes{
		ClientReqTime:      ping.ClientReqTime,
		ClientReceiveTime:  ping.ClientReceiveTime,
		ServerReceiveTime:  util.GetNowMillisecond(),
		ServerResponseTime: util.GetNowMillisecond(),
		RequestId:          "AAA",
	}
	//gateway.RequestServiceAdapter.GatewaySendMsgByUid(ping.SourceUid, "SC_Pong", &responseServerPong)
	conn.SendMsgCompressByName("Gateway", "SC_Pong", &responseServerPong)
}

// balanceFactor:负载均衡 方法
func (gateway *Gateway) HttpCallGrpc(serviceName string, funcName string, balanceFactor string, requestData []byte) (resJsonStr string, err error) {
	gateway.Log.Info("HttpCallGrpc ， serviceName:" + serviceName + " funcName:" + funcName + " balanceFactor:" + balanceFactor + " requestData:" + string(requestData))
	callGrpcResData, err := gateway.GrpcManager.CallGrpc(serviceName, funcName, balanceFactor, requestData)
	if err != nil {
		return resJsonStr, err
	}
	resJsonStrByte, err := json.Marshal(callGrpcResData)
	if err != nil {
		return resJsonStr, err
	}
	return string(resJsonStrByte), err
}
