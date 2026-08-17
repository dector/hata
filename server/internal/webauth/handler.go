package webauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/integrations"
	"hata/internal/webui"

	"github.com/go-chi/chi/v5"
)

const defaultThenPath = "/me"

// Handler serves browser auth pages.
type Handler struct {
	auth             *api.AuthHandler
	repos            db.Repositories
	deviceController api.DeviceController
	discoverDevices  func(context.Context, []string) <-chan integrations.DiscoveredDevice
}

type discoveredDeviceView struct {
	Integration string `json:"integration"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	State       string `json:"state"`
	InHouse     bool   `json:"inHouse"`
	AddURL      string `json:"addUrl"`
}

func NewHandler(auth *api.AuthHandler, repos db.Repositories) *Handler {
	return &Handler{
		auth:             auth,
		repos:            repos,
		deviceController: &api.RealDeviceController{},
		discoverDevices:  integrations.DiscoverDevices,
	}
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	renderLogin(w, r, loginPageView{
		Then:  sanitizeThen(r.URL.Query().Get("then")),
		Error: "",
	}, http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	then := sanitizeThen(r.FormValue("then"))

	result, err := h.auth.Authenticate(r.Context(), username, password)
	if err != nil {
		if errors.Is(err, api.ErrInvalidCredentials) {
			renderLogin(w, r, loginPageView{
				Then:  then,
				Error: "Invalid username or password",
			}, http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	api.SetAuthTokenCookie(w, result.Session.Token, result.Session.ValidUntil)
	redirectAfterPost(w, r, then)
}

func (h *Handler) LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.ClearAuthTokenCookie(w)
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.ClearAuthTokenCookie(w)
	redirectAfterPost(w, r, "/auth/login")
}

func (h *Handler) MePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.MePage(webui.MePageData{
		UserID:          auth.UserID,
		Username:        auth.Username,
		DisplayName:     displayNameFromAuth(auth),
		HasSessionToken: auth.Token != "",
	}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render me page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AppPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}

	devices, err := h.repos.Device().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}

	devicesByHouse := make(map[string][]webui.AppDeviceData, len(memberships))
	for _, d := range devices {
		devicesByHouse[d.HouseID] = append(devicesByHouse[d.HouseID], webui.AppDeviceData{
			ID:            d.ID,
			Name:          d.Name,
			IntegrationID: d.IntegrationID,
			State:         d.State,
			ToggleURL:     webui.HouseDeviceTogglePath(d.HouseID, d.ID),
		})
	}

	houses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		houses = append(houses, webui.AppHouseData{
			ID:          m.HouseID,
			DisplayName: m.DisplayName,
			Role:        m.Role,
			Devices:     devicesByHouse[m.HouseID],
		})
	}

	sort.Slice(houses, func(i, j int) bool {
		if houses[i].DisplayName == houses[j].DisplayName {
			return houses[i].ID < houses[j].ID
		}
		return houses[i].DisplayName < houses[j].DisplayName
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.AppPage(webui.AppPageData{DisplayName: displayNameFromAuth(auth), Houses: houses}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render app page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ToggleHouseDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil || h.deviceController == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if houseID == "" || deviceID == "" {
		http.Error(w, "house id and device id are required", http.StatusBadRequest)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}
	if findMembership(memberships, houseID) == nil {
		http.NotFound(w, r)
		return
	}

	device, err := h.repos.Device().GetByHouseAndID(r.Context(), houseID, deviceID)
	if err != nil {
		http.Error(w, "failed to load device", http.StatusInternalServerError)
		return
	}
	if device == nil {
		http.NotFound(w, r)
		return
	}

	newState := "on"
	if strings.EqualFold(strings.TrimSpace(device.State), "on") {
		newState = "off"
	}

	if err := h.deviceController.SetState(r.Context(), device, newState); err != nil {
		fmt.Printf("Error toggling device %q in house %q: %v\n", deviceID, houseID, err)
		http.Error(w, "failed to toggle device", http.StatusBadGateway)
		return
	}
	if err := h.repos.Device().UpdateState(r.Context(), houseID, deviceID, newState); err != nil {
		http.Error(w, "failed to update device", http.StatusInternalServerError)
		return
	}

	redirectAfterPost(w, r, "/app")
}

func (h *Handler) HouseManagePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		http.Error(w, "house id is required", http.StatusBadRequest)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}

	membership := findMembership(memberships, houseID)
	if membership == nil {
		http.NotFound(w, r)
		return
	}

	devices, err := h.repos.Device().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}

	viewDevices := make([]webui.AppDeviceData, 0, len(devices))
	for _, d := range devices {
		viewDevices = append(viewDevices, webui.AppDeviceData{
			ID:            d.ID,
			Name:          d.Name,
			IntegrationID: d.IntegrationID,
			State:         d.State,
		})
	}

	discoveryNetworks, err := h.repos.HouseDiscoveryNetwork().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load discovery networks", http.StatusInternalServerError)
		return
	}
	viewNetworks := make([]webui.HouseDiscoveryNetworkData, 0, len(discoveryNetworks))
	for _, n := range discoveryNetworks {
		viewNetworks = append(viewNetworks, webui.HouseDiscoveryNetworkData{ID: n.ID, CIDR: n.CIDR, Label: n.Label})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.HouseManagePage(webui.HouseManagePageData{
		DisplayName:       displayNameFromAuth(auth),
		DiscoveryNetworks: viewNetworks,
		CanManageHouse:    canManageHouse(membership.Role),
		House: webui.AppHouseData{
			ID:          membership.HouseID,
			DisplayName: membership.DisplayName,
			Role:        membership.Role,
			Devices:     viewDevices,
		},
	}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render house manage page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HouseDeviceDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil || h.discoverDevices == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		http.Error(w, "house id is required", http.StatusBadRequest)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}
	if findMembership(memberships, houseID) == nil {
		http.NotFound(w, r)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	existingDevices, err := h.repos.Device().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}
	existingIPs := deviceIPsByIntegration(existingDevices)

	discoveryNetworks, err := h.repos.HouseDiscoveryNetwork().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load discovery networks", http.StatusInternalServerError)
		return
	}
	cidrs := make([]string, 0, len(discoveryNetworks))
	for _, network := range discoveryNetworks {
		cidrs = append(cidrs, network.CIDR)
	}

	writeSSE(w, "progress", discoveryProgressMessage(discoveryNetworks))
	flusher.Flush()

	found := 0
	for device := range h.discoverDevices(r.Context(), cidrs) {
		found++
		view := discoveredDeviceView{
			Integration: device.Integration,
			Name:        device.Name,
			IP:          device.IP,
			State:       device.State,
			InHouse:     existingIPs[strings.ToLower(device.Integration)+":"+device.IP],
			AddURL:      houseManageURL(houseID) + "/devices",
		}
		payload, err := json.Marshal(view)
		if err != nil {
			continue
		}
		writeSSE(w, "device", string(payload))
		flusher.Flush()
	}

	if found == 0 {
		writeSSE(w, "done", "Discovery complete. No devices found.")
	} else {
		writeSSE(w, "done", fmt.Sprintf("Discovery complete. Found %d device(s).", found))
	}
	flusher.Flush()
}

func (h *Handler) AddHouseDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	if !canManageHouse(membership.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	integration := strings.ToLower(strings.TrimSpace(r.FormValue("integration")))
	if integration != "wiz" {
		http.Error(w, "unsupported integration", http.StatusBadRequest)
		return
	}
	ip := strings.TrimSpace(r.FormValue("ip"))
	if net.ParseIP(ip) == nil {
		http.Error(w, "valid device IP is required", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = "WiZ " + ip
	}
	state := strings.ToLower(strings.TrimSpace(r.FormValue("state")))
	if state != "on" && state != "off" {
		state = ""
	}

	integrationData := fmt.Sprintf(`{"ip":%q}`, ip)
	if _, err := h.repos.Device().Create(r.Context(), membership.HouseID, deviceIDForIntegrationIP(integration, ip), name, integration, &integrationData, state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

func (h *Handler) AddHouseDiscoveryNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	if !canManageHouse(membership.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if _, err := h.repos.HouseDiscoveryNetwork().Create(r.Context(), membership.HouseID, r.FormValue("cidr"), r.FormValue("label")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

func (h *Handler) DeleteHouseDiscoveryNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	if !canManageHouse(membership.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	networkID, err := strconv.Atoi(strings.TrimSpace(chi.URLParam(r, "networkId")))
	if err != nil || networkID <= 0 {
		http.Error(w, "invalid discovery network id", http.StatusBadRequest)
		return
	}
	if err := h.repos.HouseDiscoveryNetwork().DeleteByID(r.Context(), membership.HouseID, networkID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

func (h *Handler) authorizedHouseMembership(w http.ResponseWriter, r *http.Request) *db.HouseMembershipData {
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil
	}
	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil
	}
	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		http.Error(w, "house id is required", http.StatusBadRequest)
		return nil
	}
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return nil
	}
	membership := findMembership(memberships, houseID)
	if membership == nil {
		http.NotFound(w, r)
		return nil
	}
	return membership
}

func findMembership(memberships []*db.HouseMembershipData, houseID string) *db.HouseMembershipData {
	for _, membership := range memberships {
		if membership.HouseID == houseID {
			return membership
		}
	}
	return nil
}

func writeSSE(w http.ResponseWriter, event, data string) {
	fmt.Fprintf(w, "event: %s\n", event)
	for line := range strings.SplitSeq(data, "\n") {
		fmt.Fprintf(w, "data: %s\n", line)
	}
	fmt.Fprint(w, "\n")
}

type loginPageView struct {
	Then  string
	Error string
}

func deviceIPsByIntegration(devices []*db.DeviceData) map[string]bool {
	result := map[string]bool{}
	for _, device := range devices {
		if device.IntegrationData == nil {
			continue
		}
		var data map[string]any
		if err := json.Unmarshal([]byte(*device.IntegrationData), &data); err != nil {
			continue
		}
		ip, _ := data["ip"].(string)
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		result[strings.ToLower(device.IntegrationID)+":"+ip] = true
	}
	return result
}

func deviceIDForIntegrationIP(integration, ip string) string {
	id := strings.ToLower(strings.TrimSpace(integration)) + "-" + strings.TrimSpace(ip)
	id = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		default:
			return '-'
		}
	}, id)
	id = strings.Trim(id, "-")
	if id == "" {
		return "device"
	}
	return id
}

func canManageHouse(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}

func houseManageURL(houseID string) string {
	return "/h/" + url.PathEscape(houseID) + "/manage"
}

func discoveryProgressMessage(networks []*db.HouseDiscoveryNetworkData) string {
	parts := []string{"Scanning WiZ devices on server local networks"}
	for _, network := range networks {
		label := strings.TrimSpace(network.Label)
		if label == "" {
			parts = append(parts, network.CIDR)
			continue
		}
		parts = append(parts, label+": "+network.CIDR)
	}
	return strings.Join(parts, "; ") + "…"
}

func displayNameFromAuth(auth AuthContext) string {
	displayName := strings.TrimSpace(auth.DisplayName)
	if displayName == "" {
		displayName = auth.Username
	}
	return displayName
}

func renderLogin(w http.ResponseWriter, r *http.Request, view loginPageView, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := webui.LoginPage(webui.LoginPageData{Then: view.Then, Error: view.Error}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render login page", http.StatusInternalServerError)
		return
	}
}

func redirectAfterPost(w http.ResponseWriter, r *http.Request, target string) {
	target = sanitizeThen(target)
	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		w.Header().Set("HX-Redirect", target)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func sanitizeThen(then string) string {
	v := strings.TrimSpace(then)
	if v == "" {
		return defaultThenPath
	}
	if strings.ContainsAny(v, "\\\r\n") {
		return defaultThenPath
	}
	if !strings.HasPrefix(v, "/") || strings.HasPrefix(v, "//") {
		return defaultThenPath
	}

	u, err := url.ParseRequestURI(v)
	if err != nil || u == nil {
		return defaultThenPath
	}
	if u.IsAbs() || u.Host != "" || u.Scheme != "" {
		return defaultThenPath
	}
	if !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return defaultThenPath
	}
	if _, err := url.ParseQuery(u.RawQuery); err != nil {
		return defaultThenPath
	}
	return u.RequestURI()
}
