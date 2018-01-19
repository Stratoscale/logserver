// +build go1.7

package ctxlog

import (
	"testing"

	"bytes"
	"context"
)

func testPrefix(t *testing.T, ctx context.Context, expect string) {
	var buf bytes.Buffer

	logger := LoggerFromContext(ctx)
	logger.SetOutput(&buf)
	logger.SetFlags(0)
	Infof(ctx, "")

	if o := buf.String(); o[:len(o)-len("info: \n")] != expect {
		t.Errorf("prefix should be %q: got %q", expect, o[:len(o)-1])
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.Background(), "prefix: ")
	testPrefix(t, ctx, "prefix: ")

	{
		ctx := NewContext(ctx, "foo: ")
		testPrefix(t, ctx, "prefix: foo: ")
	}

	testPrefix(t, ctx, "prefix: ")
}
