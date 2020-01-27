package main

import (
	"github.com/chronark/charon/service/filecache/logging"
	"github.com/chronark/charon/service/filecache/filecache"
	proto "github.com/chronark/charon/service/filecache/proto/filecache"
	micro "github.com/micro/go-micro"
	"github.com/sirupsen/logrus"
)

var serviceName = "charon.srv.filecache"

var log *logrus.Logger

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

	service := micro.NewService(
		micro.Name(serviceName),
	)
	service.Init()

	err = proto.RegisterFilecacheServiceHandler(
		service.Server(),
		&handler{cache},
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
