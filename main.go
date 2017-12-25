package main

import (
	"fmt"
	"log"
	"net/http"

	"flag"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/handler"
	"github.com/gorilla/mux"
)

var options struct {
	port       int
	configFile string
}

func init() {
	flag.IntVar(&options.port, "port", 8888, "Listen port")
	flag.StringVar(&options.configFile, "config", "", "Path to a config file")
}

func main() {
	c, err := config.New([]config.SrcDesc{
		{Name: "src1", Address: "file://./example/log1"},
		{Name: "src2", Address: "file://./example/log2"},
	})
	if err != nil {
		panic(err)
	}

	h := handler.New(*c)
	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(h)

	log.Printf("serving on http://localhost:%d", options.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", options.port), r)
	if err != nil {
		panic(err)
	}
}
