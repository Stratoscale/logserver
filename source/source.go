package source

import (
	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/filesystem/targz"
	"github.com/bluele/gcache"
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

func New(c []Config, cache gcache.Cache) (Sources, error) {
	var s Sources
	for _, srcDesc := range c {
		u, err := url.Parse(srcDesc.URL)
		if err != nil {
			return s, err
		}
		var fs filesystem.FileSystem
		switch u.Scheme {
		case "file":
			fs, err = filesystem.NewLocal(u)
		case "sftp", "ssh":
			fs, err = filesystem.NewSFTP(u)
		case "nginx+http", "nginx+https":
			if srcDesc.OpenTarFiles {
				return nil, fmt.Errorf("can't have 'open_tar_files' option over http")
			}
			fs, err = filesystem.NewNginx(u)
		}
		if err != nil {
			return nil, err
		}
		log.Infof("Opened: %s", u)
		if srcDesc.OpenTarFiles {
			fs = targz.New(fs, cache, srcDesc.URL)
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
