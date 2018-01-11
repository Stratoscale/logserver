package source

import (
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/filesystem/targz"
)

var log = logrus.WithField("pkg", "config")

// Config is used to configure a filesystem source
type Config struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	OpenTarFiles bool   `json:"open_tar_files"`
}

type Sources []Source

// Source is a filesystem source
type Source struct {
	Name string
	FS   filesystem.FileSystem
}

func New(c []Config) (Sources, error) {
	var s Sources
	for _, srcDesc := range c {
		u, err := url.Parse(srcDesc.URL)
		if err != nil {
			return s, err
		}
		var fs filesystem.FileSystem
		switch u.Scheme {
		case "file":
			fs, err = filesystem.NewLocalFS(u)
		case "sftp", "ssh":
			fs, err = filesystem.NewSFTP(u)
		}
		if err != nil {
			return nil, err
		}
		log.Infof("Opened: %s", u)
		if srcDesc.OpenTarFiles {
			fs = targz.New(fs)
		}
		s = append(s, Source{srcDesc.Name, fs})
	}
	return s, nil
}

func (s Sources) CloseSources() {
	for _, src := range s {
		err := src.FS.Close()
		if err != nil {
			log.WithError(err).Errorf("Failed closing source %s", src.Name)
		}
	}
}
