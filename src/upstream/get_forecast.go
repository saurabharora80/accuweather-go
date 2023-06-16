package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/domain"
	"org.example/hello/src/upstream/model"
	"sync"
)

type ForecastConnectorInterface interface {
	GetCityForecast(locationKey string, daysOfForecast int, forecastChan chan model.Forecast, errorsChan chan error)
}

type ForecastConnector struct {
	Client *resty.Client
}

var (
	ForecastConnectorInstance      *ForecastConnector
	ForecastConnectorInstanceError error
	onceForForecastInstance        sync.Once
)

func GetForecastConnectorInstance() (*ForecastConnector, error) {
	onceForForecastInstance.Do(func() {
		client, err := NewRestyClient()
		ForecastConnectorInstanceError = err
		ForecastConnectorInstance = NewForecastConnector(client)
	})

	return ForecastConnectorInstance, ForecastConnectorInstanceError
}

func NewForecastConnector(client *resty.Client) *ForecastConnector {
	return &ForecastConnector{Client: client}
}

func (c *ForecastConnector) GetCityForecast(locationKey string, daysOfForecast int, forecastChan chan model.Forecast, errorsChan chan error) {
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
		var forecast model.Forecast
		err := json.Unmarshal(resp.Body(), &forecast)
		if err != nil {
			errorsChan <- err
			return
		}
		forecastChan <- forecast
	case http.StatusNotFound:
		forecastChan <- model.Forecast{}
	default:
		errorsChan <- &domain.HttpError{Method: "GET",
			Path:       fmt.Sprintf("/forecasts/v1/daily/%dday/%s?metric=true", daysOfForecast, locationKey),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
