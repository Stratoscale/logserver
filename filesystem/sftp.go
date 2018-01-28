package filesystem

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTP represents a filesystem over an SFTP connection
type SFTP struct {
	client   *sftp.Client
	basePath string
}

// NewSFTP returns a new SFTP filesystem
func NewSFTP(u *url.URL) (FileSystem, error) {
	config := &ssh.ClientConfig{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
	}

	// add user's public key to ssh configuration
	usr, pubKey := publicKey()
	config.User = usr
	if pubKey != nil {
		config.Auth = append(config.Auth, pubKey)
	}

	// add user/password authentication if specified in url
	if u.User != nil {
		usr = u.User.Username()
		config.User = usr
		if password, ok := u.User.Password(); ok {
			config.Auth = append(config.Auth, ssh.Password(password))
		}
	}

	hp := hostPort(u.Host)
	conn, err := ssh.Dial("tcp", hp, config)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("create sftp client: %s", err)
	}
	basePath := u.Path
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

func (s *SFTP) Open(path string) (File, error) {
	return s.client.Open(filepath.Join(s.basePath, path))
}

func (s *SFTP) Close() error {
	return s.client.Close()
}

func publicKey() (username string, pubKey ssh.AuthMethod) {
	usr, err := user.Current()
	if err != nil {
		return "", nil
	}
	key, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".ssh/id_rsa"))
	if err != nil {
		return "", nil
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", nil
	}
	return usr.Username, ssh.PublicKeys(signer)
}

func hostPort(host string) string {
	if !strings.ContainsRune(host, ':') {
		return fmt.Sprintf("%s:%d", host, 22)
	}
	return host
}
