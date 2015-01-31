package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", sayHello)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World!</h1>"))
}
