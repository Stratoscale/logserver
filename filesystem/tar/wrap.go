package tar

import (
	"os"
	"regexp"
	"strings"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/bluele/gcache"
)

var (
	reContains = regexp.MustCompile(`\.tar(\.gz)?`)
	reSuffix   = regexp.MustCompile(`\.tar(\.gz)?$`)
)

// Wrap wraps a filesystem with a tar or tar.gz file opener
// so it a tar or tar.gz file is present in a filesystem, it's inner content can be
// browsed as a filesystem
func Wrap(inner filesystem.FileSystem, cache gcache.Cache, cachePrefix string) filesystem.FileSystem {
	return &tarfs{
		inner:       inner,
		cache:       cache,
		cachePrefix: cachePrefix,
	}
}

type tarfs struct {
	inner       filesystem.FileSystem
	cache       gcache.Cache
	cachePrefix string
}

func (w *tarfs) ReadDir(dirname string) ([]os.FileInfo, error) {
	tfs, innerPath, err := w.getTarFS(dirname)
	if err != nil {
		return nil, err
	}
	if tfs == nil {
		files, err := w.inner.ReadDir(dirname)
		if err != nil {
			return nil, err
		}
		return changeTarToDir(files...), nil
	}
	return tfs.ReadDir(innerPath)
}

func (w *tarfs) Lstat(name string) (os.FileInfo, error) {
	tfs, innerPath, err := w.getTarFS(name)
	if err != nil {
		return nil, err
	}
	if tfs == nil {
		file, err := w.inner.Lstat(name)
		if err != nil {
			return nil, err
		}
		return changeTarToDir(file)[0], nil
	}
	return tfs.Lstat(innerPath)
}

func (w *tarfs) Join(elem ...string) string {
	return w.inner.Join(elem...)
}

func (w *tarfs) Open(name string) (filesystem.File, error) {
	tfs, innerPath, err := w.getTarFS(name)
	if err != nil {
		return nil, err
	}
	if tfs == nil {
		return w.inner.Open(name)
	}
	return tfs.Open(innerPath)
}

func (w *tarfs) Close() error {
	return w.inner.Close()
}

type cacheKey string

func (w *tarfs) getTarFS(dirname string) (filesystem.FileSystem, string, error) {
	tarName, innerPath := split(dirname)
	if tarName == "" {
		return nil, dirname, nil
	}

	var (
		// key for storing tar files in cache
		key = cacheKey(w.cachePrefix + tarName)
		fs  filesystem.FileSystem
	)

	if val, err := w.cache.Get(key); err == nil {
		log.Debugf("Using cache for tar %s", key)
		fs = val.(filesystem.FileSystem)
	} else { // not in cache
		log.Infof("Opening new tar %s", key)
		f, err := w.inner.Open(tarName)
		if err != nil {
			return nil, "", err
		}
		fs, err = New(f)
		if err != nil {
			return nil, "", err
		}
		err = w.cache.Set(key, fs)
		if err != nil {
			log.WithError(err).Warn("Setting cache for %s", key)
		}
	}
	return fs, innerPath, nil
}

func split(dirname string) (tarName string, innerPath string) {
	loc := reContains.FindStringIndex(dirname)
	if len(loc) == 0 {
		return "", dirname
	}
	end := loc[1]

	tarName = dirname[:end]
	innerPath = strings.Trim(dirname[end:], sep)
	return
}

// changeTarToDir exposes tar files as directories
func changeTarToDir(files ...os.FileInfo) []os.FileInfo {
	for i, file := range files {
		if reSuffix.MatchString(file.Name()) {
			files[i] = &tarFile{FileInfo: file}
		}
	}
	return files
}

type tarFile struct{ os.FileInfo }

func (d *tarFile) IsDir() bool { return true }
