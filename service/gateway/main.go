package main

import (
	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/service/gateway/handler/nominatim"
	"github.com/chronark/charon/service/gateway/handler/osm"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"net/http"
	"os"
	"time"
)

const serviceName = "charon.service.gateway"

var serviceAddress string
var log *logrus.Entry

func init() {
	serviceAddress = os.Getenv("SERVICE_ADDRESS")
	log = logging.New(serviceName)
	if serviceAddress == "" {
		log.Error("You need to specify the environment variable 'SERVICE_ADDRESS' like 'host:post'")
	}
}

func corsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "GET")
		h.ServeHTTP(w, r)
	})
}

type LogrusAdaper struct{}

func (l LogrusAdaper) Error(msg string) {
	log.Errorf(msg)
}

func (l LogrusAdaper) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func OpenTracing(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wireCtx, _ := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))

		serverSpan := opentracing.StartSpan(r.URL.Path,
			ext.RPCServerOption(wireCtx))
		defer serverSpan.Finish()
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))

		h.ServeHTTP(w, r)
	})
}

func main() {
	factory := jaegerprom.New()
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})
	time.Sleep(5 * time.Second)
	transport, err := jaeger.NewUDPTransport(("jaeger:5775"), 0)
	if err != nil {
		log.Error(err)
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
	service := web.NewService(
		web.Name(serviceName),
		web.Address(serviceAddress),
	)

	nominatimHandler := &nominatim.Handler{
		Logger: log,
		Client: geocoding.NewGeocodingService("charon.srv.geocoding.nominatim", client.DefaultClient),
	}

	osmHandler := &osm.Handler{
		Logger: log,
		Client: tiles.NewTilesService("charon.srv.tiles.osm", client.DefaultClient),
	}

	service.Handle("/geocoding/forward/", corsWrapper(http.HandlerFunc(nominatimHandler.Forward)))
	service.Handle("/geocoding/reverse/", corsWrapper(http.HandlerFunc(nominatimHandler.Reverse)))

	service.Handle("/tile/", corsWrapper(OpenTracing(http.HandlerFunc(osmHandler.Get))))

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
