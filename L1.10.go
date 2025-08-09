package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	temperatureMap := GroupTemperature(nums)

	fmt.Println(temperatureMap)
}

func GroupTemperature(nums []float64) map[int][]float64 {
	temperatureMap := make(map[int][]float64)
	for i := 0; i < len(nums); i++ {
		key := int(math.Trunc(nums[i])) / 10 * 10
		if _, exist := temperatureMap[key]; exist == false {
			temperatureMap[key] = []float64{nums[i]}
			continue
		}
		temperatureMap[key] = append(temperatureMap[key], nums[i])
	}
	return temperatureMap
}
