package upstream

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/walkerus/go-wiremock"
)

type CitiesTestSuite struct {
	suite.Suite
	AccuweatherBaseUrl string
	Container          testcontainers.Container
	WiremockClient     *wiremock.Client
	HttpClient         *resty.Client
}

func (suite *CitiesTestSuite) SetupSuite() {
	ctx := context.Background()

	containerRequest := testcontainers.ContainerRequest{
		Image:        "wiremock/wiremock",
		ExposedPorts: []string{"8080"},
		WaitingFor:   wait.ForHTTP("/__admin/mappings").WithPort("8080").WithMethod("GET")}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{ContainerRequest: containerRequest, Started: true})

	if err != nil {
		suite.T().Fatalf("Failed to start container due to %q", err.Error())
	}

	wiremockPort, _ := container.MappedPort(ctx, "8080")

	suite.Container = container

	suite.AccuweatherBaseUrl = fmt.Sprintf("http://localhost:%s", wiremockPort.Port())

	suite.WiremockClient = wiremock.NewClient(suite.AccuweatherBaseUrl)

	suite.HttpClient = resty.New().
		EnableTrace().
		SetQueryParam("apikey", "test-api-key").
		SetHeader("Accept", "application/json").
		SetBaseURL(suite.AccuweatherBaseUrl)

}

func (suite *CitiesTestSuite) TearDownSuite() {
	if err := suite.Container.Terminate(context.Background()); err != nil {
		suite.T().Fatalf("Faild to terminate container: %s", err.Error())
	}
}

func (suite *CitiesTestSuite) SetupTest() {
	err := suite.WiremockClient.Reset()
	if err != nil {
		suite.T().Errorf("Unable to reset wiremock %q", err.Error())
	}
}
