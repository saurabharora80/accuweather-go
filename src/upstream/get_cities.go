package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/domain"
	"sync"
)

type CityConnectorInterface interface {
	GetCityInCountry(countryCode string, city string, citiesChan chan domain.City, errorsChan chan error)
}

type CityConnector struct {
	Client *resty.Client
}

var (
	CityConnectorInstance      *CityConnector
	CityConnectorInstanceError error
	onceForCityConnector       sync.Once
)

func GetCityConnectorInstance() (*CityConnector, error) {
	onceForCityConnector.Do(func() {
		client, err := NewRestyClient()
		CityConnectorInstanceError = err
		CityConnectorInstance = NewCityConnector(client)
	})

	return CityConnectorInstance, CityConnectorInstanceError
}

func NewCityConnector(client *resty.Client) *CityConnector {
	return &CityConnector{Client: client}
}

func (c *CityConnector) GetCityInCountry(countryCode string, city string, citiesChan chan domain.City, errorsChan chan error) {
	resp, httpError := c.Client.R().
		SetPathParam("countryCode", countryCode).
		SetQueryParam("q", city).
		SetQueryParam("offset", "1").
		Get("/locations/v1/cities/{countryCode}/search")

	if httpError != nil {
		errorsChan <- httpError
		return
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var cities []domain.City
		err := json.Unmarshal(resp.Body(), &cities)
		if err != nil {
			errorsChan <- err
			return
		}
		citiesChan <- cities[0]
	case http.StatusNotFound:
		citiesChan <- domain.City{}
	default:
		errorsChan <- &domain.HttpError{Method: "GET",
			Path:       fmt.Sprintf("/locations/v1/cities/%s/search?q=%s", countryCode, city),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
