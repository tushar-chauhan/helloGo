package weatherutil

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Method called by Go Routine to make HTTP GET request.
func SourceWeathermap(city string) (OutputStruct, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return OutputStruct{}, err
	}

	defer resp.Body.Close()

	var data weatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return OutputStruct{}, err
	}

	if data.Name != "" {
		data.Main.Kelvin = roundFloat((data.Main.Kelvin - 273.15), .5, 2)
	}

	var output OutputStruct
	output.Weather.City = data.Name
	output.Weather.Lat = data.Coord.Lat
	output.Weather.Long = data.Coord.Lon
	output.Weather.Temp = data.Main.Kelvin

	return output, nil
}

// Process multiple cities HTTP calls using Go Routine
func ProcessCities(cities []string) []*OutputStruct {
	structSlice := []*OutputStruct{}

	// Make channels for
	weathers := make(chan *OutputStruct)
	errs := make(chan error)

	for _, city := range cities {
		go func(city string) {
			log.Println("Started processing city: " + city)
			data, err := SourceWeathermap(city)
			if err != nil {
				errs <- err
				return
			}
			weathers <- &data
		}(city)
	}

	for {
		select {
		case weather := <-weathers:
			log.Printf("%s was fetched", weather.Weather.City)
			structSlice = append(structSlice, weather)
			if len(structSlice) == len(cities) {
				return structSlice
			}
		default:
			time.Sleep(2000 * time.Millisecond)
			log.Printf(".")
		}
	}

	return structSlice
}

// Struct to hold the output JSON Object
type OutputStruct struct {
	Weather struct {
		City string  `json:"city"`
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
		Temp float64 `json:"temp"`
	} `json:"weather"`
}
