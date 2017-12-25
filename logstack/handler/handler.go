package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/ws"
)

type Config struct {
	Root     string
	MarkFile string
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullpath := filepath.Join(c.Root, r.URL.Path)
	dir, _ := ioutil.ReadDir(r.URL.Path)
	if !fileInSlice(c.MarkFile, dir) {
		http.FileServer(http.Dir(fullpath))
		return
	}
	handlerConfig := config.Config{}
	for _, file := range dir {
		if file.Mode().IsDir() {
			u := url.URL{Scheme: "file://", Path: filepath.Join(fullpath + file.Name())}
			fs, _ := filesystem.NewLocalFS(&u)
			s := config.Src{
				Name: file.Name(),
				FS:   fs,
			}
			handlerConfig.Nodes = append(handlerConfig.Nodes, s)
		}
	}
	ws.New(handlerConfig).ServeHTTP(w, r)
}

func fileInSlice(filename string, list []os.FileInfo) bool {
	for _, b := range list {
		if !b.IsDir() && b.Name() == filename {
			return true
		}
	}
	return false
}
