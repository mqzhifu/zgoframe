package util

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"zgoframe/protobuf/pb"
)

//获取一个服务的grpc client : FrameSync
func (grpcManager *GrpcManager) GetFrameSyncClient(name string, balanceFactor string) (pb.FrameSyncClient, error) {
	client, err := grpcManager.GetClientByLoadBalance(name, balanceFactor)
	if err != nil {
		return nil, err
	}

	return client.(pb.FrameSyncClient), nil
}

//获取一个服务的grpc client : GameMatch
func (grpcManager *GrpcManager) GetGameMatchClient(name string, balanceFactor string) (pb.GameMatchClient, error) {
	client, err := grpcManager.GetClientByLoadBalance(name, balanceFactor)
	if err != nil {
		return nil, err
	}

	return client.(pb.GameMatchClient), nil
}

//获取一个服务的grpc client : Gateway
func (grpcManager *GrpcManager) GetGatewayClient(name string, balanceFactor string) (pb.GatewayClient, error) {
	client, err := grpcManager.GetClientByLoadBalance(name, balanceFactor)
	if err != nil {
		return nil, err
	}

	return client.(pb.GatewayClient), nil
}

//获取一个服务的grpc client : LogSlave
func (grpcManager *GrpcManager) GetLogSlaveClient(name string, balanceFactor string) (pb.LogSlaveClient, error) {
	client, err := grpcManager.GetClientByLoadBalance(name, balanceFactor)
	if err != nil {
		return nil, err
	}

	return client.(pb.LogSlaveClient), nil
}

//获取一个服务的grpc client : Zgoframe
func (grpcManager *GrpcManager) GetZgoframeClient(name string, balanceFactor string) (pb.ZgoframeClient, error) {
	client, err := grpcManager.GetClientByLoadBalance(name, balanceFactor)
	if err != nil {
		return nil, err
	}

	return client.(pb.ZgoframeClient), nil
}

//根据服务名获取一个GRPC-CLIENT 连接(c端使用)
func (myGrpcClient *MyGrpcClient) GetGrpcClientByServiceName(serviceName string, clientConn *grpc.ClientConn) (interface{}, error) {
	var incClient interface{}
	switch serviceName {
	case "FrameSync":
		incClient = pb.NewFrameSyncClient(myGrpcClient.ClientConn)
	case "GameMatch":
		incClient = pb.NewGameMatchClient(myGrpcClient.ClientConn)
	case "Gateway":
		incClient = pb.NewGatewayClient(myGrpcClient.ClientConn)
	case "LogSlave":
		incClient = pb.NewLogSlaveClient(myGrpcClient.ClientConn)
	case "Zgoframe":
		incClient = pb.NewZgoframeClient(myGrpcClient.ClientConn)

	default:
		return incClient, errors.New("service name router failed.")
	}
	return incClient, nil
}

//动态调用服务的函数 : FrameSync
func (grpcManager *GrpcManager) CallServiceFuncFrameSync(funcName string, balanceFactor string, postData []byte) (data interface{}, err error) {
	//获取GRPC一个连接
	grpcClient, err := grpcManager.GetFrameSyncClient("FrameSync", balanceFactor)
	if err != nil {
		return data, err
	}

	ctx := context.Background()
	switch funcName {
	case "CS_PlayerReady":
		request := pb.PlayerReady{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerReady(ctx, &request)
	case "CS_PlayerOperations":
		request := pb.LogicFrame{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerOperations(ctx, &request)
	case "CS_PlayerResumeGame":
		request := pb.PlayerResumeGame{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerResumeGame(ctx, &request)
	case "CS_PlayerOver":
		request := pb.PlayerOver{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerOver(ctx, &request)
	case "CS_RoomHistory":
		request := pb.ReqRoomHistory{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_RoomHistory(ctx, &request)
	case "CS_RoomBaseInfo":
		request := pb.RoomBaseInfo{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_RoomBaseInfo(ctx, &request)
	case "SC_ReadyTimeout":
		request := pb.ReadyTimeout{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_ReadyTimeout(ctx, &request)
	case "SC_EnterBattle":
		request := pb.EnterBattle{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_EnterBattle(ctx, &request)
	case "SC_LogicFrame":
		request := pb.LogicFrame{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_LogicFrame(ctx, &request)
	case "SC_RoomHistory":
		request := pb.RoomHistoryList{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_RoomHistory(ctx, &request)
	case "SC_RoomBaseInfo":
		request := pb.RoomBaseInfo{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_RoomBaseInfo(ctx, &request)
	case "SC_OtherPlayerOffline":
		request := pb.OtherPlayerOffline{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_OtherPlayerOffline(ctx, &request)
	case "SC_OtherPlayerOver":
		request := pb.PlayerOver{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_OtherPlayerOver(ctx, &request)
	case "SC_OtherPlayerResumeGame":
		request := pb.PlayerResumeGame{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_OtherPlayerResumeGame(ctx, &request)
	case "SC_StartBattle":
		request := pb.StartBattle{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_StartBattle(ctx, &request)
	case "SC_RestartGame":
		request := pb.RestartGame{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_RestartGame(ctx, &request)
	case "SC_GameOver":
		request := pb.GameOver{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_GameOver(ctx, &request)

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}

//动态调用服务的函数 : GameMatch
func (grpcManager *GrpcManager) CallServiceFuncGameMatch(funcName string, balanceFactor string, postData []byte) (data interface{}, err error) {
	//获取GRPC一个连接
	grpcClient, err := grpcManager.GetGameMatchClient("GameMatch", balanceFactor)
	if err != nil {
		return data, err
	}

	ctx := context.Background()
	switch funcName {
	case "CS_PlayerMatchSign":
		request := pb.PlayerMatchSign{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerMatchSign(ctx, &request)
	case "CS_PlayerMatchSignCancel":
		request := pb.PlayerMatchSignCancel{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerMatchSignCancel(ctx, &request)
	case "SC_PlayerMatchSignFailed":
		request := pb.PlayerMatchSignFailed{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_PlayerMatchSignFailed(ctx, &request)
	case "SC_PlayerMatchingFailed":
		request := pb.PlayerMatchingFailed{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_PlayerMatchingFailed(ctx, &request)

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}

//动态调用服务的函数 : Gateway
func (grpcManager *GrpcManager) CallServiceFuncGateway(funcName string, balanceFactor string, postData []byte) (data interface{}, err error) {
	//获取GRPC一个连接
	grpcClient, err := grpcManager.GetGatewayClient("Gateway", balanceFactor)
	if err != nil {
		return data, err
	}

	ctx := context.Background()
	switch funcName {
	case "CS_Login":
		request := pb.Login{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Login(ctx, &request)
	case "CS_Ping":
		request := pb.Ping{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Ping(ctx, &request)
	case "CS_Pong":
		request := pb.Pong{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Pong(ctx, &request)
	case "CS_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Heartbeat(ctx, &request)
	case "SC_Login":
		request := pb.LoginRes{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Login(ctx, &request)
	case "SC_Ping":
		request := pb.Ping{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Ping(ctx, &request)
	case "SC_Pong":
		request := pb.Pong{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Pong(ctx, &request)
	case "SC_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Heartbeat(ctx, &request)
	case "SC_KickOff":
		request := pb.KickOff{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_KickOff(ctx, &request)
	case "SC_ProjectPush":
		request := pb.ProjectPush{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_ProjectPush(ctx, &request)

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}

//动态调用服务的函数 : LogSlave
func (grpcManager *GrpcManager) CallServiceFuncLogSlave(funcName string, balanceFactor string, postData []byte) (data interface{}, err error) {
	//获取GRPC一个连接
	grpcClient, err := grpcManager.GetLogSlaveClient("LogSlave", balanceFactor)
	if err != nil {
		return data, err
	}

	ctx := context.Background()
	switch funcName {
	case "Push":
		request := pb.SlavePushMsg{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.Push(ctx, &request)

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}

//动态调用服务的函数 : Zgoframe
func (grpcManager *GrpcManager) CallServiceFuncZgoframe(funcName string, balanceFactor string, postData []byte) (data interface{}, err error) {
	//获取GRPC一个连接
	grpcClient, err := grpcManager.GetZgoframeClient("Zgoframe", balanceFactor)
	if err != nil {
		return data, err
	}

	ctx := context.Background()
	switch funcName {
	case "SayHello":
		request := pb.RequestUser{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SayHello(ctx, &request)
	case "Comm":
		request := pb.RequestUser{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.Comm(ctx, &request)

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}

//动态调用一个GRPC-SERVER 的一个方法(c端使用)
func (grpcManager *GrpcManager) CallGrpc(serviceName string, funcName string, balanceFactor string, requestData []byte) (resData interface{}, err error) {
	switch serviceName {
	case "FrameSync":
		resData, err = grpcManager.CallServiceFuncFrameSync(funcName, balanceFactor, requestData)
	case "GameMatch":
		resData, err = grpcManager.CallServiceFuncGameMatch(funcName, balanceFactor, requestData)
	case "Gateway":
		resData, err = grpcManager.CallServiceFuncGateway(funcName, balanceFactor, requestData)
	case "LogSlave":
		resData, err = grpcManager.CallServiceFuncLogSlave(funcName, balanceFactor, requestData)
	case "Zgoframe":
		resData, err = grpcManager.CallServiceFuncZgoframe(funcName, balanceFactor, requestData)

	default:
		return requestData, errors.New("service name router failed.")
	}

	return resData, err
}
