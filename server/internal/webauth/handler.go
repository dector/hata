package webauth

import (
	"context"
	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/integrations"
	"hata/internal/weather"
)

const defaultThenPath = "/me"

// Handler serves browser auth pages.

type Handler struct {
	auth             *api.AuthHandler
	repos            db.Repositories
	deviceController api.DeviceController
	discoverDevices  func(context.Context, []string) <-chan integrations.DiscoveredDevice
	weather          *weather.Plugin
}

func NewHandler(auth *api.AuthHandler, repos db.Repositories) *Handler {
	return NewHandlerWithWeather(auth, repos, nil)
}

func NewHandlerWithWeather(auth *api.AuthHandler, repos db.Repositories, weatherPlugin *weather.Plugin) *Handler {
	return &Handler{
		auth:             auth,
		repos:            repos,
		deviceController: &api.RealDeviceController{},
		discoverDevices:  integrations.DiscoverDevices,
		weather:          weatherPlugin,
	}
}
