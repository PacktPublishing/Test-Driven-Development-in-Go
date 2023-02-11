package perf

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const userEndpoint = "http://localhost:3000/users"

func BenchmarkUpsertUser(b *testing.B) {
	requestBody, err := json.Marshal(map[string]string{
		"name":      "Concurrent Test User",
		"address":   "1 London Road",
		"post_code": "N1",
		"country":   "United Kingdom",
	})
	require.Nil(b, err)
	require.NotNil(b, requestBody)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req := bytes.NewBuffer(requestBody)
			resp, err := http.Post(userEndpoint, "application/json", req)
			assert.Nil(b, err)
			defer resp.Body.Close()
			assert.Equal(b, http.StatusOK, resp.StatusCode)
			assert.Nil(b, err)
			assert.NotNil(b, resp)
		}
	})
}
