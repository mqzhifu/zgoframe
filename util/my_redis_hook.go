package util

import (
	"github.com/go-redis/redis/v8"
	"context"
	"go.uber.org/zap"
)

type TracingHook struct{
	Log *zap.Logger
}

//var _ redis.Hook = (*TracingHook)(nil)
const (
	MY_REDIS_HOOK_LOG_PREFIX = "my_redis_hook "
)

func NewTracingHook(log *zap.Logger) *TracingHook {
	tracingHook :=  new(TracingHook)
	tracingHook.Log = log
	return tracingHook
}

func (tracingHook TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	tracingHook.Log.Info (MY_REDIS_HOOK_LOG_PREFIX + "BeforeProcess:" + cmd.String())
	//if !trace.SpanFromContext(ctx).IsRecording() {
	//	return ctx, nil
	//}
	//
	//ctx, span := tracer.Start(ctx, cmd.FullName())
	//span.SetAttributes(
	//	attribute.String("db.system", "redis"),
	//	attribute.String("db.statement", rediscmd.CmdString(cmd)),
	//)

	return ctx, nil
}

func (tracingHook TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	//MyPrint(MY_REDIS_HOOK_LOG_PREFIX + "AfterProcess:")
	if cmd.Err() != nil {
		tracingHook.Log.Info(MY_REDIS_HOOK_LOG_PREFIX + "AfterProcess error:" + cmd.Err().Error())
	}
	return nil
}

func (tracingHook TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	tracingHook.Log.Info(MY_REDIS_HOOK_LOG_PREFIX + "BeforeProcessPipeline:" )
	return ctx, nil
}

func (tracingHook TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	tracingHook.Log.Info(MY_REDIS_HOOK_LOG_PREFIX + "AfterProcessPipeline:" )
	return nil
}

//func recordError(ctx context.Context, span trace.Span, err error) {
//	if err != redis.Nil {
//		span.RecordError(err)
//		span.SetStatus(codes.Error, err.Error())
//	}
//}