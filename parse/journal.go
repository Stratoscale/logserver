package parse

import (
	"os"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

const sep = string(os.PathSeparator)

var reJournal = regexp.MustCompile(`(?P<time>\d{4}-\d{2}-\d{2}\W\d{2}:\d{2}:\d{2}.\d{6}\W[+-]\d{4}\W[[:alpha:]]+)\W+(?P<msg>.*)`)

// AppendJournal appends a journalctl parser to a parser list
func (p *Parse) AppendJournal(name string) error {
	g, err := glob.Compile(strings.Trim(name, sep))
	if err != nil {
		return err
	}
	*p = append(*p, parser{
		glob:   g,
		regexp: reJournal,
		Config: Config{
			TimeFormats: []string{"2006-01-02 15:04:05.000000 -0700 MST"},
		},
	})
	return nil
}
