package api

import (
	"net/http"
)

// ServerHandler provides server-level endpoints.
type ServerHandler struct{}

// NewServerHandler creates a new ServerHandler.
func NewServerHandler() *ServerHandler {
	return &ServerHandler{}
}

// Ping handles GET /api/next/ping
func (h *ServerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, PingResponse{
		ServerName: "Hata Server",
		Version:    "1.0.0",
		APIVersion: "latest",
	})
}
