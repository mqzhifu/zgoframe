package bridge

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
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

// 向3方服务发送请求的适配器，属于快捷方法，也可以在自己的服务中封装
type Bridge struct {
	Op                BridgeOption
	Prefix            string //调试使用，输出前缀字符串值
	NativeServiceList *NativeServiceList
}

func NewBridge(op BridgeOption) (*Bridge, error) {
	bridge := new(Bridge)
	bridge.Prefix = "ServiceBridge "
	bridge.Op = op
	if op.Flag <= 0 {
		err := bridge.MakeError(" NewBridge " + "op flag <=0")
		return bridge, err
	}

	nativeServiceList := NativeServiceList{
		GameMatch: make(chan pb.Msg, 100),
		FrameSync: make(chan pb.Msg, 100),
		Gateway:   make(chan pb.Msg, 100),
		TwinAgora: make(chan pb.Msg, 10),
	}

	bridge.NativeServiceList = &nativeServiceList

	bridge.DebugInfo("NewBridge success ,flag:" + strconv.Itoa(op.Flag))

	return bridge, nil
}

type CallGatewayMsg struct {
	ServiceName string
	FunName     string
	SourceUid   int32
	TargetUid   int32
	Data        interface{}
	CallMsg
}

type CallMsg struct {
	pb.Msg
	BalanceFactor string
	Flag          int
}

// 网关的呼叫略有点麻烦，多了一个 targetUid
func (bridge *Bridge) CallGateway(callGatewayMsg CallGatewayMsg) (resData interface{}, err error) {
	debugInfo := " CallGateway ,  ServiceName:" + callGatewayMsg.ServiceName + " FunName:" + callGatewayMsg.FunName + " sourceUid:" + strconv.Itoa(int(callGatewayMsg.SourceUid)) + " targetUid: " + strconv.Itoa(int(callGatewayMsg.TargetUid))
	bridge.DebugInfo(debugInfo)

	if callGatewayMsg.TargetUid <= 0 {
		info := debugInfo + " , CallGateway  callGatewayMsg.TargetUid -1 empty!"
		err = bridge.MakeError(info)
		return
	}

	serviceMapInfo, empty := bridge.Op.ProtoMap.GetServiceByName("Gateway", "SC_SendMsg")
	if empty {
		info := debugInfo + " , ProtoMap.GetServiceByName -2 empty!"
		err = bridge.MakeError(info)
		return
	}

	sourceServiceMapInfo, empty := bridge.Op.ProtoMap.GetServiceByName(callGatewayMsg.ServiceName, callGatewayMsg.FunName)
	if empty {
		info := debugInfo + " , GetServiceByName - 3 empty!"
		err = bridge.MakeError(info)
		return
	}
	contentStruct := callGatewayMsg.Data.(proto.Message)
	content, err := proto.Marshal(contentStruct)

	if callGatewayMsg.SourceUid <= 0 {
		callGatewayMsg.SourceUid = GATEWAY_ADMIN_UID
	}

	msg := pb.Msg{
		SourceUid:       callGatewayMsg.SourceUid,
		ServiceId:       int32(serviceMapInfo.ServiceId),
		FuncId:          int32(serviceMapInfo.FuncId),
		SidFid:          int32(serviceMapInfo.Id),
		Content:         string(content),
		TargetUid:       callGatewayMsg.TargetUid,
		SourceServiceId: int32(sourceServiceMapInfo.ServiceId),
		SourceFuncId:    int32(sourceServiceMapInfo.FuncId),
		//SourceContent:   data,
	}
	callMsg := CallMsg{}
	callMsg.Msg = msg
	callMsg.Flag = callGatewayMsg.Flag
	callMsg.BalanceFactor = callGatewayMsg.BalanceFactor
	return bridge.Call(callMsg)
}

// 根据 名称 调用 call
func (bridge *Bridge) CallByName(callGatewayMsg CallGatewayMsg) (resData interface{}, err error) {
	debugInfo := " CallByName , serviceName::" + callGatewayMsg.ServiceName + " , funcName:" + callGatewayMsg.FunName
	bridge.DebugInfo(debugInfo)
	serviceMapInfo, empty := bridge.Op.ProtoMap.GetServiceByName(callGatewayMsg.ServiceName, callGatewayMsg.FunName)
	if empty {
		debugInfo += " GetServiceByName empty!!!"
		err = bridge.MakeError(debugInfo)
		return resData, err
	}

	msg := pb.Msg{
		ServiceId: int32(serviceMapInfo.ServiceId),
		FuncId:    int32(serviceMapInfo.FuncId),
		SidFid:    int32(serviceMapInfo.Id),
		Content:   callGatewayMsg.Content,
	}

	callMsg := CallMsg{}
	callMsg.Msg = msg
	callMsg.Flag = callGatewayMsg.Flag
	callMsg.BalanceFactor = callGatewayMsg.BalanceFactor
	return bridge.Call(callMsg)
}

func (bridge *Bridge) RouterBack(msg pb.Msg, balanceFactor string, flag int) (data interface{}, err error) {
	callMsg := CallMsg{
		Msg:           msg,
		BalanceFactor: balanceFactor,
		Flag:          flag,
	}
	bridge.Call(callMsg)
	return data, nil
}

// 核心方法： 动态调用一个服务(方法)
func (bridge *Bridge) Call(callMsg CallMsg) (resData interface{}, err error) {
	//func (bridge *Bridge) Call(msg pb.Msg, balanceFactor string, flag int) (resData interface{}, err error) {
	debugInfo := "Call , flag:" + strconv.Itoa(callMsg.Flag) + " op.Flag:" + strconv.Itoa(bridge.Op.Flag) + " msg.ServiceId: " + strconv.Itoa(int(callMsg.ServiceId)) + " msg.FuncId: " + strconv.Itoa(int(callMsg.FuncId)) + " msg.SourceUid:" + strconv.Itoa(int(callMsg.SourceUid))
	bridge.DebugInfo(debugInfo)
	if callMsg.Flag <= 0 {
		bridge.DebugInfo(" use class default flag.")
		callMsg.Flag = bridge.Op.Flag
	}

	switch callMsg.Flag {
	case REQ_SERVICE_METHOD_NATIVE:
		bridge.CallNativeService(callMsg)
	case REQ_SERVICE_METHOD_HTTP:
	case REQ_SERVICE_METHOD_GRPC:
		bridge.CallRemoteService(callMsg)
	default:
		info := "Call , switch  flag not found :" + strconv.Itoa(callMsg.Flag)
		err = bridge.MakeError(info)
		return "", err
	}
	return "", nil

}

// 动用本地方法
func (bridge *Bridge) CallNativeService(callMsg CallMsg) (resData interface{}, err error) {
	debugInfo := "CallNativeService , flag:" + strconv.Itoa(callMsg.Flag) + " op.Flag:" + strconv.Itoa(bridge.Op.Flag) + " msg.ServiceId: " + strconv.Itoa(int(callMsg.ServiceId)) + " msg.FuncId: " + strconv.Itoa(int(callMsg.FuncId))
	bridge.DebugInfo(debugInfo)
	serviceFuncDesc, empty := bridge.Op.ProtoMap.GetServiceFuncById(int(callMsg.SidFid))
	if empty {
		errInfo := "CallNativeService GetServiceFuncById empty , SidFid:" + strconv.Itoa(int(callMsg.SidFid))
		err = bridge.MakeError(errInfo)
		return resData, err
	}
	bridge.DebugInfo("CallNativeService msg.ServiceName: " + serviceFuncDesc.ServiceName + " FuncName: " + serviceFuncDesc.FuncName + " sourceUid:" + strconv.Itoa(int(callMsg.SourceUid)))

	switch serviceFuncDesc.ServiceName {
	case "GameMatch":
		bridge.NativeServiceList.GameMatch <- callMsg.Msg
	case "Gateway":
		bridge.NativeServiceList.Gateway <- callMsg.Msg
	case "FrameSync":
		bridge.NativeServiceList.FrameSync <- callMsg.Msg
	case "TwinAgora":
		bridge.NativeServiceList.TwinAgora <- callMsg.Msg
	default:
		errInfo := "CallNativeService switch ServiceName , name: " + serviceFuncDesc.ServiceName
		err = bridge.MakeError(errInfo)
		return resData, err
	}

	return "insert chan Queue......", nil
}

func (bridge *Bridge) MakeError(info string) error {
	info = bridge.Prefix + " err , " + info
	bridge.Op.Log.Error(info)
	err := errors.New(info)
	return err
}

func (bridge *Bridge) DebugInfo(info string) string {
	info = bridge.Prefix + " , " + info
	bridge.Op.Log.Debug(info)
	return info
}

// 远程调用 GRPC/HTTP
func (bridge *Bridge) CallRemoteService(callMsg CallMsg) (resData interface{}, err error) {
	serviceFuncDesc, empty := bridge.Op.ProtoMap.GetServiceFuncById(int(callMsg.SidFid))
	if empty {

	}

	switch serviceFuncDesc.ServiceName {
	//case "FrameSync":
	//	resData , err = serviceCallBridge.CallServiceFuncFrameSync(funcName,balanceFactor,requestData)
	case "GameMatch":
		resData, err = bridge.CallServiceFuncGameMatch(callMsg)
	//case "Gateway":
	//	resData , err = serviceCallBridge.CallServiceFuncGateway(funcName,balanceFactor,requestData)
	//case "TwinAgora":
	//	resData , err = serviceCallBridge.CallServiceFuncTwinAgora(funcName,balanceFactor,requestData)

	default:
		return resData, errors.New("service name router failed.")
	}
	return resData, err
}

// 动态调用服务的函数 : GameMatch
func (bridge *Bridge) CallServiceFuncGameMatch(callMsg CallMsg) (data interface{}, err error) {
	//获取GRPC一个连接
	//grpcClient, err := bridge.Op.GrpcManager.GetGameMatchClient("GameMatch", balanceFactor)
	//if err != nil {
	//	return data, err
	//}

	serviceFuncDesc, empty := bridge.Op.ProtoMap.GetServiceFuncById(int(callMsg.SidFid))
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
		switch callMsg.Flag {
		case 1:
		case 2:
			grpcClient, err := bridge.Op.GrpcManager.GetGameMatchClient("GameMatch", callMsg.BalanceFactor)
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
