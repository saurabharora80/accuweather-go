package upstream

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Region struct {
	ID            string `json:"ID"`
	LocalizedName string `json:"LocalizedName"`
	EnglishName   string `json:"EnglishName"`
}

func GetRegions(client *resty.Client) ([]Region, error) {
	resp, err := client.R().
		Get("/locations/v1/regions")

	if err == nil && resp.StatusCode() == 200 {
		var regions []Region
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
