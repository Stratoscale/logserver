package targz

import (
	"fmt"
	"testing"
	"net/url"
	"github.com/test-go/testify/assert"
	"github.com/Stratoscale/logserver/filesystem"
)

func Test_wrapper(t *testing.T) {
	fmt.Println("bp5")
	_url := "file://./example/log3"
	parsedUrl, err := url.Parse(_url)
	fmt.Println("bp6")
	var fs filesystem.FileSystem
	fs, err = filesystem.NewLocalFS(parsedUrl)
	fmt.Println("bp7")
	if err != nil {
		panic(err)
	}
	fmt.Println("bp8")
	fs = New(fs)
	fs.Open("dir2/logs.tar.gz/first/second/third/tar_service.log")
}

func test_isInDir(t *testing.T) {
	fmt.Println("bp5")
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
