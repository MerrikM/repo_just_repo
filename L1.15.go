package main

/*
var justString string

func someFunc() {
  v := createHugeString(1 &lt;&lt; 10)
  justString = v[:100]
}

func main() {
  someFunc()
}

justString это срез, который ссылается на первые 100 байтов переменной v. Т.к. переменная v является локальной переменной,
то после того, как произойдет выход из функции, она уничтожится, но сама строка останется в памяти, т.к. на ее часть продолжает ссылаться
переменная justString

Из-за того, что justString является срезом, в памяти сохраняется вся строка, хотя по задаче нам нужны только первые 100 байтов,
что приводит к лишнему потреблению памяти.

Чтобы исправить это можно в переменную justString копировать первые 100 байт
*/

var justString string

// Первый способ
func someFunc() {
	v := createHugeString(1 << 20)
	tmp := v[:100]

	b := make([]byte, len(tmp))
	copy(b, tmp)

	justString = string(b)
}

// Второй способ
func someFunc2() {
	v := createHugeString(1 << 20)
	justString = string([]byte(v[:100]))
}
