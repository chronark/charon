package osm

import (
	"context"
	"fmt"
	"github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func parseCoordinates(r *http.Request) (*tiles.Request, error) {
	x := r.URL.Query().Get("x")
	if x == "" {
		return nil, fmt.Errorf("Parameter x was empty")
	}
	xInt, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error parsing x: %w", err)
	}

	y := r.URL.Query().Get("y")
	if x == "" {
		return nil, fmt.Errorf("Parameter y was empty")
	}
	yInt, err := strconv.ParseInt(y, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error parsing y: %w", err)
	}

	z := r.URL.Query().Get("z")
	if z == "" {
		return nil, fmt.Errorf("Parameter z was empty")
	}
	zInt, err := strconv.ParseInt(z, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error parsing z: %w", err)
	}

	return &tiles.Request{
		X: int32(xInt),
		Y: int32(yInt),
		Z: int32(zInt),
	}, nil

}

type Handler struct {
	Logger *logrus.Entry
	Client tiles.TilesService
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("User %s has requested %s", r.RemoteAddr, r.URL)

	req, err := parseCoordinates(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rsp, err := h.Client.Get(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(rsp.GetFile())
	return

}
