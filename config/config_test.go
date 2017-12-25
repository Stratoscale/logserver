package config
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestParseAddress(t *testing.T) {
    // exampleAddress := "file:///basepath"
    exampleAddress := "postgres://user:pass@host.com:5432/path?k=v#f"
    context, err := parseAddress(exampleAddress)
    if err != nil {
        panic(err)
    }

    assert.Equal(t, context.Scheme, "postgres")
}

func TestLocalFSConfig(t *testing.T) {

}
