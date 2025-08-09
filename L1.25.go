package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		writer(ch)
	}()

	for value := range ch {
		fmt.Println(value)
	}

	wg.Wait()
}

func mySleep(duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C
	timer.Stop()
}

func writer(ch chan<- int) {
	i := 0
	defer close(ch)
	for {
		i++
		ch <- i
		mySleep(1 * time.Second)
	}
}
