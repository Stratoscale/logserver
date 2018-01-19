package ctxlog

import (
	"io"
	"log"
	"os"
	"sync"
)

var Logger = log.New(os.Stderr, "", log.LstdFlags)

var (
	outputMu sync.Mutex
	output   io.Writer = os.Stderr
)

type contextKey struct {
	name string
}

func SetOutput(w io.Writer) {
	outputMu.Lock()
	defer outputMu.Unlock()
	Logger.SetOutput(w)
	output = w
}

var LoggerContextKey = &contextKey{"logger"}
var PrefixContextKey = &contextKey{"prefix"}
