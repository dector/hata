package extension

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HouseActionResult is the result of an extension action for a house.
type HouseActionResult struct {
	Status int
	Data   any
}

// HouseActionHandler handles API extension actions scoped to a house.
type HouseActionHandler interface {
	ExtensionID() string
	HandleHouseAction(ctx context.Context, houseID string, action string, r *http.Request) (HouseActionResult, error)
}

// HouseWebActionResult is the browser response for an extension action.
type HouseWebActionResult struct {
	Status      int
	ContentType string
	Render      func(context.Context, io.Writer) error
}

// HouseWebActionHandler handles browser extension actions scoped to a house.
type HouseWebActionHandler interface {
	ExtensionID() string
	HandleHouseWebAction(ctx context.Context, houseID string, action string, r *http.Request) (HouseWebActionResult, error)
}

// Registry stores extension action handlers.
type Registry struct {
	houseActions    map[string]HouseActionHandler
	houseWebActions map[string]HouseWebActionHandler
}

// NewRegistry creates an empty extension registry.
func NewRegistry() *Registry {
	return &Registry{
		houseActions:    map[string]HouseActionHandler{},
		houseWebActions: map[string]HouseWebActionHandler{},
	}
}

// RegisterHouseAction registers an API house action handler.
func (r *Registry) RegisterHouseAction(handler HouseActionHandler) error {
	if handler == nil {
		return fmt.Errorf("extension house action handler is nil")
	}
	id := handler.ExtensionID()
	if id == "" {
		return fmt.Errorf("extension ID is required")
	}
	if _, exists := r.houseActions[id]; exists {
		return fmt.Errorf("extension house action handler %q already registered", id)
	}
	r.houseActions[id] = handler
	return nil
}

// RegisterHouseWebAction registers a browser house action handler.
func (r *Registry) RegisterHouseWebAction(handler HouseWebActionHandler) error {
	if handler == nil {
		return fmt.Errorf("extension house web action handler is nil")
	}
	id := handler.ExtensionID()
	if id == "" {
		return fmt.Errorf("extension ID is required")
	}
	if _, exists := r.houseWebActions[id]; exists {
		return fmt.Errorf("extension house web action handler %q already registered", id)
	}
	r.houseWebActions[id] = handler
	return nil
}

// HouseAction finds an API house action handler.
func (r *Registry) HouseAction(extensionID string) (HouseActionHandler, bool) {
	if r == nil {
		return nil, false
	}
	handler, ok := r.houseActions[extensionID]
	return handler, ok
}

// HouseWebAction finds a browser house action handler.
func (r *Registry) HouseWebAction(extensionID string) (HouseWebActionHandler, bool) {
	if r == nil {
		return nil, false
	}
	handler, ok := r.houseWebActions[extensionID]
	return handler, ok
}
