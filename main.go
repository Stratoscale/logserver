package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/cache"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/download"
	"github.com/Stratoscale/logserver/dynamic"
	"github.com/Stratoscale/logserver/engine"
	"github.com/Stratoscale/logserver/parse"
	"github.com/Stratoscale/logserver/route"
	"github.com/Stratoscale/logserver/source"
	"github.com/bakins/logrus-middleware"
	"github.com/gorilla/mux"
)

var log = logrus.WithField("pkg", "main")

const (
	defaultConfig = "logserver.json"
	defaultAddr   = "localhost:8888"
)

var options struct {
	addr    string
	config  string
	debug   bool
	dynamic bool
}

func init() {
	flag.StringVar(&options.addr, "addr", defaultAddr, "Serving address")
	flag.StringVar(&options.config, "config", defaultConfig, "Path to a config file")
	flag.BoolVar(&options.debug, "debug", false, "Show debug logs")
	flag.BoolVar(&options.dynamic, "dynamic", false, "Run in dynamic mode")
}

type config struct {
	Global  engine.Config   `json:"global"`
	Sources []source.Config `json:"sources"`
	Parsers []parse.Config  `json:"parsers"`
	Dynamic dynamic.Config  `json:"dynamic"`
	Cache   cache.Config    `json:"cache"`
	Route   route.Config    `json:"route"`
}

func (c config) journal() string {
	if name := c.Dynamic.OpenJournal; name != "" {
		return name
	}
	for _, src := range c.Sources {
		if name := src.OpenJournal; name != "" {
			return name
		}
	}
	return ""
}

func main() {
	flag.Parse()

	// apply debug logs
	if options.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// validate address
	_, _, err := net.SplitHostPort(options.addr)
	failOnErr(err, "Bad address value: %s", options.addr)

	cfg := loadConfig(options.config)

	log.Infof("Loading parsers...")
	parser, err := parse.New(cfg.Parsers)
	failOnErr(err, "Creating parsers")

	// add journal parser if necessary
	if journalName := cfg.journal(); journalName != "" {
		log.Infof("Adding a journalctl parser")
		err := parser.AppendJournal(journalName)
		if err != nil {
			log.WithError(err).Warn("Failed adding a journalctl parser")
		}
	}

	log.Printf("Loaded with %d parsers", len(parser))

	cache := cache.New(cfg.Cache)

	r := mux.NewRouter()
	route.Static(r)

	if !options.dynamic {

		s, err := source.New(cfg.Sources, cache)
		failOnErr(err, "Creating config")
		defer s.CloseSources()

		dl := download.New(filepath.Join(cfg.Route.RootPath, "_dl"), s, cache)
		eng := engine.New(cfg.Global, s, parser, cache)

		// put websocket handler behind the root and behind the proxy path
		// it must be before the redirect handlers because it is on the proxy path
		route.Engine(r, "/", eng)
		route.Download(r, "/", dl)

		if cfg.Route.RootPath != "" && cfg.Route.RootPath != "/" {
			route.Engine(r, cfg.Route.RootPath, eng)
			route.Download(r, cfg.Route.RootPath, dl)
		}

		// add redirect of request that are sent to a proxy path with the same URL without the proxy prefix
		route.Redirect(r, cfg.Route)

		// handle with index on any route that does not match anything else
		failOnErr(route.Index(r, "/", cfg.Route), "Creating index")

	} else {
		var err error
		h, err := dynamic.New(cfg.Dynamic, cfg.Global, parser, cache)
		failOnErr(err, "Creating dynamic handler")
		logMW := logrusmiddleware.Middleware{Logger: log.Logger}
		h = logMW.Handler(h, "")
		r.PathPrefix("/").Handler(h)
	}

	// add debug handlers
	if options.debug {
		debug.PProfHandle(r)
	}

	log.Infof("Serving on http://%s", options.addr)
	err = http.ListenAndServe(options.addr, r)
	failOnErr(err, "Serving")
}

func loadConfig(fileName string) config {
	f, err := os.Open(fileName)
	failOnErr(err, fmt.Sprintf("open file %s", fileName))
	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	failOnErr(err, "Decode config file")
	return cfg
}

func failOnErr(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	log.Fatalf("%s: %s", fmt.Sprintf(msg, args...), err)
}
