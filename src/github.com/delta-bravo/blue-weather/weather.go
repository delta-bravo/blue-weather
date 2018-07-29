package main

import (
	"net/http"
	"html/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
	"log"
	"github.com/delta-bravo/blue-weather/bleservices"
)

var htmlTemplate *template.Template

var temperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "ambient_temperature",
	Help: "Ambient temperature reading",
})

var ambientTemperature = 0.0

var temperatureHandler = func(data []byte) {
	oldReading := ambientTemperature
	ambientTemperature = float64(data[0])
	temperatureGauge.Set(ambientTemperature)
	if oldReading != ambientTemperature {
		log.Printf("Got new temperature reading: %f", ambientTemperature)
	}

}

func main() {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	htmlTemplate = t
	bluetoothClient := bleservices.GetClient()

	go bleservices.StartBluetoothServices(bluetoothClient, temperatureHandler)

	prometheus.Register(temperatureGauge)
	addHttpHandlers()
}

func addHttpHandlers() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/temperature", refreshHandler)
	http.HandleFunc("/style.css", createFileHandler("style.css"))
	http.HandleFunc("/refresh.js", createFileHandler("refresh.js"))
	err := http.ListenAndServe(":8015", nil)
	if err != nil {
		log.Fatal("Unable to open port for listening")
		panic(err)
	}
}

func createFileHandler(fileName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fileName)
	}
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, ambientTemperature)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	htmlTemplate.Execute(w, map[string]interface{}{"temperature": ambientTemperature})
}
