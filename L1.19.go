package main

import (
	"fmt"
)

func main() {
	str := "💬глав рыбка✅"
	fmt.Println(stringRevers(str))
}

func stringRevers(str string) string {
	runes := []rune(str)
	left := 0
	right := len(runes) - 1

	for left < right {
		tmp := runes[left]
		runes[left] = runes[right]
		runes[right] = tmp

		left++
		right--
	}

	return string(runes)
}
