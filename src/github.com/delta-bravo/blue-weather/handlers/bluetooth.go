package handlers

import (
	"github.com/delta-bravo/blue-weather/prometheusclient"
	"log"
	"math"
)

type BluetoothHandlers interface {
	CreateTemperatureHandler() func(data []byte)
	CreateBearingHandler() func(data []byte)
	GetCurrentTemperature() float64
	GetCurrentBearing() float64
}

type GaugeHandlers struct {
	temperatureGauge float64
	bearingGauge     float64
	prometheusClient prometheusclient.Client
}

func (handlers *GaugeHandlers) CreateTemperatureHandler() func(data []byte) {
	return func(data []byte) {
		currentTemperature := float64(data[0])
		if currentTemperature != handlers.temperatureGauge {
			log.Println("New temperature reading received", currentTemperature, "degrees C")
			handlers.temperatureGauge = currentTemperature
			handlers.prometheusClient.UpdateTemperatureGaugeValue(handlers.temperatureGauge)
		}
	}
}

func (handlers *GaugeHandlers) CreateBearingHandler() func(data []byte) {
	return func(data []byte) {
		currentBearing := float64(data[0])
		if math.Abs(currentBearing-handlers.bearingGauge) > 2 {
			log.Println("New bearing reading received", currentBearing, "degrees")
			handlers.bearingGauge = currentBearing
			handlers.prometheusClient.UpdateBearing(handlers.bearingGauge)
		}
	}
}

func (handlers *GaugeHandlers) GetCurrentTemperature() float64 {
	return handlers.temperatureGauge
}

func (handlers *GaugeHandlers) GetCurrentBearing() float64 {
	return handlers.bearingGauge
}

func CreateBluetoothHandlers(prometheusClient prometheusclient.Client) BluetoothHandlers {
	return &GaugeHandlers{
		prometheusClient: prometheusClient,
	}
}
