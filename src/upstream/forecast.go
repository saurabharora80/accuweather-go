package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/domain"
)

type ForecastConnectorInterface interface {
	GetCityForecast(locationKey string, daysOfForecast int, forecastChan chan domain.Forecast, errorsChan chan error)
}

type ForecastConnector struct {
	Client *resty.Client
}

func NewForecastConnector(client *resty.Client) *ForecastConnector {
	return &ForecastConnector{Client: client}
}

func (c *ForecastConnector) GetCityForecast(locationKey string, daysOfForecast int, forecastChan chan domain.Forecast, errorsChan chan error) {
	resp, httpError := c.Client.R().
		SetPathParam("locationKey", locationKey).
		SetQueryParam("offset", "1").
		Get(fmt.Sprintf("/forecasts/v1/daily/%dday/{locationKey}?metric=true", daysOfForecast))

	if httpError != nil {
		errorsChan <- httpError
		return
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var forecasts []domain.Forecast
		err := json.Unmarshal(resp.Body(), &forecasts)
		if err != nil {
			errorsChan <- err
			return
		}
		forecastChan <- forecasts[0]
	case http.StatusNotFound:
		forecastChan <- domain.Forecast{}
	default:
		errorsChan <- &domain.HttpError{Method: "GET",
			Path:       fmt.Sprintf("/forecasts/v1/daily/%dday/%s?metric=true", daysOfForecast, locationKey),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
