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

//获取一个服务的grpc client : TwinAgora
func (grpcManager *GrpcManager) GetTwinAgoraClient(name string, balanceFactor string) (pb.TwinAgoraClient, error) {
	client, err := grpcManager.GetClientByLoadBalance(name, balanceFactor)
	if err != nil {
		return nil, err
	}

	return client.(pb.TwinAgoraClient), nil
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
	case "TwinAgora":
		incClient = pb.NewTwinAgoraClient(myGrpcClient.ClientConn)

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
	case "CS_PlayerState":
		request := pb.PlayerBase{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerState(ctx, &request)
	case "CS_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Heartbeat(ctx, &request)
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
		request := pb.RoomHistorySets{}
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
	case "SC_PlayerState":
		request := pb.PlayerState{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_PlayerState(ctx, &request)
	case "SC_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Heartbeat(ctx, &request)
	case "FdClose":
		request := pb.FDCloseEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdClose(ctx, &request)
	case "FdCreate":
		request := pb.FDCreateEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdCreate(ctx, &request)

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
		request := pb.GameMatchSign{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerMatchSign(ctx, &request)
	case "CS_PlayerMatchSignCancel":
		request := pb.GameMatchPlayerCancel{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PlayerMatchSignCancel(ctx, &request)
	case "CS_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Heartbeat(ctx, &request)
	case "SC_GameMatchOptResult":
		request := pb.GameMatchOptResult{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_GameMatchOptResult(ctx, &request)
	case "SC_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Heartbeat(ctx, &request)
	case "FdClose":
		request := pb.FDCloseEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdClose(ctx, &request)
	case "FdCreate":
		request := pb.FDCreateEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdCreate(ctx, &request)

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
	case "FdClose":
		request := pb.FDCloseEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdClose(ctx, &request)
	case "FdCreate":
		request := pb.FDCreateEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdCreate(ctx, &request)
	case "CS_Login":
		request := pb.Login{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Login(ctx, &request)
	case "CS_Ping":
		request := pb.PingReq{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Ping(ctx, &request)
	case "CS_Pong":
		request := pb.PongRes{}
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
		request := pb.PingReq{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_Ping(ctx, &request)
	case "SC_Pong":
		request := pb.PongRes{}
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
	case "SC_ProjectPushMsg":
		request := pb.ProjectPushMsg{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_ProjectPushMsg(ctx, &request)
	case "SC_SendMsg":
		request := pb.Msg{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_SendMsg(ctx, &request)

	default:
		return data, errors.New("func name router failed.")
	}
	return data, err
}

//动态调用服务的函数 : TwinAgora
func (grpcManager *GrpcManager) CallServiceFuncTwinAgora(funcName string, balanceFactor string, postData []byte) (data interface{}, err error) {
	//获取GRPC一个连接
	grpcClient, err := grpcManager.GetTwinAgoraClient("TwinAgora", balanceFactor)
	if err != nil {
		return data, err
	}

	ctx := context.Background()
	switch funcName {
	case "CS_Heartbeat":
		request := pb.Heartbeat{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_Heartbeat(ctx, &request)
	case "CS_CallPeople":
		request := pb.CallPeopleReq{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_CallPeople(ctx, &request)
	case "CS_CancelCallPeople":
		request := pb.CancelCallPeopleReq{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_CancelCallPeople(ctx, &request)
	case "CS_PeopleEntry":
		request := pb.PeopleEntry{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PeopleEntry(ctx, &request)
	case "CS_PeopleLeave":
		request := pb.PeopleLeaveRes{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_PeopleLeave(ctx, &request)
	case "CS_CallPeopleAccept":
		request := pb.CallVote{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_CallPeopleAccept(ctx, &request)
	case "CS_CallPeopleDeny":
		request := pb.CallVote{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_CallPeopleDeny(ctx, &request)
	case "CS_RoomHeartbeat":
		request := pb.RoomHeartbeatReq{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.CS_RoomHeartbeat(ctx, &request)
	case "SC_CallPeople":
		request := pb.CallPeopleRes{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_CallPeople(ctx, &request)
	case "SC_CancelCallPeople":
		request := pb.CancelCallPeopleReq{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_CancelCallPeople(ctx, &request)
	case "SC_PeopleEntry":
		request := pb.PeopleEntry{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_PeopleEntry(ctx, &request)
	case "SC_PeopleLeave":
		request := pb.PeopleLeaveRes{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_PeopleLeave(ctx, &request)
	case "SC_CallReply":
		request := pb.CallReply{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_CallReply(ctx, &request)
	case "SC_CallPeopleAccept":
		request := pb.CallVote{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_CallPeopleAccept(ctx, &request)
	case "SC_CallPeopleDeny":
		request := pb.CallVote{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.SC_CallPeopleDeny(ctx, &request)
	case "FdClose":
		request := pb.FDCloseEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdClose(ctx, &request)
	case "FdCreate":
		request := pb.FDCreateEvent{}
		err := json.Unmarshal(postData, &request)
		if err != nil {
			return data, err
		}
		data, err = grpcClient.FdCreate(ctx, &request)

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
	case "TwinAgora":
		resData, err = grpcManager.CallServiceFuncTwinAgora(funcName, balanceFactor, requestData)

	default:
		return requestData, errors.New("service name router failed.")
	}

	return resData, err
}
