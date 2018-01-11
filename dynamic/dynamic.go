package dynamic

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Stratoscale/logserver/engine"
	"github.com/Stratoscale/logserver/parse"
	"github.com/Stratoscale/logserver/router"
	"github.com/Stratoscale/logserver/source"
)

const (
	defaultMarkFile = "logstack.enable"
)

type Config struct {
	Root     string `json:"root"`
	MarkFile string `json:"mark_file"`
}

func New(c Config, engineConfig engine.Config, p parse.Parse) (http.Handler, error) {
	var err error
	c.Root, err = filepath.Abs(c.Root)
	if err != nil {
		return nil, err
	}
	h := &handler{
		Config:       c,
		parse:        p,
		engineConfig: engineConfig,
	}
	if h.MarkFile == "" {
		h.MarkFile = defaultMarkFile
	}
	return h, nil
}

type handler struct {
	Config
	parse        parse.Parse
	engineConfig engine.Config
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	root, err := h.searchRoot(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if root == "" {
		http.FileServer(http.Dir(h.Root)).ServeHTTP(w, r)
		return
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var srcConfig []source.Config
	for _, file := range files {
		if file.Mode().IsDir() {
			srcConfig = append(srcConfig, source.Config{
				Name:         file.Name(),
				URL:          "file://" + filepath.Join(root, file.Name()),
				OpenTarFiles: true,
			})
		}
	}
	src, err := source.New(srcConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer src.CloseSources()

	serverPath := root[len(h.Root):]

	rtr, err := router.New(router.Config{
		Engine:   engine.New(h.engineConfig, src, h.parse),
		BasePath: serverPath,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.StripPrefix(serverPath, rtr).ServeHTTP(w, r)
}

func (h *handler) searchRoot(path string) (string, error) {
	fullPath := ""
	parts := strings.Split(path, string(os.PathSeparator))
	parts = append([]string{h.Root}, parts...)
	for _, part := range parts {
		fullPath = filepath.Join(fullPath, part)
		isRootDir, err := h.markerDir(fullPath)
		if err != nil {
			return "", err
		}
		if isRootDir {
			return fullPath, nil
		}
	}
	return "", nil
}

func (h *handler) markerDir(dir string) (bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if !f.IsDir() && f.Name() == h.MarkFile {
			return true, nil
		}
	}
	return false, nil
}
