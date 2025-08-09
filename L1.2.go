package main

import (
	"fmt"
	"sync"
)

func main() {
	L1_2([]int{2, 4, 6, 8, 10})
}

func L1_2(nums []int) {
	var wg sync.WaitGroup

	result := make(chan int, len(nums))

	for i := 0; i < len(nums); i++ {
		wg.Add(1)
		go func(newI int) { // в новых версиях Го явное замыкание уже не нужно
			defer wg.Done()
			result <- nums[newI] * nums[newI]
		}(i) // в новых версиях Го явное замыкание уже не нужно
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for i := range result {
		fmt.Println(i)
	}
}
