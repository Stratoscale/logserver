package config

import (
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/kr/fs"
	"net/url"
)

type Config struct {
	Nodes []Src
}

type Src struct {
	Name string
	FS   fs.FileSystem
}

type SrcDesc struct {
	Name    string
	Address string
}


func parseAddress(address string) (*url.URL, error){
    u, err := url.Parse(address)

	return u, err
}

func New(sources []SrcDesc) (*Config, error) {
	c := new(Config)
	for _, srcDesc := range sources {
		fsContext, err := parseAddress(srcDesc.Address)
		if err != nil {
			return c, err
		}
		var fs fs.FileSystem
		switch fsContext.Scheme {
		case "file":
			fs = &filesystem.LocalFS{BaseFS: filesystem.BaseFS{fsContext}}
		case "ssh":
			fs = &filesystem.SshFS{BaseFS: filesystem.BaseFS{fsContext}}
		case "http":
			fs = &filesystem.HttpFS{BaseFS: filesystem.BaseFS{fsContext}}
		}
		c.Nodes = append(c.Nodes, Src{srcDesc.Name, fs})
	}
	return c, nil
}
