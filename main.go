package main

import (
	"encoding/json"
	"github.com/tushar-chauhan/helloGo/weather_lib"
	"log"
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
	service := strings.SplitN(r.URL.Path, "/", 4)[2]
	city := strings.SplitN(r.URL.Path, "/", 4)[3]
	log.Println("Service is " + service + " and city is " + city)
	switch service {
	case "yahoo":
		data, _ := weather_lib.QueryYahooWeather(city)
		// PanicErr(w, err)
		w.Header().Set("Content-Type", "application/json, charset=utf-8")
		json.NewEncoder(w).Encode(data)
	case "open":
		data, err := weather_lib.QueryOpenweathermap(city)
		PanicErr(w, err)
		w.Header().Set("Content-Type", "application/json, charset=utf-8")
		json.NewEncoder(w).Encode(data)
	default:
		log.Println("Default:", "No weather data.")
	}

}

func PanicErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
