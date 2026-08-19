package extension

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
)

// HouseExtraProvider provides API extension data scoped to a house.
type HouseExtraProvider interface {
	ExtensionID() string
	HouseExtra(ctx context.Context, houseID string) (any, error)
}

// HouseCardProvider provides browser UI cards scoped to a house.
type HouseCardProvider interface {
	ExtensionID() string
	HouseCard(ctx context.Context, houseID string, r *http.Request) (templ.Component, error)
}

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
	houseExtras     map[string]HouseExtraProvider
	houseCards      map[string]HouseCardProvider
	houseActions    map[string]HouseActionHandler
	houseWebActions map[string]HouseWebActionHandler
}

// NewRegistry creates an empty extension registry.
func NewRegistry() *Registry {
	return &Registry{
		houseExtras:     map[string]HouseExtraProvider{},
		houseCards:      map[string]HouseCardProvider{},
		houseActions:    map[string]HouseActionHandler{},
		houseWebActions: map[string]HouseWebActionHandler{},
	}
}

// RegisterHouseExtra registers an API house extra provider.
func (r *Registry) RegisterHouseExtra(provider HouseExtraProvider) error {
	if provider == nil {
		return fmt.Errorf("extension house extra provider is nil")
	}
	id := provider.ExtensionID()
	if id == "" {
		return fmt.Errorf("extension ID is required")
	}
	if _, exists := r.houseExtras[id]; exists {
		return fmt.Errorf("extension house extra provider %q already registered", id)
	}
	r.houseExtras[id] = provider
	return nil
}

// RegisterHouseCard registers a browser house card provider.
func (r *Registry) RegisterHouseCard(provider HouseCardProvider) error {
	if provider == nil {
		return fmt.Errorf("extension house card provider is nil")
	}
	id := provider.ExtensionID()
	if id == "" {
		return fmt.Errorf("extension ID is required")
	}
	if _, exists := r.houseCards[id]; exists {
		return fmt.Errorf("extension house card provider %q already registered", id)
	}
	r.houseCards[id] = provider
	return nil
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

// HouseExtras returns registered API house extra providers.
func (r *Registry) HouseExtras() []HouseExtraProvider {
	if r == nil {
		return nil
	}
	providers := make([]HouseExtraProvider, 0, len(r.houseExtras))
	for _, provider := range r.houseExtras {
		providers = append(providers, provider)
	}
	return providers
}

// HouseCards returns registered browser house card providers.
func (r *Registry) HouseCards() []HouseCardProvider {
	if r == nil {
		return nil
	}
	providers := make([]HouseCardProvider, 0, len(r.houseCards))
	for _, provider := range r.houseCards {
		providers = append(providers, provider)
	}
	return providers
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
