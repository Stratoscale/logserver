package roundtime

import (
	"testing"

	"time"
)

func TestDuration(t *testing.T) {
	cases := []struct {
		from string
		prec int
		to   string
	}{
		{from: "12.863722988s", prec: 2, to: "12.86s"},
		{from: "225.128274ms", prec: 3, to: "225.128ms"},
		{from: "13m35.436221022s", prec: 1, to: "13m35.4s"},
		{from: "13m35.436221022s", prec: 1, to: "13m35.4s"},
		{from: "2.039480015s", prec: 0, to: "2s"},
		{from: "90.112µs", prec: 5, to: "90.112µs"},
		{from: "2h56m24s", prec: 2, to: "2h56m24s"},
	}

	for _, c := range cases {
		original, _ := time.ParseDuration(c.from)
		rounded := Duration(original, c.prec)
		if got, expected := rounded.String(), c.to; got != expected {
			t.Errorf("%q precision %d -> %q, expected %q", original, c.prec, rounded, expected)
		}
	}
}
