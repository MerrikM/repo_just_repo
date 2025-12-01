package crawler

import (
	"fmt"
	"repo_just_repo/L2_16/internal/downloader"
	"sync"
)

type Crawler struct {
	d        *downloader.Downloader
	maxDepth int
	parallel int

	visited map[string]struct{}
	mu      sync.Mutex
	sem     chan struct{} // для ограничения параллельности
	wg      sync.WaitGroup
}

func NewCrawler(d *downloader.Downloader, maxDepth, parallel int) *Crawler {
	return &Crawler{
		d:        d,
		maxDepth: maxDepth,
		parallel: parallel,
		visited:  make(map[string]struct{}),
		sem:      make(chan struct{}, parallel),
	}
}

func (c *Crawler) Start(rootURL string) {
	c.enqueue(rootURL, 0)
	c.wg.Wait()
}

// Добавление URL в очередь обхода
func (c *Crawler) enqueue(url string, depth int) {
	if depth > c.maxDepth {
		return
	}

	c.mu.Lock()
	if _, ok := c.visited[url]; ok {
		c.mu.Unlock()
		return
	}
	c.visited[url] = struct{}{}
	c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.sem <- struct{}{}
		defer func() { <-c.sem }()

		links, err := c.d.Download(url)
		if err != nil {
			fmt.Println("Ошибка скачивания:", err)
			return
		}

		for _, link := range links {
			c.enqueue(link, depth+1)
		}
	}()
}
