package api

// HouseInfo contains house details for a user.
type HouseInfo struct {
	ID          string       `json:"id"`
	DisplayName string       `json:"displayName"`
	Location    *string      `json:"location,omitempty"`
	Role        string       `json:"role"`
	Weather     *WeatherInfo `json:"weather"`
}

// WeatherInfo contains current weather plugin data for a house.
type WeatherInfo struct {
	Status          string   `json:"status"`
	LocationLabel   string   `json:"locationLabel,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	TemperatureUnit string   `json:"temperatureUnit,omitempty"`
	ConditionCode   *int     `json:"conditionCode,omitempty"`
	ConditionText   string   `json:"conditionText,omitempty"`
	ConditionIcon   string   `json:"conditionIcon,omitempty"`
	HumidityPercent *float64 `json:"humidityPercent,omitempty"`
	WindSpeed       *float64 `json:"windSpeed,omitempty"`
	WindSpeedUnit   string   `json:"windSpeedUnit,omitempty"`
	ObservedAt      string   `json:"observedAt,omitempty"`
	UpdatedAt       string   `json:"updatedAt,omitempty"`
	Error           string   `json:"error,omitempty"`
}

// HouseListResponse represents the house list endpoint response.
type HouseListResponse struct {
	Houses []HouseInfo `json:"houses"`
}

// WeatherRefreshResponse represents the force weather refresh endpoint response.
type WeatherRefreshResponse struct {
	Weather *WeatherInfo `json:"weather"`
}
