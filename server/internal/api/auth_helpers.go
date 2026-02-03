package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"hata/internal/db"
)

var errUnauthorized = errors.New("unauthorized")

func userIDFromBearerToken(r *http.Request, repos db.Repositories) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errUnauthorized
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return 0, errUnauthorized
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return 0, errUnauthorized
	}

	session, err := repos.Session().GetByToken(r.Context(), token)
	if err != nil {
		return 0, err
	}
	if session == nil {
		return 0, errUnauthorized
	}
	if session.InvalidSince != nil {
		return 0, errUnauthorized
	}
	if time.Now().After(session.ValidUntil) {
		return 0, errUnauthorized
	}

	return session.UserID, nil
}
