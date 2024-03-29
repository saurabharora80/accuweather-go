package web

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/walkerus/go-wiremock"
	"net/http"
	"net/http/httptest"
	"org.example/hello/src/domain"
	"os"
	"testing"
)

func (suite *WebTestSuite) TestGetWeatherForecast() {

	_ = os.Setenv("ENABLE.RESTY.DEBUG", "true")

	setupCityStub(suite)

	setUpForecastStub(suite)

	request, err := http.NewRequest("GET", "/weather/AU/Sydney/1", nil)

	assert.NoError(suite.T(), err)

	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 200, recorder.Code)

	actualResponse := domain.DailyForecast{}

	assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &actualResponse))

	assert.Equal(suite.T(),
		domain.DailyForecast{MinimumTemp: 8.3, MaximumTemp: 19.3, TempUnit: "C", Sunrise: "2023-06-15T06:54:00+10:00", Sunset: "2023-06-15T16:56:00+10:00"},
		actualResponse)
}

func (suite *WebTestSuite) TestGetWeatherForecastCityNotFound() {

	request, err := http.NewRequest("GET", "/weather/AU/Sydney/1", nil)

	assert.NoError(suite.T(), err)

	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 404, recorder.Code)

	actualResponse := domain.DailyForecast{}

	assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &actualResponse))

	assert.Equal(suite.T(), domain.DailyForecast{}, actualResponse)
}

func (suite *WebTestSuite) TestGetWeatherForecastForecastNotFound() {

	setupCityStub(suite)

	request, err := http.NewRequest("GET", "/weather/AU/Sydney/1", nil)

	assert.NoError(suite.T(), err)

	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 404, recorder.Code)

	actualResponse := domain.DailyForecast{}

	assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &actualResponse))

	assert.Equal(suite.T(), domain.DailyForecast{}, actualResponse)
}

func (suite *WebTestSuite) TestGetWeatherForecastInvalidNoOfDays() {

	request, err := http.NewRequest("GET", "/weather/AU/Sydney/1day", nil)

	assert.NoError(suite.T(), err)

	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 400, recorder.Code)

	assert.Equal(suite.T(), "{\"message\":\"strconv.Atoi: parsing \\\"1day\\\": invalid syntax\"}", recorder.Body.String())
}

func (suite *WebTestSuite) TestGetWeatherForecastInvalidResponse() {

	citiesRequestErr := suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/locations/v1/cities/AU/search")).
		WithQueryParam("q", wiremock.EqualTo("Sydney")).
		WithQueryParam("offset", wiremock.EqualTo("1")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnJSON(
			[]map[string]interface{}{{"Key": 123, "EnglishName": "Sydney"}}, map[string]string{"Content-Type": "application/json"}, 200,
		))

	assert.NoError(suite.T(), citiesRequestErr)

	request, err := http.NewRequest("GET", "/weather/AU/Sydney/1", nil)

	assert.NoError(suite.T(), err)

	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), 500, recorder.Code)
}

func (suite *WebTestSuite) BenchmarkGetWeatherForecast(b *testing.B) {

	setupCityStub(suite)

	setUpForecastStub(suite)

	for i := 0; i < b.N; i++ {
		request, err := http.NewRequest("GET", "/weather/AU/Sydney/1", nil)

		assert.NoError(suite.T(), err)

		recorder := httptest.NewRecorder()

		suite.router.ServeHTTP(recorder, request)

		assert.Equal(suite.T(), 200, recorder.Code)

		actualResponse := domain.DailyForecast{}

		assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &actualResponse))

		assert.Equal(suite.T(),
			domain.DailyForecast{MinimumTemp: 9.0, MaximumTemp: 21.0, TempUnit: "C", Sunrise: "2019-05-15T06:44:00+10:00", Sunset: "2019-05-15T17:01:00+10:00"},
			actualResponse)
	}
}

func setUpForecastStub(suite *WebTestSuite) {
	forecastRequestErr := suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/forecasts/v1/daily/1day/123")).
		WithQueryParam("metric", wiremock.EqualTo("true")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnFileContent("upstream_responses/forecast.json", map[string]string{"Content-Type": "application/json"}, 200))

	assert.NoError(suite.T(), forecastRequestErr)
}

func setupCityStub(suite *WebTestSuite) {
	citiesRequestErr := suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/locations/v1/cities/AU/search")).
		WithQueryParam("q", wiremock.EqualTo("Sydney")).
		WithQueryParam("offset", wiremock.EqualTo("1")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnJSON(
			[]map[string]interface{}{{"Key": "123", "EnglishName": "Sydney"}}, map[string]string{"Content-Type": "application/json"}, 200,
		))

	assert.NoError(suite.T(), citiesRequestErr)
}

func TestWebSuite(t *testing.T) {
	suite.Run(t, new(WebTestSuite))
}
