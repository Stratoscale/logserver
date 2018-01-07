package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/router"
)

var log = logrus.WithField("pkg", "main")

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

	log.Infof("serving on http://localhost:%d", options.port)
	rtr, err := router.New(*c)
	failOnErr(err, "creating router")
	err = http.ListenAndServe(fmt.Sprintf(":%d", options.port), rtr)
	failOnErr(err, "serving")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", msg, err)
}
