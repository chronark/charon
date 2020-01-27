package nominatim

import (
	"context"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	Logger *logrus.Logger
	Client geocoding.GeocodingService
}

func (h *Handler) Forward(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("User %s has requested %s", r.RemoteAddr, r.URL)

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "parameter 'query' was empty'", http.StatusBadRequest)
		return
	}

	rsp, err := h.Client.Forward(context.Background(), &geocoding.Search{Query: query})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp.GetPayload())
	return

}

func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("User %s has requested %s", r.RemoteAddr, r.URL)

	latString := r.URL.Query().Get("lat")
	if latString == "" {
		http.Error(w, "parameter 'lat' was empty'", http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latString, 32)
	if err != nil {
		h.Logger.Errorf("Could not convert lat to float: %w", err)
		return
	}
	lonString := r.URL.Query().Get("lon")
	if lonString == "" {
		http.Error(w, "parameter 'lon' was empty'", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonString, 32)
	if err != nil {
		h.Logger.Errorf("Could not convert lat to float: %w", err)
		return
	}

	rsp, err := h.Client.Reverse(
		context.Background(),
		&geocoding.Coordinates{
			Lat: float32(lat),
			Lon: float32(lon),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp.GetPayload())
	return

}
