package main

import (
	"fmt"
	"sort"
)
var input map[int]string
func GetValues(dir string) []string {
	var keys []int
	for k := range input {
		keys = append(keys, k)
	}
	if dir == "asc" {
		sort.Ints(keys)
	}
	if dir == "desc" {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
	}
	var vals []string
	for _, k := range keys {
		vals = append(vals, input[k])
	}
	return vals
}

func main() {
	input = map[int]string{2:"B", 4: "D", 3: "C", 1:"A"}
	fmt.Println("Sorted asc:", GetValues("asc"))
	fmt.Println("Sorted desc:", GetValues("desc"))
}
