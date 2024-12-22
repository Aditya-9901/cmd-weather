package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type City struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type WeatherDetail struct {
	Temperature float64 `json:"temp"`
	Pressure    int64   `json:"pressure"`
}

type WeatherResponse struct {
	Main WeatherDetail `json:"main"`
}

func FindCityByName(city string) int {
	by, _ := ioutil.ReadFile("cityList.json")
	var cities []City
	json.Unmarshal(by, &cities)
	for _, c := range cities {
		if strings.ToLower(c.Name) == city {
			return c.Id
		}
	}
	return -1
}

func main() {

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("Fail to get the api key from the env")
	}

	cityName := flag.String("city", "", "--city \"<city-name>\" to get the weather")
	flag.Parse()

	cityId := FindCityByName(strings.ToLower(*cityName))
	if cityId == -1 {
		fmt.Println("Failed to find the city id")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s", cityId, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The current temperature in %v is %.2f ∘C and the pressure is %v Pa", string(*cityName), weather.Main.Temperature-273.15, weather.Main.Pressure*1000)
}
