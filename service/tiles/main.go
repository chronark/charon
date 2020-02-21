package main

import (
	"context"
	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/tiles/handler/osm"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"os"
)

const providerError = "You must set the environment variable 'TILE_PROVIDER' with either 'osm' or 'mapbox'"

var serviceName = "charon.srv.tiles"
var tileProvider string

func init() {
	tileProvider = os.Getenv("TILE_PROVIDER")
	if tileProvider != "" {
		serviceName = serviceName + "." + tileProvider
	}
}

func main() {
	logger := log.NewDefaultLogger(serviceName)

	tracer, closer := tracing.NewTracer(serviceName, logger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "main()")
	defer span.Finish()

	// New Service
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
		micro.WrapHandler(opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer())),
	)
	// optionally setup command line usage
	service.Init()

	// Register Handlers
	tiles.RegisterTilesHandler(service.Server(), &osm.Handler{
		Logger: logger,
		Client: service.Client(),
	})

	// Run server
	if err := service.Run(); err != nil {
		logger.For(ctx).Fatal(err.Error())
	}
}
