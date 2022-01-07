

    func (grpcManager *GrpcManager)GetFrameSyncClient(name string)(pb.FrameSyncClient,error){
        client, err := grpcManager.GetClientByLoadBalance(name,0)
        if err != nil{
            return nil,err
        }
    
        return client.(pb.FrameSyncClient),nil
    }
    func (grpcManager *GrpcManager)GetZgoframeClient(name string)(pb.ZgoframeClient,error){
        client, err := grpcManager.GetClientByLoadBalance(name,0)
        if err != nil{
            return nil,err
        }
    
        return client.(pb.ZgoframeClient),nil
    }
    func (grpcManager *GrpcManager)GetSyncClient(name string)(pb.SyncClient,error){
        client, err := grpcManager.GetClientByLoadBalance(name,0)
        if err != nil{
            return nil,err
        }
    
        return client.(pb.SyncClient),nil
    }


    func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string){
        var incClient interface{}
        switch serviceName {
                case "FrameSync":
        incClient = pb.NewFrameSyncClient(myGrpcClient.ClientConn)
    case "Zgoframe":
        incClient = pb.NewZgoframeClient(myGrpcClient.ClientConn)
    case "Sync":
        incClient = pb.NewSyncClient(myGrpcClient.ClientConn)

        }
    
        myGrpcClient.GrpcClientList[serviceName] = incClient
	}