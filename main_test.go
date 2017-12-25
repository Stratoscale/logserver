package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Stratoscale/logserver/config"
	"github.com/gorilla/websocket"
	"github.com/test-go/testify/require"
)

func TestRun(t *testing.T) {
	cwd, err := os.Getwd()
	require.Nil(t, err)
	cfg, err := config.New(config.FileConfig{
		Sources: []config.SrcDesc{
			{Name: "node1", URL: fmt.Sprintf("file://%s/example/log1", cwd)},
			{Name: "node2", URL: fmt.Sprintf("file://%s/example/log2", cwd)},
		},
	})
	require.Nil(t, err)

	h := router(*cfg)
	s := httptest.NewServer(h)
	defer s.Close()

	websocket.NewClient()

}
