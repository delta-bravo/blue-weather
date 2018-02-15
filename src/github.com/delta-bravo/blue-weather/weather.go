package main

import (
	"net/http"
	"html/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-ble/ble"
	"golang.org/x/net/context"
	"time"
	"fmt"
	"log"
	"github.com/pkg/errors"
	"strings"
	"math"
	"encoding/binary"
	"github.com/delta-bravo/blue-weather/bleservices"
)

var ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), time.Second*120))

var htmlTemplate *template.Template

var temperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "ambient_temperature",
	Help: "Ambient temperature reading",
})

var ambientTemperature = 0.0

var temperatureHandler = func(data []byte) {
	log.Printf("Current temperature %d\n", data)
	ambientTemperature = math.Float64frombits(binary.LittleEndian.Uint64(data))
	temperatureGauge.Set(ambientTemperature)
}

func main() {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	htmlTemplate = t
	go bleservices.StartReadingMicrobitTemperature(temperatureHandler)
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

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
