package main

import (
	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/filecache/filecache"
	"github.com/chronark/charon/service/filecache/handler"
	proto "github.com/chronark/charon/service/filecache/proto/filecache"
	micro "github.com/micro/go-micro/v2"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

var serviceName = "charon.srv.filecache"

var log *logrus.Entry

func init() {
	log = logging.New(serviceName)
}

func main() {
	log.Infof("Initializing %s\n", serviceName)

	cache, err := filecache.New("./cache")
	if err != nil {
		log.Fatalf("Error initializing cache: %w", err)
	}
	log.Info("filecache ready")

	tracer, closer, err := tracing.NewTracer(serviceName)
	if err != nil {
		log.Error("Could not connect to jaeger: " + err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

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
		&handler.Handler{Cache: cache},
	)
	if err != nil {
		log.Fatalf("Error registering filecache service handler: %w", err)
	}

	log.Infof("Service [%s] starting", serviceName)
	err = service.Run()
	if err != nil {
		log.Fatalf("Error running service: %w", err)
	}
}
