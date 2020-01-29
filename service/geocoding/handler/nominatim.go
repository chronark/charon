package handler

import (
	"context"
	"fmt"
	"github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/micro/go-micro/client"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Nominatim struct {
	Logger   *logrus.Logger
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
	h.Logger.Infof("Search: %s", req.Query)

	hashKey := filepath.Join("nominatim", "forward", req.Query + ".json")

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(context.TODO(), &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		h.Logger.Errorf("Could not get file from filecache: %v\n", err)
	}
	var geojson []byte
	if filecacheGetResponse.GetHit() {
		h.Logger.Debugf("Cache hit: %s", hashKey)
		geojson = filecacheGetResponse.File
	} else {
		h.Logger.Debugf("Cache miss: %s", hashKey)
		parameters := []string{fmt.Sprintf("q=%s", req.Query), "format=jsonv2", "polygon_geojson=1"}
		url := "https://nominatim.openstreetmap.org/search?" + strings.Join(parameters, "&")
		h.Logger.Warn(url)
		geojson, err := h.request(url)
		if err != nil {
			return fmt.Errorf("Could not request response from nominatim: %w", err)
		}
		h.Logger.Debugf("Writing %s to cache", hashKey)
		fileCacheClient.Set(context.TODO(), &filecache.SetRequest{HashKey: hashKey, File: geojson})
	}

	res.Payload = geojson
	return nil
}

func (h *Nominatim) Reverse(ctx context.Context, req *geocoding.Coordinates, res *geocoding.ReverseResponse) error {
	h.Logger.Infof("Search: lat %f, lon %f", req.Lat, req.Lon)

	hashKey := filepath.Join("nominatim", "reverse", fmt.Sprintf("%f.%f", req.Lat, req.Lon) + ".json")

	fileCacheClient := filecache.NewFilecacheService("charon.srv.filecache", h.Client)
	filecacheGetResponse, err := fileCacheClient.Get(context.TODO(), &filecache.GetRequest{HashKey: hashKey})
	if err != nil {
		h.Logger.Errorf("Could not get file from filecache: %v\n", err)
	}
	if filecacheGetResponse.GetHit() {
		res.Payload = filecacheGetResponse.File
	} else {
		parameters := []string{"format=geojson", "polygon_geojson=1","zoom=3", "limit=1", fmt.Sprintf("lon=%f", req.GetLon()), fmt.Sprintf("lat=%f", req.GetLat())}
		url := "https://nominatim.openstreetmap.org/reverse?" + strings.Join(parameters, "&")

		geojson, err := h.request(url)
		if err != nil {
			return fmt.Errorf("Could not request response from nominatim: %w", err)
		}
		
		fileCacheClient.Set(context.TODO(), &filecache.SetRequest{HashKey: hashKey, File: geojson})
		h.Logger.Infof("Payload: %v", geojson[:100])
		res.Payload = geojson
	}

	return nil

}