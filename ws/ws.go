package ws

import (
	"bufio"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/parser"
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

type ContentResponse struct {
	Metadata `json:"meta"`
	Lines    []parser.LogLine `json:"lines"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ws Request from: %s", r.RemoteAddr)
	u := new(websocket.Upgrader)
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ch := make(chan interface{})
	go reader(conn, ch)

	for {
		var r Request
		err = conn.ReadJSON(&r)
		if err != nil {
			log.Printf("read: %s", err)
			return
		}
		go h.serve(ch, r)
	}
}

func reader(conn *websocket.Conn, ch <-chan interface{}) {
	for req := range ch {
		err := conn.WriteJSON(req)
		if err != nil {
			log.Printf("write: %s", err)
		}
	}
}

func (h *handler) serve(ch chan<- interface{}, r Request) {
	defer close(ch)
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
		ch <- &ResponseFileTree{
			Metadata: Metadata{ID: r.ID, Action: r.Action},
			Tree:     fsElements,
		}

	case "get-content":
		wg := sync.WaitGroup{}
		for _, node := range h.Nodes {
			stat, err := node.FS.Lstat(path)
			if err != nil {
				log.Printf("Stat file %s: %s", path, err)
				continue
			}

			if stat.IsDir() {
				continue
			}
			wg.Add(1)
			go func(node config.Src) {
				defer wg.Done()
				h.readContent(ch, r, node, path)
			}(node)
		}
		wg.Wait()

	case "search":
		_ = r.Path
		// TODO: user basepath to get file system tree
		ch <- &ContentResponse{
			Metadata: Metadata{ID: r.ID, Action: r.Action},
			Lines: []parser.LogLine{
				{Msg: "bla bla bla", Level: "debug", FS: "node0", FileName: "bla.log", LineNumber: 1},
				{Msg: "bla bla", Level: "debug", FS: "node1", FileName: "bla.log", LineNumber: 100},
				{Msg: "harta barta", Level: "debug", FS: "node1", FileName: "harta.log", LineNumber: 1},
				{Msg: "harta barta", Level: "debug", FS: "node2", FileName: "harta.log", LineNumber: 7},
				{Msg: "panic error!", Level: "debug", FS: "node2", FileName: "harta.log", LineNumber: 7},
			},
		}
	}
}

func (h *handler) readContent(ch chan<- interface{}, r Request, src config.Src, s string) {
	rc, err := src.FS.Open(s)
	if err != nil {
		log.Printf("Open file %s: %s", s, err)
		return
	}
	defer rc.Close()

	pars := parser.GetParser(filepath.Ext(s))
	scanner := bufio.NewScanner(rc)

	var logLines []parser.LogLine
	lineNumber := 1
	fileOffset := 0
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println("reading standard input:", err)
		}

		logLine, err := pars(line)
		if err != nil {
			log.Println("Failed to pars line:", err)
		}
		logLine.FileName = s
		logLine.Offset = fileOffset
		logLine.FS = src.Name
		logLine.LineNumber = lineNumber

		logLines = append(logLines, *logLine)

		lineNumber += 1
		fileOffset += len(line)

		// if we read lines more than the defined batch size, send them to the client and continue
		if len(logLines) > h.Config.ContentBatchSize {
			ch <- &ContentResponse{
				Metadata: Metadata{ID: r.ID, Action: r.Action},
				Lines:    logLines,
			}
			logLines = nil
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Reading standard input:", err)
		return
	}

	ch <- &ContentResponse{
		Metadata: Metadata{ID: r.ID, Action: r.Action},
		Lines:    logLines,
	}
}
