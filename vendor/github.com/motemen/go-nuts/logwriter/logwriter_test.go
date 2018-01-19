package logwriter

import (
	"testing"

	"bytes"
	"fmt"
	"log"
)

func mustEqual(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("got %q but expected %q", got, expected)
	}
}

func TestLogWriter(t *testing.T) {
	var buf bytes.Buffer
	w := &LogWriter{
		Logger: log.New(&buf, "", log.Lshortfile),
		Format: "[test] %s",
	}

	fmt.Fprintln(w, "foo")
	fmt.Fprint(w, "bar-")

	mustEqual(t, buf.String(), "logwriter_test.go:24: [test] foo\n")

	fmt.Fprintln(w, "baz")

	mustEqual(t, buf.String(), "logwriter_test.go:24: [test] foo\nlogwriter_test.go:29: [test] bar-baz\n")

	fmt.Fprint(w, "qux")

	mustEqual(t, buf.String(), "logwriter_test.go:24: [test] foo\nlogwriter_test.go:29: [test] bar-baz\n")

	w.Close()

	mustEqual(t, buf.String(), "logwriter_test.go:24: [test] foo\nlogwriter_test.go:29: [test] bar-baz\nlogwriter_test.go:37: [test] qux\n")
}
