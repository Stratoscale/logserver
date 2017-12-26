package ws

import (
	"bufio"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"regexp"

	"fmt"

<<<<<<< HEAD
=======
	"os"
>>>>>>> xxx
	"strings"

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
	ID     int     `json:"id"`
	Action string  `json:"action"`
	FS     string  `json:"fs,omitempty"`
	Path   Path `json:"path,omitempty"`
}

type Request struct {
	Metadata `json:"meta"`
	Path     Path   `json:"path"`
	Regexp   string `json:"regexp"`
}

type Path []string

type Response struct {
	Metadata `json:"meta"`
	Lines    []parser.LogLine `json:"lines,omitempty"`
	Error    string           `json:"error,omitempty"`
	Tree     []FSElement      `json:"tree,omitempty"`
}

type FSElement struct {
	Key       string         `json:"key"`
	Path      Path           `json:"path"`
	IsDir     bool           `json:"is_dir"`
	Instances []FileInstance `json:"instances"`
}

type FileInstance struct {
	Size int64  `json:"size"`
	FS   string `json:"fs"`
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
	defer close(ch)
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
	switch r.Action {
	case "get-file-tree":
		var (
			fsElements []FSElement
			m          = make(map[string]*FSElement)
		)

		for _, node := range h.Sources {
			path := node.FS.Join(r.Path...)
			if path == "" {
				path = "/"
			}
			walker := fs.WalkFS(path, node.FS)
			for walker.Step() {
				if err := walker.Err(); err != nil {
					log.Printf("Walk: %s", err)
					continue
				}

				key := strings.Trim(walker.Path(), string(os.PathSeparator))
				if key == "" {
					continue
				}

				element := m[key]
				if element == nil {
					fsElements = append(fsElements, FSElement{
						Key:   key,
						Path:  strings.Split(key, string(os.PathSeparator)),
						IsDir: walker.Stat().IsDir(),
					})
					m[key] = &fsElements[len(fsElements)-1]
				}
				m[key].Instances = append(m[key].Instances, FileInstance{
					Size: walker.Stat().Size(),
					FS:   node.Name,
				})
			}
		}
		// reply
		ch <- &Response{
			Metadata: r.Metadata,
			Tree:     fsElements,
		}

	case "get-content":
		wg := sync.WaitGroup{}
		wg.Add(len(h.Sources))
		for _, node := range h.Sources {
			go func(node config.Source) {
				path := node.FS.Join(r.Path...)
				if path == "" {
					path = "/"
				}
				h.read(ch, r, node, path, nil)
				wg.Done()
			}(node)
		}
		wg.Wait()

	case "search":
		re, err := regexp.Compile(r.Regexp)
		if err != nil {
			ch <- &Response{
				Metadata: r.Metadata,
				Error:    fmt.Sprintf("Bad regexp %s: %s", r.Regexp, err),
			}
		}
		wg := sync.WaitGroup{}
		wg.Add(len(h.Sources))
		for _, node := range h.Sources {
			go func(node config.Source) {
				path := node.FS.Join(r.Path...)
				if path == "" {
					path = "/"
				}
				h.search(ch, r, node, path, re)
				wg.Done()
			}(node)
		}
		wg.Wait()
	}
}

func (h *handler) search(ch chan<- interface{}, req Request, node config.Source, path string, re *regexp.Regexp) {
	var walker = fs.WalkFS(path, node.FS)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			log.Printf("Walk: %s", err)
			continue
		}
		filePath := walker.Path()
		h.read(ch, req, node, filePath, re)
	}
}

func (h *handler) read(ch chan<- interface{}, req Request, node config.Source, path string, re *regexp.Regexp) {
	stat, err := node.FS.Lstat(path)
	if err != nil {
		return
	}
	if stat.IsDir() {
		return
	}

	r, err := node.FS.Open(path)
	if err != nil {
		log.Printf("Open %s: %s", path, err)
		return
	}
	defer r.Close()

	var (
		pars       = parser.GetParser(filepath.Ext(path))
		scanner    = bufio.NewScanner(r)
		logLines   []parser.LogLine
		lineNumber = 1
		fileOffset = 0
		respMeta   = Metadata{ID: req.Metadata.ID, Action: req.Metadata.Action, FS: node.Name, Path: strings.Split(path, "/")}
		sentAny    = false
	)

	if respMeta.Path[0] == "" {
		respMeta.Path = respMeta.Path[1:]
	}

	for scanner.Scan() {
		if re != nil && !re.Match(scanner.Bytes()) {
			lineNumber += 1
			fileOffset += len(scanner.Bytes())
			continue
		}

		logLine, err := pars(scanner.Bytes())
		if err != nil {
			log.Println("Failed to pars line:", err)
		}
		logLine.FileName = path
		logLine.Offset = fileOffset
		logLine.FS = node.Name
		logLine.LineNumber = lineNumber

		logLines = append(logLines, *logLine)
		lineNumber += 1
		fileOffset += len(scanner.Bytes())

		// if we read lines more than the defined batch size, send them to the client and continue
		if len(logLines) > h.Config.ContentBatchSize {
			sentAny = true
			ch <- &Response{Metadata: respMeta, Lines: logLines}
			logLines = nil
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Scan:", err)
		return
	}
	if (len(logLines) != 0 || !sentAny) && re != nil {
		ch <- &Response{Metadata: respMeta, Lines: logLines}
	}
}
