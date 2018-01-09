package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
)

var log = logrus.StandardLogger().WithField("pkg", "targz")

func NewFS(r io.ReadCloser) (*FileSystem, error) {
	var tarReader *tar.Reader

	if z, err := gzip.NewReader(r); err == nil {
		tarReader = tar.NewReader(z)
	} else {
		tarReader = tar.NewReader(r)
	}
	return &FileSystem{
		Reader: tarReader,
		Closer: r,
	}, nil
}

type FileSystem struct {
	Reader *tar.Reader
	Closer io.Closer
}

// ReadDir implements the FileSystem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	var content []os.FileInfo
	for {
		h, err := f.Reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !isInDir(dirname, h.Name) {
			continue
		}
		content = append(content, h.FileInfo())
	}
	return content, nil
}

// Lstat implements the FileSystem Lstat method,
// it returns fileinfo for a given path
func (f *FileSystem) Lstat(name string) (os.FileInfo, error) {
	defer debug.Time(log, "Stat %s", name)()
	for {
		h, err := f.Reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.Name != name {
			continue
		}
		return h.FileInfo(), nil
	}
	return nil, fmt.Errorf("not found: %s", name)
}

// Join implements the FileSystem Join method,
func (f *FileSystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *FileSystem) Open(name string) (io.ReadCloser, error) {
	defer debug.Time(log, "Opened: %s", name)()
	_, err := f.Lstat(name)
	if err != nil {
		return nil, err
	}
	return &readCloser{Reader: f.Reader, Closer: f.Closer}, nil
}

func (f *FileSystem) Close() error {
	return f.Closer.Close()
}

type readCloser struct {
	io.Reader
	io.Closer
}

func isInDir(dirname, name string) bool {
	if !strings.HasPrefix(name, dirname) {
		return false
	}
	after := name[len(dirname):]
	if strings.Contains(strings.Trim(after, string(os.PathSeparator)), string(os.PathListSeparator)) {
		return false
	}
	return true

}
