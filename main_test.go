package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/handler"
	"github.com/Stratoscale/logserver/parser"
	"github.com/gorilla/websocket"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestWS_GetFileTree(t *testing.T) {
	require := require.New(t)
	cwd, err := os.Getwd()
	require.Nil(err)
	cfg, err := config.New(config.FileConfig{
		Sources: []config.SrcDesc{
			{Name: "node1", URL: fmt.Sprintf("file://%s/example/log1", cwd)},
			{Name: "node2", URL: fmt.Sprintf("file://%s/example/log2", cwd)},
		},
	})
	require.Nil(err)

	h := router(*cfg)
	s := httptest.NewServer(h)
	defer s.Close()

	conn, httpResp, err := websocket.DefaultDialer.Dial("ws://"+s.Listener.Addr().String()+"/ws", nil)
	require.Nil(err)
	assert.Equal(t, httpResp.StatusCode, http.StatusSwitchingProtocols)

	require.Nil(conn.WriteMessage(1, []byte(`{"meta":{"action":"get-file-tree","id":7},"base_path":[]}`)))

	var resp handler.ResponseFileTree
	require.Nil(conn.ReadJSON(&resp))

	log.Print(resp)

}

func TestWS_GetContentStratolog(t *testing.T) {
	require := require.New(t)
	cwd, err := os.Getwd()
	require.Nil(err)
	cfg, err := config.New(config.FileConfig{
		Sources: []config.SrcDesc{
			{Name: "node1", URL: fmt.Sprintf("file://%s/example/log1", cwd)},
			{Name: "node2", URL: fmt.Sprintf("file://%s/example/log2", cwd)},
		},
	})
	require.Nil(err)

	h := router(*cfg)
	s := httptest.NewServer(h)
	defer s.Close()

	conn, httpResp, err := websocket.DefaultDialer.Dial("ws://"+s.Listener.Addr().String()+"/ws", nil)
	require.Nil(err)
	assert.Equal(t, httpResp.StatusCode, http.StatusSwitchingProtocols)

	require.Nil(conn.WriteMessage(1, []byte(`{"meta":{"action":"get-content","id":9},"path":["mancala.stratolog"]}`)))

	var resp handler.ContentResponse
	require.Nil(conn.ReadJSON(&resp))

	expResp := handler.ContentResponse{
		Metadata: handler.Metadata{ID: 9, Action: "get-content"},
		Lines: []parser.LogLine{
			{Msg: "data disk %s was found in distributionID:%s table version:%s, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "mancala.stratolog", LineNumber: 1, Offset: 0},
			{Msg: "data disk %s was found in distributionID:%s table version:%s, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "mancala.stratolog", LineNumber: 2, Offset: 699},
			{Msg: "data disk %s was found in distributionID:%s table version:%s, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "mancala.stratolog", LineNumber: 3, Offset: 1398},
		},
	}
	log.Print(resp)
	assert.Equal(t, expResp, resp)
}
