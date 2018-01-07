package router

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/ws"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	indexTemplate = template.Must(template.ParseFiles("./client/dist/index.html"))
	log           = logrus.StandardLogger().WithField("pkg", "router")
)

func New(cfg config.Config) (http.Handler, error) {
	var (
		static = handlers.LoggingHandler(os.Stderr, http.FileServer(http.Dir("./client/dist")))
	)
	index := bytes.NewBuffer(nil)
	err := indexTemplate.Execute(index, cfg)
	if err != nil {
		return nil, fmt.Errorf("executing index template: %s", err)
	}

	serveIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(index.Bytes()); err != nil {
			log.WithError(err).Errorf("Writing index")
		}
	})

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(ws.New(cfg))
	r.Methods(http.MethodGet).PathPrefix("/files").Handler(http.StripPrefix("/files", serveIndex))
	r.Methods(http.MethodGet).PathPrefix("/").Handler(static)
	return r, nil
}
