package handler

import (
	"context"
	"github.com/chronark/charon/pkg/log"
	"github.com/chronark/charon/service/filecache/filecache"
	proto "github.com/chronark/charon/service/filecache/proto/filecache"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Handler struct {
	Cache  *filecache.FileCache
	Logger log.Factory
}

func (h *Handler) Get(ctx context.Context, req *proto.GetRequest, res *proto.GetResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get()")
	defer span.Finish()
	h.Logger.For(ctx).Info("Loading from cache",
		zap.String("hash", req.GetHashKey()),
	)
	value, hit, err := h.Cache.Get(ctx, req.GetHashKey())

	if err != nil {
		h.Logger.For(ctx).Error("Could not load from cache", zap.Error(err))
		return err
	}
	res.Hit = hit
	res.File = value
	return nil
}

func (h *Handler) Set(ctx context.Context, req *proto.SetRequest, res *proto.SetResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Set()")
	defer span.Finish()

	h.Logger.For(ctx).Info("Saving to cache",
		zap.String("hash", req.GetHashKey()),
	)
	err := h.Cache.Set(ctx, req.GetHashKey(), req.GetFile())

	if err != nil {
		h.Logger.For(ctx).Error("Could not store to cache", zap.Error(err))
		return err
	}
	res.Created = true
	h.Logger.For(ctx).Info("Saved to cache",
		zap.String("hash", req.GetHashKey()),
	)
	return nil
}

func (h *Handler) Delete(ctx context.Context, req *proto.DeleteRequest, res *proto.DeleteResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete()")
	defer span.Finish()

	h.Logger.For(ctx).Info("Deleting from cache",
		zap.String("hash", req.GetHashKey()),
	)
	err := h.Cache.Delete(ctx, req.GetHashKey())

	if err != nil {
		h.Logger.For(ctx).Error("Could not delete from cache", zap.Error(err))
		return err
	}
	res.Deleted = true
	h.Logger.For(ctx).Info("Deleted from cache",
		zap.String("hash", req.GetHashKey()),
	)
	return nil
}
