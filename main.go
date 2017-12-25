package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/ws"
	"github.com/gorilla/mux"
)

const (
	defaultConfig = "logserver.json"
	port          = 8888
)

var options struct {
	port     int
	jsonFile string
}

func init() {
	flag.IntVar(&options.port, "port", port, "Listen port")
	flag.StringVar(&options.jsonFile, "json", defaultConfig, "Path to a config json file")
}

func main() {
	f, err := os.Open(options.jsonFile)
	failOnErr(err, fmt.Sprintf("open file %s", options.jsonFile))
	defer f.Close()

	var cf config.FileConfig
	err = json.NewDecoder(f).Decode(&cf)
	failOnErr(err, "decode file")

	handlerConfig, err := config.New(cf)
	failOnErr(err, "creating config")

	log.Printf("serving on http://localhost:%d", options.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", options.port), router(*handlerConfig))
	failOnErr(err, "serving")
}

func router(cfg config.Config) http.Handler {
	var static = http.FileServer(http.Dir("./client/dist"))

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(ws.New(cfg))
	r.Methods(http.MethodGet).PathPrefix("/").Handler(static)
	return r
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", msg, err)
}
