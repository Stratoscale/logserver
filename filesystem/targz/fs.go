package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/filesystem"
)

var log = logrus.StandardLogger().WithField("pkg", "targz")

const sep = string(os.PathSeparator)

func NewFS(file filesystem.File) (*FileSystem, error) {
	fs := &FileSystem{
		index:  make(map[string]os.FileInfo),
		file:   file,
		Closer: file,
	}
	return fs, fs.init()
}

type FileSystem struct {
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
		f.index[strings.Trim(h.Name, sep)] = h.FileInfo()
		for _, dir := range subDirs(h.Name) {
			if _, ok := f.index[dir]; !ok {
				f.index[dir] = &fileInfo{name: filepath.Base(dir), isDir: true}
			}
		}
	}
	return nil
}

func subDirs(path string) []string {
	var subDirs []string
	dir, _ := filepath.Split(path)
	dir = strings.Trim(dir, sep)
	for dir != "" {
		subDirs = append(subDirs, dir)
		dir, _ = filepath.Split(dir)
		dir = strings.Trim(dir, sep)
	}
	return subDirs
}

// ReadDir implements the FileSystem ReadDir method,
// It returns a list of fileinfos in a given path
func (f *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	var files []os.FileInfo
	for path, file := range f.index {
		if isInDir(dirname, path) {
			files = append(files, file)
		}
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
	dirname = strings.Trim(dirname, sep)
	name = strings.Trim(name, sep)
	if !strings.HasPrefix(name, dirname) {
		return false
	}
	after := strings.Trim(name[len(dirname):], sep)
	return len(after) != 0 && !strings.Contains(after, sep)
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

type fileInfo struct {
	name  string
	isDir bool
}

func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) IsDir() bool        { return f.isDir }
func (f *fileInfo) Size() int64        { return 0 }
func (f *fileInfo) Mode() os.FileMode  { return 0666 }
func (f *fileInfo) ModTime() time.Time { return time.Time{} }
func (f *fileInfo) Sys() interface{}   { return nil }
