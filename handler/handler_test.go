package handler

import (
	"net/http"
	"testing"

	"os"

	"fmt"

	"path/filepath"

	"encoding/json"

	"bytes"

	"github.com/Stratoscale/logserver/config"
	"github.com/posener/wstest"
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

type vfs struct {
}

func (vfs) ReadDir(dirname string) ([]os.FileInfo, error) {
	return nil, nil
}

func (vfs) Lstat(name string) (os.FileInfo, error) {
	return nil, nil
}

func (vfs) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func TestUnmarshal(t *testing.T) {
	a := `{"meta":{"action":"get-file-tree","id":1},"base_path":["a"]}`
	var r Request
	err := json.NewDecoder(bytes.NewBufferString(a)).Decode(&r)
	require.Nil(t, err)
	t.Log(r)

}
