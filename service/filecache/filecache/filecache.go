package filecache

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/chronark/charon/pkg/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Cache interface {
	Get(ctx context.Context, hashKey string) (val []byte, hit bool, err error)
	Set(ctx context.Context, hashKey string, val []byte) error
	Delete(ctx context.Context, hashKey string) error
	Destroy(ctx context.Context) error
}

type FileCache struct {
	Basepath string
	mutex    *sync.Mutex
	log      log.Factory
}

// New creates a new filecache instance. Items will be saved in the given path.
func New(basepath string, logger log.Factory) (*FileCache, error) {

	if basepath == "" {
		return nil, fmt.Errorf("basepath must not be empty")
	}
	fc := &FileCache{
		Basepath: basepath,
		mutex:    &sync.Mutex{},
		log:      logger,
	}

	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return nil, err
	}
	fmt.Printf("New Filecache initialized in %s\n", basepath)
	return fc, nil

}

func (fc *FileCache) Get(ctx context.Context, hashKey string) (value []byte, hit bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get()")
	defer span.Finish()
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	path := filepath.Join(fc.Basepath, hashKey)
	fc.log.For(ctx).Info("trying to get file",
		zap.String("path", path),
		zap.String("hashKey", hashKey),
	)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			fc.log.For(ctx).Info("Cache access",
				zap.Bool("hit", false),
				zap.String("path", path),
				zap.String("hashKey", hashKey),
			)
			return nil, false, nil
		}
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not open file", zap.Error(err))
		return nil, false, err
	}
	value, err = ioutil.ReadAll(file)
	if err != nil {
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not read file", zap.Error(err))

		return nil, false, err
	}
	fc.log.For(ctx).Info("Cache access",
		zap.Bool("hit", true),
		zap.String("path", path),
		zap.String("hashKey", hashKey),
	)
	return value, true, nil
}

func (fc *FileCache) Set(ctx context.Context, hashKey string, value []byte) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	span, ctx := opentracing.StartSpanFromContext(ctx, "Set()")
	defer span.Finish()

	destPath := filepath.Join(fc.Basepath, hashKey)
	fc.log.For(ctx).Info("trying to write to path",
		zap.String("path", destPath),
		zap.String("hashKey", hashKey),
	)
	err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
	if err != nil {
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not create directory", zap.Error(err))

		return err
	}
	file, err := os.Create(destPath)
	if err != nil {
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not create file", zap.Error(err))
		return err
	}
	_, err = file.Write(value)
	if err != nil {
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not write to file", zap.Error(err))
		file.Close()
		return err
	}
	err = file.Close()
	if err != nil {
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not close file", zap.Error(err))

		return err
	}
	fc.log.For(ctx).Info("Stored data in file",
		zap.String("hash", hashKey),
	)

	return nil

}
func (fc *FileCache) Delete(ctx context.Context, hashKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete()")
	defer span.Finish()
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.log.For(ctx).Info("Trying to delete file",
		zap.String("hash", hashKey),
	)

	path := filepath.Join(fc.Basepath, hashKey)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fc.log.For(ctx).Info("File does not exist",
				zap.String("path", path),
			)
			return nil
		}
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not get info about file", zap.Error(err))

		return err
	}
	err = os.Remove(path)
	if err != nil {
		span.SetTag("error", true)
		fc.log.For(ctx).Error("Could not remove file", zap.Error(err))
		return err
	}
	fc.log.For(ctx).Info("Deleted file",
		zap.String("path", path),
		zap.String("hash", hashKey),
	)

	return nil
}

func (fc *FileCache) Destroy(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete()")
	defer span.Finish()

	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.log.For(ctx).Info("Deleting all files")
	return os.RemoveAll(fc.Basepath)
}
