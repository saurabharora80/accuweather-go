package upstream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/walkerus/go-wiremock"
	"org.example/hello/src/upstream"
	"org.example/hello/src/upstream/model"
	"testing"
)

func (suite *ForecastTestSuite) TestGetCityWeatherForecast() {

	err := suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/forecasts/v1/daily/1day/123")).
		WithQueryParam("metric", wiremock.EqualTo("true")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnFileContent("upstream_responses/forecast.json", map[string]string{"Content-Type": "application/json"}, 200))

	assert.NoError(suite.T(), err)

	forecasts := make(chan model.Forecast)
	errors := make(chan error)

	connector := upstream.NewForecastConnector(suite.HttpClient)

	go connector.GetCityForecast("123", 1, forecasts, errors)

	select {
	case actualForecast := <-forecasts:
		fmt.Println(actualForecast.DailyForecasts[0].Temperature.Minimum.Value)
		assert.Equal(suite.T(), 19.3, actualForecast.DailyForecasts[0].Temperature.Maximum.Value)
		assert.Equal(suite.T(), 8.3, actualForecast.DailyForecasts[0].Temperature.Minimum.Value)
	case err := <-errors:
		assert.Fail(suite.T(), "Unable to get Forecast because of %s", err.Error())
	}
}

func (suite *ForecastTestSuite) TestGetCityWeatherForecastNotFound() {

	forecasts := make(chan model.Forecast)
	errors := make(chan error)

	connector := upstream.NewForecastConnector(suite.HttpClient)

	go connector.GetCityForecast("123", 1, forecasts, errors)

	select {
	case city := <-forecasts:
		assert.Equal(suite.T(), model.Forecast{}, city)
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

		forecasts := make(chan model.Forecast)
		errors := make(chan error)

		connector := upstream.NewForecastConnector(suite.HttpClient)

		go connector.GetCityForecast("123", 1, forecasts, errors)

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
