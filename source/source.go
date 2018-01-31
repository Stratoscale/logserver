package source

import (
	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/filesystem"
	"github.com/Stratoscale/logserver/filesystem/tar"
	"github.com/bluele/gcache"
)

var log = logrus.WithField("pkg", "config")

// Config is used to configure a filesystem source
type Config struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Flags
}

// Flags are configuration options for a source
type Flags struct {
	OpenTar     bool   `json:"open_tar"`
	OpenJournal string `json:"open_journal"`
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
			if srcDesc.OpenTar {
				return nil, fmt.Errorf("can't have 'open_tar' option over http")
			}
			fs, err = filesystem.NewNginx(u)
		}
		if err != nil {
			log.WithError(err).Errorf("Failed adding source %s(%s)", srcDesc.Name, srcDesc.URL)
			continue
		}
		log.Infof("Opened %s: %s", srcDesc.Name, srcDesc.URL)
		if srcDesc.OpenTar {
			fs = tar.Wrap(fs, cache, srcDesc.URL+"/")
		}
		if srcDesc.OpenJournal != "" {
			fs = filesystem.WrapJournal(fs, srcDesc.OpenJournal)
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
