package perf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter04/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const LOAD_AMOUNT = 1000

func TestUpsertUser_Load(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestUpsertUser_Load in short mode.")
	}
	userEndpoint := getTestEndpoint(t)
	requestBody, err := json.Marshal(map[string]string{
		"name":      "Concurrent Test User",
		"address":   "1 London Road",
		"post_code": "N1",
		"country":   "United Kingdom",
	})
	require.Nil(t, err)
	require.NotNil(t, requestBody)
	for i := 0; i < LOAD_AMOUNT; i++ {
		t.Run("concurrent upsert", func(t *testing.T) {
			t.Parallel()
			req := bytes.NewBuffer(requestBody)
			r, err := http.Post(userEndpoint, "application/json", req)
			assert.Nil(t, err)
			body, err := io.ReadAll(r.Body)
			r.Body.Close()
			require.Nil(t, err)

			var resp handlers.Response
			err = json.Unmarshal(body, &resp)
			require.Nil(t, err)
			assert.Equal(t, http.StatusOK, r.StatusCode)
			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.User.ID)
		})
	}
}

func getTestEndpoint(t *testing.T) string {
	t.Helper()
	baseURL, ok := os.LookupEnv("BOOKSWAP_BASE_URL")
	require.True(t, ok)
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	require.True(t, ok)

	return fmt.Sprintf("%s:%s/users", baseURL, port)
}
