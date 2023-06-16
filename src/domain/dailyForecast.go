package domain

type DailyForecast struct {
	MinimumTemp float64 `json:"minimum_temp"`
	MaximumTemp float64 `json:"maximum_temp"`
	TempUnit    string  `json:"temp_unit"`
	Sunrise     string  `json:"sunrise"`
	Sunset      string  `json:"sunset"`
}

func (f *DailyForecast) IsEmpty() bool {
	return f.MinimumTemp == 0 && f.MaximumTemp == 0 && f.Sunrise == "" && f.Sunset == ""
}
