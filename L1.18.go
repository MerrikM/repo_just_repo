package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	num := 0

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 1; j <= 1000; j++ {
				concurrencyIncrement(ctx, &mutex, &num)
			}
		}()
	}
	wg.Wait()

	fmt.Println(num)
}

func concurrencyIncrement(ctx context.Context, mutex *sync.Mutex, num *int) {
	select {
	case <-ctx.Done():
		return
	default:
		mutex.Lock()
		*num++
		mutex.Unlock()

		// Либо если без мьюетексов через:
		// atomic.AddInt64(num, 1)
		// но в таком случае num должен быть int64
	}
}
