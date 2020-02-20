package main

import (
	"github.com/chronark/charon/service/tiles/handler/osm"
	"github.com/micro/go-micro/v2"

	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"os"
)

const providerError = "You must set the environment variable 'TILE_PROVIDER' with either 'osm' or 'mapbox'"

var serviceName = "charon.srv.tiles"
var tileProvider string

var log *logrus.Entry

func init() {
	tileProvider = os.Getenv("TILE_PROVIDER")

	serviceName = serviceName + "." + tileProvider

	log = logging.New(serviceName)

}

func main() {
	log.Infof("Initializing %s", serviceName)

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
	// optionally setup command line usage
	service.Init()

	// Register Handlers
	tiles.RegisterTilesHandler(service.Server(), &osm.Handler{
		Logger: log,
		Client: service.Client(),
	})

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
