package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	loghttp "github.com/motemen/go-loghttp"
)

var nginxLine = regexp.MustCompile(`^<a\s.*>(.*)</a>\s+(\d+-\w+-\d{4}\s\d{2}:\d{2})\s+(-|[\d]+)$`)

const (
	nginxTimeFormat               = "02-Jan-2006 15:04"
	nginxLastModifiedHeaderFormat = "Mon, 02 Jan 2006 15:04:05 MST"
	nginxSchemePrefix             = "nginx+"
)

type Nginx struct {
	url *url.URL
	c   *http.Client
}

func (n *Nginx) get(path string) (*http.Response, error) {
	return n.c.Get(urlExtend(*n.url, path).String())
}

func (n *Nginx) head(path string) (*http.Response, error) {
	return n.c.Head(urlExtend(*n.url, path).String())
}

func urlExtend(u url.URL, path string) *url.URL {
	u.Path = filepath.Join(u.Path, path)
	fmt.Println("url", (&u).String())
	return &u
}

func NewNginx(u *url.URL) (FileSystem, error) {
	n := &Nginx{
		url: u,
		c: &http.Client{
			Transport: &loghttp.Transport{},
		},
	}

	// remove scheme prefix if exists
	if strings.HasPrefix(n.url.Scheme, nginxSchemePrefix) {
		n.url.Scheme = n.url.Scheme[len(nginxSchemePrefix):]
	}

	return n, nil
}

func (n *Nginx) ReadDir(dirname string) ([]os.FileInfo, error) {
	resp, err := n.get(dirname)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %d for: %s", resp.StatusCode, dirname)
	}
	switch contentType := resp.Header.Get("Content-Type"); contentType {
	case "text/html":
		return parseDirectoryHTML(resp.Body)
	case "application/json":
		return parseDirectoryJSON(resp.Body)
	default:
		return nil, fmt.Errorf("content-type %s not supported", contentType)
	}
}

func parseDirectoryHTML(body io.ReadCloser) ([]os.FileInfo, error) {
	scan := bufio.NewScanner(body)
	var files []os.FileInfo
	for scan.Scan() {
		matches := nginxLine.FindStringSubmatch(scan.Text())
		if len(matches) != 4 {
			continue
		}

		name := matches[1]
		var size int64
		if sizeStr := matches[3]; sizeStr != "-" {
			size, _ = strconv.ParseInt(sizeStr, 10, 64)
		}
		tm, err := time.Parse(nginxTimeFormat, matches[2])
		if err != nil {
			continue
		}
		files = append(files, file{
			name:  strings.TrimRight(name, "/"),
			isDir: name[len(name)-1] == '/',
			size:  size,
			time:  tm,
		})
	}
	return files, scan.Err()
}

func parseDirectoryJSON(body io.ReadCloser) ([]os.FileInfo, error) {
	var nginxFiles []nginxFile

	err := json.NewDecoder(body).Decode(&nginxFiles)
	if err != nil {
		return nil, fmt.Errorf("decoding nginx json files: %s", err)
	}

	files := make([]os.FileInfo, len(nginxFiles))
	for i, nginxFile := range nginxFiles {
		t, err := time.Parse(nginxLastModifiedHeaderFormat, nginxFile.ModTime)
		if err != nil {
			return nil, fmt.Errorf("parsing nginx date '%s': %s", nginxFile.ModTime, err)
		}
		files[i] = file{
			name:  nginxFile.Name,
			isDir: nginxFile.Type == "directory",
			size:  nginxFile.Size,
			time:  t,
		}
	}
	return files, nil
}

type nginxFile struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	ModTime string `json:"mtime"`
	Size    int64  `json:"size"`
}

func (n *Nginx) Lstat(name string) (os.FileInfo, error) {
	resp, err := n.head(name)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %d for: %s", resp.StatusCode, name)
	}
	var f file
	length := resp.Header.Get("Content-Length")
	if length == "" {
		f.isDir = true
	} else {
		f.size, err = strconv.ParseInt(length, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing content-length header: %s", length)
		}
	}
	if modTime := resp.Header.Get("Last-Modified"); modTime != "" {
		f.time, err = time.Parse(nginxLastModifiedHeaderFormat, modTime)
		if err != nil {
			return nil, fmt.Errorf("parsing last-modified header: %s", modTime)
		}
	}
	return f, nil
}

func (n *Nginx) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (n *Nginx) Open(name string) (File, error) {
	resp, err := n.get(name)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %d for: %s", resp.StatusCode, name)
	}

	return struct {
		io.ReadCloser
		io.Seeker
	}{
		ReadCloser: resp.Body,
		Seeker:     nil,
	}, nil
}

func (n *Nginx) Close() error {
	return nil
}

type file struct {
	name  string
	time  time.Time
	isDir bool
	size  int64
}

func (f file) Name() string       { return f.name }
func (f file) Size() int64        { return f.size }
func (file) Mode() os.FileMode    { return 0 }
func (f file) ModTime() time.Time { return f.time }
func (f file) IsDir() bool        { return f.isDir }
func (file) Sys() interface{}     { return nil }
