package nominatim

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/geocoding/proto/geocoding"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Handler struct {
	Logger log.Factory
	Client geocoding.GeocodingService
}

func (h *Handler) Forward(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "Forward()")
	defer span.Finish()

	h.Logger.For(ctx).Info(
		"Performing forward geocoding",
		zap.String("user", r.RemoteAddr),
		zap.String("request", r.URL.String()),
	)

	h.Logger.For(ctx).Info("request",
		zap.String("user", r.RemoteAddr),
		zap.String("url", r.URL.String()),
	)

	query := r.URL.Query().Get("query")
	if query == "" {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("No query found", zap.String("url", r.URL.String()))
		http.Error(w, "parameter 'query' was empty'", http.StatusBadRequest)
		return
	}
	rsp, err := h.Client.Forward(ctx, &geocoding.Search{Query: query})
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not get geocoding from service",
			zap.String("query", query),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(rsp.GetPayload())
	span.SetTag("http.status", 200)
	return

}

func (h *Handler) Reverse(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "Reverse()")
	defer span.Finish()
	h.Logger.For(ctx).Info(
		"Performing reverse geocoding",
		zap.String("user", r.RemoteAddr),
		zap.String("request", r.URL.String()),
	)

	latString := r.URL.Query().Get("lat")
	if latString == "" {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("lat was empty")

		http.Error(w, "parameter 'lat' was empty'", http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latString, 32)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not convert lat to float", zap.Error(err))
		return
	}
	lonString := r.URL.Query().Get("lon")
	if lonString == "" {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("lat was empty")
		http.Error(w, "parameter 'lon' was empty'", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonString, 32)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not convert lon to float", zap.Error(err))
		return
	}

	h.Logger.For(ctx).Info("lat and lon",
		zap.Float64("lat", lat),
		zap.Float64("lon", lon),
	)
	rsp, err := h.Client.Reverse(
		ctx,
		&geocoding.Coordinates{
			Lat: float32(lat),
			Lon: float32(lon),
		},
	)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Geocoding service returned error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(rsp.GetPayload())
	span.SetTag("http.status_code", 200)
	return

}
