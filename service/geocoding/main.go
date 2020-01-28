package main

import (
	"github.com/chronark/charon/pkg/logging"
	"github.com/chronark/charon/service/geocoding/handler"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	geocodingProviderNotFoundError = "The environmentvariable 'GEOCODING_PROVIDER' must be set with 'nominatim'"
)

var (
	serviceName       = "charon.srv.geocoding"
	geocodingProvider string
	logger            *logrus.Logger
)

func init() {
	geocodingProvider = os.Getenv("GEOCODING_PROVIDER")
	if geocodingProvider != "" {
		serviceName = serviceName + "." + geocodingProvider
	}
	logger = logging.New(serviceName)

}
func main() {
	// New Service
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
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
		logger.Fatal(geocodingProviderNotFoundError)
	}

	geocoding.RegisterGeocodingHandler(service.Server(), srvHandler)

	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
