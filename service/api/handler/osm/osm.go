package osm

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Handler struct {
	Logger log.Factory
	Client tiles.TilesService
}

func (h *Handler) parseCoordinates(ctx context.Context, r *http.Request) (*tiles.Request, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "parseCoordinates")
	defer span.Finish()

	x := r.URL.Query().Get("x")
	if x == "" {
		err := fmt.Errorf("Parameter x was empty")
		span.SetTag("error", true)
		h.Logger.For(ctx).Error(err.Error())
		return nil, err
	}
	xInt, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Error parsing x", zap.Error(err))
		return nil, err
	}

	y := r.URL.Query().Get("y")
	if y == "" {
		err := fmt.Errorf("Parameter y was empty")
		span.SetTag("error", true)
		h.Logger.For(ctx).Error(err.Error())
		return nil, err
	}
	yInt, err := strconv.ParseInt(y, 10, 32)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Error parsing y", zap.Error(err))
		return nil, err
	}

	z := r.URL.Query().Get("z")
	if z == "" {
		err := fmt.Errorf("Parameter z was empty")
		span.SetTag("error", true)
		h.Logger.For(ctx).Error(err.Error())
		return nil, err
	}
	zInt, err := strconv.ParseInt(z, 10, 32)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Error parsing z", zap.Error(err))
		return nil, err
	}

	coords := tiles.Request{
		X: int32(xInt),
		Y: int32(yInt),
		Z: int32(zInt),
	}
	h.Logger.For(ctx).Info(
		"Requested coordinates",
		zap.Int32("x", coords.X),
		zap.Int32("y", coords.Y),
		zap.Int32("z", coords.Z),
	)
	return &coords, nil

}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "Get")
	defer span.Finish()

	h.Logger.For(ctx).Info(
		"request",
		zap.String("user", r.RemoteAddr),
		zap.String("request", r.URL.String()),
	)

	req, err := h.parseCoordinates(ctx, r)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Error parsing coordinates", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rsp, err := h.Client.Get(ctx, req)
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Error fetching tile from tile service", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(rsp.GetFile())
	if err != nil {
		span.SetTag("error", true)
		h.Logger.For(ctx).Error("Could not write data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetTag("http.status", 200)
}
