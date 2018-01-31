package route

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	pathStatic = "/_static"
	pathWS     = "/_ws"
)

var (
	indexTemplate = template.Must(template.ParseFiles("./client/dist/index.html"))
	log           = logrus.WithField("pkg", "router")
)

type Config struct {
	BasePath string `json:"base_path"`
	RootPath string `json:"root_path"`
}

func Static(r *mux.Router, c Config) {
	var (
		static = http.FileServer(http.Dir("./client/dist"))
		path   = filepath.Join(c.RootPath, pathStatic)
	)
	log.Infof("Adding static file serving on %s", path)
	r.PathPrefix(path + "/").Handler(http.StripPrefix(path, static))
}

func Index(r *mux.Router, c Config) error {

	if c.BasePath == "" && c.RootPath != "" {
		c.BasePath = c.RootPath
	}

	var index = bytes.NewBuffer(nil)
	if err := indexTemplate.Execute(index, c); err != nil {
		return fmt.Errorf("executing index template: %s", err)
	}

	path := c.RootPath
	if len(path) == 0 || path[len(path)-1] != '/' {
		path += "/"
	}

	log.Infof("Adding index route on %s", path)
	r.Path(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(index.Bytes()); err != nil {
			log.WithError(err).Errorf("Writing index to response")
		}
	})
	return nil
}

func Engine(r *mux.Router, c Config, engine http.Handler) {
	path := filepath.Join(c.RootPath, pathWS)
	log.Debugf("Adding engine route on %s", path)
	r.Path(path).Handler(engine)
}

func Redirect(r *mux.Router, c Config) {
	if c.RootPath == "" {
		return
	}
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, c.RootPath) {
			http.NotFound(w, r)
			return
		}
		dest := filepath.Join(c.RootPath, r.URL.Path)
		if len(r.URL.Path) == 0 || r.URL.Path[len(r.URL.Path)-1] == '/' {
			dest += "/"
		}
		log.Printf("Redirecting to %s", dest)
		http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
	})
}
