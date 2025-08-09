package main

import "fmt"

func main() {
	slice := []int{2, 3, 4, 1, 2331, -24, -2, 0, 52, 228}
	fmt.Println(removeElementFromSlice(slice, 2))
}

func removeElementFromSlice(slice []int, i int) []int {
	if i < 0 || i >= len(slice) {
		return slice
	}
	copy(slice[i:], slice[i+1:]) // сдвигаем хвост на позицию i
	slice[len(slice)-1] = 0      // обнуляем последний элемент для предотвращения утечки памяти (если это срез ссылочных типов)
	return slice[:len(slice)-1]  // уменьшаем длину слайса на 1
}
