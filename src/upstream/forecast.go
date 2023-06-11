package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Forecast struct {
	MinimumTemp float64 `json:"DailyForecasts.Temperature.Minimum.Value"`
	MaximumTemp float64 `json:"DailyForecasts.Temperature.Maximum.Value"`
	TempUnit    string  `json:"DailyForecasts.Temperature.Minimum.Unit"`
	Sunrise     string  `json:"DailyForecasts.Sun.Rise"`
	Sunset      string  `json:"DailyForecasts.Sun.Set"`
}

func GetCityForecast(client *resty.Client, locationKey string, daysOfForecast int, forecastChan chan<- Forecast, errorsChan chan<- error) {
	resp, httpError := client.R().
		SetPathParam("locationKey", locationKey).
		SetQueryParam("offset", "1").
		Get(fmt.Sprintf("/forecasts/v1/daily/%dday/{locationKey}?metric=true", daysOfForecast))

	if httpError != nil {
		errorsChan <- httpError
		return
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var forecasts []Forecast
		err := json.Unmarshal(resp.Body(), &forecasts)
		if err != nil {
			errorsChan <- err
			return
		}
		forecastChan <- forecasts[0]
	case http.StatusNotFound:
		forecastChan <- Forecast{}
	default:
		errorsChan <- &HttpError{Method: "GET",
			Path:       fmt.Sprintf("/forecasts/v1/daily/%dday/%s?metric=true", daysOfForecast, locationKey),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
