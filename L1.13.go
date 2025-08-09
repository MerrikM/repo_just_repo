package main

import "fmt"

func main() {
	x := 5
	y := 10

	fmt.Println(ValuesTrade(x, y))
}

func ValuesTrade(var1 int, var2 int) (int, int) {
	var1 = var1 ^ var2
	var2 = var1 ^ var2
	var1 = var1 ^ var2

	return var1, var2
}
