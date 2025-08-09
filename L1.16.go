package main

import "fmt"

func main() {
	fmt.Println(quickSort([]int{2, 3, 4, 1, 2331, -24, -2, 0, 52, 228}))
}

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	pivotIndex := len(arr) / 2
	pivot := arr[pivotIndex]

	left := make([]int, 0, len(arr))
	right := make([]int, 0, len(arr))
	middle := make([]int, 0, len(arr)) // для элементов равных pivot

	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v > pivot {
			right = append(right, v)
		} else {
			middle = append(middle, v)
		}
	}

	leftSorted := quickSort(left)
	rightSorted := quickSort(right)

	return append(append(leftSorted, middle...), rightSorted...)
}
