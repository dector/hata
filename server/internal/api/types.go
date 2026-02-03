package api

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a successful login response.
type LoginResponse struct {
	Session SessionInfo `json:"session"`
	User    UserInfo    `json:"user"`
}

// SessionInfo contains session details.
type SessionInfo struct {
	Token      string `json:"token"`
	ValidUntil string `json:"validUntil"` // ISO 8601 format
}

// UserInfo contains user details.
type UserInfo struct {
	DisplayName string `json:"displayName"`
}

// PingResponse represents the ping endpoint response.
type PingResponse struct {
	ServerName string `json:"serverName"`
	Version    string `json:"version"`
	APIVersion string `json:"apiVersion"`
}

// HouseInfo contains house details for a user.
type HouseInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
}

// HouseListResponse represents the house list endpoint response.
type HouseListResponse struct {
	Houses []HouseInfo `json:"houses"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information.
type ErrorDetail struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
