package config

import (
	"io"

	"net/url"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/kr/fs"
)

type Config struct {
	Nodes []Src
}

type Src struct {
	Name string
	FS   FileSystem
}

type FileSystem interface {
	fs.FileSystem
	Open(path string) (io.ReadCloser, error)
}

type SrcDesc struct {
	Name    string
	Address string
}

func New(sources []SrcDesc) (*Config, error) {
	c := new(Config)
	for _, srcDesc := range sources {
		fsContext, err := url.Parse(srcDesc.Address)
		if err != nil {
			return c, err
		}
		var fs FileSystem
		switch fsContext.Scheme {
		case "file":
			fs = &filesystem.LocalFS{BaseFS: filesystem.BaseFS{fsContext}}
		case "ssh":
			// TODO
			return c, nil
			// fs, err = filesystem.NewSftp(fsContext)
		case "http":
			// TODO
			// fs, err = filesystem.NewHttp(fsContext)
			return c, nil
		}
		c.Nodes = append(c.Nodes, Src{srcDesc.Name, fs})
	}
	return c, nil
}
