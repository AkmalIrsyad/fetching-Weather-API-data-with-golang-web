package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
)

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Timestamp int64 `json:"dt"`
	Date      string
}

func getWeatherData(city string, apiKey string) (WeatherData, error) {
	var weatherData WeatherData
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return weatherData, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return weatherData, err
	}

	// Membulatkan suhu ke bilangan bulat
	weatherData.Main.Temp = math.Round(weatherData.Main.Temp)

	// Konversi Timestamp UNIX menjadi format tanggal
	weatherData.Date = time.Unix(weatherData.Timestamp, 0).Format("02 January 2006")

	return weatherData, nil
}

func renderTemplate(w http.ResponseWriter, data WeatherData) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// Load API key from .env
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not found", http.StatusInternalServerError)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		city = "jakarta"
	}

	weatherData, err := getWeatherData(city, apiKey)
	if err != nil {
		http.Error(w, "Unable to get weather data", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, weatherData)
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", weatherHandler)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
