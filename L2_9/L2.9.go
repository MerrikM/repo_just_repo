package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func stringUnpacking(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	builder := strings.Builder{}
	runes := []rune(str)
	escaped := false
	var lastRune rune

	if unicode.IsDigit(runes[0]) {
		return "", errors.New("некорректная строка: начинается с цифры")
	}

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		if escaped {
			builder.WriteRune(ch)
			lastRune = ch
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if unicode.IsDigit(ch) {
			if lastRune == 0 {
				return "", errors.New("некорректная строка: число без символа")
			}

			j := i
			for j < len(runes) && unicode.IsDigit(runes[j]) {
				j++
			}
			num, _ := strconv.Atoi(string(runes[i:j]))
			i = j - 1

			for k := 1; k < num; k++ {
				builder.WriteRune(lastRune)
			}
		} else {
			builder.WriteRune(ch)
			lastRune = ch
		}
	}

	if escaped {
		return "", errors.New("некорректная строка: строка заканчивается на '\\'")
	}

	return builder.String(), nil
}

func main() {
	tests := []string{
		"a4bc2d5e",  // aaaabccddddde
		"abcd",      // abcd
		"45",        // ошибка
		"",          // ""
		"qwe\\4\\5", // qwe45
		"qwe\\45",   // qwe44444
		"a12",       // aaaaaaaaaaaa
	}

	for _, t := range tests {
		res, err := stringUnpacking(t)
		if err != nil {
			fmt.Printf("IN: %-8s → ERROR: %v\n", t, err)
		} else {
			fmt.Printf("IN: %-8s → OUT: %s\n", t, res)
		}
	}
}
