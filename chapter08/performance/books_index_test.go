package perf

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const endpoint = "http://localhost:3000"

func BenchmarkGetIndex(b *testing.B) {
	for x := 0; x < b.N; x++ {
		bks, err := http.Get(endpoint)
		assert.Nil(b, err)
		assert.NotNil(b, bks)
	}
}
