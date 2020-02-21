package osm

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/chronark/charon/service/tiles/hash"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/micro/go-micro/v2/client"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Handler struct {
	Client client.Client
	Logger log.Factory
}

func (h *Handler) Get(ctx context.Context, req *tiles.Request, res *tiles.Response) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	h.Logger.For(ctx).Info("Requesting",
		zap.Int32("x", req.GetX()),
		zap.Int32("y", req.GetY()),
		zap.Int32("z", req.GetZ()),
	)

	



	hashKey := hash.HashRequest(ctx, req)

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(ctx, &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Info(
			"Could not get file from filecache",
			zap.Error(err),
		)
	}
	var tile []byte
	if filecacheGetResponse.GetHit() {
		tile = filecacheGetResponse.File
	} else {
		tileURL := fmt.Sprintf("https://a.tile.openstreetmap.org/%d/%d/%d.png", req.GetZ(), req.GetX(), req.GetY())
		resp, err := http.Get(tileURL)
		if err != nil {
			span.SetTag("error", true)
			h.Logger.For(ctx).Info(
				"Could not load tile from osm",
				zap.Error(err),
			)
			return err
		}
		defer resp.Body.Close()
		tile, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			span.SetTag("error", true)
			h.Logger.For(ctx).Info(
				"Could not read response body",
				zap.Error(err))
			return err
		}
		_, err = fileCacheClient.Set(ctx, &filecache.SetRequest{HashKey: hashKey, File: tile})
		if err != nil {
			h.Logger.For(ctx).Info(
				"Could not set file to filecache",
				zap.Error(err))
		}
	}
	res.File = tile
	return nil
}
