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
	// RootPath is defined for changing the root path of serving. This is useful
	// if the server is behind a proxy that changes the root path.
	RootPath string `json:"root_path"`
	// BasePath is to change the base path after the root path.
	// It is used for dynamic mode where we have different locations for the index page.
	BasePath string `json:"base_path"`
}

// Static serves static files
func Static(r *mux.Router) {
	r.PathPrefix(pathStatic + "/").Handler(http.StripPrefix(pathStatic, static))
}

// Index mounts serving of index.html on a path prefix.
// It uses a prefix since reloads of a page should give serving of the index.html page with the same url
// for the javascript frontend.
func Index(r *mux.Router, pathPrefix string, c Config) error {

	if c.BasePath == "" && c.RootPath != "" {
		c.BasePath = c.RootPath
	}

	var index = bytes.NewBuffer(nil)
	if err := indexTemplate.Execute(index, c); err != nil {
		return fmt.Errorf("executing index template: %s", err)
	}

	r.PathPrefix(pathPrefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(index.Bytes()); err != nil {
			log.WithError(err).Errorf("Writing index to response")
		}
	})
	return nil
}

// Engine mounts the websocket handler on the router
func Engine(r *mux.Router, basePath string, engine http.Handler) {
	path := filepath.Join(basePath, pathWS)
	log.Debugf("Adding engine route on %s", path)
	r.Path(path).Handler(engine)
}

// Redirect mounts a redirect handler for a proxy on the router
func Redirect(r *mux.Router, c Config) {
	if c.RootPath == "" {
		return
	}
	r.PathPrefix(c.RootPath + "/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dest := r.URL.Path[len(c.RootPath):]
		http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
	})
}
