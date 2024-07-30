package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"github.com/gofor-little/env"
)

type Report struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch    int64 `json:"time_epoch"`
				TempC        float64 `json:"temp_c"`
				Condition    struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	x := "kolkata"

	if len(os.Args) >= 2 {
		x = os.Args[1]
	}
	if err := env.Load(".env"); err != nil {
		panic(err)
	}

	API_KEY, err := env.MustGet("WEATHER_API_KEY")
	if err != nil {
		panic(err)
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key="+API_KEY+"&q="+x+"&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var report Report

	err = json.Unmarshal(body, &report)
	if err !=nil {
		panic(err)
	}

	location, current, hours := report.Location, report.Current, report.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s, %.0fC, %s\n", location.Name,location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours{
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}

		fmt.Printf("%s - %.0fC, %.0f, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
	}
}