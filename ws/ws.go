package ws

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/parser"
	"github.com/gorilla/websocket"
	"github.com/kr/fs"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger().WithField("pkg", "ws")

// New returns a new websocket handler
func New(c config.Config) http.Handler {
	return &handler{
		Config: c,
	}
}

type handler struct {
	config.Config
}

// Path describes a file path
// Each directory or file is an item in the slice
type Path []string

// Meta is request/response metadata
type Meta struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
	FS     string `json:"fs,omitempty"`
	Path   Path   `json:"path,omitempty"`
}

// Request from client
type Request struct {
	Meta   `json:"meta"`
	Path   Path   `json:"path"`
	Regexp string `json:"regexp"`
}

// Response from the server
type Response struct {
	Meta  `json:"meta"`
	Lines []parser.LogLine `json:"lines,omitempty"`
	Tree  []*File          `json:"tree,omitempty"`
	Error string           `json:"error,omitempty"`
}

// File describes a file in multiple file systems
type File struct {
	Key   string `json:"key"`
	Path  Path   `json:"path"`
	IsDir bool   `json:"is_dir"`
	// Instances are all the instances of the same file in different file systems
	Instances []FileInstance `json:"instances"`
}

// FileInstance describe a file on a filesystem
type FileInstance struct {
	Size int64  `json:"size"`
	FS   string `json:"fs"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Got ws Request from: %s", r.RemoteAddr)
	u := new(websocket.Upgrader)
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Errorf("Failed upgrade from %s", r.RemoteAddr)
		return
	}

	ch := make(chan *Response)
	defer close(ch)
	go reader(conn, ch)

	for {
		var r Request
		err = conn.ReadJSON(&r)
		if err != nil {
			log.WithError(err).Errorf("Failed read")
			return
		}
		go h.serve(ch, r)
	}
}

func reader(conn *websocket.Conn, ch <-chan *Response) {
	for req := range ch {
		err := conn.WriteJSON(req)
		if err != nil {
			log.Printf("write: %s", err)
		}
	}
}

func (h *handler) serve(ch chan<- *Response, r Request) {
	switch r.Action {
	case "get-file-tree":
		h.serveTree(r, ch)

	case "get-content":
		h.serveContent(r, ch)

	case "search":
		h.search(r, ch)
	}
}

func (h *handler) serveTree(r Request, ch chan<- *Response) {
	var (
		fsElements []*File
		m          = make(map[string]*File)
	)
	for _, node := range h.Sources {
		path := node.FS.Join(r.Path...)
		walker := fs.WalkFS(path, node.FS)
		for walker.Step() {
			if err := walker.Err(); err != nil {
				log.WithError(err).Errorf("Failed walk %s:%s", node.Name, path)
				continue
			}

			key := strings.Trim(walker.Path(), string(os.PathSeparator))
			if key == "" {
				continue
			}

			element := m[key]
			if element == nil {
				fsElements = append(fsElements, &File{
					Key:   key,
					Path:  strings.Split(key, string(os.PathSeparator)),
					IsDir: walker.Stat().IsDir(),
				})
				m[key] = fsElements[len(fsElements)-1]
			}
			m[key].Instances = append(m[key].Instances, FileInstance{
				Size: walker.Stat().Size(),
				FS:   node.Name,
			})
		}
	}
	// reply
	ch <- &Response{Meta: r.Meta, Tree: fsElements}
}

func (h *handler) serveContent(r Request, ch chan<- *Response) {
	wg := sync.WaitGroup{}
	wg.Add(len(h.Sources))
	for _, node := range h.Sources {
		go func(node config.Source) {
			path := node.FS.Join(r.Path...)
			h.read(ch, r, node, path, nil)
			wg.Done()
		}(node)
	}
	wg.Wait()
}

func (h *handler) search(r Request, ch chan<- *Response) {
	re, err := regexp.Compile(r.Regexp)
	if err != nil {
		ch <- &Response{
			Meta:  r.Meta,
			Error: fmt.Sprintf("Bad regexp %s: %s", r.Regexp, err),
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(len(h.Sources))
	for _, node := range h.Sources {
		go func(node config.Source) {
			path := node.FS.Join(r.Path...)
			h.searchNode(ch, r, node, path, re)
			wg.Done()
		}(node)
	}
	wg.Wait()
}

func (h *handler) searchNode(ch chan<- *Response, req Request, node config.Source, path string, re *regexp.Regexp) {
	var walker = fs.WalkFS(path, node.FS)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			log.WithError(err).Errorf("Failed walk %s:%s", node.Name, path)
			continue
		}
		filePath := walker.Path()
		h.read(ch, req, node, filePath, re)
	}
}

func (h *handler) read(ch chan<- *Response, req Request, node config.Source, path string, re *regexp.Regexp) {
	log := log.WithField("path", fmt.Sprintf("%s:%s", node.Name, path))
	stat, err := node.FS.Lstat(path)
	if err != nil {
		// the file might not exists in all filesystem, so just return without an error
		return
	}
	if stat.IsDir() {
		return
	}

	r, err := node.FS.Open(path)
	if err != nil {
		log.WithError(err).Error("Failed open")
		return
	}
	defer r.Close()

	var (
		pars       = parser.GetParser(filepath.Ext(path))
		scanner    = bufio.NewScanner(r)
		logLines   []parser.LogLine
		lineNumber = 1
		fileOffset = 0
		respMeta   = Meta{ID: req.Meta.ID, Action: req.Meta.Action, FS: node.Name, Path: strings.Split(path, "/")}
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
			logLine = &parser.LogLine{Msg: scanner.Text()}
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
			ch <- &Response{Meta: respMeta, Lines: logLines}
			logLines = nil
		}
	}
	if err := scanner.Err(); err != nil {
		log.WithError(err).Errorf("Failed scan")
		return
	}
	if len(logLines) == 0 && !sentAny {
		return
	}
	ch <- &Response{Meta: respMeta, Lines: logLines}
}
