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

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/debug"
	"github.com/Stratoscale/logserver/parser"
	"github.com/Stratoscale/logserver/source"
	"github.com/gorilla/websocket"
	"github.com/kr/fs"
)

var log = logrus.WithField("pkg", "ws")

// New returns a new websocket handler
func New(c source.Config) http.Handler {
	return &handler{
		Config: c,
	}
}

type handler struct {
	source.Config
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
	Meta       `json:"meta"`
	Path       Path      `json:"path"`
	Regexp     string    `json:"regexp"`
	FilterFS   []string  `json:"filter_fs"`
	FilterTime TimeRange `json:"filter_time"`
}

type TimeRange struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
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
	log.Infof("New WS Client from: %s", r.RemoteAddr)
	defer log.Info("Disconnected WS Client from: %s", r.RemoteAddr)
	u := new(websocket.Upgrader)
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Errorf("Failed upgrade from %s", r.RemoteAddr)
		return
	}

	send := make(chan *Response)
	defer close(send)
	go reader(conn, send)

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
		go h.serve(ctx, req, func(resp *Response) {
			select {
			case send <- resp:
			case <-ctx.Done():
			}
		})
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

func (h *handler) serve(ctx context.Context, req Request, send func(*Response)) {
	defer debug.Time(log, "Request %+v", req.Meta)()
	switch req.Action {
	case "get-file-tree":
		h.serveTree(ctx, req, send)

	case "get-content":
		h.serveContent(ctx, req, send)

	case "search":
		h.search(ctx, req, send)
	}
	if err := ctx.Err(); err != nil {
		log.Debugf("Request %d cancelled", req.ID)
	}
}

func (h *handler) serveTree(ctx context.Context, req Request, send func(*Response)) {
	var (
		c  = newCombiner()
		wg sync.WaitGroup
	)
	sources := filterNodes(h.Sources, req.FilterFS)
	wg.Add(len(sources))
	for _, src := range sources {
		go func(src source.Source) {
			defer wg.Done()
			h.srcTree(ctx, req, src, c)
		}(src)
	}
	wg.Wait()
	log.Debugf("Serve tree for %v with %d files", req.Path, len(c.files))
	send(&Response{Meta: req.Meta, Tree: c.files})
}

func (h *handler) srcTree(ctx context.Context, req Request, src source.Source, c *combiner) {
	const sep = string(os.PathSeparator)
	path := src.FS.Join(req.Path...)
	walker := fs.WalkFS(path, src.FS)
	for walker.Step() {
		if err := ctx.Err(); err != nil {
			return
		}

		if err := walker.Err(); err != nil {
			log.WithError(err).Errorf("Failed walk %s:%s", src.Name, path)
			continue
		}

		key := strings.Trim(walker.Path(), sep)
		if key == "" {
			continue
		}

		c.add(
			File{
				Key:   key,
				Path:  strings.Split(key, sep),
				IsDir: walker.Stat().IsDir(),
			},
			FileInstance{
				Size: walker.Stat().Size(),
				FS:   src.Name,
			},
		)
	}
}

type combiner struct {
	files []*File
	index map[string]*File
	lock  sync.Mutex
}

func newCombiner() *combiner {
	return &combiner{index: make(map[string]*File)}
}

func (c *combiner) add(f File, instance FileInstance) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.index[f.Key] == nil {
		c.files = append(c.files, &f)
		c.index[f.Key] = c.files[len(c.files)-1]
	}
	c.index[f.Key].Instances = append(c.index[f.Key].Instances, instance)
}

func (h *handler) serveContent(ctx context.Context, req Request, send func(*Response)) {
	wg := sync.WaitGroup{}
	sources := filterNodes(h.Sources, req.FilterFS)
	wg.Add(len(sources))
	for _, src := range sources {
		go func(src source.Source) {
			defer wg.Done()
			path := src.FS.Join(req.Path...)
			h.read(ctx, send, req, src, path, nil)
		}(src)
	}
	wg.Wait()
}

func (h *handler) search(ctx context.Context, req Request, send func(*Response)) {
	re, err := regexp.Compile(req.Regexp)
	if err != nil {
		send(&Response{
			Meta:  req.Meta,
			Error: fmt.Sprintf("Bad regexp %s: %s", req.Regexp, err),
		})
		return
	}
	nodes := filterNodes(h.Sources, req.FilterFS)
	wg := sync.WaitGroup{}
	wg.Add(len(nodes))
	for _, node := range nodes {
		go func(node source.Source) {
			defer wg.Done()
			path := node.FS.Join(req.Path...)
			h.searchNode(ctx, send, req, node, path, re)
		}(node)
	}
	wg.Wait()
}

func (h *handler) searchNode(ctx context.Context, send func(*Response), req Request, node source.Source, path string, re *regexp.Regexp) {
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
		h.read(ctx, send, req, node, filePath, re)
	}
}

func (h *handler) read(ctx context.Context, send func(*Response), req Request, node source.Source, path string, re *regexp.Regexp) {
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

		if filterOutTime(logLine, req.FilterTime) {
			continue
		}

		logLines = append(logLines, *logLine)
		lineNumber += 1
		fileOffset += len(scanner.Bytes())

		// if we read lines more than the defined batch size or batch time,
		// send them to the client and continue
		if len(logLines) > h.Config.ContentBatchSize || time.Now().Sub(lastRespTime) > h.ContentBatchTime {
			sentAny = true
			send(&Response{Meta: respMeta, Lines: logLines})
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
	send(&Response{Meta: respMeta, Lines: logLines})
}

func filterNodes(sources []source.Source, filterSources []string) []source.Source {
	if len(filterSources) == 0 {
		return sources
	}
	nodes := make(map[string]bool, len(filterSources))
	for _, node := range filterSources {
		nodes[node] = true
	}
	ret := make([]source.Source, 0, len(filterSources))
	for _, src := range sources {
		if nodes[src.Name] {
			ret = append(ret, src)
		}
	}
	return ret
}

func filterOutTime(line *parser.LogLine, timeRange TimeRange) bool {
	if start := timeRange.Start; start != nil {
		return line.Time == nil || start.After(*line.Time)
	}
	if end := timeRange.End; end != nil {
		return line.Time == nil || end.Before(*line.Time)
	}
	return false
}
