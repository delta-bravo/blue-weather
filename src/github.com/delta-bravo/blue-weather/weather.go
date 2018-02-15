package main

import (
	"net/http"
	"html/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"golang.org/x/net/context"
	"time"
	"fmt"
	"log"
	"github.com/pkg/errors"
	"strings"
	"math"
	"encoding/binary"
)

const temperatureServiceUuid = "E95D6100251D470AA062FA1922DFA9A8"
const temperatureCharacteristicUuid = "E95D9250251D470AA062FA1922DFA9A8"

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
	go startReadingMicrobitTemperature()
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

func filter(a ble.Advertisement) bool {
	return strings.Contains(a.LocalName(), "micro:bit");
}

func startReadingMicrobitTemperature() {
	device, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("Unable to init device: %s\n", err)
		return
	}
	ble.SetDefaultDevice(device)

	log.Println("Connecting...")

	client, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatalf("Unable to connect: %s\n", err)
		return
	}

	profile, err := client.DiscoverProfile(true)

	if err != nil {
		log.Fatalf("Unable to discover device profile: %s\n", err)
		return
	}

	discoverAndSubscribeTemperatureService(profile, client)
}

func discoverAndSubscribeTemperatureService(profile *ble.Profile, client ble.Client) {
	for _, service := range profile.Services {
		serviceUuid, _ := ble.Parse(temperatureServiceUuid)
		if service.UUID.Equal(serviceUuid) {
			log.Println("Temperature service discovered service, yay")
			subscribeForTemperatureUpdates(service, client)
		}
	}
}

func subscribeForTemperatureUpdates(service *ble.Service, client ble.Client) {
	for _, characteristic := range service.Characteristics {
		characteristicUuid, _ := ble.Parse(temperatureCharacteristicUuid)
		if characteristic.UUID.Equal(characteristicUuid) {
			log.Println("Temperature characteristic discovered!")
			client.Subscribe(characteristic, true, temperatureHandler)
		}
	}
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
