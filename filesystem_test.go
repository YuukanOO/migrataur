package migrataur

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// mockFSAdapter represents a mock for the filesystem
var mockFSAdapter = &mockFileSystem{}

func TestMain(m *testing.M) {
	// Replace the fsAdapter by the mock and run the test suite
	oldFs := fsAdapter
	fsAdapter = mockFSAdapter

	defer func() {
		fsAdapter = oldFs
	}()

	os.Exit(m.Run())
}

type mockFileSystem struct {
	files []os.FileInfo // This where are stored expected files returns by ReadDir
}

func (fs *mockFileSystem) Create(path string) (file, error) {
	fs.files = append(fs.files, mockFileInfo{name: filepath.Base(path)})

	return mockFile{}, nil
}

func (fs *mockFileSystem) Remove(path string) error {
	name := filepath.Base(path)

	for i, f := range fs.files {
		if f.Name() == name {
			fs.files = append(fs.files[:i], fs.files[i+1:]...)
			break
		}
	}
	return nil
}

func (*mockFileSystem) MkdirAll(path string, mode os.FileMode) error     { return nil }
func (fs *mockFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) { return fs.files, nil }
func (fs *mockFileSystem) ReadFile(filename string) ([]byte, error)      { return []byte{}, nil }

// empty the filesystem adapter
func (fs *mockFileSystem) empty() {
	fs.files = []os.FileInfo{}
}

// hasFiles mocks the returns of ReadDir, it determines what the filesystem should have
func (fs *mockFileSystem) hasFiles(files ...mockFileInfo) {
	arr := make([]os.FileInfo, len(files))

	for i := range files {
		arr[i] = files[i]
	}

	fs.files = arr
}

// exists checks if a file with the given name exists in the mock fs adapter
func (fs *mockFileSystem) exists(name string) bool {
	for _, f := range fs.files {
		if f.Name() == name {
			return true
		}
	}

	return false
}

type mockFile struct{}

func (mockFile) Close() error                   { return nil }
func (mockFile) Write(data []byte) (int, error) { return len(data), nil }

type mockFileInfo struct {
	name string
	size int64
	dir  bool
}

func (m mockFileInfo) Name() string     { return m.name }
func (m mockFileInfo) Size() int64      { return m.size }
func (mockFileInfo) Mode() os.FileMode  { return os.ModeDevice }
func (mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool      { return m.dir }
func (mockFileInfo) Sys() interface{}   { return nil }
