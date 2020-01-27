package main

import (
	"context"
	"fmt"
	"github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/client"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type osmHandler struct {
	Client client.Client
	Logger *logrus.Logger
}

func (h *osmHandler) Get(ctx context.Context, req *tiles.Request, res *tiles.Response) error {
	h.Logger.Infof("Requesting %+v", req)
	hashKey := hashRequest(req)

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(context.TODO(), &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		log.Errorf("Could not get file from filecache: %v\n", err)
	}
	var tile []byte
	if filecacheGetResponse.GetHit() {
		tile = filecacheGetResponse.File
	} else {
		tileURL := fmt.Sprintf("https://a.tile.openstreetmap.org/%d/%d/%d.png", req.GetZ(), req.GetX(), req.GetY())
		resp, err := http.Get(tileURL)
		if err != nil {
			log.Errorf("Could not load tile from osm: %v\n", err)
			return err
		}
		defer resp.Body.Close()
		tile, err = ioutil.ReadAll(resp.Body)
		_, err = fileCacheClient.Set(context.TODO(), &filecache.SetRequest{HashKey: hashKey, File: tile})
		if err != nil {
			log.Errorf("Could not set file to filecache: %v\n", err)
		}
	}
	res.File = tile
	return nil
}
func (h *osmHandler) Delete(ctx context.Context, req *tiles.Request, res *tiles.Response) error {
	return nil
}
