package upstream

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type RegionOrCountry struct {
	ID            string `json:"ID"`
	LocalizedName string `json:"LocalizedName"`
	EnglishName   string `json:"EnglishName"`
}

func GetRegions(client *resty.Client) ([]RegionOrCountry, error) {
	resp, err := client.R().Get("/locations/v1/regions")

	if err == nil && resp.StatusCode() == 200 {
		var regions []RegionOrCountry
		err := json.Unmarshal(resp.Body(), &regions)
		if err == nil {
			return regions, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func GetCountries(client *resty.Client, region string) ([]RegionOrCountry, error) {
	resp, err := client.R().
		SetPathParam("region", region).
		Get("/locations/v1/countries/{region}")

	if err == nil && resp.StatusCode() == 200 {
		var regionsOrCountries []RegionOrCountry
		err := json.Unmarshal(resp.Body(), &regionsOrCountries)
		if err == nil {
			return regionsOrCountries, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
