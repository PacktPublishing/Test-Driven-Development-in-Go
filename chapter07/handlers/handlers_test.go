package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/handlers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestIndexIntegration in short mode.")
	}
	// Arrange
	eb := db.Book{
		ID:     uuid.New().String(),
		Name:   "My first integration test",
		Status: db.Available.String(),
	}
	bs := db.NewBookService([]db.Book{eb}, nil)
	ha := handlers.NewHandler(bs, nil)
	svr := httptest.NewServer(http.HandlerFunc(ha.Index))
	defer svr.Close()

	// Act
	r, err := http.Get(svr.URL)

	// Assert
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)

	assert.Equal(t, 1, len(resp.Books))
	assert.Contains(t, resp.Books, eb)
}
