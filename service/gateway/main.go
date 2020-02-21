package main

import (
	//"github.com/micro/go-micro/v2"

	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/gateway/handler/nominatim"
	"github.com/chronark/charon/service/gateway/handler/osm"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	//opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
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

func main() {
	tracer, closer, err := tracing.NewTracer(serviceName)
	if err != nil {
		log.Error("Could not connect to jaeger: " + err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// New Service
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

	service.Handle("/tile/", corsWrapper(http.HandlerFunc(osmHandler.Get)))

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
