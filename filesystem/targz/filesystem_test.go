package targz

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"io/ioutil"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func Test_wrapper(t *testing.T) {
	cwd, err := os.Getwd()
	require.Nil(t, err)
	_url := fmt.Sprintf("file://%s/../../example/log3", cwd)
	parsedUrl, err := url.Parse(_url)
	require.Nil(t, err)
	var fs filesystem.FileSystem
	fs, err = filesystem.NewLocalFS(parsedUrl)
	if err != nil {
		panic(err)
	}
	fs = New(fs)
	tests := []struct {
		path        string
		wantErr     bool
		wantContent string
	}{
		{
			path:    "dir2/logs.tar.gz",
			wantErr: true,
		},
		{
			path:    "dir2/logs.tar.gz/first/second/third/tar_service_doesnt_exists.log",
			wantErr: true,
		},
		{
			path:        "dir2/logs.tar.gz/first/second/third/tar_service.log",
			wantContent: "blabla\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, err := fs.Open(tt.path)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			gotContent, err := ioutil.ReadAll(f)
			require.Nil(t, err)
			assert.Equal(t, tt.wantContent, string(gotContent))
		})
	}
}

func Test_isInDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		dirname string
		want    bool
	}{
		{
			name:    "/a/b",
			dirname: "/a",
			want:    true,
		},
	}

	for _, tt := range tests {
		got := isInDir(tt.dirname, tt.name)
		assert.Equal(t, tt.want, got)
	}

}
