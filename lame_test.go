package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/stretchr/testify/assert"
)

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

func TestAPIKey(t *testing.T) {
	// Prepare request
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rec := httptest.NewRecorder()

	// Call handler
	handlerReadiness(rec, req)

	// Validate response
	assert.Equal(t, http.StatusOK, rec.Code)
	var respBody map[string]string
	json.NewDecoder(rec.Body).Decode(&respBody)
	assert.Equal(t, "odsdasdk", respBody["status"])
}
