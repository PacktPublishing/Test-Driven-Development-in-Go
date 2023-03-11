package contract_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/require"
)

const PACTS_PATH = "./pacts/consumer-bookswap.json"

func TestProviderIndex_Local(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestConsumerIndex_Local in short mode.")
	}
	// Initialise
	pact := dsl.Pact{
		Provider: "BookSwap",
	}
	url := getTestEndpoint(t)

	// Verify
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: url,
		PactURLs:        []string{PACTS_PATH},
	})
	require.Nil(t, err)
}

func getTestEndpoint(t *testing.T) string {
	t.Helper()
	baseURL, ok := os.LookupEnv("BOOKSWAP_BASE_URL")
	require.True(t, ok)
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	require.True(t, ok)

	return fmt.Sprintf("%s:%s", baseURL, port)
}
