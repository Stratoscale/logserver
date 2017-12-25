package filesystem

import (
	"io/ioutil"
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
