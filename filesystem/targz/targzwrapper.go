package targz

import (
	"io"
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"github.com/Stratoscale/logserver/filesystem"
)

const suffix = ".tar.gz"

func New(inner filesystem.FileSystem) filesystem.FileSystem {
	return &wrap{inner: inner}
}

type wrap struct {
	inner filesystem.FileSystem
}

func split(dirname string) (string, string) {
	paths := strings.Split(dirname, suffix)
	tarName := filepath.Dir(paths[0])
	if len(paths) == 1 {
		return tarName, ""
	}
	innerPath := strings.Trim(paths[1], "/")
	return tarName, innerPath


}

func (w *wrap) getTfs(dirname string) (filesystem.FileSystem, string, error) {
	tarName, innerPath := split(dirname)
	f, err := w.inner.Open(tarName)
	defer f.Close()
	if err != nil {
		return nil, "", err
	}
	tfs, err := NewFS(f)
	if err != nil {
		return nil, "", err
	}
	return tfs, innerPath, err
}

func (w *wrap) ReadDir(dirname string) ([]os.FileInfo, error) {
	if !strings.Contains(dirname, suffix) {
		return w.inner.ReadDir(dirname)
	}
	tfs, innerPath, err := w.getTfs(dirname)
	if err != nil {
		return nil, err
	}
	return tfs.ReadDir(innerPath)
}

func (w *wrap) Lstat(name string) (os.FileInfo, error) {
	if !strings.Contains(name, suffix) {
		return w.inner.Lstat(name)
	}
	tfs, innerPath, err := w.getTfs(name)
	if err != nil {
		return nil, err
	}
	return tfs.Lstat(innerPath)
}

func (w *wrap) Join(elem ...string) string {
	return w.inner.ReadDir(name)
}

func (w *wrap) Open(name string) (io.ReadCloser, error) {
	fmt.Println("bp1")
	if !strings.Contains(name, suffix) {
		return w.inner.Open(name)
	}
	fmt.Println("bp1")
	tfs, innerPath, err := w.getTfs(name)
	if err != nil {
		return nil, err
	}
	fmt.Println("bp1")
	return tfs.Open(innerPath)
}
