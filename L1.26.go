package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(allUniqueChars("abcd"))
	fmt.Println(allUniqueChars("abCdefAaf"))
	fmt.Println(allUniqueChars("aabcd"))
}

func allUniqueChars(str string) bool {
	str = strings.ToLower(str)

	seen := make(map[rune]bool)
	for _, ch := range str {
		if seen[ch] {
			return false
		}
		seen[ch] = true
	}
	return true
}
