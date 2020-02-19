package main

import (
	"context"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"log"
	"time"

	micro "github.com/micro/go-micro/v2"
)

func main() {
	time.Sleep(2 * time.Second)
	service := micro.NewService()
	service.Init()

	client := geocoding.NewGeocodingService("charon.srv.geocoding.nominatim", service.Client())
	for i := 0; i < 10; i++ {
		r, err := client.Forward(context.Background(), &geocoding.Search{Query: "Nürnberg"})
		if err != nil {
			log.Fatalf("Could not create geocoding request: %v", err)
		}
		log.Printf("Created: %+v", r)
	}

}
