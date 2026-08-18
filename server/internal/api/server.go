package api

import (
	"net/http"

	"hata/version"
)

// ServerHandler provides server-level endpoints.
type ServerHandler struct{}

// NewServerHandler creates a new ServerHandler.
func NewServerHandler() *ServerHandler {
	return &ServerHandler{}
}

// Ping handles GET /api/latest/ping
func (h *ServerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, PingResponse{
		ServerName: "Hata Server",
		Version:    version.Get(),
		APIVersion: "latest",
	})
}
