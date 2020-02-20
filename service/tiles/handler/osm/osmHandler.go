package osm

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/chronark/charon/service/tiles/hash"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Client client.Client
	Logger *logrus.Entry
}

func (h *Handler) Get(ctx context.Context, req *tiles.Request, res *tiles.Response) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	h.Logger.Infof("Requesting %+v", req)
	hashKey := hash.HashRequest(ctx, req)

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(ctx, &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(
			log.String("message", "Could not get file from filecache"),
			log.Error(err),
		)
	}
	var tile []byte
	if filecacheGetResponse.GetHit() {
		tile = filecacheGetResponse.File
	} else {
		tileURL := fmt.Sprintf("https://a.tile.openstreetmap.org/%d/%d/%d.png", req.GetZ(), req.GetX(), req.GetY())
		resp, err := http.Get(tileURL)
		if err != nil {
			span.LogFields(
				log.String("message", "Could not load tile from osm"),
				log.Error(err),
			)
			return err
		}
		defer resp.Body.Close()
		tile, err = ioutil.ReadAll(resp.Body)
		_, err = fileCacheClient.Set(ctx, &filecache.SetRequest{HashKey: hashKey, File: tile})
		if err != nil {
			span.LogFields(
				log.String("message", "Could not set file to filecache"),
				log.Error(err))
		}
	}
	res.File = tile
	return nil
}
func (h *Handler) Delete(ctx context.Context, req *tiles.Request, res *tiles.Response) error {
	return nil
}
