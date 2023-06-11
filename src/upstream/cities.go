package upstream

import (
	"encoding/json"
	"fmt"
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

	if resp != nil {
		if resp.StatusCode() == 200 {
			var cities []City
			err := json.Unmarshal(resp.Body(), &cities)
			if err != nil {
				errorsChan <- err
			}
			citiesChan <- cities[0]
		} else if resp.StatusCode() == 404 {
			citiesChan <- City{}
		} else {
			errorsChan <- &HttpError{Method: "GET",
				Path:       fmt.Sprintf("/locations/v1/cities/%s/search?q=%s", countryCode, city),
				StatusCode: resp.StatusCode(),
				Details:    resp.Body()}
		}
	}

	errorsChan <- httpError
}
