package giturl

import "testing"

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkParseGitURL(t *testing.T, url string) {
	t.Logf("URL: %s", url)

	got := map[string]string{}
	proto, host, port, path, exotic, err := ParseGitURL(url)
	if err != nil {
		t.Fatal(err)
	}
	if exotic == true {
		t.Errorf("protocol should not be exotic: %q", url)
	}

	got["url"] = url
	got["protocol"] = proto
	got["path"] = path

	if proto == "ssh" {
		got["userandhost"] = host
	} else {
		got["hostandport"] = host
	}

	if port != 0 {
		got["port"] = fmt.Sprint(port)
	} else {
		got["port"] = "NONE"
	}

	b, err := exec.Command("git", "fetch-pack", "--diag-url", url).Output()
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "Diag: ") {
			continue
		}
		kv := strings.SplitN(line[len("Diag: "):], "=", 2)
		expected[kv[0]] = kv[1]
	}

	for k, v := range expected {
		if got[k] != v {
			t.Errorf("%s expected %q but got %q", k, v, got[k])
		}
	}
}

func TestMain(m *testing.M) {
	cmd := exec.Command("git", "version")
	b, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	version := string(b[:len(b)-1])
	if !strings.HasPrefix(version, "git version 2.") {
		fmt.Fprintf(os.Stderr, "git version ~2 required: got [%s]\n", version)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "# %s\n", b[:len(b)-1])
	os.Exit(m.Run())
}

// source: t/t5500-fetch-pack.sh
func TestParseGitURL(t *testing.T) {
	for _, repo := range []string{"repo", "re:po", "re/po"} {
		for _, proto := range []string{"ssh+git", "git+ssh", "git", "ssh"} {
			for _, host := range []string{"host", "user@host", "user@[::1]", "user@::1"} {
				checkParseGitURL(t, fmt.Sprintf("%s://%s/%s", proto, host, repo))
				checkParseGitURL(t, fmt.Sprintf("%s://%s/~%s", proto, host, repo))
			}
			for _, host := range []string{"host", "User@host", "User@[::1]"} {
				checkParseGitURL(t, fmt.Sprintf("%s://%s:22/%s", proto, host, repo))
			}
		}
		for _, proto := range []string{"file"} {
			checkParseGitURL(t, fmt.Sprintf("%s:///%s", proto, repo))
			checkParseGitURL(t, fmt.Sprintf("%s:///~%s", proto, repo))
		}
		for _, host := range []string{"nohost", "nohost:12", "[::1]", "[::1]:23", "[", "[:aa"} {
			checkParseGitURL(t, fmt.Sprintf("./%s:%s", host, repo))
			checkParseGitURL(t, fmt.Sprintf("./:%s/~%s", host, repo))
		}
		for _, host := range []string{"host", "[::1]"} {
			checkParseGitURL(t, fmt.Sprintf("%s:%s", host, repo))
			checkParseGitURL(t, fmt.Sprintf("%s:/~%s", host, repo))
		}
	}
}

func TestParseGitURL_ExtraSCPLike(t *testing.T) {
	for _, repo := range []string{"repo", "re:po", "re/po"} {
		for _, host := range []string{"user@host", "user@[::1]"} {
			checkParseGitURL(t, fmt.Sprintf("%s:%s", host, repo))
			checkParseGitURL(t, fmt.Sprintf("%s:/~%s", host, repo))
		}
	}
}

func TestParseGitURL_HTTP(t *testing.T) {
	for _, repo := range []string{"repo", "re:po", "re/po"} {
		for _, host := range []string{"host", "host:80"} {
			for _, proto := range []string{"http", "https"} {
				p, _, _, _, exotic, err := ParseGitURL(fmt.Sprintf("%s://%s/%s", proto, host, repo))
				if err != nil {
					t.Error(err)
				} else if p != proto {
					t.Errorf("expected protocol %q but got %q", proto, p)
				} else if exotic == false {
					t.Errorf("protocol %q should be exotic", p)
				}
			}
		}
	}
}
