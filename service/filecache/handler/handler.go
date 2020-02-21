package handler

import (
	"context"

	"github.com/chronark/charon/service/filecache/filecache"
	proto "github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Handler struct {
	Cache *filecache.FileCache
}

func (h *Handler) Get(ctx context.Context, req *proto.GetRequest, res *proto.GetResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get()")
	defer span.Finish()
	span.LogFields(log.String("hash", req.GetHashKey()))
	value, hit, err := h.Cache.Get(ctx, req.GetHashKey())
	if hit {
		span.LogEvent("cache-hit")
	} else {
		span.LogEvent("cache-miss")
	}

	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}
	res.Hit = hit
	res.File = value
	return nil
}

func (h *Handler) Set(ctx context.Context, req *proto.SetRequest, res *proto.SetResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Set()")
	defer span.Finish()

	span.LogFields(
		log.String("hash", req.GetHashKey()))
	err := h.Cache.Set(ctx, req.GetHashKey(), req.GetFile())

	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}
	res.Created = true
	return nil
}

func (h *Handler) Delete(ctx context.Context, req *proto.DeleteRequest, res *proto.DeleteResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete()")
	defer span.Finish()

	span.LogFields(
		log.String("hash", req.GetHashKey()))
	err := h.Cache.Delete(ctx, req.GetHashKey())

	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}
	res.Deleted = true
	return nil
}
