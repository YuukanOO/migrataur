package migrataur

import (
	"io/ioutil"
	"os"
)

// Sooooo, this stuff is built to make testing more reliable by providing a mock for the
// filesystem without actually writing to the os. See http://nf.wh3rd.net/10things/#8

var fsAdapter fileSystem = osFileSystem{}

type fileSystem interface {
	MkdirAll(path string, mode os.FileMode) error
	Create(path string) (file, error)
	Remove(path string) error
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

type file interface {
	Write(data []byte) (int, error)
	Close() error
}

// osFileSystem implements the fileSystem interface using the os provider
type osFileSystem struct{}

func (osFileSystem) MkdirAll(path string, mode os.FileMode) error  { return os.MkdirAll(path, mode) }
func (osFileSystem) Create(path string) (file, error)              { return os.Create(path) }
func (osFileSystem) Remove(path string) error                      { return os.Remove(path) }
func (osFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) { return ioutil.ReadDir(dirname) }
func (osFileSystem) ReadFile(filename string) ([]byte, error)      { return ioutil.ReadFile(filename) }
