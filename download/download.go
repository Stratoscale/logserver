package download

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/source"
	"github.com/bluele/gcache"
)

var log = logrus.WithField("pkg", "router")

func New(root string, sources source.Sources, cache gcache.Cache) http.Handler {
	return &handler{
		sources: sources,
		cache:   cache,
		root:    root,
	}
}

type handler struct {
	sources source.Sources
	cache   gcache.Cache
	root    string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// collect all wanted sources
	sources := querySources(r.URL.Query())
	var downloadSources []source.Source
	for _, src := range h.sources {
		if sources[src.Name] {
			downloadSources = append(downloadSources, src)
		}
	}

	// if no specific source was specified, collect all of them
	if len(downloadSources) == 0 {
		downloadSources = h.sources
	}

	switch {
	case len(downloadSources) == 0:
		http.NotFound(w, r)
	case len(downloadSources) == 1:
		h.downloadOne(w, r, downloadSources[0])
	default:
		if filepath.Ext(r.URL.Path) == ".zip" {
			h.downloadMany(w, r, downloadSources)
		} else {
			u := filepath.Join(h.root, r.URL.Path+".zip")
			if r.URL.RawQuery != "" {
				u += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		}
	}
}

func (h *handler) downloadOne(w http.ResponseWriter, r *http.Request, src source.Source) {
	path := r.URL.Path
	log.Debugf("Download one file: %v, source: %v", path, src.Name)

	_, err := src.FS.Lstat(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	f, err := src.FS.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", contentType(path))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func (h *handler) downloadMany(w http.ResponseWriter, r *http.Request, sources []source.Source) {
	path := strings.TrimSuffix(r.URL.Path, ".zip")
	log.Debugf("Download multiple files: %v, sources: %v", path, sources)

	// create a zip file
	f, err := ioutil.TempFile("/tmp", "logserver-dl-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer f.Close()
	defer os.Remove(f.Name())

	// create a zip achiever
	z := zip.NewWriter(f)
	for _, src := range sources {

		fsFile, err := src.FS.Open(path)
		if err != nil {
			log.Debugf("Failed opening file %v/ %v: %v", src.Name, path, err)
			continue
		}

		zipFileName := fmt.Sprintf("%s-%s", src.Name, filepath.Base(path))
		zipFile, err := z.Create(zipFileName)
		if err != nil {
			log.Debugf("Failed creating zip file: %v", err)
			continue
		}
		io.Copy(zipFile, fsFile)
	}

	err = z.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// seek to beginning of file to prepare for reading
	f.Seek(0, io.SeekStart)

	w.Header().Set("Content-Type", "application/zip")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

func contentType(path string) string {
	switch filepath.Ext(path) {
	default:
		return "application/text"
	}
}

// querySources convert the url parameters to a set of source names
func querySources(u url.Values) map[string]bool {
	sources := u["fs"]
	srcMap := make(map[string]bool, len(sources))
	for _, src := range sources {
		srcMap[src] = true
	}
	return srcMap
}
