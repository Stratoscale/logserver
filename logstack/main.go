package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/logstack/handler"
	"github.com/bakins/logrus-middleware"
	"github.com/gorilla/handlers"
)

var log = logrus.WithField("pkg", "main")

var options struct {
	rootPath string
	port     int
	markFile string
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&options.rootPath, "root", cwd, "path to root directory")
	flag.StringVar(&options.markFile, "mark-file", "logstack.enable", "file that marks test root")
	flag.IntVar(&options.port, "port", 8889, "port to listen on")
}

type logger struct{}

func (logger) Write(p []byte) (n int, err error) {
	log.Debugf(string(p))
	return len(p), nil
}

func main() {
	flag.Parse()
	var h http.Handler = &handler.Config{
		Root:     options.rootPath,
		MarkFile: options.markFile,
	}

	logMW := logrusmiddleware.Middleware{Logger: log.Logger}
	h = logMW.Handler(handlers.LoggingHandler(logger{}, h), "")

	log.Infof("serving on http://localhost:%d", options.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", options.port), h)
	if err != nil {
		panic(err)
	}
}
