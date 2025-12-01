package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Пример запуска 1: "go line1`nline2`nmatch`nline4`nline5" | go run L2.12.go -C 1 "match"
// Пример запуска 2: go run L2.12.go -C 1 match test12.txt
func main() {
	// Флаги
	after := flag.Int("A", 0, "Print N lines after match")
	before := flag.Int("B", 0, "Print N lines before match")
	context := flag.Int("C", 0, "Print N lines before and after match")
	count := flag.Bool("c", false, "Print only count of matching lines")
	ignoreCase := flag.Bool("i", false, "Ignore case")
	invert := flag.Bool("v", false, "Invert match")
	fixed := flag.Bool("F", false, "Match fixed string instead of regex")
	number := flag.Bool("n", false, "Print line numbers")

	flag.Parse()

	if *context > 0 {
		*before, *after = *context, *context
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: grep [OPTIONS] PATTERN [FILE]")
		os.Exit(1)
	}

	pattern := args[0]
	var file io.Reader = os.Stdin
	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		file = f
	}

	var matcher func(string) bool
	if *fixed {
		if *ignoreCase {
			pattern = strings.ToLower(pattern)
			matcher = func(s string) bool {
				return strings.Contains(strings.ToLower(s), pattern)
			}
		} else {
			matcher = func(s string) bool {
				return strings.Contains(s, pattern)
			}
		}
	} else {
		flags := ""
		if *ignoreCase {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regex: %v\n", err)
			os.Exit(1)
		}
		matcher = func(s string) bool {
			return re.MatchString(s)
		}
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading: %v\n", err)
		os.Exit(1)
	}

	matchedCount := 0
	seen := make([]bool, len(lines)) // для отметки строк, которые должны выводиться

	for i, line := range lines {
		matched := matcher(line)
		if *invert {
			matched = !matched
		}
		if matched {
			matchedCount++
			start := i - *before
			if start < 0 {
				start = 0
			}
			end := i + *after
			if end >= len(lines) {
				end = len(lines) - 1
			}
			for j := start; j <= end; j++ {
				seen[j] = true
			}
		}
	}

	if *count {
		fmt.Println(matchedCount)
		return
	}

	for i, line := range lines {
		if seen[i] {
			if *number {
				fmt.Printf("%d:%s\n", i+1, line)
			} else {
				fmt.Println(line)
			}
		}
	}
}
