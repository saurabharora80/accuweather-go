package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"org.example/hello/src/web"
	"org.example/hello/test/common"
	"os"
)

type WebTestSuite struct {
	common.WiremockTestSuite
	router *gin.Engine
}

func (suite *WebTestSuite) SetupTest() {
	/*
		Set the environment variable upstream.host to the value of the Wiremock server
		to allow the application upstreams to connect to the Wiremock server.
	*/
	assert.NoError(suite.T(), os.Setenv("UPSTREAM.HOST", suite.AccuweatherBaseUrl))

	suite.router = web.InitAndConfigureRouter()

	assert.NoError(suite.T(), suite.WiremockClient.Reset())
}
