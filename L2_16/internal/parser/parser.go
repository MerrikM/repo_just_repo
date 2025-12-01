package parser

import (
	"net/url"
	"path/filepath"
	"repo_just_repo/L2_16/internal/storage"

	"golang.org/x/net/html"
)

// RewriteLinks переписывает ссылки в HTML --> локальные пути и возвращает список URL для загрузки
func RewriteLinks(baseURL string, doc *html.Node) []string {
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					abs := resolveURL(baseURL, attr.Val)
					if abs == "" {
						continue
					}
					links = append(links, abs)

					// заменяем ссылку в HTML на локальный путь
					n.Attr[i].Val = filepath.ToSlash(storage.LocalPath(abs))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return unique(links)
}

// Делает относительные ссылки абсолютными
func resolveURL(base, ref string) string {
	b, err := url.Parse(base)
	if err != nil {
		return ""
	}
	r, err := url.Parse(ref)
	if err != nil {
		return ""
	}
	u := b.ResolveReference(r)

	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	if u.Host != b.Host {
		return "" // только свой домен
	}

	return u.String()
}

func unique(strs []string) []string {
	seen := make(map[string]struct{})
	var res []string
	for _, s := range strs {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			res = append(res, s)
		}
	}
	return res
}
