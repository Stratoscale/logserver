package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/router"
	"github.com/Stratoscale/logserver/source"
)

var log = logrus.WithField("pkg", "main")

const (
	defaultConfig = "logserver.json"
	port          = 8888
)

var options struct {
	port     int
	jsonFile string
	debug    bool
}

func init() {
	flag.IntVar(&options.port, "port", port, "Listen port")
	flag.StringVar(&options.jsonFile, "json", defaultConfig, "Path to a config json file")
	flag.BoolVar(&options.debug, "debug", false, "Show debug logs")
}

func main() {
	flag.Parse()

	if options.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	f, err := os.Open(options.jsonFile)
	failOnErr(err, fmt.Sprintf("open file %s", options.jsonFile))
	defer f.Close()

	var cf source.FileConfig
	err = json.NewDecoder(f).Decode(&cf)
	failOnErr(err, "decode file")

	c, err := source.New(cf)
	failOnErr(err, "creating config")

	defer c.CloseSources()

	log.Infof("serving on http://localhost:%d", options.port)
	rtr, err := router.New(*c)
	failOnErr(err, "creating router")

	m := http.NewServeMux()
	m.Handle("/", rtr)
	if options.debug {
		debug.PProfHandle(m)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", options.port), m)
	failOnErr(err, "serving")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", msg, err)
}
