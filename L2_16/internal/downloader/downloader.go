package downloader

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"repo_just_repo/L2_16/internal/parser"
	"strings"
	"time"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 10 * time.Second, // таймаут на запрос
		},
	}
}

func (d *Downloader) Download(rawURL string) ([]string, error) {
	resp, err := d.client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("запрос не удался: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка загрузки %s: %s", rawURL, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	path := localPath(rawURL)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if strings.Contains(contentType, "text/html") {
		// парсим HTML
		doc, err := html.Parse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
		}

		// собираем и переписываем ссылки
		links := parser.RewriteLinks(rawURL, doc)

		// сохраняем переписанный HTML
		if err := html.Render(f, doc); err != nil {
			return nil, err
		}

		fmt.Println("Скачано HTML:", rawURL, "→", path)
		return links, nil
	}

	// для других файлов просто сохраняем
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Скачан файл:", rawURL, "→", path)
	return nil, nil
}

func localPath(rawURL string) string {
	u, _ := url.Parse(rawURL)
	path := u.Host + u.Path

	// если путь пустой или заканчивается на /
	if strings.HasSuffix(path, "/") || path == u.Host {
		path = filepath.Join(path, "index.html")
	}

	return path
}
