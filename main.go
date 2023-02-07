package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const lattitudeRiga float64 = 56.9493977
const longitudeRiga float64 = 24.1051846
const openWeatherUrlFormat string = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s"

type WeatherResponse struct {
	City                 string
	Temperature          string
	TemperatureFeelsLike string
	Lat                  float64
	Long                 float64
}

func main() {
	http.HandleFunc("/weather", weather)
	http.ListenAndServe(":8091", nil)
}

func weather(w http.ResponseWriter, r *http.Request) {
	bytes := readBytesFromOpenWeather(lattitudeRiga, longitudeRiga)
	genericData := bytesToGenericObject(bytes)
	responseWeather := makeWeatherResponse("Riga", genericData)
	responseBytes, _ := json.Marshal(responseWeather)
	fmt.Fprint(w, string(responseBytes))
}

func readBytesFromOpenWeather(lat float64, long float64) []byte {
	key := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf(openWeatherUrlFormat, lat, long, key)
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return body
}

func bytesToGenericObject(bytes []byte) map[string]interface{} {
	var dat map[string]interface{}

	json.Unmarshal(bytes, &dat)

	return dat
}

func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

func makeWeatherResponse(city string, json map[string]interface{}) WeatherResponse {
	main := json["main"].(map[string]interface{})

	return WeatherResponse{
		City:                 city,
		Temperature:          fmt.Sprintf("%.2f", kelvinToCelsius(main["temp"].(float64))),
		TemperatureFeelsLike: fmt.Sprintf("%.2f", kelvinToCelsius(main["feels_like"].(float64))),
		Lat:                  lattitudeRiga,
		Long:                 longitudeRiga,
	}
}
