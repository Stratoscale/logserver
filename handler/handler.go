package handler

import (
	"log"
	"net/http"

	"github.com/Stratoscale/logserver/config"
	"github.com/gorilla/websocket"
)

func New(c config.Config) http.Handler {
	return &handler{
		Config: c,
	}
}

type handler struct {
	config.Config
}

type request struct {
	ID       int    `json:"id"`
	Action   string `json:"action"`
	BasePath path   `json:"base_path"`
}

type path []string

type fileType string

const (
	file fileType = "file"
	dir  fileType = "dir"
)

type fsElement struct {
	Path path     `json:"path"`
	Type fileType `json:"type"`
	Size int      `json:"size"`
	FS   string   `json:"fs"`
}

type response struct {
	ID int `json:"id"`
}

type fileTreeResponse struct {
	response
	Tree []fsElement `json:"tree"`
}

type LogLine struct {
	Msg string `json:"msg"`
	Level string `json:"level"`
	Time string `json:"time"`
	FS string `json:"fs"`
	FileName string `json:"file-name"`
	LineNumber int `json:"line-number"`
	Offset int `json:"offset"`
}

type contentResponse struct {
	response
	Lines []LogLine `json:"line"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ws request from: %s", r.RemoteAddr)
	u := new(websocket.Upgrader)
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Request upgraded to: %s", r.RemoteAddr)

	for {
		var r request
		err := conn.ReadJSON(&r)
		if err != nil {
			log.Printf("read: %s", err)
			return
		}
		go h.serve(conn, r)
	}
}

type connWriter interface {
	WriteJSON(interface{}) error
}

func (h *handler) serve(w connWriter, r request) {
	switch r.Action {
	case "get-file-tree":
		_ = r.BasePath

		// TODO: user basepath to get file system tree
		w.WriteJSON(&fileTreeResponse{
			response: response{ID: r.ID},
			Tree: []fsElement{
				{Path: []string{"var"}, Type: dir, FS: "node0"},
				{Path: []string{"var"}, Type: dir, FS: "node1"},
				{Path: []string{"var", "log"}, Type: dir, FS: "node0"},
				{Path: []string{"var", "log"}, Type: dir, FS: "node1"},
				{Path: []string{"var", "log", "mancala"}, Type: dir, FS: "node1"},
				{Path: []string{"var", "log", "keystone.log"}, Type: file, Size: 10, FS: "node0"},
				{Path: []string{"var", "log", "keystone.log"}, Type: file, Size: 15, FS: "node1"},
				{Path: []string{"var", "log", "nova.log"}, Type: file, Size: 10},
			},
		})
	case "get-content":
		_ = r.BasePath
		// TODO: user basepath to get file system tree
		w.WriteJSON(&contentResponse{
			response: response{ID: r.ID},
			Lines: []LogLine{
				{Msg: "bla bla bla", Level: "debug", FS: "node0", FileName: "bla.log", LineNumber: 1},
				{Msg: "bla bla", Level: "debug", FS: "node1", FileName: "bla.log", LineNumber: 100},
				{Msg: "harta barta", Level: "info", FS: "node1", FileName: "harta.log", LineNumber: 1},
				{Msg: "harta barta", Level: "info", FS: "node2", FileName: "harta.log", LineNumber: 7},
			},
		})
	case "search":
		_ = r.BasePath
		// TODO: user basepath to get file system tree
		w.WriteJSON(&contentResponse{
			response: response{ID: r.ID},
			Lines: []LogLine{
				{Msg: "bla bla bla", Level: "debug", FS: "node0", FileName: "bla.log", LineNumber: 1},
				{Msg: "bla bla", Level: "debug", FS: "node1", FileName: "bla.log", LineNumber: 100},
				{Msg: "harta barta", Level: "info", FS: "node1", FileName: "harta.log", LineNumber: 1},
				{Msg: "harta barta", Level: "info", FS: "node2", FileName: "harta.log", LineNumber: 7},
			},
		})
	}
}
