package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/upstream"
	"time"
)

func main() {

	config, config_err := config()

	if config_err != nil {
		fmt.Println(config_err)
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
			fmt.Println(region.ID)
		}
	} else {
		fmt.Println(err)
	}

}
