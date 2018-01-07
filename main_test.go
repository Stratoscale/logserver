package main

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Stratoscale/logserver/config"
	"github.com/Stratoscale/logserver/parser"
	"github.com/Stratoscale/logserver/router"
	"github.com/Stratoscale/logserver/ws"
	"github.com/gorilla/websocket"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func mustParseTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func TestWS(t *testing.T) {
	cfg, err := config.New(config.FileConfig{
		Sources: []config.SourceConfig{
			{Name: "node1", URL: "file://./example/log1"},
			{Name: "node2", URL: "file://./example/log2"},
		},
	})
	require.Nil(t, err)

	h, err := router.New(*cfg)
	require.Nil(t, err)

	s := httptest.NewServer(h)
	defer s.Close()

	conn, httpResp, err := websocket.DefaultDialer.Dial("ws://"+s.Listener.Addr().String()+"/ws", nil)
	require.Nil(t, err)
	assert.Equal(t, httpResp.StatusCode, http.StatusSwitchingProtocols)

	tests := []struct {
		name    string
		message string
		want    []ws.Response
	}{
		{
			name:    "get content",
			message: `{"meta":{"action":"get-content","id":9},"path":["mancala.stratolog"]}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "get-content", FS: "node1", Path: ws.Path{"mancala.stratolog"}},
					Lines: []parser.LogLine{
						{
							Msg:        "data disk <disk: hostname=stratonode1.node.strato, ID=dce9381a-cada-434d-a1ba-4e351f4afcbb, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:      "INFO",
							Time:       mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:         "node1",
							FileName:   "mancala.stratolog",
							LineNumber: 1,
							Offset:     0,
						},
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:      "INFO",
							Time:       mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:         "node1",
							FileName:   "mancala.stratolog",
							LineNumber: 2,
							Offset:     699,
						},
						{Msg: "data disk <disk: hostname=stratonode0.node.strato, ID=f3d510c7-1185-4942-b349-0de055165f78, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:      "INFO",
							Time:       mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:         "node1",
							FileName:   "mancala.stratolog",
							LineNumber: 3,
							Offset:     1398,
						},
					},
				},
			},
		},
		{
			name:    "get content / empty file",
			message: `{"meta":{"action":"get-content","id":9},"path":["service2.log"]}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "get-content", FS: "node1", Path: ws.Path{"service2.log"}},
				},
			},
		},
		{
			name:    "get content / content-file empty file combination",
			message: `{"meta":{"action":"get-content","id":9},"path":["service1.log"]}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "get-content", FS: "node1", Path: ws.Path{"service1.log"}},
					Lines: []parser.LogLine{
						{Msg: "find me", LineNumber: 1, FileName: "service1.log", FS: "node1"},
					},
				},
				{
					Meta: ws.Meta{ID: 9, Action: "get-content", FS: "node2", Path: ws.Path{"service1.log"}},
				},
			},
		},
		{
			name:    "search",
			message: `{"meta":{"action":"search","id":9},"path":[], "regexp": "2d03c436-c197-464f-9ad0-d861e650cd61"}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "search", FS: "node1", Path: ws.Path{"mancala.stratolog"}},
					Lines: []parser.LogLine{
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:      "INFO",
							Time:       mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:         "node1",
							FileName:   "mancala.stratolog",
							LineNumber: 2,
							Offset:     699,
						},
					},
				},
			},
		},
		{
			name:    "search / filter node",
			message: `{"meta":{"action":"search","id":9},"path":[], "regexp": "2d03c436-c197-464f-9ad0-d861e650cd61", "filter_fs": ["node1"]}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "search", FS: "node1", Path: ws.Path{"mancala.stratolog"}},
					Lines: []parser.LogLine{
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:      "INFO",
							Time:       mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:         "node1",
							FileName:   "mancala.stratolog",
							LineNumber: 2,
							Offset:     699,
						},
					},
				},
			},
		},
		{
			name:    "search regexp",
			message: `{"meta":{"action":"search","id":9},"path":[], "regexp": "2d03c436-[c197]+-464f-9ad0-d861e650cd61"}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "search", FS: "node1", Path: ws.Path{"mancala.stratolog"}},
					Lines: []parser.LogLine{
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:      "INFO",
							Time:       mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:         "node1",
							FileName:   "mancala.stratolog",
							LineNumber: 2,
							Offset:     699,
						},
					},
				},
			},
		},
		{
			name:    "get file tree",
			message: `{"meta":{"action":"get-file-tree","id":9},"base_path":[]}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "get-file-tree"},
					Tree: []*ws.File{
						{
							Key:       "dir1",
							Path:      ws.Path{"dir1"},
							IsDir:     true,
							Instances: []ws.FileInstance{{Size: 4096, FS: "node1"}},
						},
						{
							Key:       "dir1/service3.log",
							Path:      ws.Path{"dir1", "service3.log"},
							IsDir:     false,
							Instances: []ws.FileInstance{{Size: 0, FS: "node1"}},
						},
						{
							Key:       "mancala.stratolog",
							Path:      ws.Path{"mancala.stratolog"},
							IsDir:     false,
							Instances: []ws.FileInstance{{Size: 2100, FS: "node1"}},
						},
						{
							Key:   "service1.log",
							Path:  ws.Path{"service1.log"},
							IsDir: false,
							Instances: []ws.FileInstance{
								{Size: 7, FS: "node1"},
								{Size: 0, FS: "node2"},
							},
						},
						{
							Key:       "service2.log",
							Path:      ws.Path{"service2.log"},
							IsDir:     false,
							Instances: []ws.FileInstance{{Size: 0, FS: "node1"}},
						},
					},
				},
			},
		},
		{
			name:    "get file tree/filter node",
			message: `{"meta":{"action":"get-file-tree","id":9},"base_path":[], "filter_fs": ["node2"]}`,
			want: []ws.Response{
				{
					Meta: ws.Meta{ID: 9, Action: "get-file-tree"},
					Tree: []*ws.File{
						{
							Key:       "service1.log",
							Path:      ws.Path{"service1.log"},
							IsDir:     false,
							Instances: []ws.FileInstance{{Size: 0, FS: "node2"}},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		require.Nil(t, conn.WriteMessage(1, []byte(tt.message)))
		var got []ws.Response
		for i := 0; i < len(tt.want); i++ {
			select {
			case gotOne := <-get(t, conn):
				got = append(got, gotOne)
			case <-time.After(time.Second):
				t.Fatal("no response!")
			}
		}
		sort.Slice(got, func(i, j int) bool { return strings.Compare(got[i].Meta.FS, got[j].Meta.FS) == -1 })
		sort.Slice(tt.want, func(i, j int) bool { return strings.Compare(tt.want[i].Meta.FS, tt.want[j].Meta.FS) == -1 })
		assert.Equal(t, tt.want, got)
	}
}

func get(t *testing.T, conn *websocket.Conn) <-chan ws.Response {
	ch := make(chan ws.Response)
	go func() {
		var got ws.Response
		require.Nil(t, conn.ReadJSON(&got))
		ch <- got
	}()
	return ch
}
