package filesystem

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type Local struct {
	basePath string
}

func NewLocal(u *url.URL) (FileSystem, error) {
	fs := &Local{
		basePath: filepath.Join(u.Host, u.Path),
	}
	if _, err := fs.ReadDir(""); err != nil {
		return nil, err
	}
	return fs, nil
}

func (f *Local) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(filepath.Join(f.basePath, dirname))
}

func (f *Local) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(filepath.Join(f.basePath, name))
}

func (f *Local) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *Local) Open(name string) (File, error) {
	return os.Open(filepath.Join(f.basePath, name))
}

func (f *Local) Close() error {
	return nil
}
