package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/handler"
	"github.com/gorilla/mux"
)

var options struct {
	configFile string
}

func init() {
	flag.StringVar(&options.configFile, "config", "", "path to a config file")
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

	log.Printf("serving on http://localhost:8888")
	err = http.ListenAndServe(":8888", r)
	if err != nil {
		panic(err)
	}
}
