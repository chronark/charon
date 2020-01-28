package main

import (
	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	micro "github.com/micro/go-micro"
	"github.com/sirupsen/logrus"
	"os"
)

const providerError = "You must set the environment variable 'TILE_PROVIDER' with either 'osm' or 'mapbox'"

var serviceName = "charon.srv.tiles"
var tileProvider string

var log *logrus.Logger

func init() {
	tileProvider = os.Getenv("TILE_PROVIDER")

	serviceName = serviceName + "." + tileProvider

	log = logging.New(serviceName)

}

func main() {
	log.Infof("Initializing %s", serviceName)

	service := micro.NewService(
		micro.Name(serviceName),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	tiles.RegisterTilesHandler(service.Server(), &osmHandler{
		Logger: log,
		Client: service.Client(),
	})

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
