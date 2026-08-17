package webauth

import (
	"context"
	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/integrations"
)

const defaultThenPath = "/me"

// Handler serves browser auth pages.

type Handler struct {
	auth             *api.AuthHandler
	repos            db.Repositories
	deviceController api.DeviceController
	discoverDevices  func(context.Context, []string) <-chan integrations.DiscoveredDevice
}

func NewHandler(auth *api.AuthHandler, repos db.Repositories) *Handler {
	return &Handler{
		auth:             auth,
		repos:            repos,
		deviceController: &api.RealDeviceController{},
		discoverDevices:  integrations.DiscoverDevices,
	}
}
