package main

import "fmt"

func main() {
	strs := []string{"cat", "cat", "dog", "cat", "tree"}
	fmt.Println(stringsArray(strs))
}

func stringsArray(strs []string) []string {
	mapa := make(map[string]struct{})
	for i := 0; i < len(strs); i++ {
		mapa[strs[i]] = struct{}{}
	}

	slice := make([]string, 0, len(mapa))
	for key, _ := range mapa {
		slice = append(slice, key)
	}

	return slice
}
