package api

// DeviceIntegrationInfo contains integration details.
type DeviceIntegrationInfo struct {
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
}

// DeviceInfo contains device details.
type DeviceInfo struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Integration  DeviceIntegrationInfo  `json:"integration"`
	State        string                 `json:"state"`
	Availability string                 `json:"availability"`
	Capabilities DeviceCapabilitiesInfo `json:"capabilities"`
	Light        *DeviceLightInfo       `json:"light,omitempty"`
}

type DeviceCapabilitiesInfo struct {
	Light        bool `json:"light"`
	Brightness   bool `json:"brightness,omitempty"`
	ColorPresets bool `json:"colorPresets,omitempty"`
}

type DeviceLightInfo struct {
	Brightness  *int    `json:"brightness,omitempty"`
	ColorPreset *string `json:"colorPreset,omitempty"`
}

// HouseRef contains house references for device responses.
type HouseRef struct {
	ID string `json:"id"`
}

// DeviceInfoWithHouse contains device details with house reference.
type DeviceInfoWithHouse struct {
	DeviceInfo
	House HouseRef `json:"house"`
}

// DeviceListResponse represents the device list endpoint response.
type DeviceListResponse struct {
	Devices []DeviceInfoWithHouse `json:"devices"`
}

// DeviceListByHouseResponse represents the house device list endpoint response.
type DeviceListByHouseResponse struct {
	Devices []DeviceInfo `json:"devices"`
}

// DeviceSetStateRequest represents a request to set the desired state for a device.
type DeviceSetStateRequest struct {
	State string `json:"state"`
}

// DeviceSetStateResponse represents the result of setting a device state.
type DeviceSetStateResponse struct {
	DeviceID string   `json:"deviceId"`
	House    HouseRef `json:"house"`
	State    string   `json:"state"`
}

type DeviceSetLightRequest struct {
	Brightness  *int    `json:"brightness"`
	ColorPreset *string `json:"colorPreset"`
}

type DeviceSetLightResponse struct {
	DeviceID     string           `json:"deviceId"`
	House        HouseRef         `json:"house"`
	State        string           `json:"state"`
	Availability string           `json:"availability"`
	Light        *DeviceLightInfo `json:"light"`
}
