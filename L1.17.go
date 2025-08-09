package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{2, 3, 4, 1, 2331, -24, -2, 0, 52, 228}

	sort.Ints(nums)

	fmt.Println("Отсортированный срез:", nums)
	fmt.Println("Индекс числа:", binarySearch(nums, 0))
}

func binarySearch(nums []int, target int) int {
	left := 0
	right := len(nums) - 1

	for left <= right {
		middle := left + (right-left)/2

		if nums[middle] == target {
			return middle
		} else if nums[middle] < target {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}

	return -1
}
