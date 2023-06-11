package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/domain"
)

type CityConnector struct {
	Client     *resty.Client
	CitiesChan chan<- domain.City
	ErrorsChan chan<- error
}

func NewCityConnector(client *resty.Client, citiesChan chan<- domain.City, errorsChan chan<- error) *CityConnector {
	return &CityConnector{Client: client, CitiesChan: citiesChan, ErrorsChan: errorsChan}
}

func (c *CityConnector) GetCityInCountry(countryCode string, city string) {
	resp, httpError := c.Client.R().
		SetPathParam("countryCode", countryCode).
		SetQueryParam("q", city).
		SetQueryParam("offset", "1").
		Get("/locations/v1/cities/{countryCode}/search")

	if httpError != nil {
		c.ErrorsChan <- httpError
		return
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var cities []domain.City
		err := json.Unmarshal(resp.Body(), &cities)
		if err != nil {
			c.ErrorsChan <- err
			return
		}
		c.CitiesChan <- cities[0]
	case http.StatusNotFound:
		c.CitiesChan <- domain.City{}
	default:
		c.ErrorsChan <- &domain.HttpError{Method: "GET",
			Path:       fmt.Sprintf("/locations/v1/cities/%s/search?q=%s", countryCode, city),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
