package targz

import (
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/Stratoscale/logserver/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	u, err := url.Parse("file://../../example/log3")
	require.Nil(t, err)
	fs, err := filesystem.NewLocalFS(u)
	require.Nil(t, err)
	fs = New(fs)

	openTests := []struct {
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

	for _, tt := range openTests {
		t.Run("open/"+tt.path, func(t *testing.T) {
			f, err := fs.Open(tt.path)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				require.Nil(t, err)
				gotContent, err := ioutil.ReadAll(f)
				require.Nil(t, err)
				assert.Equal(t, tt.wantContent, string(gotContent))
			}
		})
	}

	dirTests := []struct {
		path string
		want []fileInfo
	}{
		{
			path: "/",
			want: []fileInfo{
				{name: "dir1", isDir: true},
				{name: "dir2", isDir: true},
				{name: "service1.log", isDir: false},
				{name: "service2.log", isDir: false},
			},
		},
		{
			path: "dir2",
			want: []fileInfo{
				{name: "logs.tar.gz", isDir: true},
			},
		},
		{
			path: "dir2/logs.tar.gz",
			want: []fileInfo{
				{name: "first", isDir: true},
			},
		},
		{
			path: "dir2/logs.tar.gz/first/second/",
			want: []fileInfo{
				{name: "third", isDir: true},
			},
		},
		{
			path: "dir2/logs.tar.gz/first/second/third",
			want: []fileInfo{
				{name: "tar_service.log", isDir: false},
			},
		},
	}

	for _, tt := range dirTests {
		t.Run("dir/"+tt.path, func(t *testing.T) {
			files, err := fs.ReadDir(tt.path)
			require.Nil(t, err)
			var gotFileInfos []fileInfo
			for _, f := range files {
				gotFileInfos = append(gotFileInfos, fileInfo{f.Name(), f.IsDir()})
			}
			sort.Slice(gotFileInfos, func(i, j int) bool { return strings.Compare(gotFileInfos[i].name, gotFileInfos[j].name) == -1 })

			assert.Equal(t, tt.want, gotFileInfos)
		})
	}
}

func Test_isInDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dirname string
		name    string
		want    bool
	}{
		{dirname: "/a/", name: "/a/b", want: true},
		{dirname: "/a", name: "/a/b/", want: true},
		{dirname: "/a/", name: "/a/b/", want: true},
		{dirname: "/a", name: "/a/b", want: true},
		{dirname: "a", name: "/a/b", want: true},
		{dirname: "/a", name: "a/b", want: true},
		{dirname: "/a", name: "/a", want: false},
		{dirname: "/a/b", name: "/a", want: false},
		{dirname: "/a/b", name: "/a", want: false},
	}

	for _, tt := range tests {
		got := isInDir(tt.dirname, tt.name)
		assert.Equal(t, tt.want, got)
	}

}

func Test_subDirs(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{path: ""},
		{path: "xxx"},
		{path: "xxx/yyy", want: []string{"xxx"}},
		{path: "xxx/yyy/zzz", want: []string{"xxx/yyy", "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, subDirs(tt.path))
		})
	}
}
