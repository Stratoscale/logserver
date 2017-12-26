package targz

import (
	"os"
	"fmt"
	"testing"
	"net/url"
	"github.com/test-go/testify/assert"
	"github.com/Stratoscale/logserver/filesystem"
)

func Test_wrapper(t *testing.T) {
	fmt.Println("bp5")
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_url := fmt.Sprintf("file://%s/../../example/log3", cwd)
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
	// _, err = fs.Open("dir2/logs.tar.gz/first/second/third/tar_service.log")
	_, err = fs.Open("dir2/logs.tar.gz")
	if err != nil {
		panic(err)
	}
	_, err = fs.Open("dir2/logs.tar.gz/first/second/third/tar_service_doesnt_exists.log")
	if err == nil {
		fmt.Errorf("succeed to open non existent file")
		assert.NotNil(t, err)
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
