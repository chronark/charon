package main

import (
	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/service/geocoding/handler"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"os"
	"time"
)

const (
	geocodingProviderNotFoundError = "The environmentvariable 'GEOCODING_PROVIDER' must be set with 'nominatim'"
)

var (
	serviceName       = "charon.srv.geocoding"
	geocodingProvider string
	logger            *logrus.Entry
)

func init() {
	geocodingProvider = os.Getenv("GEOCODING_PROVIDER")
	if geocodingProvider != "" {
		serviceName = serviceName + "." + geocodingProvider
	}
	logger = logging.New(serviceName)

}

type LogrusAdaper struct{}

func (l LogrusAdaper) Error(msg string) {
	logger.Errorf(msg)
}

func (l LogrusAdaper) Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func main() {
	factory := jaegerprom.New()
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})
	time.Sleep(5 * time.Second)
	transport, err := jaeger.NewUDPTransport(("jaeger:5775"), 0)
	if err != nil {
		logger.Error(err)
	}

	logAdapt := LogrusAdaper{}
	reporter := jaeger.NewCompositeReporter(
		jaeger.NewLoggingReporter(logAdapt),
		jaeger.NewRemoteReporter(transport,
			jaeger.ReporterOptions.Metrics(metrics),
			jaeger.ReporterOptions.Logger(logAdapt),
		),
	)
	defer reporter.Close()

	sampler := jaeger.NewConstSampler(true)

	tracer, closer := jaeger.NewTracer("geocoding",
		sampler, reporter, jaeger.TracerOptions.Metrics(metrics),
	)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
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
		logger.Error(geocodingProviderNotFoundError)
	}

	geocoding.RegisterGeocodingHandler(service.Server(), srvHandler)

	// Run service
	if err := service.Run(); err != nil {
		logger.Error(err.Error())
	}
}
