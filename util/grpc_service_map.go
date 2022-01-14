package util

import (
    "encoding/json"
    "errors"
    "google.golang.org/grpc"
    "zgoframe/protobuf/pb"
    "context"
)


//获取一个服务的grpc client : FrameSync
func (grpcManager *GrpcManager)GetFrameSyncClient(name string,balanceFactor string)(pb.FrameSyncClient,error){
    client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
    if err != nil{
        return nil,err
    }

    return client.(pb.FrameSyncClient),nil
}
//获取一个服务的grpc client : Gateway
func (grpcManager *GrpcManager)GetGatewayClient(name string,balanceFactor string)(pb.GatewayClient,error){
    client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
    if err != nil{
        return nil,err
    }

    return client.(pb.GatewayClient),nil
}
//获取一个服务的grpc client : LogSlave
func (grpcManager *GrpcManager)GetLogSlaveClient(name string,balanceFactor string)(pb.LogSlaveClient,error){
    client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
    if err != nil{
        return nil,err
    }

    return client.(pb.LogSlaveClient),nil
}
//获取一个服务的grpc client : Zgoframe
func (grpcManager *GrpcManager)GetZgoframeClient(name string,balanceFactor string)(pb.ZgoframeClient,error){
    client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
    if err != nil{
        return nil,err
    }

    return client.(pb.ZgoframeClient),nil
}


//根据服务名获取一个GRPC-CLIENT 连接(c端使用)
func  (myGrpcClient *MyGrpcClient)GetGrpcClientByServiceName(serviceName string,clientConn *grpc.ClientConn)(interface{},error){
    var incClient interface{}
    switch serviceName {
    case "FrameSync":
        incClient = pb.NewFrameSyncClient(myGrpcClient.ClientConn)
    case "Gateway":
        incClient = pb.NewGatewayClient(myGrpcClient.ClientConn)
    case "LogSlave":
        incClient = pb.NewLogSlaveClient(myGrpcClient.ClientConn)
    case "Zgoframe":
        incClient = pb.NewZgoframeClient(myGrpcClient.ClientConn)

    default:
        return incClient,errors.New("service name router failed.")
    }
    return incClient,nil
}
//动态调用服务的函数 : FrameSync
func (grpcManager *GrpcManager) CallServiceFuncFrameSync(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
    //获取GRPC一个连接
    grpcClient,err := grpcManager.GetFrameSyncClient("FrameSync",balanceFactor)
    if err != nil{
        return data,err
    }

    ctx := context.Background()
    switch funcName {
    case "PlayerOperations":
        request := pb.RequestPlayerOperations{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerOperations(ctx,&request)
    case "PlayerResumeGame":
        request := pb.RequestPlayerResumeGame{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerResumeGame(ctx,&request)
    case "PlayerReady":
        request := pb.RequestPlayerReady{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerReady(ctx,&request)
    case "PlayerOver":
        request := pb.RequestPlayerOver{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerOver(ctx,&request)
    case "RoomHistory":
        request := pb.RequestRoomHistory{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.RoomHistory(ctx,&request)
    case "GetRoom":
        request := pb.RequestGetRoom{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.GetRoom(ctx,&request)
    case "PlayerMatchSign":
        request := pb.RequestPlayerMatchSign{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerMatchSign(ctx,&request)
    case "PlayerMatchSignCancel":
        request := pb.RequestPlayerMatchSignCancel{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerMatchSignCancel(ctx,&request)
    case "EnterBattle":
        request := pb.ResponseEnterBattle{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.EnterBattle(ctx,&request)
    case "PushLogicFrame":
        request := pb.ResponsePushLogicFrame{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PushLogicFrame(ctx,&request)
    case "OtherPlayerOffline":
        request := pb.ResponseOtherPlayerOffline{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.OtherPlayerOffline(ctx,&request)
    case "ReadyTimeout":
        request := pb.ResponseReadyTimeout{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ReadyTimeout(ctx,&request)
    case "PushRoomHistory":
        request := pb.ResponsePushRoomHistory{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PushRoomHistory(ctx,&request)
    case "GameOver":
        request := pb.ResponseGameOver{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.GameOver(ctx,&request)
    case "PushRoomInfo":
        request := pb.ResponsePushRoomInfo{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PushRoomInfo(ctx,&request)
    case "StartBattle":
        request := pb.ResponseStartBattle{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.StartBattle(ctx,&request)
    case "OtherPlayerOver":
        request := pb.ResponseOtherPlayerOver{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.OtherPlayerOver(ctx,&request)
    case "RestartGame":
        request := pb.ResponseRestartGame{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.RestartGame(ctx,&request)
    case "PlayerMatchSignFailed":
        request := pb.ResponsePlayerMatchSignFailed{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerMatchSignFailed(ctx,&request)
    case "PlayerMatchingFailed":
        request := pb.ResponsePlayerMatchingFailed{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.PlayerMatchingFailed(ctx,&request)
    case "OtherPlayerResumeGame":
        request := pb.ResponseOtherPlayerResumeGame{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.OtherPlayerResumeGame(ctx,&request)

    default:
        return data,errors.New("func name router failed.")
    }
    return data,err
}
//动态调用服务的函数 : Gateway
func (grpcManager *GrpcManager) CallServiceFuncGateway(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
    //获取GRPC一个连接
    grpcClient,err := grpcManager.GetGatewayClient("Gateway",balanceFactor)
    if err != nil{
        return data,err
    }

    ctx := context.Background()
    switch funcName {
    case "ClientLogin":
        request := pb.RequestLogin{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ClientLogin(ctx,&request)
    case "ClientPing":
        request := pb.RequestClientPing{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ClientPing(ctx,&request)
    case "ClientPong":
        request := pb.RequestClientPong{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ClientPong(ctx,&request)
    case "ClientHeartbeat":
        request := pb.RequestClientHeartbeat{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ClientHeartbeat(ctx,&request)
    case "ServerPing":
        request := pb.ResponseServerPing{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ServerPing(ctx,&request)
    case "ServerPong":
        request := pb.ResponseServerPong{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ServerPong(ctx,&request)
    case "ServerHeartbeat":
        request := pb.RequestClientHeartbeat{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ServerHeartbeat(ctx,&request)
    case "ServerLogin":
        request := pb.ResponseLoginRes{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ServerLogin(ctx,&request)
    case "KickOff":
        request := pb.ResponseKickOff{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.KickOff(ctx,&request)
    case "ProjectPush":
        request := pb.RequestProjectPush{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.ProjectPush(ctx,&request)

    default:
        return data,errors.New("func name router failed.")
    }
    return data,err
}
//动态调用服务的函数 : LogSlave
func (grpcManager *GrpcManager) CallServiceFuncLogSlave(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
    //获取GRPC一个连接
    grpcClient,err := grpcManager.GetLogSlaveClient("LogSlave",balanceFactor)
    if err != nil{
        return data,err
    }

    ctx := context.Background()
    switch funcName {
    case "Push":
        request := pb.SlavePushMsg{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.Push(ctx,&request)

    default:
        return data,errors.New("func name router failed.")
    }
    return data,err
}
//动态调用服务的函数 : Zgoframe
func (grpcManager *GrpcManager) CallServiceFuncZgoframe(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
    //获取GRPC一个连接
    grpcClient,err := grpcManager.GetZgoframeClient("Zgoframe",balanceFactor)
    if err != nil{
        return data,err
    }

    ctx := context.Background()
    switch funcName {
    case "SayHello":
        request := pb.RequestUser{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.SayHello(ctx,&request)
    case "Comm":
        request := pb.RequestUser{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.Comm(ctx,&request)

    default:
        return data,errors.New("func name router failed.")
    }
    return data,err
}
//动态调用一个GRPC-SERVER 的一个方法(c端使用)
func (grpcManager *GrpcManager) CallGrpc(serviceName string,funcName string,balanceFactor string,requestData []byte)( resData interface{},err error){
    switch serviceName {
    case "FrameSync":
        resData , err = grpcManager.CallServiceFuncFrameSync(funcName,balanceFactor,requestData)
    case "Gateway":
        resData , err = grpcManager.CallServiceFuncGateway(funcName,balanceFactor,requestData)
    case "LogSlave":
        resData , err = grpcManager.CallServiceFuncLogSlave(funcName,balanceFactor,requestData)
    case "Zgoframe":
        resData , err = grpcManager.CallServiceFuncZgoframe(funcName,balanceFactor,requestData)

    default:
        return requestData,errors.New("service name router failed.")
    }

    return resData,err
}