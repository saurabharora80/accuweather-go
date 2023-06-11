package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"org.example/hello/src/domain"
	"org.example/hello/src/service"
	"testing"
)

type MockForecastConnector struct {
	mock.Mock
}

type MockCityConnector struct {
	mock.Mock
}

func (m *MockForecastConnector) GetCityForecast(cityKey string, daysOfForecast int, forecastChan chan domain.Forecast, errorsChan chan error) {
	m.Called(cityKey, daysOfForecast)
	if cityKey == "VALID_CTY_KEY" {
		forecastChan <- domain.Forecast{MinimumTemp: 9.0, MaximumTemp: 20.0}
	} else if cityKey == "IN_VALID_CTY_KEY" {
		errorsChan <- &domain.HttpError{Method: "GET", Path: "/forecasts/v1/daily/5day/IN_VALID_CTY_KEY?metric=true", StatusCode: 404, Details: []byte("Not Found")}
	} else {
		forecastChan <- domain.Forecast{}
	}
}

func (m *MockCityConnector) GetCityInCountry(countryCode string, city string, citiesChan chan domain.City, errorsChan chan error) {
	m.Called(countryCode, city)
	if city == "VALID" {
		citiesChan <- domain.City{Key: "VALID_CTY_KEY"}
	} else if city == "IN_VALID" {
		errorsChan <- &domain.HttpError{Method: "GET", Path: "/locations/v1/cities/AU/search?q=IN_VALID", StatusCode: 404, Details: []byte("Not Found")}
	} else if city == "IN_VALID_CITY_KEY" {
		citiesChan <- domain.City{Key: "IN_VALID_CTY_KEY"}
	} else {
		citiesChan <- domain.City{}
	}
}

func TestWeatherValid(t *testing.T) {

	mockForecastConnector := new(MockForecastConnector)
	mockCityConnector := new(MockCityConnector)

	mockCityConnector.On("GetCityInCountry", "AU", "VALID").Return()
	mockForecastConnector.On("GetCityForecast", "VALID_CTY_KEY", 5).Return()

	weatherService := service.NewWeatherService(mockCityConnector, mockForecastConnector)

	forecast, _ := weatherService.GetCityForecast("AU", "VALID", 5)

	assert.Equal(t, domain.Forecast{MinimumTemp: 9.0, MaximumTemp: 20.0}, forecast)

}

func TestWeatherInValidCity(t *testing.T) {

	mockForecastConnector := new(MockForecastConnector)
	mockCityConnector := new(MockCityConnector)

	mockCityConnector.On("GetCityInCountry", "AU", "IN_VALID").Return()

	weatherService := service.NewWeatherService(mockCityConnector, mockForecastConnector)

	_, err := weatherService.GetCityForecast("AU", "IN_VALID", 5)

	assert.Equal(t, "GET /locations/v1/cities/AU/search?q=IN_VALID failed with 404 => Not Found", err.Error())

}

func TestWeatherInValidForecastCity(t *testing.T) {

	mockForecastConnector := new(MockForecastConnector)
	mockCityConnector := new(MockCityConnector)

	mockCityConnector.On("GetCityInCountry", "AU", "IN_VALID_CITY_KEY").Return()
	mockForecastConnector.On("GetCityForecast", "IN_VALID_CTY_KEY", 5).Return()

	weatherService := service.NewWeatherService(mockCityConnector, mockForecastConnector)

	_, err := weatherService.GetCityForecast("AU", "IN_VALID_CITY_KEY", 5)

	assert.Equal(t, "GET /forecasts/v1/daily/5day/IN_VALID_CTY_KEY?metric=true failed with 404 => Not Found", err.Error())

}
