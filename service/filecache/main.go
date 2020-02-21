package main

import (
	"context"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/filecache/filecache"
	"github.com/chronark/charon/service/filecache/handler"
	proto "github.com/chronark/charon/service/filecache/proto/filecache"
	micro "github.com/micro/go-micro/v2"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

var serviceName = "charon.srv.filecache"

func main() {
	logger := log.NewDefaultLogger(serviceName)

	tracer, closer := tracing.NewTracer(serviceName, logger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "main()")
	defer span.Finish()
	cache, err := filecache.New("./cache", logger)
	if err != nil {
		logger.For(ctx).Fatal("Error initializing cache", zap.Error(err))
	}
	logger.For(ctx).Info("filecache ready")

	// New Service
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
		micro.WrapHandler(opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer())),
	)

	service.Init()

	err = proto.RegisterFilecacheServiceHandler(
		service.Server(),
		&handler.Handler{Logger: logger, Cache: cache},
	)
	if err != nil {
		logger.For(ctx).Error("Error registering filecache service handler", zap.Error(err))
	}

	logger.For(ctx).Info("Service starting", zap.String("service name", serviceName))
	err = service.Run()
	if err != nil {
		logger.For(ctx).Error("Error running service", zap.Error(err))
	}
}
