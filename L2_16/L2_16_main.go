package main

import (
	"flag"
	"fmt"
	"os"
	"repo_just_repo/L2_16/internal/crawler"
	"repo_just_repo/L2_16/internal/downloader"
)

// Пример запуска: go run . -url=https://example.com -depth=1 -parallel=3
func main() {
	startURL := flag.String("url", "", "Начальный URL для скачивания")
	depth := flag.Int("depth", 2, "Максимальная глубина обхода")
	parallel := flag.Int("parallel", 5, "Максимальное количество параллельных загрузок")

	flag.Parse()

	if *startURL == "" {
		fmt.Println("Ошибка: необходимо указать начальный URL с помощью -url")
		os.Exit(1)
	}

	fmt.Println("Запуск с параметрами:")
	fmt.Println("URL:", *startURL)
	fmt.Println("Глубина:", *depth)
	fmt.Println("Параллельность:", *parallel)

	d := downloader.NewDownloader()

	c := crawler.NewCrawler(d, *depth, *parallel)
	fmt.Println("Начинаем обход сайта...")
	c.Start(*startURL)
	fmt.Println("Обход завершён.")
}
