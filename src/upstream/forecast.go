package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/domain"
)

type ForecastConnector struct {
	Client       *resty.Client
	ForecastChan chan<- domain.Forecast
	ErrorsChan   chan<- error
}

func NewForecastConnector(client *resty.Client, forecastChan chan<- domain.Forecast, errorsChan chan<- error) *ForecastConnector {
	return &ForecastConnector{Client: client, ForecastChan: forecastChan, ErrorsChan: errorsChan}
}

func (c *ForecastConnector) GetCityForecast(locationKey string, daysOfForecast int) {
	resp, httpError := c.Client.R().
		SetPathParam("locationKey", locationKey).
		SetQueryParam("offset", "1").
		Get(fmt.Sprintf("/forecasts/v1/daily/%dday/{locationKey}?metric=true", daysOfForecast))

	if httpError != nil {
		c.ErrorsChan <- httpError
		return
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var forecasts []domain.Forecast
		err := json.Unmarshal(resp.Body(), &forecasts)
		if err != nil {
			c.ErrorsChan <- err
			return
		}
		c.ForecastChan <- forecasts[0]
	case http.StatusNotFound:
		c.ForecastChan <- domain.Forecast{}
	default:
		c.ErrorsChan <- &domain.HttpError{Method: "GET",
			Path:       fmt.Sprintf("/forecasts/v1/daily/%dday/%s?metric=true", daysOfForecast, locationKey),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
