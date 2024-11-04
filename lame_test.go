package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/stretchr/testify/assert"

	"github.com/bootdotdev/learn-cicd-starter/internal/auth"
)

type t_apiConfig struct {
	apiConfig
	DB MockDBInterface // Use MockDBInterface instead of *database.Queries
}

// Define a minimal interface with only the methods you need
type MockDBInterface interface {
	CreateUser(ctx context.Context, name string) (string, error)
	GetUser(ctx context.Context, apiKey string) (map[string]string, error)
}

// MockDB struct implementing MockDBInterface without using database-specific types
type MockDB struct {
	CreateUserFunc func(ctx context.Context, name string) (string, error)
	GetUserFunc    func(ctx context.Context, apiKey string) (map[string]string, error)
}

func (m *MockDB) CreateUser(ctx context.Context, name string) (string, error) {
	return m.CreateUserFunc(ctx, name)
}

func (m *MockDB) GetUser(ctx context.Context, apiKey string) (map[string]string, error) {
	return m.GetUserFunc(ctx, apiKey)
}

// MockQueries implements the required methods of database.Queries.
type MockQueries struct {
	CreateUserFunc func(database.CreateUserParams) error
	GetUserFunc    func(string) (database.User, error)
}

func (m *MockQueries) CreateUser(params database.CreateUserParams) error {
	return m.CreateUserFunc(params)
}

func (m *MockQueries) GetUser(apiKey string) (database.User, error) {
	return m.GetUserFunc(apiKey)
}

func TestHandlerReadiness(t *testing.T) {
	// Prepare request
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	// Call handler
	handlerReadiness(rec, req)

	// Validate response
	assert.Equal(t, http.StatusOK, rec.Code)
	var respBody map[string]string
	json.NewDecoder(rec.Body).Decode(&respBody)
	assert.Equal(t, "ok", respBody["status"])
}

func TestGenerateRandomSHA256Hash(t *testing.T) {
	hash1, err := generateRandomSHA256Hash()
	assert.NoError(t, err)
	assert.Len(t, hash1, 64) // SHA-256 hash length in hexadecimal

	hash2, err := generateRandomSHA256Hash()
	assert.NoError(t, err)
	assert.Len(t, hash2, 64)

	// The hashes should be different since they are randomly generated
	assert.NotEqual(t, hash1, hash2)
}

func TestGetAPIKey(t *testing.T) {
	headers := http.Header{}

	// Case 1: No Authorization header
	_, err := auth.GetAPIKey(headers)
	assert.Equal(t, auth.ErrNoAuthHeaderIncluded, err)

	// Case 2: Malformed Authorization header
	headers.Set("Authorization", "Bearer some_token")
	_, err = auth.GetAPIKey(headers)
	assert.EqualError(t, err, "malformed authorization header")

	// Case 3: Correct Authorization header
	headers.Set("Authorization", "ApiKey valid_api_key")
	apiKey, err := auth.GetAPIKey(headers)
	assert.NoError(t, err)
	assert.Equal(t, "valid_api_key", apiKey)
}

func TestIntentional_fail(t *testing.T) {
	headers := http.Header{}

	// Case 1: No Authorization header
	_, err := auth.GetAPIKey(headers)
	assert.Equal(t, auth.ErrNoAuthHeaderIncluded, err)

	// Case 2: Malformed Authorization header
	headers.Set("Authorization", "Bearer some_token")
	_, err = auth.GetAPIKey(headers)
	assert.EqualError(t, err, "malformed authorization header")

	// Case 3: Correct Authorization header
	headers.Set("Authorization", "ApiKey valid_api_key")
	apiKey, err := auth.GetAPIKey(headers)
	assert.NoError(t, err)
	assert.Equal(t, "valid_api_key", apiKey)

}
