package upstream

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"org.example/hello/src/config"
	"os"
	"strconv"
	"time"
)

func NewRestyClient() *resty.Client {
	c, err := config.GetConfig()

	ForecastConnectorInstanceError = err

	parseBool, boolParseErr := strconv.ParseBool(os.Getenv("ENABLE.RESTY.DEBUG"))

	if boolParseErr != nil {
		parseBool = false
	}

	return resty.New().
		EnableTrace().
		SetDebug(parseBool).
		SetTransport(&http.Transport{
			MaxIdleConns:    c.Upstream.MaxIdleConnections,
			IdleConnTimeout: c.Upstream.IdleConnectionTimeoutSeconds * time.Second}).
		SetQueryParam("apikey", c.Upstream.Key).
		SetHeader("Accept", "application/json").
		SetBaseURL(c.Upstream.Host)
}
