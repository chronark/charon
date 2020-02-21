package filecache

import (
	"context"
	"fmt"
	"github.com/chronark/charon/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
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
	logger   *logrus.Entry
}

// New creates a new filecache instance. Items will be saved in the given path.
func New(basepath string) (*FileCache, error) {
	if basepath == "" {
		return nil, fmt.Errorf("Basepath must not be empty")
	}
	fc := &FileCache{
		Basepath: basepath,
		mutex:    &sync.Mutex{},
		logger:   logging.New("charon.srv.filecache"),
	}

	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return nil, err
	}
	fmt.Printf("New Filecache initialized in %s\n", basepath)
	return fc, nil

}

func (fc *FileCache) Get(ctx context.Context, hashKey string) (value []byte, hit bool, err error) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get()")
	defer span.Finish()
	span.LogFields(log.String("hash", hashKey))

	path := filepath.Join(fc.Basepath, hashKey)
	span.LogFields(log.String("path", path))

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			span.LogFields(log.Error(err), log.String("cache", "miss"))
			span.SetTag("error", true)
			fc.logger.Infof("Miss %s", hashKey)
			return nil, false, nil
		}
		return nil, false, err
	}
	span.LogFields(log.String("cache", "hit"))

	value, err = ioutil.ReadAll(file)
	if err != nil {
		span.SetTag("error", true)
		span.LogFields(log.Error(err))

		return nil, false, err
	}
	fc.logger.Infof("Hit %s", hashKey)
	return value, true, nil
}

func (fc *FileCache) Set(ctx context.Context, hashKey string, value []byte) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	span, ctx := opentracing.StartSpanFromContext(ctx, "Set()")
	defer span.Finish()
	span.LogFields(log.String("hash", hashKey))

	destPath := filepath.Join(fc.Basepath, hashKey)
	err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
	if err != nil {
		span.LogFields(log.Error(err))
		span.SetTag("error", true)

		return fmt.Errorf("Could not create directory for the destination: %w", err)
	}
	file, err := os.Create(destPath)
	if err != nil {
		span.SetTag("error", true)

		span.LogFields(
			log.String("message", "Could not create file"),
			log.Error(err),
		)
		return fmt.Errorf("Could not create file: %w", err)
	}
	_, err = file.Write(value)
	if err != nil {
		span.SetTag("error", err)
		file.Close()
		return fmt.Errorf("Could not write to file: %w", err)
	}
	err = file.Close()
	if err != nil {
		span.SetTag("error", err)
		return fmt.Errorf("Could not close file: %w", err)
	}
	fc.logger.Infof("Stored %s", hashKey)
	return nil

}
func (fc *FileCache) Delete(ctx context.Context, hashKey string) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete()")
	defer span.Finish()
	span.LogFields(log.String("hash", hashKey))

	path := filepath.Join(fc.Basepath, hashKey)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		span.SetTag("error", err)
		return err
	}
	err = os.Remove(path)
	if err == nil {
		fc.logger.Infof("Removed %s", hashKey)
	}
	return err
}

func (fc *FileCache) Destroy() error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	return os.RemoveAll(fc.Basepath)
}
