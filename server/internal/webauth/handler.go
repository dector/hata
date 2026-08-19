package webauth

import (
	"context"
	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/extension"
	"hata/internal/integrations"
)

const defaultThenPath = "/me"

// Handler serves browser auth pages.

type Handler struct {
	auth             *api.AuthHandler
	repos            db.Repositories
	deviceController api.DeviceController
	discoverDevices  func(context.Context, []string) <-chan integrations.DiscoveredDevice
	extensions       *extension.Registry
}

func NewHandler(auth *api.AuthHandler, repos db.Repositories) *Handler {
	return NewHandlerWithExtensions(auth, repos, nil)
}

func NewHandlerWithExtensions(auth *api.AuthHandler, repos db.Repositories, extensions *extension.Registry) *Handler {
	if extensions == nil {
		extensions = extension.NewRegistry()
	}
	return &Handler{
		auth:             auth,
		repos:            repos,
		deviceController: &api.RealDeviceController{},
		discoverDevices:  integrations.DiscoverDevices,
		extensions:       extensions,
	}
}
