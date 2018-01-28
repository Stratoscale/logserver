package filesystem

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/kr/fs"
)

var log = logrus.StandardLogger().WithField("pkg", "filesystem")

type File interface {
	io.Reader
	io.Closer
	io.Seeker
}

// Filesystem represents a filesystem, which can be local or remote
type FileSystem interface {
	fs.FileSystem
	// Open opens a file in the filesystem
	Open(path string) (File, error)
	// Close closes the filesystem.
	// This is useful for remote filesystems, like http, or sftp
	Close() error
}
