package main

import (
	"fmt"
	"org.example/hello/src/service"
)

func main() {

	weatherService, errors := service.GetWeatherInstance()

	if errors != nil {
		return
	}

	cityForecast, err := weatherService.GetCityForecast("AU", "Sydney", 5)

	if err == nil {
		fmt.Println(cityForecast)
	} else {
		fmt.Printf("failed to weather forecast: %s", err.Error())
	}

}
