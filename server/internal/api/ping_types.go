package api

// PingResponse represents the ping endpoint response.
type PingResponse struct {
	ServerName string `json:"serverName"`
	Version    string `json:"version"`
	APIVersion string `json:"apiVersion"`
}
