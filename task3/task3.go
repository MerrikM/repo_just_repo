package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// парсинг диапазона "start:end"
func parseRange(r string) (int, int, error) {
	parts := strings.Split(r, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат диапазона: %s", r)
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("некорректные числа в диапазоне: %s", r)
	}
	if start > end {
		start, end = end, start
	}
	return start, end, nil
}

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func main() {
	fileName := flag.String("file", "output.txt", "имя выходного файла")
	timeout := flag.Int("timeout", 10, "таймаут в секундах")
	var ranges multiFlag
	flag.Var(&ranges, "range", "числовой диапазон в формате start:end (может повторяться)")

	flag.Parse()

	if len(ranges) == 0 {
		fmt.Println("ошибка: укажите хотя бы один диапазон через --range")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	primesCh := make(chan int, 100)
	var wg sync.WaitGroup

	for _, r := range ranges {
		start, end, err := parseRange(r)
		if err != nil {
			fmt.Println("ошибка:", err)
			return
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			for i := s; i <= e; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					if isPrime(i) {
						primesCh <- i
					}
				}
			}
		}(start, end)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		f, err := os.Create(*fileName)
		if err != nil {
			fmt.Println("ошибка при создании файла:", err)
			cancel()
			return
		}
		defer f.Close()

		for prime := range primesCh {
			_, _ = f.WriteString(fmt.Sprintf("%d\n", prime))
		}
	}()

	go func() {
		wg.Wait()
		close(primesCh)
	}()

	select {
	case <-done:
		fmt.Println("завершено успешно, результаты записаны в файл:", *fileName)
	case <-ctx.Done():
		fmt.Println("время ожидания истекло, выполнение остановлено")
	}
}
