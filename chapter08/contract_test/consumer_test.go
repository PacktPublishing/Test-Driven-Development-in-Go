package contract_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/handlers"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerIndex_Local(t *testing.T) {
	// Initialize
	pact := dsl.Pact{
		Consumer: "Consumer",
		Provider: "BookSwap",
	}
	pact.Setup(true)

	// Test case - makes the call to the provider
	var test = func() (err error) {
		url := fmt.Sprintf("http://localhost:%d/", pact.Server.Port)
		req, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		return
	}

	t.Run("get index", func(t *testing.T) {
		pact.
			AddInteraction().
			Given("BookSwap is up").
			UponReceiving("GET / request").
			WithRequest(dsl.Request {
				Method: "GET",
				Path:   dsl.String("/"),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: http.StatusOK,
				Body: dsl.Like(handlers.Response{
					Message: "Welcome to the BookSwap Service!",
				}),
			})
		require.Nil(t, pact.Verify(test))
	})

	// Clean up
	require.Nil(t, pact.WritePact())
	pact.Teardown()
}
