package api

import (
	"net/http"
	"os"
	"strings"
	"time"
)

const AuthTokenCookieName = "hata_auth_token"

func SetAuthTokenCookie(w http.ResponseWriter, token string, validUntil time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   authCookieSecure(),
		Expires:  validUntil,
	})
}

func ClearAuthTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    "-",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   authCookieSecure(),
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func authCookieSecure() bool {
	if v := strings.TrimSpace(os.Getenv("HATA_COOKIE_SECURE")); v != "" {
		return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
	}

	return os.Getenv("HATA_DEV") != "1" && os.Getenv("AIR") != "1"
}
