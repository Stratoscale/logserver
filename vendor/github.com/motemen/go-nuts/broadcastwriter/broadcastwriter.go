package broadcastwriter

import (
	"bytes"
)

const chanBufSize = 256

type BroadcastWriter struct {
	backlog *bytes.Buffer

	listeners    []chan<- []byte
	addListenerC chan chan<- []byte

	writeC chan []byte

	closeC chan struct{}
	closed bool
}

func NewBroadcastWriter() *BroadcastWriter {
	bw := &BroadcastWriter{
		backlog:      new(bytes.Buffer),
		listeners:    []chan<- []byte{},
		addListenerC: make(chan chan<- []byte),
		writeC:       make(chan []byte),
		closeC:       make(chan struct{}),
	}
	go bw.loop()
	return bw
}

func (bw *BroadcastWriter) loop() {
	for {
		select {
		case l := <-bw.addListenerC:
			if bw.backlog.Len() > 0 {
				l <- bw.backlog.Bytes()
			}
			if bw.closed {
				close(l)
			}
			bw.listeners = append(bw.listeners, l)

		case buf := <-bw.writeC:
			bw.backlog.Write(buf)
			for _, l := range bw.listeners {
				l <- buf
			}

		case <-bw.closeC:
			for _, l := range bw.listeners {
				close(l)
			}
			bw.closed = true
		}
	}
}

func (bw *BroadcastWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	bw.writeC <- buf
	return len(p), nil
}

func (bw *BroadcastWriter) Close() error {
	bw.closeC <- struct{}{}
	return nil
}

func (bw *BroadcastWriter) NewListener() <-chan []byte {
	l := make(chan []byte, chanBufSize)
	bw.addListenerC <- l
	return l
}
