package config

import "github.com/kr/fs"

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

func New(sources []SrcDesc) (*Config, error) {
	c := new(Config)
	//for _, src := range sources {
	//
	//}
	return c, nil
}
