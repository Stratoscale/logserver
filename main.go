package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/dynamic"
	"github.com/Stratoscale/logserver/engine"
	"github.com/Stratoscale/logserver/parse"
	"github.com/Stratoscale/logserver/router"
	"github.com/Stratoscale/logserver/source"
	"github.com/bakins/logrus-middleware"
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
	dynamic  bool
}

type config struct {
	Global  engine.Config   `json:"global"`
	Sources []source.Config `json:"sources"`
	Parsers []parse.Config  `json:"parsers"`
	Dynamic dynamic.Config  `json:"dynamic"`
}

func init() {
	flag.IntVar(&options.port, "port", port, "Listen port")
	flag.StringVar(&options.jsonFile, "json", defaultConfig, "Path to a config json file")
	flag.BoolVar(&options.debug, "debug", false, "Show debug logs")
	flag.BoolVar(&options.dynamic, "dynamic", false, "Run dynamic mod")
}

func main() {
	flag.Parse()

	if options.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg := loadConfig(options.jsonFile)

	parser, err := parse.New(cfg.Parsers)
	failOnErr(err, "creating parsers")

	var h http.Handler

	if !options.dynamic {
		s, err := source.New(cfg.Sources)
		failOnErr(err, "creating config")
		defer s.CloseSources()

		h, err = router.New(router.Config{
			Engine: engine.New(cfg.Global, s, parser),
		})
		failOnErr(err, "creating router")
	} else {
		var err error
		h, err = dynamic.New(cfg.Dynamic, cfg.Global, parser)
		failOnErr(err, "creating dynamic handler")
		logMW := logrusmiddleware.Middleware{Logger: log.Logger}
		h = logMW.Handler(h, "")
	}

	log.Infof("serving on http://localhost:%d", options.port)
	m := http.NewServeMux()
	m.Handle("/", h)
	if options.debug {
		debug.PProfHandle(m)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", options.port), m)
	failOnErr(err, "serving")
}

func loadConfig(fileName string) config {
	f, err := os.Open(fileName)
	failOnErr(err, fmt.Sprintf("open file %s", fileName))
	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	failOnErr(err, "decode file")
	return cfg
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", msg, err)
}

type logger struct{}

func (logger) Write(p []byte) (n int, err error) {
	log.Debugf(string(p))
	return len(p), nil
}
