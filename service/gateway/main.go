package main

import (
	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/service/gateway/handler/nominatim"
	"github.com/chronark/charon/service/gateway/handler/osm"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/chronark/charon/service/tiles/proto/tiles"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/web"
	"github.com/sirupsen/logrus"
	"os"
)

const serviceName = "charon.service.gateway"

var serviceAddress string
var log *logrus.Logger

func init() {
	serviceAddress = os.Getenv("SERVICE_ADDRESS")
	log = logging.New(serviceName)
	if serviceAddress == "" {
		log.Error("You need to specify the environment variable 'SERVICE_ADDRESS' like 'host:post'")
	}
}

func main() {
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

	service.HandleFunc("/geocoding/forward/", nominatimHandler.Forward)
	service.HandleFunc("/geocoding/reverse/", nominatimHandler.Reverse)

	service.HandleFunc("/tiles/", osmHandler.Get)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
