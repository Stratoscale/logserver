package filesystem

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTP represents a filesystem over an SFTP connection
type SFTP struct {
	client   *sftp.Client
	basePath string
}

// NewSFTP returns a new SFTP filesystem
func NewSFTP(u *url.URL) (*SFTP, error) {
	config := &ssh.ClientConfig{
		User: u.User.Username(),
	}
	if password, ok := u.User.Password(); ok {
		config.Auth = append(config.Auth, ssh.Password(password))
	}

	conn, err := ssh.Dial("tcp", u.Hostname(), config)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %s", u.Hostname(), err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("create sftp client: %s", err)
	}
	return &SFTP{
		client:   client,
		basePath: u.Path,
	}, nil
}

func (s *SFTP) ReadDir(dirname string) ([]os.FileInfo, error) {
	return s.client.ReadDir(filepath.Join(s.basePath, dirname))
}

func (s *SFTP) Lstat(name string) (os.FileInfo, error) {
	return s.client.Lstat(filepath.Join(s.basePath, name))
}

func (s *SFTP) Join(elem ...string) string {
	return s.client.Join(elem...)
}

func (s *SFTP) Open(path string) (io.ReadCloser, error) {
	return s.client.Open(filepath.Join(s.basePath, path))
}

func (s *SFTP) Close() error {
	return s.client.Close()
}
