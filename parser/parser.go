package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type debugLevel string

var keyword = regexp.MustCompile(`(%\(([^)]+\))[diouxXeEfFgGcrs])`)

type LogLine struct {
	Msg        string     `json:"msg"`
	Level      debugLevel `json:"level"`
	Time       string     `json:"time"`
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

type stratologFormat struct {
	Msg   string      `json:"msg"`
	Level debugLevel  `json:"levelname"`
	Time  float64     `json:"created"`
	Args  interface{} `json:"args"`
}

func stratologParser(line []byte) (*LogLine, error) {
	var stratoFormat stratologFormat
	err := json.Unmarshal(line, &stratoFormat)
	if err != nil {
		log.Printf("Failed to pars line: %s", line)
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

	return &LogLine{
		Msg:   msg,
		Level: stratoFormat.Level,
		Time:  time.Unix(int64(stratoFormat.Time), int64(stratoFormat.Time-float64(int64(stratoFormat.Time)))).String(),
	}, nil
}

func defaultParser(line []byte) (*LogLine, error) {
	// TODO: try to parser time
	return &LogLine{
		Msg: string(line),
	}, nil
}
