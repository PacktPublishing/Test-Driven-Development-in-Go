package contract_test

import (
	//
	// some imports
	//
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/require"
)

const URL = "http://localhost:3000"
const PACTS_PATH = "./pacts/consumer-bookswap.json"

func TestProviderIndex_Local(t *testing.T) {
	// Initialise
	pact := dsl.Pact{
		Provider: "BookSwap",
	}

	// Verify
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: URL,
		PactURLs:        []string{PACTS_PATH},
	})
	require.Nil(t, err)
}
