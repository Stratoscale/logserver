package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/logstack/handler"
	"github.com/bakins/logrus-middleware"
	"github.com/gorilla/handlers"
)

var log = logrus.WithField("pkg", "main")

var options struct {
	rootPath string
	port     int
	markFile string
	debug    bool
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&options.rootPath, "root", cwd, "Path to root directory")
	flag.StringVar(&options.markFile, "mark-file", "logstack.enable", "File that marks test root")
	flag.IntVar(&options.port, "port", 8889, "Port to listen on")
	flag.BoolVar(&options.debug, "debug", false, "Show debug logs")
}

type logger struct{}

func (logger) Write(p []byte) (n int, err error) {
	log.Debugf(string(p))
	return len(p), nil
}

func main() {
	flag.Parse()

	if options.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var h http.Handler = &handler.Config{
		Root:     options.rootPath,
		MarkFile: options.markFile,
	}

	logMW := logrusmiddleware.Middleware{Logger: log.Logger}
	h = logMW.Handler(handlers.LoggingHandler(logger{}, h), "")

	m := http.NewServeMux()
	m.Handle("/", h)
	if options.debug {
		debug.PProfHandle(m)
	}

	log.Infof("serving on http://localhost:%d", options.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", options.port), m)
	if err != nil {
		panic(err)
	}
}
