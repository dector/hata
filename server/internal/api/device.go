package api

import "hata/internal/db"

// DeviceHandler provides device endpoints.
type DeviceHandler struct {
	repos      db.Repositories
	controller DeviceController
}

// NewDeviceHandler creates a new DeviceHandler.
func NewDeviceHandler(repos db.Repositories) *DeviceHandler {
	return &DeviceHandler{
		repos:      repos,
		controller: NewRealDeviceController(),
	}
}

// NewDeviceHandlerWithController creates a new DeviceHandler with explicit device controller.
func NewDeviceHandlerWithController(repos db.Repositories, controller DeviceController) *DeviceHandler {
	if controller == nil {
		controller = NewRealDeviceController()
	}
	return &DeviceHandler{repos: repos, controller: controller}
}
