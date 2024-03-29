package service

import (
	"fmt"
	"org.example/hello/src/domain"
	"org.example/hello/src/upstream"
	"org.example/hello/src/upstream/model"
	"sync"
)

type WeatherService struct {
	CityConnector     upstream.CityConnectorInterface
	ForecastConnector upstream.ForecastConnectorInterface
}

var (
	WeatherServiceInstance       *WeatherService
	WeatherServiceInstanceErrors []error
	once                         sync.Once
)

func GetWeatherInstance() (*WeatherService, []error) {
	once.Do(func() {
		forecastConnectorInstance, fError := upstream.GetForecastConnectorInstance()

		if fError != nil {
			fmt.Printf("failed to get forecast connector instance: %s", fError.Error())
			WeatherServiceInstanceErrors = append(WeatherServiceInstanceErrors, fError)
		}

		cityConnectorInstance, cError := upstream.GetCityConnectorInstance()

		if cError != nil {
			fmt.Printf("failed to get city connector instance: %s", cError.Error())
			WeatherServiceInstanceErrors = append(WeatherServiceInstanceErrors, cError)
		}

		WeatherServiceInstance = NewWeatherService(cityConnectorInstance, forecastConnectorInstance)
	})

	return WeatherServiceInstance, WeatherServiceInstanceErrors

}

func NewWeatherService(cityConnector upstream.CityConnectorInterface,
	forecastConnector upstream.ForecastConnectorInterface) *WeatherService {
	return &WeatherService{CityConnector: cityConnector, ForecastConnector: forecastConnector}
}

func (s *WeatherService) GetCityForecast(countryCode string, city string, daysOfForecast int) (domain.DailyForecast, error) {
	citiesChan := make(chan domain.City)
	forecastChan := make(chan model.Forecast)
	errorsChan := make(chan error)

	go s.CityConnector.GetCityInCountry(countryCode, city, citiesChan, errorsChan)

	select {
	case city := <-citiesChan:
		go s.ForecastConnector.GetCityForecast(city.Key, daysOfForecast, forecastChan, errorsChan)
		select {
		case forecast := <-forecastChan:
			return from(forecast.DailyForecasts), nil
		case err := <-errorsChan:
			return domain.DailyForecast{}, err
		}
	case err := <-errorsChan:
		return domain.DailyForecast{}, err
	}
}

func from(forecasts []model.DailyForecast) domain.DailyForecast {
	if len(forecasts) == 0 {
		return domain.DailyForecast{}
	}
	return forecasts[0].To()
}
