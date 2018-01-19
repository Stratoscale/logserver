package filesystem

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_readDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		want    []file
		wantErr bool
	}{
		{
			name: "directory",
			body: `<html>
<head><title>Index of /a/</title></head>
<body bgcolor="white">
<h1>Index of /a/</h1><hr><pre><a href="../">../</a>
<a href="c/">c/</a>                                                 18-Jan-2018 23:32                   -
<a href="xxx">xxx</a>                                                18-Jan-2018 23:32                   17
<a href="yyy">yyy</a>                                                18-Jan-2018 23:32                   0
</pre><hr></body>
</html>
`,
			want: []file{
				{name: "c", isDir: true, time: time.Date(2018, time.January, 18, 23, 32, 0, 0, time.UTC)},
				{name: "xxx", isDir: false, size: 17, time: time.Date(2018, time.January, 18, 23, 32, 0, 0, time.UTC)},
				{name: "yyy", isDir: false, size: 0, time: time.Date(2018, time.January, 18, 23, 32, 0, 0, time.UTC)},
			},
		},
		{
			name: "bad time format",
			body: `<html>
<head><title>Index of /a/</title></head>
<body bgcolor="white">
<h1>Index of /a/</h1><hr><pre><a href="../">../</a>
<a href="c/">c/</a>                                                 18-Jan-2018 23:32                   -
<a href="xxx">xxx</a>                                                181-Jan-2018 23:32                   17
<a href="yyy">yyy</a>                                                18-Jan-2018 23:32                   0
</pre><hr></body>
</html>
`,
			want: []file{
				{name: "c", isDir: true, time: time.Date(2018, time.January, 18, 23, 32, 0, 0, time.UTC)},
				{name: "yyy", isDir: false, size: 0, time: time.Date(2018, time.January, 18, 23, 32, 0, 0, time.UTC)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(tt.body)
			got, err := parseDirectoryHTML(ioutil.NopCloser(body))
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, len(tt.want), len(got))
			for i := 0; i < len(tt.want); i++ {
				assert.Equal(t, tt.want[i].name, got[i].Name())
				assert.Equal(t, tt.want[i].isDir, got[i].IsDir())
				assert.Equal(t, tt.want[i].time, got[i].ModTime())
				assert.Equal(t, tt.want[i].size, got[i].Size())
			}
		})
	}

}
