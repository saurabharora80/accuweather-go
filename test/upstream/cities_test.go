package upstream

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/walkerus/go-wiremock"
	"org.example/hello/src/upstream"
	"testing"
)

func setUpSuite(t *testing.T) (func(t *testing.T), string) {
	ctx := context.Background()

	containerRequest := testcontainers.ContainerRequest{
		Image:        "wiremock/wiremock",
		ExposedPorts: []string{"8080"},
		WaitingFor:   wait.ForHTTP("/__admin/mappings").WithPort("8080").WithMethod("GET")}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{ContainerRequest: containerRequest, Started: true})

	if err != nil {
		t.Fatalf("Failed to start container due to %q", err.Error())
	}

	wiremockPort, err := container.MappedPort(ctx, "8080")

	return func(t *testing.T) {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Faild to terminate container: %s", err.Error())
		}
	}, wiremockPort.Port()
}

func TestGetCityInCountry(t *testing.T) {
	tearDownSuite, wiremockPort := setUpSuite(t)
	defer tearDownSuite(t)

	wiremockUrl := fmt.Sprintf("http://localhost:%s", wiremockPort)

	wiremockClient := wiremock.NewClient(wiremockUrl)

	defer func(wiremockClient *wiremock.Client) {
		err := wiremockClient.Reset()
		if err != nil {
			t.Errorf("Unable to reset wiremock %q", err.Error())
		}
	}(wiremockClient)

	getCityStubRule := wiremock.Get(wiremock.
		URLPathEqualTo("/locations/v1/cities/AU/search")).
		WithQueryParam("q", wiremock.EqualTo("sydney")).
		WithQueryParam("offset", wiremock.EqualTo("1")).
		WithQueryParam("apikey", wiremock.EqualTo("test-api-key")).
		WillReturnJSON(
			[]map[string]interface{}{{"Key": "123", "Name": "Sydney"}}, map[string]string{"Content-Type": "application/json"}, 200,
		)

	if err := wiremockClient.StubFor(getCityStubRule.AtPriority(1)); err != nil {
		t.Errorf("Failed to configure wiremock stub due to %q", err)
		return
	}

	client := resty.New().
		EnableTrace().
		SetQueryParam("apikey", "test-api-key").
		SetHeader("Accept", "application/json").
		SetBaseURL(wiremockUrl)

	country, err := upstream.GetCityInCountry(client, "AU", "sydney")

	if err != nil {
		t.Errorf("Get Countries Failed with error %s", err.Error())
	} else if country.Key != "123" && country.Name != "Sydney" {
		t.Errorf("Expected country with Key 123 instead Got (%s,%s)", country.Key, country.Name)
	}
}
