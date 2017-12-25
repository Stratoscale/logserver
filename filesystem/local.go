package filesystem

import (
	"io/ioutil"
	"io"
	"os"
	"path/filepath"
)

type LocalFS struct {
	BaseFS
}

func (f *LocalFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(filepath.Join(f.Url.Path, dirname))
}

func (f *LocalFS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(filepath.Join(f.Url.Path, name))
}

func (f *LocalFS) Join(elem ...string) string {
	return filepath.Join(f.Url.Path, filepath.Join(elem...))
}

func (f *LocalFS) Open(name string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(f.Url.Path, name))
}

