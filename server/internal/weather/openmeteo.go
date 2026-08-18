package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultOpenMeteoBaseURL = "https://api.open-meteo.com/v1/forecast"

// OpenMeteoProvider fetches current weather from Open-Meteo.
type OpenMeteoProvider struct {
	Client  *http.Client
	BaseURL string
}

// NewOpenMeteoProvider creates an Open-Meteo backed weather provider.
func NewOpenMeteoProvider(client *http.Client) *OpenMeteoProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &OpenMeteoProvider{
		Client:  client,
		BaseURL: defaultOpenMeteoBaseURL,
	}
}

// Current fetches current weather for the requested location and unit system.
func (p *OpenMeteoProvider) Current(ctx context.Context, req CurrentRequest) (*Current, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	req = req.Normalize()

	endpoint, err := p.currentURL(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("fetch open-meteo current weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("open-meteo returned status %d", resp.StatusCode)
	}

	var payload openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode open-meteo response: %w", err)
	}

	current, err := payload.toCurrent()
	if err != nil {
		return nil, err
	}
	return current, nil
}

func (p *OpenMeteoProvider) currentURL(req CurrentRequest) (string, error) {
	base := p.BaseURL
	if base == "" {
		base = defaultOpenMeteoBaseURL
	}

	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse open-meteo base url: %w", err)
	}

	q := u.Query()
	q.Set("latitude", strconv.FormatFloat(req.Location.Latitude, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(req.Location.Longitude, 'f', -1, 64))
	q.Set("current", "temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m")
	q.Set("timezone", "UTC")

	if req.Units == UnitSystemImperial {
		q.Set("temperature_unit", "fahrenheit")
		q.Set("wind_speed_unit", "mph")
	} else {
		q.Set("temperature_unit", "celsius")
		q.Set("wind_speed_unit", "kmh")
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

type openMeteoResponse struct {
	CurrentUnits openMeteoCurrentUnits `json:"current_units"`
	Current      *openMeteoCurrent     `json:"current"`
}

type openMeteoCurrentUnits struct {
	Temperature string `json:"temperature_2m"`
	WindSpeed   string `json:"wind_speed_10m"`
}

type openMeteoCurrent struct {
	Time            string   `json:"time"`
	Temperature     *float64 `json:"temperature_2m"`
	HumidityPercent *float64 `json:"relative_humidity_2m"`
	WeatherCode     *int     `json:"weather_code"`
	WindSpeed       *float64 `json:"wind_speed_10m"`
}

func (r openMeteoResponse) toCurrent() (*Current, error) {
	if r.Current == nil || r.Current.Temperature == nil || r.Current.WeatherCode == nil || r.Current.Time == "" {
		return nil, ErrMissingCurrentWeather
	}

	observedAt, err := parseOpenMeteoTime(r.Current.Time)
	if err != nil {
		return nil, fmt.Errorf("parse open-meteo current time: %w", err)
	}

	text, icon := conditionFromWMOCode(*r.Current.WeatherCode)
	return &Current{
		Temperature:     *r.Current.Temperature,
		TemperatureUnit: r.CurrentUnits.Temperature,
		ConditionCode:   *r.Current.WeatherCode,
		ConditionText:   text,
		ConditionIcon:   icon,
		HumidityPercent: r.Current.HumidityPercent,
		WindSpeed:       r.Current.WindSpeed,
		WindSpeedUnit:   r.CurrentUnits.WindSpeed,
		ObservedAt:      observedAt,
	}, nil
}

func parseOpenMeteoTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format %q", value)
}

func conditionFromWMOCode(code int) (text string, icon string) {
	switch code {
	case 0:
		return "Clear sky", "clear"
	case 1:
		return "Mainly clear", "mostly_clear"
	case 2:
		return "Partly cloudy", "partly_cloudy"
	case 3:
		return "Overcast", "cloudy"
	case 45, 48:
		return "Fog", "fog"
	case 51, 53, 55:
		return "Drizzle", "drizzle"
	case 56, 57:
		return "Freezing drizzle", "freezing_drizzle"
	case 61, 63, 65:
		return "Rain", "rain"
	case 66, 67:
		return "Freezing rain", "freezing_rain"
	case 71, 73, 75:
		return "Snow", "snow"
	case 77:
		return "Snow grains", "snow"
	case 80, 81, 82:
		return "Rain showers", "rain_showers"
	case 85, 86:
		return "Snow showers", "snow_showers"
	case 95:
		return "Thunderstorm", "thunderstorm"
	case 96, 99:
		return "Thunderstorm with hail", "thunderstorm_hail"
	default:
		return "Unknown", "unknown"
	}
}
