package main

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", sayHello)

	http.HandleFunc("/forecast/", getForecast)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html, charset=utf-8")
	w.Write([]byte("<h3>Hello World!</h3>"))
}

func getForecast(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var data weatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherData{}, err
	}

	if data.Name != "" {
		data.Main.Kelvin = roundFloat((data.Main.Kelvin - 273.15), .5, 2)
		data.Message = "success"
	}

	return data, nil
}

func roundFloat(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

type weatherData struct {
	Name  string `json:"name"`
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Message string `json:"message"`
}
