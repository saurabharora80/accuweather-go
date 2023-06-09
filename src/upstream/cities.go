package upstream

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type City struct {
	Key  string `json:"Key"`
	Name string `json:"EnglishName"`
}

func GetCityInCountry(client *resty.Client, countryCode string, city string) (City, error) {
	resp, err := client.R().
		SetPathParam("countryCode", countryCode).
		SetQueryParam("q", city).
		SetQueryParam("offset", "1").
		Get("/locations/v1/cities/{countryCode}/search")

	if err == nil && resp.StatusCode() == 200 {
		var cities []City
		err := json.Unmarshal(resp.Body(), &cities)
		return cities[0], err
	}

	return City{}, err
}
