package upstream

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type City struct {
	Key  string `json:"Key"`
	Name string `json:"EnglishName"`
}

func GetCityInCountry(client *resty.Client, countryCode string, city string, citiesChan chan<- City, errorsChan chan<- error) {
	resp, httpError := client.R().
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
		var cities []City
		err := json.Unmarshal(resp.Body(), &cities)
		if err != nil {
			errorsChan <- err
			return
		}
		citiesChan <- cities[0]
	case http.StatusNotFound:
		citiesChan <- City{}
	default:
		errorsChan <- &HttpError{Method: "GET",
			Path:       fmt.Sprintf("/locations/v1/cities/%s/search?q=%s", countryCode, city),
			StatusCode: resp.StatusCode(),
			Details:    resp.Body()}
	}
}
