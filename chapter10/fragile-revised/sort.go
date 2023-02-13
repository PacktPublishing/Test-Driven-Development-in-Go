package fragilerevised

import (
	"fmt"
	"sort"
)

type SortDirection int

const (
	ASC SortDirection = iota
	DESC
)

// GetSortedValues returns the key-sorted values of a given input map.
func GetSortedValues(input map[int]string, dir SortDirection) ([]string, error) {
	if input == nil {
		return nil, fmt.Errorf("cannot sort nil input map")
	}
	keys := make([]int, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	switch dir {
	case ASC:
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
	case DESC:
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
	default:
		return nil, fmt.Errorf("sort direction not recognised")
	}
	vals := make([]string, 0, len(input))
	for _, k := range keys {
		vals = append(vals, input[k])
	}
	return vals, nil
}
