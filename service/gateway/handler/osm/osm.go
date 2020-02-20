package osm

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/sirupsen/logrus"
)

func parseCoordinates(ctx context.Context, r *http.Request) (*tiles.Request, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "parseCoordinates")
	defer span.Finish()

	x := r.URL.Query().Get("x")
	if x == "" {
		err := fmt.Errorf("Parameter x was empty")
		span.LogFields(log.Error(err))
		return nil, err
	}
	xInt, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, fmt.Errorf("Error parsing x: %w", err)
	}

	y := r.URL.Query().Get("y")
	if y == "" {
		err := fmt.Errorf("Parameter y was empty")
		span.LogFields(log.Error(err))
		return nil, err
	}
	yInt, err := strconv.ParseInt(y, 10, 32)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, fmt.Errorf("Error parsing y: %w", err)
	}

	z := r.URL.Query().Get("z")
	if z == "" {
		err := fmt.Errorf("Parameter z was empty")
		span.LogFields(log.Error(err))
		return nil, err
	}
	zInt, err := strconv.ParseInt(z, 10, 32)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, fmt.Errorf("Error parsing z: %w", err)
	}

	coords := tiles.Request{
		X: int32(xInt),
		Y: int32(yInt),
		Z: int32(zInt),
	}
	span.LogFields(
		log.Int32("x", coords.X),
		log.Int32("y", coords.Y),
		log.Int32("z", coords.Z),
	)
	return &coords, nil

}

type Handler struct {
	Logger *logrus.Entry
	Client tiles.TilesService
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Get")
	defer span.Finish()

	span.LogFields(
		log.String("user", r.RemoteAddr),
		log.String("request", r.URL.String()),
	)

	h.Logger.Infof("User %s has requested %s", r.RemoteAddr, r.URL)

	req, err := parseCoordinates(ctx, r)
	if err != nil {
		span.LogFields(log.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rsp, err := h.Client.Get(ctx, req)
	if err != nil {
		span.LogFields(log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	span.SetTag("http.status", 200)
	w.Header().Set("Content-Type", "image/png")
	w.Write(rsp.GetFile())
	return

}
