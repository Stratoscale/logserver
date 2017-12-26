package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/router"
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
	flag.Parse()
	f, err := os.Open(options.jsonFile)
	failOnErr(err, fmt.Sprintf("open file %s", options.jsonFile))
	defer f.Close()

	var cf config.FileConfig
	err = json.NewDecoder(f).Decode(&cf)
	failOnErr(err, "decode file")

	c, err := config.New(cf)
	failOnErr(err, "creating config")

	defer c.CloseSources()

	log.Printf("serving on http://localhost:%d", options.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", options.port), router.New(*c))
	failOnErr(err, "serving")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", msg, err)
}
