package handler

import (
	"context"
	"fmt"
	"github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/micro/go-micro/v2/client"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Nominatim struct {
	Logger   *logrus.Entry
	Throttle <-chan time.Time
	Client   client.Client
}

func (h *Nominatim) request(url string) ([]byte, error) {

	// Obey the Ratelimit of 1 req / s
	<-h.Throttle

	// Call nominatim
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not make request: %w", err)
	}
	request.Header.Set("User-Agent", "www.hochschuljobboerse.de")

	h.Logger.Infof("Request %s", request.URL)
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Could not get nominatim response: %w", err)
	}
	defer response.Body.Close()

	// Return payload
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read body of response: %w", err)
	}
	return body, nil
}

// Forward makes a call to Nominatim for forward geocoding.
func (h *Nominatim) Forward(ctx context.Context, req *geocoding.Search, res *geocoding.ForwardResponse) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Forward()")
	defer span.Finish()

	span.LogFields(log.String("search", req.GetQuery()))
	h.Logger.Infof("Search: %s", req.Query)
	hashKey := filepath.Join("nominatim", "forward", req.Query+".json")

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(ctx, &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		span.LogFields(log.Error(err))
		h.Logger.Errorf("Could not get file from filecache: %v\n", err)
	}
	if filecacheGetResponse.GetHit() {
		span.LogFields(
			log.String("Cache hit", hashKey),
		)
		h.Logger.Debugf("Cache hit: %s", hashKey)
		res.Payload = filecacheGetResponse.File
	} else {
		span.LogFields(
			log.String("Cache miss", hashKey),
		)
		h.Logger.Debugf("Cache miss: %s", hashKey)
		parameters := []string{fmt.Sprintf("q=%s", req.Query), "format=geojson", "polygon_geojson=1", "limit=1"}
		url := "https://nominatim.openstreetmap.org/search?" + strings.Join(parameters, "&")
		geojson, err := h.request(url)
		if err != nil {
			span.LogFields(log.Error(err))

			return fmt.Errorf("Could not request response from nominatim: %w", err)
		}
		res.Payload = geojson
		span.LogFields(
			log.String("Writing to cache", hashKey),
		)
		h.Logger.Debugf("Writing %s to cache", hashKey)
		go fileCacheClient.Set(ctx, &filecache.SetRequest{HashKey: hashKey, File: geojson})
	}

	return nil
}

func (h *Nominatim) Reverse(ctx context.Context, req *geocoding.Coordinates, res *geocoding.ReverseResponse) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Forward()")
	defer span.Finish()

	span.LogFields(log.Float32("lat", req.Lat), log.Float32("lon", req.Lon))

	hashKey := filepath.Join("nominatim", "reverse", fmt.Sprintf("%f.%f", req.Lat, req.Lon)+".json")

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(ctx, &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		span.LogFields(log.Error(err))
		h.Logger.Errorf("Could not get file from filecache: %v\n", err)
	}
	if filecacheGetResponse.GetHit() {
		res.Payload = filecacheGetResponse.File
	} else {
		parameters := []string{"format=geojson", "polygon_geojson=1", "zoom=3", "limit=1", fmt.Sprintf("lon=%f", req.GetLon()), fmt.Sprintf("lat=%f", req.GetLat())}
		url := "https://nominatim.openstreetmap.org/reverse?" + strings.Join(parameters, "&")

		geojson, err := h.request(url)
		if err != nil {
			span.LogFields(log.Error(err))
			return fmt.Errorf("Could not request response from nominatim: %w", err)
		}

		res.Payload = geojson
		go fileCacheClient.Set(ctx, &filecache.SetRequest{HashKey: hashKey, File: geojson})
	}

	return nil

}
