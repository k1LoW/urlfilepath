package urlfilepath

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

const queryDelim = "?"
const schemeDelim = ":/"
const pathRoot = "___"
const pathTrailing = "____"

// Encode a URL to a file path.
func Encode(u *url.URL) (string, error) {
	p := []string{}
	if u.Scheme != "" || u.RawQuery != "" {
		p = []string{url.PathEscape(fmt.Sprintf("%s%s%s%s%s", queryDelim, u.RawQuery, queryDelim, u.Scheme, schemeDelim))}
	}
	if u.Host != "" {
		p = append(p, url.PathEscape(u.Host))
	}
	if u.Path != "" {
		splitted := strings.Split(u.Path, "/")
		for i, up := range splitted {
			if up == "" {
				switch i {
				case 0:
					p = append(p, url.PathEscape(pathRoot))
				case len(splitted) - 1:
					p = append(p, url.PathEscape(pathTrailing))
				}
				continue
			}
			p = append(p, url.PathEscape(up))
		}
	}
	return filepath.Join(p...), nil
}

// Decode a URL from a file path.
func Decode(pathstr string) (*url.URL, error) {
	uu := []string{}
	var rq string
	for i, pp := range strings.Split(pathstr, string(filepath.Separator)) {
		if pp == pathRoot {
			switch i {
			case 0, 1:
				uu = append(uu, "")
			}
			continue
		}
		if pp == pathTrailing {
			uu = append(uu, "")
			continue
		}
		u, err := url.PathUnescape(pp)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(u, queryDelim) && strings.HasSuffix(u, schemeDelim) {
			splitted := strings.Split(u, queryDelim)
			rq = splitted[1]
			if splitted[2] != schemeDelim {
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
