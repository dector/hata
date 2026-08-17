package webauth

import (
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"hata/internal/db"
)

const activeHouseCookieName = "active_house_id"

func activeHouseFromRequest(r *http.Request, memberships []*db.HouseMembershipData) *db.HouseMembershipData {
	if len(memberships) == 0 {
		return nil
	}

	if c, err := r.Cookie(activeHouseCookieName); err == nil {
		if membership := findMembership(memberships, strings.TrimSpace(c.Value)); membership != nil {
			return membership
		}
	}

	return defaultHouseMembership(memberships)
}

func defaultHouseMembership(memberships []*db.HouseMembershipData) *db.HouseMembershipData {
	if len(memberships) == 0 {
		return nil
	}
	sortMemberships(memberships)
	return memberships[0]
}

func sortMemberships(memberships []*db.HouseMembershipData) {
	sort.Slice(memberships, func(i, j int) bool {
		if memberships[i].DisplayName == memberships[j].DisplayName {
			return memberships[i].HouseID < memberships[j].HouseID
		}
		return memberships[i].DisplayName < memberships[j].DisplayName
	})
}

func setActiveHouseCookie(w http.ResponseWriter, houseID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     activeHouseCookieName,
		Value:    houseID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   webCookieSecure(),
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
}

func clearActiveHouseCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     activeHouseCookieName,
		Value:    "-",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   webCookieSecure(),
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func webCookieSecure() bool {
	if v := strings.TrimSpace(os.Getenv("HATA_COOKIE_SECURE")); v != "" {
		return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
	}
	return os.Getenv("HATA_DEV") != "1" && os.Getenv("AIR") != "1"
}
