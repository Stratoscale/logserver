package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"strings"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/filesystem/targz"
	"github.com/Stratoscale/logserver/router"
)

type Config struct {
	Root     string
	MarkFile string
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request for %s", r.URL.Path)

	root, err := c.searchRoot(r.URL.Path)
	if err != nil {
		log.Printf("not found: %s", err)
		http.NotFound(w, r)
		return
	}
	if root == "" {
		log.Printf("serving regular file: %s", r.URL.Path)
		http.FileServer(http.Dir(c.Root)).ServeHTTP(w, r)
		return
	}

	log.Printf("Serving root: %s", root)
	files, err := ioutil.ReadDir(root)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	handlerConfig := config.Config{}
	for _, file := range files {
		if file.Mode().IsDir() {
			u := url.URL{Scheme: "file://", Path: filepath.Join(root, file.Name())}
			var (
				fs  filesystem.FileSystem
				err error
			)
			fs, err = filesystem.NewLocalFS(&u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Opened FS on %s", u.Path)
			fs = targz.New(fs)
			handlerConfig.Sources = append(handlerConfig.Sources, config.Source{Name: file.Name(), FS: fs})
		}
	}
	defer handlerConfig.CloseSources()
	log.Printf("stripping: %s", root[len(c.Root):])
	http.StripPrefix(root[len(c.Root):], router.New(handlerConfig)).ServeHTTP(w, r)
}

func (c *Config) searchRoot(path string) (string, error) {
	fullPath := ""
	parts := strings.Split(path, string(os.PathSeparator))
	parts = append([]string{c.Root}, parts...)
	for _, part := range parts {
		fullPath = filepath.Join(fullPath, part)
		isRootDir, err := c.markerDir(fullPath)
		if err != nil {
			return "", err
		}
		if isRootDir {
			return fullPath, nil
		}
	}
	return "", nil
}

func (c *Config) markerDir(dir string) (bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if !f.IsDir() && f.Name() == c.MarkFile {
			return true, nil
		}
	}
	return false, nil
}
