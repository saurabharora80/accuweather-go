package upstream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/walkerus/go-wiremock"
	"org.example/hello/src/domain"
	"org.example/hello/src/upstream"
	"testing"
)

func (suite *ForecastTestSuite) TestGetCityWeatherForecast() {

	err := suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/forecasts/v1/daily/1day/123")).
		WithQueryParam("metric", wiremock.EqualTo("true")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnJSON(
			[]map[string]interface{}{{
				"DailyForecasts.Temperature.Minimum.Value": 9.0,
				"DailyForecasts.Temperature.Maximum.Value": 21.0,
				"DailyForecasts.Temperature.Minimum.Unit":  "C",
				"DailyForecasts.Temperature.Maximum.Unit":  "C",
				"DailyForecasts.Sun.Rise":                  "2019-05-15T06:44:00+10:00",
				"DailyForecasts.Sun.Set":                   "2019-05-15T17:01:00+10:00",
			}}, map[string]string{"Content-Type": "application/json"}, 200,
		))

	if err != nil {
		assert.Fail(suite.T(), "Failed to configure wiremock stub due to %s", err.Error())
		return
	}

	forecasts := make(chan domain.Forecast)
	errors := make(chan error)

	connector := upstream.NewForecastConnector(suite.HttpClient, forecasts, errors)

	go connector.GetCityForecast("123", 1)

	select {
	case city := <-forecasts:
		assert.Equal(suite.T(),
			domain.Forecast{
				MinimumTemp: 9.0,
				MaximumTemp: 21.0,
				TempUnit:    "C",
				Sunrise:     "2019-05-15T06:44:00+10:00",
				Sunset:      "2019-05-15T17:01:00+10:00"},
			city)
	case err := <-errors:
		assert.Fail(suite.T(), "Unable to get Forecast because of %s", err.Error())
	}
}

func (suite *ForecastTestSuite) TestGetCityWeatherForecastNotFound() {

	forecasts := make(chan domain.Forecast)
	errors := make(chan error)

	connector := upstream.NewForecastConnector(suite.HttpClient, forecasts, errors)

	go connector.GetCityForecast("123", 1)

	select {
	case city := <-forecasts:
		assert.Equal(suite.T(), domain.Forecast{}, city)
	case err := <-errors:
		assert.Fail(suite.T(), "Unable to get Forecast because of %s", err.Error())
	}
}

func (suite *ForecastTestSuite) TestGetCityWeatherForecastFailed() {
	invalidStatusCodes := []int64{403, 500, 501, 503}

	for _, invalidStatusCode := range invalidStatusCodes {
		_ = suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/forecasts/v1/daily/1day/123")).
			WithQueryParam("metric", wiremock.EqualTo("true")).
			WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
			WillReturnJSON(
				[]map[string]interface{}{{}}, map[string]string{"Content-Type": "application/json"}, invalidStatusCode,
			))

		forecasts := make(chan domain.Forecast)
		errors := make(chan error)

		connector := upstream.NewForecastConnector(suite.HttpClient, forecasts, errors)

		go connector.GetCityForecast("123", 1)

		select {
		case _ = <-forecasts:
			assert.Fail(suite.T(), "Shouldn't return forecast")
		case err := <-errors:
			assert.Equal(suite.T(), fmt.Sprintf("GET /forecasts/v1/daily/1day/123?metric=true failed with %d => [{}]", invalidStatusCode), err.Error())
		}
	}
}

func TestForecastSuite(t *testing.T) {
	suite.Run(t, new(ForecastTestSuite))
}
