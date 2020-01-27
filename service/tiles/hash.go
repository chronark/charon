package main

import (
	"fmt"
	tiles "github.com/chronark/charon/service/tiles/proto/tiles"
	"path/filepath"
)

func hashRequest(req *tiles.Request) string {
	concatenated := fmt.Sprintf("%d/%d/%d", req.GetZ(), req.GetX(), req.GetY())
	return filepath.Join("tiles", concatenated)
}
