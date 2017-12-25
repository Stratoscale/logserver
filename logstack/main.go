package main

import (
	"flag"
	"log"
	"net/http"

	"os"

	"regexp"

	"fmt"

	"github.com/Stratoscale/logserver/logstack/handler"
	"github.com/gorilla/mux"
)

var options struct {
	rootPath string
	port     int
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&options.rootPath, "root", cwd, "path to root directory")
	flag.IntVar(&options.port, "port", 8888, "port to listen on")
}

func main() {
	h := &handler.Config{
		Re:   regexp.MustCompile(""),
		Root: options.rootPath,
	}

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(h)

	log.Printf("serving on http://localhost:%d", options.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", options.port), r)
	if err != nil {
		panic(err)
	}
}
