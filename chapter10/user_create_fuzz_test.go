package bookswapfuzzhttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzTestUserCreation(f *testing.F) {
	if os.Getenv("LONG") == "" {
		f.Skip("Skipping FuzzTestUserCreation in short mode.")
	}
	userEndpoint := getTestEndpoint(f)
	f.Add("test user", "1 London Road", "N1", "UK")
	f.Fuzz(func(t *testing.T, name string, address string,
		postCode string, country string) {
		requestBody, err := json.Marshal(map[string]string{
			"name":      name,
			"address":   address,
			"post_code": postCode,
			"country":   country,
		})
		require.Nil(t, err)
		req := bytes.NewBuffer(requestBody)
		require.Nil(t, err)
		resp, err := http.Post(userEndpoint, "application/json", req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func getTestEndpoint(f *testing.F) string {
	f.Helper()
	baseURL, ok := os.LookupEnv("BOOKSWAP_BASE_URL")
	require.True(f, ok)
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	require.True(f, ok)

	return fmt.Sprintf("%s:%s/users", baseURL, port)
}
