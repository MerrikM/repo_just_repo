package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*

Вывести 1-й и 3-й столбцы
go run L2.13.go -f 1,3 -d "," test13.txt
Результат:
a,c
1,3
x,z

Вывести 2-4 столбцы
go run L2.13.go -f 2-4 -d "," test13.txt
Результат:
b,c,d
2,3,4
y,z

Только строки с разделителем
go run L2.13.go -f 1 -d "," -s test13.txt
Результат:
a
1
x

*/

// fieldsSet хранит уникальные номера колонок (0-based)
type fieldsSet map[int]struct{}

// parseFields парсит строку вида "1,3-5,7" и возвращает set колонок
func parseFields(spec string) (fieldsSet, error) {
	set := make(fieldsSet)
	parts := strings.Split(spec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			// диапазон
			bounds := strings.SplitN(part, "-", 2)
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start <= 0 || end <= 0 || start > end {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := start - 1; i <= end-1; i++ { // перевод в 0-based
				set[i] = struct{}{}
			}
		} else {
			// отдельное число
			num, err := strconv.Atoi(part)
			if err != nil || num <= 0 {
				return nil, fmt.Errorf("invalid field: %s", part)
			}
			set[num-1] = struct{}{} // 0-based
		}
	}
	return set, nil
}

// processLine обрабатывает одну строку, возвращает выбранные колонки
func processLine(line string, delim string, onlySeparated bool, fields fieldsSet) (string, bool) {
	if onlySeparated && !strings.Contains(line, delim) {
		return "", false
	}

	cols := strings.Split(line, delim)
	var out []string
	for i := 0; i < len(cols); i++ {
		if _, ok := fields[i]; ok {
			out = append(out, cols[i])
		}
	}
	if len(out) == 0 {
		return "", false
	}
	return strings.Join(out, delim), true
}

func main() {
	fieldsArg := flag.String("f", "", "Fields to extract, e.g. 1,3-5")
	delim := flag.String("d", "\t", "Field delimiter")
	onlySeparated := flag.Bool("s", false, "Only print lines containing delimiter")

	flag.Parse()

	if *fieldsArg == "" {
		fmt.Fprintln(os.Stderr, "error: -f flag is required")
		os.Exit(1)
	}

	fields, err := parseFields(*fieldsArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -f argument: %v\n", err)
		os.Exit(1)
	}

	var file io.Reader = os.Stdin
	args := flag.Args()
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		file = f
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		out, ok := processLine(line, *delim, *onlySeparated, fields)
		if ok {
			fmt.Println(out)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
}
