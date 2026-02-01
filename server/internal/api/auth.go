package api

import (
	"fmt"
	"net/http"
	"time"

	"hata/internal/db"
	"hata/internal/util"
)

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

// Login handles POST /api/latest/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req LoginRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format", "invalid-request")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "Invalid username or password", "invalid-credentials")
		return
	}

	ctx := r.Context()

	// Get user by username
	user, err := h.repos.User().GetByUsername(ctx, req.Username)
	if err != nil {
		// Internal error
		fmt.Printf("Error getting user: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	// User not found OR password doesn't match - return SAME error
	if user == nil || !util.VerifyPassword(req.Password, user.PasswordHash) {
		WriteError(w, http.StatusBadRequest, "Invalid username or password", "invalid-credentials")
		return
	}

	// Generate session token
	token, err := util.GenerateToken(40)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	// Create session (valid for 30 days)
	validUntil := time.Now().Add(30 * 24 * time.Hour)
	session, err := h.repos.Session().Create(ctx, user.ID, token, validUntil)
	if err != nil {
		fmt.Printf("Error creating session: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	// Return success response
	WriteJSON(w, http.StatusOK, LoginResponse{
		Session: SessionInfo{
			Token:      session.Token,
			ValidUntil: session.ValidUntil.Format(time.RFC3339),
		},
		User: UserInfo{
			DisplayName: user.DisplayName,
		},
	})
}
