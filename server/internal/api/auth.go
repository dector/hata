package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"hata/internal/db"
	"hata/internal/util"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthResult contains successful authentication result.
type AuthResult struct {
	User    *db.UserData
	Session *db.SessionData
}

// AuthHandler provides authentication endpoints.
type AuthHandler struct {
	repos db.Repositories
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(repos db.Repositories) *AuthHandler {
	return &AuthHandler{
		repos: repos,
	}
}

// Authenticate validates credentials and creates a session.
func (h *AuthHandler) Authenticate(ctx context.Context, username, password string) (*AuthResult, error) {
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := h.repos.User().GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if user == nil || !util.VerifyPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := util.GenerateToken(40)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	validUntil := time.Now().Add(30 * 24 * time.Hour)
	session, err := h.repos.Session().Create(ctx, user.ID, token, validUntil)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &AuthResult{User: user, Session: session}, nil
}

// Login handles POST /api/latest/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format", "invalid-request")
		return
	}

	result, err := h.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			WriteError(w, http.StatusBadRequest, "Invalid username or password", "invalid-credentials")
			return
		}
		fmt.Printf("Error authenticating user: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	WriteJSON(w, http.StatusOK, LoginResponse{
		Session: SessionInfo{
			Token:      result.Session.Token,
			ValidUntil: result.Session.ValidUntil.Format(time.RFC3339),
		},
		User: UserInfo{
			DisplayName: result.User.DisplayName,
		},
	})
}
