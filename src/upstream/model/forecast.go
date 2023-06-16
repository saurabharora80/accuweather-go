package model

import "org.example/hello/src/domain"

type DailyForecast struct {
	Temperature struct {
		Minimum struct {
			Value float64 `json:"Value"`
			Unit  string  `json:"Unit"`
		} `json:"Minimum"`
		Maximum struct {
			Value float64 `json:"Value"`
			Unit  string  `json:"Unit"`
		} `json:"Maximum"`
	} `json:"Temperature"`
	Sun struct {
		Set  string `json:"Set"`
		Rise string `json:"Rise"`
	} `json:"Sun"`
}

func (f DailyForecast) To() domain.DailyForecast {
	return domain.DailyForecast{
		MinimumTemp: f.Temperature.Minimum.Value,
		MaximumTemp: f.Temperature.Maximum.Value,
		TempUnit:    f.Temperature.Maximum.Unit,
		Sunrise:     f.Sun.Rise,
		Sunset:      f.Sun.Set,
	}
}

type Forecast struct {
	DailyForecasts []DailyForecast `json:"DailyForecasts"`
}

func (f *Forecast) IsEmpty() bool {
	return f.DailyForecasts == nil || len(f.DailyForecasts) == 0
}
