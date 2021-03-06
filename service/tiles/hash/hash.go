package hash

import (
	"context"
	"fmt"
	"path/filepath"

	tiles "github.com/chronark/charon/service/tiles/proto/tiles"
	"github.com/opentracing/opentracing-go"
)

func Request(ctx context.Context, req *tiles.Request) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "hashRequest")
	defer span.Finish()

	concatenated := fmt.Sprintf("%d/%d/%d", req.GetZ(), req.GetX(), req.GetY())
	return filepath.Join("tiles", concatenated)
}
