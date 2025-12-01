package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	dict := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	anagrams := findAnagrams(dict)

	fmt.Println(anagrams)
}

func findAnagrams(words []string) map[string][]string {
	anagramMap := make(map[string][]string) // ключ -- отсортированные буквы, значение -- список слов

	for _, word := range words {
		lower := strings.ToLower(word)
		letters := strings.Split(lower, "")
		sort.Strings(letters)
		key := strings.Join(letters, "")
		anagramMap[key] = append(anagramMap[key], lower)
	}

	result := make(map[string][]string)

	for _, group := range anagramMap {
		if len(group) > 1 {
			sort.Strings(group)
			result[group[0]] = group
		}
	}

	return result
}
