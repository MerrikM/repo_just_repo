package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	keyFlag     = flag.Int("k", 0, "sort via a key (column number, 1-based)")
	numericFlag = flag.Bool("n", false, "compare according to string numerical value")
	reverseFlag = flag.Bool("r", false, "reverse the result of comparisons")
	uniqueFlag  = flag.Bool("u", false, "output only the first of an equal run")
	monthFlag   = flag.Bool("M", false, "compare as months")
	ignoreSpace = flag.Bool("b", false, "ignore leading and trailing blanks")
	checkFlag   = flag.Bool("c", false, "check for sorted input")
	humanFlag   = flag.Bool("h", false, "compare human readable numbers (1K 1M etc.)")
)

var monthMap = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
	"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
	"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

// parseHuman преобразует строки типа 1K, 2M в float64
func parseHuman(s string) float64 {
	s = strings.TrimSpace(s)
	mult := 1.0
	if len(s) == 0 {
		return 0
	}
	last := s[len(s)-1]
	switch last {
	case 'K', 'k':
		mult = 1e3
		s = s[:len(s)-1]
	case 'M', 'm':
		mult = 1e6
		s = s[:len(s)-1]
	case 'G', 'g':
		mult = 1e9
		s = s[:len(s)-1]
	case 'T', 't':
		mult = 1e12
		s = s[:len(s)-1]
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val * mult
}

// getKey возвращает значение столбца на основе флага -k
func getKey(line string) string {
	if *keyFlag == 0 {
		return line
	}
	fields := strings.Fields(line)
	if *keyFlag-1 < len(fields) {
		return fields[*keyFlag-1]
	}
	return ""
}

// cmpLines сравнивает две строки по флагам
func cmpLines(a, b string) bool {
	if *ignoreSpace {
		a = strings.TrimSpace(a)
		b = strings.TrimSpace(b)
	}
	keyA := getKey(a)
	keyB := getKey(b)

	var valA, valB float64
	if *numericFlag {
		valA, _ = strconv.ParseFloat(keyA, 64)
		valB, _ = strconv.ParseFloat(keyB, 64)
	} else if *humanFlag {
		valA = parseHuman(keyA)
		valB = parseHuman(keyB)
	} else if *monthFlag {
		valA = float64(monthMap[keyA])
		valB = float64(monthMap[keyB])
	}

	// сравнение
	var less bool
	if *numericFlag || *humanFlag || *monthFlag {
		less = valA < valB
		if valA == valB {
			less = keyA < keyB // возврат к строке
		}
	} else {
		less = keyA < keyB
	}

	if *reverseFlag {
		return !less
	}
	return less
}

// checkSorted проверяет, отсортированы ли строки
func checkSorted(lines []string) bool {
	for i := 1; i < len(lines); i++ {
		if !cmpLines(lines[i-1], lines[i]) {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()
	args := flag.Args()

	var scanner *bufio.Scanner
	if len(args) > 0 {
		// читаем из файла
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка чтения файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		// читаем из STDIN
		scanner = bufio.NewScanner(os.Stdin)
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	if *checkFlag {
		if checkSorted(lines) {
			fmt.Println("Ввод отсортирован")
			os.Exit(0)
		} else {
			fmt.Println("Ввод не отсортирован")
			os.Exit(1)
		}
	}

	sort.SliceStable(lines, func(i, j int) bool {
		return cmpLines(lines[i], lines[j])
	})

	if *uniqueFlag && len(lines) > 0 {
		uniq := []string{lines[0]}
		for i := 1; i < len(lines); i++ {
			if lines[i] != lines[i-1] {
				uniq = append(uniq, lines[i])
			}
		}
		lines = uniq
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
