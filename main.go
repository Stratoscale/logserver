package main

import (
	"net/http"

	"log"

	"github.com/Stratoscale/logserver/handler"
	"github.com/gorilla/mux"
)

func main() {

	h := handler.New(handler.Config{})

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ws").Handler(h)

	log.Printf("serving on http://localhost:8888")
	err := http.ListenAndServe(":8888", r)
	if err != nil {
		panic(err)
	}
}
