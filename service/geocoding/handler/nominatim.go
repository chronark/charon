package handler

import (
	"context"
	"fmt"
	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/micro/go-micro/v2/client"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Nominatim struct {
	Logger   log.Factory
	Throttle <-chan time.Time
	Client   client.Client
}

func (h *Nominatim) request(ctx context.Context, url string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "request")
	defer span.Finish()
	// Obey the Ratelimit of 1 req / s
	<-h.Throttle

	// Call nominatim
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not make request", zap.Error(err))
		return nil, err
	}
	request.Header.Set("User-Agent", "www.hochschuljobboerse.de")

	h.Logger.For(ctx).Info("Requesting", zap.String("url", request.URL.String()))
	response, err := client.Do(request)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not get nominatim response", zap.Error(err))
		return nil, err
	}
	defer response.Body.Close()

	// Return payload
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not read response body", zap.Error(err))
		return nil, err
	}
	return body, nil
}

// Forward makes a call to Nominatim for forward geocoding.
func (h *Nominatim) Forward(ctx context.Context, req *geocoding.Search, res *geocoding.ForwardResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Forward()")
	defer span.Finish()

	h.Logger.For(ctx).Info("searching", zap.String("query", req.GetQuery()))
	hashKey := filepath.Join("nominatim", "forward", req.Query+".json")

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(ctx, &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not get file from filecache", zap.Error(err))
	}

	if filecacheGetResponse.GetHit() {
		h.Logger.For(ctx).Info("Cache access",
			zap.Bool("hit", true),
			zap.String("hashKey", hashKey),
		)
		res.Payload = filecacheGetResponse.File
	} else {
		h.Logger.For(ctx).Info("Cache access",
			zap.Bool("hit", false),
			zap.String("hashKey", hashKey),
		)
		parameters := []string{fmt.Sprintf("q=%s", req.Query), "format=geojson", "polygon_geojson=1", "limit=1"}
		url := "https://nominatim.openstreetmap.org/search?" + strings.Join(parameters, "&")
		geojson, err := h.request(ctx, url)
		if err != nil {
			span.SetTag("error", true)
			h.Logger.For(ctx).Error("Could not get response from nominatim", zap.Error(err))
			return err
		}
		res.Payload = geojson
		_, err = fileCacheClient.Set(ctx, &filecache.SetRequest{HashKey: hashKey, File: geojson})
		if err != nil {
			h.Logger.For(ctx).Info(
				"Could not set file to filecache",
				zap.Error(err))
		}
	}

	return nil
}

func (h *Nominatim) Reverse(ctx context.Context, req *geocoding.Coordinates, res *geocoding.ReverseResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Forward()")
	defer span.Finish()

	h.Logger.For(ctx).Info("requesting", zap.Float32("lat", req.Lat), zap.Float32("lon", req.Lon))

	hashKey := filepath.Join("nominatim", "reverse", fmt.Sprintf("%f.%f", req.Lat, req.Lon)+".json")

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(ctx, &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not get response from filecache", zap.Error(err))
		return err
	}
	if filecacheGetResponse.GetHit() {
		res.Payload = filecacheGetResponse.File
	} else {
		parameters := []string{"format=geojson", "polygon_geojson=1", "zoom=3", "limit=1", fmt.Sprintf("lon=%f", req.GetLon()), fmt.Sprintf("lat=%f", req.GetLat())}
		url := "https://nominatim.openstreetmap.org/reverse?" + strings.Join(parameters, "&")

		geojson, err := h.request(ctx, url)
		if err != nil {
			span.SetTag("error", true)
			h.Logger.For(ctx).Error("Could not request response from nominatim", zap.Error(err))
			return err
		}

		res.Payload = geojson
		_, err = fileCacheClient.Set(ctx, &filecache.SetRequest{HashKey: hashKey, File: geojson})
		if err != nil {
			h.Logger.For(ctx).Info(
				"Could not set file to filecache",
				zap.Error(err))
		}
	}

	return nil

}
