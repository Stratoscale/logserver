package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/handler"
)

type Config struct {
	Re   *regexp.Regexp
	Root string
}

func fileInSlice(filename string, list []os.FileInfo) bool {
	for _, b := range list {
		if !b.IsDir() && b.Name() == filename {
			return true
		}
	}
	return false
}

func (l *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fullpath := filepath.Join(l.Root, r.URL.Path)

	dir, _ := ioutil.ReadDir(r.URL.Path)

	if !fileInSlice("logstack.enable", dir) {
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
	handler.New(handlerConfig).ServeHTTP(w, r)

}
