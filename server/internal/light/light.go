package light

// Preset describes a supported light color preset.
type Preset struct {
	ID    string
	Label string
	Hex   string
}

var Presets = []Preset{
	{ID: "warm_white", Label: "Warm White", Hex: "#FFDCA8"},
	{ID: "soft_white", Label: "Soft White", Hex: "#FFF1D6"},
	{ID: "daylight_white", Label: "Daylight White", Hex: "#F7FBFF"},
	{ID: "cold_white", Label: "Cold White", Hex: "#DDEBFF"},
	{ID: "red", Label: "Red", Hex: "#FF3B30"},
	{ID: "orange", Label: "Orange", Hex: "#FF9500"},
	{ID: "yellow", Label: "Yellow", Hex: "#FFCC00"},
	{ID: "green", Label: "Green", Hex: "#34C759"},
	{ID: "blue", Label: "Blue", Hex: "#007AFF"},
	{ID: "purple", Label: "Purple", Hex: "#AF52DE"},
}

func IsValidBrightness(value int) bool {
	return value >= 0 && value <= 100 && value%5 == 0
}

func IsValidPreset(id string) bool {
	for _, preset := range Presets {
		if preset.ID == id {
			return true
		}
	}
	return false
}
