package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type HttpFS struct {
	BaseFS
}

func (f *HttpFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(filepath.Join(f.Url.Path, dirname))
}

func (f *HttpFS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(filepath.Join(f.Url.Path, name))
}

func (f *HttpFS) Join(elem ...string) string {
	return filepath.Join(f.Url.Path, filepath.Join(elem...))
}
