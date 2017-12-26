package filesystem

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"net"

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
	config := &ssh.ClientConfig{}
	if u.User != nil {
		config.User = u.User.Username()
		if password, ok := u.User.Password(); ok {
			config.Auth = append(config.Auth, ssh.Password(password))
		}
	}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }

	conn, err := ssh.Dial("tcp", u.Host, config)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %s", u.Hostname(), err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("create sftp client: %s", err)
	}
	basePath := u.Path

	log.Printf("Opened local: %s", basePath)
	return &SFTP{
		client:   client,
		basePath: basePath,
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
