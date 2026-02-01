package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hata/internal/db"
	"hata/internal/util"
)

// Test database setup helpers

// setupAuthTest initializes a test database with all migrations.
// Returns Repositories interface for the test and a cleanup function.
func setupAuthTest(t *testing.T) (db.Repositories, func()) {
	ctx := context.Background()

	// Open in-memory test database
	dbInst, err := db.OpenTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	cleanup := func() {
		if err := dbInst.Close(); err != nil {
			t.Logf("Warning: failed to close test database: %v", err)
		}
	}

	return dbInst.Repos(), cleanup
}

func TestLogin_Success(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user with hashed password (password: "testpass")
	passwordHash, err := util.HashPassword("testpass")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	_, err = repos.User().Create(ctx, "testuser", passwordHash, "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	handler := NewAuthHandler(repos)

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Session.Token == "" {
		t.Error("Expected session token, got empty string")
	}

	if resp.User.DisplayName != "Test User" {
		t.Errorf("Expected displayName 'Test User', got '%s'", resp.User.DisplayName)
	}

	if resp.Session.ValidUntil == "" {
		t.Error("Expected validUntil timestamp, got empty string")
	}
}

func TestLogin_InvalidCredentials_UserNotFound(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewAuthHandler(repos)

	reqBody := LoginRequest{
		Username: "nonexistent",
		Password: "anypass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Code != "invalid-credentials" {
		t.Errorf("Expected error code 'invalid-credentials', got '%s'", resp.Error.Code)
	}

	if resp.Error.Message != "Invalid username or password" {
		t.Errorf("Expected generic error message, got '%s'", resp.Error.Message)
	}
}

func TestLogin_InvalidCredentials_WrongPassword(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user
	passwordHash, err := util.HashPassword("testpass")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	_, err = repos.User().Create(ctx, "testuser", passwordHash, "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	handler := NewAuthHandler(repos)

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Code != "invalid-credentials" {
		t.Errorf("Expected error code 'invalid-credentials', got '%s'", resp.Error.Code)
	}

	// Generic error message should be the same as user not found
	if resp.Error.Message != "Invalid username or password" {
		t.Errorf("Expected generic error message, got '%s'", resp.Error.Message)
	}
}

func TestLogin_MissingUsername(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewAuthHandler(repos)

	reqBody := LoginRequest{
		Username: "",
		Password: "testpass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Code != "invalid-credentials" {
		t.Errorf("Expected error code 'invalid-credentials', got '%s'", resp.Error.Code)
	}
}

func TestLogin_MissingPassword(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewAuthHandler(repos)

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Code != "invalid-credentials" {
		t.Errorf("Expected error code 'invalid-credentials', got '%s'", resp.Error.Code)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewAuthHandler(repos)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error.Code != "invalid-request" {
		t.Errorf("Expected error code 'invalid-request', got '%s'", resp.Error.Code)
	}
}

func TestLogin_SessionCreated(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test user
	passwordHash, err := util.HashPassword("testpass")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	_, err = repos.User().Create(ctx, "testuser", passwordHash, "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	handler := NewAuthHandler(repos)

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify session token is 40 characters
	var resp LoginResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp.Session.Token) != 40 {
		t.Errorf("Expected token length 40, got %d", len(resp.Session.Token))
	}

	// Verify session was created in the database by querying it
	session, err := repos.Session().GetByToken(ctx, resp.Session.Token)
	if err != nil {
		t.Fatalf("Failed to verify session: %v", err)
	}
	if session == nil {
		t.Error("Session was not created in database")
	}
}
