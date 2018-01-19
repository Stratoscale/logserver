package broadcastwriter

import (
	"bytes"
	"sync"
	"testing"

	"fmt"
)

func TestBroadcastWriter(t *testing.T) {
	var wg sync.WaitGroup
	consumeListenerC := func(c <-chan []byte) <-chan string {
		outc := make(chan string)
		wg.Add(1)
		go func() {
			var buf bytes.Buffer
			for b := range c {
				buf.Write(b)
			}
			outc <- buf.String()
			wg.Done()
		}()
		return outc
	}

	bw := NewBroadcastWriter()

	l1 := bw.NewListener()
	c1 := consumeListenerC(l1)

	fmt.Fprintln(bw, "foo")

	l2 := bw.NewListener()
	c2 := consumeListenerC(l2)

	fmt.Fprintln(bw, "bar")

	l3 := bw.NewListener()
	c3 := consumeListenerC(l3)

	bw.Close()

	l4 := bw.NewListener()
	c4 := consumeListenerC(l4)

	check := func(name string, c <-chan string) {
		s := <-c
		expected := "foo\nbar\n"
		if s != expected {
			t.Errorf("%s: got %q but expected %q", name, s, expected)
		}
	}
	check("c1", c1)
	check("c2", c2)
	check("c3", c3)
	check("c4", c4)

	wg.Wait()
}
