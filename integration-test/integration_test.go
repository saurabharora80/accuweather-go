package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"org.example/hello/src/domain"
	"org.example/hello/src/web"
	"os"
	"testing"
)

func TestGetForecastFromAccuweather(t *testing.T) {
	pwd, _ := os.Getwd()

	_ = os.Setenv("CONFIG_PATH", fmt.Sprintf("%s/../src/config", pwd))
	_ = os.Setenv("ENABLE.RESTY.DEBUG", "true")

	router := web.InitAndConfigureRouter()

	request, err := http.NewRequest("GET", "/weather/AU/Sydney/1", nil)

	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Printf("==================BODY=================== %s\n", recorder.Body.String())

	assert.Equal(t, 200, recorder.Code)

	actualResponse := domain.DailyForecast{}

	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &actualResponse))

	assert.Condition(t,
		func() bool {
			return !actualResponse.IsEmpty()
		},
		actualResponse)

}
