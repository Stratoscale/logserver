package logwriter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
)

type LogWriter struct {
	Logger     *log.Logger
	Format     string
	FormatArgs []interface{}
	Calldepth  int

	buf []byte
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	var buf []byte
	if lw.buf == nil {
		buf = make([]byte, len(lw.buf)+len(p))
		copy(buf, p)
	} else {
		buf = append(lw.buf, p...)
	}

	for len(buf) > 0 {
		n := bytes.IndexByte(buf, '\n')
		if n == -1 {
			lw.buf = buf
			break
		}

		lw.writeln(buf[0:n+1], 0)
		buf = buf[n+1:]
	}

	lw.buf = buf

	return len(p), nil
}

func (lw *LogWriter) Close() error {
	if len(lw.buf) > 0 {
		lw.writeln(lw.buf, -1)
		lw.buf = nil
	}

	return nil
}

func (lw *LogWriter) ReadFrom(r io.Reader) (n int64, err error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		lw.writeln(s.Bytes(), 0)
	}

	return 0, s.Err()
}

func (lw LogWriter) writeln(line []byte, delta int) {
	var s string
	if lw.Format == "" {
		s = fmt.Sprintln(string(line))
	} else {
		args := lw.FormatArgs
		if args == nil {
			args = []interface{}{string(line)}
		} else {
			args = append(args, string(line))
		}
		s = fmt.Sprintf(lw.Format, args...)
	}

	calldepth := lw.Calldepth
	if calldepth == 0 {
		calldepth = 4
	}

	lw.Logger.Output(calldepth+delta, s)
}
