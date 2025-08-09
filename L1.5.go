package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	timer := time.Duration(5 * time.Second)
	ch := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(2)
	go producer(ch, timer, &wg)
	go consumer(ch, timer, &wg)

	wg.Wait()
}

func producer(ch chan<- interface{}, t time.Duration, wg *sync.WaitGroup) {
	timeout := time.After(t)
	defer wg.Done()

	for {
		select {
		case <-timeout:
			fmt.Println("время вышло, завершаю отправку данных")
			close(ch)
			return
		default:
			ch <- "message"
		}
	}
}

func consumer(ch <-chan interface{}, t time.Duration, wg *sync.WaitGroup) {
	timeout := time.After(t)
	defer wg.Done()

	for {
		select {
		case <-timeout:
			fmt.Println("время вышло, завершаю чтение данных")
			return
		case data, ok := <-ch:
			if ok != true {
				return
			}
			fmt.Println(data)
		}
	}
}
