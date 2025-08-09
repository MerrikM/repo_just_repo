package main

import "fmt"

func main() {
	A := []int{1, 2, 3}
	B := []int{2, 3, 4}

	fmt.Println(intersection(A, B))
}

func intersection(nums1 []int, nums2 []int) []int {
	hashMap := make(map[int]bool)

	for _, value := range nums1 {
		if hashMap[value] == false {
			hashMap[value] = true
		}
	}

	slice := make([]int, 0)
	for _, value := range nums2 {
		if hashMap[value] == true {
			slice = append(slice, value)
			hashMap[value] = false
		}
	}

	return slice
}
