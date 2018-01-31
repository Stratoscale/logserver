package dynamic

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Stratoscale/logserver/engine"
	"github.com/Stratoscale/logserver/parse"
	"github.com/Stratoscale/logserver/route"
	"github.com/Stratoscale/logserver/source"
	"github.com/bluele/gcache"
	"github.com/gorilla/mux"
)

const defaultMarkFile = "logstack.enable"

// Config is dynamic configuration
type Config struct {
	Root     string `json:"root"`
	MarkFile string `json:"mark_file"`
	source.Flags
}

func New(c Config, engineCfg engine.Config, routeCfg route.Config, p parse.Parse, cache gcache.Cache) (http.Handler, error) {
	var err error
	c.Root, err = filepath.Abs(c.Root)
	if err != nil {
		return nil, err
	}
	h := &handler{
		Config:    c,
		parse:     p,
		cache:     cache,
		engineCfg: engineCfg,
		routeCfg:  routeCfg,
	}
	if h.MarkFile == "" {
		h.MarkFile = defaultMarkFile
	}
	return h, nil
}

type handler struct {
	Config
	parse     parse.Parse
	cache     gcache.Cache
	route     route.Config
	engineCfg engine.Config
	routeCfg  route.Config
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
				Name:  file.Name(),
				URL:   "file://" + filepath.Join(root, file.Name()),
				Flags: h.Config.Flags,
			})
		}
	}
	src, err := source.New(srcConfig, h.cache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer src.CloseSources()

	serverPath := root[len(h.Root):]
	rtr := mux.NewRouter()

	// add index.html serving at the serverPath which is the dynamic root
	err = route.Index(rtr, route.Config{
		// BasePath is used to determined the websocket path
		BasePath: filepath.Join(h.routeCfg.RootPath, serverPath),
		// RootPath is for taking the static files, which are served by the root handler
		RootPath: ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add websocket handler on the server root
	route.Engine(
		rtr,
		route.Config{RootPath: ""},
		engine.New(h.engineCfg, src, h.parse, h.cache),
	)

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
