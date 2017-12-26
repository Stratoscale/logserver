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
	"github.com/gorilla/websocket"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestWS(t *testing.T) {
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
