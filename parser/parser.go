package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type debugLevel string

const (
	levelDebug   debugLevel = "debug"
	levelInfo    debugLevel = "info"
	levelError   debugLevel = "error"
	levelWarning debugLevel = "warning"
)

type LogLine struct {
	Msg        string     `json:"msg"`
	Level      debugLevel `json:"level"`
	Time       string     `json:"time"`
	FS         string     `json:"fs"`
	FileName   string     `json:"file_name"`
	LineNumber int        `json:"line_number"`
	Offset     int        `json:"offset"`
}

type Parser func(line string) (*LogLine, error)

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
	Msg   string        `json:"msg"`
	Level debugLevel    `json:"levelname"`
	Time  float64       `json:"created"`
	Args  []interface{} `json:"args"`
}

func stratologParser(line string) (*LogLine, error) {
	var stratoFormat stratologFormat
	err := json.Unmarshal([]byte(line), &stratoFormat)
	if err != nil {
		log.Printf("Failed to pars line: %s", line)
		return nil, err
	}

	r := strings.NewReplacer("%s", "%v")
	stratoFormat.Msg = r.Replace(stratoFormat.Msg)

	return &LogLine{
		Msg:   fmt.Sprintf(stratoFormat.Msg, stratoFormat.Args...),
		Level: stratoFormat.Level,
		Time:  time.Unix(int64(stratoFormat.Time), int64(stratoFormat.Time-float64(int64(stratoFormat.Time)))).String(),
	}, nil
}

func defaultParser(line string) (*LogLine, error) {
	return &LogLine{
		Msg: line,
	}, nil
}
