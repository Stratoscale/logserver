package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Stratoscale/logserver/logstack/handler"
)

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

func main() {
	flag.Parse()
	h := &handler.Config{
		Root:     options.rootPath,
		MarkFile: options.markFile,
	}

	log.Printf("serving on http://localhost:%d", options.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", options.port), h)
	if err != nil {
		panic(err)
	}
}
