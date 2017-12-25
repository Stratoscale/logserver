package config

import (
	"net/url"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/filesystem/targz"
)

const (
	defaultContentBatchSize = 20
)

// Config is configuration for logserver handler
type Config struct {
	GlobalConfig
	Sources []Source
}

// GlobalConfig are global configuration parameter for logserver
type GlobalConfig struct {
	ContentBatchSize int `json:"content_batch_size"`
}

// Source is a filesystem source
type Source struct {
	Name string
	FS   filesystem.FileSystem
}

// FileConfig is logserver configuration in a file
type FileConfig struct {
	Global  GlobalConfig   `json:"global"`
	Sources []SourceConfig `json:"sources"`
}

// SourceConfig is used to configure a filesystem source
type SourceConfig struct {
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
		var fs filesystem.FileSystem
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
		fs = targz.New(fs)
		c.Sources = append(c.Sources, Source{srcDesc.Name, fs})
	}
	c.GlobalConfig = fc.Global

	if c.GlobalConfig.ContentBatchSize == 0 {
		c.GlobalConfig.ContentBatchSize = defaultContentBatchSize
	}
	return c, nil
}
