package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/upstream"
	"time"
)

func main() {

	config, configErr := config()

	if configErr != nil {
		fmt.Println(configErr)
		return
	}

	client := resty.New().
		EnableTrace().
		SetTransport(&http.Transport{
			MaxIdleConns:    config.Upstream.MaxIdleConnections,
			IdleConnTimeout: config.Upstream.IdleConnectionTimeoutSeconds * time.Second}).
		SetQueryParam("apikey", config.Upstream.Key).
		SetHeader("Accept", "application/json").
		SetBaseURL(config.Upstream.Host)

	regions, err := upstream.GetRegions(client)

	if err == nil {
		for _, region := range regions {
			countries, err := upstream.GetCountries(client, region.ID)
			if err == nil {
				for _, country := range countries {
					fmt.Printf("%v -> %v\n", region.ID, country.ID)
				}
			}
		}
	} else {
		fmt.Println(err)
	}

}
