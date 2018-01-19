// Package giturl provides ParseGitURL which parses remote URLs under the way
// that git does.
package giturl

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	rxURLLike     = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9+.-]*://`)
	rxHostAndPort = regexp.MustCompile(`^([^:]+|\[.+?\]):([0-9]+)$`)
	rxSCPLikeV6   = regexp.MustCompile(`^(.+?@)?\[(.+?)\]:(.*)`)
)

func ParseGitURL(giturl string) (proto string, host string, port uint, path string, exotic bool, err error) {
	// ref: parse_connect_url() in connect.c

	if rxURLLike.MatchString(giturl) {
		var u *url.URL
		u, err = url.Parse(giturl)
		if err != nil {
			return
		}

		proto = u.Scheme
		if proto == "git+ssh" || proto == "ssh+git" {
			proto = "ssh"
		}

		host = u.Host
		path = u.Path

		if proto == "ssh" {
			if m := rxHostAndPort.FindStringSubmatch(host); m != nil {
				if port64, err := strconv.ParseUint(m[2], 10, 16); err == nil {
					host = m[1]
					port = uint(port64)
				}
			}
			if host[0] == '[' && host[len(host)-1] == ']' {
				host = host[1 : len(host)-1]
			}
		}

		if u.User != nil {
			host = u.User.String() + "@" + host
		}

		if proto == "git" || proto == "ssh" {
			if path[1] == '~' {
				path = path[1:]
			}
		} else if proto == "file" {
			host = ""
			path = u.Host + u.Path
		} else {
			exotic = true
		}
	} else {
		colon := strings.IndexByte(giturl, ':')
		slash := strings.IndexByte(giturl, '/')

		if colon > -1 && (slash == -1 || colon < slash) /*&& !hasDosDrivePrefix(giturl)*/ {
			// For SCP-like URLs, colon must appear and be before any slashes
			// - user@host.xyz:path/to/repo.git/
			// - host.xyz:path/to/repo.git/
			// - user@[::1]:path/to/repo.git/
			// - [::1]:path/to/repo.git/
			proto = "ssh"
			m := rxSCPLikeV6.FindStringSubmatch(giturl)
			if m != nil {
				host = m[1] + m[2]
				path = m[3]
			} else {
				host = giturl[:colon]
				path = giturl[colon+1:]
			}
			if path[1] == '~' {
				path = path[1:]
			}
		} else {
			proto = "file"
			path = giturl
		}
	}

	return
}
