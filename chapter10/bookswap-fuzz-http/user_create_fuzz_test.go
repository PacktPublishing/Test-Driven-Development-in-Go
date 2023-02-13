package bookswapfuzzhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const userEndpoint = "http://localhost:3000/users"

func FuzzTestUserCreation(f *testing.F) {
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
		resp, err := http.Post(userEndpoint, "application/json", req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
