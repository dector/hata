package webauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"hata/internal/api"
	"hata/internal/db"
)

var errUnauthorized = errors.New("unauthorized")

type authContextKey struct{}

// AuthContext stores auth data for page handlers.
type AuthContext struct {
	UserID      int
	Username    string
	DisplayName string
	Token       string
}

// AuthFromContext returns auth context attached by the auth middleware.
func AuthFromContext(ctx context.Context) (AuthContext, bool) {
	v := ctx.Value(authContextKey{})
	if v == nil {
		return AuthContext{}, false
	}
	auth, ok := v.(AuthContext)
	return auth, ok
}

// TryAuthFromRequest validates cookie auth without redirecting.
// Returns (auth, true, nil) when authenticated, (zero, false, nil) when anonymous,
// and (zero, false, err) on internal errors.
func TryAuthFromRequest(r *http.Request, repos db.Repositories) (AuthContext, bool, error) {
	auth, err := authFromCookie(r, repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			return AuthContext{}, false, nil
		}
		return AuthContext{}, false, err
	}
	return auth, true, nil
}

// RequirePageAuth validates cookie auth and redirects anonymous users to login.
func RequirePageAuth(repos db.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, ok, err := TryAuthFromRequest(r, repos)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if !ok {
				then := sanitizeThen(currentRequestThen(r))
				login := "/auth/login?then=" + url.QueryEscape(then)
				http.Redirect(w, r, login, http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), authContextKey{}, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authFromCookie(r *http.Request, repos db.Repositories) (AuthContext, error) {
	c, err := r.Cookie(api.AuthTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return AuthContext{}, errUnauthorized
		}
		return AuthContext{}, fmt.Errorf("read cookie: %w", err)
	}

	session, err := repos.Session().GetByToken(r.Context(), c.Value)
	if err != nil {
		return AuthContext{}, fmt.Errorf("get session: %w", err)
	}
	if session == nil || session.InvalidSince != nil || time.Now().After(session.ValidUntil) {
		return AuthContext{}, errUnauthorized
	}

	user, err := repos.User().GetByID(r.Context(), session.UserID)
	if err != nil {
		return AuthContext{}, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return AuthContext{}, errUnauthorized
	}

	return AuthContext{
		UserID:      session.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Token:       session.Token,
	}, nil
}

func currentRequestThen(r *http.Request) string {
	if r.URL == nil {
		return defaultThenPath
	}
	return r.URL.RequestURI()
}
