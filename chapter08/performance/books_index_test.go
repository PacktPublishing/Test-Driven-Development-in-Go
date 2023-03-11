package perf

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkGetIndex(b *testing.B) {
	endpoint := getTestEndpoint(b)
	for x := 0; x < b.N; x++ {
		bks, err := http.Get(endpoint)
		assert.Nil(b, err)
		assert.NotNil(b, bks)
	}
}

func getTestEndpoint(b *testing.B) string {
	b.Helper()
	baseURL, ok := os.LookupEnv("BOOKSWAP_BASE_URL")
	require.True(b, ok)
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	require.True(b, ok)

	return fmt.Sprintf("%s:%s", baseURL, port)
}
