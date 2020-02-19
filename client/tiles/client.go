package main

import (
	"context"
	proto "github.com/chronark/charon/service/tiles/proto/tiles"
	micro "github.com/micro/go-micro/v2"
	"log"
	"math"
	"math/rand"
	"time"
)

func main() {
	time.Sleep(2 * time.Second)
	service := micro.NewService()
	service.Init()

	client := proto.NewTilesService("charon.srv.tiles.osm", service.Client())

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)

		z := rand.Intn(19)
		maxXY := int(math.Max(1, math.Pow(2, float64(z)-1)))
		x := rand.Intn(maxXY)
		y := rand.Intn(maxXY)
		log.Printf("Requesting X: %d, Y: %d, Z: %d.\n", x, y, z)
		_, err := client.Get(context.Background(), &proto.Request{
			X: int32(x),
			Y: int32(y),
			Z: int32(z),
		})
		if err != nil {
			log.Fatalf("Could not create tile request: %v", err)
		}
	}

}
