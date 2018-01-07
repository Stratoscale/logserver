package filesystem

import (
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type LocalFS struct {
	basePath string
}

func NewLocalFS(u *url.URL) (*LocalFS, error) {
	fs := &LocalFS{
		basePath: filepath.Join(u.Host, u.Path),
	}
	if _, err := fs.ReadDir(""); err != nil {
		return nil, err
	}
	log.Printf("Opened local: %s", fs.basePath)
	return fs, nil
}

func (f *LocalFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(filepath.Join(f.basePath, dirname))
}

func (f *LocalFS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(filepath.Join(f.basePath, name))
}

func (f *LocalFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *LocalFS) Open(name string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(f.basePath, name))
}

func (f *LocalFS) Close() error {
	return nil
}
