package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"org.example/hello/src/service"
	"strconv"
)

func getWeather(c *gin.Context) {
	weatherService, errors := service.GetWeatherInstance()

	if errors != nil {
		c.JSON(500, gin.H{
			"message": fmt.Errorf("failed to get weather service instance: %v", errors).Error(),
		})
		return
	}

	countryCode := c.Param("countryCode")
	city := c.Param("city")
	daysOfForecast := c.Param("daysOfForecast")

	days, err2 := strconv.Atoi(daysOfForecast)

	if err2 != nil {
		c.JSON(400, gin.H{
			"message": err2.Error(),
		})
		return
	}

	cityForecast, err := weatherService.GetCityForecast(countryCode, city, days)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if cityForecast.IsEmpty() {
		c.JSON(404, gin.H{
			"message": fmt.Sprintf("Forecast for %s@%s not found", city, countryCode),
		})
		return
	}

	c.JSON(200, cityForecast)
}
