package main

import (
	"context"
	"net/http"
	"os"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/pkg/tracing"
	"github.com/chronark/charon/service/api/handler/nominatim"
	"github.com/chronark/charon/service/api/handler/osm"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-micro/v2"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

const serviceName = "charon.api"

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
	api := web.NewService(
		web.Name(serviceName),
		web.Address(serviceAddress),
	)
	service := micro.NewService(
		micro.Name("charon.srv.api"),
		micro.Version("latest"),
		micro.WrapHandler(opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer())),
	)



	nominatimHandler := &nominatim.Handler{
		Logger: logger,
		Client: geocoding.NewGeocodingService("charon.srv.geocoding.nominatim", service.Client()),
	}

	osmHandler := &osm.Handler{
		Logger: logger,
		Client: tiles.NewTilesService("charon.srv.tiles.osm", service.Client()),
	}

	api.Handle("/geocoding/forward/", corsWrapper(http.HandlerFunc(nominatimHandler.Forward)))
	api.Handle("/geocoding/reverse/", corsWrapper(http.HandlerFunc(nominatimHandler.Reverse)))

	api.Handle("/tile/", corsWrapper(http.HandlerFunc(osmHandler.Get)))

	if err := api.Init(); err != nil {
		logger.For(ctx).Fatal(err.Error())
	}

	if err := api.Run(); err != nil {
		logger.For(ctx).Fatal(err.Error())
	}

}
