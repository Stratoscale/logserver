package main

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/Stratoscale/logserver/engine"
	"github.com/Stratoscale/logserver/parse"
	"github.com/Stratoscale/logserver/source"
	"github.com/bluele/gcache"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func TestHandler(t *testing.T) {
	t.Parallel()
	if testing.Verbose() {
		logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	}
	cfg := loadConfig("./example/logserver.json")
	cache := gcache.New(0).Build()

	sources, err := source.New(cfg.Sources, cache)
	require.Nil(t, err)
	parser, err := parse.New(cfg.Parsers)
	require.Nil(t, err)

	s := httptest.NewServer(engine.New(engine.Config{}, sources, parser, cache))
	defer s.Close()

	tests := []struct {
		name    string
		message string
		want    []engine.Response
	}{
		{
			name:    "get content",
			message: `{"meta":{"action":"get-content","id":1},"path":["mancala.stratolog"]}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 1, Action: "get-content", FS: "node1", Path: engine.Path{"mancala.stratolog"}},
					Lines: []parse.Log{
						{
							Msg:      "data disk <disk: hostname=stratonode1.node.strato, ID=dce9381a-cada-434d-a1ba-4e351f4afcbb, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:    "INFO",
							Time:     mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:       "node1",
							FileName: "mancala.stratolog",
							Line:     1,
							Offset:   0,
						},
						{
							Msg:      "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:    "INFO",
							Time:     mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:       "node1",
							FileName: "mancala.stratolog",
							Line:     2,
							Offset:   699,
						},
						{
							Msg:      "data disk <disk: hostname=stratonode0.node.strato, ID=f3d510c7-1185-4942-b349-0de055165f78, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:    "INFO",
							Time:     mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:       "node1",
							FileName: "mancala.stratolog",
							Line:     3,
							Offset:   1398,
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 1, Action: "get-content"},
					Finished: true,
				},
			},
		},
		{
			name:    "get content / empty file",
			message: `{"meta":{"action":"get-content","id":2},"path":["service2.log"]}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 2, Action: "get-content", FS: "node1", Path: engine.Path{"service2.log"}},
				},
				{
					Meta: engine.Meta{ID: 2, Action: "get-content", FS: "node3", Path: engine.Path{"service2.log"}},
				},
				{
					Meta:     engine.Meta{ID: 2, Action: "get-content"},
					Finished: true,
				},
			},
		},
		{
			name:    "get content / content-file empty file combination",
			message: `{"meta":{"action":"get-content","id":3},"path":["service1.log"]}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 3, Action: "get-content", FS: "node1", Path: engine.Path{"service1.log"}},
					Lines: []parse.Log{
						{Msg: "find me", Line: 1, FileName: "service1.log", FS: "node1"},
					},
				},
				{
					Meta: engine.Meta{ID: 3, Action: "get-content", FS: "node2", Path: engine.Path{"service1.log"}},
				},
				{
					Meta: engine.Meta{ID: 3, Action: "get-content", FS: "node3", Path: engine.Path{"service1.log"}},
				},
				{
					Meta:     engine.Meta{ID: 3, Action: "get-content"},
					Finished: true,
				},
			},
		},
		{
			name:    "search",
			message: `{"meta":{"action":"search","id":4},"path":[], "regexp": "2d03c436-c197-464f-9ad0-d861e650cd61"}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 4, Action: "search", FS: "node1", Path: engine.Path{"mancala.stratolog"}},
					Lines: []parse.Log{
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:    "INFO",
							Time:     mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:       "node1",
							FileName: "mancala.stratolog",
							Line:     2,
							Offset:   699,
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 4, Action: "search"},
					Finished: true,
				},
			},
		},
		{
			name:    "search/long file",
			message: `{"meta":{"action":"search","id":5},"path":[], "regexp": "zzzzzzz"}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 5, Action: "search", FS: "node1", Path: engine.Path{"dir1", "service3.log"}},
					Lines: []parse.Log{
						{
							Msg:      `{"msg": "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"}`,
							FS:       "node1",
							FileName: "dir1/service3.log",
							Line:     8965,
							Offset:   977076,
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 5, Action: "search"},
					Finished: true,
				},
			},
		},
		{
			name:    "search/filter node",
			message: `{"meta":{"action":"search","id":6},"path":[],"regexp":"2d03c436-c197-464f-9ad0-d861e650cd61","filter_fs":["node1"]}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 6, Action: "search", FS: "node1", Path: engine.Path{"mancala.stratolog"}},
					Lines: []parse.Log{
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:    "INFO",
							Time:     mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:       "node1",
							FileName: "mancala.stratolog",
							Line:     2,
							Offset:   699,
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 6, Action: "search"},
					Finished: true,
				},
			},
		},
		{
			name:    "search/regexp",
			message: `{"meta":{"action":"search","id":7},"path":[], "regexp": "2d03c436-[c197]+-464f-9ad0-d861e650cd61"}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 7, Action: "search", FS: "node1", Path: engine.Path{"mancala.stratolog"}},
					Lines: []parse.Log{
						{Msg: "data disk <disk: hostname=stratonode2.node.strato, ID=2d03c436-c197-464f-9ad0-d861e650cd61, path=/dev/sdc, type=mancala> was found in distributionID:0 table version:1, setting inTable=True",
							Level:    "INFO",
							Time:     mustParseTime("2017-12-25T16:23:05+02:00"),
							FS:       "node1",
							FileName: "mancala.stratolog",
							Line:     2,
							Offset:   699,
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 7, Action: "search"},
					Finished: true,
				},
			},
		},
		{
			name:    "search/not-found",
			message: `{"meta":{"action":"search","id":8},"path":[], "regexp": "value that you won't found'"}`,
			want: []engine.Response{
				{
					Meta:     engine.Meta{ID: 8, Action: "search"},
					Finished: true,
				},
			},
		},
		{
			name:    "get file tree",
			message: `{"meta":{"action":"get-file-tree","id":9},"base_path":[],"filter_fs":["node1","node2"]}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 9, Action: "get-file-tree"},
					Files: []*engine.File{
						{
							Key:       "dir1",
							Path:      engine.Path{"dir1"},
							IsDir:     true,
							Instances: []engine.FileInstance{{Size: 4096, FS: "node1"}},
						},
						{
							Key:       "dir1/service3.log",
							Path:      engine.Path{"dir1", "service3.log"},
							IsDir:     false,
							Instances: []engine.FileInstance{{Size: 986150, FS: "node1"}},
						},
						{
							Key:       "mancala.stratolog",
							Path:      engine.Path{"mancala.stratolog"},
							IsDir:     false,
							Instances: []engine.FileInstance{{Size: 2100, FS: "node1"}},
						},
						{
							Key:   "service1.log",
							Path:  engine.Path{"service1.log"},
							IsDir: false,
							Instances: []engine.FileInstance{
								{Size: 7, FS: "node1"},
								{Size: 0, FS: "node2"},
							},
						},
						{
							Key:       "service2.log",
							Path:      engine.Path{"service2.log"},
							IsDir:     false,
							Instances: []engine.FileInstance{{Size: 0, FS: "node1"}},
						},
						{
							Key:       "journal",
							Path:      engine.Path{"journal"},
							IsDir:     false,
							Instances: []engine.FileInstance{{Size: 4096, FS: "node2"}},
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 9, Action: "get-file-tree"},
					Finished: true,
				},
			},
		},
		{
			name:    "get file tree/filter node",
			message: `{"meta":{"action":"get-file-tree","id":10},"base_path":[],"filter_fs":["node2"]}`,
			want: []engine.Response{
				{
					Meta: engine.Meta{ID: 10, Action: "get-file-tree"},
					Files: []*engine.File{
						{
							Key:       "service1.log",
							Path:      engine.Path{"service1.log"},
							IsDir:     false,
							Instances: []engine.FileInstance{{Size: 0, FS: "node2"}},
						},
						{
							Key:       "journal",
							Path:      engine.Path{"journal"},
							IsDir:     false,
							Instances: []engine.FileInstance{{Size: 4096, FS: "node2"}},
						},
					},
				},
				{
					Meta:     engine.Meta{ID: 10, Action: "get-file-tree"},
					Finished: true,
				},
			},
		},
	}

	addr := "ws://" + s.Listener.Addr().String()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			conn, httpResp, err := websocket.DefaultDialer.Dial(addr, nil)
			require.Nil(t, err)
			assert.Equal(t, httpResp.StatusCode, http.StatusSwitchingProtocols)

			t.Parallel()

			require.Nil(t, conn.WriteMessage(1, []byte(tt.message)))
			var got []engine.Response
			for i := 0; i < len(tt.want); i++ {
				gotOne := <-get(t, conn)
				got = append(got, gotOne)
			}
			sortResp(got)
			sortResp(tt.want)
			assert.Equal(t, tt.want, got)
		})
	}
}

func sortResp(responses []engine.Response) {
	sort.Slice(responses, func(i, j int) bool { return strings.Compare(responses[i].Meta.FS, responses[j].Meta.FS) == -1 })
	for _, resp := range responses {
		sort.Slice(resp.Files, func(i, j int) bool { return strings.Compare(resp.Files[i].Key, resp.Files[j].Key) == -1 })
		for _, file := range resp.Files {
			sort.Slice(file.Instances, func(i, j int) bool { return strings.Compare(file.Instances[i].FS, file.Instances[j].FS) == -1 })
		}
	}
}

func get(t *testing.T, conn *websocket.Conn) <-chan engine.Response {
	ch := make(chan engine.Response)
	go func() {
		var got engine.Response
		require.Nil(t, conn.ReadJSON(&got))
		ch <- got
	}()
	return ch
}
