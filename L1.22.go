package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := big.NewInt(1 << 21) // 2097152
	b := big.NewInt(1 << 22) // 4194304

	fmt.Println("a =", a)
	fmt.Println("b =", b)
	fmt.Println("a + b =", bigAdd(a, b))
	fmt.Println("b - a =", bigSub(b, a))
	fmt.Println("a * b =", bigMul(a, b))
	fmt.Println("b / a =", bigDiv(b, a))
}

func bigAdd(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func bigSub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func bigMul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func bigDiv(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}
