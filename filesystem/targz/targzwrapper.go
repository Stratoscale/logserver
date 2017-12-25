package tarfswrap

import (
    "io"
    "os"

    "strings"

    "github.com/Stratoscale/logserver/config"
    "github.com/posener/tarfs"
)

const suffix = ".tar.gz"

func New(inner config.FileSystem) config.FileSystem {
    return &wrap{inner: inner}
}

type wrap struct {
    inner config.FileSystem
}

func (w *wrap) ReadDir(dirname string) ([]os.FileInfo, error) {
    if !strings.Contains(dirname, suffix) {
        return w.inner.ReadDir(dirname)
    }
    tarName, innerPath := split(dirname)
    f, err := w.inner.Open(tarName)
    defer f.Close()
    if err != nil {
        return nil, err
    }
    tfs, err := tarfs.New(f)
    if err != nil {
        return nil, err
    }
    return tfs.ReadDir(innerPath)
}

func (w *wrap) Lstat(name string) (os.FileInfo, error) {
    panic("implement me")
}

func (w *wrap) Join(elem ...string) string {
    panic("implement me")
}

func (w *wrap) Open(path string) (io.ReadCloser, error) {
    panic("implement me")
}
