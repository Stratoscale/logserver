package config

import (
	"io"
	"net/url"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/kr/fs"
)

const (
	defaultContentBatchSize = 20
)

type Config struct {
	GlobalConfig
	Nodes []Src
}

type GlobalConfig struct {
	ContentBatchSize int `json:"content_batch_size"`
}

type Src struct {
	Name string
	FS   FileSystem
}

type FileSystem interface {
	fs.FileSystem
	Open(path string) (io.ReadCloser, error)
}

type FileConfig struct {
	Global  GlobalConfig `json:"global"`
	Sources []SrcDesc    `json:"source"`
}

type SrcDesc struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func New(fc FileConfig) (*Config, error) {
	c := new(Config)
	for _, srcDesc := range fc.Sources {
		u, err := url.Parse(srcDesc.URL)
		if err != nil {
			return c, err
		}
		var fs FileSystem
		switch u.Scheme {
		case "file":
			fs, err = filesystem.NewLocalFS(u)
		case "sftp", "ssh":
			fs, err = filesystem.NewSFTP(u)
		case "http":
			// TODO
			// fs, err = filesystem.NewHttp(u)
			return c, nil
		}
		if err != nil {
			return nil, err
		}
		c.Nodes = append(c.Nodes, Src{srcDesc.Name, fs})
	}
	c.GlobalConfig = fc.Global

	if c.GlobalConfig.ContentBatchSize == 0 {
		c.GlobalConfig.ContentBatchSize = defaultContentBatchSize
	}
	return c, nil
}
