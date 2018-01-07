package ws

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/parser"
	"github.com/gorilla/websocket"
	"github.com/kr/fs"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("pkg", "ws")

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
	Path   Path     `json:"path"`
	Regexp string   `json:"regexp"`
	Nodes  []string `json:"nodes"`
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

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	for {
		var req Request
		err = conn.ReadJSON(&req)
		if err != nil {
			log.WithError(err).Errorf("Failed read")
			return
		}
		// cancel the last serving up on a new request
		if cancel != nil {
			cancel()
		}
		ctx, cancel = context.WithCancel(r.Context())
		go h.serve(ctx, ch, req)
	}
	// cancel last serving if exists
	if cancel != nil {
		cancel()
	}
}

func reader(conn *websocket.Conn, ch <-chan *Response) {
	for req := range ch {
		err := conn.WriteJSON(req)
		if err != nil {
			log.WithError(err).Errorf("Failed write")
		}
	}
}

func (h *handler) serve(ctx context.Context, ch chan<- *Response, req Request) {
	switch req.Action {
	case "get-file-tree":
		h.serveTree(ctx, req, ch)

	case "get-content":
		h.serveContent(ctx, req, ch)

	case "search":
		h.search(ctx, req, ch)
	}
}

func (h *handler) serveTree(ctx context.Context, req Request, ch chan<- *Response) {
	var (
		files   []*File
		fileMap = make(map[string]*File)
	)
	for _, node := range filterNodes(h.Sources, req.Nodes) {
		path := node.FS.Join(req.Path...)
		walker := fs.WalkFS(path, node.FS)
		for walker.Step() {
			if err := ctx.Err(); err != nil {
				return
			}

			if err := walker.Err(); err != nil {
				log.WithError(err).Errorf("Failed walk %s:%s", node.Name, path)
				continue
			}

			key := strings.Trim(walker.Path(), string(os.PathSeparator))
			if key == "" {
				continue
			}

			element := fileMap[key]
			if element == nil {
				files = append(files, &File{
					Key:   key,
					Path:  strings.Split(key, string(os.PathSeparator)),
					IsDir: walker.Stat().IsDir(),
				})
				fileMap[key] = files[len(files)-1]
			}
			fileMap[key].Instances = append(fileMap[key].Instances, FileInstance{
				Size: walker.Stat().Size(),
				FS:   node.Name,
			})
		}
	}
	// reply
	ch <- &Response{Meta: req.Meta, Tree: files}
}

func (h *handler) serveContent(ctx context.Context, req Request, ch chan<- *Response) {
	wg := sync.WaitGroup{}
	nodes := filterNodes(h.Sources, req.Nodes)
	wg.Add(len(nodes))
	for _, node := range nodes {
		go func(node config.Source) {
			defer wg.Done()
			path := node.FS.Join(req.Path...)
			h.read(ctx, ch, req, node, path, nil)
		}(node)
	}
	wg.Wait()
}

func (h *handler) search(ctx context.Context, req Request, ch chan<- *Response) {
	re, err := regexp.Compile(req.Regexp)
	if err != nil {
		ch <- &Response{
			Meta:  req.Meta,
			Error: fmt.Sprintf("Bad regexp %s: %s", req.Regexp, err),
		}
		return
	}
	nodes := filterNodes(h.Sources, req.Nodes)
	wg := sync.WaitGroup{}
	wg.Add(len(nodes))
	for _, node := range nodes {
		go func(node config.Source) {
			defer wg.Done()
			path := node.FS.Join(req.Path...)
			h.searchNode(ctx, ch, req, node, path, re)
		}(node)
	}
	wg.Wait()
}

func (h *handler) searchNode(ctx context.Context, ch chan<- *Response, req Request, node config.Source, path string, re *regexp.Regexp) {
	var walker = fs.WalkFS(path, node.FS)
	for walker.Step() {
		if err := ctx.Err(); err != nil {
			return
		}

		if err := walker.Err(); err != nil {
			log.WithError(err).Errorf("Failed walk %s:%s", node.Name, path)
			continue
		}
		filePath := walker.Path()
		h.read(ctx, ch, req, node, filePath, re)
	}
}

func (h *handler) read(ctx context.Context, ch chan<- *Response, req Request, node config.Source, path string, re *regexp.Regexp) {
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
		pars         = parser.GetParser(filepath.Ext(path))
		scanner      = bufio.NewScanner(r)
		logLines     []parser.LogLine
		lastRespTime = time.Now()
		lineNumber   = 1
		fileOffset   = 0
		respMeta     = Meta{
			ID:     req.Meta.ID,
			Action: req.Meta.Action,
			FS:     node.Name,
			Path:   strings.Split(path, "/"),
		}
		sentAny = false
	)

	if respMeta.Path[0] == "" {
		respMeta.Path = respMeta.Path[1:]
	}

	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return
		}
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

		// if we read lines more than the defined batch size or batch time,
		// send them to the client and continue
		if len(logLines) > h.Config.ContentBatchSize || time.Now().Sub(lastRespTime) > h.ContentBatchTime {
			sentAny = true
			ch <- &Response{Meta: respMeta, Lines: logLines}
			logLines = nil
			lastRespTime = time.Now()
		}
		// max search lines exceeded
		if re != nil && lineNumber > h.SearchMaxSize {
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.WithError(err).Errorf("Failed scan")
		return
	}
	if len(logLines) == 0 && (sentAny || re != nil) {
		return
	}
	ch <- &Response{Meta: respMeta, Lines: logLines}
}

func filterNodes(sources []config.Source, filterNodes []string) []config.Source {
	if len(filterNodes) == 0 {
		return sources
	}
	nodes := make(map[string]bool, len(filterNodes))
	for _, node := range filterNodes {
		nodes[node] = true
	}
	ret := make([]config.Source, 0, len(filterNodes))
	for _, src := range sources {
		if nodes[src.Name] {
			ret = append(ret, src)
		}
	}
	return ret
}
