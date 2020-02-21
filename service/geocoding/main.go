package main

import (
	"context"
	"os"
	"time"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/geocoding/handler"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

const (
	geocodingProviderNotFoundError = "The environmentvariable 'GEOCODING_PROVIDER' must be set with 'nominatim'"
)

var (
	serviceName       = "charon.srv.geocoding"
	geocodingProvider string
)

func init() {
	geocodingProvider = os.Getenv("GEOCODING_PROVIDER")
	if geocodingProvider != "" {
		serviceName = serviceName + "." + geocodingProvider
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

	// Initialise service
	service.Init()

	// Register Handler
	var srvHandler geocoding.GeocodingHandler

	switch geocodingProvider {
	case "nominatim":
		srvHandler = &handler.Nominatim{Logger: logger, Throttle: time.Tick(time.Second), Client: client.DefaultClient}
		break
	default:
		logger.For(ctx).Fatal(geocodingProviderNotFoundError)
	}

	geocoding.RegisterGeocodingHandler(service.Server(), srvHandler)

	// Run service
	if err := service.Run(); err != nil {
		logger.For(ctx).Fatal(err.Error())
	}
}
