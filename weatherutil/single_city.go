package weatherutil

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
)

// Method to make API call to Openweathermap
func QueryOpenweathermap(city string) (weatherData, error) {
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

//Method to make call to Yahoo Weather API using YQL
func QueryYahooWeather(city string) (yahooData, error) {

	u, err := url.Parse("https://query.yahooapis.com/v1/public/yql")
	q := u.Query()
	q.Set("q", "select location.city, item.lat, item.long, item.condition.temp from weather.forecast where woeid in (select woeid from geo.places(1) where text='"+city+"') and u='c'")
	q.Set("format", "json")
	q.Set("appid", "yIIX1S4o")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		return yahooData{}, err
	}

	defer resp.Body.Close()

	var data yahooData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return yahooData{}, err
	}

	return data, nil
}

// Function to round of the floating point values
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

// Struct to store the required Openweathermap JSON response
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

// Struct to store the required Yahoo Weather YQL JSON response
type yahooData struct {
	Query struct {
		Result struct {
			Channel struct {
				Location struct {
					City string `json:"city"`
				} `json:"location"`
				Item struct {
					Lat       string `json:"lat"`
					Lon       string `json:"long"`
					Condition struct {
						Temp string `json:"temp"`
					} `json:"condition"`
				} `json:"item"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
}
