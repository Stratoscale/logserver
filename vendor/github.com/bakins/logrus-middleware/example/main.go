package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/bakins/logrus-middleware"
)

func main() {

	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.JSONFormatter{}

	l := logrusmiddleware.Middleware{
		Name:   "example",
		Logger: logger,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "hello world\n")
	}

	http.Handle("/", l.Handler(http.HandlerFunc(handler), "homepage"))

	http.ListenAndServe(":8080", nil)
}
