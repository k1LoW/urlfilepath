package urlfilepath_test

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/k1LoW/urlfilepath"
)

func TestUrlfilepath(t *testing.T) {
	tests := []struct {
		urlstr string
	}{
		{"https://github.com/k1LoW/urlfilepath"},
		{"https://github.com/k1LoW/urlfilepath/"},
		{"https://api.github.com/repositories/1300192/issues?page=515"},
		{"https://api.github.com/repositories/1300192/issues/?page=515"},
		{"k1LoW/urlfilepath"},
		{"/k1LoW/urlfilepath"},
		{"repositories/1300192/issues?page=515"},
		{"s3://testbucket/path/to/urlfilepath"},
		{"file:///root/path/to/urlfilepath"},
	}
	for _, tt := range tests {
		t.Run(tt.urlstr, func(t *testing.T) {
			dir := t.TempDir()
			u, err := url.Parse(tt.urlstr)
			if err != nil {
				t.Error(err)
			}
			pathstr, err := urlfilepath.Encode(u)
			if err != nil {
				t.Error(err)
			}
			if err := os.MkdirAll(filepath.Join(dir, pathstr), os.ModePerm); err != nil {
				t.Error(err)
			}
			if err := os.WriteFile(filepath.Join(dir, pathstr, "test.txt"), []byte(tt.urlstr), os.ModePerm); err != nil {
				t.Error(err)
			}
			got, err := urlfilepath.Decode(pathstr)
			if err != nil {
				t.Error(err)
			}
			if got.String() != tt.urlstr {
				t.Errorf("got %v\nwant %v", got, tt.urlstr)
			}
		})
	}
}
