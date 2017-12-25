package filesystem

import (
	"io"
	"github.com/kr/fs"
)

type FileSystem interface {
	fs.FileSystem
	Open(path string) (io.ReadCloser, error)
}

