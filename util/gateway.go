package util

import (
	"go.uber.org/zap"
)

type Gateway struct {
	GrpcManager  *GrpcManager
	Log          *zap.Logger
	NetWayOption NetWayOption
}
