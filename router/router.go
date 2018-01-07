package router

import (
	"net/http"
	"os"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/ws"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func New(cfg config.Config) http.Handler {
	var static = handlers.LoggingHandler(os.Stderr, http.FileServer(http.Dir("./client/dist")))

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(ws.New(cfg))
	r.Methods(http.MethodGet).PathPrefix("/files").Handler(http.StripPrefix("/files", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./client/dist/index.html") })))
	r.Methods(http.MethodGet).PathPrefix("/").Handler(static)
	return r
}
