package urlfilepath

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// Convert a URL to a file path.
func Convert(u *url.URL) (string, error) {
	p := []string{}
	if u.Scheme != "" || u.RawQuery != "" {
		p = []string{url.PathEscape(fmt.Sprintf("?%s?%s:/", u.RawQuery, u.Scheme))}
	}
	if u.Host != "" {
		p = append(p, url.PathEscape(u.Host))
	}
	if u.Path != "" {
		for _, up := range strings.Split(u.Path, "/") {
			p = append(p, url.PathEscape(up))
		}
	}
	return filepath.Join(p...), nil
}

// Restore a URL from a file path.
func Restore(pathstr string) (*url.URL, error) {
	uu := []string{}
	var rq string
	for _, pp := range strings.Split(pathstr, string(filepath.Separator)) {
		u, err := url.PathUnescape(pp)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(u, "?") && strings.HasSuffix(u, ":/") {
			splitted := strings.Split(u, "?")
			rq = splitted[1]
			if splitted[2] != ":/" {
				uu = append(uu, splitted[2])
			}
			continue
		}
		uu = append(uu, u)
	}
	u, err := url.Parse(strings.Join(uu, "/"))
	if err != nil {
		return nil, err
	}
	if rq != "" {
		u.RawQuery = rq
	}
	return u, nil
}
