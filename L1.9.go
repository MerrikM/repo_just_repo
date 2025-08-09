package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	firstChan := make(chan int)
	secondChan := make(chan int)

	nums := []int{1, 2, 3, 4, 5}

	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		Generator(ctx, firstChan, nums)
	}()

	go func() {
		defer wg.Done()
		Handler(ctx, firstChan, secondChan)
	}()

	go func() {
		defer wg.Done()
		Consumer(ctx, secondChan)
	}()

	wg.Wait()
}

func Generator(ctx context.Context, firstChan chan<- int, nums []int) {
	defer close(firstChan)
	for i := 0; i < len(nums); i++ {
		select {
		case <-ctx.Done():
			return
		case firstChan <- nums[i]:
		}
	}
}

func Handler(ctx context.Context, firstChan <-chan int, secondChan chan<- int) {
	defer close(secondChan)
	for value := range firstChan {
		select {
		case <-ctx.Done():
			return
		case secondChan <- value * 2:

		}
	}
}

func Consumer(ctx context.Context, secondChan <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-secondChan:
			if ok != true {
				return
			}
			fmt.Println(value)
		}
	}
}
