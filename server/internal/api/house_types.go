package api

// HouseInfo contains house details for a user.
type HouseInfo struct {
	ID          string         `json:"id"`
	DisplayName string         `json:"displayName"`
	Location    *string        `json:"location,omitempty"`
	Role        string         `json:"role"`
	Extras      map[string]any `json:"extras,omitempty"`
}

// HouseListResponse represents the house list endpoint response.
type HouseListResponse struct {
	Houses []HouseInfo `json:"houses"`
}
