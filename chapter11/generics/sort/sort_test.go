package genSort_test

import (
	"sort"
	"testing"

	gs "github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/generics/sort"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase[K ~int, V comparable] struct {
	input map[K]V
}

type CustomI int

func TestGetSortedValues(t *testing.T) {
	t.Run("[int]string", func(t *testing.T) {
		testStrings := map[string]testCase[int, string]{
			"unordered":       {input: map[int]string{99: "A", 50: "X"}},
			"empty map":       {input: map[int]string{}},
			"negative values": {input: map[int]string{-99: "A", -1: "X"}},
		}
		runTests(t, testStrings)
	})
	t.Run("[CustomI]float64", func(t *testing.T) {
		testFloats := map[string]testCase[CustomI, float64]{
			"unordered":     {input: map[CustomI]float64{99: 1.23, 0: 4.6}},
			"empty map":     {input: map[CustomI]float64{}},
			"negative keys": {input: map[CustomI]float64{-99: 1.23, 0: 4.6}},
		}
		runTests(t, testFloats)
	})
}

func runTests[K ~int, V comparable](t *testing.T, tests map[string]testCase[K, V]) {
	t.Helper()
	for name, rtc := range tests {
		tc := rtc
		t.Run(name, func(t *testing.T) {
			sortedValues, err := gs.GetSortedValues(tc.input, gs.ASC)
			require.Nil(t, err)
			require.NotNil(t, sortedValues)
			AssertMapOrderedByKeys(t, tc.input, sortedValues)
		})
	}
}

func AssertMapOrderedByKeys[K ~int, V comparable](t *testing.T, input map[K]V, want []V) {
	t.Helper()
	keys := make([]K, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for index, v := range want {
		key := keys[index]
		assert.Equal(t, input[key], v)
	}
}
