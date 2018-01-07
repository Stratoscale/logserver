package config

import (
	"net/url"
	"time"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/filesystem/targz"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("pkg", "config")

const (
	defaultContentBatchSize = 2000
	defaultContentBatchTime = time.Second * 2
	defaultSearchMaxSize    = 5000
)

// Config is configuration for logserver handler
type Config struct {
	GlobalConfig
	Sources  []Source
	BasePath string
}

// GlobalConfig are global configuration parameter for logserver
type GlobalConfig struct {
	ContentBatchSize int           `json:"content_batch_size"`
	ContentBatchTime time.Duration `json:"content_batch_time"`
	SearchMaxSize    int           `json:"search_max_size"`
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
	Name         string `json:"name"`
	URL          string `json:"url"`
	OpenTarFiles bool   `json:"open_tar_files"`
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
		log.Infof("Opened: %s", u)
		if srcDesc.OpenTarFiles {
			fs = targz.New(fs)
		}
		c.Sources = append(c.Sources, Source{srcDesc.Name, fs})
	}
	c.GlobalConfig = fc.Global

	if c.GlobalConfig.ContentBatchSize == 0 {
		c.GlobalConfig.ContentBatchSize = defaultContentBatchSize
	}
	if c.GlobalConfig.ContentBatchTime == 0 {
		c.GlobalConfig.ContentBatchTime = defaultContentBatchTime
	}
	if c.GlobalConfig.SearchMaxSize == 0 {
		c.GlobalConfig.SearchMaxSize = defaultSearchMaxSize
	}
	return c, nil
}

func (c *Config) CloseSources() {
	for _, src := range c.Sources {
		err := src.FS.Close()
		if err != nil {
			log.WithError(err).Errorf("Failed closing source %s", src.Name)
		}
	}
}
