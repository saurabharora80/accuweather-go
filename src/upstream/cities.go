package upstream

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type City struct {
	Key  string `json:"Key"`
	Name string `json:"EnglishName"`
}

func GetCityInCountry(client *resty.Client, countryCode string, city string, citiesChan chan City, errorsChan chan error) {
	resp, httpError := client.R().
		SetPathParam("countryCode", countryCode).
		SetQueryParam("q", city).
		SetQueryParam("offset", "1").
		Get("/locations/v1/cities/{countryCode}/search")

	if resp != nil && resp.StatusCode() == 200 {
		var cities []City
		err := json.Unmarshal(resp.Body(), &cities)
		if err != nil {
			errorsChan <- err
			return
		}
		citiesChan <- cities[0]
		return
	} else if resp != nil && resp.StatusCode() == 404 {
		citiesChan <- City{}
		return
	}
	errorsChan <- httpError
}
