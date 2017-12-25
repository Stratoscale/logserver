package targz

import (
	"testing"

	"github.com/test-go/testify/assert"
)

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
