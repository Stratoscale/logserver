package handler

import (
	"log"
	"net/http"
	"path/filepath"

	"bufio"

	"github.com/Stratoscale/logserver/config"
	"github.com/gorilla/websocket"
	"github.com/kr/fs"
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
	Path     pathArr `json:"path"`
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
	Key       string         `json:"key"`
	Path      pathArr        `json:"path"`
	IsDir     bool           `json:"is_dir"`
	Instances []fileInstance `json:"instances"`
}

type fileInstance struct {
	Size int64  `json:"size"`
	FS   string `json:"fs"`
}

type ResponseFileTree struct {
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
	path := filepath.Join(r.Path...)
	if path == "" {
		path = "/"
	}
	switch r.Action {
	case "get-file-tree":
		var (
			fsElements []fsElement
			m          = make(map[string]*fsElement)
		)

		for _, node := range h.Nodes {
			walker := fs.WalkFS(path, node.FS)
			for walker.Step() {
				if err := walker.Err(); err != nil {
					log.Println(err)
					continue
				}

				key := walker.Path()
				element := m[key]
				if element == nil {
					fsElements = append(fsElements, fsElement{
						Key:   key,
						Path:  filepath.SplitList(key),
						IsDir: walker.Stat().IsDir(),
					})
					m[key] = &fsElements[len(fsElements)-1]
				}
				m[key].Instances = append(m[key].Instances, fileInstance{
					Size: walker.Stat().Size(),
					FS:   node.Name,
				})
			}
		}
		// reply
		w.WriteJSON(&ResponseFileTree{
			Metadata: Metadata{ID: r.ID, Action: r.Action},
			Tree:     fsElements,
		})

	case "get-content":
		for _, node := range h.Nodes {
			stat, err := node.FS.Lstat(path)
			if err != nil {
				log.Printf("Stat file %s: %s", path, err)
				continue
			}
			if stat.IsDir() {
				continue
			}
			go h.readContent(w, r, node, path)
		}

	case "search":
		_ = r.Path
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

func (h *handler) readContent(writer connWriter, r Request, src config.Src, s string) {
	rc, err := src.FS.Open(s)
	if err != nil {
		log.Printf("Open file %s: %s", s, err)
		return
	}
	defer rc.Close()

	// TODO: use specific parser by file suffix to populate logLine
	scanner := bufio.NewScanner(rc)

	var logLines []LogLine
	for scanner.Scan() {
		lineNumber := 1
		fileOffset := 0
		for scanner.Scan() {
			msg := scanner.Text()
			if err := scanner.Err(); err != nil {
				log.Println("reading standard input:", err)
			}
			logLines = append(logLines, LogLine{
				FS:         src.Name,
				FileName:   s,
				Level:      levelInfo, // TODO: read from file
				LineNumber: lineNumber,
				Msg:        msg,
				Offset:     fileOffset,
				Time:       "13:37",
			})

			lineNumber += 1
			fileOffset += len(msg)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Reading standard input:", err)
		return
	}

	writer.WriteJSON(&contentResponse{
		Metadata: Metadata{ID: r.ID, Action: r.Action},
		Lines:    logLines,
	})
}
