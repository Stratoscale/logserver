package sftp

import (
	"encoding"
	"io"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"

	"github.com/pkg/errors"
)

var maxTxPacket uint32 = 1 << 15

type handleHandler func(string) string

// Handlers contains the 4 SFTP server request handlers.
type Handlers struct {
	FileGet  FileReader
	FilePut  FileWriter
	FileCmd  FileCmder
	FileList FileLister
}

// RequestServer abstracts the sftp protocol with an http request-like protocol
type RequestServer struct {
	*serverConn
	Handlers        Handlers
	pktMgr          *packetManager
	openRequests    map[string]*Request
	openRequestLock sync.RWMutex
	handleCount     int
}

// NewRequestServer creates/allocates/returns new RequestServer.
// Normally there there will be one server per user-session.
func NewRequestServer(rwc io.ReadWriteCloser, h Handlers) *RequestServer {
	svrConn := &serverConn{
		conn: conn{
			Reader:      rwc,
			WriteCloser: rwc,
		},
	}
	return &RequestServer{
		serverConn:   svrConn,
		Handlers:     h,
		pktMgr:       newPktMgr(svrConn),
		openRequests: make(map[string]*Request),
	}
}

// New Open packet/Request
func (rs *RequestServer) nextRequest(r *Request) string {
	rs.openRequestLock.Lock()
	defer rs.openRequestLock.Unlock()
	rs.handleCount++
	handle := strconv.Itoa(rs.handleCount)
	rs.openRequests[handle] = r
	return handle
}

// Returns Request from openRequests, bool is false if it is missing
// If the method is different, save/return a new Request w/ that Method.
//
// The Requests in openRequests work essentially as open file descriptors that
// you can do different things with. What you are doing with it are denoted by
// the first packet of that type (read/write/etc). We create a new Request when
// it changes to set the request.Method attribute in a thread safe way.
func (rs *RequestServer) getRequest(handle, method string) (*Request, bool) {
	rs.openRequestLock.RLock()
	r, ok := rs.openRequests[handle]
	rs.openRequestLock.RUnlock()
	if !ok || r.Method == method {
		return r, ok
	}
	// if we make it here we need to replace the request
	rs.openRequestLock.Lock()
	defer rs.openRequestLock.Unlock()
	r, ok = rs.openRequests[handle]
	if !ok || r.Method == method { // re-check needed b/c lock race
		return r, ok
	}
	r = &Request{Method: method, Filepath: r.Filepath, state: r.state}
	rs.openRequests[handle] = r
	return r, ok
}

func (rs *RequestServer) closeRequest(handle string) error {
	rs.openRequestLock.Lock()
	defer rs.openRequestLock.Unlock()
	if r, ok := rs.openRequests[handle]; ok {
		delete(rs.openRequests, handle)
		return r.close()
	}
	return syscall.EBADF
}

// Close the read/write/closer to trigger exiting the main server loop
func (rs *RequestServer) Close() error { return rs.conn.Close() }

// Serve requests for user session
func (rs *RequestServer) Serve() error {
	var wg sync.WaitGroup
	runWorker := func(ch requestChan) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := rs.packetWorker(ch); err != nil {
				rs.conn.Close() // shuts down recvPacket
			}
		}()
	}
	pktChan := rs.pktMgr.workerChan(runWorker)

	var err error
	var pkt requestPacket
	var pktType uint8
	var pktBytes []byte
	for {
		pktType, pktBytes, err = rs.recvPacket()
		if err != nil {
			break
		}

		pkt, err = makePacket(rxPacket{fxp(pktType), pktBytes})
		if err != nil {
			debug("makePacket err: %v", err)
			rs.conn.Close() // shuts down recvPacket
			break
		}

		pktChan <- pkt
	}

	close(pktChan) // shuts down sftpServerWorkers
	wg.Wait()      // wait for all workers to exit

	return err
}

func (rs *RequestServer) packetWorker(pktChan chan requestPacket) error {
	for pkt := range pktChan {
		var rpkt responsePacket
		switch pkt := pkt.(type) {
		case *sshFxInitPacket:
			rpkt = sshFxVersionPacket{sftpProtocolVersion, nil}
		case *sshFxpClosePacket:
			handle := pkt.getHandle()
			rpkt = statusFromError(pkt, rs.closeRequest(handle))
		case *sshFxpRealpathPacket:
			rpkt = cleanPacketPath(pkt)
		case isOpener:
			request := requestFromPacket(pkt)
			handle := rs.nextRequest(request)
			p, ok := pkt.(*sshFxpOpenPacket)
			if ok && p.hasPflags(ssh_FXF_CREAT) {
				request.call(rs.Handlers, pkt)
			}
			rpkt = sshFxpHandlePacket{pkt.id(), handle}
		case *sshFxpFstatPacket:
			handle := pkt.getHandle()
			request, ok := rs.getRequest(handle, requestMethod(pkt))
			if !ok {
				rpkt = statusFromError(pkt, syscall.EBADF)
			} else {
				request = requestFromPacket(
					&sshFxpStatPacket{ID: pkt.id(), Path: request.Filepath})
				rpkt = request.call(rs.Handlers, pkt)
			}
		case *sshFxpFsetstatPacket:
			handle := pkt.getHandle()
			request, ok := rs.getRequest(handle, requestMethod(pkt))
			if !ok {
				rpkt = statusFromError(pkt, syscall.EBADF)
			} else {
				request = requestFromPacket(
					&sshFxpSetstatPacket{ID: pkt.id(), Path: request.Filepath,
						Flags: pkt.Flags, Attrs: pkt.Attrs,
					})
				rpkt = request.call(rs.Handlers, pkt)
			}
		case hasHandle:
			handle := pkt.getHandle()
			request, ok := rs.getRequest(handle, requestMethod(pkt))
			if !ok {
				rpkt = statusFromError(pkt, syscall.EBADF)
			} else {
				rpkt = request.call(rs.Handlers, pkt)
			}
		case hasPath:
			request := requestFromPacket(pkt)
			rpkt = request.call(rs.Handlers, pkt)
		default:
			return errors.Errorf("unexpected packet type %T", pkt)
		}

		err := rs.sendPacket(rpkt)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanPacketPath(pkt *sshFxpRealpathPacket) responsePacket {
	path := cleanPath(pkt.getPath())
	return &sshFxpNamePacket{
		ID: pkt.id(),
		NameAttrs: []sshFxpNameAttr{{
			Name:     path,
			LongName: path,
			Attrs:    emptyFileStat,
		}},
	}
}

// Makes sure we have a clean POSIX (/) absolute path to work with
func cleanPath(p string) string {
	p = filepath.ToSlash(p)
	if !path.IsAbs(p) {
		p = "/" + p
	}
	return path.Clean(p)
}

// Wrap underlying connection methods to use packetManager
func (rs *RequestServer) sendPacket(m encoding.BinaryMarshaler) error {
	if pkt, ok := m.(responsePacket); ok {
		rs.pktMgr.readyPacket(pkt)
	} else {
		return errors.Errorf("unexpected packet type %T", m)
	}
	return nil
}

func (rs *RequestServer) sendError(p ider, err error) error {
	return rs.sendPacket(statusFromError(p, err))
}
