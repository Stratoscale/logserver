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
		u, err := url.Parse(srcDesc.Address)
		if err != nil {
			return c, err
		}
		var fs FileSystem
		switch u.Scheme {
		case "file":
			fs, err = filesystem.NewLocalFS(u)
		case "ssh":
			// TODO
			return c, nil
			// fs, err = filesystem.NewSftp(u)
		case "http":
			// TODO
			// fs, err = filesystem.NewHttp(u)
			return c, nil
		}
		c.Nodes = append(c.Nodes, Src{srcDesc.Name, fs})
	}
	return c, nil
}
