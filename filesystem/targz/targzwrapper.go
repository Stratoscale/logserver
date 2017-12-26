package targz

import (
	"io"
	"os"
	"strings"

	"fmt"

	"github.com/Stratoscale/logserver/filesystem"
)

const suffix = ".tar.gz"

func New(inner filesystem.FileSystem) filesystem.FileSystem {
	return &wrap{inner: inner}
}

type wrap struct {
	inner filesystem.FileSystem
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
	return w.inner.Join(elem...)
}

func (w *wrap) Open(name string) (io.ReadCloser, error) {
	tfs, innerPath, err := w.getTfs(name)
	if err != nil {
		return nil, err
	}
	if tfs == nil { // we are outside the tar filesystem
		return w.inner.Open(name)
	}
	if innerPath == "" {
		return nil, fmt.Errorf("no such file ''")
	}
	return tfs.Open(innerPath)
}

func (w *wrap) getTfs(dirname string) (filesystem.FileSystem, string, error) {
	tarName, innerPath := split(dirname)
	if tarName == "" {
		return nil, dirname, nil
	}
	f, err := w.inner.Open(tarName)
	if err != nil {
		return nil, "", err
	}
	tfs, err := NewFS(f)
	if err != nil {
		return nil, "", err
	}
	return tfs, innerPath, err
}

func split(dirname string) (tarName string, innerPath string) {
	i := strings.Index(dirname, suffix)
	if i == -1 {
		return "", dirname
	}

	i += len(suffix)

	tarName = dirname[:i]
	innerPath = strings.Trim(dirname[i:], string(os.PathSeparator))
	return
}
