package main

import (
	"context"
	"github.com/chronark/charon/service/filecache/filecache"
	proto "github.com/chronark/charon/service/filecache/proto/filecache"
)

type handler struct {
	cache *filecache.FileCache
}

func (h *handler) Get(ctx context.Context, req *proto.GetRequest, res *proto.GetResponse) error {
	value, hit, err := h.cache.Get(req.GetHashKey())

	if err != nil {
		return err
	}
	res.Hit = hit
	res.File = value
	return nil
}

func (h *handler) Set(ctx context.Context, req *proto.SetRequest, res *proto.SetResponse) error {
	err := h.cache.Set(req.GetHashKey(), req.GetFile())

	if err != nil {
		return err
	}
	res.Created = true
	return nil
}

func (h *handler) Delete(ctx context.Context, req *proto.DeleteRequest, res *proto.DeleteResponse) error {
	err := h.cache.Delete(req.GetHashKey())

	if err != nil {
		return err
	}
	res.Deleted = true
	return nil
}
