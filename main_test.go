package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"time"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/parser"
	"github.com/Stratoscale/logserver/ws"
	"github.com/gorilla/websocket"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestWS_GetFileTree(t *testing.T) {
	require := require.New(t)
	cwd, err := os.Getwd()
	require.Nil(err)
	cfg, err := config.New(config.FileConfig{
		Sources: []config.SourceConfig{
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

	var resp ws.Response
	require.Nil(conn.ReadJSON(&resp))

	log.Print(resp)

}

func TestWS_GetContentStratolog(t *testing.T) {
	cwd, err := os.Getwd()
	require.Nil(t, err)
	cfg, err := config.New(config.FileConfig{
		Sources: []config.SourceConfig{
			{Name: "node1", URL: fmt.Sprintf("file://%s/example/log1", cwd)},
			{Name: "node2", URL: fmt.Sprintf("file://%s/example/log2", cwd)},
		},
	})
	require.Nil(t, err)

	h := router(*cfg)
	s := httptest.NewServer(h)
	defer s.Close()

	conn, httpResp, err := websocket.DefaultDialer.Dial("ws://"+s.Listener.Addr().String()+"/ws", nil)
	require.Nil(t, err)
	assert.Equal(t, httpResp.StatusCode, http.StatusSwitchingProtocols)

	tests := []struct {
		name    string
		message string
		want    ws.Response
	}{
		{
			name:    "get content",
			message: `{"meta":{"action":"get-content","id":9},"path":["mancala.stratolog"]}`,
			want: ws.Response{
				Metadata: ws.Metadata{ID: 9, Action: "get-content", FS: "node1", Path: "mancala.stratolog"},
				Lines: []parser.LogLine{
					{Msg: "data disk <disk: hostname=stratonode1.node.strato, ID=dce9381a-cada-434d-a1ba-4e351f4afcbb, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "mancala.stratolog", LineNumber: 1, Offset: 0},
					{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "mancala.stratolog", LineNumber: 2, Offset: 699},
					{Msg: "data disk <disk: hostname=stratonode0.node.strato, ID=f3d510c7-1185-4942-b349-0de055165f78, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "mancala.stratolog", LineNumber: 3, Offset: 1398},
				},
				Error: "",
			},
		},
		{
			name:    "search",
			message: `{"meta":{"action":"search","id":9},"path":[], "regexp": "2d03c436-c197-464f-9ad0-d861e650cd61"}`,
			want: ws.Response{
				Metadata: ws.Metadata{ID: 9, Action: "search", FS: "node1", Path: "/mancala.stratolog"},
				Lines: []parser.LogLine{
					{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "/mancala.stratolog", LineNumber: 2, Offset: 699},
				},
			},
		},
		{
			name:    "search regexp",
			message: `{"meta":{"action":"search","id":9},"path":[], "regexp": "2d03c436-[c197]+-464f-9ad0-d861e650cd61"}`,
			want: ws.Response{
				Metadata: ws.Metadata{ID: 9, Action: "search", FS: "node1", Path: "/mancala.stratolog"},
				Lines: []parser.LogLine{
					{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True", Level: "INFO", Time: "2017-12-25 16:23:05 +0200 IST", FS: "node1", FileName: "/mancala.stratolog", LineNumber: 2, Offset: 699},
				},
			},
		},
	}

	for _, tt := range tests {
		// TODO: use t.Run(tt.name, func(t *testing.T), there is a bug here if one test fail on timeout the next test may fail as well
		require.Nil(t, conn.WriteMessage(1, []byte(tt.message)))
		select {
		case got := <-get(t, conn):
			assert.Equal(t, tt.want, got)
		case <-time.After(time.Millisecond * 1000):
			t.Fatal("no response!")
		}
	}
}

func get(t *testing.T, conn *websocket.Conn) <-chan ws.Response {
	ch := make(chan ws.Response)
	go func() {
		var got ws.Response
		require.Nil(t, conn.ReadJSON(&got))
		ch <- got
		//close(ch)
	}()
	return ch
}
