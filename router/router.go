package router

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/source"
	"github.com/Stratoscale/logserver/ws"
	"github.com/gorilla/mux"
)

var (
	indexTemplate = template.Must(template.ParseFiles("./client/dist/index.html"))
	log           = logrus.WithField("pkg", "router")
)

func New(cfg source.Config) (http.Handler, error) {
	var (
		static = http.FileServer(http.Dir("./client/dist"))
	)
	index := bytes.NewBuffer(nil)
	if err := indexTemplate.Execute(index, &cfg); err != nil {
		return nil, fmt.Errorf("executing index template: %s", err)
	}

	serveIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(index.Bytes()); err != nil {
			log.WithError(err).Errorf("Writing index to response")
		}
	})

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(ws.New(cfg))
	r.Methods(http.MethodGet).Path("/").Handler(serveIndex)
	r.Methods(http.MethodGet).Path("/index.html").Handler(serveIndex)
	r.Methods(http.MethodGet).PathPrefix("/files").Handler(http.StripPrefix("/files", serveIndex))
	r.Methods(http.MethodGet).PathPrefix("/").Handler(static)

	return r, nil
}

type logger struct{}

func (logger) Write(b []byte) (int, error) {
	log.Printf(string(b))
	return len(b), nil
}
