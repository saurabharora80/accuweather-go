package upstream

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/walkerus/go-wiremock"
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

	forecasts := make(chan upstream.Forecast)
	errors := make(chan error)

	go upstream.GetCityForecast(suite.HttpClient, "123", 1, forecasts, errors)

	select {
	case city := <-forecasts:
		assert.Equal(suite.T(),
			upstream.Forecast{
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

	forecasts := make(chan upstream.Forecast)
	errors := make(chan error)

	go upstream.GetCityForecast(suite.HttpClient, "123", 1, forecasts, errors)

	select {
	case city := <-forecasts:
		assert.Equal(suite.T(), upstream.Forecast{}, city)
	case err := <-errors:
		assert.Fail(suite.T(), "Unable to get Forecast because of %s", err.Error())
	}
}

func TestForecastSuite(t *testing.T) {
	suite.Run(t, new(ForecastTestSuite))
}
