package main

import (
	"context"
	//"github.com/micro/go-micro/v2"
	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/gateway/handler/nominatim"
	"github.com/chronark/charon/service/gateway/handler/osm"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	//opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"os"
)

const serviceName = "charon.service.gateway"

func corsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "GET")
		h.ServeHTTP(w, r)
	})
}

func main() {
	logger := log.NewDefaultLogger(serviceName)

	tracer, closer := tracing.NewTracer(serviceName, logger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "main()")
	defer span.Finish()

	// New Service
	serviceAddress := os.Getenv("SERVICE_ADDRESS")
	if serviceAddress == "" {
		logger.For(ctx).Error("You need to specify the environment variable 'SERVICE_ADDRESS' like 'host:post'")
	}
	service := web.NewService(
		web.Name(serviceName),
		web.Address(serviceAddress),
	)

	nominatimHandler := &nominatim.Handler{
		Logger: logger,
		Client: geocoding.NewGeocodingService("charon.srv.geocoding.nominatim", client.DefaultClient),
	}

	osmHandler := &osm.Handler{
		Logger: logger,
		Client: tiles.NewTilesService("charon.srv.tiles.osm", client.DefaultClient),
	}

	service.Handle("/geocoding/forward/", corsWrapper(http.HandlerFunc(nominatimHandler.Forward)))
	service.Handle("/geocoding/reverse/", corsWrapper(http.HandlerFunc(nominatimHandler.Reverse)))

	service.Handle("/tile/", corsWrapper(http.HandlerFunc(osmHandler.Get)))

	if err := service.Init(); err != nil {
		logger.For(ctx).Fatal(err.Error())
	}

	if err := service.Run(); err != nil {
		logger.For(ctx).Fatal(err.Error())
	}

}
