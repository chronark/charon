package nominatim

import (
	"net/http"
	"strconv"

	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Logger *logrus.Entry
	Client geocoding.GeocodingService
}

func (h *Handler) Forward(w http.ResponseWriter, r *http.Request) {
	span := opentracing.GlobalTracer().StartSpan("Forward()")
	ctx := opentracing.ContextWithSpan(r.Context(), span)
	defer span.Finish()
	span.LogFields(
		log.String("user", r.RemoteAddr),
		log.String("request", r.URL.String()),
	)

	h.Logger.Infof("User %s has requested %s", r.RemoteAddr, r.URL)

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "parameter 'query' was empty'", http.StatusBadRequest)
		return
	}
	span.LogFields(
		log.String("query", query),
	)
	rsp, err := h.Client.Forward(ctx, &geocoding.Search{Query: query})
	if err != nil {
		span.LogFields(log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(rsp.GetPayload())
	span.SetTag("http.status", 200)
	return

}

func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Reverse()")
	defer span.Finish()
	span.LogFields(
		log.String("user", r.RemoteAddr),
		log.String("request", r.URL.String()),
	)
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

	span.LogFields(
		log.Float64("lat", lat),
		log.Float64("lon", lon),
	)
	rsp, err := h.Client.Reverse(
		ctx,
		&geocoding.Coordinates{
			Lat: float32(lat),
			Lon: float32(lon),
		},
	)
	if err != nil {
		span.LogFields(log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(rsp.GetPayload())
	span.SetTag("http.status_code", 200)
	return

}
