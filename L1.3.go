package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Для запуска используй:
// go run L1.3.go -workers=10
// флаг -workers -- кол-во воркеров
// В данной задаче уже реализован graceful shutdown, необходимый для L1.4
func main() {
	var workersCount int
	flag.IntVar(&workersCount, "workers", 5, "number of workers")
	flag.Parse()

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	dataChannel := make(chan interface{})

	var wg sync.WaitGroup

	go func() {
		Workers(shutdownCtx, workersCount, dataChannel, &wg)
	}()

	go func() {
		defer close(dataChannel)
		for i := 1; ; i++ {
			select {
			case <-shutdownCtx.Done():
				return
			case dataChannel <- fmt.Sprintf("сообщение %d", i):
				time.Sleep(1 * time.Second)
			}
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalCh:
		fmt.Println("получен сигнал:", sig, ", завершаем работу")
		signal.Stop(signalCh) // чтобы не получать сигнал ещё раз
		shutdownCancel()      // вызов прямо в селекте обусловлен тем, что после получения сигнала есть еще логика (wg.Wait())
	case <-shutdownCtx.Done():
		fmt.Println("контекст отменен, завершаем работу")
	}

	wg.Wait()
	fmt.Println("программа завершена")
}

// Workers
// Начиная с версии 1.19 не требуется явное замыкание, посиму и не стал делать это
func Workers(ctx context.Context, workersCount int, ch <-chan interface{}, wg *sync.WaitGroup) {
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if ok != true {
						return
					}
					fmt.Printf("Worker %d: %v\n", i, data)
				}
			}
		}()
	}
}
