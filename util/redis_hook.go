package util

import (
	"github.com/go-redis/redis/v8"
	"context"
)

type TracingHook struct{}

//var _ redis.Hook = (*TracingHook)(nil)

func NewTracingHook() *TracingHook {
	return new(TracingHook)
}

func (TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	MyPrint("redis BeforeProcess:",cmd.String())
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

func (TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	MyPrint("redis AfterProcess:",cmd.String())
	//span := trace.SpanFromContext(ctx)
	//if err := cmd.Err(); err != nil {
	//	recordError(ctx, span, err)
	//}
	//span.End()
	return nil
}

func (TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	MyPrint("redis BeforeProcessPipeline:" )
	//if !trace.SpanFromContext(ctx).IsRecording() {
	//	return ctx, nil
	//}
	//
	//summary, cmdsString := rediscmd.CmdsString(cmds)
	//
	//ctx, span := tracer.Start(ctx, "pipeline "+summary)
	//span.SetAttributes(
	//	attribute.String("db.system", "redis"),
	//	attribute.Int("db.redis.num_cmd", len(cmds)),
	//	attribute.String("db.statement", cmdsString),
	//)

	return ctx, nil
}

func (TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	MyPrint("redis AfterProcessPipeline:" )
	//span := trace.SpanFromContext(ctx)
	//if err := cmds[0].Err(); err != nil {
	//	recordError(ctx, span, err)
	//}
	//span.End()
	return nil
}

//func recordError(ctx context.Context, span trace.Span, err error) {
//	if err != redis.Nil {
//		span.RecordError(err)
//		span.SetStatus(codes.Error, err.Error())
//	}
//}