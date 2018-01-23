package router

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	indexTemplate = template.Must(template.ParseFiles("./client/dist/index.html"))
	log           = logrus.WithField("pkg", "router")
)

type Config struct {
	BasePath string
	Engine   http.Handler
}

func New(c Config) (http.Handler, error) {
	var static = http.FileServer(http.Dir("./client/dist"))
	index := bytes.NewBuffer(nil)
	if err := indexTemplate.Execute(index, c); err != nil {
		return nil, fmt.Errorf("executing index template: %s", err)
	}

	serveIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(index.Bytes()); err != nil {
			log.WithError(err).Errorf("Writing index to response")
		}
	})

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/_ws").Handler(c.Engine)
	r.Methods(http.MethodGet).PathPrefix("/_static").Handler(http.StripPrefix("/_static", static))
	r.Methods(http.MethodGet).PathPrefix("/").Handler(serveIndex)

	return r, nil
}

type logger struct{}

func (logger) Write(b []byte) (int, error) {
	log.Printf(string(b))
	return len(b), nil
}
