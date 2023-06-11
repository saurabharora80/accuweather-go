package upstream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/walkerus/go-wiremock"
	"org.example/hello/src/upstream"
	"testing"
)

func (suite *CitiesTestSuite) TestGetCityInCountry() {

	err := suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/locations/v1/cities/AU/search")).
		WithQueryParam("q", wiremock.EqualTo("sydney")).
		WithQueryParam("offset", wiremock.EqualTo("1")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnJSON(
			[]map[string]interface{}{{"Key": "123", "EnglishName": "Sydney"}}, map[string]string{"Content-Type": "application/json"}, 200,
		))

	if err != nil {
		assert.Fail(suite.T(), "Failed to configure wiremock stub due to %s", err.Error())
		return
	}

	cities := make(chan upstream.City)
	errors := make(chan error)

	go upstream.GetCityInCountry(suite.HttpClient, "AU", "sydney", cities, errors)

	select {
	case city := <-cities:
		assert.Equal(suite.T(), upstream.City{Key: "123", Name: "Sydney"}, city)
	case err := <-errors:
		assert.Fail(suite.T(), "Unable to get City because of %s", err.Error())
	}
}

func (suite *CitiesTestSuite) TestGetCityInCountryNotFound() {
	cities := make(chan upstream.City)
	errors := make(chan error)

	go upstream.GetCityInCountry(suite.HttpClient, "AU", "melbourne", cities, errors)

	select {
	case city := <-cities:
		assert.Equal(suite.T(), upstream.City{}, city)
	case err := <-errors:
		assert.Fail(suite.T(), "Unable to get City because of %s", err.Error())
	}
}

func (suite *CitiesTestSuite) TestGetCityInCountryFailed() {
	invalidStatusCodes := []int64{403, 500, 501, 503}

	for _, invalidStatusCode := range invalidStatusCodes {
		_ = suite.WiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/locations/v1/cities/AU/search")).
			WithQueryParam("q", wiremock.EqualTo("melbourne")).
			WithQueryParam("offset", wiremock.EqualTo("1")).
			WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
			WillReturnJSON(
				[]map[string]interface{}{{}}, map[string]string{"Content-Type": "application/json"}, invalidStatusCode,
			))

		cities := make(chan upstream.City)
		errors := make(chan error)

		go upstream.GetCityInCountry(suite.HttpClient, "AU", "melbourne", cities, errors)

		select {
		case _ = <-cities:
			assert.Fail(suite.T(), "Shouldn't return city")
		case err := <-errors:
			assert.Equal(suite.T(), fmt.Sprintf("GET /locations/v1/cities/AU/search?q=melbourne failed with %d => [{}]", invalidStatusCode), err.Error())
		}
	}
}

func TestCalculatorTestSuite(t *testing.T) {
	suite.Run(t, new(CitiesTestSuite))
}
