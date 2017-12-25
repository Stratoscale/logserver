package handler

import (
	"log"
	"net/http"

	"fmt"

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

type Metadata struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
}

type request struct {
	Metadata `json:"meta"`
	BasePath path `json:"base_path"`
}

type path []string

type fileType string

const (
	file fileType = "file"
	dir  fileType = "dir"
)

type debugLevel string

const (
	levelDebug   debugLevel = "debug"
	levelInfo    debugLevel = "info"
	levelError   debugLevel = "error"
	levelWarning debugLevel = "warning"
)

type fsElement struct {
	Path path     `json:"path"`
	Type fileType `json:"type"`
	Size int      `json:"size"`
	FS   string   `json:"fs"`
}

type fileTreeResponse struct {
	Metadata `json:"meta"`
	Tree     []fsElement `json:"tree"`
}

type LogLine struct {
	Msg        string     `json:"msg"`
	Level      debugLevel `json:"level"`
	Time       string     `json:"time"`
	FS         string     `json:"fs"`
	FileName   string     `json:"file-name"`
	LineNumber int        `json:"line-number"`
	Offset     int        `json:"offset"`
}

type contentResponse struct {
	Metadata `json:"meta"`
	Lines    []LogLine `json:"line"`
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
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read: %s", err)
			return
		}
		fmt.Print("got: %s", string(b))
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
			Metadata: Metadata{ID: r.ID, Action: r.Action},
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
			Metadata: Metadata{ID: r.ID, Action: r.Action},
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
			Metadata: Metadata{ID: r.ID, Action: r.Action},
			Lines: []LogLine{
				{Msg: "bla bla bla", Level: levelDebug, FS: "node0", FileName: "bla.log", LineNumber: 1},
				{Msg: "bla bla", Level: levelDebug, FS: "node1", FileName: "bla.log", LineNumber: 100},
				{Msg: "harta barta", Level: levelWarning, FS: "node1", FileName: "harta.log", LineNumber: 1},
				{Msg: "harta barta", Level: levelInfo, FS: "node2", FileName: "harta.log", LineNumber: 7},
				{Msg: "panic error!", Level: levelError, FS: "node2", FileName: "harta.log", LineNumber: 7},
			},
		})
	}
}
