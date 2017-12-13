package migrataur

import (
	"io/ioutil"
	"os"
)

type fileSystem interface {
	MkdirAll(path string, mode os.FileMode) error
	Create(path string) (*os.File, error)
	Remove(path string) error
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

// osFileSystem implements the fileSystem interface using the os provider
type osFileSystem struct{}

func (osFileSystem) MkdirAll(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}

func (osFileSystem) Create(path string) (*os.File, error) {
	return os.Create(path)
}

func (osFileSystem) Remove(path string) error {
	return os.Remove(path)
}

func (osFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (osFileSystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
