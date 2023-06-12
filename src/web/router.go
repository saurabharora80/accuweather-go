package web

import "github.com/gin-gonic/gin"

func InitAndConfigureRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/weather/:countryCode/:city/:daysOfForecast", getWeather)

	return router
}
