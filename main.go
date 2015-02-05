package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tushar-chauhan/helloGo/weatherutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", sayHello)

	router.HandleFunc("/forecast/{service}/{city}", getForecast)
	router.HandleFunc("/weather/{cities}", getForecast)
	http.Handle("/", router)
	log.Println("Server listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html, charset=utf-8")
	w.Write([]byte("<h3>Hello World!</h3>"))
}

func getForecast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	city := vars["city"]
	var cities []string = strings.Split(vars["cities"], ",")

	if service != "" {
		log.Printf("Received single request for Service: '%s' and City: '%s'", service, city)
		switch service {
		case "yahoo":
			data, err := weatherutil.QueryYahooWeather(city)
			PanicErr(w, err)
			w.Header().Set("Content-Type", "application/json, charset=utf-8")
			json.NewEncoder(w).Encode(data)
		case "open":
			data, err := weatherutil.QueryOpenweathermap(city)
			PanicErr(w, err)
			w.Header().Set("Content-Type", "application/json, charset=utf-8")
			json.NewEncoder(w).Encode(data)
		default:
			log.Println("Default:", "No weather Service opted..")
		}
	}

	if cities != nil && service == "" {
		log.Printf("Received single request for Multiple cities: %s", vars["cities"])
		structSlice := weatherutil.ProcessCities(cities)
		w.Header().Set("Content-Type", "application/json, charset=utf-8")
		json.NewEncoder(w).Encode(structSlice)
	}

}

func PanicErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
