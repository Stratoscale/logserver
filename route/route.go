package route

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

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
	static        = http.FileServer(http.Dir("./client/dist"))
)

type Config struct {
	BasePath string `json:"base_path"`
	RootPath string `json:"root_path"`
}

func Static(r *mux.Router) {
	r.PathPrefix(pathStatic + "/").Handler(http.StripPrefix(pathStatic, static))
}

func Index(r *mux.Router, basePath string, c Config) error {

	if c.BasePath == "" && c.RootPath != "" {
		c.BasePath = c.RootPath
	}

	var index = bytes.NewBuffer(nil)
	if err := indexTemplate.Execute(index, c); err != nil {
		return fmt.Errorf("executing index template: %s", err)
	}

	r.Path(basePath).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(index.Bytes()); err != nil {
			log.WithError(err).Errorf("Writing index to response")
		}
	})
	return nil
}

func Engine(r *mux.Router, basePath string, engine http.Handler) {
	path := filepath.Join(basePath, pathWS)
	log.Debugf("Adding engine route on %s", path)
	r.Path(path).Handler(engine)
}

func Redirect(r *mux.Router, c Config) {
	if c.RootPath == "" {
		return
	}
	r.PathPrefix(c.RootPath + "/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dest := r.URL.Path[len(c.RootPath):]
		log.Printf("Redirecting to %s", dest)
		http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
	})
}
