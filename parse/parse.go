package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/gobwas/glob"
)

type Type string

const (
	KeyTime  = "time"
	KeyLevel = "level"
	KeyMsg   = "msg"
	KeyArgs  = "args"
	KeyThreadName  = "threadName"
)

// noParserAfter determines how many line should be parsed before choosing a no-parser
// for a specific file
const noParserAfter = 1000

type Config struct {
	Glob        string            `json:"glob"`
	JsonMapping map[string]string `json:"json_mapping"`
	Regexp      string            `json:"regexp"`
	TimeFormats []string          `json:"time_formats"`
	AppendArgs  bool              `json:"append_args"`
}

type Parse []parser

func New(configs []Config) (Parse, error) {
	var ps Parse
	for _, c := range configs {
		if c.Regexp != "" && len(c.JsonMapping) != 0 {
			return nil, fmt.Errorf("can't specify both 'regexp' and 'json_mapping', got: %+v", c)
		}
		if c.Regexp == "" && len(c.JsonMapping) == 0 {
			return nil, fmt.Errorf("must specify 'regexp' or 'json_mapping', got: %+v", c)
		}

		var (
			p   = parser{Config: c}
			err error
		)

		if c.Regexp != "" {
			p.regexp, err = regexp.Compile(c.Regexp)
			if err != nil {
				return nil, fmt.Errorf("compiling regexp: %s", err)
			}
		}
		if c.Glob == "" {
			c.Glob = "*"
		}
		p.glob, err = glob.Compile(c.Glob)
		if err != nil {
			return nil, fmt.Errorf("compiling glob: %s", err)
		}
		ps = append(ps, p)
	}
	return ps, nil
}

type parser struct {
	Config
	regexp *regexp.Regexp
	glob   glob.Glob
}

// Memory is used to remember which parser applied for a file
// the same Memory object should be passed to Parse when parsing lines from the same file
type Memory struct {
	parser *parser
	count  int
}

func (ps Parse) Parse(logName string, line []byte, mem *Memory) *Log {

	// check for memory for file that was already parsed with a parser
	if mem.parser != nil {
		parsed := mem.parser.parse(line)
		if parsed != nil {
			return parsed
		} else {
			return &Log{Msg: string(line)}
		}
	}

	for _, p := range ps {
		if !p.glob.Match(logName) {
			continue
		}
		log := p.parse(line)
		if log != nil {
			mem.parser = &p
			return log
		}
	}
	mem.count++
	if mem.count >= noParserAfter {
		mem.parser = &parser{}
	}
	return &Log{Msg: string(line)}
}

func (p *parser) parse(line []byte) *Log {
	switch {
	case len(p.JsonMapping) > 0:
		return p.parseJson(line)
	case p.regexp != nil:
		return p.parseRegexp(line)
	default:
		// default no-parser
		return &Log{Msg: string(line)}
	}
}

func (p *parser) parseJson(line []byte) *Log {
	var j map[string]interface{}
	err := json.Unmarshal(line, &j)
	if err != nil {
		return nil
	}
	log := new(Log)
	var ok bool

	msgKey := p.JsonMapping[KeyMsg]
	if log.Msg, ok = j[msgKey].(string); !ok {
		return nil
	}
	delete(j, msgKey)

	if jsonKey, ok := p.JsonMapping[KeyLevel]; ok {
		if log.Level, ok = j[jsonKey].(string); !ok {
			return nil
		}
		delete(j, jsonKey)
	}
	if jsonKey, ok := p.JsonMapping[KeyThreadName]; ok {
		if log.threadName, ok = j[jsonKey].(string); !ok {
			return nil
		}
		fmt.Printf("thread: %s\n", log.threadName)
		delete(j, jsonKey)
	}

	if jsonKey, ok := p.JsonMapping[KeyTime]; ok {
		switch t := j[jsonKey].(type) {
		case float64:
			tt := time.Unix(int64(t), int64(t-float64(int64(t))))
			log.Time = &tt
		case int64:
			tt := time.Unix(t, 0)
			log.Time = &tt
		case string:
			log.parseTime(p.TimeFormats, t)
		}
		delete(j, jsonKey)
	}
	if jsonKey, ok := p.JsonMapping[KeyArgs]; ok {
		log.injectArgs(j[jsonKey])
	}

	if p.AppendArgs {
		log.Msg += argsToMessage(j)
	}

	return log
}

func (p *parser) parseRegexp(line []byte) *Log {
	var (
		match = p.regexp.FindSubmatch(line)
		log   = new(Log)
	)
	if len(match) == 0 {
		return nil
	}
	for i, key := range p.regexp.SubexpNames() {
		if i == 0 || i >= len(match) {
			continue
		}
		value := string(match[i])
		switch key {
		case KeyMsg:
			log.Msg = value
		case KeyLevel:
			log.Level = value
		case KeyTime:
			log.parseTime(p.TimeFormats, value)
		case KeyArgs:
			log.injectArgs(value)
		}
	}
	return log
}

func argsToMessage(j map[string]interface{}) string {
	buf := bytes.NewBuffer(nil)
	for key, val := range j {
		buf.WriteString(fmt.Sprintf(" %v=%v", key, val))
	}
	return buf.String()
}
