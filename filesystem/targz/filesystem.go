package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/filesystem"
)

var log = logrus.StandardLogger().WithField("pkg", "targz")

func NewFS(file filesystem.File) (*FileSystem, error) {
	fs := &FileSystem{
		index:  make(map[string]os.FileInfo),
		file:   file,
		Closer: file,
	}
	return fs, fs.init()
}

type FileSystem struct {
	list   []os.FileInfo
	index  map[string]os.FileInfo
	file   filesystem.File
	Closer io.Closer
}

func (f *FileSystem) init() error {
	tarReader := f.tarReader()
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f.list = append(f.list, h.FileInfo())
		f.index[h.Name] = f.list[len(f.list)-1]
	}
	return nil
}

// ReadDir implements the FileSystem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	files := make([]os.FileInfo, 0, len(f.index))
	for path, file := range f.index {
		if !isInDir(dirname, path) {
			continue
		}
		files = append(files, file)
	}
	return files, nil
}

// Lstat implements the FileSystem Lstat method,
// it returns fileinfo for a given path
func (f *FileSystem) Lstat(name string) (os.FileInfo, error) {
	file := f.index[name]
	if file == nil {
		return nil, notFound(name)
	}
	return file, nil
}

// Join implements the FileSystem Join method,
func (f *FileSystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (f *FileSystem) Open(name string) (filesystem.File, error) {
	if name == "" {
		return nil, notFound(name)
	}
	defer debug.Time(log, "Opened: %s", name)()

	if _, ok := f.index[name]; !ok {
		return nil, notFound(name)
	}

	f.file.Seek(0, io.SeekStart)
	tarReader := f.tarReader()
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.Name == name {
			return &file{ReadCloser: ioutil.NopCloser(tarReader), Seeker: f.file}, nil
		}
	}
	return nil, notFound(name)
}

func (f *FileSystem) Close() error {
	return f.Closer.Close()
}

func isInDir(dirname, name string) bool {
	const sep = string(os.PathSeparator)
	if !strings.HasPrefix(name, dirname) {
		return false
	}
	after := name[len(dirname):]
	return !strings.Contains(strings.Trim(after, sep), sep)
}

func (f *FileSystem) tarReader() *tar.Reader {
	if z, err := gzip.NewReader(f.file); err == nil {
		return tar.NewReader(z)
	}
	return tar.NewReader(f.file)
}

type file struct {
	io.ReadCloser
	io.Seeker
}

func notFound(name string) error {
	return fmt.Errorf("not found: '%s'", name)
}
