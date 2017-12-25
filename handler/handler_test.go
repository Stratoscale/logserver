package handler

import (
	"io"
	"net/http"
	"testing"

	"fmt"
	"os"

	"path/filepath"

	"github.com/Stratoscale/logserver/config"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	h := New(config.Config{
		Nodes: []config.Src{
			{
				Name: "node1",
				FS:   new(vfs),
			},
		},
	})
	d := wstest.NewDialer(h, nil) // or t.Log instead of nil

	c, httpResp, err := d.Dial("ws://"+"whatever"+"/ws", nil)
	require.Nil(t, err)
	defer c.Close()

	if got, want := httpResp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	err = c.WriteMessage(1, []byte(`{"meta":{"action":"get-file-tree","id":1},"base_path":["a"]}`))
	require.Nil(t, err)

	resp := fileTreeResponse{}

	err = c.ReadJSON(&resp)
	require.Nil(t, err)

	fmt.Print(resp)
}

type mockedReadClose struct {
	io.ReadCloser
	mock.Mock
}

type mockFileInfo struct {
	os.FileInfo
	mock.Mock
}

func (m mockFileInfo) IsDir() bool {
	return false
}

func (m mockFileInfo) Size() int64 {
	return 0
}

type vfs struct {
}

func (vfs) Open(path string) (io.ReadCloser, error) {
	return new(mockedReadClose), nil
}

func (vfs) ReadDir(dirname string) ([]os.FileInfo, error) {
	var infos = []os.FileInfo{mockFileInfo{}}
	return infos, nil
}

func (vfs) Lstat(name string) (os.FileInfo, error) {
	return mockFileInfo{}, nil
}

func (vfs) Join(elem ...string) string {
	return filepath.Join(elem...)
}
