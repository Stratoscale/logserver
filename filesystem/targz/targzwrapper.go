package targz

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Stratoscale/logserver/filesystem"
)

var (
	reContains = regexp.MustCompile(`\.tar(\.gz)?`)
	reSuffix   = regexp.MustCompile(`\.tar(\.gz)?$`)
)

func New(inner filesystem.FileSystem) filesystem.FileSystem {
	return &wrap{inner: inner}
}

type wrap struct {
	inner filesystem.FileSystem
}

func (w *wrap) ReadDir(dirname string) ([]os.FileInfo, error) {
	if !reContains.MatchString(dirname) {
		files, err := w.inner.ReadDir(dirname)
		if err != nil {
			return nil, err
		}
		return changeTarToDir(files), nil
	}
	tfs, innerPath, err := w.getTfs(dirname)
	if err != nil {
		return nil, err
	}
	return tfs.ReadDir(innerPath)
}

// changeTarToDir exposes tar files as directories
func changeTarToDir(files []os.FileInfo) []os.FileInfo {
	for i, file := range files {
		if reSuffix.MatchString(file.Name()) {
			files[i] = &tarFile{FileInfo: file}
		}
	}
	return files
}

func (w *wrap) Lstat(name string) (os.FileInfo, error) {
	if !reContains.MatchString(name) {
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

func (w *wrap) Close() error {
	return w.inner.Close()
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
	loc := reContains.FindStringIndex(dirname)
	if len(loc) == 0 {
		return "", dirname
	}
	end := loc[1]

	tarName = dirname[:end]
	innerPath = strings.Trim(dirname[end:], string(os.PathSeparator))
	return
}

type tarFile struct{ os.FileInfo }

func (d *tarFile) IsDir() bool { return true }
