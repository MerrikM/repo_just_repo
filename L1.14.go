package main

import (
	"fmt"
	"reflect"
)

func main() {
	cha := make(chan interface{})
	VarTypeDetector(cha)
}

// VarTypeDetector
// Обычный кейс для канала не подошел, т.к. придется делать много кейсов для конкретных типов канала,
// поэтому можно использовать рефлексию
func VarTypeDetector(variable interface{}) {
	switch variable.(type) {
	case int:
		fmt.Println("это число")
	case string:
		fmt.Println("это строка")
	case bool:
		fmt.Println("это логический тип")
	default:
		t := reflect.TypeOf(variable)
		if t != nil && t.Kind() == reflect.Chan {
			fmt.Println("это канал")
		} else {
			fmt.Println("для других типов не предусмотренно проверок")
		}
	}
}
