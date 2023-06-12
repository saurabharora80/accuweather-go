package model

type Forecast struct {
	MinimumTemp float64 `json:"minimum-temp"`
	MaximumTemp float64 `json:"maximum-temp"`
	TempUnit    string  `json:"temp-unit"`
	Sunrise     string  `json:"sunrise"`
	Sunset      string  `json:"sunset"`
}
