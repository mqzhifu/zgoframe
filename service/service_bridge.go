package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

const (
	BRIDGE_SLEEP_TIME = 10   //内部调用时，每个服务要监听 管道里的消息，睡眠时间
	GATEWAY_ADMIN_UID = 9999 //后端反向给前端推送消息时，最好加上一个来源UID
)

type NativeServiceList struct {
	GameMatch chan pb.Msg
	FrameSync chan pb.Msg
	Gateway   chan pb.Msg
	TwinAgora chan pb.Msg
}

type BridgeOption struct {
	ProtoMap         *util.ProtoMap
	ProjectId        int
	ServiceDiscovery *util.ServiceDiscovery
	GrpcManager      *util.GrpcManager
	Flag             int
	Log              *zap.Logger
}

//向3方服务发送请求的适配器，属于快捷方法，也可以在自己的服务中封装
type Bridge struct {
	Op                BridgeOption
	NativeServiceList *NativeServiceList
}

func NewBridge(op BridgeOption) (*Bridge, error) {
	bridge := new(Bridge)
	bridge.Op = op
	if op.Flag <= 0 {
		return bridge, errors.New("op flag <=0 ")
	}

	nativeServiceList := NativeServiceList{
		GameMatch: make(chan pb.Msg, 100),
		FrameSync: make(chan pb.Msg, 100),
		Gateway:   make(chan pb.Msg, 100),
		TwinAgora: make(chan pb.Msg, 100),
	}

	bridge.NativeServiceList = &nativeServiceList
	return bridge, nil
}

//网关的呼叫略有点麻烦，多了一个 targetUid
func (bridge *Bridge) CallGateway(sourceServiceName string, sourceFunName string, sourceUid int32, targetUid int32, data string, balanceFactor string, flag int) (resData interface{}, err error) {
	serviceMapInfo, empty := bridge.Op.ProtoMap.GetServiceByName("Gateway", "SC_SendMsg")
	if empty {

	}

	sourceServiceMapInfp, empty := bridge.Op.ProtoMap.GetServiceByName(sourceServiceName, sourceFunName)
	msg := pb.Msg{
		SourceUid:       sourceUid,
		ServiceId:       int32(serviceMapInfo.ServiceId),
		FuncId:          int32(serviceMapInfo.FuncId),
		Content:         data,
		TargetUid:       targetUid,
		SourceServiceId: int32(sourceServiceMapInfp.ServiceId),
		SourceFuncId:    int32(sourceServiceMapInfp.FuncId),
	}
	return bridge.Call(msg, balanceFactor, flag)
}

func (bridge *Bridge) CallByName(serviceName, funcName string, data string, balanceFactor string, flag int) (resData interface{}, err error) {
	util.MyPrint("bridge CallByName , serviceNameL:"+serviceName+" , funcName:", funcName)
	serviceMapInfo, empty := bridge.Op.ProtoMap.GetServiceByName(serviceName, funcName)
	if empty {
		util.ExitPrint("bridge CallByName empty!")
	}

	msg := pb.Msg{
		ServiceId: int32(serviceMapInfo.ServiceId),
		FuncId:    int32(serviceMapInfo.FuncId),
		Content:   data,
	}
	return bridge.Call(msg, balanceFactor, flag)
}

//动态调用一个GRPC-SERVER 的一个方法 ( c端使用 )
func (bridge *Bridge) Call(msg pb.Msg, balanceFactor string, flag int) (resData interface{}, err error) {
	bridge.Op.Log.Debug("Bridge Call , flag:" + strconv.Itoa(flag) + " op.Flag:" + strconv.Itoa(bridge.Op.Flag) + " msg.ServiceId: " + strconv.Itoa(int(msg.ServiceId)) + " msg.FuncId: " + strconv.Itoa(int(msg.FuncId)))
	if flag <= 0 {
		flag = bridge.Op.Flag
	}

	switch flag {
	case REQ_SERVICE_METHOD_NATIVE:
		bridge.CallNativeService(msg, balanceFactor, flag)
	case REQ_SERVICE_METHOD_HTTP:
	case REQ_SERVICE_METHOD_GRPC:
		bridge.CallRemoteService(msg, balanceFactor, flag)
	default:
		util.MyPrint("bridge Call flag err:" + strconv.Itoa(flag))
		return "", errors.New("flag err")
	}
	return "", nil

}

func (bridge *Bridge) CallNativeService(msg pb.Msg, balanceFactor string, flag int) (resData interface{}, err error) {
	serviceFuncDesc, empty := bridge.Op.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	if empty {
		util.ExitPrint("CallNativeService GetServiceFuncById empty , SidFid:" + strconv.Itoa(int(msg.SidFid)))
	}
	bridge.Op.Log.Debug("CallNativeService msg.ServiceName: " + serviceFuncDesc.ServiceName + " FuncName: " + serviceFuncDesc.FuncName)

	switch serviceFuncDesc.ServiceName {
	case "GameMatch":
		bridge.NativeServiceList.GameMatch <- msg
	case "Gateway":
		bridge.NativeServiceList.Gateway <- msg
	case "FrameSync":
		bridge.NativeServiceList.FrameSync <- msg
	case "TwinAgora":
		bridge.NativeServiceList.TwinAgora <- msg
	default:
		util.ExitPrint("CallNativeService router err.")
	}

	return "insert chan Queue......", nil
}

func (bridge *Bridge) CallRemoteService(msg pb.Msg, balanceFactor string, flag int) (resData interface{}, err error) {
	serviceFuncDesc, empty := bridge.Op.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	if empty {

	}

	switch serviceFuncDesc.ServiceName {
	//case "FrameSync":
	//	resData , err = serviceCallBridge.CallServiceFuncFrameSync(funcName,balanceFactor,requestData)
	case "GameMatch":
		resData, err = bridge.CallServiceFuncGameMatch(msg, balanceFactor, flag)
	//case "Gateway":
	//	resData , err = serviceCallBridge.CallServiceFuncGateway(funcName,balanceFactor,requestData)
	//case "TwinAgora":
	//	resData , err = serviceCallBridge.CallServiceFuncTwinAgora(funcName,balanceFactor,requestData)

	default:
		return resData, errors.New("service name router failed.")
	}
	return resData, err
}

//动态调用服务的函数 : GameMatch
func (bridge *Bridge) CallServiceFuncGameMatch(msg pb.Msg, balanceFactor string, flag int) (data interface{}, err error) {
	//获取GRPC一个连接
	//grpcClient, err := bridge.Op.GrpcManager.GetGameMatchClient("GameMatch", balanceFactor)
	//if err != nil {
	//	return data, err
	//}

	serviceFuncDesc, empty := bridge.Op.ProtoMap.GetServiceFuncById(int(msg.SidFid))
	if empty {

	}

	ctx := context.Background()
	switch serviceFuncDesc.FuncName {
	case "CS_PlayerMatchSign":
		request := pb.GameMatchSign{}
		//err := json.Unmarshal(postData, &request)
		//if err != nil {
		//	return data, err
		//}
		//
		switch flag {
		case 1:
		case 2:
			grpcClient, err := bridge.Op.GrpcManager.GetGameMatchClient("GameMatch", balanceFactor)
			if err != nil {
				return data, err
			}
			data, err = grpcClient.CS_PlayerMatchSign(ctx, &request)
		case 3:
			//bridge.Op.NativeServiceList.GameMatch.PlayerJoin(msg.Content)
		}

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}
