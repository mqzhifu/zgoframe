package test

import (
	"zgoframe/core/global"
	"zgoframe/util"
)

func Gateway(){
	gateway := util.NewGateway(global.V.Grpc,global.V.Zap)
	gateway.StartSocket()
}
