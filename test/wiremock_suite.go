package test

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/walkerus/go-wiremock"
	"os"
	"path/filepath"
)

type WiremockTestSuite struct {
	suite.Suite
	AccuweatherBaseUrl string
	Container          testcontainers.Container
	WiremockClient     *wiremock.Client
	HttpClient         *resty.Client
}

func (suite *WiremockTestSuite) SetupSuite() {
	ctx := context.Background()

	dir, _ := os.Getwd()
	testDirectory := filepath.Dir(dir)

	containerRequest := testcontainers.ContainerRequest{
		Image:        "wiremock/wiremock",
		ExposedPorts: []string{"8080"},
		WaitingFor:   wait.ForHTTP("/__admin/mappings").WithPort("8080").WithMethod("GET"),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = map[nat.Port][]nat.PortBinding{"8080": {{HostIP: "0.0.0.0", HostPort: "8080"}}}
		},
	}

	wiremockContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true})

	assert.NoError(suite.T(), err)

	err = wiremockContainer.CopyDirToContainer(ctx, fmt.Sprintf("%s/resources/upstream_responses", testDirectory), "/home/wiremock/__files/", 700)

	assert.NoError(suite.T(), err)

	wiremockPort, _ := wiremockContainer.MappedPort(ctx, "8080")

	fmt.Println("Wiremock running on port", wiremockPort.Port())

	suite.Container = wiremockContainer

	suite.AccuweatherBaseUrl = fmt.Sprintf("http://localhost:%s", wiremockPort.Port())

	suite.WiremockClient = wiremock.NewClient(suite.AccuweatherBaseUrl)

	suite.HttpClient = resty.New().
		EnableTrace().
		SetQueryParam("apikey", "test-api-key").
		SetHeader("Accept", "application/json").
		SetBaseURL(suite.AccuweatherBaseUrl)

}

func (suite *WiremockTestSuite) TearDownSuite() {
	assert.NoError(suite.T(), suite.Container.Terminate(context.Background()))
}

func (suite *WiremockTestSuite) SetupTest() {
	assert.NoError(suite.T(), suite.WiremockClient.Reset())
}
