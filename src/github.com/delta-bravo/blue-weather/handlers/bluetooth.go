package handlers

import (
	"github.com/delta-bravo/blue-weather/prometheusclient"
	"log"
)

type BluetoothHandlers interface {
	CreateTemperatureHandler() func(data []byte)
	GetCurrentTemperature() float64
}

type GaugeHandlers struct {
	temperatureGauge float64
	prometheusClient prometheusclient.Client
}

func (handlers *GaugeHandlers) CreateTemperatureHandler() func(data []byte) {
	return func(data []byte) {
		currentTemperature := float64(data[0])
		if currentTemperature != handlers.temperatureGauge {
			log.Println("New temperature reading received", currentTemperature, " degrees C")
			handlers.temperatureGauge = currentTemperature
			handlers.prometheusClient.UpdateTemperatureGaugeValue(handlers.temperatureGauge)
		}
	}
}

func (handlers *GaugeHandlers) GetCurrentTemperature() float64 {
	return handlers.temperatureGauge
}

func CreateBluetoothHandlers(prometheusClient prometheusclient.Client) BluetoothHandlers {
	return &GaugeHandlers{
		prometheusClient: prometheusClient,
	}
}
