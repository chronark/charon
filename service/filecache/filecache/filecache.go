package filecache

import (
	"fmt"
	"github.com/chronark/charon/service/filecache/logging"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Cache interface {
	Get(hashKey string) (val []byte, hit bool, err error)
	Set(hashKey string, val []byte) error
	Delete(hashKey string) error
	Destroy() error
}

type FileCache struct {
	Basepath string
	mutex    *sync.Mutex
	logger   *logrus.Logger
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
	return fc, nil

}

func (fc *FileCache) Get(hashKey string) (value []byte, hit bool, err error) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	path := filepath.Join(fc.Basepath, hashKey)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			fc.logger.Infof("Miss %s", hashKey)
			return nil, false, nil
		}
		return nil, false, err
	}
	value, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, false, err
	}
	fc.logger.Infof("Hit %s", hashKey)
	return value, true, nil
}

func (fc *FileCache) Set(hashKey string, value []byte) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	destPath := filepath.Join(fc.Basepath, hashKey)
	err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	_, err = file.Write(value)
	if err != nil {
		file.Close()
		return err
	}
	err = file.Close()
	if err == nil {
		logrus.Infof("Stored %s", hashKey)
	}
	return err

}
func (fc *FileCache) Delete(hashKey string) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	path := filepath.Join(fc.Basepath, hashKey)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
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
