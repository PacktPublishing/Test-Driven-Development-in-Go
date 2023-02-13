package fragilerevised_test

import (
	"sort"
	"testing"
	"testing/quick"

	fr "github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter10/fragile-revised"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzGetSortedValues_ASC(f *testing.F) {
	input := map[int]string{
		99: "B",
		0:  "A",
	}
	f.Add(3, "W")

	f.Fuzz(func(t *testing.T, k int, v string) {
		input[k] = v
		keys := make([]int, 0, len(input))
		for k := range input {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		sortedValues, err := fr.GetSortedValues(input, fr.ASC)
		require.Nil(t, err)
		require.NotNil(t, sortedValues)
		for index, v := range sortedValues {
			key := keys[index]
			assert.Equal(t, input[key], v)
		}
	})
}

func TestGetSortedValues_ASC(t *testing.T) {
	input := map[int]string{
		99: "B",
		0:  "A",
	}
	isSorted := func(k int, val string) bool {
		input[k] = val
		keys := make([]int, 0, len(input))
		for k := range input {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		sortedValues, err := fr.GetSortedValues(input, fr.ASC)
		if err != nil || sortedValues == nil {
			return false
		}
		for index, v := range sortedValues {
			key := keys[index]
			if input[key] != v {
				return false
			}
		}
		return true
	}
	if err := quick.Check(isSorted, nil); err != nil {
		t.Error(err)
	}
}
