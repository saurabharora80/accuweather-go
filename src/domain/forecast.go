package domain

type Forecast struct {
	MinimumTemp float64 `json:"DailyForecasts.Temperature.Minimum.Value"`
	MaximumTemp float64 `json:"DailyForecasts.Temperature.Maximum.Value"`
	TempUnit    string  `json:"DailyForecasts.Temperature.Minimum.Unit"`
	Sunrise     string  `json:"DailyForecasts.Sun.Rise"`
	Sunset      string  `json:"DailyForecasts.Sun.Set"`
}
