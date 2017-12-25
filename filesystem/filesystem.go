package filesystem

import (
	"io"
	"github.com/kr/fs"
)

// Filesystem represents a filesystem, which can be local or remote
type FileSystem interface {
	fs.FileSystem
	// Open opens a file in the filesystem
	Open(path string) (io.ReadCloser, error)
	// Close closes the filesystem.
	// This is useful for remote filesystems, like http, or sftp
	Close() error
}

