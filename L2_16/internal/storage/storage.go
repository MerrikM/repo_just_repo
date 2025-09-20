package storage

import (
	"net/url"
	"path/filepath"
	"strings"
)

// LocalPath формирует локальный путь для сохранения URL
func LocalPath(rawURL string) string {
	u, _ := url.Parse(rawURL)
	path := u.Host + u.Path

	if strings.HasSuffix(path, "/") || path == u.Host {
		path = filepath.Join(path, "index.html")
	}
	return path
}
