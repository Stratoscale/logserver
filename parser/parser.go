package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type debugLevel string

var keyword = regexp.MustCompile(`(%\(([^)]+\))[diouxXeEfFgGcrs])`)

type LogLine struct {
	Msg        string     `json:"msg"`
	Level      debugLevel `json:"level"`
	Time       *time.Time `json:"time,omitempty"`
	FS         string     `json:"fs"`
	FileName   string     `json:"file_name"`
	LineNumber int        `json:"line_number"`
	Offset     int        `json:"offset"`
}

type Parser func(line []byte) (*LogLine, error)

func GetParser(suffix string) Parser {
	parser := parsers[suffix]
	if parser == nil {
		parser = defaultParser
	}
	return parser
}

var parsers = map[string]Parser{
	".stratolog": stratologParser,
}

type stratoFormat struct {
	Msg   string      `json:"msg"`
	Level debugLevel  `json:"levelname"`
	Time  float64     `json:"created"`
	Args  interface{} `json:"args"`
}

func stratologParser(line []byte) (*LogLine, error) {
	var stratoFormat stratoFormat
	err := json.Unmarshal(line, &stratoFormat)
	if err != nil {
		return nil, err
	}

	r := strings.NewReplacer("%s", "%v")
	stratoFormat.Msg = r.Replace(stratoFormat.Msg)

	msg := ""
	switch args := stratoFormat.Args.(type) {
	case []interface{}:
		msg = fmt.Sprintf(stratoFormat.Msg, args...)
	case map[string]interface{}:
		msg = keyword.ReplaceAllStringFunc(stratoFormat.Msg, func(src string) string {
			key := src[2 : len(src)-2]
			val, ok := args[key]
			if !ok {
				return src
			}
			return fmt.Sprintf("%v", val)
		})
	}

	t := time.Unix(int64(stratoFormat.Time), int64(stratoFormat.Time-float64(int64(stratoFormat.Time))))

	return &LogLine{
		Msg:   msg,
		Level: stratoFormat.Level,
		Time:  &t,
	}, nil
}

func defaultParser(line []byte) (*LogLine, error) {
	return &LogLine{
		Msg:  string(line),
		Time: tryParseTime(line),
	}, nil
}
func tryParseTime(b []byte) *time.Time {
	const format = "2018-01-06 11:18:45,880"
	var prefixLen = len([]byte(format))
	if len(b) < prefixLen {
		return nil
	}
	timePhrase := b[:prefixLen]
	if t, err := time.Parse("2006-01-02 15:04:05.999", string(timePhrase)); err == nil {
		return &t
	}
	return nil
}
