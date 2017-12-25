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

type Metadata struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
}

type Request struct {
	Metadata `json:"meta"`
	BasePath pathArr `json:"base_path"`
}

type pathArr []string

type debugLevel string

const (
	levelDebug   debugLevel = "debug"
	levelInfo    debugLevel = "info"
	levelError   debugLevel = "error"
	levelWarning debugLevel = "warning"
)

type fsElement struct {
	Path  pathArr `json:"path"`
	IsDir bool    `json:"is_dir"`
	Size  int64   `json:"size"`
	FS    string  `json:"fs"`
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
	FileName   string     `json:"file_name"`
	LineNumber int        `json:"line_number"`
	Offset     int        `json:"offset"`
}

type contentResponse struct {
	Metadata `json:"meta"`
	Lines    []LogLine `json:"line"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ws Request from: %s", r.RemoteAddr)
	u := new(websocket.Upgrader)
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Request upgraded to: %s", r.RemoteAddr)

	for {
		var r Request
		err = conn.ReadJSON(&r)
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

func (h *handler) serve(w connWriter, r Request) {
	switch r.Action {
	case "get-file-tree":
		_ = r.BasePath
		//var fsElements []fsElement
		//for _, node := range h.Nodes {
		//	walker := fs.WalkFS(filepath.Join(r.BasePath...), node.FS)
		//	for walker.Step() {
		//		if err := walker.Err(); err != nil {
		//			log.Println(err)
		//			continue
		//		}
		//
		//		fsElements = append(fsElements, fsElement{
		//			Path:  filepath.SplitList(walker.Path()),
		//			IsDir: walker.Stat().IsDir(),
		//			Size:  walker.Stat().Size(),
		//			FS:    node.Name,
		//		})
		//	}
		//}
		//// reply
		//w.WriteJSON(&fileTreeResponse{
		//	Metadata: Metadata{ID: r.ID, Action: r.Action},
		//	Tree:     fsElements,
		//})

		// TODO: user basepath to get file system tree
		w.WriteJSON(&fileTreeResponse{
			Metadata: Metadata{ID: r.ID, Action: r.Action},
			Tree: []fsElement{
				{Path: []string{"var"}, IsDir: true, FS: "node0"},
				{Path: []string{"var"}, IsDir: true, FS: "node1"},
				{Path: []string{"var", "log"}, IsDir: true, FS: "node0"},
				{Path: []string{"var", "log"}, IsDir: true, FS: "node1"},
				{Path: []string{"var", "log", "mancala"}, IsDir: true, FS: "node1"},
				{Path: []string{"var", "log", "keystone.log"}, IsDir: false, Size: 10, FS: "node0"},
				{Path: []string{"var", "log", "keystone.log"}, IsDir: false, Size: 15, FS: "node1"},
				{Path: []string{"var", "log", "nova.log"}, IsDir: false, Size: 10},
			},
		})
	case "get-content":
		_ = r.BasePath
		//var logLines []LogLine
		//for _, node := range h.Nodes {
		//	walker := fs.WalkFS(filepath.Join(r.BasePath...), node.FS)
		//	for walker.Step() {
		//		if err := walker.Err(); err != nil {
		//			log.Println(err)
		//			continue
		//		}
		//
		//		logLines = append(LogLine, fsElement{
		//			Path:  filepath.SplitList(walker.Path()),
		//			IsDir: walker.Stat().IsDir(),
		//			Size:  walker.Stat().Size(),
		//			FS:    node.Name,
		//		})
		//	}
		//}

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
