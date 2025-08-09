package main

import "fmt"

func main() {
	input := "snow dog sun"
	result := WordsReverseInPlace(input)
	fmt.Println(result)
}

func RunesReverse(runes []rune, left int, right int) {
	for left < right {
		tmp := runes[left]
		runes[left] = runes[right]
		runes[right] = tmp

		left++
		right--
	}
}

func WordsReverseInPlace(str string) string {
	runes := []rune(str)

	RunesReverse(runes, 0, len(str)-1)

	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || runes[i] == ' ' {
			RunesReverse(runes, start, i-1)
			start = i + 1
		}
	}

	return string(runes)
}
