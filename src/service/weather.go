package service

import (
	"org.example/hello/src/domain"
	"org.example/hello/src/upstream"
)

type WeatherService struct {
	CityConnector     upstream.CityConnectorInterface
	ForecastConnector upstream.ForecastConnectorInterface
}

func NewWeatherService(cityConnector upstream.CityConnectorInterface,
	forecastConnector upstream.ForecastConnectorInterface) *WeatherService {
	return &WeatherService{CityConnector: cityConnector, ForecastConnector: forecastConnector}
}

func (s *WeatherService) GetCityForecast(countryCode string, city string, daysOfForecast int) (domain.Forecast, error) {
	citiesChan := make(chan domain.City)
	forecastChan := make(chan domain.Forecast)
	errorsChan := make(chan error)

	go s.CityConnector.GetCityInCountry(countryCode, city, citiesChan, errorsChan)

	select {
	case city := <-citiesChan:
		go s.ForecastConnector.GetCityForecast(city.Key, daysOfForecast, forecastChan, errorsChan)
		select {
		case forecast := <-forecastChan:
			return forecast, nil
		case err := <-errorsChan:
			return domain.Forecast{}, err
		}
	case err := <-errorsChan:
		return domain.Forecast{}, err
	}
}
