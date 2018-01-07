package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/router"
)

type Config struct {
	Root     string
	MarkFile string
}

func (c *Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	root, err := c.searchRoot(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if root == "" {
		http.FileServer(http.Dir(c.Root)).ServeHTTP(w, r)
		return
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cfg := config.FileConfig{}
	for _, file := range files {
		if file.Mode().IsDir() {
			cfg.Sources = append(cfg.Sources, config.SourceConfig{
				Name:         file.Name(),
				URL:          "file://" + filepath.Join(root, file.Name()),
				OpenTarFiles: true,
			})
		}
	}

	handlerCfg, err := config.New(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer handlerCfg.CloseSources()

	serverPath := root[len(c.Root):]
	handlerCfg.BasePath = serverPath
	rtr, err := router.New(*handlerCfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.StripPrefix(serverPath, rtr).ServeHTTP(w, r)
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
