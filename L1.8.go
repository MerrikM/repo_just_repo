package main

import "fmt"

func main() {
	var num int64 = 5 // 0101 (2)

	result1 := SetBit(num, 0, 0)
	fmt.Printf("%b (%d) → %b (%d)\n", num, num, result1, result1)

	result2 := SetBit(num, 1, 1)
	fmt.Printf("%b (%d) → %b (%d)\n", num, num, result2, result2)
}

func SetBit(n int64, i uint, bit uint) int64 {
	if bit == 1 {
		return n | (1 << i) // Установка бита в 1 (OR с маской)
	} else {
		return n &^ (1 << i) // Установка бита в 0 (AND NOT с маской)
	}
}
